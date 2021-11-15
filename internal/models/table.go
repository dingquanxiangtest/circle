package models

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

// Table schema info
type Table struct {
	// id pk
	ID string `json:"id" bson:"_id"`
	// table id
	TableID string `json:"tableID" bson:"table_id"`
	// table design json schema
	Schema map[string]interface{} `json:"schema" bson:"schema"`
	// table page config json schema
	Config map[string]interface{} `json:"config" bson:"config"`
}

// TableRepo schema
type TableRepo interface {
	// Create save schema
	Create(ctx context.Context, db *mongo.Database, table *Table) error
	// GetByID GetByID
	GetByID(ctx context.Context, db *mongo.Database, table *Table) (*Table, error)
	// Update schema
	Update(ctx context.Context, db *mongo.Database, table *Table) error
	// UpdateConfig page config
	UpdateConfig(ctx context.Context, db *mongo.Database, table *Table) error
	// Delete schema
	Delete(ctx context.Context, db *mongo.Database, tabID string) error
	// DeleteConfig config
	DeleteConfig(ctx context.Context, db *mongo.Database, table *Table) error
}
