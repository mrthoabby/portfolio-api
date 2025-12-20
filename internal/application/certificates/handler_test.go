package certificates

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler_GetByProfileID_MissingID(t *testing.T) {
	service := &Service{}
	handler := NewHandler(service)

	req := httptest.NewRequest("GET", "/api/v1/profiles//certificates", nil)
	req = req.WithContext(context.Background())
	w := httptest.NewRecorder()

	handler.GetByProfileID(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNewHandler(t *testing.T) {
	service := &Service{}
	handler := NewHandler(service)
	assert.NotNil(t, handler)
	assert.Equal(t, service, handler.service)
}

