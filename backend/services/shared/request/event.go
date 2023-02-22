package request

import "encoding/json"

type EventRequest struct {
	Event   string                 `json:"event"`
	Ts      int                    `json:"event_ts"`
	Payload map[string]interface{} `json:"payload"`
}

func (e EventRequest) ToJSON() []byte {
	buf, _ := json.Marshal(e)
	return buf
}

type WebhookEventPayload struct {
	PlainToken string `json:"plainToken" mapstructure:"plainToken"`
}

func (e WebhookEventPayload) ToJSON() []byte {
	buf, _ := json.Marshal(e)
	return buf
}
