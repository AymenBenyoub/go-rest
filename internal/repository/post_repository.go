package repository

import (
	"database/sql"
	"errors"

	"rest/internal/db"
)

var ErrPostNotFound = errors.New("post not found")

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) GetAllPosts() ([]*db.Post, error) {
	query := "SELECT id, title, text, poster, posted_at FROM posts"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*db.Post
	for rows.Next() {
		var post db.Post
		if err := rows.Scan(&post.Id, &post.Title, &post.Text, &post.PosterID, &post.PostedAt); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) GetPostByID(id int) (*db.Post, error) {
	query := "SELECT id, title, text, poster, posted_at FROM posts WHERE id = ?"
	var post db.Post
	err := r.db.QueryRow(query, id).Scan(&post.Id, &post.Title, &post.Text, &post.PosterID, &post.PostedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) GetPostsByUserID(uid int) ([]*db.Post, error) {
	query := "SELECT id, title, text, poster, posted_at FROM posts WHERE poster = ?"
	rows, err := r.db.Query(query, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*db.Post
	for rows.Next() {
		var post db.Post
		if err := rows.Scan(&post.Id, &post.Title, &post.Text, &post.PosterID, &post.PostedAt); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) CreatePost(post *db.Post) error {
	query := "INSERT INTO posts (title, text, poster) VALUES (?, ?, ?)"
	_, err := r.db.Exec(query, post.Title, post.Text, post.PosterID)
	return err
}

func (r *PostRepository) UpdatePost(post *db.Post) error {
	query := "UPDATE posts SET title = ?, text = ? WHERE id = ?"
	res, err := r.db.Exec(query, post.Title, post.Text, post.Id)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrPostNotFound
	}

	return nil
}

func (r *PostRepository) DeletePost(id int) error {
	query := "DELETE FROM posts WHERE id = ?"
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrPostNotFound
	}

	return nil
}
