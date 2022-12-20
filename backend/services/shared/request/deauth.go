package request

import (
	"encoding/json"
	"time"
)

type DeauthorizationEventRequest struct {
	EventRequest
	Payload DeauthorizationPayload `json:"payload"`
}

func (dr DeauthorizationEventRequest) ToJSON() []byte {
	buf, _ := json.Marshal(dr)
	return buf
}

type DeauthorizationPayload struct {
	Uid                 string    `json:"user_id"`
	Aid                 string    `json:"account_id"`
	Cid                 string    `json:"client_id"`
	DeauthorizationTime time.Time `json:"deauthorization_time"`
	Signature           string    `json:"signature"`
}

func (p DeauthorizationPayload) ToJSON() []byte {
	buf, _ := json.Marshal(p)
	return buf
}
