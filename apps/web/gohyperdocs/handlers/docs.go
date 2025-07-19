package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/dunamismax/go-stdlib/apps/web/gohyperdocs/models"
	"github.com/dunamismax/go-stdlib/pkg/utils"
)

// Demo state for counters
var (
	visitorCount int64
)

type DocsHandler struct {
	docsService *models.DocsService
	templates   *template.Template
}

func NewDocsHandler(docsService *models.DocsService, templates *template.Template) *DocsHandler {
	return &DocsHandler{
		docsService: docsService,
		templates:   templates,
	}
}

func (h *DocsHandler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	sections, err := h.docsService.GetAllSections()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to load documentation")
		return
	}

	// Group sections by category
	categories := make(map[string][]models.DocSection)
	for _, section := range sections {
		categories[section.Category] = append(categories[section.Category], section)
	}

	data := map[string]interface{}{
		"Title":      "GoHyperDocs - The Ultimate Go Hypermedia Documentation",
		"Categories": categories,
		"Sections":   sections,
	}

	w.Header().Set("Content-Type", "text/html")
	if err := h.templates.ExecuteTemplate(w, "home.html", data); err != nil {
		utils.Error(w, http.StatusInternalServerError, fmt.Sprintf("Failed to render template: %v", err))
		return
	}
}

func (h *DocsHandler) SectionHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		utils.Error(w, http.StatusBadRequest, "Section slug is required")
		return
	}

	section, err := h.docsService.GetSectionBySlug(slug)
	if err != nil {
		utils.Error(w, http.StatusNotFound, "Section not found")
		return
	}

	// Get all sections for navigation
	allSections, err := h.docsService.GetAllSections()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to load navigation")
		return
	}

	// Group sections by category for navigation
	categories := make(map[string][]models.DocSection)
	for _, s := range allSections {
		categories[s.Category] = append(categories[s.Category], s)
	}

	data := map[string]interface{}{
		"Title":      fmt.Sprintf("%s - GoHyperDocs", section.Title),
		"Section":    section,
		"Categories": categories,
	}

	w.Header().Set("Content-Type", "text/html")
	if err := h.templates.ExecuteTemplate(w, "section.html", data); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to render template")
		return
	}
}

func (h *DocsHandler) CategoryHandler(w http.ResponseWriter, r *http.Request) {
	category := r.PathValue("category")
	if category == "" {
		utils.Error(w, http.StatusBadRequest, "Category is required")
		return
	}

	sections, err := h.docsService.GetSectionsByCategory(category)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to load category sections")
		return
	}

	// Get all sections for navigation
	allSections, err := h.docsService.GetAllSections()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to load navigation")
		return
	}

	// Group sections by category for navigation
	categories := make(map[string][]models.DocSection)
	for _, s := range allSections {
		categories[s.Category] = append(categories[s.Category], s)
	}

	data := map[string]interface{}{
		"Title":      fmt.Sprintf("%s - GoHyperDocs", strings.Title(strings.ReplaceAll(category, "-", " "))),
		"Category":   category,
		"Sections":   sections,
		"Categories": categories,
	}

	w.Header().Set("Content-Type", "text/html")
	if err := h.templates.ExecuteTemplate(w, "category.html", data); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to render template")
		return
	}
}

func (h *DocsHandler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		// Return empty results template
		data := map[string]interface{}{
			"Query":   "",
			"Results": []models.DocSection{},
		}
		w.Header().Set("Content-Type", "text/html")
		h.templates.ExecuteTemplate(w, "search-results.html", data)
		return
	}

	sections, err := h.docsService.SearchSections(query)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Search failed")
		return
	}

	data := map[string]interface{}{
		"Query":   query,
		"Results": sections,
	}

	w.Header().Set("Content-Type", "text/html")
	if err := h.templates.ExecuteTemplate(w, "search-results.html", data); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to render search results")
		return
	}
}

func (h *DocsHandler) LazyLoadDetailsHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		utils.Error(w, http.StatusBadRequest, "Section slug is required")
		return
	}

	section, err := h.docsService.GetSectionBySlug(slug)
	if err != nil {
		utils.Error(w, http.StatusNotFound, "Section not found")
		return
	}

	data := map[string]interface{}{
		"Section": section,
	}

	w.Header().Set("Content-Type", "text/html")
	if err := h.templates.ExecuteTemplate(w, "section-details.html", data); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to render details")
		return
	}
}

func (h *DocsHandler) CopyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// This endpoint is called when copy is successful
	// Return updated button state
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<button class="copy-btn copied" disabled>Copied!</button>`)
}

// Interactive demonstration handlers
func (h *DocsHandler) LiveCounterHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate a live counter that increments
	counterValue := 42 + (time.Now().Unix() % 100)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
	<div class="live-demo-card">
		<h4>Live Server Stats</h4>
		<div class="stats-grid">
			<div class="stat-item">
				<span class="stat-label">Active Users</span>
				<span class="stat-value">%d</span>
			</div>
			<div class="stat-item">
				<span class="stat-label">Requests/min</span>
				<span class="stat-value">%d</span>
			</div>
			<div class="stat-item">
				<span class="stat-label">Response Time</span>
				<span class="stat-value">%.1fms</span>
			</div>
		</div>
		<div class="update-timestamp">Updated: %s</div>
	</div>`,
		counterValue,
		150+(time.Now().Unix()%50),
		0.2+float64(time.Now().Unix()%10)/10.0,
		time.Now().Format("15:04:05"))
}

