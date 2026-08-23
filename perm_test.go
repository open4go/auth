package auth

import (
	"testing"

	"github.com/open4go/req5rsp/cst"
)

func TestIsWSEndpoint(t *testing.T) {
	cases := []struct {
		path string
		want bool
	}{
		{"/v1/system/data/ws", true},
		{"/ws", true},
		{"/ws/event", true},
		{"/v1/foo/websocket", true},
		{"/v1/system/data/ws/", true},
		{"/v1/system/order", false},
		{"/workspace", false},
		{"/v1/ws-admin/order", false},
		{"", false},
	}
	for _, tc := range cases {
		if got := isWSEndpoint(tc.path); got != tc.want {
			t.Errorf("isWSEndpoint(%q) = %v, want %v", tc.path, got, tc.want)
		}
	}
}

func TestHasGETPermission(t *testing.T) {
	if hasGETPermission(0) {
		t.Fatal("none should not grant GET")
	}
	if !hasGETPermission(int(cst.List)) {
		t.Fatal("List should grant GET")
	}
	if !hasGETPermission(int(cst.Read)) {
		t.Fatal("Read should grant GET")
	}
	if !hasGETPermission(int(cst.List | cst.Read)) {
		t.Fatal("List|Read should grant GET")
	}
	if hasGETPermission(int(cst.Create | cst.Update | cst.Delete)) {
		t.Fatal("write bits should not grant GET")
	}
}

func TestTranslateHTTPMethodToPermission(t *testing.T) {
	if p := translateHTTPMethodToPermission("GET", false); p != cst.List {
		t.Errorf("collection GET = %v, want List", p)
	}
	if p := translateHTTPMethodToPermission("GET", true); p != cst.Read {
		t.Errorf("single GET = %v, want Read", p)
	}
	if p := translateHTTPMethodToPermission("POST", false); p != cst.Create {
		t.Errorf("collection POST = %v, want Create", p)
	}
	if p := translateHTTPMethodToPermission("CONNECT", false); p != 0 {
		t.Errorf("unknown method = %v, want 0", p)
	}
}

func TestFormatPermission(t *testing.T) {
	if got := formatPermission(0); got != "none(0)" {
		t.Errorf("got %q", got)
	}
	if got := formatPermission(int(cst.List)); got != "List(16)" {
		t.Errorf("got %q", got)
	}
	if got := formatPermission(int(cst.Read | cst.List)); got != "Read|List(24)" {
		t.Errorf("got %q", got)
	}
}
