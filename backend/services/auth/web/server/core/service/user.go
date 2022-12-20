package service

import (
	"context"
	"errors"
	"strings"

	plog "github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/web/server/core/domain"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/web/server/core/port"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/crypto"
)

var _ErrOperationTimeout = errors.New("operation timeout")

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
	errChan := make(chan error)
	doneChan := make(chan bool)

	go func() {
		s.logger.Debugf("validating user %s to perform a persist action", user.ID)
		if err := user.Validate(); err != nil {
			errChan <- err
			return
		}

		aToken, err := s.encryptor.Encrypt(user.AccessToken)
		if err != nil {
			errChan <- err
			return
		}

		rToken, err := s.encryptor.Encrypt(user.RefreshToken)
		if err != nil {
			errChan <- err
			return
		}

		s.logger.Debugf("user %s is valid. Persisting to database: %s", user.ID, user.AccessToken)
		if err := s.adapter.InsertUser(ctx, domain.UserAccess{
			ID:           user.ID,
			AccessToken:  aToken,
			RefreshToken: rToken,
			TokenType:    user.TokenType,
			Scope:        user.Scope,
			ExpiresAt:    user.ExpiresAt,
		}); err != nil {
			if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
				errChan <- err
			}
			return
		}

		doneChan <- true
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return _ErrOperationTimeout
	case <-doneChan:
		return nil
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

	errChan := make(chan error)
	doneChan := make(chan domain.UserAccess)

	go func() {
		user, err := s.adapter.SelectUserByID(ctx, id)
		if err != nil {
			if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
				errChan <- err
			}
			return
		}

		aToken, err := s.encryptor.Decrypt(user.AccessToken)
		if err != nil {
			errChan <- err
			return
		}

		rToken, err := s.encryptor.Decrypt(user.RefreshToken)
		if err != nil {
			errChan <- err
			return
		}

		s.logger.Debugf("found a user: %v", user)
		doneChan <- domain.UserAccess{
			ID:           user.ID,
			AccessToken:  aToken,
			RefreshToken: rToken,
			TokenType:    user.TokenType,
			Scope:        user.Scope,
			ExpiresAt:    user.ExpiresAt,
		}
	}()

	select {
	case err := <-errChan:
		return domain.UserAccess{}, err
	case <-ctx.Done():
		return domain.UserAccess{}, _ErrOperationTimeout
	case usr := <-doneChan:
		return usr, nil
	}
}

func (s userService) UpdateUser(ctx context.Context, user domain.UserAccess) (domain.UserAccess, error) {
	s.logger.Debugf("validating user %s to perform an update action", user.ID)
	if err := user.Validate(); err != nil {
		return domain.UserAccess{}, err
	}

	errChan := make(chan error)
	doneChan := make(chan bool)

	go func() {
		aToken, err := s.encryptor.Encrypt(user.AccessToken)
		if err != nil {
			errChan <- err
			return
		}

		rToken, err := s.encryptor.Encrypt(user.RefreshToken)
		if err != nil {
			errChan <- err
			return
		}

		s.logger.Debugf("user %s is valid to perform an update action", user.ID)
		if _, err := s.adapter.UpsertUser(ctx, domain.UserAccess{
			ID:           user.ID,
			AccessToken:  aToken,
			RefreshToken: rToken,
			TokenType:    user.TokenType,
			Scope:        user.Scope,
			ExpiresAt:    user.ExpiresAt,
		}); err != nil {
			if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
				errChan <- err
			}
			return
		}

		doneChan <- true
	}()

	select {
	case err := <-errChan:
		return user, err
	case <-ctx.Done():
		return user, _ErrOperationTimeout
	case <-doneChan:
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
