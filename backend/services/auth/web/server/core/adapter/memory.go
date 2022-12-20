package adapter

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/web/server/core/domain"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/web/server/core/port"
)

type memoryUserAdapter struct {
	kvs map[string][]byte
}

func NewMemoryUserAdapter() port.UserAccessServiceAdapter {
	return &memoryUserAdapter{
		kvs: make(map[string][]byte),
	}
}

func (m *memoryUserAdapter) save(user domain.UserAccess) error {
	buffer, err := json.Marshal(user)

	if err != nil {
		return err
	}

	m.kvs[user.ID] = buffer

	return nil
}

func (m *memoryUserAdapter) InsertUser(ctx context.Context, user domain.UserAccess) error {
	return m.save(user)
}

func (m *memoryUserAdapter) SelectUserByID(ctx context.Context, uid string) (domain.UserAccess, error) {
	buffer, ok := m.kvs[uid]
	var user domain.UserAccess

	if !ok {
		return user, errors.New("user with this id doesn't exist")
	}

	if err := json.Unmarshal(buffer, &user); err != nil {
		return user, err
	}

	return user, nil
}

func (m *memoryUserAdapter) UpsertUser(ctx context.Context, user domain.UserAccess) (domain.UserAccess, error) {
	if err := m.save(user); err != nil {
		return domain.UserAccess{}, err
	}

	return user, nil
}

func (m *memoryUserAdapter) DeleteUserByID(ctx context.Context, uid string) error {
	if _, ok := m.kvs[uid]; !ok {
		return errors.New("user with this id doesn't exist")
	}

	delete(m.kvs, uid)

	return nil
}
