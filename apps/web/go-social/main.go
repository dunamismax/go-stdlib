package main

import (
	_ "embed"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/dunamismax/go-stdlib/apps/web/go-social/handlers"
	"github.com/dunamismax/go-stdlib/apps/web/go-social/models"
	"github.com/dunamismax/go-stdlib/pkg/database"
)

//go:embed static/htmx.min.js
var htmxJS []byte

//go:embed static/styles.css
var stylesCSS []byte

//go:embed templates/layout.html
var layoutTemplate string

//go:embed templates/login.html
var loginTemplate string

//go:embed templates/register.html
var registerTemplate string

//go:embed templates/home.html
var homeTemplate string

func main() {
	db, err := database.NewDB("./data")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
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

	// Static files
	mux.HandleFunc("/static/htmx.min.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(htmxJS)
	})
	mux.HandleFunc("/static/styles.css", func(w http.ResponseWriter, r *http.Request) {
		handler.ServeCSS(w, r, stylesCSS)
	})

	// Authentication routes
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.LoginHandler(w, r)
		} else {
			handler.LoginPageHandler(w, r)
		}
	})
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.RegisterHandler(w, r)
		} else {
			handler.RegisterPageHandler(w, r)
		}
	})

	// Other routes
	mux.HandleFunc("/logout", handler.LogoutHandler)
	mux.HandleFunc("/post", handler.CreatePostHandler)
	mux.HandleFunc("/like/", handler.LikePostHandler)

	// API endpoints
	mux.HandleFunc("/api/posts", handler.GetPostsHandler)
	mux.HandleFunc("/api/user/me", handler.GetCurrentUserHandler)

	// Pages
	mux.HandleFunc("/", handler.HomeHandler)

	// Apply basic logging middleware
	finalHandler := loggerMiddleware(mux)

	server := &http.Server{
		Addr:         ":8081",
		Handler:      finalHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("GoSocial server starting on :8081")
	log.Fatal(server.ListenAndServe())
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s %v", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}