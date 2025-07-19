package main

import (
	_ "embed"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/dunamismax/go-stdlib/apps/web/go-social/handlers"
	"github.com/dunamismax/go-stdlib/apps/web/go-social/models"
	"github.com/dunamismax/go-stdlib/pkg/database"
)

//go:embed dist
var distFS embed.FS

//go:embed static/htmx.min.js
var htmxJS []byte

//go:embed templates/layout.html
var layoutTemplate string

//go:embed templates/login.html
var loginTemplate string

//go:embed templates/register.html
var registerTemplate string

//go:embed templates/home.html
var homeTemplate string

func main() {
	// Setup structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	db, err := database.NewDB("./data")
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		slog.Error("Failed to run migrations", "error", err)
		log.Fatal("Failed to run migrations:", err)
	}

	userService := models.NewUserService(db)

	// Create templates
	templates := template.New("").Funcs(template.FuncMap{
		"formatTime": func(t interface{}) string {
			return "Jan 2, 2006"
		},
	})
	templates = template.Must(templates.Parse(layoutTemplate))
	templates = template.Must(templates.Parse(loginTemplate))
	templates = template.Must(templates.Parse(registerTemplate))
	templates = template.Must(templates.Parse(homeTemplate))

	handler := handlers.NewHandler(userService, templates)

	mux := http.NewServeMux()

	// Static files - serve both old and new assets
	mux.HandleFunc("GET /static/htmx.min.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(htmxJS)
	})

	// Serve built assets from Vite
	assetsHandler := http.FileServer(http.FS(distFS))
	mux.Handle("GET /assets/", http.StripPrefix("/", assetsHandler))

	// Authentication routes
	mux.HandleFunc("GET /login", handler.LoginPageHandler)
	mux.HandleFunc("POST /login", handler.LoginHandler)
	mux.HandleFunc("GET /register", handler.RegisterPageHandler)
	mux.HandleFunc("POST /register", handler.RegisterHandler)

	// Other routes
	mux.HandleFunc("POST /logout", handler.LogoutHandler)
	mux.HandleFunc("POST /post", handler.CreatePostHandler)
	mux.HandleFunc("POST /like/{postId}", handler.LikePostHandler)

	// API endpoints
	mux.HandleFunc("GET /api/posts", handler.GetPostsHandler)
	mux.HandleFunc("GET /api/user/me", handler.GetCurrentUserHandler)

	// Pages
	mux.HandleFunc("GET /", handler.HomeHandler)

	// Apply basic logging middleware
	finalHandler := loggerMiddleware(mux)

	server := &http.Server{
		Addr:         ":8081",
		Handler:      finalHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	slog.Info("GoSocial server starting", "port", ":8081")
	if err := server.ListenAndServe(); err != nil {
		slog.Error("Server failed to start", "error", err)
		log.Fatal(err)
	}
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("HTTP request",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"duration", time.Since(start),
			"user_agent", r.UserAgent(),
		)
	})
}
