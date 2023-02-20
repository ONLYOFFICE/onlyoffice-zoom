package message

import (
	"context"
	"time"

	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/web/server/core/domain"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/web/server/core/port"
	"github.com/mitchellh/mapstructure"
)

type InsertMessageHandler struct {
	service port.UserAccessService
}

func BuildInsertMessageHandler(service port.UserAccessService) InsertMessageHandler {
	return InsertMessageHandler{
		service: service,
	}
}

func (i InsertMessageHandler) GetHandler() func(context.Context, interface{}) error {
	return func(ctx context.Context, payload interface{}) error {
		tctx, cancel := context.WithTimeout(ctx, 4*time.Second)
		defer cancel()
		var user domain.UserAccess
		if err := mapstructure.Decode(payload, &user); err != nil {
			return err
		}
		_, err := i.service.UpdateUser(tctx, user)
		return err
	}
}
