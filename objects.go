package auth

import (
	"github.com/r2day/auth"
	"github.com/r2day/collections"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PermissionsModel 模型
// 记录角色对接口的操作细节
type PermissionsModel struct {
	// 角色编号
	RoleID string `json:"role_id" bson:"role_id"`
	// 所属应用编号
	AppID string `json:"app_id" bson:"app_id"`
	// 请求路径
	Path string `json:"path" bson:"path"`
	// 操作
	Operation OperationModel `json:"operation" bson:"operation"`
}

// OperationModel 操作模型
type OperationModel struct {
	// 读
	Read bool `json:"read" bson:"read"`
	// 写
	Write bool `json:"write" bson:"write"`
	// 改
	Update bool `json:"update" bson:"update"`
	// 详情
	Detail bool `json:"detail" bson:"detail"`
	// 删除
	Delete bool `json:"delete" bson:"delete"`
}

// RoleModel 模型
type RoleModel struct {
	// 基本的数据库模型字段，一般情况所有model都应该包含如下字段
	// 创建时（用户上传的数据为空，所以默认可以不传该值)
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	// 基本的数据库模型字段，一般情况所有model都应该包含如下字段
	Meta auth.MetaModel `json:"meta" bson:"meta"`
	// 名称
	Name string `json:"name" bson:"name"`
	// 工具列表
	Toolbar int `json:"toolbar" bson:"toolbar"`
	// 应用列表 toolbar
	// 存储应用的id
	// 通过应用id 快速获得应用列表
	Apps []string `json:"apps" bson:"apps"`
	// 权限
	Permissions []PermissionsModel `json:"permissions" bson:"permissions"`
}

// AppModel 模型
type AppModel struct {
	// AccessApi 可访问的api列表
	AccessAPI []collections.APIInfo `json:"access_api"  bson:"access_api"`
}
