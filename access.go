package auth

import (
	"context"
	"encoding/json"
	"fmt"
	//roleManage "github.com/open4go/auth/model/role"
	"github.com/open4go/log"
)

const (
	defaultRoleMemberKey       = "access:role:members:"
	defaultPermissionKey       = "access:permissions:"
	defaultPath2Method2RoleKey = "access:path:method:role:"
)

type MyRole struct {
	Name          string
	PermissionsV4 []MyPermission `json:"permissions_v4"`
}

type MyPermission struct {
	Path      string   `json:"path"`
	Operation []string `json:"operation"`
}

type MyACL struct {
	Ctx context.Context
	// debug, release
	Mode string
	// 角色存储的key (集合 <role:members> + <accountId> roles )
	RoleMemberKey string
	// 角色
	Roles []MyRole
	//
	App           MyApp
	PermissionKey string
	//
	Methods           []string
	Path2Method4Roles string
	// 角色权限

}

func NewMyACL(ctx context.Context, mode string,
	myApp MyApp, roles []MyRole) MyACL {
	return MyACL{
		ctx,
		mode,
		defaultRoleMemberKey,
		roles,
		myApp,
		defaultPermissionKey,
		[]string{
			"read",
			"write",
			"update",
			"delete",
			"list",
		},
		defaultPath2Method2RoleKey,
	}
}

// IsDebug 默认开启
func (a MyACL) isDebug() bool {
	// 查看是否是公开接口
	if a.Mode == "debug" {
		return true
	}
	return false
}

// IsPublicPath 默认开启
func (a MyACL) isPublicPath(path string) bool {
	// 查看是否是公开接口
	// 查看是否是公开接口
	if a.App.getApiAttribute(path, "public") == "1" {
		return true
	}
	return false
}

// IsDisablePath 默认关闭
func (a MyACL) isDisablePath(path string) bool {
	// 查看是否被禁用接口
	if a.App.getApiAttribute(path, "disable") == "1" {
		return true
	}
	return false
}

// myRoles 获取角色
// 如果不存在则返回false
// 默认为false
func (a MyACL) myRoles() []string {
	members, err := RDB.SMembers(a.Ctx, a.RoleMemberKey).Result()
	if err != nil {
		return nil
	}
	for _, name := range members {
		a.Roles = append(a.Roles, MyRole{
			Name:          name,
			PermissionsV4: make([]MyPermission, 0),
		})
	}
	// 如果没有则返回空
	return members
}

// SetPermission 设置权限
// 如果不存在则返回false
// 默认为false
func (a MyACL) setPermission(path string, method string, role string) error {
	secondKey := fmt.Sprintf("%s:%s", path, method)
	roleArr := a.loadRole(path, method)
	roleArr = append(roleArr, method)
	// 存储
	newPayload, err := json.Marshal(roleArr)
	if err != nil {
		return err
	}

	err = RDB.HSet(a.Ctx, a.PermissionKey, secondKey, newPayload).Err()
	if err != nil {
		return err
	}
	return nil
}

func (a MyACL) loadRole(path, method string) []string {
	secondKey := fmt.Sprintf("%s:%s", path, method)
	roleArr := make([]string, 0)

	role, err := RDB.HGet(a.Ctx, a.Path2Method4Roles, secondKey).Result()
	if err != nil {
		return roleArr
	}

	err = json.Unmarshal([]byte(role), &roleArr)
	if err != nil {
		return roleArr
	}
	// 如果没有则返回空
	return roleArr
}

// SetAllPermission 设置权限
// 如果不存在则返回false
// 默认为false
//func (a MyACL) SetAllPermission(roles []*roleManage.Model) error {
//	// 所有角色
//	for _, role := range roles {
//		// 读取角色里绑定的基本信息
//		for _, app := range role.PermissionsV2 {
//			// 遍历每个path的属性
//			for _, method := range app.Operation {
//				// 设定权限给角色名称
//				err := a.setPermission(app.Path, method, role.Name)
//				if err != nil {
//					continue
//				}
//			}
//		}
//	}
//	return nil
//}

//func (a MyACL) Path2Operations() error {
//	for path, _ := range a.App.GetAllPath() {
//		for _, method := range a.Methods {
//			role := a.loadRole(path, method)
//			err := a.setPermission(path, method, role)
//			if err != nil {
//				continue
//			}
//		}
//	}
//	return nil
//}

// myPermission 获取
// 如果不存在则返回false
// 默认为false
func (a MyACL) hasPermission(path string, method string) bool {
	for _, role := range a.Roles {
		secondKey := fmt.Sprintf("%s:%s", path, method)
		// load roleHas
		rolePayload, err := RDB.HGet(a.Ctx, a.PermissionKey, secondKey).Result()
		if err != nil {
			continue
		}
		//
		roleArr := make([]string, 0)
		err = json.Unmarshal([]byte(rolePayload), &roleArr)
		if err != nil {
			continue
		}

		// 如果当前的角色已经包含则可直接返回
		// 否则继续尝试对比下一个角色
		if contains(roleArr, role.Name) {
			return true
		}
	}
	// 如果没有则返回空
	return false
}

// CanVisit 判断是否能够访问接口
func (a MyACL) CanVisit(path string, method string, accountId string) bool {

	// 查看是否已经启用接口
	if a.isDisablePath(path) {
		log.Log().WithField("path", path).Info("has set disable")
		return false
	}

	// 查看是否是公开接口
	if a.isPublicPath(path) {
		log.Log().WithField("path", path).Info("is open path")
		return true
	}

	// 如果开启了调试模式则不校验请求权限
	if a.isDebug() {
		log.Log().WithField("path", path).Info("debug mode")
		return true
	}

	//if len(a.Roles) == 0 {
	//	log.Log().WithField("path", path).WithField("accountId", accountId).
	//		Warning("no set any role for this account")
	//	return false
	//}

	// 通过判断角色是否包含信息
	if a.hasPermission(path, method) {
		log.Log().WithField("path", path).WithField("accountId", accountId).
			WithField("role", a.Roles).Warning("no set any role for this account")
		return true
	}
	// 默认是无权限访问
	return false
}

func contains(have []string, target string) bool {
	for _, a := range have {
		if a == target {
			return true
		}
	}
	return false
}
