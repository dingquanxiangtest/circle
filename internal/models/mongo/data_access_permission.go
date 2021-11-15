package mongo

import (
	"context"

	"git.internal.yunify.com/qxp/molecule/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type dataAccessPermissionRepo struct {
}

// NewDataAccessPermissionRepo new repo
func NewDataAccessPermissionRepo() models.DataAccessPermissionRepo {
	return &dataAccessPermissionRepo{}
}
func (d *dataAccessPermissionRepo) getCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection("data_access_permission")
}

// Create 创建用户组
func (d *dataAccessPermissionRepo) Create(ctx context.Context, db *mongo.Database, dataAccessPer *models.DataAccessPermission) error {
	collection := d.getCollection(db)
	_, err := collection.InsertOne(ctx, dataAccessPer)
	if err != nil {
		return err
	}
	return nil
}

// Delete delete by id
func (d *dataAccessPermissionRepo) Delete(ctx context.Context, db *mongo.Database, formID, perGroupID string) error {
	_, err := d.getCollection(db).DeleteOne(ctx, bson.M{
		"per_group_id": perGroupID,
		"form_id":      formID,
	})
	if err != nil {
		return err
	}
	return nil
}

// Get get dataAccessPermission
func (d *dataAccessPermissionRepo) Get(ctx context.Context, db *mongo.Database, formID, perGroupID string) (*models.DataAccessPermission, error) {
	// 创建一个DataAccessPermission变量用来接收查询的结果
	resp := new(models.DataAccessPermission)
	filter := bson.M{
		"per_group_id": perGroupID,
		"form_id":      formID,
	}
	err := d.getCollection(db).FindOne(ctx, filter).Decode(resp)
	if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *dataAccessPermissionRepo) Update(ctx context.Context, db *mongo.Database, dataAccessPer *models.DataAccessPermission) error {
	update := bson.M{
		"$set": bson.M{
			"conditions": dataAccessPer.Conditions,
		},
	}
	filter := bson.M{
		"per_group_id": dataAccessPer.PerGroupID,
		"form_id":      dataAccessPer.FormID,
	}
	_, err := d.getCollection(db).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return err
}
