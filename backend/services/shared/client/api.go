package client

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/client/model"
	resty "github.com/go-resty/resty/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type ZoomApi interface {
	GetMe(ctx context.Context, token string) (model.User, error)
	GetFilesFromMessages(ctx context.Context, token string, params map[string]string) (model.ZoomFileMessage, error)
	GetFileFromMessage(ctx context.Context, token, uid, mid, fid string) (model.ZoomFile, error)
}

type zoomApiClient struct {
	client *resty.Client
}

func NewZoomApiClient() ZoomApi {
	otelClient := otelhttp.DefaultClient
	otelClient.Transport = otelhttp.NewTransport(&http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   15 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})
	return zoomApiClient{
		client: resty.NewWithClient(otelClient).
			SetHostURL("https://api.zoom.us").
			SetRetryCount(3).
			SetTimeout(3 * time.Second).
			SetRetryWaitTime(120 * time.Millisecond).
			SetRetryMaxWaitTime(900 * time.Millisecond).
			SetLogger(log.NewEmptyLogger()).
			AddRetryCondition(func(r *resty.Response, err error) bool {
				return r.StatusCode() == http.StatusTooManyRequests
			}),
	}
}

func (c zoomApiClient) GetMe(ctx context.Context, token string) (model.User, error) {
	var resp model.User
	res, err := c.client.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetResult(&resp).
		Get("/v2/users/me")

	if err != nil {
		return resp, err
	}

	if res.StatusCode() != http.StatusOK {
		return resp, &UnexpectedStatusCodeError{
			Action: "get me",
			Code:   res.StatusCode(),
		}
	}

	return resp, nil
}

func (c zoomApiClient) GetFilesFromMessages(ctx context.Context, token string, params map[string]string) (model.ZoomFileMessage, error) {
	var resp model.ZoomFileMessage
	if val, ok := params["to_contact"]; !ok || val == "" {
		return resp, &InvalidQueryParameterError{
			Parameter: "to_contact",
		}
	}

	params["from"] = "2000-02-10T21:39:50Z"
	params["search_type"] = "file"
	if _, ok := params["search_key"]; !ok {
		params["search_key"] = " "
	}

	res, err := c.client.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetPathParam("to_contact", params["to_contact"]).
		SetQueryParams(params).
		SetResult(&resp).
		Get("/v2/chat/users/{to_contact}/messages")

	if err != nil {
		return resp, err
	}

	if res.StatusCode() != http.StatusOK {
		return resp, &UnexpectedStatusCodeError{
			Action: "get files from message",
			Code:   res.StatusCode(),
		}
	}

	return resp, nil
}

func (c zoomApiClient) GetFileFromMessage(ctx context.Context, token, uid, mid, fid string) (model.ZoomFile, error) {
	var empty model.ZoomFile
	var resp model.ZoomMessage

	uid = strings.TrimSpace(uid)
	if uid == "" {
		return empty, &InvalidQueryParameterError{
			Parameter: "userID",
		}
	}

	mid = strings.TrimSpace(mid)
	if mid == "" {
		return empty, &InvalidQueryParameterError{
			Parameter: "messageID",
		}
	}

	fid = strings.TrimSpace(fid)
	if fid == "" {
		return empty, &InvalidQueryParameterError{
			Parameter: "fileID",
		}
	}

	_, err := c.client.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetPathParam("user", uid).
		SetPathParam("message", mid).
		SetQueryParams(map[string]string{
			"to_contact": uid,
		}).
		SetResult(&resp).
		Get("/v2/chat/users/{user}/messages/{message}")

	if err != nil {
		return empty, err
	}

	for _, f := range resp.Files {
		if f.FileID == fid {
			return f, nil
		}
	}

	return empty, ErrFileDoesNotExist
}
