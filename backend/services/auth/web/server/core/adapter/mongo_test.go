package adapter

import (
	"context"
	"testing"
	"time"

	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/web/server/core/domain"
	"github.com/stretchr/testify/assert"
)

var user = domain.UserAccess{
	ID:           "mock",
	AccessToken:  "mock",
	RefreshToken: "mock",
	TokenType:    "mock",
	Scope:        "mock",
	ExpiresAt:    123456,
}

func TestMongoAdapter(t *testing.T) {
	adapter := NewMongoUserAdapter("mongodb://localhost:27017")

	t.Run("save user with timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 0*time.Second)
		defer cancel()
		assert.Error(t, adapter.InsertUser(ctx, user))
	})

	t.Run("save user", func(t *testing.T) {
		assert.NoError(t, adapter.InsertUser(context.Background(), user))
	})

	t.Run("save the same user", func(t *testing.T) {
		assert.NoError(t, adapter.InsertUser(context.Background(), user))
	})

	t.Run("get user by id with timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 0*time.Second)
		defer cancel()
		_, err := adapter.SelectUserByID(ctx, "mock")
		assert.Error(t, err)
	})

	t.Run("get user by id", func(t *testing.T) {
		u, err := adapter.SelectUserByID(context.Background(), "mock")
		assert.NoError(t, err)
		assert.Equal(t, user, u)
	})

	t.Run("delete user by id with timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 0*time.Second)
		defer cancel()
		assert.Error(t, adapter.DeleteUserByID(ctx, "mock"))
	})

	t.Run("delete user by id", func(t *testing.T) {
		assert.NoError(t, adapter.DeleteUserByID(context.Background(), "mock"))
	})

	t.Run("get invalid user", func(t *testing.T) {
		_, err := adapter.SelectUserByID(context.Background(), "mock")
		assert.Error(t, err)
	})

	t.Run("invald user update", func(t *testing.T) {
		_, err := adapter.UpsertUser(context.Background(), domain.UserAccess{
			ID:          "mock",
			AccessToken: "BRuh",
		})
		assert.Error(t, err)
	})

	t.Run("update user with timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 0*time.Second)
		defer cancel()
		_, err := adapter.UpsertUser(ctx, domain.UserAccess{
			ID:           "mock",
			AccessToken:  "BRuh",
			RefreshToken: "BRUH",
			TokenType:    "mock",
			Scope:        "mock",
			ExpiresAt:    123456,
		})
		assert.Error(t, err)
	})

	t.Run("update user", func(t *testing.T) {
		_, err := adapter.UpsertUser(context.Background(), domain.UserAccess{
			ID:           "mock",
			AccessToken:  "BRuh",
			RefreshToken: "BRUH",
			TokenType:    "mock",
			Scope:        "mock",
			ExpiresAt:    123456,
		})
		assert.NoError(t, err)
	})

	t.Run("get updated user", func(t *testing.T) {
		u, err := adapter.SelectUserByID(context.Background(), "mock")
		assert.NoError(t, err)
		assert.Equal(t, "BRuh", u.AccessToken)
	})

	adapter.DeleteUserByID(context.Background(), "mock")
}
