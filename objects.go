package auth


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
