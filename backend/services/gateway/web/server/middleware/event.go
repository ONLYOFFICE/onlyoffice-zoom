package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	plog "github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/gateway/web/server/middleware/security"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/constants"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/request"
)

func BuildHandleZoomEventMiddleware(
	logger plog.Logger,
	secret string,
) func(next http.Handler) http.Handler {
	logger.Debugf("zoom event middleware has been built with secret: %s", secret)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Set("Content-Type", "application/json")
			signature := r.Header.Get(constants.ZOOM_EVENT_SIGNATURE_HEADER)
			ts := r.Header.Get(constants.ZOOM_EVENT_TIMESTAMP_HEADER)

			if signature == "" || ts == "" {
				logger.Errorf("an unauthorized access to deauth endpoint")
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			b, err := io.ReadAll(r.Body)
			if err != nil {
				logger.Errorf("could not decode a deauth body")
				rw.WriteHeader(http.StatusForbidden)
				return
			}

			message := []byte(fmt.Sprintf("v0:%s:%v", ts, string(b)))
			h := hmac.New(sha256.New, []byte(secret))
			h.Write(message)
			sig := "v0=" + hex.EncodeToString(h.Sum(nil))

			if signature != sig {
				logger.Errorf("deauth signatures don't match: %s", r.Header.Get("x-real-ip"))
				rw.WriteHeader(http.StatusForbidden)
				return
			}

			var body request.EventRequest
			if err := json.Unmarshal(b, &body); err != nil {
				logger.Errorf("could not unmarshal a request body: %s. Reason: %s", string(b), err.Error())
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			next.ServeHTTP(rw, r.WithContext(context.WithValue(r.Context(), security.ZoomContext{}, body)))
		})
	}
}
