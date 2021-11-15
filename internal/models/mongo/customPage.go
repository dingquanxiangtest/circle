package mongo

import (
	"context"
	"git.internal.yunify.com/qxp/molecule/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// DefaultPageSize if PageSize is less than or equal to zero, set pageSize to DefaultPageSize
	DefaultPageSize int64 = 10

	// NonePageSize used to determine empty PageSize
	NonePageSize int64 = 0

	// DefaultCurrentPage  if CurrentPage is less than or equal to zero, set CurrentPage to DefaultCurrentPage
	DefaultCurrentPage int64 = 1

	// NoneCurrentPage used to determine empty CurrentPage
	NoneCurrentPage int64 = 0
)

// NewCustomPageRepo 创建自定义页面存储服务
func NewCustomPageRepo() models.CustomPageRepo {
	return &customPageRepo{}
}

type customPageRepo struct {
}

func (c *customPageRepo) getCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection("custom")
}

func (c *customPageRepo) CreateCustomPage(ctx context.Context, db *mongo.Database, custom *models.CustomPage) error {
	_, err := c.getCollection(db).InsertOne(ctx, custom)
	return err
}

func (c *customPageRepo) DeleteCustomPage(ctx context.Context, db *mongo.Database, id string) error {
	filter := bson.M{
		"_id": id,
	}
	_, err := c.getCollection(db).DeleteOne(ctx, filter)
	return err
}

func (c *customPageRepo) UpdateCustomPage(ctx context.Context, db *mongo.Database, custom *models.CustomPage) error {
	update := bson.M{
		"$set": bson.M{
			"file_url":     custom.FileURL,
			"file_size":    custom.FileSize,
			"updated_at":   custom.UpdatedAt,
			"updated_name": custom.UpdatedName,
			"updated_by":   custom.UpdatedBy,
		},
	}
	_, err := c.getCollection(db).UpdateByID(ctx, custom.ID, update)
	return err
}

func (c *customPageRepo) FindOne(ctx context.Context, db *mongo.Database, id string) (*models.CustomPage, error) {
	filter := bson.M{
		"_id": id,
	}
	result := &models.CustomPage{}
	err := c.getCollection(db).FindOne(ctx, filter).Decode(result)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}
