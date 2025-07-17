package utils

import (
	"testing"
)

func TestRandomInt(t *testing.T) {
	min := 1
	max := 100
	for i := 0; i < 100; i++ {
		got := RandomInt(min, max)
		if got < min || got > max {
			t.Errorf("RandomInt() = %v, want between %v and %v", got, min, max)
		}
	}
}

func TestRandomString(t *testing.T) {
	length := 10
	got := RandomString(length)
	if len(got) != length {
		t.Errorf("RandomString() length = %v, want %v", len(got), length)
	}
}
