package message

import "encoding/json"

type SessionMessage struct {
	MID       string `json:"mid"`
	InSession bool   `json:"in_session"`
}

func (s SessionMessage) ToJSON() []byte {
	buf, _ := json.Marshal(s)
	return buf
}
