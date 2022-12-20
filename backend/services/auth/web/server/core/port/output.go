package port

import (
	"context"

	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/web/server/core/domain"
)

type UserAccessServiceAdapter interface {
	InsertUser(ctx context.Context, user domain.UserAccess) error
	SelectUserByID(ctx context.Context, uid string) (domain.UserAccess, error)
	UpsertUser(ctx context.Context, user domain.UserAccess) (domain.UserAccess, error)
	DeleteUserByID(ctx context.Context, uid string) error
}
