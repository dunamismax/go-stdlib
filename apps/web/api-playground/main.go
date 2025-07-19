package main

import (
	"crypto/rand"
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//go:embed dist/*
var distFS embed.FS

//go:embed templates/index.html
var indexHTML string

//go:embed templates/result.html
var resultHTML string

var (
	indexTemplate  *template.Template
	resultTemplate *template.Template
)

func init() {
	indexTemplate = template.Must(template.New("index").Parse(indexHTML))
	resultTemplate = template.Must(template.New("result").Parse(resultHTML))
}

type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	indexTemplate.Execute(w, nil)
}

func (a *App) analyzeText(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")

	chars := len(text)
	words := len(strings.Fields(text))
	sentences := strings.Count(text, ".") + strings.Count(text, "!") + strings.Count(text, "?")
	paragraphs := len(strings.Split(text, "\n\n"))

	data := map[string]interface{}{
		"Title":      "Analysis Results:",
		"Characters": chars,
		"Words":      words,
		"Sentences":  sentences,
		"Paragraphs": paragraphs,
	}

	w.Header().Set("Content-Type", "text/html")
	resultTemplate.Execute(w, data)
}

func (a *App) reverseText(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")

	runes := []rune(text)
	n := len(runes)
	for i := range n / 2 {
		runes[i], runes[n-1-i] = runes[n-1-i], runes[i]
	}

	data := map[string]interface{}{
		"Title":  "Reversed Text:",
		"Result": string(runes),
	}

	w.Header().Set("Content-Type", "text/html")
	resultTemplate.Execute(w, data)
}

func (a *App) titleCase(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")

	data := map[string]interface{}{
		"Title":  "Title Case:",
		"Result": strings.Title(text),
	}

	w.Header().Set("Content-Type", "text/html")
	resultTemplate.Execute(w, data)
}

func (a *App) countVowels(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")

	vowels := "aeiouAEIOU"
	count := 0
	for _, char := range text {
		if strings.ContainsRune(vowels, char) {
			count++
		}
	}

	data := map[string]interface{}{
		"Title":  "Vowel Count:",
		"Result": fmt.Sprintf("Vowels: %d", count),
	}

	w.Header().Set("Content-Type", "text/html")
	resultTemplate.Execute(w, data)
}

func (a *App) randomNumber(w http.ResponseWriter, r *http.Request) {
	minStr := r.FormValue("min")
	maxStr := r.FormValue("max")

	min, _ := strconv.Atoi(minStr)
	max, _ := strconv.Atoi(maxStr)

	if min > max {
		min, max = max, min
	}

	// Generate random number using crypto/rand for better randomness
	rangeSize := max - min + 1
	b := make([]byte, 8)
	rand.Read(b)
	randInt := int(b[0])<<56 | int(b[1])<<48 | int(b[2])<<40 | int(b[3])<<32 |
		int(b[4])<<24 | int(b[5])<<16 | int(b[6])<<8 | int(b[7])
	if randInt < 0 {
		randInt = -randInt
	}
	num := (randInt % rangeSize) + min

	data := map[string]interface{}{
		"Title":  "Random Number:",
		"Result": fmt.Sprintf("%d", num),
	}

	w.Header().Set("Content-Type", "text/html")
	resultTemplate.Execute(w, data)
}

func (a *App) randomString(w http.ResponseWriter, r *http.Request) {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 10

	result := make([]byte, length)
	randomBytes := make([]byte, length)
	rand.Read(randomBytes)

	for i := range result {
		result[i] = chars[int(randomBytes[i])%len(chars)]
	}

	data := map[string]interface{}{
		"Title":  "Random String:",
		"Result": string(result),
	}

	w.Header().Set("Content-Type", "text/html")
	resultTemplate.Execute(w, data)
}

