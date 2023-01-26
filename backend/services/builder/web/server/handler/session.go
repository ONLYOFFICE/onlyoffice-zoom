package handler

import (
	"context"
	"fmt"
	"strings"

	plog "github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/builder/web/server/core/domain"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/builder/web/server/core/port"
	"golang.org/x/sync/singleflight"
)

type SessionHandler struct {
	logger  plog.Logger
	service port.SessionService
	group   singleflight.Group
}

func NewSessionHandler(
	logger plog.Logger,
	service port.SessionService,
) SessionHandler {
	return SessionHandler{
		logger:  logger,
		service: service,
	}
}

func (s SessionHandler) getSession(ctx context.Context, mid string) (domain.Session, error) {
	mid = strings.TrimSpace(mid)

	if mid == "" {
		s.logger.Error("invalid mid to fetch a session")
		return domain.Session{}, ErrEmptyIdValue
	}

	session, err, _ := s.group.Do(mid, func() (interface{}, error) {
		s.logger.Debugf("trying to find a session for mid %s", mid)
		session, err := s.service.GetSession(ctx, mid)
		if err != nil {
			s.logger.Debugf("could not find any session for mid %s. Error: %s", mid, err.Error())
			return nil, err
		}
		return session, nil
	})

	if err != nil {
		return domain.Session{}, err
	}

	if sess, ok := session.(domain.Session); ok {
		s.logger.Debugf("session for mid %s has been found", mid)
		return sess, nil
	}

	return domain.Session{}, fmt.Errorf("could not find session with id %s", mid)
}

func (s SessionHandler) GetRealSession(ctx context.Context, mid *string, response *bool) error {
	sess, err := s.getSession(ctx, *mid)
	if err != nil {
		return err
	}
	*response = sess.Initial
	return nil
}

func (s SessionHandler) GetSessionOwner(ctx context.Context, mid *string, response *string) error {
	session, err := s.getSession(ctx, *mid)
	if err != nil {
		return err
	}

	*response = session.Owner
	return nil
}

func (s SessionHandler) RefreshSession(ctx context.Context, mid *string, response *bool) error {
	*mid = strings.TrimSpace(*mid)

	if *mid == "" {
		s.logger.Error("invalid mid to fetch a session")
		return ErrEmptyIdValue
	}

	if _, err, _ := s.group.Do(*mid, func() (interface{}, error) {
		s.logger.Debugf("trying to find a session for mid %s", *mid)
		session, err := s.service.GetSession(ctx, *mid)
		if err != nil {
			s.logger.Errorf("could not find any session for mid %s to refresh it. Error: %s", *mid, err.Error())
			return nil, err
		}

		if session.Initial {
			s.logger.Debugf("refreshing initial session with key %s", session.DocKey)
			_, err = s.service.CreateSession(ctx, *mid, domain.Session{
				Owner:    session.Owner,
				Filename: session.Filename,
				FileURL:  session.FileURL,
				DocKey:   session.DocKey,
			})

			return nil, err
		}

		return nil, nil
	}); err != nil {
		return err
	}

	return nil
}
