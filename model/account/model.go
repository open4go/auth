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
	// 模型继承
	model.Model `json:"_" bson:"_"`
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
	// 手机号（新记录存空，接口返回脱敏值）
	Phone       string `json:"phone" bson:"phone"`
	PhoneEnc    string `json:"-" bson:"phone_enc,omitempty"`
	PhoneHash   string `json:"-" bson:"phone_hash,omitempty"`
	PhoneSuffix string `json:"phone_suffix,omitempty" bson:"phone_suffix,omitempty"`
	PhoneSec    string `json:"-" bson:"phone_sec,omitempty"`
	// 用户id 在这个表中使用userId
	// 但是在内部的meta.account_id 是一样的值
	// 但是meta 值是登陆者操作的人的id
	UserID string `json:"user_id"  bson:"user_id"`
	// 密码
	Password string `json:"password"  bson:"password"`
	// 是否开启审核
	IsRequiredApprove bool `json:"is_required_approve" bson:"is_required_approve"`
	// 更多信息
	// 账号名称
	Name string `json:"name"`
	// 描述
	Desc string `json:"desc"  bson:"desc"`
	// 邮箱（新记录存空，接口返回脱敏值）
	Email     string `json:"email" bson:"email"`
	EmailEnc  string `json:"-" bson:"email_enc,omitempty"`
	EmailHash string `json:"-" bson:"email_hash,omitempty"`
	EmailSec  string `json:"-" bson:"email_sec,omitempty"`
	// 角色id列表
	RoleID []string `json:"roles"  bson:"roles"`
	// 角色名称列表
	Roles []string `json:"role"  bson:"role"`
	// IsOnline 是否在线
	IsOnline bool `json:"is_online" bson:"-"`

	OtpEnabled  bool   `json:"opt_enabled" bson:"opt_enabled"`
	OtpVerified bool   `json:"opt_verified" bson:"opt_verified"`
	OtpSecret   string `json:"opt_secret" bson:"opt_secret"`
	OtpAuthUrl  string `json:"opt_auth_url" bson:"opt_auth_url"`

	// 	登陆信息
	LoginInfo LoginInfoType `json:"login_info" bson:"login_info"`

	// 安全信息（不直接存储到库里）
	SafeInfo SafeInfoType `json:"safe" bson:"-"`

	// 头像信息
	Avatar string `json:"avatar"  bson:"avatar"`

	// Prefs 后台界面偏好，跟账号走（语言、主题、导航、字体）
	Prefs Prefs `json:"prefs" bson:"prefs"`
}

// Prefs 管理员后台的语言、主题、导航和字体，登录后各端共用
type Prefs struct {
	Locale       string `json:"locale" bson:"locale"`
	ThemeName    string `json:"theme_name" bson:"theme_name"`
	ThemeMode    string `json:"theme_mode" bson:"theme_mode"` // light | dark | auto
	SidebarStyle string `json:"sidebar_style" bson:"sidebar_style"`
	Font         string `json:"font" bson:"font"` // system | song | kai
}

type LoginInfoType struct {
	// CurrentIPAddr 当前登陆ip
	CurrentIPAddr string `json:"current_ip_addr"  bson:"current_ip_addr"`
	// LastIPAddr 上一次登陆ip
	LastIPAddr string `json:"last_ip_addr"  bson:"last_ip_addr"`
	// LoginCount 登陆次数
	LoginCount uint `json:"login_count"  bson:"login_count"`
	// LoginTime 登陆时长（秒）
	LoginTime uint `json:"login_time"  bson:"login_time"`
	// LastDeviceInfo 登陆时长（秒）
	CurrentDeviceInfo string `json:"current_device_info"  bson:"current_device_info"`
	// LastDeviceInfo 登陆时长（秒）
	LastDeviceInfo string `json:"last_device_info"  bson:"last_device_info"`
}

type SafeInfoType struct {
	// 仅用于读取前端参数
	Password string `json:"password"  bson:"-"`
}
