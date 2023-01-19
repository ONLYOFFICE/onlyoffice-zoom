package port

import (
	"context"
	"time"

	"github.com/ONLYOFFICE/zoom-onlyoffice/services/builder/web/server/core/domain"
)

type SessionServiceAdapter interface {
	InsertSession(ctx context.Context, mid string, session domain.Session, expiresIn time.Duration) (domain.Session, error)
	SelectSessionByMettingID(ctx context.Context, mid string) (domain.Session, error)
	DeleteSessionByMeetingID(ctx context.Context, mid string) error
}
