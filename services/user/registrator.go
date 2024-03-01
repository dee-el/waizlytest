package userservice

import "context"

type RegisterRequest struct {
	FullName string `json:"fullname"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

type Registrator interface {
	Register(ctx context.Context, params RegisterRequest) (int64, error)
}
