package handler

import (
	"context"

	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/web/server/core/port"
)

type UserDeleteHandler struct {
	service port.UserAccessService
}

func NewUserDeleteHandler(service port.UserAccessService) UserDeleteHandler {
	return UserDeleteHandler{
		service: service,
	}
}

func (i UserDeleteHandler) DeleteUser(ctx context.Context, uid *string, res *interface{}) error {
	return i.service.DeleteUser(ctx, *uid)
}
