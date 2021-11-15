package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"git.internal.yunify.com/qxp/molecule/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewGroupFormRepo NewGroupFormRepo
func NewGroupFormRepo() models.GroupFormRepo {
	return &groupForm{}
}

type groupForm struct {
}

func (g *groupForm) getCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection("group_form")
}
func (g *groupForm) Create(ctx context.Context, db *mongo.Database, form *models.GroupForm) error {
	_, err := g.getCollection(db).InsertOne(ctx, form)
	if err != nil {
		return err
	}
	return nil
}

func (g *groupForm) GetByGroupID(ctx context.Context, db *mongo.Database, PerGroupID string) ([]*models.GroupForm, error) {
	findOptions := &options.FindOptions{}
	results := make([]*models.GroupForm, 0)
	filter := bson.M{"per_group_id": PerGroupID}
	cur, err := g.getCollection(db).Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	cur.All(ctx, &results)
	return results, nil
}

func (g *groupForm) DeleteByGroupID(ctx context.Context, db *mongo.Database, PerGroupID string) error {
	// 删除id 为指定的那个
	_, err := g.getCollection(db).DeleteMany(ctx, bson.M{"per_group_id": PerGroupID})
	if err != nil {
		return err
	}
	return nil
}

func (g *groupForm) GetByPerAndForm(ctx context.Context, db *mongo.Database, perGroupID, formID string) (*models.GroupForm, error) {
	resp := &models.GroupForm{}
	findOptions := &options.FindOneOptions{}
	filter := bson.M{
		"form_id":      formID,
		"per_group_id": perGroupID,
	}
	err := g.getCollection(db).FindOne(ctx, filter, findOptions).Decode(resp)
	if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (g *groupForm) DeleteByPerIDAndFormID(ctx context.Context, db *mongo.Database, perGroupID string, formID string) error {
	filter := bson.M{
		"form_id":      formID,
		"per_group_id": perGroupID,
	}
	_, err := g.getCollection(db).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (g *groupForm) FindGroupByFormID(ctx context.Context, db *mongo.Database, formID string) ([]*models.GroupForm, error) {
	result := make([]*models.GroupForm, 0)
	filter := bson.M{
		"form_id": formID,
	}
	cur, err := g.getCollection(db).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cur.All(ctx, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
