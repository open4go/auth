package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/open4go/auth/model/app"
	"github.com/open4go/auth/model/role"
	"github.com/open4go/log"
	"github.com/open4go/req5rsp/cst"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RoleManager 角色管理
type RoleManager struct {
	RedisPrefix string
	App         []app.Model `json:"app"`
}

// fetchRoleList 加载所有角色到内存中
// roleID 可以直接通过jwt token解析出来得到
// 其他场景也可以通过账户id直接查询redis缓存
func (r *RoleManager) fetchRoleList(ctx context.Context) ([]*role.Model, error) {
	// 初始化模型
	m := &role.Model{}
	// 定义列表
	s := make([]*role.Model, 0)

	// 获取表操作handler
	h := m.Init(ctx, GetDBHandler(ctx), m.CollectionName())
	// 执行查询
	// counter 表示在该过滤条件下的总数
	//objIds := covertSliceToObjectID(ctx, roleID)
	//log.Log(ctx).WithField("ids", objIds).Info("=============222")
	// s 拉取到的列表绑定到s
	_, err := h.GetListWithOpt(bson.M{
		"status": true,
	}, &s, nil)
	if err != nil {
		log.Log(ctx).Error(err)
		return nil, err
	}
	return s, nil
}

// fetchAppList 加载所有应用到内存中
// roleID 可以直接通过jwt token解析出来得到
// 其他场景也可以通过账户id直接查询redis缓存
func (r *RoleManager) fetchRoleListByIDs(ctx context.Context, roleID []string) ([]*role.Model, error) {
	// 初始化模型
	m := &role.Model{}
	// 定义列表
	s := make([]*role.Model, 0)

	// 获取表操作handler
	h := m.Init(ctx, GetDBHandler(ctx), m.CollectionName())
	// 执行查询
	// counter 表示在该过滤条件下的总数
	objIds := covertSliceToObjectID(ctx, roleID)
	log.Log(ctx).WithField("ids", objIds).Info("=============222")
	// s 拉取到的列表绑定到s
	_, err := h.GetListWithOpt(bson.M{
		"_id": bson.M{"$in": objIds},
	}, &s, nil)
	if err != nil {
		log.Log(ctx).Error(err)
		return nil, err
	}
	return s, nil
}

// reloadRoles 重新加载所有角色的权限（推荐全量替换）
func (r *RoleManager) reloadRoles(ctx context.Context, roles []*role.Model) error {
	redisClient := GetRedisAuthHandler(ctx)

	for _, roleModel := range roles {
		roleKey := fmt.Sprintf("%s:roles:permissions:%s", r.RedisPrefix, roleModel.ID.Hex())

		// 【关键】先删除旧的整个 Hash，避免残留已删除的权限
		if err := redisClient.Del(ctx, roleKey).Err(); err != nil {
			log.Log(ctx).WithField("key", roleKey).Error(err)
			// 可以继续，不一定返回错误
		}

		// 写入新的权限
		for _, api := range roleModel.ApiDetail {
			err := redisClient.HSet(ctx, roleKey, api.Path, uint(api.Attr)).Err()
			if err != nil {
				log.Log(ctx).WithField("key", roleKey).WithField("path", api.Path).Error(err)
				return err
			}
		}

		// 可选：设置过期时间（防止数据长期不更新）
		//redisClient.Expire(ctx, roleKey, time.Hour*24)
	}
	return nil
}

// fetchPathsByRoleID 设置缓存避免重复查询用户角色
// 仅当用户角色发生更新后再进行查询同步到redis
func (r *RoleManager) fetchPathsByRoleID(ctx context.Context, role string) ([]string, error) {
	roleKey := fmt.Sprintf("%s:roles:permissions:%s", r.RedisPrefix, role)
	results, err := GetRedisAuthHandler(ctx).HGetAll(ctx, roleKey).Result()
	if err != nil {
		log.Log(ctx).WithField("key", roleKey).
			Error(err)
		return nil, err
	}

	paths := make([]string, 0)
	for path := range results {
		paths = append(paths, path)
	}
	return paths, nil
}

