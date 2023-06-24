package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/r2day/collections"
	log "github.com/sirupsen/logrus"
)

// 遍历账号中所拥有的所有Roles
// 1. 通过角色找到其应用列表，并且决定是否显示其在导航menu上或者是否禁用api
// 2. 通过permissions 找到每一个接口的具体操作方法，是否具备相关操作权限
func iterRoles(c *gin.Context, roles []*RoleModel, logCtx *log.Entry, permissions []PermissionsModel,
	rolesName []string, maxAccessLevel int, myRolesKey string, keyPrefix string, keyNames string,
	path2roles map[string][]string) ([]string, []PermissionsModel, int, int) {
	toolBar := 0
	for _, role := range roles {

		if !role.Meta.Status {
			logCtx.WithField("role_name", role.Name).
				Warning("current role is disable, so no going to use it")
			break
		}

		if role.Toolbar > toolBar {
			toolBar = role.Toolbar
		}

		// 将所有角色下的接口的权限管理进行统一管理
		permissions = append(permissions, role.Permissions...)

		// 加入key
		// 以便退出登陆后删除
		rolesName = append(rolesName, role.Name)

		// 遍历寻找最大的用户角色等级
		if role.Meta.AccessLevel > uint(maxAccessLevel) {
			maxAccessLevel = int(role.Meta.AccessLevel)
		}

		// 获取当前用户角色的应用列表
		//appM := &AppModel{}
		//apps, err := appM.GetMany(c.Request.Context(), role.Apps)
		//if err != nil {
		//	logCtx.WithField("role", role.Name).WithField("apps", apps).
		//		Warning("no found any apps, try next role")
		//	continue
		//}

		// 将处理好的角色名称也加入到缓存中
		// 使用角色id 避免用户输入特殊字符无法作为redis key
		err = RDB.SAdd(c.Request.Context(), myRolesKey, role.ID.Hex()).Err()
		if err != nil {
			logCtx.Error(err)
			continue
		}

		accessAPIList := make([]collections.APIInfo, 0)

		// 遍历所有授权的应用
		for _, app := range apps {
			// 获得所有应用下的api列表
			accessAPIList = append(accessAPIList, app.AccessAPI...)
		}

		//setMenu(c, accessAPIList, keyPrefix, err, logCtx, keyNames, path2roles, role)

	}

	return rolesName, permissions, maxAccessLevel, toolBar
}
