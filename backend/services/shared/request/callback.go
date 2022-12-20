package request

import (
	"encoding/json"
	"fmt"
	"strings"
)

type MissingRequestFieldsError struct {
	Request string
	Field   string
	Reason  string
}

func (e *MissingRequestFieldsError) Error() string {
	return fmt.Sprintf("missing %s's field %s. Reason: %s", e.Request, e.Field, e.Reason)
}

type CallbackRequest struct {
	Actions []struct {
		Type   int    `json:"type"`
		UserID string `json:"userid"`
	} `json:"actions"`
	Key    string   `json:"key"`
	Status int      `json:"status"`
	Users  []string `json:"users"`
	URL    string   `json:"url"`
	Token  string   `json:"token"`
}

func (cr CallbackRequest) ToJSON() []byte {
	buf, _ := json.Marshal(cr)
	return buf
}

func (c *CallbackRequest) Validate() error {
	c.Key = strings.TrimSpace(c.Key)
	c.Token = strings.TrimSpace(c.Token)

	if c.Key == "" {
		return &MissingRequestFieldsError{
			Request: "Callback",
			Field:   "Key",
			Reason:  "Should not be empty",
		}
	}

	if c.Token == "" {
		return &MissingRequestFieldsError{
			Request: "Callback",
			Field:   "Token",
			Reason:  "Should not be empty",
		}
	}

	if c.Status <= 0 || c.Status > 7 {
		return &MissingRequestFieldsError{
			Request: "Callback",
			Field:   "Status",
			Reason:  "Invalid status. Exptected 0 < status <= 7",
		}
	}

	return nil
}
