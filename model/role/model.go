package role

import (
	"github.com/open4go/auth/model/app"
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
	collectionNameSuffix = "_detail"
	// 这个需要用户根据具体业务完成设定
	modelName = "role"
)

// Model 订单信息
type Model struct {
	// 模型继承
	model.Model `json:"_" bson:"_"`
	// 基本的数据库模型字段，一般情况所有model都应该包含如下字段
	// 创建时（用户上传的数据为空，所以默认可以不传该值)
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	// 订单来源(系统根据订单来源终端自动赋值）
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Status   bool   `json:"status"`
	Category string `json:"category"`
	// 通过应用id 快速获得应用列表
	Apps []string `json:"app" bson:"app"`
	// 权限设置
	// 用于展示当角色加入对应的app后自动渲染出每一个子应用的path
	// 可以单独设置每一个path在该角色的权限
	ApiDetail []app.ApiDetail `json:"api" bson:"api"`
}

// ResourceName 返回资源名称
func (m *Model) ResourceName() string {
	return modelName
}

// CollectionName 返回表名称
func (m *Model) CollectionName() string {
	return collectionNamePrefix + modelName + collectionNameSuffix
}