// SignIn 设置缓存避免重复查询用户角色
// 仅当用户角色发生更新后再进行查询同步到redis
func (r *RoleManager) SignIn(ctx context.Context, accountId string, roles []string) error {
	roleModel, err := r.fetchRoleListByIDs(ctx, roles)
	if err != nil {
		log.Log(ctx).Error(err)
		return err
	}
	log.Log(ctx).WithField("roleIds", roles).WithField("roleModel", roleModel).Debug("before setRoles")
	err = r.setRoles(ctx, accountId, roleModel)
	if err != nil {
		log.Log(ctx).WithField("roleIds", roles).WithField("roleModel", roleModel).Error("call setRoles")
		return err
	}
	return nil
}

func (r *RoleManager) SignOut(ctx context.Context, accountId string) error {
	roleKey := fmt.Sprintf("%s:account:to:role:%s", r.RedisPrefix, accountId)
	err := GetRedisAuthHandler(ctx).Expire(ctx, roleKey, time.Second*10).Err()
	if err != nil {
		return err
	}
	return nil
}

// Verify 校验账户是否有权访问 path。
// WebSocket 接口强制按 GET 校验：具备 List 或 Read 任一即可。
func (r *RoleManager) Verify(ctx context.Context, path string, accountId string, method string, isSingleResource bool) (bool, error) {
	isWS := wsFlagFromCtx(ctx) || isWSEndpoint(path)
	checkMethod := method
	if isWS {
		checkMethod = "GET"
	}

	rolesOfThisAccount, err := r.fetchRolesFromCache(ctx, accountId)
	if err != nil {
		log.Log(ctx).
			WithField("path", path).
			WithField("accountId", accountId).
			WithField("method", method).
			WithField("is_ws", isWS).
			WithField("error", err).
			Error("Verify: failed to load account roles from cache")
		return false, err
	}

	p := translateHTTPMethodToPermission(checkMethod, isSingleResource)
	log.Log(ctx).
		WithField("path", path).
		WithField("accountId", accountId).
		WithField("request_method", method).
		WithField("check_method", checkMethod).
		WithField("is_ws", isWS).
		WithField("is_single_resource", isSingleResource).
		WithField("expect_perm", formatPermission(int(p))).
		WithField("roles", rolesOfThisAccount).
		Debug("Verify: start permission check")

	if len(rolesOfThisAccount) == 0 {
		log.Log(ctx).
			WithField("path", path).
			WithField("accountId", accountId).
			WithField("is_ws", isWS).
			Error("Verify: account has no roles in cache (not signed in, or role cache expired)")
		return false, nil
	}

	isCanAccess, err := r.canAccess(ctx, rolesOfThisAccount, path, p, isWS)
	if err != nil {
		log.Log(ctx).
			WithField("path", path).
			WithField("accountId", accountId).
			WithField("is_ws", isWS).
			WithField("error", err).
			Error("Verify: canAccess failed")
		return false, err
	}
	return isCanAccess, nil
}

// 设置缓存避免重复查询用户角色
// 仅当用户角色发生更新后再进行查询同步到redis
// 当后台重新调整该账户的角色时应该刷新该缓存
func (r *RoleManager) setRoles(ctx context.Context, account string, roles []*role.Model) error {
	roleKey := fmt.Sprintf("%s:account:to:role:%s", r.RedisPrefix, account)
	for _, i := range roles {
		err := GetRedisAuthHandler(ctx).SAdd(ctx, roleKey, i.ID.Hex()).Err()
		if err != nil {
			log.Log(ctx).WithField("key", roleKey).
				Error(err)
			return err
		}
	}
	// 默认两个小时过期
	err := GetRedisAuthHandler(ctx).Expire(ctx, roleKey, time.Hour*2).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RoleManager) IsOnline(ctx context.Context, account string) bool {
	roleKey := fmt.Sprintf("%s:account:to:role:%s", r.RedisPrefix, account)
	result, err := GetRedisAuthHandler(ctx).Exists(ctx, roleKey).Result()
	if err != nil {
		return false
	}
	// 返回1 则表示存在
	return result == 1
}

// KeepAlive 当用户有操作行为则，继续保持key的生命周期延续2小时
func (r *RoleManager) KeepAlive(ctx context.Context, account string) error {
	roleKey := fmt.Sprintf("%s:account:to:role:%s", r.RedisPrefix, account)
	err := GetRedisAuthHandler(ctx).Expire(ctx, roleKey, time.Hour*2).Err()
	if err != nil {
		return err
	}
	// 返回1 则表示存在
	return nil
}

