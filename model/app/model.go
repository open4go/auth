package app

import (
	"github.com/open4go/model"
	"github.com/open4go/req5rsp/cst"
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
	modelName = "app"
)

// Model 订单信息
type Model struct {
	// 模型继承
	model.Model `json:"_" bson:"_"`
	// 基本的数据库模型字段，一般情况所有model都应该包含如下字段
	// 创建时（用户上传的数据为空，所以默认可以不传该值)
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	// 订单来源(系统根据订单来源终端自动赋值）
	Name      string      `json:"name"`
	Desc      string      `json:"desc"`
	Category  string      `json:"category"`
	Version   string      `json:"version"`
	Status    bool        `json:"status"`
	Namespace string      `json:"namespace"`
	Server    string      `json:"server"`
	Prefix    string      `json:"prefix"`
	AccessApi []ApiDetail `json:"api" bson:"api"`
}

// ApiDetail 接口详情
type ApiDetail struct {
	Path string `json:"path"`
	Name string `json:"name"`
	Desc string `json:"desc"`
	// 属性: 增, 删, 改， 查，列，展示，禁用
	// 属性: Create, Delete, Update, Read, List, Show, Disable,
	//对于你提到的权限（Create, Delete, Update, Read, List, Show, Disable），可以用二进制位来表示：
	//
	//1.	Create: 2^0 = 1 (00000001)
	//2.	Delete: 2^1 = 2 (00000010)
	//3.	Update: 2^2 = 4 (00000100)
	//4.	Read: 2^3 = 8 (00001000)
	//5.	List: 2^4 = 16 (00010000)
	//6.	Show: 2^5 = 32 (00100000)
	//7.	Disable: 2^6 = 64 (01000000)
	Attr cst.Permission `json:"attr"`
	// 不存储，仅用于展示效果及用于选择
	AttrBits []cst.Permission `json:"bits" bson:"-"`
	//Disable       bool   `json:"disable"`
	//CanViewDetail bool   `json:"can_view_detail"`
	//HideOnSidebar bool   `json:"hide_on_sidebar"`
}

// CheckAttr 检查是否启用了特定的权限
func (a ApiDetail) CheckAttr(bit cst.Permission) bool {
	return a.Attr&(1<<bit) != 0
}

// GetEnabledPermissions 根据权限值返回启用的权限名称
func (a ApiDetail) GetEnabledPermissions(attr cst.Permission) []string {
	var enabled []string
	for bit, name := range cst.PermissionNames {
		if attr&bit != 0 {
			enabled = append(enabled, name)
		}
	}
	return enabled
}

// ResourceName 返回资源名称
func (m *Model) ResourceName() string {
	return modelName
}

// CollectionName 返回表名称
func (m *Model) CollectionName() string {
	return collectionNamePrefix + modelName + collectionNameSuffix
}
