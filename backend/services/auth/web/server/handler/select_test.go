package handler

import (
	"context"
	"testing"
	"time"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/web/server/core/adapter"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/web/server/core/domain"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/web/server/core/service"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/client"
	"github.com/stretchr/testify/assert"
)

type mockEncryptor struct{}

func (e mockEncryptor) Encrypt(text string) (string, error) {
	return string(text), nil
}

func (e mockEncryptor) Decrypt(ciphertext string) (string, error) {
	return string(ciphertext), nil
}

func TestSelectCaching(t *testing.T) {
	adapter := adapter.NewMemoryUserAdapter()
	service := service.NewUserService(log.NewDefaultLogger(), adapter, mockEncryptor{})
	zm := client.NewZoomClient("clientID", "clientSecret")

	sel := NewUserSelectHandler(service, nil, zm, log.NewEmptyLogger())

	service.CreateUser(context.Background(), domain.UserAccess{
		ID:           "mock",
		AccessToken:  "mock",
		RefreshToken: "mock",
		TokenType:    "mock",
		Scope:        "mock",
		ExpiresAt:    time.Now().Add(24 * time.Hour).UnixMilli(),
	})

	t.Run("get user", func(t *testing.T) {
		var res domain.UserAccess
		id := "mock"
		assert.NoError(t, sel.GetUser(context.Background(), &id, &res))
		assert.NotEmpty(t, res)
	})
}
