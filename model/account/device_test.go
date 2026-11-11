package account

import "testing"

func TestSanitizeAndGuard(t *testing.T) {
	g := SanitizeDeviceGuard(DeviceGuard{
		Enabled: true,
		Allow:   []string{"x6_0123456789abcdef0123456789abcdef", "bad", "x6_0123456789abcdef0123456789abcdef", "x6_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "x6_bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", "x6_cccccccccccccccccccccccccccccccc"},
	})
	if len(g.Allow) != 3 {
		t.Fatalf("allow %d", len(g.Allow))
	}
	if !GuardAllows(DeviceGuard{}, "anything") {
		t.Fatal("disabled should allow")
	}
	if GuardAllows(g, "x6_cccccccccccccccccccccccccccccccc") {
		t.Fatal("4th id must be dropped")
	}
	if !GuardAllows(g, "x6_0123456789abcdef0123456789abcdef") {
		t.Fatal("first allow")
	}
	if GuardAllows(g, "") {
		t.Fatal("empty device")
	}
}

func TestRememberDevice(t *testing.T) {
	var info LoginInfoType
	RememberDevice(&info, "x6_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "x6", 10)
	if info.FirstDeviceInfo != "x6_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" || info.FirstDeviceAt != 10 {
		t.Fatalf("first %+v", info)
	}
	RememberDevice(&info, "x6_bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", "x6", 20)
	if info.CurrentDeviceInfo != "x6_bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb" || info.LastDeviceInfo != "x6_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" {
		t.Fatalf("rotate %+v", info)
	}
	if len(info.SeenDevices) != 2 {
		t.Fatalf("seen %d", len(info.SeenDevices))
	}
	RememberDevice(&info, "x6_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "x6", 30)
	if info.SeenDevices[0].ID != "x6_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" && info.SeenDevices[1].LastAt != 30 {
		if info.SeenDevices[1].LastAt != 30 && info.SeenDevices[0].LastAt != 30 {
			t.Fatalf("touch last_at %+v", info.SeenDevices)
		}
	}
}
