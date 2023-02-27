package service

import (
	"context"
	"testing"
	"time"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/builder/web/server/core/domain"
	"github.com/stretchr/testify/assert"
)

type mockAdapter struct {
}

var session = domain.Session{
	Filename: "file-name",
	FileURL:  "https://example.com",
	Owner:    "owner-id",
	DocKey:   "doc-key",
}

func (m mockAdapter) InsertSession(ctx context.Context, mid string, session domain.Session, expiresIn time.Duration) (domain.Session, error) {
	return session, nil
}

func (m mockAdapter) SelectSessionByMettingID(ctx context.Context, mid string) (domain.Session, error) {
	return session, nil
}

func (m mockAdapter) DeleteSessionByMeetingID(ctx context.Context, mid string) error {
	return nil
}

func TestUserService(t *testing.T) {
	service := NewSessionService(mockAdapter{}, log.NewDefaultLogger())

	t.Run("create session", func(t *testing.T) {
		s, err := service.CreateSession(context.Background(), "meeting-uuid", session)
		assert.NoError(t, err)
		assert.Equal(t, session, s)
	})

	t.Run("get session", func(t *testing.T) {
		u, err := service.GetSession(context.Background(), "meeting-uuid")
		assert.NoError(t, err)
		assert.Equal(t, session, u)
	})

	t.Run("delete session", func(t *testing.T) {
		assert.NoError(t, service.DeleteSession(context.Background(), "meeting-uuid"))
	})
}