func (a *App) generateUUID(w http.ResponseWriter, r *http.Request) {
	// Generate UUID using crypto/rand
	uuid := make([]byte, 16)
	rand.Read(uuid)

	// Set version (4) and variant bits
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant bits

	uuidStr := fmt.Sprintf("%x-%x-%x-%x-%x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16])

	data := map[string]interface{}{
		"Title":  "UUID:",
		"Result": uuidStr,
	}

	w.Header().Set("Content-Type", "text/html")
	resultTemplate.Execute(w, data)
}

func (a *App) sha256Hash(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")

	hash := sha256.Sum256([]byte(text))

	data := map[string]interface{}{
		"Title":  "SHA256 Hash:",
		"Result": fmt.Sprintf("%x", hash),
	}

	w.Header().Set("Content-Type", "text/html")
	resultTemplate.Execute(w, data)
}

func (a *App) base64Encode(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")

	encoded := base64.StdEncoding.EncodeToString([]byte(text))

	data := map[string]interface{}{
		"Title":  "Base64 Encoded:",
		"Result": encoded,
	}

	w.Header().Set("Content-Type", "text/html")
	resultTemplate.Execute(w, data)
}

func (a *App) base64Decode(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")

	decoded, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		data := map[string]interface{}{
			"Title":  "Error:",
			"Result": "Invalid base64 string",
		}
		w.Header().Set("Content-Type", "text/html")
		resultTemplate.Execute(w, data)
		return
	}

	data := map[string]interface{}{
		"Title":  "Base64 Decoded:",
		"Result": string(decoded),
	}

	w.Header().Set("Content-Type", "text/html")
	resultTemplate.Execute(w, data)
}

// Fun APIs
var jokes = []string{
	"Why don't scientists trust atoms? Because they make up everything!",
	"Why did the programmer quit his job? He didn't get arrays!",
	"How many programmers does it take to change a light bulb? None, that's a hardware problem!",
	"Why do Java developers wear glasses? Because they can't C#!",
	"What's the object-oriented way to become wealthy? Inheritance!",
	"Why did the developer go broke? Because he used up all his cache!",
	"What do you call a programmer from Finland? Nerdic!",
	"Why do programmers prefer dark mode? Because light attracts bugs!",
	"What's a programmer's favorite hangout place? Foo Bar!",
	"Why don't programmers like nature? It has too many bugs!",
}

var quotes = []string{
	"The only way to do great work is to love what you do. - Steve Jobs",
	"Code is like humor. When you have to explain it, it's bad. - Cory House",
	"Programs must be written for people to read, and only incidentally for machines to execute. - Harold Abelson",
	"The best error message is the one that never shows up. - Thomas Fuchs",
	"Talk is cheap. Show me the code. - Linus Torvalds",
	"Any fool can write code that a computer can understand. Good programmers write code that humans can understand. - Martin Fowler",
	"First, solve the problem. Then, write the code. - John Johnson",
	"Experience is the name everyone gives to their mistakes. - Oscar Wilde",
	"The most important property of a program is whether it accomplishes the intention of its user. - C.A.R. Hoare",
	"Simplicity is the ultimate sophistication. - Leonardo da Vinci",
}

var facts = []string{
	"The first computer bug was an actual bug found in a Harvard Mark II computer in 1947.",
	"The term 'debugging' was coined by Admiral Grace Hopper.",
	"The first programming language was created in 1843 by Ada Lovelace.",
	"COBOL was one of the first programming languages and is still used today in banking systems.",
	"The @ symbol was used in email addresses because it was the only preposition available on the keyboard.",
	"The first computer virus was created in 1971 and was called 'The Creeper'.",
	"Java was originally called Oak, but had to be renamed due to trademark issues.",
	"The term 'WiFi' doesn't actually stand for anything - it's just a brand name.",
	"The first 1GB hard drive cost $40,000 and weighed over 500 pounds.",
	"Python was named after Monty Python's Flying Circus, not the snake.",
}

func (a *App) randomJoke(w http.ResponseWriter, r *http.Request) {
	// Generate random index
	b := make([]byte, 1)
	rand.Read(b)
	index := int(b[0]) % len(jokes)

	data := map[string]interface{}{
		"Title":  "Random Joke:",
		"Result": jokes[index],
	}

	w.Header().Set("Content-Type", "text/html")
	resultTemplate.Execute(w, data)
}

func (a *App) randomQuote(w http.ResponseWriter, r *http.Request) {
	// Generate random index
	b := make([]byte, 1)
	rand.Read(b)
	index := int(b[0]) % len(quotes)

	data := map[string]interface{}{
		"Title":  "Inspirational Quote:",
		"Result": quotes[index],
	}

	w.Header().Set("Content-Type", "text/html")
	resultTemplate.Execute(w, data)
}

func (a *App) randomFact(w http.ResponseWriter, r *http.Request) {
	// Generate random index
	b := make([]byte, 1)
	rand.Read(b)
	index := int(b[0]) % len(facts)

	data := map[string]interface{}{
		"Title":  "Random Tech Fact:",
		"Result": facts[index],
	}

	w.Header().Set("Content-Type", "text/html")
	resultTemplate.Execute(w, data)
}

