package request

import "encoding/json"

type OwnerRemoveSessionRequest struct {
	Uid string `json:"uid" mapstructure:"uid"`
	Mid string `json:"mid" mapstructure:"mid"`
}

func (r OwnerRemoveSessionRequest) ToJSON() []byte {
	buf, _ := json.Marshal(r)
	return buf
}
