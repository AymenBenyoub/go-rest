package repository

import (
	"database/sql"
	"errors"
	"log"
	"rest/internal/db"
)

// PostRepository defines the interface for post-related database operations.

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {

	return &PostRepository{
		db: db}
}

func (r *PostRepository) GetAllPosts() ([]*db.Post, error) {

	query := "SELECT * FROM posts"
	rows, err := r.db.Query(query)

	if err != nil {
		log.Printf("Failed to retrieve posts: %v", err)
		return nil, err
	}
	defer rows.Close()
	var posts []*db.Post
	for rows.Next() {
		var post db.Post
		if err := rows.Scan(&post.Id, &post.Title, &post.Text, &post.PosterID, &post.PostedAt); err != nil {
			log.Printf("Failed to scan post: %v", err)
			return nil, err
		}
		posts = append(posts, &post)
	}
	return posts, nil
}

func (r *PostRepository) GetPostByID(id int) (*db.Post, error) {
	query := "SELECT * FROM posts WHERE id =?"

	var post db.Post
	err := r.db.QueryRow(query, id).Scan(&post.Id, &post.Title, &post.Text, &post.PostedAt, &post.PosterID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Post with ID %d not found", id)
			return nil, errors.New("post now found")
		}
		log.Printf("Failed to retrieve post with ID %d: %v", id, err)
		return nil, err
	}
	return &post, nil

}

func (r *PostRepository) GetPostsByUserID(uid int) ([]*db.Post, error) {

	query := "SELECT * FROM posts WHERE poster = ?"
	var posts []*db.Post
	rows, err := r.db.Query(query, uid)
	if err != nil {
		log.Printf("Failed to retrieve posts for user ID %d: %v", uid, err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var post db.Post
		if err := rows.Scan(&post.Id, &post.Title, &post.Text, &post.PostedAt, &post.PosterID); err != nil {
			log.Printf("Failed to scan post for user ID %d: %v", uid, err)
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *PostRepository) CreatePost(post *db.Post) error {
	query := "INSERT INTO posts (title,text,poster) VALUES (?,?,?)"

	_, err := r.db.Exec(query, post.Title, post.Text, post.PosterID)
	if err != nil {
		log.Printf("Failed to create post: %v", err)
		return err
	}
	return nil
}

func (r *PostRepository) UpdatePost(post *db.Post) error {
	query := "UPDATE posts SET title = ?, text = ? WHERE id = ?"

	_, err := r.db.Exec(query, post.Title, post.Text, post.Id)
	if err != nil {
		log.Printf("Failed to update post with ID %d: %v", post.Id, err)
		return err
	}
	return nil
}

func (r *PostRepository) DeletePost(id int) error {
	query := "DELETE FROM posts WHERE id = ? "

	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("Failed to delete post with ID %d: %v", id, err)
		return err
	}
	return nil
}
