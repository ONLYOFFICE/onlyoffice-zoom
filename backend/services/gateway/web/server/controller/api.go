package controller

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	plog "github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/gateway/web/server/middleware/security"
	zclient "github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/client"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/request"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/response"
	"go-micro.dev/v4/client"
)

type apiController struct {
	logger  plog.Logger
	client  client.Client
	zoomAPI zclient.ZoomApi
}

// TODO: Distributed cache
func NewAPIController(
	logger plog.Logger,
	client client.Client,
	zoomAPI zclient.ZoomApi,
) *apiController {
	return &apiController{
		logger:  logger,
		client:  client,
		zoomAPI: zoomAPI,
	}
}

func (c apiController) BuildGetFiles() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		query := r.URL.Query()
		searchKey := strings.TrimSpace(query.Get("search_key"))
		pageSize := strings.TrimSpace(query.Get("page_size"))
		nextPage := strings.TrimSpace(query.Get("next_page_token"))

		zctx, ok := r.Context().Value(security.ZoomContext{}).(security.ZoomContext)
		if !ok {
			rw.WriteHeader(http.StatusForbidden)
			c.logger.Error("could not extract zoom context from the context")
			return
		}

		params := map[string]string{
			"to_contact": zctx.Uid,
			"page_size":  "10",
		}

		if searchKey != "" {
			c.logger.Debugf("appending search_key to the request: %s", searchKey)
			params["search_key"] = searchKey
		}

		if nextPage != "" {
			c.logger.Debugf("appending next_page_token to the request: %s", nextPage)
			params["next_page_token"] = nextPage
		}

		if _, err := strconv.ParseInt(pageSize, 0, 8); err == nil {
			c.logger.Debugf("appending page_size to the request: %s", pageSize)
			params["page_size"] = pageSize
		}

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		var ures response.UserResponse
		if err := c.client.Call(r.Context(), c.client.NewRequest("onlyoffice:auth", "UserHandler.GetUser", zctx.Uid), &ures); err != nil {
			c.logger.Errorf("could not get user access info: %s", err.Error())
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				rw.WriteHeader(http.StatusRequestTimeout)
				return
			}

			microErr := response.MicroError{}
			if err := json.Unmarshal([]byte(err.Error()), &microErr); err != nil {
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			rw.WriteHeader(microErr.Code)
			return
		}

		c.logger.Debugf("got a user response: %v", ures)
		res, err := c.zoomAPI.GetFilesFromMessages(ctx, ures.AccessToken, params)
		if err != nil {
			c.logger.Errorf("could not get files messages: %s", err.Error())
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				rw.WriteHeader(http.StatusRequestTimeout)
				return
			}

			microErr := response.MicroError{}
			if err := json.Unmarshal([]byte(err.Error()), &microErr); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			c.logger.Errorf("get files micro error: %s", microErr.Detail)
			rw.WriteHeader(microErr.Code)
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write(res.ToJSON())
	}
}

func (c apiController) BuildGetConfig() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		query := r.URL.Query()
		fileName, fileURL := strings.TrimSpace(query.Get("file_name")), strings.TrimSpace(query.Get("file_url"))

		zctx, ok := r.Context().Value(security.ZoomContext{}).(security.ZoomContext)
		if !ok {
			rw.WriteHeader(http.StatusForbidden)
			c.logger.Error("could not extract zoom context from the context")
			return
		}

		if fileName == "" {
			rw.WriteHeader(http.StatusBadRequest)
			c.logger.Error("could not extract file_name from URL Query")
			return
		}

		if fileURL == "" {
			rw.WriteHeader(http.StatusBadRequest)
			c.logger.Error("could not extract file_url from URL Query")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		var resp response.BuildConfigResponse
		if err := c.client.Call(ctx, c.client.NewRequest("onlyoffice:builder", "ConfigHandler.BuildConfig", request.BuildConfigRequest{
			Uid:       zctx.Uid,
			Mid:       zctx.Mid,
			UserAgent: r.UserAgent(),
			Filename:  fileName,
			FileURL:   fileURL,
		}), &resp); err != nil {
			c.logger.Errorf("could not build onlyoffice config: %s", err.Error())
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				rw.WriteHeader(http.StatusRequestTimeout)
				return
			}

			microErr := response.MicroError{}
			if err := json.Unmarshal([]byte(err.Error()), &microErr); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			c.logger.Errorf("build config micro error: %s", microErr.Detail)
			rw.WriteHeader(microErr.Code)
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write(resp.ToJSON())
	}
}

func (c apiController) BuildDeleteSession() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		zctx, ok := r.Context().Value(security.ZoomContext{}).(security.ZoomContext)
		if !ok {
			rw.WriteHeader(http.StatusForbidden)
			c.logger.Error("could not extract zoom context from the context")
			return
		}

		if zctx.Mid != "" {
			req := c.client.NewRequest("onlyoffice:builder", "SessionHandler.OwnerRemoveSession", request.OwnerRemoveSessionRequest{
				Mid: zctx.Mid,
				Uid: zctx.Uid,
			})
			var resp interface{}

			rctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
			defer cancel()

			if err := c.client.Call(rctx, req, &resp); err != nil {
				c.logger.Errorf("could not build remove owner session: %s", err.Error())
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					rw.WriteHeader(http.StatusRequestTimeout)
					return
				}

				microErr := response.MicroError{}
				if err := json.Unmarshal([]byte(err.Error()), &microErr); err != nil {
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}

				c.logger.Errorf("delete session micro error: %s", microErr.Detail)
				rw.WriteHeader(microErr.Code)
				return
			}
		}

		rw.WriteHeader(http.StatusOK)
	}
}
