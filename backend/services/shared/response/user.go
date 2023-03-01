package response

import "encoding/json"

type UserResponse struct {
	ID           string `json:"id" mapstructure:"ID"`
	AccessToken  string `json:"access_token" mapstructure:"AccessToken"`
	RefreshToken string `json:"refresh_token" mapstructure:"RefreshToken"`
	TokenType    string `json:"token_type" mapstructure:"TokenType"`
	Scope        string `json:"scope" mapstructure:"Scope"`
	ExpiresAt    int64  `json:"expires_at" mapstructure:"ExpiresAt"`
}

func (ur UserResponse) ToJSON() []byte {
	buf, _ := json.Marshal(ur)
	return buf
}

type UserTokenResponse struct {
	ID          string `json:"id"`
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

func (ut UserTokenResponse) ToJSON() []byte {
	buf, _ := json.Marshal(ut)
	return buf
}
