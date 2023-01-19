package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	plog "github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	zclient "github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/client"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/crypto"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/request"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/response"
	"github.com/gocraft/work"
	"go-micro.dev/v4/client"
)

type callbackController struct {
	maxSize       int64
	logger        plog.Logger
	client        client.Client
	zoomFilestore zclient.ZoomFilestore
	jwtManager    crypto.JwtManager
}

func NewCallbackController(
	maxSize int64,
	logger plog.Logger,
	client client.Client,
	jwtManager crypto.JwtManager,
) *callbackController {
	return &callbackController{
		maxSize:       maxSize,
		logger:        logger,
		client:        client,
		zoomFilestore: zclient.NewZoomFilestoreClient(),
		jwtManager:    jwtManager,
	}
}

func (c callbackController) getSessionOwner(ctx context.Context, mid string) (string, error) {
	req := c.client.NewRequest("onlyoffice:builder", "SessionHandler.GetSessionOwner", mid)
	var resp string
	if err := c.client.Call(ctx, req, &resp); err != nil {
		c.logger.Errorf("could not get session owner: %s", err.Error())
		return "", err
	}

	return resp, nil
}

func (c callbackController) BuildPostHandleCallback(enqueuer *work.Enqueuer) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		mid := strings.TrimSpace(r.URL.Query().Get("mid"))
		rw.Header().Set("Content-Type", "application/json")

		var body request.CallbackRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			c.logger.Errorf("could not decode a callback body")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(response.CallbackResponse{
				Error: 1,
			}.ToJSON())
			return
		}

		if body.Token == "" {
			c.logger.Error("invalid callback body token")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(response.CallbackResponse{
				Error: 1,
			}.ToJSON())
			return
		}

		if err := c.jwtManager.Verify(body.Token, &body); err != nil {
			c.logger.Errorf("could not verify callback jwt (%s). Reason: %s", body.Token, err.Error())
			rw.WriteHeader(http.StatusForbidden)
			rw.Write(response.CallbackResponse{
				Error: 1,
			}.ToJSON())
			return
		}

		if err := body.Validate(); err != nil {
			c.logger.Errorf("invalid callback body. Reason: %s", err.Error())
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(response.CallbackResponse{
				Error: 1,
			}.ToJSON())
			return
		}

		if body.Status == 1 && len(body.Users) == 1 && mid != "" {
			rctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
			defer cancel()

			req := c.client.NewRequest("onlyoffice:builder", "SessionHandler.RefreshSession", mid)
			var res interface{}
			if err := c.client.Call(rctx, req, &res); err != nil {
				c.logger.Errorf("could not refresh initial session with key %s", body.Key)
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write(response.CallbackResponse{
					Error: 1,
				}.ToJSON())
				return
			}
		}

		if body.Status == 4 {
			if mid != "" {
				rctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
				defer cancel()

				if err := c.client.Publish(rctx, client.NewMessage("remove-session", mid)); err != nil {
					rw.WriteHeader(http.StatusInternalServerError)
					c.logger.Errorf("remove session error: %s", err.Error())
					rw.Write(response.CallbackResponse{
						Error: 1,
					}.ToJSON())
					return
				}
			}
		}

		if body.Status == 2 {
			filename := strings.TrimSpace(r.URL.Query().Get("filename"))
			if filename == "" {
				rw.WriteHeader(http.StatusInternalServerError)
				c.logger.Errorf("callback request %s does not contain a filename", body.Key)
				rw.Write(response.CallbackResponse{
					Error: 1,
				}.ToJSON())
				return
			}

			errChan := make(chan error, 2)
			var wg sync.WaitGroup

			ctx, cancel := context.WithTimeout(r.Context(), 7*time.Second)
			defer cancel()

			mid := strings.TrimSpace(r.URL.Query().Get("mid"))
			usr := body.Users[0]

			if mid != "" {
				var err error
				if usr, err = c.getSessionOwner(ctx, mid); err != nil {
					rw.WriteHeader(http.StatusInternalServerError)
					c.logger.Errorf("could not extract meeting owner for %s", body.Key)
					rw.Write(response.CallbackResponse{
						Error: 1,
					}.ToJSON())
					return
				}
			}

			if usr != "" {
				wg.Add(1)
				go func() {
					defer wg.Done()
					if err := c.zoomFilestore.ValidateFileSize(ctx, c.maxSize, body.URL); err != nil {
						c.logger.Errorf("could not validate file %s: %s", filename, err.Error())
						errChan <- err
						return
					}

					if _, err := enqueuer.Enqueue("callback-upload", map[string]interface{}{
						"uid":      usr,
						"filename": filename,
						"url":      body.URL,
					}); err != nil {
						c.logger.Errorf("could not enqueue a new job")
						errChan <- err
						return
					}
				}()
			}

			if mid != "" {
				wg.Add(1)
				go func() {
					defer wg.Done()
					if err := c.client.Publish(ctx, client.NewMessage("remove-session", mid)); err != nil {
						errChan <- err
					}
				}()
			}

			wg.Wait()
			defer close(errChan)

			select {
			case <-errChan:
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write(response.CallbackResponse{
					Error: 1,
				}.ToJSON())
				return
			default:
			}
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write(response.CallbackResponse{
			Error: 0,
		}.ToJSON())
	}
}
