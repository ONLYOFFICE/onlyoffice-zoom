package domain

import (
	"net/url"
	"strings"
)

type Session struct {
	Owner    string `json:"owner_id"`
	Filename string `json:"file_name"`
	FileURL  string `json:"file_url"`
	DocKey   string `json:"doc_key"`
	Initial  bool   `json:"initial"`
}

func (u *Session) Validate() error {
	u.Owner = strings.TrimSpace(u.Owner)
	u.Filename = strings.TrimSpace(u.Filename)
	u.FileURL = strings.TrimSpace(u.FileURL)
	u.DocKey = strings.TrimSpace(u.DocKey)

	if u.Owner == "" {
		return &InvalidModelFieldError{
			Model:  "Session",
			Field:  "Owner id",
			Reason: "Should not be empty",
		}
	}

	if u.Filename == "" {
		return &InvalidModelFieldError{
			Model:  "Session",
			Field:  "File name",
			Reason: "Should not be empty",
		}
	}

	if _, err := url.ParseRequestURI(u.FileURL); err != nil {
		return &InvalidModelFieldError{
			Model:  "Session",
			Field:  "File URL",
			Reason: "Invalid URL format",
		}
	}

	if u.DocKey == "" {
		return &InvalidModelFieldError{
			Model:  "Session",
			Field:  "Doc key",
			Reason: "Should not be empty",
		}
	}

	return nil
}
