package message

import (
	"context"

	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/web/server/core/port"
)

type DeleteMessageHandler struct {
	service port.UserAccessService
}

func BuildDeleteMessageHandler(service port.UserAccessService) DeleteMessageHandler {
	return DeleteMessageHandler{
		service: service,
	}
}

func (i DeleteMessageHandler) GetHandler() func(context.Context, interface{}) error {
	return func(ctx context.Context, payload interface{}) error {
		if uid, ok := payload.(string); !ok {
			return _ErrInvalidHandlerPayload
		} else {
			return i.service.DeleteUser(ctx, uid)
		}
	}
}
