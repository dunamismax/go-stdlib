package utils

import (
	"regexp"
	"sort"
	"strings"
)

type TextAnalysis struct {
	WordCount        int            `json:"word_count"`
	CharCount        int            `json:"char_count"`
	CharCountNoSpace int            `json:"char_count_no_space"`
	SentenceCount    int            `json:"sentence_count"`
	ParagraphCount   int            `json:"paragraph_count"`
	WordFrequency    map[string]int `json:"word_frequency"`
	ReadingTime      int            `json:"reading_time_minutes"`
}

func AnalyzeText(text string) TextAnalysis {
	analysis := TextAnalysis{
		WordFrequency: make(map[string]int),
	}

	analysis.CharCount = len(text)
	analysis.CharCountNoSpace = len(strings.ReplaceAll(text, " ", ""))

	words := strings.Fields(strings.ToLower(text))
	analysis.WordCount = len(words)

	for _, word := range words {
		cleaned := cleanWord(word)
		if cleaned != "" {
			analysis.WordFrequency[cleaned]++
		}
	}

	analysis.SentenceCount = countSentences(text)
	analysis.ParagraphCount = len(strings.Split(strings.TrimSpace(text), "\n\n"))
	analysis.ReadingTime = analysis.WordCount / 200 // ~200 words per minute
	if analysis.ReadingTime < 1 {
		analysis.ReadingTime = 1
	}

	return analysis
}

func cleanWord(word string) string {
	reg := regexp.MustCompile(`[^\p{L}\p{N}]+`)
	return reg.ReplaceAllString(word, "")
}

func countSentences(text string) int {
	reg := regexp.MustCompile(`[.!?]+`)
	return len(reg.FindAllString(text, -1))
}

func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func CountVowels(s string) int {
	vowels := "aeiouAEIOU"
	count := 0
	for _, char := range s {
		if strings.ContainsRune(vowels, char) {
			count++
		}
	}
	return count
}

func IsPalindrome(s string) bool {
	cleaned := strings.ToLower(regexp.MustCompile(`[^a-zA-Z0-9]`).ReplaceAllString(s, ""))
	return cleaned == ReverseString(cleaned)
}

func ToTitleCase(s string) string {
	return strings.Title(s)
}

func ExtractWords(text string) []string {
	reg := regexp.MustCompile(`\b\w+\b`)
	return reg.FindAllString(strings.ToLower(text), -1)
}

func SortWords(words []string) []string {
	sorted := make([]string, len(words))
	copy(sorted, words)
	sort.Strings(sorted)
	return sorted
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return e.Message
}

type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return ""
	}
	return ve[0].Message
}

func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

func ValidateUsername(username string) *ValidationError {
	if username == "" {
		return &ValidationError{Field: "username", Message: "Username is required"}
	}
	if len(username) < 3 {
		return &ValidationError{Field: "username", Message: "Username must be at least 3 characters"}
	}
	if len(username) > 20 {
		return &ValidationError{Field: "username", Message: "Username must be less than 20 characters"}
	}

	validUsernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !validUsernameRegex.MatchString(username) {
		return &ValidationError{Field: "username", Message: "Username can only contain letters, numbers, and underscores"}
	}

	return nil
}

func ValidateEmail(email string) *ValidationError {
	if email == "" {
		return &ValidationError{Field: "email", Message: "Email is required"}
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return &ValidationError{Field: "email", Message: "Invalid email format"}
	}

	if len(email) > 100 {
		return &ValidationError{Field: "email", Message: "Email must be less than 100 characters"}
	}

	return nil
}

func ValidatePassword(password string) *ValidationError {
	if password == "" {
		return &ValidationError{Field: "password", Message: "Password is required"}
	}
	if len(password) < 8 {
		return &ValidationError{Field: "password", Message: "Password must be at least 8 characters"}
	}
	if len(password) > 128 {
		return &ValidationError{Field: "password", Message: "Password must be less than 128 characters"}
	}

	return nil
}

func ValidateDisplayName(displayName string) *ValidationError {
	if displayName == "" {
		return &ValidationError{Field: "display_name", Message: "Display name is required"}
	}
	if len(displayName) < 1 {
		return &ValidationError{Field: "display_name", Message: "Display name is required"}
	}
	if len(displayName) > 50 {
		return &ValidationError{Field: "display_name", Message: "Display name must be less than 50 characters"}
	}

	return nil
}

func ValidatePostContent(content string) *ValidationError {
	if content == "" {
		return &ValidationError{Field: "content", Message: "Post content cannot be empty"}
	}
	if len(content) > 280 {
		return &ValidationError{Field: "content", Message: "Post content must be less than 280 characters"}
	}

	return nil
}

func SanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	return strings.ReplaceAll(trimmed, "\x00", "")
}

func ValidateUserRegistration(username, email, displayName, password string) ValidationErrors {
	var errors ValidationErrors

	if err := ValidateUsername(username); err != nil {
		errors = append(errors, *err)
	}

	if err := ValidateEmail(email); err != nil {
		errors = append(errors, *err)
	}

	if err := ValidateDisplayName(displayName); err != nil {
		errors = append(errors, *err)
	}

	if err := ValidatePassword(password); err != nil {
		errors = append(errors, *err)
	}

	return errors
}