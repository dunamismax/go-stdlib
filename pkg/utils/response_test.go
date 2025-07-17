package utils

import (
	"net/http/httptest"
	"testing"
)

func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"foo": "bar"}
	JSON(w, 200, data)

	if w.Code != 200 {
		t.Errorf("JSON() status = %v, want %v", w.Code, 200)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("JSON() Content-Type = %v, want %v", w.Header().Get("Content-Type"), "application/json")
	}
}

func TestError(t *testing.T) {
	w := httptest.NewRecorder()
	Error(w, 400, "bad request")

	if w.Code != 400 {
		t.Errorf("Error() status = %v, want %v", w.Code, 400)
	}
}

func TestSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	Success(w, "ok")

	if w.Code != 200 {
		t.Errorf("Success() status = %v, want %v", w.Code, 200)
	}
}
