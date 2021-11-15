package mongo

import (
	"context"

	"git.internal.yunify.com/qxp/molecule/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMenuRepo 创建菜单存储服务
func NewMenuRepo() models.MenuRepo {
	return &menuRepo{}
}

type menuRepo struct {
}

func (m *menuRepo) getCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection("menu")
}

func (m *menuRepo) FindSameMenuName(ctx context.Context, db *mongo.Database, menu *models.Menu) (string, error) {
	filter := bson.M{
		"$and": []bson.M{
			{"app_id": menu.AppID},
			{"name": menu.Name},
			{"group_id": menu.GroupID},
		},
	}
	var result models.Menu
	err := m.getCollection(db).FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return result.ID, err
}

func (m *menuRepo) FindMaxSortFromGroup(ctx context.Context, db *mongo.Database, menu *models.Menu) (int, error) {
	var result models.Menu
	filter := bson.M{
		"$and": []bson.M{
			{"app_id": menu.AppID},
			{"group_id": menu.GroupID},
		},
	}
	opts := &options.FindOneOptions{
		Sort: bson.D{
			{
				Key:   "sort",
				Value: -1, // 降序
			},
		},
	}
	err := m.getCollection(db).FindOne(ctx, filter, opts).Decode(&result)
	if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
		return 0, nil
	}
	if err != nil {
		return -1, err
	}

	return result.Sort, err
}

func (m *menuRepo) InsertMenu(ctx context.Context, db *mongo.Database, menu *models.Menu) error {

	_, err := m.getCollection(db).InsertOne(ctx, menu)
	return err

}

func (m *menuRepo) DeleteMenuFromGroup(ctx context.Context, db *mongo.Database, menu *models.Menu) error {

	filter := bson.M{
		"_id": menu.ID,
	}
	_, err := m.getCollection(db).DeleteOne(ctx, filter)
	return err

}

func (m *menuRepo) UpdateSortFromGroup(ctx context.Context, db *mongo.Database, menu *models.Menu) error {

	filter := bson.M{
		"group_id": menu.GroupID,
		"sort": bson.M{
			"$gt": menu.Sort,
		},
		"app_id": menu.AppID,
	}
	update := bson.M{
		"$inc": bson.M{
			"sort": -1,
		},
	}
	_, err := m.getCollection(db).UpdateMany(ctx, filter, update)
	return err
}

func (m *menuRepo) UpdateMenuByID(ctx context.Context, db *mongo.Database, menu *models.Menu) error {

	update := bson.M{
		"$set": bson.M{
			"name":     menu.Name,
			"icon":     menu.Icon,
			"describe": menu.Describe,
		},
	}
	_, err := m.getCollection(db).UpdateByID(ctx, menu.ID, update)
	return err
}

func (m *menuRepo) FindAllFromGroup(ctx context.Context, db *mongo.Database, menu *models.Menu) ([]*models.Menu, error) {
	result := make([]*models.Menu, 0)
	filter := bson.M{
		"$and": []bson.M{
			{"app_id": menu.AppID},
			{"group_id": menu.ID},
		},
	}
	opts := &options.FindOptions{
		Sort: bson.D{
			{
				Key:   "sort",
				Value: 1, // 升序
			},
		},
	}
	cursor, err := m.getCollection(db).Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &result)
	return result, err
}

func (m *menuRepo) DeleteGroupByID(ctx context.Context, db *mongo.Database, id string) error {
	filter := bson.M{
		"_id": id,
	}
	_, err := m.getCollection(db).DeleteOne(ctx, filter)
	return err
}

func (m *menuRepo) ListAllGroup(ctx context.Context, db *mongo.Database, menu *models.Menu) ([]*models.Menu, error) {
	menus := make([]*models.Menu, 0)
	filter := bson.M{
		"app_id": menu.AppID,
	}
	switch menu.MenuType {
	case models.GroupType:
		filter["menu_type"] = menu.MenuType
	case models.MenuType:
		filter["menu_type"] = menu.MenuType
	case models.None:
	}
	opts := &options.FindOptions{
		Sort: bson.M{
			"sort": 1,
		},
	}
	cursor, err := m.getCollection(db).Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &menus)
	return menus, err
}

func (m *menuRepo) FindAllFromRange(ctx context.Context, db *mongo.Database, from, to int, appID, groupID string) ([]*models.Menu, error) {
	var filter bson.M
	element := map[string]interface{}{
		"group_id": groupID,
		"app_id":   appID,
	}

	if from > to {
		element["sort"] = bson.M{"$gte": to, "$lt": from}
	}
	if from < to {
		element["sort"] = bson.M{"$gt": from, "$lte": to}
	}

	filter = bson.M{
		"$and": []bson.M{
			element,
		},
	}
	opts := &options.FindOptions{
		Sort: bson.M{
			"sort": 1, // 升序
		},
	}
	result := make([]*models.Menu, 0)
	cursor, err := m.getCollection(db).Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &result)
	return result, err
}

func (m *menuRepo) UpdateSortByID(ctx context.Context, db *mongo.Database, menu *models.Menu) error {
	update := bson.M{
		"$set": bson.M{
			"sort":     menu.Sort,
			"group_id": menu.GroupID,
		},
	}
	_, err := m.getCollection(db).UpdateByID(ctx, menu.ID, update)
	return err
}

func (m *menuRepo) BatchUpdateSortByID(ctx context.Context, db *mongo.Database, number int, ids []string) error {
	filter := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}
	update := bson.M{
		"$inc": bson.M{
			"sort": number,
		},
	}
	_, err := m.getCollection(db).UpdateMany(ctx, filter, update)
	return err
}

func (m *menuRepo) FindAllBySortAndGroup(ctx context.Context, db *mongo.Database, menu *models.Menu) ([]*models.Menu, error) {
	filter := bson.M{
		"$and": []bson.M{
			{
				"app_id": menu.AppID,
			},
			{
				"group_id": menu.GroupID,
			},
			{
				"sort": bson.M{
					"$gte": menu.Sort,
				},
			},
		},
	}
	opts := &options.FindOptions{
		Sort: bson.M{
			"sort": 1, // 升序
		},
	}
	cursor, err := m.getCollection(db).Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	menus := make([]*models.Menu, 0)

	err = cursor.All(ctx, &menus)
	return menus, err
}

func (m *menuRepo) BatchFindMenus(ctx context.Context, db *mongo.Database, ids []string) ([]*models.Menu, error) {
	filter := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}
	opts := &options.FindOptions{
		Sort: bson.M{
			"sort": 1,
		},
	}
	cursor, err := m.getCollection(db).Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	result := make([]*models.Menu, cursor.RemainingBatchLength())
	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (m *menuRepo) ModifyMenuTypeByID(ctx context.Context, db *mongo.Database, id string, menuType models.Type) error {
	update := bson.M{
		"$set": bson.M{
			"menu_type":     menuType,
			"binding_state": models.Bound,
		},
	}
	_, err := m.getCollection(db).UpdateByID(ctx, id, update)
	return err
}

func (m *menuRepo) FindByID(ctx context.Context, db *mongo.Database, id string) (*models.Menu, error) {
	result := &models.Menu{}
	filer := bson.M{
		"_id": id,
	}
	err := m.getCollection(db).FindOne(ctx, filer).Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *menuRepo) UpdateBindingStateByID(ctx context.Context, db *mongo.Database, menu *models.Menu) error {

	update := bson.M{
		"$set": bson.M{
			"binding_state": menu.BindingState,
		},
	}
	_, err := m.getCollection(db).UpdateByID(ctx, menu.ID, update)
	return err
}
