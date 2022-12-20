package service

import (
	"context"
	"strings"

	plog "github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/builder/web/server/core/domain"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/builder/web/server/core/port"
)

type sessionService struct {
	logger  plog.Logger
	adapter port.SessionServiceAdapter
}

func NewSessionService(
	logger plog.Logger,
	adapter port.SessionServiceAdapter,
) port.SessionService {
	return sessionService{
		logger:  logger,
		adapter: adapter,
	}
}

func (s sessionService) CreateSession(ctx context.Context, mid string, session domain.Session) (domain.Session, error) {
	mid = strings.TrimSpace(mid)
	s.logger.Debugf("validating mid %s to create a new session", mid)
	if mid == "" {
		return session, &InvalidServiceParameterError{
			Name:   "MeetingID",
			Reason: "Should not be blank",
		}
	}

	s.logger.Debugf("validating a session %s intance to create a new session", session.DocKey)
	if err := session.Validate(); err != nil {
		return session, err
	}

	s.logger.Debugf("session %s is valid", session.DocKey)
	return s.adapter.InsertSession(ctx, mid, session)
}

func (s sessionService) GetSession(ctx context.Context, mid string) (domain.Session, error) {
	mid = strings.TrimSpace(mid)
	s.logger.Debugf("validating mid %s to get an existing session", mid)
	if mid == "" {
		return domain.Session{}, &InvalidServiceParameterError{
			Name:   "MeetingID",
			Reason: "Should not be blank",
		}
	}

	s.logger.Debugf("mid %s is valid", mid)
	session, err := s.adapter.SelectSessionByMettingID(ctx, mid)
	if err != nil {
		return session, err
	}

	s.logger.Debugf("found a session: %v", session)
	return session, nil
}

func (s sessionService) DeleteSession(ctx context.Context, mid string) error {
	mid = strings.TrimSpace(mid)
	s.logger.Debugf("validating mid %s to delete a session", mid)
	if mid == "" {
		return &InvalidServiceParameterError{
			Name:   "MeetingID",
			Reason: "Should not be blank",
		}
	}

	s.logger.Debugf("mid %s is valid", mid)
	return s.adapter.DeleteSessionByMeetingID(ctx, mid)
}
