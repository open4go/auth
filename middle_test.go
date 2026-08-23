package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestIsWebSocketRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("upgrade headers", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/v1/data/event", nil)
		req.Header.Set("Connection", "Upgrade")
		req.Header.Set("Upgrade", "websocket")
		c.Request = req
		if !isWebSocketRequest(c) {
			t.Fatal("expected websocket handshake to be detected")
		}
	})

	t.Run("ws path without upgrade", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/v1/data/ws", nil)
		c.Request = req
		if !isWebSocketRequest(c) {
			t.Fatal("expected /ws path to be treated as websocket endpoint")
		}
	})

	t.Run("normal rest", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/v1/system/order", nil)
		c.Request = req
		if isWebSocketRequest(c) {
			t.Fatal("did not expect REST path to be websocket")
		}
	})
}
