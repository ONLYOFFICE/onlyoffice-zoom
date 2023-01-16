package model

import (
	"strings"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
}

func (t *Token) Validate() error {
	t.AccessToken = strings.TrimSpace(t.AccessToken)
	t.RefreshToken = strings.TrimSpace(t.RefreshToken)
	t.TokenType = strings.TrimSpace(t.TokenType)
	t.Scope = strings.TrimSpace(t.Scope)

	if t.AccessToken == "" {
		return ErrInvalidTokenFormat
	}

	if t.RefreshToken == "" {
		return ErrInvalidTokenFormat
	}

	if t.TokenType == "" {
		return ErrInvalidTokenFormat
	}

	if t.Scope == "" {
		return ErrInvalidTokenFormat
	}

	if t.ExpiresIn < 1 {
		return ErrInvalidTokenFormat
	}

	return nil
}
