package ws

import (
	"encoding/json"
	"strings"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/shared/wsmessage"
	"github.com/olahol/melody"
	"go-micro.dev/v4/broker"
)

func NewNotifyOnMessage(logger log.Logger, m *melody.Melody) func(p broker.Event) error {
	return func(p broker.Event) error {
		var msg wsmessage.SessionMessage
		if err := json.Unmarshal(p.Message().Body, &msg); err != nil {
			logger.Errorf("could not notify ws clients. Reason: %s", err.Error())
			return err
		}

		ss, err := m.Sessions()
		if err != nil {
			logger.Fatalf("could not get ws sessions. Reason: %s", err.Error())
			return err
		}

		for _, o := range ss {
			if strings.Contains(o.Request.URL.Path, msg.MID) {
				o.Set(msg.MID, msg)
				o.Write(msg.ToJSON())
			}
		}

		return nil
	}
}
