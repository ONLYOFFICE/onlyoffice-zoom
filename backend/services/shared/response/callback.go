package response

import "encoding/json"

type CallbackResponse struct {
	Error int `json:"error"`
}

func (c CallbackResponse) ToJSON() []byte {
	buf, _ := json.Marshal(c)
	return buf
}
