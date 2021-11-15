package mongo

import (
	"context"
	"git.internal.yunify.com/qxp/molecule/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dataBaseSchemaRepo struct {
}

func (d *dataBaseSchemaRepo) Delete(ctx context.Context, db *mongo.Database, tableID string) error {
	filter := bson.M{"table_id": tableID}
	// 目前是硬删除，存在删除表单设计schema，就会连同页面配置一起删除。
	// 是否需要考虑两者单独删除的情况?
	_, err := d.getCollection(db).DeleteOne(ctx, filter)
	return err
}

func (d *dataBaseSchemaRepo) getCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection("database_schema")
}

// NewDataBaseSchemaRepo NewDataBaseSchemaRepo
func NewDataBaseSchemaRepo() models.DataBaseSchemaRepo {
	return &dataBaseSchemaRepo{}
}

func (d *dataBaseSchemaRepo) Create(ctx context.Context, db *mongo.Database, schema *models.DataBaseSchema) error {
	_, err := d.getCollection(db).InsertOne(ctx, schema)
	return err
}

func (d *dataBaseSchemaRepo) GetByTableID(ctx context.Context, db *mongo.Database, tableID string) (*models.DataBaseSchema, error) {
	database := &models.DataBaseSchema{}
	err := d.getCollection(db).FindOne(ctx,
		bson.M{
			"table_id": tableID,
		},
	).Decode(database)
	if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return database, nil
}

func (d *dataBaseSchemaRepo) GetByCondition(ctx context.Context, db *mongo.Database, appID, tableID, title string) (*models.DataBaseSchema, error) {
	database := &models.DataBaseSchema{}
	filter := bson.M{
		"app_id": appID,
	}
	if tableID != "" {
		filter["table_id"] = tableID
	}
	if title != "" {
		filter["title"] = title
	}
	err := d.getCollection(db).FindOne(ctx,
		filter,
	).Decode(database)
	if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return database, nil
}
func (d *dataBaseSchemaRepo) Update(ctx context.Context, db *mongo.Database, schema *models.DataBaseSchema) error {
	filter := bson.M{"table_id": schema.TableID}
	update := bson.M{
		"$set": bson.M{
			"title":       schema.Title,
			"field_len":   schema.FieldLen,
			"description": schema.Description,
			"schema":      schema.Schema,
			"updated_at":  schema.UpdatedAt,
			"editor_id":   schema.EditorID,
			"editor_name": schema.EditorName,
		},
	}
	_, err := d.getCollection(db).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

const (
	defaultLimit = 999
	defaultPage  = 1
)

func (d *dataBaseSchemaRepo) Search(ctx context.Context, db *mongo.Database, appID, title string, source models.SourceType, size int64, page int64) ([]*models.DataBaseSchema, int64, error) {
	opt := &options.FindOptions{
		Sort: bson.M{
			"created_at": 1,
		},
	}
	if size != 0 && page != 0 {
		opt = opt.SetLimit(size)
		opt = opt.SetSkip((page - 1) * size)
	} else {
		opt = opt.SetLimit(defaultLimit)
		opt = opt.SetSkip((defaultPage - 1) * size)
	}
	filter := bson.M{"app_id": appID}
	if title != "" {
		filter["title"] = bson.M{"$regex": title}
	}
	if source != 0 {
		filter["source"] = source
	}
	cursor, err := d.getCollection(db).Find(ctx, filter, opt)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*models.DataBaseSchema, 0)
	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, 0, err
	}
	count, err := d.getCollection(db).CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return result, count, err
}
