package request

import "encoding/json"

type OwnerRemoveSessionRequest struct {
	Uid string `json:"uid"`
	Mid string `json:"mid"`
}

func (r OwnerRemoveSessionRequest) ToJSON() []byte {
	buf, _ := json.Marshal(r)
	return buf
}
