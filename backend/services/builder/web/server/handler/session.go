package handler

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"time"

	plog "github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/builder/web/server/core/domain"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/builder/web/server/core/port"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/request"
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

func (s SessionHandler) GetSession(ctx context.Context, mid *string, response *domain.Session) error {
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

	if sess, ok := session.(domain.Session); ok {
		s.logger.Debugf("session for mid %s has been found", *mid)
		*response = sess
		return nil
	}

	return err
}

func (s SessionHandler) OwnerRemoveSession(ctx context.Context, request *request.OwnerRemoveSessionRequest, reponse *interface{}) error {
	dctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	if request.Mid == "" {
		return nil
	}

	md := md5.Sum([]byte(request.Mid))
	mid := hex.EncodeToString(md[:])

	sess, err := s.service.GetSession(dctx, mid)
	if sess.Owner == request.Uid {
		return s.service.DeleteSession(dctx, mid)
	}

	return err
}

func (s SessionHandler) RemoveSession(ctx context.Context, mid *string, response *interface{}) error {
	*mid = strings.TrimSpace(*mid)

	if *mid == "" {
		s.logger.Error("invalid mid to fetch a session")
		return ErrEmptyIdValue
	}

	s.logger.Debugf("trying to remove a session for mid %s", *mid)
	dctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()
	return s.service.DeleteSession(dctx, *mid)
}
