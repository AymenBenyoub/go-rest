package handlers

import (

	"encoding/json"
	"errors"
	"log"
	"net/http"
	"rest/internal/db"
	"rest/internal/repository"
	"strconv"
)

type PostHandler struct {
	repo *repository.PostRepository
}

func NewPostHandler(db *repository.PostRepository) *PostHandler {
	return &PostHandler{
		repo: db,
	}
}

func (h *PostHandler) RegisterRoutes(mux *http.ServeMux) {

	mux.HandleFunc("GET /posts", h.GetAllPosts)
	mux.HandleFunc("GET /posts/{id}", h.GetPostByID)
	mux.HandleFunc("GET /posts/user/{uid}", h.GetPostsByUserID)
	mux.HandleFunc("POST /posts/create", h.CreatePost)
	mux.HandleFunc("PUT /posts/update/{id}", h.UpdatePost)
	mux.HandleFunc("DELETE /posts/delete/{id}", h.DeletePost)
}


func (h *PostHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	posts, err := h.repo.GetAllPosts(ctx)
	if err != nil {
		log.Printf("Error retrieving posts: %v", err)
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
func (h *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	if id == "" {

		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(id)
	if err != nil {

		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	post, err := h.repo.GetPostByID(ctx,postID)
	if err != nil {
		if errors.Is(err, repository.ErrPostNotFound) {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		log.Printf("Error retrieving post with ID %d: %v", postID, err)
		http.Error(w, "Failed to retrieve post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
func (h *PostHandler) GetPostsByUserID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid := r.PathValue("uid")
	if uid == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	posts, err := h.repo.GetPostsByUserID(ctx,uid)
	if err != nil {
		log.Printf("Error retrieving posts for user %s: %v", uid, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var post db.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.CreatePost(ctx,&post); err != nil {
		log.Printf("Error creating post: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post created successfully"})
}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	var post db.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	post.Id = postID

	if err := h.repo.UpdatePost(ctx,&post); err != nil {
		if errors.Is(err, repository.ErrPostNotFound) {
			log.Printf("Post with ID %d not found", postID)
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		log.Printf("Error updating post with ID %d: %v", postID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post updated successfully"})
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeletePost(ctx,postID); err != nil {
		if errors.Is(err, repository.ErrPostNotFound) {
			log.Printf("Post with ID %d not found", postID)
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		log.Printf("Error deleting post with ID %d: %v", postID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post deleted successfully"})
}
