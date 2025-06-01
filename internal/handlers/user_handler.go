package handlers

import (
	
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"rest/internal/db"
	"rest/internal/repository"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /users/{id}", h.GetUser)
	mux.HandleFunc("POST /users/create", h.CreateUser)
	mux.HandleFunc("PUT /users/update/username/{id}", h.UpdateUsername)
	mux.HandleFunc("PUT /users/update/password/{id}", h.UpdatePassword)
	mux.HandleFunc("DELETE /users/delete/{id}", h.DeleteUser)
	mux.HandleFunc("GET /users", h.GetAllUsers)

}

// will add logging later
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetUserByID(ctx,id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			log.Printf("User with ID %s not found", id)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Printf("Error retrieving user with ID %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var user db.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.repo.CreateUser(ctx,&user); err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

func (h *UserHandler) UpdateUsername(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var newUsername string
	id := r.PathValue("id")

	if err := json.NewDecoder(r.Body).Decode(&newUsername); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.repo.UpdateUsername(ctx,id, newUsername); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			log.Printf("User with ID %s not found", id)
			return
		}
		log.Printf("Error updating username for user %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Username updated successfully"})
}

func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	var newPassword string
	if err := json.NewDecoder(r.Body).Decode(&newPassword); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdatePassword(ctx,id, newPassword); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			log.Printf("User with ID %s not found", id)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Printf("Error updating password for user %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password updated successfully"})
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubID := r.PathValue("id")
	if pubID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteUser(ctx,pubID); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			log.Printf("User with ID %s not found", pubID)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		log.Printf("Error deleting user with ID %s: %v", pubID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})

}



//protected endpoint
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := h.repo.GetAllUsers(ctx)
	if err != nil {
		
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}
