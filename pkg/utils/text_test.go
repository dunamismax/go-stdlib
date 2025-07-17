package utils

import (
	"testing"
)

func TestReverseString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty string", "", ""},
		{"a single character", "a", "a"},
		{"a word", "hello", "olleh"},
		{"a sentence", "hello world", "dlrow olleh"},
		{"a palindrome", "madam", "madam"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReverseString(tt.input); got != tt.want {
				t.Errorf("ReverseString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty string", "", true},
		{"a single character", "a", true},
		{"a palindrome", "madam", true},
		{"a non-palindrome", "hello", false},
		{"a sentence palindrome", "A man, a plan, a canal: Panama", true},
		{"a sentence non-palindrome", "hello world", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPalindrome(tt.input); got != tt.want {
				t.Errorf("IsPalindrome() = %v, want %v", got, tt.want)
			}
		})
	}
}
