package contacts

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler_Create_MissingID(t *testing.T) {
	service := &Service{}
	handler := NewHandler(service)

	req := httptest.NewRequest("POST", "/api/v1/profiles//contacts", nil)
	req = req.WithContext(context.Background())
	w := httptest.NewRecorder()

	handler.Create(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Create_InvalidJSON(t *testing.T) {
	service := &Service{}
	handler := NewHandler(service)

	req := httptest.NewRequest("POST", "/api/v1/profiles/test-id/contacts", strings.NewReader("invalid json"))
	req = req.WithContext(context.Background())
	req.SetPathValue("id", "test-id")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Create_InvalidRequest(t *testing.T) {
	service := &Service{}
	handler := NewHandler(service)

	contactReq := Request{
		Name:    "T",
		Email:   "invalid-email",
		Message: "Short",
	}

	body, _ := json.Marshal(contactReq)
	req := httptest.NewRequest("POST", "/api/v1/profiles/test-id/contacts", strings.NewReader(string(body)))
	req = req.WithContext(context.Background())
	req.SetPathValue("id", "test-id")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNewHandler(t *testing.T) {
	service := &Service{}
	handler := NewHandler(service)
	assert.NotNil(t, handler)
	assert.Equal(t, service, handler.service)
	assert.NotNil(t, handler.validator)
}

