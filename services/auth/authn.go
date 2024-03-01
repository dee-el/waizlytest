package authservice

import (
	"context"
)

type LoginRequest struct {
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

type LoginResponse struct {
	ID    int64  `json:"-"`
	Token string `json:"token"`
}

type Authn interface {
	Login(ctx context.Context, params LoginRequest) (LoginResponse, error)
}
