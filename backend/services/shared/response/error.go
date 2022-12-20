package response

import "encoding/json"

type MicroError struct {
	ID     string `json:"id"`
	Code   int    `json:"code"`
	Detail string `json:"detail"`
	Status string `json:"status"`
}

func (e MicroError) ToJSON() []byte {
	buf, _ := json.Marshal(e)
	return buf
}
