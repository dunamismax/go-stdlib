package handlers

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/dunamismax/go-stdlib/apps/web/go-social/models"
	"github.com/dunamismax/go-stdlib/pkg/utils"
)



type Handler struct {
	userService *models.UserService
	templates   *template.Template
}

type PageData struct {
	Title      string
	IsLoggedIn bool
	Username   string
	Posts      []*models.Post
	User       *models.User
}

func NewHandler(userService *models.UserService, templates *template.Template) *Handler {
	return &Handler{
		userService: userService,
		templates:   templates,
	}
}

func (h *Handler) ServeCSS(w http.ResponseWriter, r *http.Request, cssData []byte) {
	w.Header().Set("Content-Type", "text/css")
	w.Write(cssData)
}

func isHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := h.getCurrentUser(r)
	userID := 0
	isLoggedIn := false

	if currentUser != nil {
		userID = currentUser.ID
		isLoggedIn = true
	}

	posts, err := h.userService.GetRecentPosts(userID, 20)
	if err != nil {
		posts = []*models.Post{}
	}

	data := PageData{
		Title:      "GoSocial",
		IsLoggedIn: isLoggedIn,
		Posts:      posts,
		User:       currentUser,
	}

	if currentUser != nil {
		data.Username = currentUser.Username
	}

	if err := h.templates.ExecuteTemplate(w, "home.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := h.getCurrentUser(r)
	if currentUser == nil {
		if isHTMXRequest(r) {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, `<div class="error">Must be logged in to post</div>`)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		if isHTMXRequest(r) {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, `<div class="error">Invalid form data</div>`)
			return
		}
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	content := utils.SanitizeInput(r.FormValue("content"))

	if validationErr := utils.ValidatePostContent(content); validationErr != nil {
		if isHTMXRequest(r) {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, `<div class="error">Post content is invalid</div>`)
			return
		}
		http.Redirect(w, r, "/?error=validation_failed", http.StatusSeeOther)
		return
	}

	post, err := h.userService.CreatePost(currentUser.ID, content)
	if err != nil {
		if isHTMXRequest(r) {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, `<div class="error">Failed to create post</div>`)
			return
		}
		http.Redirect(w, r, "/?error=post_failed", http.StatusSeeOther)
		return
	}

	if isHTMXRequest(r) {
		w.Header().Set("Content-Type", "text/html")
		
		// Create a new post with the current user data
		newPost := &models.Post{
			ID:        post.ID,
			Username:  currentUser.Username,
			Content:   post.Content,
			CreatedAt: post.CreatedAt,
			LikeCount: 0,
			IsLiked:   false,
		}

		// Render just the post template
		postData := struct {
			*models.Post
			IsLoggedIn bool
		}{
			Post:       newPost,
			IsLoggedIn: true,
		}

		if err := h.templates.ExecuteTemplate(w, "post", postData); err != nil {
			fmt.Fprint(w, `<div class="error">Failed to render post</div>`)
		}
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) LikePostHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := h.getCurrentUser(r)
	if currentUser == nil {
		if isHTMXRequest(r) {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, `<button class="like-btn">❤️ Must login</button>`)
			return
		}
		utils.Error(w, http.StatusUnauthorized, "Must be logged in")
		return
	}

	// Extract post ID from URL path using Go 1.22+ path parameters
	postIDStr := r.PathValue("postId")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		if isHTMXRequest(r) {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, `<button class="like-btn">❤️ Error</button>`)
			return
		}
		utils.Error(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	// Check if already liked
	post, err := h.userService.GetPostByID(postID, currentUser.ID)
	if err != nil {
		if isHTMXRequest(r) {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, `<button class="like-btn">❤️ Not found</button>`)
			return
		}
		utils.Error(w, http.StatusNotFound, "Post not found")
		return
	}

	if post.IsLiked {
		err = h.userService.UnlikePost(currentUser.ID, postID)
		post.IsLiked = false
		if post.LikeCount > 0 {
			post.LikeCount--
		}
	} else {
		err = h.userService.LikePost(currentUser.ID, postID)
		post.IsLiked = true
		post.LikeCount++
	}

	if err != nil {
		if isHTMXRequest(r) {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, `<button class="like-btn">❤️ Error</button>`)
			return
		}
		utils.Error(w, http.StatusInternalServerError, "Failed to update like")
		return
	}

	if isHTMXRequest(r) {
		w.Header().Set("Content-Type", "text/html")

		likeClass := "like-btn"
		if post.IsLiked {
			likeClass = "like-btn liked"
		}

		fmt.Fprintf(w, `<button hx-post="/like/%d" hx-target="this" hx-swap="outerHTML" class="%s">❤️ %d</button>`,
			postID, likeClass, post.LikeCount)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := h.getCurrentUser(r)
	userID := 0

	if currentUser != nil {
		userID = currentUser.ID
	}

	posts, err := h.userService.GetRecentPosts(userID, 20)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to load posts")
		return
	}

	utils.Success(w, posts)
}

func (h *Handler) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := h.getCurrentUser(r)
	if currentUser == nil {
		utils.Error(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	utils.Success(w, currentUser)
}