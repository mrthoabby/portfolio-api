package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mrthoabby/portfolio-api/internal/common"
)

func Test_getClientIP_XForwardedFor(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Forwarded-For", "10.0.0.1")
	req.RemoteAddr = "192.168.1.1:1234"

	ip := getClientIP(req)
	assert.Equal(t, "10.0.0.1", ip)
}

func Test_getClientIP_XRealIP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Real-IP", "10.0.0.2")
	req.RemoteAddr = "192.168.1.1:1234"

	ip := getClientIP(req)
	assert.Equal(t, "10.0.0.2", ip)
}

func Test_getClientIP_RemoteAddr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.1:1234"

	ip := getClientIP(req)
	assert.Equal(t, "192.168.1.1:1234", ip)
}

func TestClientIP_Middleware(t *testing.T) {
	handler := ClientIP(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify IP is in context using common package
		ip := common.ClientIPFromContext(r.Context())
		assert.Equal(t, "192.168.1.1:1234", ip)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

