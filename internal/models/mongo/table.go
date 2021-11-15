package mongo

import (
	"context"
	"git.internal.yunify.com/qxp/molecule/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type tableRepo struct{}

// NewTableRepo NewTableRepo
func NewTableRepo() models.TableRepo {
	return &tableRepo{}
}

func (t *tableRepo) TableName() string {
	return "table_schema"
}

func (t *tableRepo) Create(ctx context.Context, db *mongo.Database, table *models.Table) error {
	//filter := bson.M{"table_id": table.TableID}
	//r := db.Collection(t.TableName()).FindOne(ctx, filter)
	//an := models.Table{}
	//err := r.Decode(&an)
	//if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
	//	_, err := db.Collection(t.TableName()).InsertOne(ctx, table)
	//	return err
	//}
	//update := bson.D{
	//	{Key: "$set",
	//		Value: bson.D{
	//			{Key: "schema", Value: table.Schema},
	//		},
	//	},
	//}
	//_, err = db.Collection(t.TableName()).UpdateOne(ctx, filter, update)
	//return err
	_, err := db.Collection(t.TableName()).InsertOne(ctx, table)
	return err

}

func (t *tableRepo) GetByID(ctx context.Context, db *mongo.Database, table *models.Table) (*models.Table, error) {
	query := bson.M{"table_id": table.TableID}
	r := db.Collection(t.TableName()).FindOne(ctx, query)
	an := models.Table{}
	err := r.Decode(&an)
	if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &an, err
}

func (t *tableRepo) Update(ctx context.Context, db *mongo.Database, table *models.Table) error {
	query := bson.M{"table_id": table.TableID}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "schema", Value: table.Schema},
			},
		},
	}
	_, err := db.Collection(t.TableName()).UpdateOne(ctx, query, update)
	return err
}

func (t *tableRepo) Delete(ctx context.Context, db *mongo.Database, tableID string) error {
	filter := bson.M{"table_id": tableID}
	// 目前是硬删除，存在删除表单设计schema，就会连同页面配置一起删除。
	// 是否需要考虑两者单独删除的情况?
	_, err := db.Collection(t.TableName()).DeleteOne(ctx, filter)
	return err
}

func (t *tableRepo) UpdateConfig(ctx context.Context, db *mongo.Database, table *models.Table) error {
	query := bson.M{"table_id": table.TableID}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "config", Value: table.Config},
			},
		},
	}
	_, err := db.Collection(t.TableName()).UpdateOne(ctx, query, update)
	return err
}

func (t *tableRepo) DeleteConfig(ctx context.Context, db *mongo.Database, table *models.Table) error {
	query := bson.M{"_id": table.ID}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "config", Value: table.Config},
			},
		},
	}
	_, err := db.Collection(t.TableName()).UpdateOne(ctx, query, update)
	return err
}
