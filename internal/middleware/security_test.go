package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityHeaders(t *testing.T) {
	handler := SecurityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "nosniff", rr.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", rr.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", rr.Header().Get("X-XSS-Protection"))
	assert.Equal(t, "strict-origin-when-cross-origin", rr.Header().Get("Referrer-Policy"))
	assert.NotEmpty(t, rr.Header().Get("Content-Security-Policy"))
	assert.Equal(t, "no-store, no-cache, must-revalidate, proxy-revalidate", rr.Header().Get("Cache-Control"))
}

func TestMaxBodySize_LimitsPostRequests(t *testing.T) {
	maxSize := int64(100)
	handler := MaxBodySize(maxSize)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to read body
		buf := make([]byte, 200)
		_, err := r.Body.Read(buf)
		if err != nil && err.Error() == "http: request body too large" {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))

	// Small body - should pass
	smallBody := strings.NewReader("small body")
	req := httptest.NewRequest(http.MethodPost, "/test", smallBody)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Large body - should fail
	largeBody := strings.NewReader(strings.Repeat("x", 200))
	req = httptest.NewRequest(http.MethodPost, "/test", largeBody)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusRequestEntityTooLarge, rr.Code)
}

func TestMaxBodySize_DoesNotLimitGetRequests(t *testing.T) {
	maxSize := int64(10)
	handler := MaxBodySize(maxSize)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRecoverPanic(t *testing.T) {
	handler := RecoverPanic(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	// Should not panic
	assert.NotPanics(t, func() {
		handler.ServeHTTP(rr, req)
	})

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

