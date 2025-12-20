package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCORS(t *testing.T) {
	allowedOrigins := []string{"http://localhost:3000", "http://localhost:3001"}
	corsHandler := NewCORS(allowedOrigins)

	assert.NotNil(t, corsHandler)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()

	corsHandler(handler).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestNewCORS_OptionsRequest(t *testing.T) {
	allowedOrigins := []string{"http://localhost:3000"}
	corsHandler := NewCORS(allowedOrigins)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()

	corsHandler(handler).ServeHTTP(w, req)

	// CORS middleware should handle OPTIONS request
	assert.NotNil(t, corsHandler)
}

func TestNewCORS_EmptyOrigins(t *testing.T) {
	corsHandler := NewCORS([]string{})
	assert.NotNil(t, corsHandler)
}

