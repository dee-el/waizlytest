package stduserservice

import (
	"context"
	"fmt"
	"time"

	commonerr "waizlytest/common/errors"

	"waizlytest/repositories"
	userservice "waizlytest/services/user"
)

func (s *StdService) Register(ctx context.Context, params userservice.RegisterRequest) (int64, error) {
	usr, err := s.userReader.FindUserByPhone(ctx, params.Phone)
	if err != nil {
		return 0, err
	}

	if usr != nil {
		e := commonerr.Conflicted("")
		e.AddField("phone", fmt.Sprintf("[%s] already registered", params.Phone))
		return 0, e
	}

	hashed, err := hashed(params.Password)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	return s.userWriter.CreateUser(ctx, &repositories.User{
		FullName:  params.FullName,
		Password:  hashed,
		Phone:     params.Phone,
		CreatedAt: now,
		UpdatedAt: now,
	})
}