func (h *DocsHandler) FormValidationHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	if email == "" {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<div class="validation-message"></div>`)
		return
	}

	// Simple email validation
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)

	w.Header().Set("Content-Type", "text/html")
	if matched {
		fmt.Fprint(w, `<div class="validation-message success">Valid email address</div>`)
	} else {
		fmt.Fprint(w, `<div class="validation-message error">Please enter a valid email address</div>`)
	}
}

func (h *DocsHandler) LoadMoreContentHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate loading additional content
	time.Sleep(500 * time.Millisecond) // Simulate network delay

	items := []string{
		"Blazing fast server-side rendering",
		"Zero-dependency JavaScript interactivity",
		"Built-in XSS protection with html/template",
		"SQLite for zero-latency data access",
		"Single binary deployment",
		"Progressive enhancement by design",
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<div class="loaded-content">`)
	for i, item := range items {
		fmt.Fprintf(w, `<div class="feature-item" style="animation-delay: %dms">%s</div>`, i*100, item)
	}
	fmt.Fprint(w, `</div>`)
}

func (h *DocsHandler) TabDemoHandler(w http.ResponseWriter, r *http.Request) {
	tab := r.URL.Query().Get("tab")
	if tab == "" {
		tab = "overview"
	}

	content := map[string]string{
		"overview": `
		<div class="tab-content">
			<h3>Go Hypermedia Stack Overview</h3>
			<p>A revolutionary approach to web development that combines the power of Go's standard library with the elegance of hypermedia-driven interfaces.</p>
			<ul>
				<li>Zero external dependencies</li>
				<li>Single binary deployment</li>
				<li>Progressive enhancement</li>
				<li>Maximum performance</li>
			</ul>
		</div>`,
		"features": `
		<div class="tab-content">
			<h3>Key Features</h3>
			<div class="feature-grid">
				<div class="feature-box">
					<h4>Lightning Fast</h4>
					<p>Sub-millisecond response times with compiled Go code</p>
				</div>
				<div class="feature-box">
					<h4>Secure by Default</h4>
					<p>Built-in XSS protection and secure templating</p>
				</div>
				<div class="feature-box">
					<h4>Self-Contained</h4>
					<p>Everything embedded in a single binary</p>
				</div>
			</div>
		</div>`,
		"performance": `
		<div class="tab-content">
			<h3>Performance Metrics</h3>
			<div class="metrics">
				<div class="metric">
					<span class="metric-value">0.2ms</span>
					<span class="metric-label">Avg Response Time</span>
				</div>
				<div class="metric">
					<span class="metric-value">15MB</span>
					<span class="metric-label">Memory Usage</span>
				</div>
				<div class="metric">
					<span class="metric-value">10K+</span>
					<span class="metric-label">Concurrent Users</span>
				</div>
			</div>
		</div>`,
	}

	w.Header().Set("Content-Type", "text/html")
	if tabContent, exists := content[tab]; exists {
		fmt.Fprint(w, tabContent)
	} else {
		fmt.Fprint(w, content["overview"])
	}
}

func (h *DocsHandler) APIDocsHandler(w http.ResponseWriter, r *http.Request) {
	sections, err := h.docsService.GetAllSections()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to load documentation")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"sections": sections,
	})
}

func (h *DocsHandler) APISectionHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		utils.Error(w, http.StatusBadRequest, "Section slug is required")
		return
	}

	section, err := h.docsService.GetSectionBySlug(slug)
	if err != nil {
		utils.Error(w, http.StatusNotFound, "Section not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(section)
}

func (h *DocsHandler) ServeStatic(w http.ResponseWriter, r *http.Request, content []byte, contentType string) {
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=31536000") // Cache for 1 year
	w.Write(content)
}

// Enhanced Live Demo Handlers

// Enhanced Todo Demo Handler
func (h *DocsHandler) TodoDemoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Return the enhanced todo form and list
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `
		<div class="todo-demo">
			<form hx-post="/api/demo/todo" hx-target="#todo-list" hx-swap="beforeend" hx-on::after-request="this.reset()">
				<input name="task" type="text" placeholder="Add a new task..." required maxlength="100">
				<button type="submit">Add Task</button>
			</form>
			<div id="todo-list" class="todo-list">
				<div class="todo-item">
					<span class="todo-text">Example: Learn HTMX hypermedia patterns</span>
					<button class="todo-delete" hx-delete="/api/demo/todo/1" hx-target="closest .todo-item" hx-swap="outerHTML">×</button>
				</div>
				<div class="todo-item">
					<span class="todo-text">Example: Build with Go standard library</span>
					<button class="todo-delete" hx-delete="/api/demo/todo/2" hx-target="closest .todo-item" hx-swap="outerHTML">×</button>
				</div>
			</div>
		</div>`)

	case "POST":
		task := r.FormValue("task")
		if task == "" {
			utils.Error(w, http.StatusBadRequest, "Task is required")
			return
		}

		// Generate a random ID for this demo
		id := time.Now().UnixNano()

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
		<div class="todo-item" style="animation: slideIn 0.3s ease-out">
			<span class="todo-text">%s</span>
			<button class="todo-delete" hx-delete="/api/demo/todo/%d" hx-target="closest .todo-item" hx-swap="outerHTML">×</button>
		</div>`, template.HTMLEscapeString(task), id)

	case "DELETE":
		// Just return empty content to remove the item
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(""))
	}
}
