package roleCategory

import (
	"context"
	"time"

	db "github.com/r2day/auth"
	rtime "github.com/r2day/base/time"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ResourceName 返回资源名称
func (m *Model) ResourceName() string {
	return modelName
}

// CollectionName 返回表名称
func (m *Model) CollectionName() string {
	return collectionNamePrefix + modelName + collectionNameSuffix
}

// IncrementReference 更新
// https://www.mongodb.com/docs/manual/reference/operator/update/inc/
func (m *Model) IncrementReference(ctx context.Context, id string) error {
	coll := db.MDB.Collection(m.CollectionName())
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: objID}}
	// 设定更新时间
	m.Meta.UpdatedAt = rtime.FomratTimeAsReader(time.Now().Unix())

	result, err := coll.UpdateOne(ctx, filter,
		bson.D{{Key: "$set", Value: bson.D{{"reference", 1}}}})
	if err != nil {
		log.WithField("id", id).Error(err)
		return err
	}

	if result.MatchedCount < 1 {
		log.WithField("id", id).Warning("no matched record")
		return nil
	}
	return nil
}

// DecrementReference 更新
// https://www.mongodb.com/docs/manual/reference/operator/update/inc/
func (m *Model) DecrementReference(ctx context.Context, id string) error {
	coll := db.MDB.Collection(m.CollectionName())
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: objID}}
	// 设定更新时间
	m.Meta.UpdatedAt = rtime.FomratTimeAsReader(time.Now().Unix())

	result, err := coll.UpdateOne(ctx, filter,
		bson.D{{Key: "$set", Value: bson.D{{"reference", -1}}}})
	if err != nil {
		log.WithField("id", id).Error(err)
		return err
	}

	if result.MatchedCount < 1 {
		log.WithField("id", id).Warning("no matched record")
		return nil
	}
	return nil
}
