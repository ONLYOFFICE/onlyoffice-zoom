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
	"sync"
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
	}
}

func emptyMultipartSize(fieldname, filename string) int64 {
	body := &bytes.Buffer{}
	form := multipart.NewWriter(body)
	form.CreateFormFile(fieldname, filename)
	form.Close()
	return int64(body.Len())
}

func doRequest(ctx context.Context, address, method string, body io.Reader, contentType string, contentLength int64, desiredStatus int, token string) (*http.Response, error) {
	targetURL, err := url.Parse(address)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, targetURL.String(), body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}

	if contentLength > 0 {
		request.ContentLength = contentLength
	}

	response, err := otelhttp.DefaultClient.Do(request.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	if response.StatusCode != desiredStatus {
		response.Body.Close()
		return nil, fmt.Errorf("unexpected '%s' from: %s %s", response.Status, method, targetURL.String())
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

	return fileResp.RawBody(), nil
}

func (c zoomFilestoreClient) UploadFile(ctx context.Context, url, token, uid, filename string, size int64) error {
	token = strings.TrimSpace(token)
	uid = strings.TrimSpace(uid)
	if token == "" || uid == "" {
		return _ErrInvalidClientToken
	}

	var wg sync.WaitGroup
	fileChan := make(chan io.ReadCloser)
	urlChan := make(chan string)
	errorsChan := make(chan error, 2)

	go func() {
		wg.Add(1)
		defer wg.Done()
		file, err := c.getFile(ctx, url)
		if err != nil {
			errorsChan <- err
			return
		}

		fileChan <- file
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		res, err := c.client.R().
			SetContext(ctx).
			SetAuthToken(token).
			SetFileReader("file", filename, bytes.NewReader([]byte{})).
			SetPathParam("user", uid).
			Post(fmt.Sprintf("%s/chat/users/{user}/files", constants.ZOOM_FILE_API_HOST))

		if res.StatusCode() != 307 {
			errorsChan <- err
			return
		}

		urlChan <- res.Header().Get("Location")
	}()

	wg.Wait()

	select {
	case err := <-errorsChan:
		return err
	default:
	}

	file, url := <-fileChan, <-urlChan
	defer file.Close()

	contentLength := emptyMultipartSize("file", filename) + size
	readBody, writeBody := io.Pipe()
	defer readBody.Close()
	form := multipart.NewWriter(writeBody)
	errChan := make(chan error, 1)

	go func() {
		defer writeBody.Close()

		part, err := form.CreateFormFile("file", filename)
		if err != nil {
			errChan <- err
			return
		}

		if _, err := io.CopyN(part, file, size); err != nil {
			errChan <- err
			return
		}

		errChan <- form.Close()
	}()

	if _, err := doRequest(ctx, url, "POST", readBody, form.FormDataContentType(), contentLength, http.StatusCreated, token); err != nil {
		<-errChan
		return err
	}

	return <-errChan
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
