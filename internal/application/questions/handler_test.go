package questions

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Create_InvalidProfileID(t *testing.T) {
	handler := &Handler{
		service: nil,
	}

	body := `{"message": "Test question message"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/profiles/invalid-uuid/questions", bytes.NewBufferString(body))
	req.SetPathValue("id", "invalid-uuid")
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	errorObj, ok := response["error"].(map[string]interface{})
	if !ok {
		t.Fatal("expected error object in response")
	}

	if errorObj["code"] != "BAD_REQUEST" {
		t.Errorf("expected error code BAD_REQUEST, got %s", errorObj["code"])
	}
}

func TestHandler_Create_MissingProfileID(t *testing.T) {
	handler := &Handler{
		service: nil,
	}

	body := `{"message": "Test question message"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/profiles//questions", bytes.NewBufferString(body))
	req.SetPathValue("id", "")
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandler_Create_InvalidJSON(t *testing.T) {
	handler := &Handler{
		service: nil,
	}

	body := `{invalid json}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/profiles/550e8400-e29b-41d4-a716-446655440000/questions", bytes.NewBufferString(body))
	req.SetPathValue("id", "550e8400-e29b-41d4-a716-446655440000")
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
