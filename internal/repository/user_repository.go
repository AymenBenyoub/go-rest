package repository

import (
	"context"
	"database/sql"
	"errors"

	"rest/internal/db"

	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByID(ctx context.Context,pubID string) (*db.User, error) {
	var user db.User

	query := "SELECT public_id, username, email FROM users WHERE public_id = ?"
	err := r.db.QueryRowContext(ctx,query, pubID).Scan(&user.PublicID, &user.Username, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// no password crypto for now...
func (r *UserRepository) CreateUser(ctx context.Context,user *db.User) error {
	userUUID := uuid.New().String()

	query := "INSERT INTO users (public_id, username, email, password) VALUES (?, ?, ?, ?)"
	_, err := r.db.ExecContext(ctx,query, userUUID, user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}

	user.PublicID = userUUID // Assign generated UUID back to user
	return nil
}

func (r *UserRepository) UpdateUsername(ctx context.Context,pubID, newUsername string) error {
	query := "UPDATE users SET username = ? WHERE public_id = ?"
	res, err := r.db.ExecContext(ctx,query, newUsername, pubID)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// no password crypto for now...
func (r *UserRepository) UpdatePassword(ctx context.Context,pubID, newPassword string) error {
	query := "UPDATE users SET password = ? WHERE public_id = ?"
	res, err := r.db.ExecContext(ctx,query, newPassword, pubID)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context,pubID string) error {
	query := "DELETE FROM users WHERE public_id = ?"
	res, err := r.db.ExecContext(ctx,query, pubID)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
// protected endpoint.
func (r *UserRepository) GetAllUsers(ctx context.Context,) ([]db.User, error) {
	query := "SELECT public_id, username, email FROM users"
	rows, err := r.db.QueryContext(ctx,query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []db.User
	for rows.Next() {
		var user db.User
		if err := rows.Scan( &user.PublicID, &user.Username, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
