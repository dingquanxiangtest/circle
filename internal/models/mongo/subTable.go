package mongo

import (
	"context"
	"git.internal.yunify.com/qxp/molecule/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type subTableRepo struct {}

// NewSubTableRepo NewSubTableRepo
func NewSubTableRepo() models.SubTableRepo {
	return &subTableRepo{}
}

func (s *subTableRepo) TableName() string {
	return "sub_table_relation"
}

func (s *subTableRepo)Create(ctx context.Context,db *mongo.Database,table *models.SubTable)error  {
	query := bson.M{
		"table_id"  :    table.TableID,
		"field_name":    table.FieldName,
	}
	r := db.Collection(s.TableName()).FindOne(ctx,query)
	an := models.SubTable{}
	err := r.Decode(&an)
	if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
		_,err := db.Collection(s.TableName()).InsertOne(ctx,table)
		return err
	}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "filter", Value: table.Filter},
				{Key: "sub_table_id", Value: table.SubTableID},
				{Key: "sub_table_type",Value: table.SubTableType},
				{Key: "agg_condition",Value: table.AggCondition},
			},
		},
	}
	_,err = db.Collection(s.TableName()).UpdateOne(ctx,query,update)
	return err
}

func (s *subTableRepo)GetByID(ctx context.Context,db *mongo.Database,table *models.SubTable)([]*models.SubTable,error) {
	query := bson.M{"table_id":table.TableID}
	result := make([]*models.SubTable, 0)
	cursor,err := db.Collection(s.TableName()).Find(ctx,query)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &result)
	return result,err
}

func (s *subTableRepo)GetByCondition(ctx context.Context,db *mongo.Database,table *models.SubTable)(*models.SubTable,error) {
	query := bson.M{
		"table_id"     : table.TableID,
		"sub_table_id" : table.SubTableID,
		"field_name"   : table.FieldName,
	}
	if table.SubTableID == ""{
		query = bson.M{
			"table_id"   : table.TableID,
			"field_name" : table.FieldName,
		}
	}
	r := db.Collection(s.TableName()).FindOne(ctx,query)
	an := models.SubTable{}
	err := r.Decode(&an)
	if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &an,err
}

func (s *subTableRepo)Update(ctx context.Context,db *mongo.Database,table *models.SubTable) error  {
	query := bson.M{
		"table_id"     : table.TableID,
		"sub_table_id" : table.SubTableID,
	}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{},
		},
	}
	_,err := db.Collection(s.TableName()).UpdateOne(ctx,query,update)
	return err
}

func (s *subTableRepo)Delete(ctx context.Context,db *mongo.Database,table *models.SubTable) error  {
	filter := bson.M{
		"table_id"     : table.TableID,
		"sub_table_id" : table.SubTableID,
	}
	_,err := db.Collection(s.TableName()).DeleteOne(ctx,filter)
	return err
}

func (s *subTableRepo)GetSubTableByType(ctx context.Context,db *mongo.Database,table *models.SubTable)([]*models.SubTable,error) {
	result := make([]*models.SubTable, 0)
	query := bson.M{
		"table_id"   : table.TableID,
		"sub_table_type" : table.SubTableType,
	}
	cursor, err := db.Collection(s.TableName()).Find(ctx,query)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &result)
	return result, err
}
