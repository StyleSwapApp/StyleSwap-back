package model

import (
	"errors"
	"net/http"
)

type MessageRequest struct {
	Content    string `json:"content"`
}

func (m *MessageRequest) Bind(r *http.Request) error {
	if m.Content == "" {
		return errors.New("missing required Content field")
	}
	return nil
}