package client

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/client/model"
	resty "github.com/go-resty/resty/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type ZoomAuth interface {
	RefreshZoomAccessToken(ctx context.Context, refreshToken string) (model.Token, error)
	GetZoomAccessToken(ctx context.Context, code, verifier, redirect string) (model.Token, error)
	GetZoomUser(ctx context.Context, token string) (model.User, error)
	GetDeeplink(ctx context.Context, token model.Token) (string, error)
}

type zoomAuthClient struct {
	client       *resty.Client
	clientID     string
	clientSecret string
}

func NewZoomClient(clientID, clientSecret string) ZoomAuth {
	otelClient := otelhttp.DefaultClient
	otelClient.Transport = otelhttp.NewTransport(&http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   30 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})
	return zoomAuthClient{
		client: resty.NewWithClient(otelClient).
			SetHostURL("https://zoom.us").
			SetRetryCount(3).
			SetRetryWaitTime(120 * time.Millisecond).
			SetRetryMaxWaitTime(900 * time.Millisecond).
			SetLogger(log.NewEmptyLogger()).
			AddRetryCondition(func(r *resty.Response, err error) bool {
				return r.StatusCode() == http.StatusTooManyRequests
			}),
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (c zoomAuthClient) RefreshZoomAccessToken(ctx context.Context, refreshToken string) (model.Token, error) {
	var resp model.Token

	res, err := c.client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetContext(ctx).
		SetBody(strings.NewReader(url.Values{
			"client_id":     []string{c.clientID},
			"client_secret": []string{c.clientSecret},
			"grant_type":    []string{"refresh_token"},
			"refresh_token": []string{refreshToken},
		}.Encode())).
		SetResult(&resp).
		Post("/oauth/token")

	if err != nil {
		return resp, err
	}

	if res.StatusCode() != http.StatusOK {
		return resp, &UnexpectedStatusCodeError{
			Action: "refresh access token",
			Code:   res.StatusCode(),
		}
	}

	return resp, resp.Validate()
}

func (c zoomAuthClient) GetZoomAccessToken(ctx context.Context, code, verifier, redirect string) (model.Token, error) {
	var resp model.Token
	if _, err := url.ParseRequestURI(redirect); err != nil {
		return resp, ErrInvalidUrlFormat
	}

	res, err := c.client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetContext(ctx).
		SetBody(strings.NewReader(url.Values{
			"grant_type":    []string{"authorization_code"},
			"code":          []string{code},
			"redirect_uri":  []string{redirect},
			"code_verifier": []string{verifier},
		}.Encode())).
		SetBasicAuth(c.clientID, c.clientSecret).
		SetResult(&resp).
		Post("/oauth/token")

	if err != nil {
		return resp, err
	}

	if res.StatusCode() != http.StatusOK {
		return resp, &UnexpectedStatusCodeError{
			Action: "get access token",
			Code:   res.StatusCode(),
		}
	}

	return resp, resp.Validate()
}

func (c zoomAuthClient) GetZoomUser(ctx context.Context, token string) (model.User, error) {
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
			Action: "get current zoom user",
			Code:   res.StatusCode(),
		}
	}

	return resp, resp.Validate()
}

func (c zoomAuthClient) GetDeeplink(ctx context.Context, token model.Token) (string, error) {
	var resp map[string]string

	if err := token.Validate(); err != nil {
		return "", err
	}

	_, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(token.AccessToken).
		SetContext(ctx).
		SetBody(map[string]string{
			"action": "client",
		}).
		SetResult(&resp).
		Post("/v2/zoomapp/deeplink")

	if err != nil {
		return "", err
	}

	if val, ok := resp["deeplink"]; !ok {
		return "", ErrEmptyDeeplinkResponse
	} else {
		return val, nil
	}
}