// OnlineTime 在线时长统计
func (r *RoleManager) OnlineTime(ctx context.Context, account string) (time.Duration, error) {
	roleKey := fmt.Sprintf("%s:account:to:role:%s", r.RedisPrefix, account)
	t, err := GetRedisAuthHandler(ctx).TTL(ctx, roleKey).Result()
	if err != nil {
		return 0, err
	}
	// 返回1 则表示存在
	return t, nil
}

func (r *RoleManager) fetchRolesFromCache(ctx context.Context, account string) ([]string, error) {
	roleKey := fmt.Sprintf("%s:account:to:role:%s", r.RedisPrefix, account)
	roles, err := GetRedisAuthHandler(ctx).SMembers(ctx, roleKey).Result()
	if err != nil {
		log.Log(ctx).WithField("key", roleKey).
			Error(err)
		return nil, err
	}
	return roles, nil
}

func (r *RoleManager) MyRoles(ctx context.Context, account string) ([]string, error) {
	return r.fetchRolesFromCache(ctx, account)
}

// FetchAllPaths 获取路径key
func (r *RoleManager) FetchAllPaths(ctx context.Context, account string) (map[string]int, []string, error) {
	roles, err := r.fetchRolesFromCache(ctx, account)
	if err != nil {
		log.Log(ctx).WithField("account", account).
			Error(err)
		return nil, nil, err
	}
	allPaths := make([]string, 0)
	path2attr := make(map[string]int)
	for _, i := range roles {
		paths, err := r.fetchPathsByRoleID(ctx, i)
		if err != nil {
			log.Log(ctx).WithField("account", account).WithField("role", i).
				Error(err)
			continue
		}
		// 计算当前路径的权限并与map中该路径在其他角色的权限进行位运算（如果其他角色也有该路径，则合并）
		for _, path := range paths {
			roleAttr, _, err2 := r.getAttrByPathAndRole(ctx, i, path)
			if err2 != nil {
				//log.Log(ctx).WithField("i", i).WithField("path", path).Debug("this role no match, try next")
				continue
			}
			path2attr[path] = roleAttr | path2attr[path]
		}
		allPaths = append(allPaths, paths...)
	}
	return path2attr, allPaths, nil
}

func (r *RoleManager) Menu(ctx context.Context, paths2Attr map[string]int) []MenuTree {
	return r.convertPathsToStructure(paths2Attr)
}

// convertPathsToStructure 将路径数组转换为所需的结构
func (r *RoleManager) convertPathsToStructure(paths2Attr map[string]int) []MenuTree {
	menuMap := make(map[string][]string)

	for path, attr := range paths2Attr {

		if !isShow(attr) {
			continue
		}

		// 去除路径前的斜杠，并拆分路径
		parts := strings.Split(strings.TrimPrefix(path, "/"), "/")

		if len(parts) != 4 {
			log.Log(context.TODO()).WithField("path", path).Debug("Invalid path format so ignore")
			continue
		}

		// 构建主菜单名称
		// 只取 /版本/命名空间/服务名 （作为主目录）
		// 例如：/v1/system/auth/account
		// 那么：v1.system.auth 作为服务名/即主目录名
		mainMenuName := fmt.Sprintf("menu.%s._name", strings.Join(parts[:len(parts)-1], "."))

		// 构建子菜单名称
		subMenuName := fmt.Sprintf("menu.%s", strings.Join(parts, "."))

		// 将子菜单名称添加到主菜单名称对应的列表中
		menuMap[mainMenuName] = append(menuMap[mainMenuName], subMenuName)
	}

	// 构建最终的 MTree 列表
	var mTree []MenuTree
	for mainMenuName, subMenuList := range menuMap {
		mTree = append(mTree, MenuTree{
			Name:    mainMenuName,
			SubMenu: subMenuList,
		})
	}
	return mTree
}

