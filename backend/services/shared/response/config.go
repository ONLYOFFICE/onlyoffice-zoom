package response

import (
	"encoding/json"

	"github.com/golang-jwt/jwt"
)

type BuildConfigResponse struct {
	jwt.StandardClaims
	Document     Document     `json:"document"`
	DocumentType string       `json:"documentType"`
	EditorConfig EditorConfig `json:"editorConfig"`
	Type         string       `json:"type"`
	Token        string       `json:"token,omitempty"`
	Session      bool         `json:"is_session,omitempty"`
	Owner        bool         `json:"is_owner,omitempty"`
}

func (r BuildConfigResponse) ToJSON() []byte {
	buf, _ := json.Marshal(r)
	return buf
}

type Permissions struct {
	Download                bool `json:"download"`
	Edit                    bool `json:"edit"`
	Print                   bool `json:"print"`
}

type Document struct {
	FileType    string      `json:"fileType"`
	Key         string      `json:"key"`
	Title       string      `json:"title"`
	URL         string      `json:"url"`
	Permissions Permissions `json:"permissions"`
}

type EditorConfig struct {
	User          User          `json:"user"`
	CallbackURL   string        `json:"callbackUrl"`
	Customization Customization `json:"customization,omitempty"`
	Lang          string        `json:"lang,omitempty"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Customization struct {
	Goback        Goback `json:"goback"`
	Plugins       bool   `json:"plugins"`
	HideRightMenu bool   `json:"hideRightMenu"`
}

type Goback struct {
	RequestClose bool `json:"requestClose"`
}
