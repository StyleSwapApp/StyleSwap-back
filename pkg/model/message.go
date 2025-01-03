package model

import (
	"errors"
	"net/http"
)

type AuthedRequest struct {
	UserID string `json:"userID"`
}

func (a *AuthedRequest) Bind(r *http.Request) error {
	if a.UserID == "" {
		return errors.New("missing required UserID field")
	}
	return nil
}

type MessageRequest struct {
	ReceiverID string `json:"receiverID"`
	Content    string `json:"content"`
}

func (a *MessageRequest) Bind(r *http.Request) error {
	if a.ReceiverID == "" || a.Content == "" {
		return errors.New("missing required MessageRequest fields")
	}
	return nil
}

type MessageResponse struct {
	SenderID   string `json:"senderID"`
	ReceiverID string `json:"receiverID"`
	Content    string `json:"content"`
}