func (r *RoleManager) canAccess(ctx context.Context, roles []string, path string, expect cst.Permission, isWS bool) (bool, error) {
	pathHits := make([]string, 0)
	missRoles := make([]string, 0)
	bitMismatch := make([]string, 0)
	var lastRedisErr error
	redisErrCount := 0

	for _, roleID := range roles {
		roleAttr, _, err2 := r.getAttrByPathAndRole(ctx, roleID, path)
		if err2 != nil {
			if !errors.Is(err2, redis.Nil) {
				redisErrCount++
				lastRedisErr = err2
				log.Log(ctx).
					WithField("role", roleID).
					WithField("path", path).
					WithField("error", err2).
					Error("canAccess: failed to read role permission from redis")
			} else {
				missRoles = append(missRoles, roleID)
			}
			continue
		}

		pathHits = append(pathHits, fmt.Sprintf("%s=%s", roleID, formatPermission(roleAttr)))

		allowed := false
		if isWS {
			// WS 只需 GET：List（集合 GET）或 Read（单资源 GET）任一即可
			allowed = hasGETPermission(roleAttr)
		} else {
			allowed = roleAttr&int(expect) == int(expect)
		}
		if allowed {
			log.Log(ctx).
				WithField("role", roleID).
				WithField("path", path).
				WithField("is_ws", isWS).
				WithField("role_attr", formatPermission(roleAttr)).
				WithField("expect_perm", formatPermission(int(expect))).
				Debug("canAccess: permission granted")
			return true, nil
		}
		bitMismatch = append(bitMismatch, fmt.Sprintf("%s=%s", roleID, formatPermission(roleAttr)))
	}

	if redisErrCount == len(roles) && lastRedisErr != nil {
		return false, lastRedisErr
	}

	reason := "path not found in any role permission hash"
	if len(pathHits) > 0 {
		if isWS {
			reason = "path found but role has neither List nor Read (GET) permission"
		} else {
			reason = "path found but required permission bits are missing"
		}
	}

	log.Log(ctx).
		WithField("path", path).
		WithField("is_ws", isWS).
		WithField("expect_perm", formatPermission(int(expect))).
		WithField("roles", roles).
		WithField("path_hits", pathHits).
		WithField("bit_mismatch", bitMismatch).
		WithField("roles_without_path", missRoles).
		WithField("reason", reason).
		Error("canAccess: access denied")
	return false, nil
}

func (r *RoleManager) getAttrByPathAndRole(ctx context.Context, role string, path string) (int, bool, error) {
	roleKey := fmt.Sprintf("%s:roles:permissions:%s", r.RedisPrefix, role)
	roleAttr, err := GetRedisAuthHandler(ctx).HGet(ctx, roleKey, path).Int()
	if err != nil {
		return 0, false, err
	}
	return roleAttr, false, nil
}

// Reload 完全重新加载角色权限
func (r *RoleManager) Reload(ctx context.Context) {
	log.Log(ctx).Info("Starting RoleManager full reload...")

	allRoles, err := r.fetchRoleList(ctx)
	if err != nil {
		log.Log(ctx).Error(err)
		return
	}

	if err = r.reloadRoles(ctx, allRoles); err != nil {
		log.Log(ctx).Error("loadRoles failed")
		return
	}

	log.Log(ctx).WithField("roleCount", len(allRoles)).Info("RoleManager reload completed successfully")
}

func covertSliceToObjectID(ctx context.Context, slice []string) []*primitive.ObjectID {
	objIds := make([]*primitive.ObjectID, 0)
	for _, i := range slice {
		objID, err := primitive.ObjectIDFromHex(i)
		if err != nil {
			log.Log(ctx).WithField("originID", i).
				Error(err)
			continue
		}
		objIds = append(objIds, &objID)
	}
	return objIds
}

// translateHTTPMethodToPermission 根据HTTP方法和路径来确定权限
func translateHTTPMethodToPermission(method string, isSingleResource bool) cst.Permission {
	method = strings.ToUpper(method)

	// 如果路径包含 ':id' 或者其他参数化部分，按特定方式处理
	if isSingleResource {
		switch method {
		case "DELETE":
			return cst.Delete
		case "PUT", "PATCH":
			return cst.Update
		case "GET":
			//log.Log(context.TODO()).WithField("method", "GET").Debug("hit method ============")
			return cst.Read
		}
	} else {
		// 如果路径不包含 ':id'，按集合资源方式处理
		switch method {
		case "POST":
			return cst.Create
		case "GET":
			return cst.List
		}
	}

	// 对于不匹配的方法，返回0或其他合适的值
	return 0
}

// 读取redis中的权限信息
func isShow(roleAttr int) bool {
	return roleAttr&int(cst.Show) == int(cst.Show)
}
