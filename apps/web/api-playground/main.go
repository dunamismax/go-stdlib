package main

import (
	"crypto/rand"
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//go:embed static/htmx.min.js
var htmxJS []byte

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

var stylesCSS = `
body {
	font-family: system-ui, -apple-system, sans-serif;
	background-color: #0f0f0f;
	color: #e4e4e4;
	margin: 0;
	padding: 1rem;
	line-height: 1.6;
}

.container {
	max-width: 1200px;
	margin: 0 auto;
	padding: 2rem;
}

.card {
	background: #1a1a1a;
	border-radius: 8px;
	padding: 1.5rem;
	margin: 1rem 0;
	border: 1px solid #333;
}

input, textarea, select {
	background: #2a2a2a;
	border: 1px solid #444;
	color: #e4e4e4;
	padding: 0.75rem;
	border-radius: 4px;
	font-size: 1rem;
	width: 100%;
	margin: 0.5rem 0;
}

button {
	background: #2563eb;
	color: white;
	border: none;
	padding: 0.75rem 1.5rem;
	border-radius: 4px;
	cursor: pointer;
	font-size: 1rem;
	margin: 0.5rem 0.5rem 0.5rem 0;
}

button:hover {
	background: #1d4ed8;
}

.result {
	background: #2a2a2a;
	border-radius: 4px;
	padding: 1rem;
	margin: 1rem 0;
	border: 1px solid #444;
	word-wrap: break-word;
}

.grid {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
	gap: 1rem;
	margin: 1rem 0;
}
`

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
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
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

func (a *App) serveHTMX(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	w.Write(htmxJS)
}

func (a *App) serveCSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	w.Write([]byte(stylesCSS))
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s %v", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}

func main() {
	app := NewApp()

	mux := http.NewServeMux()
	
	// Static files
	mux.HandleFunc("/static/htmx.min.js", app.serveHTMX)
	mux.HandleFunc("/static/styles.css", app.serveCSS)
	
	// Routes
	mux.HandleFunc("/", app.homePage)
	mux.HandleFunc("/analyze", app.analyzeText)
	mux.HandleFunc("/reverse", app.reverseText)
	mux.HandleFunc("/titlecase", app.titleCase)
	mux.HandleFunc("/count-vowels", app.countVowels)
	mux.HandleFunc("/random-number", app.randomNumber)
	mux.HandleFunc("/random-string", app.randomString)
	mux.HandleFunc("/uuid", app.generateUUID)
	mux.HandleFunc("/sha256", app.sha256Hash)
	mux.HandleFunc("/base64-encode", app.base64Encode)
	mux.HandleFunc("/base64-decode", app.base64Decode)

	handler := loggerMiddleware(mux)

	fmt.Println("API Playground starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}