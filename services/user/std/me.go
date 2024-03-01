package stduserservice

import (
	"context"
	"fmt"

	commonerr "waizlytest/common/errors"
	userservice "waizlytest/services/user"
)

func (s *StdService) GetProfile(ctx context.Context, id int64) (userservice.Profile, error) {
	resp := userservice.Profile{}

	usr, err := s.userReader.FindUserByID(ctx, id)
	if err != nil {
		return resp, err
	}

	if usr == nil {
		e := commonerr.BadRequest("user not found")
		return resp, e
	}

	resp.FullName = usr.FullName
	resp.Phone = usr.Phone
	return resp, nil
}

func (s *StdService) UpdateProfile(ctx context.Context, id int64, params userservice.UpdateRequest) error {
	usr, err := s.userReader.FindUserByID(ctx, id)
	if err != nil {
		return err
	}

	if usr == nil {
		e := commonerr.BadRequest("user not found")
		return e
	}

	if params.Phone != nil {
		if usr.Phone != *params.Phone {
			pUsr, err := s.userReader.FindUserByPhone(ctx, *params.Phone)
			if err != nil {
				return err
			}

			if pUsr != nil {
				e := commonerr.Conflicted("")
				e.AddField("phone", fmt.Sprintf("[%s] already registered", *params.Phone))
				return e
			}
		}

		usr.Phone = *params.Phone
	}

	if params.FullName != nil && *params.FullName != "" {
		usr.FullName = *params.FullName
	}

	return s.userWriter.UpdateUser(ctx, usr)
}
