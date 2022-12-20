package request

import "encoding/json"

type EventRequest struct {
	Event string `json:"event"`
	Ts    int    `json:"event_ts"`
}

func (e EventRequest) ToJSON() []byte {
	buf, _ := json.Marshal(e)
	return buf
}
