package client

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/constants"
	resty "github.com/go-resty/resty/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var _ErrInvalidClientToken = errors.New("could not perform zoom filestore action due to invalid access token")

type ZoomFilestore interface {
	UploadFile(ctx context.Context, token, uid, filename string, file io.Reader) error
	GetFile(ctx context.Context, url string) (io.Reader, error)
	ValidateFileSize(ctx context.Context, limit int64, url string) error
}

type zoomFilestoreClient struct {
	client *resty.Client
}

func NewZoomFilestoreClient() ZoomFilestore {
	otelClient := otelhttp.DefaultClient
	otelClient.Transport = otelhttp.NewTransport(&http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   15 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})

	client := resty.NewWithClient(otelClient).
		SetHostURL(constants.ZOOM_FILE_API_HOST).
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(20)).
		SetRetryCount(0).
		SetRetryWaitTime(100 * time.Millisecond).
		SetRetryMaxWaitTime(700 * time.Millisecond).
		SetLogger(log.NewEmptyLogger()).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			return false
		})
	client.GetClient().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		req.Header.Set("Authorization", via[0].Header.Get("Authorization"))
		return nil
	}

	return zoomFilestoreClient{
		client: client,
	}
}

func (c zoomFilestoreClient) UploadFile(ctx context.Context, token, uid, filename string, file io.Reader) error {
	token = strings.TrimSpace(token)
	uid = strings.TrimSpace(uid)
	if token == "" || uid == "" {
		return _ErrInvalidClientToken
	}

	res, err := c.client.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetFileReader("file", filename, file).
		SetPathParam("user", uid).
		Post("/chat/users/{user}/files")

	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusCreated {
		return &UnexpectedStatusCodeError{
			Action: "upload file",
			Code:   res.StatusCode(),
		}
	}

	return nil
}

func (c zoomFilestoreClient) GetFile(ctx context.Context, url string) (io.Reader, error) {
	fileResp, err := c.client.R().
		SetContext(ctx).
		Get(url)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(fileResp.Body()), nil
}

func (c zoomFilestoreClient) ValidateFileSize(ctx context.Context, limit int64, url string) error {
	headResp, err := c.client.R().
		SetContext(ctx).
		Head(url)

	if err != nil {
		return err
	}

	if val, err := strconv.ParseInt(headResp.Header().Get("Content-Length"), 10, 0); val > limit || err != nil {
		return _ErrInvalidClientToken
	}

	return nil
}
