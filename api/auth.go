package api

import "fmt"

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (l LoginRequest) Validate() error {
	if l.Username == "" || l.Password == "" {
		return fmt.Errorf("empty username or password")
	}

	return nil
}

type LoginResponse struct {
	Token string `json:"token"`
	Iss   string `json:"iss"`
}
