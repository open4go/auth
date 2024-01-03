package account

import (
	"github.com/open4go/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	// CollectionNamePrefix 数据库表前缀
	// 可以根据具体业务的需要进行定义
	// 例如: sys_, scm_, customer_, order_ 等
	collectionNamePrefix = "auth_"
	// CollectionNameSuffix 后缀
	// 例如, _log, _config, _flow,
	collectionNameSuffix = "_manage"
	// 这个需要用户根据具体业务完成设定
	modelName = "account"
)

// 每一个应用表示一个大的模块，通常其子模块是一个个接口
// 是有系统默认设定，用户无需修改
// 用户只需要在创建角色的时候选择好需要的应用即可
// 用户选择所需要的应用后->完成角色创建->系统自动拷贝应用具体信息到角色下
// 此时用户可以针对当前的角色中具体的项再自行选择是否移除部分接口，从而进行更精细的权限管理

// Model 模型
type Model struct {
	// 继承
	model.Model
	// 基本的数据库模型字段，一般情况所有model都应该包含如下字段
	// 创建时（用户上传的数据为空，所以默认可以不传该值)
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	// 分类/ 亦或则是分组等
	Category string `json:"category" bson:"category"`
	// IsAdmin 是否是管理员
	// 一般standalone 模式下，是通过读取部署时配置的 ADMIN_PHONE
	// 如果与配置的手机号匹配，那么就可以定义为管理员
	// 如果是其他的模式，一般需要超级商户平台授权后才能成为管理员
	IsAdmin bool `json:"is_admin"  bson:"is_admin"`
	// 关键信息
	// 手机号
	Phone string `json:"phone"`
	// 密码
	Password string `json:"password"  bson:"password"`
	// 是否开启审核
	IsRequiredApprove bool `json:"is_required_approve" bson:"is_required_approve"`
	// 更多信息
	// 账号名称
	Name string `json:"name"`
	// 描述
	Desc string `json:"desc"  bson:"desc"`
	// 邮箱
	Email string `json:"email"  bson:"email"`
	// 角色名称列表
	Roles []string `json:"role"  bson:"role"`
	// IsOnline 是否在线
	IsOnline bool `json:"is_online" bson:"-"`

	OtpEnabled  bool   `json:"opt_enabled" bson:"opt_enabled"`
	OtpVerified bool   `json:"opt_verified" bson:"opt_verified"`
	OtpSecret   string `json:"opt_secret" bson:"opt_secret"`
	OtpAuthUrl  string `json:"opt_auth_url" bson:"opt_auth_url"`
}
