package models

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

// Filter filter
type Filter struct {
	ID         string `json:"id" bson:"_id"` // //主键
	FormID     string `json:"formID" bson:"form_id"`
	PerGroupID string `json:"per_group_id" bson:"per_group_id"`
	FieldJSON  string `json:"fieldJSON" bson:"field_json"`
	WebSchema  string `json:"webSchema" bson:"web_schema"`
}

// FilterRepo 数据层接口
type FilterRepo interface {
	Insert(c context.Context, mongo *mongo.Database, req *Filter) error
	Update(c context.Context, mongo *mongo.Database, req *Filter) error
	Delete(c context.Context, mongo *mongo.Database, perGroupID, formID string) error
	SelectByID(c context.Context, mongo *mongo.Database, id string) (*Filter, error)
	GetByCondition(c context.Context, mongo *mongo.Database, perGroupID, formID string) (*Filter, error)
}
