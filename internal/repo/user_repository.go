package repo

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

func (r *UserRepository) GetUserByID(id int) (*db.User, error) {


	query := "SELECT public_id, username, email FROM users where public_id = ?"
	var user db.User 
	err := r.db.QueryRow(query,id).Scan(&user.PublicID, &user.Username, &user.Email)
	if err!= nil {
		if errors.Is(err,sql.ErrNoRows ){
			return nil, errors.New("user not found")
		}
		return nil, err
	}
 
	return &user,nil 
}
// no password crypto for now.
func (r *UserRepository) CreateUser (user *db.User) error {

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

func (r *UserRepository) UpdateUsername (id int, newUsername string ) error {
	query := "UPDATE users SET username = ? WHERE public_id = ?"
	_ ,err := r.db.Exec(query, newUsername, id)
	if err != nil {
		log.Printf("Failed to update username for user with ID %d: %v", id, err)
		return err
	}
	log.Printf("username updated successfully for user with ID %d", id)
	return nil
}

// no crypto for now.
func (r *UserRepository) UpdatePassword(id int, newPassword string) error {
	query:= "UPDATE users SET password = ? WHERE public_id = ?"
	_, err := r.db.Exec(query, newPassword,id)

	if err != nil {
		log.Printf("Failed to update password for user with ID %d: %v", id, err)
		return err
	}
	log.Printf("Password updated successfully for user with ID %d", id)
	return nil
}