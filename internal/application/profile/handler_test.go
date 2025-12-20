package profile

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler_GetByID_MissingID(t *testing.T) {
	service := &Service{}
	handler := NewHandler(service)

	req := httptest.NewRequest("GET", "/api/v1/profiles/", nil)
	req = req.WithContext(context.Background())
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNewHandler(t *testing.T) {
	service := &Service{}
	handler := NewHandler(service)
	assert.NotNil(t, handler)
	assert.Equal(t, service, handler.service)
}

