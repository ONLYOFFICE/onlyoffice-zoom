package request

import "encoding/json"

type BuildConfigRequest struct {
	Uid       string `json:"uid"`
	Mid       string `json:"mid"`
	UserAgent string `json:"user_agent"`
	Filename  string `json:"file_name"`
	FileURL   string `json:"file_url"`
	Language  string `json:"language"`
}

func (c BuildConfigRequest) ToJSON() []byte {
	buf, _ := json.Marshal(c)
	return buf
}
