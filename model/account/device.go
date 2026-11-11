package account

import (
	"regexp"
	"strings"
	"time"
)

const (
	MaxDeviceAllow = 3
	maxSeenDevices = 8
)

var deviceIDRe = regexp.MustCompile(`(?i)^x6_[0-9a-f]{16,64}$`)

func ValidDeviceID(id string) bool {
	return deviceIDRe.MatchString(strings.TrimSpace(id))
}

func SanitizeDeviceGuard(g DeviceGuard) DeviceGuard {
	seen := make([]string, 0, MaxDeviceAllow)
	uniq := map[string]struct{}{}
	for _, id := range g.Allow {
		id = strings.TrimSpace(id)
		if !ValidDeviceID(id) {
			continue
		}
		key := strings.ToLower(id)
		if _, ok := uniq[key]; ok {
			continue
		}
		uniq[key] = struct{}{}
		seen = append(seen, id)
		if len(seen) >= MaxDeviceAllow {
			break
		}
	}
	g.Allow = seen
	return g
}

func GuardAllows(g DeviceGuard, deviceID string) bool {
	if !g.Enabled {
		return true
	}
	deviceID = strings.TrimSpace(deviceID)
	if deviceID == "" {
		return false
	}
	for _, id := range g.Allow {
		if strings.EqualFold(id, deviceID) {
			return true
		}
	}
	return false
}

func RememberDevice(info *LoginInfoType, deviceID, client string, now int64) {
	if info == nil {
		return
	}
	deviceID = strings.TrimSpace(deviceID)
	if deviceID == "" {
		return
	}
	if now <= 0 {
		now = time.Now().Unix()
	}
	if info.FirstDeviceInfo == "" {
		info.FirstDeviceInfo = deviceID
		info.FirstDeviceAt = now
	}
	if info.CurrentDeviceInfo != deviceID {
		info.LastDeviceInfo = info.CurrentDeviceInfo
		info.CurrentDeviceInfo = deviceID
	} else {
		info.CurrentDeviceInfo = deviceID
	}
	found := false
	for i := range info.SeenDevices {
		if strings.EqualFold(info.SeenDevices[i].ID, deviceID) {
			info.SeenDevices[i].LastAt = now
			if client != "" {
				info.SeenDevices[i].Client = client
			}
			found = true
			break
		}
	}
	if !found {
		info.SeenDevices = append([]SeenDevice{{
			ID:      deviceID,
			FirstAt: now,
			LastAt:  now,
			Client:  client,
		}}, info.SeenDevices...)
		if len(info.SeenDevices) > maxSeenDevices {
			info.SeenDevices = info.SeenDevices[:maxSeenDevices]
		}
	}
}
