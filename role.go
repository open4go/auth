package auth

import (
	"fmt"
	"github.com/open4go/auth/model/app"
	"github.com/open4go/auth/model/role"
	"github.com/open4go/log"
	"github.com/open4go/req5rsp/cst"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"strings"
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

// 设置缓存避免重复查询用户角色
// 仅当用户角色发生更新后再进行查询同步到redis
func (r *RoleManager) loadRoles(ctx context.Context, roles []*role.Model) error {
	for _, i := range roles {
		roleKey := fmt.Sprintf("%s:roles:permissions:%s", r.RedisPrefix, i.ID.Hex())
		for _, j := range i.ApiDetail {
			err := GetRedisAuthHandler(ctx).HSet(ctx, roleKey, j.Path, j.Attr).Err()
			if err != nil {
				log.Log(ctx).WithField("key", roleKey).
					Error(err)
				return err
			}
		}
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
	for path, _ := range results {
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
	err = r.setRoles(ctx, accountId, roleModel)
	if err != nil {
		return err
	}
	return nil
}

// Verify 设置缓存避免重复查询用户角色
// 仅当用户角色发生更新后再进行查询同步到redis
func (r *RoleManager) Verify(ctx context.Context, path string, accountId string, method string, isSingleResource bool) (bool, error) {
	rolesOfThisAccount, err := r.fetchRolesFromCache(ctx, accountId)
	if err != nil {
		return false, err
	}

	log.Log(ctx).WithField("method", method).WithField("path", path).Debug("check the params =====")
	p := translateHTTPMethodToPermission(method, isSingleResource)
	// TODO 这里需要将p的值进行<< 左移运算
	isCanAccess, err := r.canAccess(ctx, rolesOfThisAccount, path, p)
	if err != nil {
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
	return nil
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

// FetchAllPaths 获取路径key
func (r *RoleManager) FetchAllPaths(ctx context.Context, account string) ([]string, error) {
	roles, err := r.fetchRolesFromCache(ctx, account)
	if err != nil {
		log.Log(ctx).WithField("account", account).
			Error(err)
		return nil, err
	}
	allPaths := make([]string, 0)
	for _, i := range roles {
		paths, err := r.fetchPathsByRoleID(ctx, i)
		if err != nil {
			log.Log(ctx).WithField("account", account).
				Error(err)
			return nil, err
		}
		allPaths = append(allPaths, paths...)
	}
	return allPaths, nil
}

func (r *RoleManager) canAccess(ctx context.Context, roles []string, path string, expect cst.Permission) (bool, error) {
	for _, i := range roles {
		roleKey := fmt.Sprintf("%s:roles:permissions:%s", r.RedisPrefix, i)
		roleAttr, err := GetRedisAuthHandler(ctx).HGet(ctx, roleKey, path).Int()
		if err != nil {
			log.Log(ctx).WithField("key", roleKey).
				Error(err)
			// try next role
			// 检查下一个角色是否有操作权限
			continue
		}

		// 如果当前角色的权限已经满足通行
		// 则不必继续检查其他角色，直接返回
		// 使用位操作检查权限位是否匹配
		if roleAttr&int(expect) == int(expect) {
			return true, nil
		}
	}
	// 没有任何一个角色满足
	// 因此返回false
	return false, nil
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
			log.Log(context.TODO()).WithField("method", "GET").Info("hit method ============")
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
