package models

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

// SubTable table relation info
type SubTable struct {
	// id pk
	ID              string `json:"id" bson:"_id"`
	// app id
	AppID           string `json:"appID" bson:"app_id"`
	// table id
	TableID         string `json:"tableID" bson:"table_id"`
	// table key name
	FieldName        string `json:"fieldName" bson:"field_name"`
	// sub table id
	SubTableID      string `json:"subTableID" bson:"sub_table_id"`
	// table type
	SubTableType    string  `json:"subTableType" bson:"sub_table_type"`
	// filter
	Filter         []string `json:"filter" bson:"filter"`
	AggCondition   *AggregationCondition  `json:"aggCondition" bson:"agg_condition"`
}

// AggregationCondition AggregationCondition
type AggregationCondition struct {
	FieldName     string                  `json:"fieldName"`
	AggType       string                  `json:"aggType"`
	Conditions    map[string]interface{}  `json:"condition"`
}

// SubTableRepo schema
type SubTableRepo interface {
	// create save schema
	Create(ctx context.Context,db *mongo.Database,table *SubTable) error
	// query
	GetByID(ctx context.Context,db *mongo.Database,table *SubTable)([]*SubTable,error)
	// find one
	GetByCondition(ctx context.Context,db *mongo.Database,table *SubTable)(*SubTable,error)
	// update schema
	Update(ctx context.Context,db *mongo.Database,table *SubTable) error
	// delete schema
	Delete(ctx context.Context,db *mongo.Database,table *SubTable) error
	// aggregation record
	GetSubTableByType(ctx context.Context,db *mongo.Database,table *SubTable)([]*SubTable,error)
}
