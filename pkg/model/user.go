package model

import (
	"encoding/json"
	"errors"
	"net/http"
)

type UserRequest struct {
	UserID	  int    `json:"userID"`
	UserFName string `json:"userfname"`
	UserLName string `json:"userlname"`
	Civilite  string `json:"civilite"`
	Address   string `json:"address"`
	City      string `json:"city"`
	Country   string `json:"country"`
	UserEmail string `json:"useremail"`
	UserPW    string `json:"userpw"`
	Pseudo    string `json:"pseudo"`
	BirthDate string `json:"birthdate"`
}

func (a *UserRequest) Bind(r *http.Request) error {
	if a.UserFName == "" {
		return errors.New("missing required UserFName fields")
	}
	if a.UserLName == "" {
		return errors.New("missing required UserLName fields")
	}
	if a.UserEmail == "" {
		return errors.New("missing required UserEmail fields")
	}
	if a.UserPW == "" {
		return errors.New("missing required UserPW fields")
	}
	if a.Pseudo == "" {
		return errors.New("missing required Pseudo fields")
	}
	if a.BirthDate == "" {
		return errors.New("missing required BirthDate fields")
	}
	if a.Civilite == "" {
		return errors.New("missing required Civilite fields")
	}
	if a.Address == "" {
		return errors.New("missing required Address fields")
	}
	if a.City == "" {
		return errors.New("missing required City fields")
	}
	if a.Country == "" {
		return errors.New("missing required Country fields")
	}
	return nil
}

type LoginRequest struct {
	UserEmail string `json:"useremail"`
	Pseudo    string `json:"pseudo"`
	UserPW    string `json:"userpw"`
}

func (b *LoginRequest) Bind(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(b)
	if err != nil {
		return errors.New("failed to decode JSON")
	}

	if b.UserEmail == "" && b.Pseudo == "" {
		return errors.New("missing required fields ( UserEmail or Pseudo )")
	}
	if b.UserPW == "" {
		return errors.New("missing required UserPW fields")
	}
	return nil
}

type UserDeleteRequest struct {
	UserID int `json:"userID"`
}

func (a *UserDeleteRequest) Bind(r *http.Request) error {
	if a.UserID == 0 {
		return errors.New("missing required UserID fields")
	}
	return nil
}

type UserPasswordRequest struct {
	UserID    int    `json:"userID"`
}
type UserResponse struct {
	UserFName string   `json:"userfname"`
	UserLName string   `json:"userlname"`
	UserEmail string   `json:"useremail"`
	BirthDate string   `json:"birthdate"`
	Article   []string `json:"articles"`
}