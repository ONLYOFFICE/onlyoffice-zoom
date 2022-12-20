package adapter

import (
	"context"
	"testing"

	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/web/server/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestMemoryAdapter(t *testing.T) {
	adapter := NewMemoryUserAdapter()

	t.Run("save user", func(t *testing.T) {
		assert.NoError(t, adapter.InsertUser(context.Background(), user))
	})

	t.Run("save the same user", func(t *testing.T) {
		assert.NoError(t, adapter.InsertUser(context.Background(), user))
	})

	t.Run("get user by id", func(t *testing.T) {
		u, err := adapter.SelectUserByID(context.Background(), "mock")
		assert.NoError(t, err)
		assert.Equal(t, user, u)
	})

	t.Run("update user by id", func(t *testing.T) {
		u, err := adapter.UpsertUser(context.Background(), domain.UserAccess{
			ID:          "mock",
			AccessToken: "BRuh",
		})
		assert.NoError(t, err)
		assert.NotNil(t, u)
	})

	t.Run("delete user by id", func(t *testing.T) {
		assert.NoError(t, adapter.DeleteUserByID(context.Background(), "mock"))
	})

	t.Run("get invalid user", func(t *testing.T) {
		_, err := adapter.SelectUserByID(context.Background(), "mock")
		assert.Error(t, err)
	})
}
