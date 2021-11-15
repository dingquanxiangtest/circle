package mongo

import (
	"context"

	"git.internal.yunify.com/qxp/molecule/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type operatePermissionRepo struct {
}

// NewOperatePermissionRepo new repo
func NewOperatePermissionRepo() models.OperatePermissionRepo {
	return &operatePermissionRepo{}
}
func (o *operatePermissionRepo) getCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection("operate_permission")
}

// Delete delete by id
func (o *operatePermissionRepo) Delete(ctx context.Context, db *mongo.Database, formID, perGroupID string) error {
	_, err := o.getCollection(db).DeleteOne(ctx, bson.M{
		"per_group_id": perGroupID,
		"form_id":      formID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (o *operatePermissionRepo) Create(ctx context.Context, db *mongo.Database, operatePer *models.OperatePermission) error {
	_, err := o.getCollection(db).InsertOne(ctx, operatePer)
	if err != nil {
		return err
	}
	return nil
}
func (o *operatePermissionRepo) Get(ctx context.Context, db *mongo.Database, formID, perGroupID string) (*models.OperatePermission, error) {
	resp := &models.OperatePermission{}
	filter := bson.M{
		"per_group_id": perGroupID,
		"form_id":      formID,
	}
	err := o.getCollection(db).FindOne(ctx, filter).Decode(resp)
	if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *operatePermissionRepo) Update(ctx context.Context, db *mongo.Database, per *models.OperatePermission) error {
	filter := bson.M{"_id": per.ID}
	update := bson.M{
		"$set": bson.M{"authority": per.Authority},
	}
	_, err := o.getCollection(db).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
