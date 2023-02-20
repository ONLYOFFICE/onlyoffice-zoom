package message

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"time"

	plog "github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/builder/web/server/core/port"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/request"
	"github.com/mitchellh/mapstructure"
)

type OwnerRemoveSessionMessageHandler struct {
	logger  plog.Logger
	service port.SessionService
}

func BuildOwnerRemoveSessionMessageHandler(logger plog.Logger, service port.SessionService) OwnerRemoveSessionMessageHandler {
	return OwnerRemoveSessionMessageHandler{
		logger:  logger,
		service: service,
	}
}

func (i OwnerRemoveSessionMessageHandler) GetHandler() func(context.Context, interface{}) error {
	return func(ctx context.Context, payload interface{}) error {
		tctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		var request request.OwnerRemoveSessionRequest

		if err := mapstructure.Decode(payload, &request); err != nil {
			return _ErrInvalidHandlerPayload
		}

		if request.Mid == "" {
			return nil
		}

		md := md5.Sum([]byte(request.Mid))
		mid := hex.EncodeToString(md[:])

		sess, err := i.service.GetSession(tctx, mid)
		if err != nil {
			return nil
		}

		if sess.Owner == request.Uid {
			return i.service.DeleteSession(tctx, mid)
		}

		return nil
	}
}
