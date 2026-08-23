package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/open4go/req5rsp/cst"
)

type authCtxKey int

const isWSCtxKey authCtxKey = 1

func withWSFlag(ctx context.Context, isWS bool) context.Context {
	if !isWS {
		return ctx
	}
	return context.WithValue(ctx, isWSCtxKey, true)
}

func wsFlagFromCtx(ctx context.Context) bool {
	v, _ := ctx.Value(isWSCtxKey).(bool)
	return v
}

// isWSEndpoint 判断路由是否为 websocket 接口（路径中含独立的 ws / websocket 段）。
func isWSEndpoint(path string) bool {
	if path == "" {
		return false
	}
	p := strings.ToLower(strings.TrimSuffix(path, "/"))
	for _, seg := range strings.Split(p, "/") {
		if seg == "ws" || seg == "websocket" {
			return true
		}
	}
	return false
}

// hasGETPermission WS 接口只需 GET：具备 List（集合 GET）或 Read（单资源 GET）任一即可。
func hasGETPermission(attr int) bool {
	return attr&int(cst.List) != 0 || attr&int(cst.Read) != 0
}

func formatPermission(attr int) string {
	if attr == 0 {
		return "none(0)"
	}
	names := make([]string, 0, 7)
	for _, bit := range []cst.Permission{
		cst.Create, cst.Delete, cst.Update, cst.Read, cst.List, cst.Show, cst.Disable,
	} {
		if attr&int(bit) != 0 {
			names = append(names, cst.PermissionNames[bit])
		}
	}
	if len(names) == 0 {
		return fmt.Sprintf("unknown(%d)", attr)
	}
	return fmt.Sprintf("%s(%d)", strings.Join(names, "|"), attr)
}
