package message

import (
	"context"
	"time"

	plog "github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/builder/web/server/core/port"
)

type RemoveSessionMessageHandler struct {
	service port.SessionService
	logger  plog.Logger
}

func BuildRemoveSessionMessageHandler(service port.SessionService, logger plog.Logger) RemoveSessionMessageHandler {
	return RemoveSessionMessageHandler{
		service: service,
		logger:  logger,
	}
}

func (i RemoveSessionMessageHandler) GetHandler() func(context.Context, interface{}) error {
	return func(ctx context.Context, payload interface{}) error {
		tctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if mid, ok := payload.(string); !ok {
			return _ErrInvalidHandlerPayload
		} else {
			return i.service.DeleteSession(tctx, mid)
		}
	}
}
