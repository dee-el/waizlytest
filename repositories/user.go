package repositories

import "time"

type User struct {
	ID        int64
	FullName  string
	Password  string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
type UserAttendance struct {
	UserID  int64
	LoginAt time.Time
}

type UserAttendanceSummary struct {
	UserID        int64
	TotalLoggedIn int
}
