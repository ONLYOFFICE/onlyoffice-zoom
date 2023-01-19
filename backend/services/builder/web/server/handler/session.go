package handler

import (
	"context"
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

func (s SessionHandler) GetSessionOwner(ctx context.Context, mid *string, response *string) error {
	*mid = strings.TrimSpace(*mid)

	if *mid == "" {
		s.logger.Error("invalid mid to fetch a session")
		return ErrEmptyIdValue
	}

	session, err, _ := s.group.Do(*mid, func() (interface{}, error) {
		s.logger.Debugf("trying to find a session for mid %s", *mid)
		session, err := s.service.GetSession(ctx, *mid)
		if err != nil {
			s.logger.Errorf("could not find any session for mid %s. Error: %s", *mid, err.Error())
			return nil, err
		}
		return session, nil
	})

	if err != nil {
		return err
	}

	if sess, ok := session.(domain.Session); ok {
		s.logger.Debugf("session for mid %s has been found", *mid)
		*response = sess.Owner
		return nil
	}

	return nil
}
