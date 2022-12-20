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
	Comment                 bool `json:"comment,omitempty"`
	Copy                    bool `json:"copy,omitempty"`
	DeleteCommentAuthorOnly bool `json:"deleteCommentAuthorOnly,omitempty"`
	Download                bool `json:"download,omitempty"`
	Edit                    bool `json:"edit"`
	EditCommentAuthorOnly   bool `json:"editCommentAuthorOnly,omitempty"`
	FillForms               bool `json:"fillForms,omitempty"`
	ModifyContentControl    bool `json:"modifyContentControl,omitempty"`
	ModifyFilter            bool `json:"modifyFilter,omitempty"`
	Print                   bool `json:"print,omitempty"`
	Review                  bool `json:"review,omitempty"`
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
	Goback Goback `json:"goback"`
}

type Goback struct {
	RequestClose bool `json:"requestClose"`
}
