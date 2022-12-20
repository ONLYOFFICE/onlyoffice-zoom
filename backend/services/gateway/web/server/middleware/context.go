package middleware

import (
	"context"
	"net/http"
	"time"

	plog "github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/gateway/web/server/middleware/security"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/constants"
)

func BuildHandleZoomContextMiddleware(
	logger plog.Logger,
	secret string,
) func(next http.Handler) http.Handler {
	logger.Debugf("zoom context middleware has been built with secret: %s", secret)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Set("Content-Type", "application/json")
			zctx := r.Header.Get(constants.ZOOM_CONTEXT_HEADER)

			if zctx == "" {
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx, err := security.ExtractZoomContext(zctx, secret)

			if err != nil {
				rw.WriteHeader(http.StatusForbidden)
				return
			}

			now := int(time.Now().UnixMilli()) - 4*60*1000
			logger.Debugf("zoom context expiration: %d. Time now: %d", ctx.Exp, now)
			if ctx.Exp < now {
				logger.Errorf("zoom context has expired: %d. Time now: %d", ctx.Exp, now)
				rw.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(rw, r.WithContext(context.WithValue(r.Context(), security.ZoomContext{}, ctx)))
		})
	}
}
