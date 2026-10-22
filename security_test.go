package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/open4go/model"
	"github.com/open4go/req5rsp/cst"
	"github.com/spf13/viper"
)

func setCfg(t *testing.T, key string, val interface{}) {
	t.Helper()
	prev := viper.Get(key)
	viper.Set(key, val)
	t.Cleanup(func() { viper.Set(key, prev) })
}

func TestTranslateHTTPMethodToPermissionMapping(t *testing.T) {
	cases := []struct {
		method string
		single bool
		want   cst.Permission
	}{
		{"GET", false, cst.List},
		{"HEAD", false, cst.List},
		{"POST", false, cst.Create},
		{"PUT", false, cst.Update},
		{"PATCH", false, cst.Update},
		{"DELETE", false, cst.Delete},
		{"GET", true, cst.Read},
		{"HEAD", true, cst.Read},
		{"PUT", true, cst.Update},
		{"PATCH", true, cst.Update},
		{"POST", true, cst.Update},
		{"DELETE", true, cst.Delete},
		{"OPTIONS", false, 0},
		{"OPTIONS", true, 0},
		{"TRACE", false, 0},
		{"CONNECT", true, 0},
		{"PROPFIND", false, 0},
	}
	for _, tc := range cases {
		if got := translateHTTPMethodToPermission(tc.method, tc.single); got != tc.want {
			t.Errorf("%s single=%v = %v, want %v", tc.method, tc.single, got, tc.want)
		}
	}
}

func TestPermAllowedDefaultDeny(t *testing.T) {
	all := int(cst.Create | cst.Delete | cst.Update | cst.Read | cst.List | cst.Show | cst.Disable)
	if permAllowed(all, 0, false) {
		t.Fatal("unmapped permission (0) must be denied even with all bits")
	}
	if permAllowed(0, 0, false) {
		t.Fatal("unmapped permission (0) must be denied")
	}
	if !permAllowed(int(cst.Update), cst.Update, false) {
		t.Fatal("matching bit should be allowed")
	}
	if permAllowed(int(cst.Read), cst.Update, false) {
		t.Fatal("missing bit must be denied")
	}
}

func TestVerifyDeniesUnmappedMethodWithoutTouchingRoles(t *testing.T) {
	r := RoleManager{"system:auth", nil}
	for _, m := range []string{"OPTIONS", "TRACE", "PROPFIND"} {
		ok, err := r.Verify(context.Background(), "/v1/x", "acc", m, false)
		if ok || err != nil {
			t.Fatalf("%s: ok=%v err=%v, want denied", m, ok, err)
		}
	}
}

func TestResolveCheckMethodWSOnlyForGET(t *testing.T) {
	ctx := withWSFlag(context.Background(), true)
	for _, m := range []string{"POST", "PUT", "DELETE", "PATCH"} {
		cm, isWS := resolveCheckMethod(ctx, "/v1/order/ws", m)
		if isWS || cm != m {
			t.Fatalf("%s on ws path: method=%s isWS=%v, want original method and not WS", m, cm, isWS)
		}
	}
	if _, isWS := resolveCheckMethod(context.Background(), "/v1/order/ws", "GET"); !isWS {
		t.Fatal("GET on ws path should use WS rule")
	}
	if _, isWS := resolveCheckMethod(ctx, "/v1/order", "GET"); !isWS {
		t.Fatal("GET with ws flag should use WS rule")
	}
	if _, isWS := resolveCheckMethod(context.Background(), "/v1/order", "GET"); isWS {
		t.Fatal("plain GET must not use WS rule")
	}
}

func TestIsWebSocketRequestRequiresGET(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, m := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		req := httptest.NewRequest(m, "/v1/data/ws", nil)
		req.Header.Set("Connection", "Upgrade")
		req.Header.Set("Upgrade", "websocket")
		c.Request = req
		if isWebSocketRequest(c) {
			t.Fatalf("%s with upgrade headers must not be treated as websocket", m)
		}
	}
}

func newAccessEngine(pre gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(pre, AccessMiddleware())
	r.GET("/v1/x", func(c *gin.Context) { c.Status(http.StatusOK) })
	return r
}

func TestAccessMiddlewareIgnoresUserIDHeaderByDefault(t *testing.T) {
	setCfg(t, cfgIdentitySource, "")
	r := newAccessEngine(func(c *gin.Context) { c.Next() })
	req := httptest.NewRequest(http.MethodGet, "/v1/x", nil)
	req.Header.Set("UserID", "someone-else")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("code=%d, want 401 when only a client header is present", w.Code)
	}
}

func TestResolveUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mk := func() *gin.Context {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		req := httptest.NewRequest(http.MethodGet, "/v1/x", nil)
		req.Header.Set("UserID", "from-header")
		//nolint:staticcheck
		ctx := context.WithValue(req.Context(), model.OperatorKey, "from-operator")
		c.Request = req.WithContext(ctx)
		c.Set(VerifiedUserIDKey, "from-token")
		return c
	}

	setCfg(t, cfgIdentitySource, "")
	if id, src := resolveUserID(mk()); id != "from-token" || src != IdentityVerified {
		t.Fatalf("default: %q/%q", id, src)
	}
	setCfg(t, cfgIdentitySource, "bogus")
	if id, _ := resolveUserID(mk()); id != "from-token" {
		t.Fatalf("unknown value must fall back to verified, got %q", id)
	}
	setCfg(t, cfgIdentitySource, IdentityContext)
	if id, _ := resolveUserID(mk()); id != "from-operator" {
		t.Fatalf("context: %q", id)
	}
	setCfg(t, cfgIdentitySource, IdentityHeader)
	if id, _ := resolveUserID(mk()); id != "from-header" {
		t.Fatalf("header: %q", id)
	}

	setCfg(t, cfgIdentitySource, "")
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	req := httptest.NewRequest(http.MethodGet, "/v1/x", nil)
	//nolint:staticcheck
	c.Request = req.WithContext(context.WithValue(req.Context(), VerifiedUserIDKey, "ctx-token"))
	if id, _ := resolveUserID(c); id != "ctx-token" {
		t.Fatalf("verified via request context: %q", id)
	}
}
