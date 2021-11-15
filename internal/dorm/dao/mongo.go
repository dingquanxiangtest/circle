package dao

import (
	"context"
	"strings"

	"git.internal.yunify.com/qxp/molecule/internal/dorm/clause"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MONGO mongo dao
type MONGO struct {
	C *mongo.Collection
}

// FindOne find one entity
func (m *MONGO) FindOne(ctx context.Context, builder clause.Builder) (Data, error) {
	mongoBuilder, err := instantiationMONGO(builder)
	if err != nil {
		return nil, err
	}
	singleResult := m.C.FindOne(ctx, mongoBuilder.Vars)
	var result = make(Data)
	err = singleResult.Decode(&result)
	return result, err
}

// Find find entities
func (m *MONGO) Find(ctx context.Context, builder clause.Builder, findOpt FindOptions) ([]Data, error) {
	mongoBuilder, err := instantiationMONGO(builder)
	if err != nil {
		return nil, err
	}
	if len(mongoBuilder.Agg) != 0 {
		return m.aggregation(ctx, mongoBuilder)
	}
	opt := &options.FindOptions{}
	opt = opt.SetLimit(findOpt.Size)
	opt = opt.SetSkip((findOpt.Page - 1) * findOpt.Size)
	opt = opt.SetSort(mongoSort(findOpt.Sort...))
	cursor, err := m.C.Find(ctx, mongoBuilder.Vars, opt)
	if err != nil {
		return nil, err
	}
	result := make([]Data, 0)
	err = cursor.All(ctx, &result)

	return result, err
}

// Count count entities
func (m *MONGO) Count(ctx context.Context, builder clause.Builder) (int64, error) {
	mongoBuilder, err := instantiationMONGO(builder)
	if err != nil {
		return 0, err
	}
	return m.C.CountDocuments(ctx, mongoBuilder.Vars)
}

// Insert insert entities
func (m *MONGO) Insert(ctx context.Context, entity ...interface{}) error {
	_, err := m.C.InsertMany(ctx, entity)
	return err
}

// Update update entities
func (m *MONGO) Update(ctx context.Context, builder clause.Builder, entity interface{}) (int64, error) {
	mongoBuilder, err := instantiationMONGO(builder)
	if err != nil {
		return 0, err
	}
	updateResult, err := m.C.UpdateMany(ctx, mongoBuilder.Vars, bson.M{"$set": entity})

	return updateResult.ModifiedCount, err
}

// Delete delete entities with condition
func (m *MONGO) Delete(ctx context.Context, builder clause.Builder) (int64, error) {
	mongoBuilder, err := instantiationMONGO(builder)
	if err != nil {
		return 0, err
	}
	deleteResult, err := m.C.DeleteMany(ctx, mongoBuilder.Vars)
	return deleteResult.DeletedCount, err
}

func (m *MONGO) aggregation(ctx context.Context, mongoBuilder *clause.MONGO) ([]Data, error) {
	bsons := make([]bson.M, 0)
	if mongoBuilder.Vars != nil || len(mongoBuilder.Vars) != 0 {
		bsons = append(bsons, bson.M{"$match": mongoBuilder.Vars})
	}
	bsons = append(bsons, mongoBuilder.Agg, bson.M{"$project": bson.M{"_id": 0}})
	cursor, err := m.C.Aggregate(ctx, bsons, nil)
	if err != nil {
		return nil, err
	}
	result := make([]Data, 0)
	err = cursor.All(ctx, &result)
	return result, err
}

func instantiationMONGO(builder clause.Builder) (*clause.MONGO, error) {
	mongo, ok := builder.(*clause.MONGO)
	if !ok {
		return nil, ErrAssertBuilder
	}
	return mongo, nil
}

func mongoSort(array ...string) bson.D {
	sort := make(bson.D, 0, len(array))
	for _, elem := range array {
		if strings.HasPrefix(elem, "-") {
			sort = append(sort, bson.E{Key: elem[1:], Value: -1})
			continue
		}
		sort = append(sort, bson.E{Key: elem, Value: 1})
	}
	return sort
}
