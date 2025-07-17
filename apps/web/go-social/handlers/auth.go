package handlers

import (
	"net/http"
	"os"

	"github.com/dunamismax/go-stdlib/apps/web/go-social/models"
	"github.com/dunamismax/go-stdlib/pkg/utils"
)

func (h *Handler) LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:      "Login - GoSocial",
		IsLoggedIn: false,
	}

	if err := h.templates.ExecuteTemplate(w, "login.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	username := utils.SanitizeInput(r.FormValue("username"))
	password := r.FormValue("password")

	if username == "" || password == "" {
		http.Redirect(w, r, "/login?error=missing_fields", http.StatusSeeOther)
		return
	}

	if len(username) > 20 || len(password) > 128 {
		http.Redirect(w, r, "/login?error=invalid_credentials", http.StatusSeeOther)
		return
	}

	user, err := h.userService.AuthenticateUser(username, password)
	if err != nil {
		http.Redirect(w, r, "/login?error=invalid_credentials", http.StatusSeeOther)
		return
	}

	h.setSession(w, user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:      "Register - GoSocial",
		IsLoggedIn: false,
	}

	if err := h.templates.ExecuteTemplate(w, "register.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	username := utils.SanitizeInput(r.FormValue("username"))
	email := utils.SanitizeInput(r.FormValue("email"))
	displayName := utils.SanitizeInput(r.FormValue("display_name"))
	password := r.FormValue("password")

	validationErrors := utils.ValidateUserRegistration(username, email, displayName, password)
	if validationErrors.HasErrors() {
		http.Redirect(w, r, "/register?error=validation_failed", http.StatusSeeOther)
		return
	}

	user, err := h.userService.CreateUser(username, email, password, displayName)
	if err != nil {
		http.Redirect(w, r, "/register?error=creation_failed", http.StatusSeeOther)
		return
	}

	h.setSession(w, user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	h.clearSession(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) getSessionSecret() string {
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "default-secret-key-change-in-production"
	}
	return secret
}

func (h *Handler) setSession(w http.ResponseWriter, userID int) {
	secretKey := h.getSessionSecret()
	token, err := utils.CreateSessionToken(userID, secretKey, 168) // 7 days
	if err != nil {
		return
	}

	isSecure := os.Getenv("HTTPS") == "true"
	utils.SetSecureCookie(w, "session", token, 86400*7, isSecure)
}

func (h *Handler) clearSession(w http.ResponseWriter) {
	utils.ClearCookie(w, "session")
}

func (h *Handler) getCurrentUser(r *http.Request) *models.User {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil
	}

	secretKey := h.getSessionSecret()
	sessionToken, err := utils.ValidateSessionToken(cookie.Value, secretKey)
	if err != nil {
		return nil
	}

	user, err := h.userService.GetUserByID(sessionToken.UserID)
	if err != nil {
		return nil
	}

	return user
}