package app

import (
	"github.com/open4go/model"
	"github.com/r2day/collections"
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
	modelName = "app"
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
	// 名称
	Name string `json:"name" bson:"name"`
	// 应用描述
	Desc string `json:"desc" bson:"desc"`
	// 分类/ 亦或则是分组等
	Category string `json:"category" bson:"category"`
	// 前缀
	UrlPrefix string `json:"url_prefix" bson:"url_prefix"`
	// AccessApi 可访问的api列表
	AccessAPI []collections.APIInfo `json:"access_api"  bson:"access_api"`

	// ApiAttr 可访问的api ApiAttr列表
	ApiAttr []ApiAttribute `json:"api_attr"  bson:"api_attr"`
}

type ApiAttribute struct {
	// 名称
	Name string `json:"name" bson:"name"`
	// 应用描述
	Value string `json:"value" bson:"value"`
}
