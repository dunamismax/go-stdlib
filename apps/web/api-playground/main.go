package main

import (
	"crypto/rand"
	"crypto/sha256"
	"embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed dist/*
var distFS embed.FS

//go:embed templates/index.html
var indexHTML string

//go:embed templates/result.html
var resultHTML string

type TemplateRenderer struct {
	templates map[string]*template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Template not found")
	}
	return tmpl.Execute(w, data)
}

func NewTemplateRenderer() *TemplateRenderer {
	return &TemplateRenderer{
		templates: map[string]*template.Template{
			"index":  template.Must(template.New("index").Parse(indexHTML)),
			"result": template.Must(template.New("result").Parse(resultHTML)),
		},
	}
}

type App struct {
	echo *echo.Echo
}

func NewApp() *App {
	e := echo.New()
	
	// Configure Echo
	e.HideBanner = true
	e.Renderer = NewTemplateRenderer()
	
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Secure())
	e.Use(middleware.RequestID())
	
	// Custom middleware for structured logging
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${id} ${remote_ip} ${host} ${method} ${uri} ${user_agent} ${status} ${error} ${latency} ${latency_human} ${bytes_in} ${bytes_out}\n",
	}))
	
	return &App{echo: e}
}

func (a *App) homePage(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

func (a *App) analyzeText(c echo.Context) error {
	text := c.FormValue("text")

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

	return c.Render(http.StatusOK, "result", data)
}

func (a *App) reverseText(c echo.Context) error {
	text := c.FormValue("text")

	runes := []rune(text)
	n := len(runes)
	for i := range n / 2 {
		runes[i], runes[n-1-i] = runes[n-1-i], runes[i]
	}

	data := map[string]interface{}{
		"Title":  "Reversed Text:",
		"Result": string(runes),
	}

	return c.Render(http.StatusOK, "result", data)
}

func (a *App) titleCase(c echo.Context) error {
	text := c.FormValue("text")

	data := map[string]interface{}{
		"Title":  "Title Case:",
		"Result": strings.Title(text),
	}

	return c.Render(http.StatusOK, "result", data)
}

func (a *App) countVowels(c echo.Context) error {
	text := c.FormValue("text")

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

	return c.Render(http.StatusOK, "result", data)
}

func (a *App) randomNumber(c echo.Context) error {
	minStr := c.FormValue("min")
	maxStr := c.FormValue("max")

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

	return c.Render(http.StatusOK, "result", data)
}

func (a *App) randomString(c echo.Context) error {
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

	return c.Render(http.StatusOK, "result", data)
}

func (a *App) generateUUID(c echo.Context) error {
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

	return c.Render(http.StatusOK, "result", data)
}

func (a *App) sha256Hash(c echo.Context) error {
	text := c.FormValue("text")

	hash := sha256.Sum256([]byte(text))

	data := map[string]interface{}{
		"Title":  "SHA256 Hash:",
		"Result": fmt.Sprintf("%x", hash),
	}

	return c.Render(http.StatusOK, "result", data)
}

func (a *App) base64Encode(c echo.Context) error {
	text := c.FormValue("text")

	encoded := base64.StdEncoding.EncodeToString([]byte(text))

	data := map[string]interface{}{
		"Title":  "Base64 Encoded:",
		"Result": encoded,
	}

	return c.Render(http.StatusOK, "result", data)
}

func (a *App) base64Decode(c echo.Context) error {
	text := c.FormValue("text")

	decoded, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		data := map[string]interface{}{
			"Title":  "Error:",
			"Result": "Invalid base64 string",
		}
		return c.Render(http.StatusOK, "result", data)
	}

	data := map[string]interface{}{
		"Title":  "Base64 Decoded:",
		"Result": string(decoded),
	}

	return c.Render(http.StatusOK, "result", data)
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

func (a *App) randomJoke(c echo.Context) error {
	// Generate random index
	b := make([]byte, 1)
	rand.Read(b)
	index := int(b[0]) % len(jokes)

	data := map[string]interface{}{
		"Title":  "Random Joke:",
		"Result": jokes[index],
	}

	return c.Render(http.StatusOK, "result", data)
}

func (a *App) randomQuote(c echo.Context) error {
	// Generate random index
	b := make([]byte, 1)
	rand.Read(b)
	index := int(b[0]) % len(quotes)

	data := map[string]interface{}{
		"Title":  "Inspirational Quote:",
		"Result": quotes[index],
	}

	return c.Render(http.StatusOK, "result", data)
}

func (a *App) randomFact(c echo.Context) error {
	// Generate random index
	b := make([]byte, 1)
	rand.Read(b)
	index := int(b[0]) % len(facts)

	data := map[string]interface{}{
		"Title":  "Random Tech Fact:",
		"Result": facts[index],
	}

	return c.Render(http.StatusOK, "result", data)
}

// Time and Date APIs
func (a *App) currentTime(c echo.Context) error {
	timezone := c.FormValue("timezone")
	if timezone == "" {
		timezone = "UTC"
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		data := map[string]interface{}{
			"Title":  "Error:",
			"Result": "Invalid timezone: " + timezone,
		}
		return c.Render(http.StatusOK, "result", data)
	}

	now := time.Now().In(loc)
	formatted := now.Format("Monday, January 2, 2006 at 3:04:05 PM MST")

	data := map[string]interface{}{
		"Title":  "Current Time (" + timezone + "):",
		"Result": formatted,
	}

	return c.Render(http.StatusOK, "result", data)
}

func (a *App) unixTimestamp(c echo.Context) error {
	timezone := c.FormValue("timezone")
	if timezone == "" {
		timezone = "UTC"
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		data := map[string]interface{}{
			"Title":  "Error:",
			"Result": "Invalid timezone: " + timezone,
		}
		return c.Render(http.StatusOK, "result", data)
	}

	now := time.Now().In(loc)
	timestamp := now.Unix()

	data := map[string]interface{}{
		"Title":  "Unix Timestamp:",
		"Result": fmt.Sprintf("%d", timestamp),
	}

	return c.Render(http.StatusOK, "result", data)
}

func (a *App) dateFormats(c echo.Context) error {
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

	return c.Render(http.StatusOK, "result", data)
}

func (a *App) setupRoutes() {
	// Static files using Echo's StaticFS
	a.echo.StaticFS("/static", echo.MustSubFS(distFS, "dist"))

	// Routes
	a.echo.GET("/", a.homePage)
	a.echo.POST("/analyze", a.analyzeText)
	a.echo.POST("/reverse", a.reverseText)
	a.echo.POST("/titlecase", a.titleCase)
	a.echo.POST("/count-vowels", a.countVowels)
	a.echo.POST("/random-number", a.randomNumber)
	a.echo.POST("/random-string", a.randomString)
	a.echo.POST("/uuid", a.generateUUID)
	a.echo.POST("/sha256", a.sha256Hash)
	a.echo.POST("/base64-encode", a.base64Encode)
	a.echo.POST("/base64-decode", a.base64Decode)

	// Fun APIs
	a.echo.GET("/joke", a.randomJoke)
	a.echo.GET("/quote", a.randomQuote)
	a.echo.GET("/fact", a.randomFact)

	// Time and Date APIs
	a.echo.POST("/current-time", a.currentTime)
	a.echo.POST("/timestamp", a.unixTimestamp)
	a.echo.GET("/date-formats", a.dateFormats)
}

func main() {
	app := NewApp()
	app.setupRoutes()

	// Start server
	app.echo.Logger.Info("API Playground starting on port :8080")
	if err := app.echo.Start(":8080"); err != nil && err != http.ErrServerClosed {
		app.echo.Logger.Fatal("Failed to start server:", err)
		log.Fatal(err)
	}
}