package api

import "fmt"

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Id       string `json:"id"`
	Password string `json:"password" validate:"required"`
}

func (l LoginRequest) Validate() error {
	if (l.Username == "" && l.Id == "" && l.Email == "") || l.Password == "" {
		return fmt.Errorf("empty username or password")
	}

	return nil
}

type LoginResponse struct {
	Token string `json:"token"`
	Iss   string `json:"iss"`
}
