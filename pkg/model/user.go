package model

import (
	"errors"
	"net/http"
)

type UserRequest struct {
	UserFName 	string 	`json:"userfname"`
	UserLName  	string  `json:"userlname"`
	UserEmail 	string 	`json:"useremail"`
	UserPW 		string 	`json:"userpw"`
	Pseudo      string  `json:"pseudo"`
	BirthDate 	string 	`json:"birthdate"`
}

type LoginRequest struct {
	UserEmail 	string 	`json:"useremail"`
	Pseudo 		string  `json:"pseudo"`
	UserPW 		string 	`json:"userpw"`
}

func (a *LoginRequest) Bind(r *http.Request) error {
	if a.UserEmail == "" || a.Pseudo == "" {
		return errors.New("missing required fields ( UserEmail or Pseudo )")
	}	
	if a.UserPW == "" {
		return errors.New("missing required UserPW fields")
	}
	return nil
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
	return nil
}

type UserResponse struct {
	UserFName 	string 	`json:"userfname"`
	UserLName  	string  `json:"userlname"`
	UserEmail 	string 	`json:"useremail"`
	BirthDate 	string 	`json:"birthdate"`
	Article     []string `json:"articles"`
}