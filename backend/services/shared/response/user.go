package response

import "encoding/json"

type UserResponse struct {
	ID           string `json:"id" mapstructure:"id"`
	AccessToken  string `json:"access_token" mapstructure:"access_token"`
	RefreshToken string `json:"refresh_token" mapstructure:"refresh_token"`
	TokenType    string `json:"token_type" mapstructure:"token_type"`
	Scope        string `json:"scope" mapstructure:"scope"`
	ExpiresAt    int64  `json:"expires_at" mapstructure:"expires_at"`
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
