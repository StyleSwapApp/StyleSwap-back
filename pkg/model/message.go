package model

import (
	"errors"
	"net/http"
)

type AuthedRequest struct {
	UserID string `json:"userID"`
	ClientID string `json:"clientID"`
}

func (a *AuthedRequest) Bind(r *http.Request) error {
	if a.UserID == "" || a.ClientID == "" {
		return errors.New("missing required UserID field")
	}
	return nil
}

type MessageRequest struct {
	UserID string `json:"userID"`
	Content    string `json:"content"`
}

// func (a *MessageRequest) Bind(r *http.Request) error {
// 	if a.Content == "" {
// 		return errors.New("missing required MessageRequest fields")
// 	}
// 	return nil
// }