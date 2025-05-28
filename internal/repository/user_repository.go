package repository

import (
	"database/sql"
	"errors"
	"log"

	"rest/internal/db"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user-related database operations.

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetUserByID(pub_id string) (*db.User, error) {

	query := "SELECT public_id, username, email FROM users where public_id = ?"
	var user db.User
	err := r.db.QueryRow(query, pub_id).Scan(&user.PublicID, &user.Username, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("User with ID %d not found", pub_id)
			return nil, errors.New("user not found")
		}
		log.Printf("Failed to retrieve user with ID %d: %v", pub_id, err)
		return nil, err
	}
	log.Printf("User with ID %d retrieved successfully: %s", pub_id, user.Username)
	return &user, nil
}

// no password crypto for now.
func (r *UserRepository) CreateUser(user *db.User) error {

	user_uuid := uuid.New().String()
	query := "INSERT INTO users (public_id,username,email,password) VALUES (?,?,?,?)"

	_, err := r.db.Exec(query, user_uuid, user.Username, user.Email, user.Password)
	if err != nil {

		log.Printf("Failed to create user %s: %v", user.Username, err)
		return err
	}
	log.Printf("User %s created successfully, id: %s", user.Username, user_uuid)
	return nil
}

func (r *UserRepository) UpdateUsername(pub_id string, newUsername string) error {
	query := "UPDATE users SET username = ? WHERE public_id = ?"
	_, err := r.db.Exec(query, newUsername, pub_id)
	if err != nil {
		log.Printf("Failed to update username for user with ID %d: %v", pub_id, err)
		return err
	}
	log.Printf("username updated successfully for user with ID %d", pub_id)
	return nil
}

// no crypto for now.
func (r *UserRepository) UpdatePassword(pub_id string, newPassword string) error {
	query := "UPDATE users SET password = ? WHERE public_id = ?"
	_, err := r.db.Exec(query, newPassword, pub_id)

	if err != nil {
		log.Printf("Failed to update password for user with ID %d: %v", pub_id, err)
		return err
	}
	log.Printf("Password updated successfully for user with ID %d", pub_id)
	return nil
}

func (r *UserRepository) DeleteUser(pub_id string) error {
	query := "DELETE FROM users WHERE public_id = ?"
	_, err := r.db.Exec(query, pub_id)
	if err != nil {
		log.Printf("Failed to delete user with ID %d: %v", pub_id, err)
		return err
	}
	log.Printf("User with ID %d deleted successfully", pub_id)
	return nil
}

