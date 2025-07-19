package main

import (
	_ "embed"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/dunamismax/go-stdlib/apps/web/gohyperdocs/handlers"
	"github.com/dunamismax/go-stdlib/apps/web/gohyperdocs/models"
	"github.com/dunamismax/go-stdlib/pkg/database"
)

//go:embed static/htmx.min.js
var htmxJS []byte

//go:embed static/styles.css
var stylesCSS []byte

//go:embed templates/home.html
var homeTemplate string

//go:embed templates/section.html
var sectionTemplate string

//go:embed templates/category.html
var categoryTemplate string

//go:embed templates/search-results.html
var searchResultsTemplate string

//go:embed templates/section-details.html
var sectionDetailsTemplate string

//go:embed templates/layout.html
var layoutTemplate string

func main() {
	// Setup structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Initialize database
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

	// Initialize docs service and create tables
	docsService := models.NewDocsService(db)
	if err := docsService.CreateTables(); err != nil {
		slog.Error("Failed to create tables", "error", err)
		log.Fatal("Failed to create tables:", err)
	}

	// Seed data
	if err := docsService.SeedData(); err != nil {
		slog.Error("Failed to seed data", "error", err)
		log.Fatal("Failed to seed data:", err)
	}

	// Create templates with custom functions
	templates := template.New("").Funcs(template.FuncMap{
		"formatTime": func(t interface{}) string {
			return "Jan 2, 2006"
		},
		"titleCase": func(s string) string {
			return template.HTMLEscapeString(s)
		},
	})

	templates = template.Must(templates.Parse(layoutTemplate))
	templates = template.Must(templates.Parse(homeTemplate))
	templates = template.Must(templates.Parse(sectionTemplate))
	templates = template.Must(templates.Parse(categoryTemplate))
	templates = template.Must(templates.Parse(searchResultsTemplate))
	templates = template.Must(templates.Parse(sectionDetailsTemplate))

	// Initialize handlers
	docsHandler := handlers.NewDocsHandler(docsService, templates)

	// Setup routes
	mux := http.NewServeMux()

	// Static files
	mux.HandleFunc("GET /static/htmx.min.js", func(w http.ResponseWriter, r *http.Request) {
		docsHandler.ServeStatic(w, r, htmxJS, "application/javascript")
	})
	mux.HandleFunc("GET /static/styles.css", func(w http.ResponseWriter, r *http.Request) {
		docsHandler.ServeStatic(w, r, stylesCSS, "text/css")
	})

	// API endpoints for HTMX
	mux.HandleFunc("GET /api/search", docsHandler.SearchHandler)
	mux.HandleFunc("GET /api/details/{slug}", docsHandler.LazyLoadDetailsHandler)
	mux.HandleFunc("POST /api/copy", docsHandler.CopyHandler)

	// Interactive demo endpoints
	mux.HandleFunc("GET /api/demo/live-counter", docsHandler.LiveCounterHandler)
	mux.HandleFunc("GET /api/demo/validate-email", docsHandler.FormValidationHandler)
	mux.HandleFunc("GET /api/demo/load-more", docsHandler.LoadMoreContentHandler)
	mux.HandleFunc("GET /api/demo/todo", docsHandler.TodoDemoHandler)
	mux.HandleFunc("POST /api/demo/todo", docsHandler.TodoDemoHandler)
	mux.HandleFunc("DELETE /api/demo/todo/{id}", docsHandler.TodoDemoHandler)
	mux.HandleFunc("GET /api/demo/tabs", docsHandler.TabDemoHandler)

	// JSON API endpoints
	mux.HandleFunc("GET /api/docs", docsHandler.APIDocsHandler)
	mux.HandleFunc("GET /api/section/{slug}", docsHandler.APISectionHandler)

	// Page routes
	mux.HandleFunc("GET /", docsHandler.HomeHandler)
	mux.HandleFunc("GET /section/{slug}", docsHandler.SectionHandler)
	mux.HandleFunc("GET /category/{category}", docsHandler.CategoryHandler)

	// Apply middleware
	finalHandler := corsMiddleware(loggerMiddleware(mux))

	server := &http.Server{
		Addr:         ":8082",
		Handler:      finalHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	slog.Info("GoHyperDocs server starting", "port", ":8082", "url", "http://localhost:8082")
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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
