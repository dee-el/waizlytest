package pgrepositories

import (
	"context"
	"database/sql"

	"waizlytest/repositories"
)

var _ (repositories.UserWriter) = (*UserRepository)(nil)
var _ (repositories.UserReader) = (*UserRepository)(nil)

// UserRepository implementation both of `repositories.UserWriter` & `repositories.UserReader`
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, u *repositories.User) (int64, error) {
	query := `
		INSERT INTO users (fullname, password, phone, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		u.FullName,
		u.Password,
		u.Phone,
		u.CreatedAt,
		u.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UserRepository) FindUserByPhone(ctx context.Context, phone string) (*repositories.User, error) {
	query := `
		SELECT id, fullname, password, phone, created_at, updated_at, deleted_at
		FROM users
		WHERE phone = $1
		LIMIT 1
	`

	var user repositories.User
	err := r.db.QueryRowContext(ctx, query, phone).Scan(
		&user.ID,
		&user.FullName,
		&user.Password,
		&user.Phone,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindUserByID(ctx context.Context, id int64) (*repositories.User, error) {
	query := `
		SELECT id, fullname, password,  phone, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1
		LIMIT 1
	`

	var user repositories.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.FullName,
		&user.Password,
		&user.Phone,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) CreateUserAttendance(ctx context.Context, ua *repositories.UserAttendance) error {
	query := `
		INSERT INTO user_attendance_logs (user_id, login_at)
		VALUES ($1, $2)
	`

	_, err := r.db.ExecContext(ctx, query, ua.UserID, ua.LoginAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, u *repositories.User) error {
	query := `
		UPDATE users
		SET fullname = $2, phone = $3, updated_at = $4, deleted_at = $5
		WHERE id = $1
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		u.ID,
		u.FullName,
		u.Phone,
		u.UpdatedAt,
		u.DeletedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) SaveUserAttendanceSummary(ctx context.Context, userID int64) error {
	query := `
		INSERT INTO user_attendance_summaries (user_id, total_login)
		VALUES ($1, 1)
		ON CONFLICT (user_id)
		DO UPDATE SET total_login = user_attendance_summaries.total_login + EXCLUDED.total_login
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}
