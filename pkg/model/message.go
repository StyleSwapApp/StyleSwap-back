package model

import (
	"errors"
	"net/http"
)

// type AuthedRequest struct {
// 	ClientID string `json:"clientID"`
// }

// func (a *AuthedRequest) Bind(r *http.Request) error {
// 	if a.ClientID == "" {
// 		return errors.New("missing required UserID field")
// 	}
// 	return nil
// }

type MessageRequest struct {
	UserID string `json:"Destinataire"`
	Content    string `json:"content"`
}

func (m *MessageRequest) Bind(r *http.Request) error {
	if m.Content == "" {
		return errors.New("missing required Content field")
	}
	return nil
}