package controller

import (
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
	"go-micro.dev/v4/client"
)

type authController struct {
	logger  plog.Logger
	store   sessions.Store
	client  client.Client
	zoomAPI zclient.ZoomAuth
}

func NewAuthController(
	logger plog.Logger,
	store sessions.Store,
	client client.Client,
	zoomAPI zclient.ZoomAuth,
) *authController {
	return &authController{
		logger:  logger,
		store:   store,
		client:  client,
		zoomAPI: zoomAPI,
	}
}

func (c authController) BuildGetAuth(redirectURL string) http.HandlerFunc {
	c.logger.Debugf("building a get auth endpoint with redirectURL: %s", redirectURL)
	return func(rw http.ResponseWriter, r *http.Request) {
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

		token, err := c.zoomAPI.GetZoomAccessToken(r.Context(), code, vefifier, redirectURL)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			c.logger.Errorf("get zoom access token error: %s", err.Error())
			return
		}

		c.logger.Debugf("got zoom access token: %s", token.AccessToken)

		user, err := c.zoomAPI.GetZoomUser(r.Context(), token.AccessToken)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			c.logger.Errorf("get zoom user error: %s", err.Error())
			return
		}

		c.logger.Debugf("got zoom user with id: %s", user.ID)

		if err := c.client.Publish(r.Context(), client.NewMessage("insert-auth", response.UserResponse{
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

		if deepLink, err := c.zoomAPI.GetDeeplink(r.Context(), token); err != nil {
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
		event, ok := r.Context().Value(security.ZoomContext{}).(request.DeauthorizationEventRequest)

		if !ok {
			rw.WriteHeader(http.StatusForbidden)
			c.logger.Error("could not extract zoom deauth event from the context")
			return
		}

		c.logger.Debugf("got a deauth event: %v", event)

		if err := c.client.Publish(r.Context(), client.NewMessage("delete-auth", event.Payload.Uid)); err != nil {
			c.logger.Errorf("could not publish a user %s delete message", event.Payload.Uid)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		c.logger.Debugf("successfully published a deauth message: %s", event.Payload.Uid)
		rw.WriteHeader(http.StatusOK)
	}
}
