package repository

import (
	"database/sql"
	"errors"
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
		return nil, err
	}
	defer rows.Close()
	var posts []*db.Post
	for rows.Next() {
		var post db.Post
		if err := rows.Scan(&post.Id, &post.Title, &post.Text, &post.PostedAt, &post.PosterID); err != nil {
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
			return nil, errors.New("post now found")
		}
		return nil, err
	}
	return &post, nil

}

func (r *PostRepository) GetPostsByUserID(uid int) ([]*db.Post, error) {

	query := "SELECT * FROM posts WHERE poster = ?"
	var posts []*db.Post
	rows, err := r.db.Query(query, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var post db.Post
		if err := rows.Scan(&post.Id, &post.Title, &post.Text, &post.PostedAt, &post.PosterID); err != nil {
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
		return err
	}
	return nil
}

func (r *PostRepository) UpdatePost(post *db.Post) error {
	query := "UPDATE posts SET title = ?, text = ? WHERE id = ?"

	_, err := r.db.Exec(query, post.Title, post.Text, post.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostRepository) DeletePost(id int) error {
	query := "DELETE FROM posts WHERE id = ? "

	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
