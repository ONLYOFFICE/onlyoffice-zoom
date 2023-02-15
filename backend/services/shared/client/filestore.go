package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/constants"
	resty "github.com/go-resty/resty/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var _ErrInvalidClientToken = errors.New("could not perform zoom filestore action due to invalid access token")
var _ErrInvalidContentLength = errors.New("could not perform zoom filestore actions due to exceeding content-length")

type ZoomFilestore interface {
	UploadFile(ctx context.Context, url, token, uid, filename string, size int64) error
	ValidateFileSize(ctx context.Context, limit int64, url string) error
}

type zoomFilestoreClient struct {
	client *resty.Client
	logger log.Logger
}

func NewZoomFilestoreClient(logger log.Logger) ZoomFilestore {
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
		SetRedirectPolicy(resty.NoRedirectPolicy()).
		SetRetryCount(0).
		SetRetryWaitTime(100 * time.Millisecond).
		SetRetryMaxWaitTime(700 * time.Millisecond).
		SetLogger(log.NewEmptyLogger()).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			return false
		})

	return zoomFilestoreClient{
		client: client,
		logger: logger,
	}
}

func emptyMultipartSize(fieldname, filename string) int64 {
	body := &bytes.Buffer{}
	form := multipart.NewWriter(body)
	form.CreateFormFile(fieldname, filename)
	form.Close()
	return int64(body.Len())
}

func (c zoomFilestoreClient) doRequest(ctx context.Context, address, method string, body io.Reader, contentType string, contentLength int64, desiredStatus int, token string) (*http.Response, error) {
	targetURL, err := url.Parse(address)
	if err != nil {
		c.logger.Error(err.Error())
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, method, targetURL.String(), body)
	if err != nil {
		c.logger.Error(err.Error())
		return nil, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}

	if contentLength > 0 {
		request.ContentLength = contentLength
	}

	c.logger.Debugf("POST content-length: %d", request.ContentLength)
	response, err := otelhttp.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != desiredStatus {
		response.Body.Close()
		c.logger.Warnf("unexpected '%s' from: %s %s", response.Status, method, targetURL.String())
		return response, fmt.Errorf("unexpected '%s' from: %s %s", response.Status, method, targetURL.String())
	}

	return response, nil
}

func (c zoomFilestoreClient) getFile(ctx context.Context, url string) (io.ReadCloser, error) {
	fileResp, err := c.client.R().
		SetContext(ctx).
		SetDoNotParseResponse(true).
		Get(url)
	if err != nil {
		return nil, err
	}

	c.logger.Debugf("got a file response form document server with length %s", fileResp.Header().Get("Content-Length"))
	return fileResp.RawBody(), nil
}

func (c zoomFilestoreClient) UploadFile(ctx context.Context, url, token, uid, filename string, size int64) error {
	token = strings.TrimSpace(token)
	uid = strings.TrimSpace(uid)
	if token == "" || uid == "" {
		return _ErrInvalidClientToken
	}

	c.logger.Debugf("got an upload job with token %s and uid %s", token, uid)

	contentLength := emptyMultipartSize("file", filename) + size
	fReader, fWriter := io.Pipe()
	sReader, sWriter := io.Pipe()
	formOne := multipart.NewWriter(fWriter)
	formTwo := multipart.NewWriter(sWriter)
	defer fReader.Close()
	defer sReader.Close()

	go func() {
		defer fWriter.Close()
		file, err := c.getFile(ctx, url)
		if err != nil {
			return
		} else if file == nil {
			return
		}
		defer file.Close()

		part, err := formOne.CreateFormFile("file", filename)
		if err != nil {
			return
		}

		if _, err := io.CopyN(part, file, size); err != nil {
			c.logger.Errorf("could not pipe data to writer: %s", err.Error())
			return
		}

		formOne.Close()
	}()

	go func() {
		defer sWriter.Close()
		file, err := c.getFile(ctx, url)
		if err != nil {
			return
		} else if file == nil {
			return
		}
		defer file.Close()

		part, err := formTwo.CreateFormFile("file", filename)
		if err != nil {
			return
		}

		if _, err := io.CopyN(part, file, size); err != nil {
			c.logger.Errorf("could not pipe data to writer: %s", err.Error())
			return
		}

		formTwo.Close()
	}()

	if resp, err := c.doRequest(ctx, fmt.Sprintf("%s/chat/users/%s/files", constants.ZOOM_FILE_API_HOST, uid), "POST", fReader, formOne.FormDataContentType(), contentLength, http.StatusCreated, token); err != nil {
		if resp != nil && resp.Header.Get("Location") != "" {
			if _, err = c.doRequest(ctx, resp.Header.Get("Location"), "POST", sReader, formTwo.FormDataContentType(), contentLength, http.StatusCreated, token); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	return nil
}

func (c zoomFilestoreClient) ValidateFileSize(ctx context.Context, limit int64, url string) error {
	headResp, err := c.client.R().
		SetContext(ctx).
		Head(url)

	if err != nil {
		return err
	}

	if val, err := strconv.ParseInt(headResp.Header().Get("Content-Length"), 10, 0); val > limit || err != nil {
		return _ErrInvalidContentLength
	}

	return nil
}
