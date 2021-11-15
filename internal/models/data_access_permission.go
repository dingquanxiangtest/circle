package models

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// DataAccessPermission dataAccessPermission
type DataAccessPermission struct {
	ID         string                  `bson:"_id"`
	PerGroupID string                  `bson:"per_group_id"`
	FormID     string                  `bson:"form_id"`
	Conditions map[string]*ConditionVO `bson:"conditions"`
}

// DataAccessPermissionRepo DataAccessPermissionRepo
type DataAccessPermissionRepo interface {
	Create(ctx context.Context, db *mongo.Database, dataAccessPer *DataAccessPermission) error
	// Delete delete by id
	Delete(ctx context.Context, db *mongo.Database, formID, perGroupID string) error

	Get(ctx context.Context, db *mongo.Database, formID, perGroupID string) (*DataAccessPermission, error)

	Update(ctx context.Context, db *mongo.Database, dataAccessPer *DataAccessPermission) error
}
