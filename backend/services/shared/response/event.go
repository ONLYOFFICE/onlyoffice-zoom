package response

import "encoding/json"

type WebhookEventResponse struct {
	PlainToken     string `json:"plainToken"`
	EncryptedToken string `json:"encryptedToken"`
}

func (e WebhookEventResponse) ToJSON() []byte {
	buf, _ := json.Marshal(e)
	return buf
}
