package app

import (
	"errors"
	"github.com/open4go/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *Model) GetByOrderID(id string) ([]*Model, error) {

	results := make([]*Model, 0)
	coll := m.Context.Handler.Collection(m.Context.Collection)
	filter := bson.D{{Key: "identity.order_no", Value: id}}

	// 获取数据列表
	cursor, err := coll.Find(m.Context.Context, filter)
	if errors.Is(err, mongo.ErrNoDocuments) {
		log.Log(m.Context.Context).WithField("id", id).Error(err)
		return nil, err
	}
	if err != nil {
		log.Log(m.Context.Context).WithField("id", id).Error(err)
		return nil, err
	}
	if err = cursor.All(m.Context.Context, &results); err != nil {
		log.Log(m.Context.Context).WithField("id", id).Error(err)
		return nil, err
	}
	return results, nil
}
