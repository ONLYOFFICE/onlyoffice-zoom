package ws

import (
	"fmt"

	"github.com/ONLYOFFICE/zoom-onlyoffice/services/gateway/web/server/middleware/security"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/wsmessage"
	"github.com/go-chi/chi/v5"
	"github.com/olahol/melody"
	"go-micro.dev/v4/client"
)

func NewOnConnectHandler(namespace string, m *melody.Melody, client client.Client) func(s *melody.Session) {
	return func(s *melody.Session) {
		zctx, ok := s.Request.Context().Value(security.ZoomContext{}).(security.ZoomContext)
		if !ok || zctx.Mid == "" {
			s.Close()
			return
		}

		mid := chi.URLParam(s.Request, "mid")
		ss, err := m.Sessions()
		if err != nil || mid != zctx.Mid {
			s.Close()
			return
		}

		for _, o := range ss {
			value, exists := o.Get(mid)
			if exists {
				if sm, ok := value.(wsmessage.SessionMessage); ok {
					s.Set(mid, sm)
					s.Write(sm.ToJSON())
					return
				}
			}
		}

		req := client.NewRequest(fmt.Sprintf("%s:builder", namespace), "SessionHandler.GetSessionOwner", mid)
		var resp interface{}
		if err := client.Call(s.Request.Context(), req, &resp); err == nil {
			msg := wsmessage.SessionMessage{
				MID:       mid,
				InSession: true,
			}
			s.Set(mid, msg)
			s.Write(msg.ToJSON())
			return
		}

		msg := wsmessage.SessionMessage{
			MID:       mid,
			InSession: false,
		}
		s.Set(mid, msg)
		s.Write(msg.ToJSON())
	}
}
