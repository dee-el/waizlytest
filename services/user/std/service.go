package stduserservice

import (
	"waizlytest/repositories"
	userservice "waizlytest/services/user"
)

var _ userservice.Registrator = (*StdService)(nil)
var _ userservice.Me = (*StdService)(nil)

type StdService struct {
	userReader repositories.UserReader
	userWriter repositories.UserWriter
}

func NewService(
	userWriter repositories.UserWriter,
	userReader repositories.UserReader,
) *StdService {
	return &StdService{
		userReader: userReader,
		userWriter: userWriter,
	}
}
