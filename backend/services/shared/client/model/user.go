package model

import (
	"strings"
)

type User struct {
	ID        string `json:"id"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	Email     string `json:"email"`
	Language  string `json:"language"`
}

func (u *User) Validate() error {
	u.ID = strings.TrimSpace(u.ID)
	u.Firstname = strings.TrimSpace(u.Firstname)
	u.Lastname = strings.TrimSpace(u.Lastname)
	u.Email = strings.TrimSpace(u.Email)
	u.Language = strings.TrimSpace(u.Language)

	if u.ID == "" {
		return ErrInvalidTokenFormat
	}

	if u.Firstname == "" {
		return ErrInvalidTokenFormat
	}

	if u.Language == "" {
		u.Language = "en"
	}

	return nil
}
