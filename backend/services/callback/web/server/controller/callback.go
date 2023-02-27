package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	plog "github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/worker"
	zclient "github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/client"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/crypto"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/message"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/request"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/response"
	"go-micro.dev/v4/client"
)

type callbackController struct {
	namespace     string
	maxSize       int64
	client        client.Client
	enqueuer      worker.BackgroundEnqueuer
	zoomFilestore zclient.ZoomFilestore
	jwtManager    crypto.JwtManager
	logger        plog.Logger
}

func NewCallbackController(
	namespace string,
	maxSize int64,
	client client.Client,
	enqueuer worker.BackgroundEnqueuer,
	jwtManager crypto.JwtManager,
	logger plog.Logger,
) *callbackController {
	return &callbackController{
		namespace:     namespace,
		maxSize:       maxSize,
		enqueuer:      enqueuer,
		client:        client,
		zoomFilestore: zclient.NewZoomFilestoreClient(logger),
		jwtManager:    jwtManager,
		logger:        logger,
	}
}

func (c callbackController) getSessionOwner(ctx context.Context, mid string) (string, error) {
	req := c.client.NewRequest(fmt.Sprintf("%s:builder", c.namespace), "SessionHandler.GetSessionOwner", mid)
	var resp string
	if err := c.client.Call(ctx, req, &resp); err != nil {
		c.logger.Errorf("could not get session owner: %s", err.Error())
		return "", err
	}

	return resp, nil
}

func (c callbackController) BuildPostHandleCallback() http.HandlerFunc {
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
			rctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
			defer cancel()

			req := c.client.NewRequest(fmt.Sprintf("%s:builder", c.namespace), "SessionHandler.RefreshSession", mid)
			var res interface{}
			if err := c.client.Call(rctx, req, &res); err != nil {
				c.logger.Errorf("could not refresh initial session with key %s", body.Key)
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write(response.CallbackResponse{
					Error: 1,
				}.ToJSON())
				return
			}

			if err := c.client.Publish(rctx, client.NewMessage("notify-session", message.SessionMessage{
				MID:       mid,
				InSession: true,
			})); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				c.logger.Errorf("remove session error: %s", err.Error())
				rw.Write(response.CallbackResponse{
					Error: 1,
				}.ToJSON())
				return
			}
		}

		if body.Status == 4 {
			if mid != "" {
				rctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
				defer cancel()

				if err := c.client.Publish(rctx, client.NewMessage("remove-session", mid)); err != nil {
					rw.WriteHeader(http.StatusInternalServerError)
					c.logger.Errorf("remove session error: %s", err.Error())
					rw.Write(response.CallbackResponse{
						Error: 1,
					}.ToJSON())
					return
				}

				time.Sleep(200 * time.Millisecond)
				rctx, cancel = context.WithTimeout(r.Context(), 5*time.Second)
				defer cancel()
				if err := c.client.Publish(rctx, client.NewMessage("notify-session", message.SessionMessage{
					MID:       mid,
					InSession: false,
				})); err != nil {
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

			ferrChan := make(chan error)
			serrChan := make(chan error)
			var wg sync.WaitGroup

			ctx, cancel := context.WithTimeout(r.Context(), 16*time.Second)
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
						ferrChan <- err
						return
					}

					if err := c.enqueuer.Enqueue("callback-upload", message.JobMessage{
						UID:      usr,
						Filename: filename,
						Url:      body.URL,
					}.ToJSON(), worker.WithMaxRetry(3)); err != nil {
						c.logger.Errorf("could not enqueue a new task: %s", err.Error())
						ferrChan <- err
						return
					}
				}()
			}

			if mid != "" {
				wg.Add(1)
				go func() {
					defer wg.Done()
					if err := c.client.Publish(ctx, client.NewMessage("remove-session", mid)); err != nil {
						serrChan <- err
						return
					}

					if err := c.client.Publish(ctx, client.NewMessage("notify-session", message.SessionMessage{
						MID:       mid,
						InSession: false,
					})); err != nil {
						serrChan <- err
						return
					}
				}()
			}

			wg.Wait()

			select {
			case <-ferrChan:
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write(response.CallbackResponse{
					Error: 1,
				}.ToJSON())
				return
			case <-serrChan:
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
