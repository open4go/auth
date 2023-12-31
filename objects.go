package auth

import (
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
	// 请求名称
	Name string `json:"name" bson:"name"`
	// 操作包含read/write/update/delete/list
	Operation []string `json:"operation" bson:"operation"`
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

// MetaModel 元模型
type MetaModel struct {
	// 商户号
	MerchantID string `json:"merchant_id" bson:"merchant_id"`
	// 创建者
	AccountID string `json:"account_id" bson:"account_id"`
	// 创建时间
	CreatedAt string `json:"created_at" bson:"created_at"`
	// 更新时间
	UpdatedAt string `json:"updated_at" bson:"updated_at"`
	// 状态
	Status bool `json:"status"`
	// 根据角色的最低级别写入
	AccessLevel uint `json:"access_level" bson:"access_level"`
}

// RoleModel 模型
type RoleModel struct {
	// 基本的数据库模型字段，一般情况所有model都应该包含如下字段
	// 创建时（用户上传的数据为空，所以默认可以不传该值)
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	// 基本的数据库模型字段，一般情况所有model都应该包含如下字段
	Meta MetaModel `json:"meta" bson:"meta"`
	// 名称
	Name string `json:"name" bson:"name"`
	// 工具列表
	Toolbar int `json:"toolbar" bson:"toolbar"`
	// 应用列表 toolbar
	// 存储应用的id
	// 通过应用id 快速获得应用列表
	Apps []string `json:"app" bson:"app"`
	// 权限
	Permissions []PermissionsModel `json:"permissions" bson:"permissions"`
}

// AppModel 模型
type AppModel struct {
	// AccessApi 可访问的api列表
	AccessAPI []collections.APIInfo `json:"access_api"  bson:"access_api"`
}

type LoginInfo struct {
	// 命名空间
	// 可是商户号
	Namespace string `json:"namespace"`
	// 账号id
	AccountId string `json:"account_id"  bson:"account_id"`
	// 可以是手机号
	UserId string `json:"user_id"  bson:"user_id"`
	// 用户名
	UserName string `json:"user_name"  bson:"user_name"`
	// Avatar 用户头像
	Avatar string `json:"avatar"`
	// LoginType 登陆类型
	LoginType string `json:"login_type"  bson:"login_type"`
	// LoginLevel 登陆用户等级
	LoginLevel string `json:"login_level"  bson:"login_level"`
}
