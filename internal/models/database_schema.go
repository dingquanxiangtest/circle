package models

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

// ModelSource modelSource
const (
	FormSource  SourceType = 1
	ModelSource SourceType = 2
)

// SourceType SourceType
type SourceType int64

// DataBaseSchema DataBaseSchema
type DataBaseSchema struct {
	ID          string                 `bson:"_id"`
	Title       string                 `bson:"title"`
	AppID       string                 `bson:"app_id"`
	TableID     string                 `bson:"table_id"`
	FieldLen    int64                  `bson:"field_len"`
	Description string                 `bson:"description"`
	Source      SourceType             `bson:"source"`
	CreatedAt   int64                  `bson:"created_at"`
	UpdatedAt   int64                  `bson:"updated_at"`
	CreatorID   string                 `bson:"creator_id"`
	CreatorName string                 `bson:"creator_name"`
	EditorID    string                 `bson:"editor_id"`
	EditorName  string                 `bson:"editor_name"`
	Schema      map[string]interface{} `bson:"schema"` // 过滤之后的代码
}

// DataBaseSchemaRepo DataBaseSchemaRepo
type DataBaseSchemaRepo interface {
	Create(ctx context.Context, db *mongo.Database, schema *DataBaseSchema) error
	GetByTableID(ctx context.Context, db *mongo.Database, tableID string) (*DataBaseSchema, error)
	Update(ctx context.Context, db *mongo.Database, schema *DataBaseSchema) error
	// Search Search
	Search(ctx context.Context, db *mongo.Database, appID, title string, source SourceType, size int64, page int64) ([]*DataBaseSchema, int64, error)

	Delete(ctx context.Context, db *mongo.Database, tableID string) error
	GetByCondition(ctx context.Context, db *mongo.Database, appID, tableID, title string) (*DataBaseSchema, error)
}
