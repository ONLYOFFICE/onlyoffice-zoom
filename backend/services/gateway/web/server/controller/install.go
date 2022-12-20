package controller

import (
	"fmt"
	"net/http"
	"net/url"

	plog "github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/gateway/web/server/middleware/security"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/constants"
	"github.com/gorilla/sessions"
	cv "github.com/nirasan/go-oauth-pkce-code-verifier"
)

type installController struct {
	logger   plog.Logger
	store    sessions.Store
	clientID string
}

func NewInstallController(
	logger plog.Logger,
	store sessions.Store,
	clientID string,
) *installController {
	return &installController{
		logger:   logger,
		store:    store,
		clientID: clientID,
	}
}

func (c installController) BuildGetInstall(redirectURL string) http.HandlerFunc {
	c.logger.Debugf("building a get install endpoint with redirectURL: %s", redirectURL)
	return func(rw http.ResponseWriter, r *http.Request) {
		v, _ := cv.CreateCodeVerifier()
		verifier := v.String()

		session, err := c.store.Get(r, constants.SESSION_KEY)
		if err != nil {
			c.logger.Errorf("could not get a session. Reason: %s", err.Error())
			http.Redirect(rw, r, "https://onlyoffice.com", http.StatusMovedPermanently)
			return
		}

		session.Values[constants.SESSION_KEY_VERIFIER] = verifier

		state, err := security.GenerateState(verifier)
		if err != nil {
			c.logger.Errorf("could not generate a new state. Reason: %s", err.Error())
			http.Redirect(rw, r, "https://onlyoffice.com", http.StatusMovedPermanently)
			return
		}

		session.Values[constants.SESSION_KEY_STATE] = state
		if err := session.Save(r, rw); err != nil {
			c.logger.Errorf("could not save session. Reason: %s", err.Error())
			http.Redirect(rw, r, "https://onlyoffice.com", http.StatusMovedPermanently)
			return
		}

		params := url.Values{
			"redirect_uri":          {redirectURL},
			"response_type":         {"code"},
			"client_id":             {c.clientID},
			"state":                 {state},
			"code_challenge":        {v.CodeChallengeS256()},
			"code_challenge_method": {"S256"},
		}

		redirectURL := fmt.Sprintf("%s/oauth/authorize?%s", constants.ZOOM_HOST, params.Encode())
		c.logger.Debugf("redirecting from installation to %s", redirectURL)
		http.Redirect(rw, r, redirectURL, http.StatusMovedPermanently)
	}
}
