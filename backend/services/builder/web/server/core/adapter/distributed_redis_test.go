package adapter

import (
	"context"
	"testing"
	"time"

	"github.com/ONLYOFFICE/zoom-onlyoffice/services/builder/web/server/core/domain"
	"github.com/stretchr/testify/assert"
)

var session = domain.Session{
	Filename: "mock-name",
	FileURL:  "https://example.com",
	Owner:    "mock-owner",
	DocKey:   "mock-dockey",
}

func TestRedisAdapter(t *testing.T) {
	adapter, err := NewRedisSessionAdapter(
		WithBufferSize(100),
		WithRedisAddresses([]string{"0.0.0.0:6379"}),
	)

	assert.NoError(t, err)

	t.Run("get invalid session", func(t *testing.T) {
		adapter.DeleteSessionByMeetingID(context.Background(), "mock-id")
		session, err := adapter.SelectSessionByMettingID(context.Background(), "mock-id")
		assert.Error(t, err)
		assert.Empty(t, session)
	})

	t.Run("save a new session with timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 0*time.Second)
		defer cancel()
		_, err := adapter.InsertSession(ctx, "mock-id", session, 1*time.Hour)
		assert.Error(t, err)
	})

	t.Run("save a new session", func(t *testing.T) {
		s, err := adapter.InsertSession(context.Background(), "mock-id", session, 1*time.Hour)
		assert.NoError(t, err)
		assert.NotEmpty(t, s)
	})

	t.Run("save an existing session", func(t *testing.T) {
		s, err := adapter.InsertSession(context.Background(), "mock-id", session, 1*time.Hour)
		assert.Error(t, err)
		assert.NotEmpty(t, s)
	})

	t.Run("get an existing session", func(t *testing.T) {
		s, err := adapter.SelectSessionByMettingID(context.Background(), "mock-id")
		assert.NoError(t, err)
		assert.NotEmpty(t, session)
		assert.Equal(t, session, s)
	})

	t.Run("remove an existing session", func(t *testing.T) {
		err := adapter.DeleteSessionByMeetingID(context.Background(), "mock-id")
		assert.NoError(t, err)
	})
}
