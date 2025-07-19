package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
}

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	DisplayName  string    `json:"display_name"`
	Bio          string    `json:"bio"`
	AvatarURL    string    `json:"avatar_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Follow struct {
	ID          int       `json:"id"`
	FollowerID  int       `json:"follower_id"`
	FollowingID int       `json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type Like struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	PostID    int       `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

func NewDB(dataDir string) (*DB, error) {
	if err := os.MkdirAll(dataDir, 0750); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	dbPath := filepath.Join(dataDir, "app.db")
	slog.Info("Opening database", "path", dbPath)
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection
	conn.SetMaxOpenConns(1)
	conn.SetMaxIdleConns(1)
	conn.SetConnMaxLifetime(0)

	return &DB{
		conn: conn,
	}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) Migrate() error {
	slog.Info("Running database migrations")
	// Create tables
	schema := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			display_name TEXT DEFAULT '',
			bio TEXT DEFAULT '',
			avatar_url TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id)
		);

		CREATE TABLE IF NOT EXISTS follows (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			follower_id INTEGER NOT NULL,
			following_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (follower_id) REFERENCES users (id),
			FOREIGN KEY (following_id) REFERENCES users (id),
			UNIQUE (follower_id, following_id)
		);

		CREATE TABLE IF NOT EXISTS likes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			post_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id),
			FOREIGN KEY (post_id) REFERENCES posts (id),
			UNIQUE (user_id, post_id)
		);

		CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);
		CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
		CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts (user_id);
		CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts (created_at DESC);
		CREATE INDEX IF NOT EXISTS idx_follows_follower_id ON follows (follower_id);
		CREATE INDEX IF NOT EXISTS idx_follows_following_id ON follows (following_id);
		CREATE INDEX IF NOT EXISTS idx_likes_user_id ON likes (user_id);
		CREATE INDEX IF NOT EXISTS idx_likes_post_id ON likes (post_id);
	`

	_, err := db.conn.Exec(schema)
	if err != nil {
		slog.Error("Failed to create tables", "error", err)
		return fmt.Errorf("failed to create tables: %w", err)
	}

	slog.Info("Database migrations completed successfully")
	return nil
}

func (db *DB) GetUserByUsername(username string) (*User, error) {
	query := `SELECT id, username, email, password_hash, display_name, bio, avatar_url, created_at, updated_at 
			 FROM users WHERE username = ?`

	var user User
	err := db.conn.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.DisplayName, &user.Bio, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (db *DB) GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, username, email, password_hash, display_name, bio, avatar_url, created_at, updated_at 
			 FROM users WHERE email = ?`

	var user User
	err := db.conn.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.DisplayName, &user.Bio, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (db *DB) CreateUser(username, email, passwordHash string) (*User, error) {
	query := `INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)`

	result, err := db.conn.Exec(query, username, email, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	// Return the created user
	return db.GetUserByID(int(id))
}

func (db *DB) GetUserByID(id int) (*User, error) {
	query := `SELECT id, username, email, password_hash, display_name, bio, avatar_url, created_at, updated_at 
			 FROM users WHERE id = ?`

	var user User
	err := db.conn.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.DisplayName, &user.Bio, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (db *DB) GetRecentPosts(limit int) ([]Post, error) {
	query := `SELECT id, user_id, content, created_at, updated_at 
			 FROM posts ORDER BY created_at DESC LIMIT ?`

	rows, err := db.conn.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating posts: %w", err)
	}

	return posts, nil
}

func (db *DB) CreatePost(userID int, content string) (*Post, error) {
	query := `INSERT INTO posts (user_id, content) VALUES (?, ?)`

	result, err := db.conn.Exec(query, userID, content)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get post ID: %w", err)
	}

	// Return the created post
	return db.GetPostByID(int(id))
}

func (db *DB) GetPostByID(id int) (*Post, error) {
	query := `SELECT id, user_id, content, created_at, updated_at FROM posts WHERE id = ?`

	var post Post
	err := db.conn.QueryRow(query, id).Scan(
		&post.ID, &post.UserID, &post.Content, &post.CreatedAt, &post.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post not found")
		}
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return &post, nil
}

func (db *DB) LikePost(userID, postID int) error {
	query := `INSERT INTO likes (user_id, post_id) VALUES (?, ?) ON CONFLICT DO NOTHING`

	_, err := db.conn.Exec(query, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to like post: %w", err)
	}

	return nil
}

func (db *DB) UnlikePost(userID, postID int) error {
	query := `DELETE FROM likes WHERE user_id = ? AND post_id = ?`

	_, err := db.conn.Exec(query, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to unlike post: %w", err)
	}

	return nil
}

func (db *DB) GetLikeCount(postID int) (int, error) {
	query := `SELECT COUNT(*) FROM likes WHERE post_id = ?`

	var count int
	err := db.conn.QueryRow(query, postID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get like count: %w", err)
	}

	return count, nil
}

func (db *DB) IsPostLiked(userID, postID int) (bool, error) {
	query := `SELECT COUNT(*) FROM likes WHERE user_id = ? AND post_id = ?`

	var count int
	err := db.conn.QueryRow(query, userID, postID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check like status: %w", err)
	}

	return count > 0, nil
}
