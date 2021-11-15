package mongo

import (
	"context"
	"git.internal.yunify.com/qxp/molecule/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type filterRepo struct {
}

func (a *filterRepo) Insert(c context.Context, db *mongo.Database, req *models.Filter) error {

	_, err := a.getCollection(db).InsertOne(c, req)
	return err
}

func (a *filterRepo) Update(c context.Context, db *mongo.Database, req *models.Filter) error {

	update := bson.M{
		"$set": bson.M{
			"field_json": req.FieldJSON,
			"web_schema": req.WebSchema,
		},
	}
	filter := bson.M{
		"per_group_id": req.PerGroupID,
		"form_id":      req.FormID,
	}
	_, err := a.getCollection(db).UpdateOne(c, filter, update)
	if err != nil {
		return err
	}
	return err
}

func (a *filterRepo) SelectByID(c context.Context, db *mongo.Database, permissionGroupID string) (*models.Filter, error) {
	filter := models.Filter{}
	err := a.getCollection(db).FindOne(c, bson.M{"_id": permissionGroupID}).Decode(&filter)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &filter, nil
}

func (a *filterRepo) GetByCondition(c context.Context, db *mongo.Database, perGroupID, formID string) (*models.Filter, error) {
	filter := models.Filter{}
	err := a.getCollection(db).FindOne(c,
		bson.M{
			"per_group_id": perGroupID,
			"form_id":      formID,
		},
	).Decode(&filter)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &filter, nil
}

// Delete delete
func (a *filterRepo) Delete(ctx context.Context, db *mongo.Database, perGroupID, formID string) error {
	filter := bson.M{
		"per_group_id": perGroupID,
		"form_id":      formID,
	}
	_, err := a.getCollection(db).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (a *filterRepo) TableName() string {
	return "filter"
}

func (a *filterRepo) getCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection(a.TableName())
}

// NewFilterRepo 实例化
func NewFilterRepo() models.FilterRepo {
	return &filterRepo{}
}
