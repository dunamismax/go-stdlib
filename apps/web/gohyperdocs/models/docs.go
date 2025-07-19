package models

import (
	"fmt"
	"time"

	"github.com/dunamismax/go-stdlib/pkg/database"
)

type DocSection struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Content     string    `json:"content"`
	CodeExample string    `json:"code_example"`
	Category    string    `json:"category"`
	Order       int       `json:"order"`
	Searchable  string    `json:"searchable"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DocsService struct {
	db *database.DB
}

func NewDocsService(db *database.DB) *DocsService {
	return &DocsService{db: db}
}

func (s *DocsService) CreateTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS doc_sections (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		slug TEXT UNIQUE NOT NULL,
		content TEXT NOT NULL,
		code_example TEXT,
		category TEXT NOT NULL,
		order_num INTEGER DEFAULT 0,
		searchable TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_doc_sections_slug ON doc_sections(slug);
	CREATE INDEX IF NOT EXISTS idx_doc_sections_category ON doc_sections(category);
	CREATE INDEX IF NOT EXISTS idx_doc_sections_order ON doc_sections(order_num);
	`

	_, err := s.db.GetConnection().Exec(query)
	return err
}

func (s *DocsService) GetAllSections() ([]DocSection, error) {
	query := `
	SELECT id, title, slug, content, COALESCE(code_example, ''), category, 
	       order_num, COALESCE(searchable, ''), created_at, updated_at
	FROM doc_sections 
	ORDER BY category, order_num, title
	`

	rows, err := s.db.GetConnection().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sections []DocSection
	for rows.Next() {
		var section DocSection
		err := rows.Scan(
			&section.ID, &section.Title, &section.Slug, &section.Content,
			&section.CodeExample, &section.Category, &section.Order,
			&section.Searchable, &section.CreatedAt, &section.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		sections = append(sections, section)
	}

	return sections, rows.Err()
}

func (s *DocsService) GetSectionBySlug(slug string) (*DocSection, error) {
	query := `
	SELECT id, title, slug, content, COALESCE(code_example, ''), category, 
	       order_num, COALESCE(searchable, ''), created_at, updated_at
	FROM doc_sections 
	WHERE slug = ?
	`

	var section DocSection
	err := s.db.GetConnection().QueryRow(query, slug).Scan(
		&section.ID, &section.Title, &section.Slug, &section.Content,
		&section.CodeExample, &section.Category, &section.Order,
		&section.Searchable, &section.CreatedAt, &section.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &section, nil
}

func (s *DocsService) SearchSections(query string) ([]DocSection, error) {
	searchQuery := `
	SELECT id, title, slug, content, COALESCE(code_example, ''), category, 
	       order_num, COALESCE(searchable, ''), created_at, updated_at
	FROM doc_sections 
	WHERE title LIKE ? OR content LIKE ? OR searchable LIKE ?
	ORDER BY 
		CASE 
			WHEN title LIKE ? THEN 1
			WHEN content LIKE ? THEN 2
			ELSE 3
		END,
		category, order_num, title
	LIMIT 50
	`

	searchTerm := "%" + query + "%"
	titleMatch := "%" + query + "%"

	rows, err := s.db.GetConnection().Query(searchQuery, searchTerm, searchTerm, searchTerm, titleMatch, titleMatch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sections []DocSection
	for rows.Next() {
		var section DocSection
		err := rows.Scan(
			&section.ID, &section.Title, &section.Slug, &section.Content,
			&section.CodeExample, &section.Category, &section.Order,
			&section.Searchable, &section.CreatedAt, &section.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		sections = append(sections, section)
	}

	return sections, rows.Err()
}

func (s *DocsService) GetSectionsByCategory(category string) ([]DocSection, error) {
	query := `
	SELECT id, title, slug, content, COALESCE(code_example, ''), category, 
	       order_num, COALESCE(searchable, ''), created_at, updated_at
	FROM doc_sections 
	WHERE category = ?
	ORDER BY order_num, title
	`

	rows, err := s.db.GetConnection().Query(query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sections []DocSection
	for rows.Next() {
		var section DocSection
		err := rows.Scan(
			&section.ID, &section.Title, &section.Slug, &section.Content,
			&section.CodeExample, &section.Category, &section.Order,
			&section.Searchable, &section.CreatedAt, &section.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		sections = append(sections, section)
	}

	return sections, rows.Err()
}

func (s *DocsService) SeedData() error {
	// Check if data already exists
	var count int
	err := s.db.GetConnection().QueryRow("SELECT COUNT(*) FROM doc_sections").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // Data already exists
	}

	sections := []DocSection{
		{
			Title:      "Introduction to Go Hypermedia Stack",
			Slug:       "introduction",
			Content:    "Welcome to GoHyperDocs - the definitive documentation platform demonstrating the power of a pure Go hypermedia stack. This revolutionary approach to web development combines the simplicity and performance of Go's standard library with the elegance of hypermedia-driven user interfaces. Experience the future of web development where complexity is eliminated, not hidden.",
			Category:   "getting-started",
			Order:      1,
			Searchable: "introduction getting started go hypermedia stack overview",
		},
		{
			Title:      "Why Choose This Stack?",
			Slug:       "why-choose-stack",
			Content:    "Traditional web development has become unnecessarily complex with multiple build tools, frameworks, and dependencies. This stack returns to fundamentals: a single Go binary that serves everything. No Node.js, no webpack, no complex build chains. Just pure, fast, maintainable code that deploys anywhere and runs forever.",
			Category:   "getting-started",
			Order:      2,
			Searchable: "why choose advantages benefits simplicity performance maintainability",
		},
		{
			Title:      "Core Technologies Overview",
			Slug:       "core-technologies",
			Content:    "Our stack leverages Go's net/http for routing, html/template for rendering, SQLite for data persistence, and HTMX for dynamic interactions. Each technology was chosen for stability, performance, and minimal dependencies. The result is a stack that will remain relevant and maintainable for decades.",
			Category:   "getting-started",
			Order:      3,
			Searchable: "technologies tech stack net/http html/template sqlite htmx go:embed core",
		},
		{
			Title:   "Performance & Benchmarks",
			Slug:    "performance-benchmarks",
			Content: "This stack delivers exceptional performance: sub-millisecond response times, minimal memory usage, and the ability to handle 50K+ concurrent connections on modest hardware. Compiled Go code is inherently fast, and SQLite provides zero-latency data access without network overhead.",
			CodeExample: `// Benchmark results on a modest VPS:
// Response time: 0.2ms average
// Memory usage: 15MB total
// Concurrent users: 50,000+
// Database queries: 50,000+ QPS
// Binary size: 12MB (includes all assets)

package main

import (
    "testing"
    "net/http/httptest"
    "net/http"
)

func BenchmarkHomePage(b *testing.B) {
    req := httptest.NewRequest("GET", "/", nil)
    for i := 0; i < b.N; i++ {
        w := httptest.NewRecorder()
        handler.ServeHTTP(w, req)
    }
}

// BenchmarkHomePage-8   500000   0.2ms per operation`,
			Category:   "getting-started",
			Order:      4,
			Searchable: "performance benchmarks speed fast response time memory usage",
		},
		{
			Title:   "Quick Start Guide",
			Slug:    "quick-start",
			Content: "Get your Go hypermedia application running in minutes. This guide walks you through creating a minimal but powerful web application that demonstrates all the key concepts of our stack.",
			CodeExample: `// main.go - Complete working application
package main

import (
    "embed"
    "html/template"
    "net/http"
    "log"
)

//go:embed templates/*
var templates embed.FS

//go:embed static/*
var static embed.FS

func main() {
    // Parse templates
    tmpl := template.Must(template.ParseFS(templates, "templates/*"))
    
    // Setup routes
    mux := http.NewServeMux()
    
    // Static files
    mux.Handle("/static/", http.FileServer(http.FS(static)))
    
    // Pages
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        tmpl.ExecuteTemplate(w, "home.html", map[string]string{
            "Title": "My Go App",
            "Message": "Hello, Hypermedia!",
        })
    })
    
    // Start server
    log.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}`,
			Category:   "getting-started",
			Order:      5,
			Searchable: "quick start guide tutorial setup installation",
		},
		{
			Title:   "HTMX Progressive Enhancement",
			Slug:    "htmx-progressive-enhancement",
			Content: "HTMX follows the principle of progressive enhancement - your application works without JavaScript and becomes enhanced with it. This ensures accessibility, SEO compatibility, and graceful degradation across all browsers and devices.",
			CodeExample: `<!-- Works without JavaScript -->
<form action="/submit" method="POST">
  <input name="email" type="email" required>
  <button type="submit">Subscribe</button>
</form>

<!-- Enhanced with HTMX -->
<form hx-post="/submit" 
      hx-target="#result" 
      hx-swap="innerHTML">
  <input name="email" type="email" required>
  <button type="submit">Subscribe</button>
</form>
<div id="result"></div>`,
			Category:   "htmx-features",
			Order:      1,
			Searchable: "htmx progressive enhancement accessibility SEO graceful degradation",
		},
		{
			Title:   "HTMX Navigation with hx-boost",
			Slug:    "htmx-navigation",
			Content: "Experience blazing-fast, SPA-like navigation by applying hx-boost to the body element. This converts all standard navigation links into asynchronous requests for seamless page transitions while maintaining browser history and bookmarkability.",
			CodeExample: `<body hx-boost="true">
  <nav>
    <a href="/">Home</a>
    <a href="/docs">Documentation</a>
    <a href="/api">API Reference</a>
  </nav>
  <main id="content">
    <!-- Content swapped here -->
  </main>
  
  <!-- Loading indicator -->
  <div id="loading" class="htmx-indicator">
    Loading...
  </div>
</body>

<!-- CSS for smooth transitions -->
<style>
.htmx-indicator { display: none; }
.htmx-request .htmx-indicator { display: block; }
</style>`,
			Category:   "htmx-features",
			Order:      2,
			Searchable: "htmx navigation hx-boost spa single page application routing",
		},
		{
			Title:   "Real-time Updates with hx-trigger",
			Slug:    "htmx-real-time",
			Content: "Create real-time, interactive interfaces with HTMX's powerful trigger system. Poll for updates, respond to events, and create live dashboards without complex WebSocket management.",
			CodeExample: `<!-- Auto-refresh every 5 seconds -->
<div hx-get="/api/stats" 
     hx-trigger="every 5s"
     hx-target="this">
  <h3>Server Stats</h3>
  <p>CPU: 45% | Memory: 2.1GB</p>
</div>

<!-- Trigger on multiple events -->
<input hx-get="/api/validate" 
       hx-trigger="input changed delay:500ms, blur"
       hx-target="#validation">

<!-- Load on viewport intersection -->
<div hx-get="/api/lazy-content"
     hx-trigger="intersect once"
     hx-target="this">
  Loading...
</div>`,
			Category:   "htmx-features",
			Order:      3,
			Searchable: "htmx real-time polling triggers events live updates",
		},
		{
			Title:   "Lazy Loading with hx-get",
			Slug:    "lazy-loading",
			Content: "Optimize performance by loading content on-demand. This technique reduces initial page load times and improves user experience by fetching expensive content only when needed.",
			CodeExample: `<button hx-get="/api/details/123" 
        hx-target="#details-container"
        hx-indicator="#loading-spinner"
        hx-swap="innerHTML settle:1s">
  Show Details
</button>

<div id="loading-spinner" class="htmx-indicator">
  <svg class="spinner" viewBox="0 0 24 24">
    <circle cx="12" cy="12" r="10" stroke="currentColor" 
            stroke-width="2" fill="none" 
            stroke-dasharray="31.416" 
            stroke-dashoffset="31.416">
      <animate attributeName="stroke-dasharray" 
               dur="2s" 
               values="0 31.416;15.708 15.708;0 31.416" 
               repeatCount="indefinite"/>
    </svg>
  </div>

<div id="details-container">
  <!-- Details loaded here with smooth transition -->
</div>`,
			Category:   "htmx-features",
			Order:      4,
			Searchable: "lazy loading hx-get on-demand performance optimization",
		},
		{
			Title:   "Form Validation & Feedback",
			Slug:    "form-validation",
			Content: "Create responsive forms with real-time validation using HTMX. Get instant feedback without page refreshes while maintaining accessibility and progressive enhancement.",
			CodeExample: `<form hx-post="/api/register" hx-target="#form-result">
  <div class="field">
    <input name="email" 
           type="email" 
           required
           hx-get="/api/validate-email"
           hx-trigger="input changed delay:500ms"
           hx-target="#email-validation"
           placeholder="Enter your email">
    <div id="email-validation"></div>
  </div>
  
  <div class="field">
    <input name="password" 
           type="password" 
           required
           hx-get="/api/check-password-strength"
           hx-trigger="input changed delay:500ms"
           hx-target="#password-strength"
           placeholder="Enter password">
    <div id="password-strength"></div>
  </div>
  
  <button type="submit" hx-indicator="#submit-spinner">
    Register
    <span id="submit-spinner" class="htmx-indicator">⏳</span>
  </button>
</form>

<div id="form-result"></div>`,
			Category:   "htmx-features",
			Order:      5,
			Searchable: "form validation real-time feedback progressive enhancement accessibility",
		},
		{
			Title:   "Copy to Clipboard Integration",
			Slug:    "copy-clipboard",
			Content: "Interactive code snippets with server-coordinated clipboard functionality. Copy buttons update their state through HTMX calls to Go handlers, providing visual feedback.",
			CodeExample: `<div class="code-block">
  <div class="code-header">
    <span class="code-title">main.go</span>
    <button class="copy-btn" 
            data-code="package main&#10;&#10;import &quot;fmt&quot;&#10;&#10;func main() {&#10;    fmt.Println(&quot;Hello, Go!&quot;)&#10;}"
            hx-post="/api/copy" 
            hx-target="this"
            hx-swap="outerHTML"
            onclick="copyToClipboard(this)">
      Copy
    </button>
  </div>
  <pre><code>package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}</code></pre>
</div>

<script>
function copyToClipboard(btn) {
  navigator.clipboard.writeText(btn.dataset.code);
}
</script>`,
			Category:   "htmx-features",
			Order:      6,
			Searchable: "copy clipboard code snippets javascript htmx interaction",
		},
		{
			Title:   "Intelligent Search with hx-select",
			Slug:    "search-functionality",
			Content: "Real-time search with surgical DOM updates using hx-select. Results are fetched from SQLite and only the relevant content is swapped, maintaining page context and performance.",
			CodeExample: `<div class="search-container">
  <input type="search" 
         hx-get="/api/search" 
         hx-target="#search-results"
         hx-select="#search-results-content"
         hx-trigger="input changed delay:300ms"
         hx-indicator="#search-spinner"
         placeholder="Search documentation..."
         class="search-input">
  
  <div id="search-spinner" class="htmx-indicator">
    Searching...
  </div>
</div>

<div id="search-results">
  <div id="search-results-content">
    <!-- Search results appear here -->
    <!-- Only this div's content is replaced -->
  </div>
  <div class="search-footer">
    <!-- This stays intact during searches -->
    Powered by SQLite FTS
  </div>
</div>`,
			Category:   "htmx-features",
			Order:      7,
			Searchable: "search functionality hx-select real-time sqlite dom updates",
		},
		{
			Title:   "Animation & Transitions",
			Slug:    "htmx-animations",
			Content: "HTMX provides built-in support for smooth transitions and animations. Create engaging user experiences with CSS transitions that automatically trigger during content swaps.",
			CodeExample: `<!-- CSS for smooth transitions -->
<style>
.content-area {
  transition: all 0.3s ease;
}

.htmx-swapping .content-area {
  opacity: 0;
  transform: translateY(20px);
}

.htmx-settling .content-area {
  opacity: 1;
  transform: translateY(0);
}

.fade-in {
  animation: fadeIn 0.5s ease-in;
}

@keyframes fadeIn {
  from { opacity: 0; transform: scale(0.95); }
  to { opacity: 1; transform: scale(1); }
}
</style>

<!-- Animated content loading -->
<button hx-get="/api/content" 
        hx-target="#content" 
        hx-swap="innerHTML settle:300ms"
        class="load-btn">
  Load Content
</button>

<div id="content" class="content-area">
  <div class="fade-in">
    <!-- Content appears with smooth animation -->
  </div>
</div>`,
			Category:   "htmx-features",
			Order:      8,
			Searchable: "htmx animations transitions CSS effects smooth",
		},
		{
			Title:   "Go Standard Library Routing",
			Slug:    "go-routing",
			Content: "Leverage net/http.ServeMux for powerful, zero-dependency routing with method-aware patterns and path parameters. Modern Go routing is more capable than many realize, supporting complex patterns without external dependencies.",
			CodeExample: `mux := http.NewServeMux()

// Static file serving with proper MIME types
mux.HandleFunc("GET /static/", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, r.URL.Path[1:])
})

// RESTful API endpoints with method restrictions
mux.HandleFunc("GET /api/docs", handleGetDocs)
mux.HandleFunc("POST /api/search", handleSearch)
mux.HandleFunc("PUT /api/docs/{id}", handleUpdateDoc)
mux.HandleFunc("DELETE /api/docs/{id}", handleDeleteDoc)

// Path parameters (Go 1.22+)
mux.HandleFunc("GET /api/section/{slug}", func(w http.ResponseWriter, r *http.Request) {
    slug := r.PathValue("slug")
    section := getSection(slug)
    renderJSON(w, section)
})

// Nested path parameters
mux.HandleFunc("GET /api/category/{cat}/section/{slug}", handleCategorySection)

// Wildcard patterns
mux.HandleFunc("GET /files/{path...}", serveUserFiles)

// Pages with optional parameters
mux.HandleFunc("GET /{$}", handleHomePage)  // Exact match for root
mux.HandleFunc("GET /docs/{category}", handleDocsPage)`,
			Category:   "go-backend",
			Order:      1,
			Searchable: "go routing net/http servemux standard library patterns path parameters",
		},
		{
			Title:   "Middleware Architecture",
			Slug:    "middleware-architecture",
			Content: "Build robust middleware chains using Go's standard library. Create reusable, composable middleware for logging, authentication, CORS, and more without external dependencies.",
			CodeExample: `// Middleware type
type Middleware func(http.Handler) http.Handler

// Logging middleware
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Wrap ResponseWriter to capture status code
        wrapped := &ResponseRecorder{ResponseWriter: w, StatusCode: 200}
        
        next.ServeHTTP(wrapped, r)
        
        log.Printf("%s %s %d %v",
            r.Method, r.URL.Path, wrapped.StatusCode, time.Since(start))
    })
}

// CORS middleware
func CORSMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

// Chain middlewares
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        h = middlewares[i](h)
    }
    return h
}

// Usage
handler := Chain(mux, LoggingMiddleware, CORSMiddleware)
http.ListenAndServe(":8080", handler)`,
			Category:   "go-backend",
			Order:      2,
			Searchable: "middleware chain logging authentication CORS standard library",
		},
		{
			Title:   "Template Rendering & Composition",
			Slug:    "template-rendering",
			Content: "Use html/template for type-safe HTML generation with automatic XSS protection and powerful template composition. Create reusable components and layouts without a framework.",
			CodeExample: `//go:embed templates/*.html
var templateFiles embed.FS

// Create template with custom functions
var templates = template.Must(
    template.New("").Funcs(template.FuncMap{
        "formatTime": func(t time.Time) string {
            return t.Format("Jan 2, 2006")
        },
        "truncate": func(s string, length int) string {
            if len(s) <= length {
                return s
            }
            return s[:length] + "..."
        },
        "safeHTML": func(s string) template.HTML {
            return template.HTML(s)
        },
    }).ParseFS(templateFiles, "templates/*.html"))

// Template data structure
type PageData struct {
    Title       string
    User        *User
    Content     interface{}
    CSRFToken   string
    Messages    []Message
}

// Render with layout
func renderPage(w http.ResponseWriter, name string, data PageData) error {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    return templates.ExecuteTemplate(w, name, data)
}

// Render partial for HTMX
func renderPartial(w http.ResponseWriter, name string, data interface{}) error {
    return templates.ExecuteTemplate(w, name, data)
}`,
			Category:   "go-backend",
			Order:      3,
			Searchable: "templates html/template rendering xss protection embed composition",
		},
		{
			Title:   "SQLite Integration & Performance",
			Slug:    "sqlite-integration",
			Content: "Embedded SQLite database with CGO-free driver for true single-binary deployment and zero-latency data access. SQLite is incredibly fast and perfect for most web applications.",
			CodeExample: `import (
    "database/sql"
    "time"
    _ "modernc.org/sqlite"  // CGO-free SQLite driver
)

type DB struct {
    conn *sql.DB
}

func NewDB(dataDir string) (*DB, error) {
    // Enable WAL mode for better concurrency
    dsn := fmt.Sprintf("file:%s/app.db?_journal_mode=WAL&_timeout=5000&_fk=1", dataDir)
    
    conn, err := sql.Open("sqlite", dsn)
    if err != nil {
        return nil, err
    }
    
    // Configure connection pool
    conn.SetMaxOpenConns(25)
    conn.SetMaxIdleConns(25)
    conn.SetConnMaxLifetime(5 * time.Minute)
    
    return &DB{conn: conn}, nil
}

// Optimized query with prepared statements
func (db *DB) GetSectionsByCategory(category string) ([]Section, error) {
    query := ` + "`" + `
    SELECT id, title, slug, content, created_at 
    FROM sections 
    WHERE category = ? 
    ORDER BY order_num, title
    LIMIT 50` + "`" + `
    
    rows, err := db.conn.Query(query, category)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var sections []Section
    for rows.Next() {
        var s Section
        err := rows.Scan(&s.ID, &s.Title, &s.Slug, &s.Content, &s.CreatedAt)
        if err != nil {
            return nil, err
        }
        sections = append(sections, s)
    }
    
    return sections, rows.Err()
}

// Migration system
func (db *DB) Migrate() error {
    migrations := []string{
        ` + "`" + `CREATE TABLE IF NOT EXISTS sections (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            slug TEXT UNIQUE NOT NULL,
            content TEXT NOT NULL,
            category TEXT NOT NULL,
            order_num INTEGER DEFAULT 0,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )` + "`" + `,
        ` + "`" + `CREATE INDEX IF NOT EXISTS idx_sections_category ON sections(category)` + "`" + `,
        ` + "`" + `CREATE INDEX IF NOT EXISTS idx_sections_slug ON sections(slug)` + "`" + `,
    }
    
    for _, migration := range migrations {
        if _, err := db.conn.Exec(migration); err != nil {
            return err
        }
    }
    
    return nil
}`,
			Category:   "go-backend",
			Order:      4,
			Searchable: "sqlite database modernc.org cgo-free single binary performance WAL mode",
		},
		{
			Title:   "Error Handling & Recovery",
			Slug:    "error-handling",
			Content: "Implement robust error handling and recovery mechanisms using Go's standard library. Gracefully handle panics, log errors appropriately, and maintain service availability.",
			CodeExample: `// Custom error types
type AppError struct {
    Code    int
    Message string
    Err     error
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}

// Error handling middleware
func ErrorHandlerMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("PANIC: %v\n%s", err, debug.Stack())
                http.Error(w, "Internal Server Error", 500)
            }
        }()
        
        next.ServeHTTP(w, r)
    })
}

// Error response helper
func writeError(w http.ResponseWriter, err error) {
    var appErr *AppError
    if errors.As(err, &appErr) {
        http.Error(w, appErr.Message, appErr.Code)
        log.Printf("App error: %v", appErr)
    } else {
        http.Error(w, "Internal Server Error", 500)
        log.Printf("Unexpected error: %v", err)
    }
}

// Handler with proper error handling
func handleGetSection(w http.ResponseWriter, r *http.Request) {
    slug := r.PathValue("slug")
    if slug == "" {
        writeError(w, &AppError{
            Code:    400,
            Message: "Section slug is required",
        })
        return
    }
    
    section, err := sectionService.GetBySlug(slug)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            writeError(w, &AppError{
                Code:    404,
                Message: "Section not found",
            })
        } else {
            writeError(w, err)
        }
        return
    }
    
    if err := renderPage(w, "section.html", section); err != nil {
        writeError(w, err)
    }
}`,
			Category:   "go-backend",
			Order:      5,
			Searchable: "error handling recovery panic middleware logging standard library",
		},
		{
			Title:   "Asset Embedding & Static Files",
			Slug:    "asset-embedding",
			Content: "Embed all static assets directly into the Go binary using go:embed for simplified deployment and enhanced security. No external file dependencies, no CDN requirements.",
			CodeExample: `import (
    "crypto/sha256"
    "embed"
    "fmt"
    "net/http"
    "strings"
)

//go:embed static/htmx.min.js
var htmxJS []byte

//go:embed static/styles.css  
var stylesCSS []byte

//go:embed static/fonts/*.woff2
var fontFiles embed.FS

//go:embed templates/*
var templates embed.FS

// Serve embedded static files with proper caching
func serveStaticFile(w http.ResponseWriter, r *http.Request, content []byte, contentType string) {
    w.Header().Set("Content-Type", contentType)
    w.Header().Set("Cache-Control", "public, max-age=31536000") // 1 year
    hash := sha256.Sum256(content)
    w.Header().Set("ETag", fmt.Sprintf("\"%x\"", hash[:]))
    
    // Check if client has cached version
    if match := r.Header.Get("If-None-Match"); match != "" {
        if strings.Contains(match, w.Header().Get("ETag")) {
            w.WriteHeader(http.StatusNotModified)
            return
        }
    }
    
    w.Write(content)
}

// Setup static file routes
func setupStaticRoutes(mux *http.ServeMux) {
    mux.HandleFunc("GET /static/htmx.min.js", func(w http.ResponseWriter, r *http.Request) {
        serveStaticFile(w, r, htmxJS, "application/javascript")
    })
    
    mux.HandleFunc("GET /static/styles.css", func(w http.ResponseWriter, r *http.Request) {
        serveStaticFile(w, r, stylesCSS, "text/css")
    })
    
    // Serve embedded font files
    mux.Handle("GET /static/fonts/", http.FileServer(http.FS(fontFiles)))
}`,
			Category:   "deployment",
			Order:      1,
			Searchable: "asset embedding go:embed static files deployment caching fonts",
		},
		{
			Title:   "Single Binary Deployment",
			Slug:    "single-binary-deployment",
			Content: "Deploy your entire application as a single executable binary. Perfect for Ubuntu servers behind Caddy reverse proxy. No runtime dependencies, no installation scripts, no complex deployment pipelines. Just copy and run.",
			CodeExample: `# Build for production
go build -ldflags="-s -w" -o myapp main.go

# Cross-compilation for different platforms
GOOS=linux GOARCH=amd64 go build -o myapp-linux-amd64 main.go
GOOS=windows GOARCH=amd64 go build -o myapp-windows-amd64.exe main.go
GOOS=darwin GOARCH=arm64 go build -o myapp-darwin-arm64 main.go

# Ultra-minimal Docker image
FROM scratch
COPY myapp /
COPY data/ /data/
EXPOSE 8080
ENTRYPOINT ["/myapp"]

# Or use distroless for debugging capabilities
FROM gcr.io/distroless/static-debian11
COPY myapp /
COPY data/ /data/
EXPOSE 8080
ENTRYPOINT ["/myapp"]

# Ubuntu + Caddy deployment example
# 1. Upload binary to Ubuntu server
scp myapp-linux-amd64 user@server:/opt/myapp/myapp

# 2. Create systemd service
sudo systemctl enable myapp
sudo systemctl start myapp

# 3. Caddy reverse proxy (Caddyfile)
myapp.example.com {
    reverse_proxy localhost:8080
}

# Dockerfile is only 2 lines + your binary!
# Final image size: ~15MB (including your app)`,
			Category:   "deployment",
			Order:      2,
			Searchable: "single binary deployment docker cross-compilation production",
		},
		{
			Title:   "Production Configuration",
			Slug:    "production-config",
			Content: "Configure your Go application for production with proper logging, graceful shutdown, TLS support, and environment-based configuration.",
			CodeExample: `// Production server configuration
func main() {
    // Environment-based configuration
    port := getEnv("PORT", "8080")
    dbPath := getEnv("DB_PATH", "./data/app.db")
    tlsCert := getEnv("TLS_CERT", "")
    tlsKey := getEnv("TLS_KEY", "")
    
    // Setup structured logging for production
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))
    slog.SetDefault(logger)
    
    // Initialize database with connection pooling
    db, err := database.NewDB(dbPath)
    if err != nil {
        slog.Error("Failed to connect to database", "error", err)
        os.Exit(1)
    }
    defer db.Close()
    
    // Setup handlers with middleware
    handler := setupMiddleware(setupRoutes(db))
    
    // Configure server for production
    server := &http.Server{
        Addr:         ":" + port,
        Handler:      handler,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
        ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
    }
    
    // Graceful shutdown
    go func() {
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
        <-sigChan
        
        slog.Info("Shutting down server...")
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        
        if err := server.Shutdown(ctx); err != nil {
            slog.Error("Server shutdown failed", "error", err)
        }
    }()
    
    slog.Info("Starting server", "port", port, "tls", tlsCert != "")
    
    // Start server with optional TLS
    if tlsCert != "" && tlsKey != "" {
        err = server.ListenAndServeTLS(tlsCert, tlsKey)
    } else {
        err = server.ListenAndServe()
    }
    
    if err != nil && err != http.ErrServerClosed {
        slog.Error("Server failed", "error", err)
        os.Exit(1)
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}`,
			Category:   "deployment",
			Order:      3,
			Searchable: "production configuration logging graceful shutdown TLS environment",
		},
		{
			Title:      "Performance Optimization",
			Slug:       "performance-optimization",
			Content:    "Optimize your Go hypermedia application for maximum performance. Learn about connection pooling, caching strategies, and Go-specific optimizations.",
			Category:   "optimization",
			Order:      1,
			Searchable: "performance optimization caching connection pooling go specific",
			CodeExample: `// Response caching middleware
func CacheMiddleware(duration time.Duration) func(http.Handler) http.Handler {
    type CacheEntry struct {
        Data      []byte
        Headers   http.Header
        ExpiresAt time.Time
    }
    
    cache := sync.Map{}
    
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Only cache GET requests
            if r.Method != "GET" {
                next.ServeHTTP(w, r)
                return
            }
            
            key := r.URL.Path + "?" + r.URL.RawQuery
            
            // Check cache
            if entry, ok := cache.Load(key); ok {
                cached := entry.(CacheEntry)
                if time.Now().Before(cached.ExpiresAt) {
                    // Serve from cache
                    for k, v := range cached.Headers {
                        w.Header()[k] = v
                    }
                    w.Header().Set("X-Cache", "HIT")
                    w.Write(cached.Data)
                    return
                }
                // Expired, remove from cache
                cache.Delete(key)
            }
            
            // Capture response
            recorder := &ResponseRecorder{
                ResponseWriter: w,
                Body:          bytes.NewBuffer(nil),
                Headers:       make(http.Header),
            }
            
            next.ServeHTTP(recorder, r)
            
            // Cache successful responses
            if recorder.StatusCode == 200 {
                cache.Store(key, CacheEntry{
                    Data:      recorder.Body.Bytes(),
                    Headers:   recorder.Headers,
                    ExpiresAt: time.Now().Add(duration),
                })
            }
            
            w.Header().Set("X-Cache", "MISS")
            w.Write(recorder.Body.Bytes())
        })
    }
}

// Connection pooling for external services
var httpClient = &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxConnsPerHost:     10,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}

// Memory-efficient JSON streaming
func streamJSONResponse(w http.ResponseWriter, data interface{}) error {
    w.Header().Set("Content-Type", "application/json")
    encoder := json.NewEncoder(w)
    return encoder.Encode(data)
}`,
		},
		{
			Title:      "Security Best Practices",
			Slug:       "security-best-practices",
			Content:    "Implement comprehensive security measures for your Go web application. Learn about CSRF protection, secure headers, input validation, and more.",
			Category:   "security",
			Order:      1,
			Searchable: "security CSRF XSS input validation secure headers authentication",
			CodeExample: `// Security middleware
func SecurityMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Security headers
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        w.Header().Set("Content-Security-Policy", 
            "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")
        
        next.ServeHTTP(w, r)
    })
}

// CSRF protection
type CSRFToken struct {
    Value     string
    ExpiresAt time.Time
}

var csrfTokens = sync.Map{}

func generateCSRFToken() string {
    token := make([]byte, 32)
    rand.Read(token)
    return base64.URLEncoding.EncodeToString(token)
}

func CSRFMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "POST" || r.Method == "PUT" || r.Method == "DELETE" {
            token := r.Header.Get("X-CSRF-Token")
            if token == "" {
                token = r.FormValue("csrf_token")
            }
            
            if !validateCSRFToken(token) {
                http.Error(w, "Invalid CSRF token", http.StatusForbidden)
                return
            }
        }
        
        next.ServeHTTP(w, r)
    })
}

// Input validation
func ValidateInput(input string, maxLength int, allowedChars string) error {
    if len(input) > maxLength {
        return fmt.Errorf("input too long: max %d characters", maxLength)
    }
    
    if matched, _ := regexp.MatchString(allowedChars, input); !matched {
        return fmt.Errorf("input contains invalid characters")
    }
    
    return nil
}

// Rate limiting
func RateLimitMiddleware(requestsPerMinute int) func(http.Handler) http.Handler {
    type client struct {
        lastSeen time.Time
        requests int
    }
    
    clients := sync.Map{}
    
    // Cleanup old entries
    go func() {
        for {
            time.Sleep(time.Minute)
            clients.Range(func(key, value interface{}) bool {
                client := value.(client)
                if time.Since(client.lastSeen) > time.Hour {
                    clients.Delete(key)
                }
                return true
            })
        }
    }()
    
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ip := r.RemoteAddr
            now := time.Now()
            
            value, _ := clients.LoadOrStore(ip, client{lastSeen: now, requests: 0})
            c := value.(client)
            
            // Reset counter if more than a minute has passed
            if now.Sub(c.lastSeen) > time.Minute {
                c.requests = 0
            }
            
            c.requests++
            c.lastSeen = now
            clients.Store(ip, c)
            
            if c.requests > requestsPerMinute {
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}`,
		},
	}

	for _, section := range sections {
		_, err := s.db.GetConnection().Exec(`
			INSERT INTO doc_sections (title, slug, content, code_example, category, order_num, searchable)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, section.Title, section.Slug, section.Content, section.CodeExample,
			section.Category, section.Order, section.Searchable)

		if err != nil {
			return fmt.Errorf("failed to insert section %s: %w", section.Title, err)
		}
	}

	return nil
}
