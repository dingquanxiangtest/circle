package mongo

import (
	"context"

	"git.internal.yunify.com/qxp/molecule/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type datasetRepo struct {
}

// NewDataSetRepo NewDataSetRepo
func NewDataSetRepo() models.DataSetRepo {
	return &datasetRepo{}
}

func (d *datasetRepo) getCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection("dataset")

}

// Insert Insert
func (d *datasetRepo) Insert(ctx context.Context, db *mongo.Database, dataset *models.DataSet) error {
	_, err := d.getCollection(db).InsertOne(ctx, dataset)
	if err != nil {
		return err
	}
	return nil
}

// Update Update
func (d *datasetRepo) Update(ctx context.Context, db *mongo.Database, dataset *models.DataSet) error {
	update := bson.M{
		"$set": bson.M{
			"name":    dataset.Name,
			"tag":     dataset.Tag,
			"type":    dataset.Type,
			"content": dataset.Content,
		},
	}
	_, err := d.getCollection(db).UpdateByID(ctx, dataset.ID, update)
	if err != nil {
		return err
	}
	return err
}

// Delete Delete
func (d *datasetRepo) Delete(ctx context.Context, db *mongo.Database, id string) error {
	// 删除id 为指定的那个
	_, err := d.getCollection(db).DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

// GetByID GetByID
func (d *datasetRepo) GetByID(ctx context.Context, db *mongo.Database, id string) (*models.DataSet, error) {
	dataset := &models.DataSet{}

	err := d.getCollection(db).FindOne(ctx, bson.M{"_id": id}).Decode(dataset)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return dataset, nil
}

// GetByCondition GetByCondition
func (d *datasetRepo) GetByCondition(ctx context.Context, db *mongo.Database, tag, name string, types int64) ([]*models.DataSet, error) {
	findOptions := &options.FindOptions{
		Sort: bson.M{
			"created_at": 1,
		},
	}
	results := make([]*models.DataSet, 0)
	var filter = bson.M{}
	if types != 0 {
		filter["type"] = types
	}
	if tag != "" {
		filter["tag"] = tag
	}
	if name != "" {
		filter["name"] = bson.M{"$regex": name}
	}
	cur, err := d.getCollection(db).Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	cur.All(ctx, &results)
	return results, nil
}

func (d *datasetRepo) GetByName(ctx context.Context, db *mongo.Database, name, id string) (bool, error) {
	filter := bson.M{
		"name": name,
		"_id":  bson.M{"$ne": id},
	}
	count, err := d.getCollection(db).CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
