package models

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

// GroupForm GroupForm
type GroupForm struct {
	ID         string `bson:"_id"`
	PerGroupID string `bson:"per_group_id"`
	FormID     string `bson:"form_id"`
	FormName   string `bson:"form_name"`
}

// GroupFormRepo GroupFormRepo
type GroupFormRepo interface {
	Create(ctx context.Context, db *mongo.Database, form *GroupForm) error
	GetByGroupID(ctx context.Context, db *mongo.Database, perGroupID string) ([]*GroupForm, error)
	DeleteByGroupID(ctx context.Context, db *mongo.Database, perGroupID string) error
	GetByPerAndForm(ctx context.Context, db *mongo.Database, perGroupID, formID string) (*GroupForm, error)
	DeleteByPerIDAndFormID(ctx context.Context, mongodb *mongo.Database, perGroupID string, formID string) error
	FindGroupByFormID(ctx context.Context, db *mongo.Database, formID string) ([]*GroupForm, error)
}
