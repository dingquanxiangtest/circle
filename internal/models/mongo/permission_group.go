package mongo

import (
	"context"

	"git.internal.yunify.com/qxp/molecule/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type permissionGroupRepo struct {
}

// NewPermissionGroupRepo new method repo
func NewPermissionGroupRepo() models.PermissionGroupRepo {
	return &permissionGroupRepo{}
}
func (a *permissionGroupRepo) getCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection("permission_group")

}

// Create create
func (a *permissionGroupRepo) Create(ctx context.Context, db *mongo.Database, permissionGroup *models.PermissionGroup) error {
	_, err := a.getCollection(db).InsertOne(ctx, permissionGroup)
	if err != nil {
		return err
	}
	return nil
}

// Update update
func (a *permissionGroupRepo) Update(ctx context.Context, db *mongo.Database, permissionGroup *models.PermissionGroup) error {
	filter := bson.M{"_id": permissionGroup.ID}
	var update primitive.M
	if permissionGroup.Scopes != nil {
		update = bson.M{
			"$set": bson.M{"scopes": permissionGroup.Scopes},
		}
	}
	if permissionGroup.Name != "" || permissionGroup.Description != "" {
		update = bson.M{
			"$set": bson.M{"name": permissionGroup.Name, "description": permissionGroup.Description},
		}
	}
	_, err := a.getCollection(db).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// Delete delete
func (a *permissionGroupRepo) Delete(ctx context.Context, db *mongo.Database, id string) error {
	// 删除id 为指定的那个
	_, err := a.getCollection(db).DeleteOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}})
	if err != nil {
		return err
	}
	return nil
}
func (a *permissionGroupRepo) GetByIDUserGroup(ctx context.Context, db *mongo.Database, id string) (*models.PermissionGroup, error) {
	// 创建一个UserGroup变量用来接收查询的结果
	resp := &models.PermissionGroup{}
	filter := bson.M{"_id": id}

	err := a.getCollection(db).FindOne(ctx, filter).Decode(resp)
	if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetListUserGroup GetListUserGroup
func (a *permissionGroupRepo) GetListUserGroup(ctx context.Context, db *mongo.Database, appID string) ([]*models.PermissionGroup, error) {
	findOptions := &options.FindOptions{
		Sort: bson.M{
			"created_at": 1,
		},
	}
	results := make([]*models.PermissionGroup, 0)
	filter := bson.M{"app_id": appID}
	cur, err := a.getCollection(db).Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	cur.All(ctx, &results)
	return results, nil
}

// GetByName GetByName
func (a *permissionGroupRepo) GetByName(ctx context.Context, db *mongo.Database, name, id, appID string) (bool, error) {
	filter := bson.M{
		"name":   name,
		"_id":    bson.M{"$ne": id},
		"app_id": appID,
	}
	count, err := a.getCollection(db).CountDocuments(ctx, filter)

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (a *permissionGroupRepo) GetBYScopeIDs(ctx context.Context, db *mongo.Database, userID, depID, appID string) (*models.PermissionGroup, error) {
	resp := &models.PermissionGroup{}

	findOptions := &options.FindOneOptions{
		Sort: bson.M{
			"created_at": 1,
		},
	}
	idArr := make([]string, 0)
	if userID != "" {
		idArr = append(idArr, userID)
	}
	if depID != "" {
		idArr = append(idArr, depID)
	}
	filter := bson.M{
		"app_id":    appID,
		"scopes.id": bson.M{"$in": idArr},
	}
	err := a.getCollection(db).FindOne(ctx, filter, findOptions).Decode(resp)
	if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *permissionGroupRepo) GetByUserInfo(ctx context.Context, db *mongo.Database, userID, depID string, appID string) ([]*models.PermissionGroup, error) {
	findOptions := &options.FindOptions{
		Sort: bson.M{
			"created_at": 1,
		},
	}
	results := make([]*models.PermissionGroup, 0)
	idArr := make([]string, 0)
	if userID != "" {
		idArr = append(idArr, userID)
	}
	if depID != "" {
		idArr = append(idArr, depID)
	}
	filter := bson.M{
		"app_id":    appID,
		"scopes.id": bson.M{"$in": idArr},
	}
	cur, err := a.getCollection(db).Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	cur.All(ctx, &results)
	return results, nil
}

func (a *permissionGroupRepo) VisibilityByAppID(ctx context.Context, mongodb *mongo.Database, appID string) ([]string, error) {
	filter := bson.M{
		"app_id": appID,
	}
	vo, err := a.getCollection(mongodb).Distinct(ctx, "scopes.id", filter)
	if err != nil {
		return nil, err
	}
	resp := make([]string, 0)
	for _, value := range vo {
		v, ok := value.(string)
		if ok {
			resp = append(resp, v)
		}
	}
	return resp, nil
}

func (a *permissionGroupRepo) UpdatePagePermission(ctx context.Context, db *mongo.Database, id string, pageIDs []*string) error {
	update := bson.M{
		"$set": bson.M{
			"pages": pageIDs,
		},
	}
	_, err := a.getCollection(db).UpdateByID(ctx, id, update)
	if err != nil {
		return err
	}
	return nil
}

func (a *permissionGroupRepo) DeletePagePermissionByPageID(ctx context.Context, db *mongo.Database, appID, pageID string) error {
	filter := bson.M{
		"app_id": appID,
	}
	update := bson.M{
		"$pull": bson.M{
			"pages": pageID,
		},
	}

	_, err := a.getCollection(db).UpdateMany(ctx, filter, update)
	return err
}

func (a *permissionGroupRepo) DeletePagePermissionByID(ctx context.Context, db *mongo.Database, groupID, pageID string) error {
	update := bson.M{
		"$pull": bson.M{
			"pages": pageID,
		},
	}
	_, err := a.getCollection(db).UpdateByID(ctx, groupID, update)
	return err
}

func (a *permissionGroupRepo) AddPagePermission(ctx context.Context, db *mongo.Database, groupID, pageID string) error {
	update := bson.M{
		"$addToSet": bson.M{
			"pages": pageID,
		},
	}
	_, err := a.getCollection(db).UpdateByID(ctx, groupID, update)
	return err
}

func (a *permissionGroupRepo) FindGroupByIDs(ctx context.Context, db *mongo.Database, groupIDs []*string) ([]*models.PermissionGroup, error) {
	result := make([]*models.PermissionGroup, 0)
	filter := bson.M{
		"_id": bson.M{
			"$in": groupIDs,
		},
	}
	cur, err := a.getCollection(db).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cur.All(ctx, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *permissionGroupRepo) FindGroupByPageID(ctx context.Context, db *mongo.Database, pageID string) ([]*models.PermissionGroup, error) {
	result := make([]*models.PermissionGroup, 0)
	filter := bson.M{
		"pages": pageID,
	}
	cur, err := a.getCollection(db).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cur.All(ctx, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *permissionGroupRepo) GetInitGroup(ctx context.Context, db *mongo.Database, AppID string) (*models.PermissionGroup, error) {
	resp := &models.PermissionGroup{}
	filter := bson.M{
		"app_id": AppID,
		"types":  models.InitType,
	}
	err := a.getCollection(db).FindOne(ctx, filter).Decode(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