// Time and Date APIs
func (a *App) currentTime(w http.ResponseWriter, r *http.Request) {
	timezone := r.FormValue("timezone")
	if timezone == "" {
		timezone = "UTC"
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		data := map[string]interface{}{
			"Title":  "Error:",
			"Result": "Invalid timezone: " + timezone,
		}
		w.Header().Set("Content-Type", "text/html")
		resultTemplate.Execute(w, data)
		return
	}

	now := time.Now().In(loc)
	formatted := now.Format("Monday, January 2, 2006 at 3:04:05 PM MST")

	data := map[string]interface{}{
		"Title":  "Current Time (" + timezone + "):",
		"Result": formatted,
	}

	w.Header().Set("Content-Type", "text/html")
	resultTemplate.Execute(w, data)
}

func (a *App) unixTimestamp(w http.ResponseWriter, r *http.Request) {
	timezone := r.FormValue("timezone")
	if timezone == "" {
		timezone = "UTC"
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		data := map[string]interface{}{
			"Title":  "Error:",
			"Result": "Invalid timezone: " + timezone,
		}
		w.Header().Set("Content-Type", "text/html")
		resultTemplate.Execute(w, data)
		return
	}

	now := time.Now().In(loc)
	timestamp := now.Unix()

	data := map[string]interface{}{
		"Title":  "Unix Timestamp:",
		"Result": fmt.Sprintf("%d", timestamp),
	}

	w.Header().Set("Content-Type", "text/html")
	resultTemplate.Execute(w, data)
}

func (a *App) dateFormats(w http.ResponseWriter, r *http.Request) {
	now := time.Now().UTC()

	formats := []string{
		"ISO 8601: " + now.Format("2006-01-02T15:04:05Z"),
		"RFC 3339: " + now.Format(time.RFC3339),
		"US Format: " + now.Format("01/02/2006 03:04 PM"),
		"EU Format: " + now.Format("02/01/2006 15:04"),
		"Long Format: " + now.Format("Monday, January 2, 2006"),
		"Short Date: " + now.Format("Jan 2, 2006"),
		"Time Only: " + now.Format("15:04:05"),
		"12H Time: " + now.Format("3:04:05 PM"),
	}

	result := ""
	for _, format := range formats {
		result += format + "\n"
	}

	data := map[string]interface{}{
		"Title":  "Date & Time Formats (UTC):",
		"Result": result,
	}

	w.Header().Set("Content-Type", "text/html")
	resultTemplate.Execute(w, data)
}

func (a *App) serveStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/static/" {
		path = "/static/index.html"
	}

	// Remove /static/ prefix and add dist/ prefix
	filePath := "dist" + strings.TrimPrefix(path, "/static")

	// Read file from embedded filesystem
	data, err := distFS.ReadFile(filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Set content type based on file extension
	if strings.HasSuffix(filePath, ".css") {
		w.Header().Set("Content-Type", "text/css")
	} else if strings.HasSuffix(filePath, ".js") {
		w.Header().Set("Content-Type", "application/javascript")
	} else if strings.HasSuffix(filePath, ".html") {
		w.Header().Set("Content-Type", "text/html")
	}

	// Set cache headers for static assets
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	w.Write(data)
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

func main() {
	// Setup structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	app := NewApp()

	mux := http.NewServeMux()

	// Static files (Vite build output)
	mux.HandleFunc("GET /static/", app.serveStatic)

	// Routes
	mux.HandleFunc("GET /", app.homePage)
	mux.HandleFunc("POST /analyze", app.analyzeText)
	mux.HandleFunc("POST /reverse", app.reverseText)
	mux.HandleFunc("POST /titlecase", app.titleCase)
	mux.HandleFunc("POST /count-vowels", app.countVowels)
	mux.HandleFunc("POST /random-number", app.randomNumber)
	mux.HandleFunc("POST /random-string", app.randomString)
	mux.HandleFunc("POST /uuid", app.generateUUID)
	mux.HandleFunc("POST /sha256", app.sha256Hash)
	mux.HandleFunc("POST /base64-encode", app.base64Encode)
	mux.HandleFunc("POST /base64-decode", app.base64Decode)

	// Fun APIs
	mux.HandleFunc("GET /joke", app.randomJoke)
	mux.HandleFunc("GET /quote", app.randomQuote)
	mux.HandleFunc("GET /fact", app.randomFact)

	// Time and Date APIs
	mux.HandleFunc("POST /current-time", app.currentTime)
	mux.HandleFunc("POST /timestamp", app.unixTimestamp)
	mux.HandleFunc("GET /date-formats", app.dateFormats)

	handler := loggerMiddleware(mux)

	slog.Info("API Playground starting", "port", ":8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		slog.Error("Server failed to start", "error", err)
		log.Fatal(err)
	}
}
