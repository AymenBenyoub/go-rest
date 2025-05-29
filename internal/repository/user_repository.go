package repository

import (
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

func (r *UserRepository) GetUserByID(pubID string) (*db.User, error) {
	var user db.User

	query := "SELECT public_id, username, email FROM users WHERE public_id = ?"
	err := r.db.QueryRow(query, pubID).Scan(&user.PublicID, &user.Username, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
// no password crypto for now...
func (r *UserRepository) CreateUser(user *db.User) error {
	userUUID := uuid.New().String()

	query := "INSERT INTO users (public_id, username, email, password) VALUES (?, ?, ?, ?)"
	_, err := r.db.Exec(query, userUUID, user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}

	user.PublicID = userUUID // Assign generated UUID back to user
	return nil
}

func (r *UserRepository) UpdateUsername(pubID, newUsername string) error {
	query := "UPDATE users SET username = ? WHERE public_id = ?"
	res, err := r.db.Exec(query, newUsername, pubID)
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
func (r *UserRepository) UpdatePassword(pubID, newPassword string) error {
	query := "UPDATE users SET password = ? WHERE public_id = ?"
	res, err := r.db.Exec(query, newPassword, pubID)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) DeleteUser(pubID string) error {
	query := "DELETE FROM users WHERE public_id = ?"
	res, err := r.db.Exec(query, pubID)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
