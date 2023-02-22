package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	plog "github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/gateway/web/server/middleware/security"
	zclient "github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/client"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/constants"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/request"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/response"
	"github.com/gorilla/sessions"
	"github.com/mitchellh/mapstructure"
	"go-micro.dev/v4/client"
)

type authController struct {
	namespace string
	logger    plog.Logger
	store     sessions.Store
	client    client.Client
	zoomAPI   zclient.ZoomAuth
	timeout   int
}

func NewAuthController(
	namespace string,
	logger plog.Logger,
	store sessions.Store,
	client client.Client,
	zoomAPI zclient.ZoomAuth,
	timeout int,
) *authController {
	return &authController{
		namespace: namespace,
		logger:    logger,
		store:     store,
		client:    client,
		zoomAPI:   zoomAPI,
		timeout:   timeout,
	}
}

func (c authController) BuildGetAuth(redirectURL string) http.HandlerFunc {
	c.logger.Debugf("building a get auth endpoint with redirectURL: %s", redirectURL)
	return func(rw http.ResponseWriter, r *http.Request) {
		tctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()
		query := r.URL.Query()
		code := strings.TrimSpace(query.Get("code"))

		if code == "" {
			rw.WriteHeader(http.StatusBadRequest)
			c.logger.Debug("empty auth code parameter")
			return
		}

		c.logger.Debugf("auth code is valid: %s", code)

		session, err := c.store.Get(r, constants.SESSION_KEY)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			c.logger.Errorf("could not get session. Reason: %s", err.Error())
			return
		}

		state := strings.TrimSpace(query.Get("state"))
		if state != session.Values["state"] {
			http.Redirect(rw, r, "/oauth/install", http.StatusMovedPermanently)
			return
		}

		c.logger.Debugf("auth state is valid: %s", state)

		vefifier, ok := session.Values["verifier"].(string)
		if !ok {
			rw.WriteHeader(http.StatusInternalServerError)
			c.logger.Debugf("could not cast verifier: %v", vefifier)
			return
		}

		c.logger.Debugf("verifier is valid: %s", vefifier)

		session.Options.MaxAge = -1
		if err := session.Save(r, rw); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			c.logger.Errorf("could not remove session. Reason: %s", err.Error())
			return
		}

		token, err := c.zoomAPI.GetZoomAccessToken(tctx, code, vefifier, redirectURL)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			c.logger.Errorf("get zoom access token error: %s", err.Error())
			return
		}

		c.logger.Debugf("got zoom access token: %s", token.AccessToken)

		user, err := c.zoomAPI.GetZoomUser(tctx, token.AccessToken)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			c.logger.Errorf("get zoom user error: %s", err.Error())
			return
		}

		c.logger.Debugf("got zoom user with id: %s", user.ID)
		if err := c.client.Publish(tctx, client.NewMessage("insert-auth", response.UserResponse{
			ID:           user.ID,
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			TokenType:    token.TokenType,
			Scope:        token.Scope,
			ExpiresAt:    time.Now().Local().Add(time.Second * time.Duration(token.ExpiresIn-700)).UnixMilli(),
		})); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			c.logger.Errorf("insert user error: %s", err.Error())
			return
		}

		if deepLink, err := c.zoomAPI.GetDeeplink(tctx, token); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			c.logger.Errorf("get zoom deeplink error: %s", err.Error())
			return
		} else {
			c.logger.Debugf("redirecting to deeplink: %s", deepLink)
			http.Redirect(rw, r, deepLink, http.StatusMovedPermanently)
		}
	}
}

func (c authController) BuildPostDeauth() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		c.logger.Debug("got a deauth request")
		event, ok := r.Context().Value(security.ZoomContext{}).(request.EventRequest)

		if !ok {
			rw.WriteHeader(http.StatusForbidden)
			c.logger.Error("could not extract zoom deauth event from the context")
			return
		}

		var deauthEvent request.DeauthorizationPayload
		if err := mapstructure.Decode(event.Payload, &deauthEvent); err != nil {
			rw.WriteHeader(http.StatusForbidden)
			c.logger.Errorf("could not decode an event payload: %s", err.Error())
			return
		}

		c.logger.Debugf("got a deauth event: %v", event)

		tctx, cancel := context.WithTimeout(r.Context(), time.Duration(c.timeout)*time.Millisecond)
		defer cancel()
		var res interface{}
		if err := c.client.Call(tctx, c.client.NewRequest(fmt.Sprintf("%s:auth", c.namespace), "UserDeleteHandler.DeleteUser", deauthEvent.Uid), &res); err != nil {
			c.logger.Errorf("could not delete user: %s", err.Error())
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

		c.logger.Debugf("deleted user %s", deauthEvent.Uid)
		rw.WriteHeader(http.StatusOK)
	}
}
