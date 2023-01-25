package ws

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ONLYOFFICE/zoom-onlyoffice/services/gateway/web/server/middleware/security"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/wsmessage"
	"github.com/go-chi/chi/v5"
	"github.com/olahol/melody"
	"go-micro.dev/v4/client"
)

func NewOnConnectHandler(namespace, secret string, m *melody.Melody, client client.Client) func(s *melody.Session) {
	return func(s *melody.Session) {
		mid := chi.URLParam(s.Request, "mid")
		if mid == "" {
			s.Close()
			return
		}

		token := s.Request.URL.Query().Get("token")
		if token == "" {
			s.Close()
			return
		}

		zctx, err := security.ExtractZoomContext(token, secret)

		if err != nil {
			s.Close()
			return
		}

		now := int(time.Now().UnixMilli()) - 4*60*1000
		if zctx.Exp < now {
			s.Close()
			return
		}

		ss, err := m.Sessions()
		if err != nil {
			s.Close()
			return
		}

		md := md5.Sum([]byte(zctx.Mid))
		zmid := hex.EncodeToString(md[:])

		if zmid != mid {
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

		req := client.NewRequest(fmt.Sprintf("%s:builder", namespace), "SessionHandler.GetRealSession", mid)
		var resp bool
		if err := client.Call(s.Request.Context(), req, &resp); err == nil && resp {
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
