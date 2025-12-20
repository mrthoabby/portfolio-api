package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestID_GeneratesNewID(t *testing.T) {
	var capturedID string

	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID = GetRequestID(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.NotEmpty(t, capturedID)
	assert.NotEmpty(t, rr.Header().Get("X-Request-ID"))
	assert.Equal(t, capturedID, rr.Header().Get("X-Request-ID"))
}

func TestRequestID_UsesExistingID(t *testing.T) {
	existingID := "existing-request-id-123"
	var capturedID string

	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID = GetRequestID(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Request-ID", existingID)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, existingID, capturedID)
	assert.Equal(t, existingID, rr.Header().Get("X-Request-ID"))
}

func TestGetRequestID_EmptyContext(t *testing.T) {
	ctx := context.Background()
	id := GetRequestID(ctx)
	assert.Empty(t, id)
}

