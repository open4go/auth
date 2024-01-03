package account

import (
	"go.mongodb.org/mongo-driver/bson"
)

// ResourceName 返回资源名称
func (m *Model) ResourceName() string {
	return modelName
}

// CollectionName 返回表名称
func (m *Model) CollectionName() string {
	return collectionNamePrefix + modelName + collectionNameSuffix
}

// FindByPhone 通过手机号查找到账号信息
func (m *Model) FindByPhone(phone string) (*Model, error) {

	result := &Model{}
	filter := bson.D{{Key: "phone", Value: phone}}
	err := m.GetBy(result, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindByAccountId 通过手机号查找到账号信息
func (m *Model) FindByAccountId(accountID string) (*Model, error) {
	result := &Model{}
	filter := bson.D{{Key: "meta.account_id", Value: accountID}}
	err := m.GetBy(result, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}
