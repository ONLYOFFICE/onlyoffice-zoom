package model

import "encoding/json"

type ZoomMessage struct {
	ID        string     `json:"id"`
	Message   string     `json:"message"`
	Timestamp int        `json:"timestamp"`
	Files     []ZoomFile `json:"files"`
}

func (m ZoomMessage) ToJSON() []byte {
	buf, _ := json.Marshal(m)
	return buf
}

type ZoomFile struct {
	DownloadURL string `json:"download_url"`
	FileID      string `json:"file_id"`
	Filename    string `json:"file_name"`
	Size        int    `json:"file_size"`
	Timestamp   int    `json:"timestamp,omitempty"`
}

type ZoomFileMessage struct {
	Page     int        `json:"page_size"`
	Next     string     `json:"next_page_token"`
	Messages []ZoomFile `json:"messages,omitempty"`
}

func (fm ZoomFileMessage) ToJSON() []byte {
	buf, _ := json.Marshal(fm)
	return buf
}
