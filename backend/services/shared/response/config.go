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
	Permissions Permissions `json:"permissions"`
	Title       string      `json:"title"`
	URL         string      `json:"url"`
}

type EditorConfig struct {
	CallbackURL   string        `json:"callbackUrl"`
	Customization Customization `json:"customization,omitempty"`
	Lang          string        `json:"lang,omitempty"`
	User          User          `json:"user"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Customization struct {
	Goback        Goback `json:"goback"`
	HideRightMenu bool   `json:"hideRightMenu"`
	Plugins       bool   `json:"plugins"`
}

type Goback struct {
	RequestClose bool `json:"requestClose"`
}
