package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/open4go/auth/model/account"
	"github.com/open4go/log"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	deviceGuardRedisPrefix = "pos:devguard:"
	deviceGuardTTL         = 10 * time.Minute
	deviceNotAllowed       = "device_not_allowed"
)

type deviceGuardCache struct {
	Enabled bool     `json:"enabled"`
	Allow   []string `json:"allow"`
}

func deviceGuardKey(userID string) string {
	return deviceGuardRedisPrefix + strings.TrimSpace(userID)
}

func WriteDeviceGuardCache(ctx context.Context, userID string, g account.DeviceGuard) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return
	}
	body, err := json.Marshal(deviceGuardCache{Enabled: g.Enabled, Allow: g.Allow})
	if err != nil {
		return
	}
	_ = GetRedisAuthHandler(ctx).Set(ctx, deviceGuardKey(userID), body, deviceGuardTTL).Err()
}

func loadDeviceGuard(ctx context.Context, userID string) (account.DeviceGuard, error) {
	var g account.DeviceGuard
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return g, nil
	}
	raw, err := GetRedisAuthHandler(ctx).Get(ctx, deviceGuardKey(userID)).Bytes()
	if err == nil && len(raw) > 0 {
		var cached deviceGuardCache
		if json.Unmarshal(raw, &cached) == nil {
			return account.DeviceGuard{Enabled: cached.Enabled, Allow: cached.Allow}, nil
		}
	}
	m := &account.Model{}
	h := m.Init(ctx, GetDBHandler(ctx), m.CollectionName())
	if err := h.GetBy(m, bson.M{"user_id": userID}); err != nil {
		return g, err
	}
	WriteDeviceGuardCache(ctx, userID, m.DeviceGuard)
	return m.DeviceGuard, nil
}

func checkX6DeviceGuard(c *gin.Context, userID string) bool {
	if !strings.EqualFold(strings.TrimSpace(c.GetHeader("X-Client-Id")), "x6") {
		return true
	}
	if c.GetBool("OPEN4GO_VERIFIED_SUPER_ADMIN") {
		return true
	}
	g, err := loadDeviceGuard(c.Request.Context(), userID)
	if err != nil {
		log.Log(c.Request.Context()).WithError(err).Warning("device guard lookup failed")
		return true
	}
	deviceID := strings.TrimSpace(c.GetHeader("X-Device-ID"))
	if account.GuardAllows(g, deviceID) {
		return true
	}
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"error":   deviceNotAllowed,
		"message": "本机未加入点餐设备白名单，请联系管理员",
	})
	return false
}
