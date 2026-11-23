package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	r2redis "github.com/open4go/db/redis"
)

func decodeCode(t *testing.T, w *httptest.ResponseRecorder) string {
	t.Helper()
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("body is not json: %q", w.Body.String())
	}
	if body["message"] == "" {
		t.Fatalf("body has no message: %q", w.Body.String())
	}
	return body["code"]
}

func TestAccessMiddlewareNoIdentityIsSessionExpired(t *testing.T) {
	setCfg(t, cfgIdentitySource, "")
	r := newAccessEngine(func(c *gin.Context) { c.Next() })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/v1/x", nil))
	if w.Code != http.StatusUnauthorized || decodeCode(t, w) != CodeSessionExpired {
		t.Fatalf("code=%d body=%s, want 401 session_expired", w.Code, w.Body.String())
	}
}

// 需要本地 redis：AUTH_TEST_REDIS=127.0.0.1:6379 go test ./...
func TestAccessMiddlewareDenyVsExpired(t *testing.T) {
	addr := os.Getenv("AUTH_TEST_REDIS")
	if addr == "" {
		t.Skip("AUTH_TEST_REDIS not set")
	}
	ctx := context.Background()
	r2redis.NewDataBasePool()
	if _, err := r2redis.DBPool.GetClient(ctx, addr, "auth", 15); err != nil {
		t.Fatal(err)
	}
	rdb := GetRedisAuthHandler(ctx)
	rdb.FlushDB(ctx)

	setCfg(t, cfgIdentitySource, "")
	as := func(uid string) gin.HandlerFunc {
		return func(c *gin.Context) { c.Set(VerifiedUserIDKey, uid); c.Next() }
	}

	// 账号没有角色 → 403 permission_denied（不能当登录失效，否则无角色账号会被反复踢回登录页）
	w := httptest.NewRecorder()
	newAccessEngine(as("u-none")).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/v1/x", nil))
	if w.Code != http.StatusForbidden || decodeCode(t, w) != CodePermissionDenied {
		t.Fatalf("no roles: code=%d body=%s", w.Code, w.Body.String())
	}

	// 有角色但该角色没有这个接口 → 403 permission_denied
	rdb.SAdd(ctx, "system:auth:account:to:role:u-1", "role-x")
	w = httptest.NewRecorder()
	newAccessEngine(as("u-1")).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/v1/x", nil))
	if w.Code != http.StatusForbidden || decodeCode(t, w) != CodePermissionDenied {
		t.Fatalf("no permission: code=%d body=%s", w.Code, w.Body.String())
	}

	// redis 读失败 → 503 permission_unavailable
	_ = rdb.Close()
	w = httptest.NewRecorder()
	newAccessEngine(as("u-1")).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/v1/x", nil))
	if w.Code != http.StatusServiceUnavailable || decodeCode(t, w) != CodePermissionUnavailable {
		t.Fatalf("redis down: code=%d body=%s", w.Code, w.Body.String())
	}
}
