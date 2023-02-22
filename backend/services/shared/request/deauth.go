package request

import (
	"encoding/json"
)

type DeauthorizationPayload struct {
	Uid                 string `json:"user_id" mapstructure:"user_id"`
	Aid                 string `json:"account_id" mapstructure:"account_id"`
	Cid                 string `json:"client_id" mapstructure:"client_id"`
	DeauthorizationTime string `json:"deauthorization_time" mapstructure:"deauthorization_time"`
	Signature           string `json:"signature" mapstructure:"signature"`
}

func (p DeauthorizationPayload) ToJSON() []byte {
	buf, _ := json.Marshal(p)
	return buf
}
