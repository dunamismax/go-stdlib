package models

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dunamismax/go-stdlib/pkg/database"
)

type User struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	Bio         string    `json:"bio"`
	AvatarURL   string    `json:"avatar_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LikeCount int       `json:"like_count"`
	IsLiked   bool      `json:"is_liked"`
}

type UserService struct {
	db *database.DB
}

func NewUserService(db *database.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(username, email, password, displayName string) (*User, error) {
	hashedPassword := hashPassword(password)
	
	user, err := s.db.CreateUser(username, email, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	
	return &User{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: displayName,
		Bio:         user.Bio,
		AvatarURL:   user.AvatarURL,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}

func (s *UserService) GetUserByID(id int) (*User, error) {
	user, err := s.db.GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	return &User{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Bio:         user.Bio,
		AvatarURL:   user.AvatarURL,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}

func (s *UserService) GetUserByUsername(username string) (*User, error) {
	user, err := s.db.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	return &User{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Bio:         user.Bio,
		AvatarURL:   user.AvatarURL,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}

func (s *UserService) AuthenticateUser(username, password string) (*User, error) {
	user, err := s.db.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	
	if !verifyPassword(password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid credentials")
	}
	
	return &User{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Bio:         user.Bio,
		AvatarURL:   user.AvatarURL,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}

func (s *UserService) CreatePost(userID int, content string) (*Post, error) {
	post, err := s.db.CreatePost(userID, content)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}
	
	// Get username for the post
	user, err := s.db.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	return &Post{
		ID:        post.ID,
		UserID:    post.UserID,
		Content:   post.Content,
		Username:  user.Username,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		LikeCount: 0,
		IsLiked:   false,
	}, nil
}

func (s *UserService) GetPostByID(postID, userID int) (*Post, error) {
	post, err := s.db.GetPostByID(postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	// Get username for the post
	user, err := s.db.GetUserByID(post.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get like count and check if current user liked it
	likeCount, err := s.db.GetLikeCount(postID)
	if err != nil {
		likeCount = 0
	}
	
	isLiked, err := s.db.IsPostLiked(userID, postID)
	if err != nil {
		isLiked = false
	}

	return &Post{
		ID:        post.ID,
		UserID:    post.UserID,
		Content:   post.Content,
		Username:  user.Username,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		LikeCount: likeCount,
		IsLiked:   isLiked,
	}, nil
}

func (s *UserService) GetRecentPosts(userID int, limit int) ([]*Post, error) {
	posts, err := s.db.GetRecentPosts(limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}
	
	var result []*Post
	for _, post := range posts {
		// Get username for each post
		user, err := s.db.GetUserByID(post.UserID)
		if err != nil {
			continue // Skip posts with invalid users
		}
		
		// Calculate like count and check if current user liked it
		likeCount, err := s.db.GetLikeCount(post.ID)
		if err != nil {
			likeCount = 0
		}
		
		isLiked, err := s.db.IsPostLiked(userID, post.ID)
		if err != nil {
			isLiked = false
		}
		
		result = append(result, &Post{
			ID:        post.ID,
			UserID:    post.UserID,
			Content:   post.Content,
			Username:  user.Username,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
			LikeCount: likeCount,
			IsLiked:   isLiked,
		})
	}
	
	return result, nil
}

func (s *UserService) LikePost(userID, postID int) error {
	return s.db.LikePost(userID, postID)
}

func (s *UserService) UnlikePost(userID, postID int) error {
	return s.db.UnlikePost(userID, postID)
}

// Helper functions
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

func verifyPassword(password, hash string) bool {
	return hashPassword(password) == hash
}


// Simple session management using a global map (not production-ready)
var sessions = make(map[string]*User)

func (s *UserService) CreateSession(user *User) string {
	sessionID := generateSessionID(user.ID)
	sessions[sessionID] = user
	return sessionID
}

func (s *UserService) GetUserFromSession(sessionID string) *User {
	return sessions[sessionID]
}

func (s *UserService) DeleteSession(sessionID string) {
	delete(sessions, sessionID)
}

func generateSessionID(userID int) string {
	timestamp := time.Now().UnixNano()
	data := fmt.Sprintf("%d-%d", userID, timestamp)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)[:32]
}

func ExtractUserIDFromPath(path string) (int, error) {
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return 0, fmt.Errorf("invalid path")
	}
	
	idStr := parts[len(parts)-1]
	return strconv.Atoi(idStr)
}