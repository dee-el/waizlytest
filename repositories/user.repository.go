package repositories

import "context"

type UserWriter interface {
	CreateUser(ctx context.Context, u *User) (int64, error)
	UpdateUser(ctx context.Context, u *User) error

	CreateUserAttendance(ctx context.Context, ua *UserAttendance) error
	SaveUserAttendanceSummary(ctx context.Context, userID int64) error
}

type UserReader interface {
	FindUserByPhone(ctx context.Context, phone string) (*User, error)
	FindUserByID(ctx context.Context, id int64) (*User, error)
}
