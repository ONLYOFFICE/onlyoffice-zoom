package port

import (
	"context"

	"github.com/ONLYOFFICE/zoom-onlyoffice/services/builder/web/server/core/domain"
)

type SessionService interface {
	CreateSession(ctx context.Context, mid string, session domain.Session) (domain.Session, error)
	GetSession(ctx context.Context, mid string) (domain.Session, error)
	DeleteSession(ctx context.Context, mid string) error
}
