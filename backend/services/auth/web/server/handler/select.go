package handler

import (
	"context"
	"time"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/web/server/core/domain"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/web/server/core/port"
	zclient "github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/client"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/response"
	"go-micro.dev/v4/client"
	"golang.org/x/sync/singleflight"
)

var group singleflight.Group

type UserSelectHandler struct {
	service port.UserAccessService
	client  client.Client
	zoomAPI zclient.ZoomAuth
	logger  log.Logger
}

// TODO: Distributed cache
func NewUserSelectHandler(
	service port.UserAccessService,
	client client.Client,
	zoomAPI zclient.ZoomAuth,
	logger log.Logger,
) UserSelectHandler {
	return UserSelectHandler{
		service: service,
		client:  client,
		zoomAPI: zoomAPI,
		logger:  logger,
	}
}

func (u UserSelectHandler) GetUser(ctx context.Context, uid *string, res *domain.UserAccess) error {
	user, err, _ := group.Do(*uid, func() (interface{}, error) {
		user, err := u.service.GetUser(ctx, *uid)
		if err != nil {
			u.logger.Errorf("could not get user with id: %s. Reason: %s", *uid, err.Error())
			return nil, err
		}

		if user.ExpiresAt <= time.Now().UnixMilli() {
			u.logger.Debug("user token has expired. Trying to refresh!")
			token, terr := u.zoomAPI.RefreshZoomAccessToken(ctx, user.RefreshToken)
			if terr != nil {
				u.logger.Errorf("could not refresh user's %s token. Reason: %s", *uid, terr.Error())
				return nil, terr
			}

			u.logger.Debugf("user's %s token has been refreshed", *uid)
			access := domain.UserAccess{
				ID:           user.ID,
				AccessToken:  token.AccessToken,
				RefreshToken: token.RefreshToken,
				TokenType:    token.TokenType,
				Scope:        token.Scope,
				ExpiresAt:    time.Now().Local().Add(time.Second * time.Duration(token.ExpiresIn-700)).UnixMilli(),
			}

			_, err := u.service.UpdateUser(ctx, access)
			if err != nil {
				u.logger.Debugf("could not persist a new user's %s token. Reason: %s. Sending a fallback message!", *uid, err.Error())

				pctx, pcancel := context.WithTimeout(ctx, 3*time.Second)
				defer pcancel()

				if err := u.client.Publish(pctx, client.NewMessage("insert-auth", response.UserResponse{
					ID:           user.ID,
					AccessToken:  access.AccessToken,
					RefreshToken: access.RefreshToken,
					TokenType:    access.TokenType,
					Scope:        access.Scope,
					ExpiresAt:    access.ExpiresAt,
				})); err != nil {
					u.logger.Errorf("fallback message to update user's %s token has failed. Reason: %s", *uid, err.Error())
					return nil, err
				}
			}

			u.logger.Debugf("user's %s token has been updated", *uid)
			return access, nil
		}

		return user, nil
	})

	if usr, ok := user.(domain.UserAccess); ok {
		*res = usr
		return nil
	}

	return err
}
