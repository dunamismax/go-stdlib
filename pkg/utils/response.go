package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		// Log the error, but don't return an error to the caller as the headers have already been written.
		fmt.Printf("Error encoding JSON response: %v\n", err)
	}
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]string{"error": message})
}

func Success(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": data})
}

type SessionToken struct {
	UserID    int    `json:"user_id"`
	ExpiresAt int64  `json:"expires_at"`
	Nonce     string `json:"nonce"`
}

func CreateSessionToken(userID int, secretKey string, expirationHours int) (string, error) {
	expiresAt := time.Now().Add(time.Duration(expirationHours) * time.Hour).Unix()
	nonce, err := SecureRandomHex(16)
	if err != nil {
		return "", err
	}

	token := SessionToken{
		UserID:    userID,
		ExpiresAt: expiresAt,
		Nonce:     nonce,
	}

	tokenJSON, err := json.Marshal(token)
	if err != nil {
		return "", err
	}

	tokenB64 := base64.URLEncoding.EncodeToString(tokenJSON)
	signature := createSignature(tokenB64, secretKey)

	return fmt.Sprintf("%s.%s", tokenB64, signature), nil
}

func ValidateSessionToken(tokenString string, secretKey string) (*SessionToken, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid token format")
	}

	tokenB64, signature := parts[0], parts[1]

	expectedSignature := createSignature(tokenB64, secretKey)
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return nil, fmt.Errorf("invalid token signature")
	}

	tokenJSON, err := base64.URLEncoding.DecodeString(tokenB64)
	if err != nil {
		return nil, fmt.Errorf("invalid token encoding")
	}

	var token SessionToken
	if err := json.Unmarshal(tokenJSON, &token); err != nil {
		return nil, fmt.Errorf("invalid token data")
	}

	if time.Now().Unix() > token.ExpiresAt {
		return nil, fmt.Errorf("token expired")
	}

	return &token, nil
}

func createSignature(data, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func SetSecureCookie(w http.ResponseWriter, name, value string, maxAge int, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	})
}

func ClearCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}
