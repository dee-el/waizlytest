package mockrepositories

import (
	"context"

	"github.com/stretchr/testify/mock"

	"waizlytest/repositories"
)

var _ (repositories.UserWriter) = (*MockUserRepository)(nil)
var _ (repositories.UserReader) = (*MockUserRepository)(nil)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, u *repositories.User) (int64, error) {
	args := m.Called(ctx, u)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, u *repositories.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepository) CreateUserAttendance(ctx context.Context, ua *repositories.UserAttendance) error {
	args := m.Called(ctx, ua)
	return args.Error(0)
}

func (m *MockUserRepository) SaveUserAttendanceSummary(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRepository) FindUserByPhone(ctx context.Context, phone string) (*repositories.User, error) {
	args := m.Called(ctx, phone)
	return args.Get(0).(*repositories.User), args.Error(1)
}

func (m *MockUserRepository) FindUserByID(ctx context.Context, id int64) (*repositories.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*repositories.User), args.Error(1)
}
