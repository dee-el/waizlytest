package userservice

import "context"

type UpdateRequest struct {
	Phone    *string `json:"phone"`
	FullName *string `json:"fullname"`
}

type Profile struct {
	FullName string `json:"fullname"`
	Phone    string `json:"phone"`
}

type Me interface {
	GetProfile(ctx context.Context, id int64) (Profile, error)
	UpdateProfile(ctx context.Context, id int64, params UpdateRequest) error
}
