package service

import (
	"context"
	"errors"
	"strings"
	"sync"

	plog "github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/web/server/core/domain"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/web/server/core/port"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/crypto"
)

var _ErrOperationTimeout = errors.New("user operation timeout")

type userService struct {
	logger    plog.Logger
	adapter   port.UserAccessServiceAdapter
	encryptor crypto.Encryptor
}

func NewUserService(
	logger plog.Logger,
	adapter port.UserAccessServiceAdapter,
	encryptor crypto.Encryptor,
) port.UserAccessService {
	return userService{
		logger:    logger,
		adapter:   adapter,
		encryptor: encryptor,
	}
}

func (s userService) CreateUser(ctx context.Context, user domain.UserAccess) error {
	s.logger.Debugf("validating user %s to perform a persist action", user.ID)
	if err := user.Validate(); err != nil {
		return err
	}

	var wg sync.WaitGroup
	atokenErrChan := make(chan error)
	rtokenErrChan := make(chan error)
	atokenChan := make(chan string)
	rtokenChan := make(chan string)

	go func() {
		wg.Add(1)
		defer wg.Done()
		defer close(atokenChan)
		defer close(atokenErrChan)
		aToken, err := s.encryptor.Encrypt(user.AccessToken)
		if err != nil {
			atokenErrChan <- err
			return
		}
		atokenChan <- aToken
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		defer close(rtokenChan)
		defer close(rtokenErrChan)
		rToken, err := s.encryptor.Encrypt(user.RefreshToken)
		if err != nil {
			rtokenErrChan <- err
			return
		}
		rtokenChan <- rToken
	}()

	wg.Wait()

	select {
	case err := <-atokenErrChan:
		return err
	case err := <-rtokenErrChan:
		return err
	case <-ctx.Done():
		return _ErrOperationTimeout
	default:
		return s.adapter.InsertUser(ctx, domain.UserAccess{
			ID:           user.ID,
			AccessToken:  <-atokenChan,
			RefreshToken: <-rtokenChan,
			TokenType:    user.TokenType,
			Scope:        user.Scope,
			ExpiresAt:    user.ExpiresAt,
		})
	}
}

func (s userService) GetUser(ctx context.Context, uid string) (domain.UserAccess, error) {
	s.logger.Debugf("trying to select user with id: %s", uid)
	id := strings.TrimSpace(uid)

	if id == "" {
		return domain.UserAccess{}, &InvalidServiceParameterError{
			Name:   "UID",
			Reason: "Should not be blank",
		}
	}

	user, err := s.adapter.SelectUserByID(ctx, id)
	if err != nil {
		return domain.UserAccess{}, err
	}

	var wg sync.WaitGroup
	atokenErrChan := make(chan error)
	rtokenErrChan := make(chan error)
	atokenChan := make(chan string)
	rtokenChan := make(chan string)

	go func() {
		wg.Add(1)
		defer wg.Done()
		defer close(atokenChan)
		defer close(atokenErrChan)
		aToken, err := s.encryptor.Decrypt(user.AccessToken)
		if err != nil {
			atokenErrChan <- err
			return
		}
		atokenChan <- aToken
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		defer close(rtokenChan)
		defer close(rtokenErrChan)
		rToken, err := s.encryptor.Decrypt(user.RefreshToken)
		if err != nil {
			rtokenErrChan <- err
			return
		}
		rtokenChan <- rToken
	}()

	wg.Wait()

	select {
	case err := <-atokenErrChan:
		return domain.UserAccess{}, err
	case err := <-rtokenErrChan:
		return domain.UserAccess{}, err
	case <-ctx.Done():
		return domain.UserAccess{}, _ErrOperationTimeout
	default:
		return domain.UserAccess{
			ID:           user.ID,
			AccessToken:  <-atokenChan,
			RefreshToken: <-rtokenChan,
			TokenType:    user.TokenType,
			Scope:        user.Scope,
			ExpiresAt:    user.ExpiresAt,
		}, nil
	}
}

func (s userService) UpdateUser(ctx context.Context, user domain.UserAccess) (domain.UserAccess, error) {
	s.logger.Debugf("validating user %s to perform an update action", user.ID)
	if err := user.Validate(); err != nil {
		return domain.UserAccess{}, err
	}

	s.logger.Debugf("user %s is valid to perform an update action", user.ID)

	var wg sync.WaitGroup
	atokenErrChan := make(chan error)
	rtokenErrChan := make(chan error)
	atokenChan := make(chan string)
	rtokenChan := make(chan string)

	go func() {
		wg.Add(1)
		defer wg.Done()
		defer close(atokenChan)
		defer close(atokenErrChan)
		aToken, err := s.encryptor.Encrypt(user.AccessToken)
		if err != nil {
			atokenErrChan <- err
			return
		}
		atokenChan <- aToken
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		defer close(rtokenChan)
		defer close(rtokenErrChan)
		rToken, err := s.encryptor.Encrypt(user.RefreshToken)
		if err != nil {
			rtokenErrChan <- err
			return
		}
		rtokenChan <- rToken
	}()

	wg.Wait()

	select {
	case err := <-atokenErrChan:
		return user, err
	case err := <-rtokenErrChan:
		return user, err
	case <-ctx.Done():
		return user, _ErrOperationTimeout
	default:
		if _, err := s.adapter.UpsertUser(ctx, domain.UserAccess{
			ID:           user.ID,
			AccessToken:  <-atokenChan,
			RefreshToken: <-rtokenChan,
			TokenType:    user.TokenType,
			Scope:        user.Scope,
			ExpiresAt:    user.ExpiresAt,
		}); err != nil {
			return user, err
		}
		return user, nil
	}
}

func (s userService) DeleteUser(ctx context.Context, uid string) error {
	id := strings.TrimSpace(uid)
	s.logger.Debugf("validating uid %s to perform a delete action", id)

	if id == "" {
		return &InvalidServiceParameterError{
			Name:   "UID",
			Reason: "Should not be blank",
		}
	}

	s.logger.Debugf("uid %s is valid to perform a delete action", id)
	return s.adapter.DeleteUserByID(ctx, uid)
}
