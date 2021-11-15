package models

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// OPRead find or findOne
	OPRead = 1 << iota
	// OPCreate create
	OPCreate
	// OPUpdate update
	OPUpdate
	// OPDelete delete
	OPDelete
)

// OperatePermission operatePermission
type OperatePermission struct {
	ID         string `bson:"_id"`
	PerGroupID string `bson:"per_group_id"`
	FormID     string `bson:"form_id"`
	Authority  int64  `bson:"authority"`
}

// OperatePermissionRepo operatePermissionRepo
type OperatePermissionRepo interface {
	Delete(ctx context.Context, db *mongo.Database, formID, perGroupID string) error

	Create(ctx context.Context, db *mongo.Database, operatePer *OperatePermission) error

	Get(ctx context.Context, db *mongo.Database, formID, perGroupID string) (*OperatePermission, error)
	// Update Update
	Update(ctx context.Context, mongodb *mongo.Database, per *OperatePermission) error
}
