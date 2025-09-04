package tenant

import (
	"github.com/open4go/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	// CollectionNamePrefix 数据库表前缀
	// 可以根据具体业务的需要进行定义
	// 例如: sys_, scm_, customer_, order_ 等
	collectionNamePrefix = "auth_"
	// CollectionNameSuffix 后缀
	// 例如, _log, _config, _flow,
	collectionNameSuffix = "_detail"
	// 这个需要用户根据具体业务完成设定
	modelName = "tenant"
)

// Model 订单信息
type Model struct {
	// 模型继承
	model.Model `json:"_" bson:"_"`
	// 基本的数据库模型字段，一般情况所有model都应该包含如下字段
	// 创建时（用户上传的数据为空，所以默认可以不传该值)
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name         string             `bson:"name" json:"name"`
	MerchantID   string             `bson:"merchant_id" json:"merchant_id"` // 唯一标识符
	Status       string             `bson:"status" json:"status"`           // active, suspended, pending
	CustomDomain string             `bson:"custom_domain" json:"custom_domain"`
	Plan         string             `bson:"plan" json:"plan"` // 套餐类型
	Settings     Setting            `bson:"settings" json:"settings"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
	Phone        string             `bson:"phone" json:"phone"` // 手机号/用户等了
}

// Setting 租户设置
type Setting struct {
	Theme        string   `bson:"theme" json:"theme"`
	Locale       string   `bson:"locale" json:"locale"`
	AllowedHosts []string `bson:"allowed_hosts" json:"allowed_hosts"`
	MaxUsers     int      `bson:"max_users" json:"max_users"`
	StorageLimit int64    `bson:"storage_limit" json:"storage_limit"` // 存储限制(字节)
}

// ResourceName 返回资源名称
func (m *Model) ResourceName() string {
	return modelName
}

// CollectionName 返回表名称
func (m *Model) CollectionName() string {
	return collectionNamePrefix + modelName + collectionNameSuffix
}
