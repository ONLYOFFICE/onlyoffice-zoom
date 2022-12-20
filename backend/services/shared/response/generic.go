package response

import "encoding/json"

type GenericReponse struct {
	Error  int    `json:"error"`
	Reason string `json:"reason,omitempty"`
}

func (r GenericReponse) ToJSON() []byte {
	buf, _ := json.Marshal(r)
	return buf
}
