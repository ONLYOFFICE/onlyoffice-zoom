package handler

import (
	"context"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/web/server/core/domain"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/web/server/core/port"
)

type UserInsertHandler struct {
	service port.UserAccessService
	logger  log.Logger
}

func NewUserInsertHandler(
	service port.UserAccessService,
	logger log.Logger,
) UserInsertHandler {
	return UserInsertHandler{
		service: service,
		logger:  logger,
	}
}

func (u UserInsertHandler) InsertUser(ctx context.Context, req *domain.UserAccess, res *interface{}) error {
	u.logger.Debugf("trying to insert a new user: %s", req.ID)
	if _, err := u.service.UpdateUser(ctx, *req); err != nil {
		u.logger.Errorf("could not insert a new user: %s", err.Error())
		return err
	}

	return nil
}
