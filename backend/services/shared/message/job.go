package message

import "encoding/json"

type JobMessage struct {
	UID      string `json:"uid"`
	Filename string `json:"filename"`
	Url      string `json:"url"`
}

func (s JobMessage) ToJSON() []byte {
	buf, _ := json.Marshal(s)
	return buf
}
