package models

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// InitType 默认初始化的权限
	InitType PerType = 1
	// CreateType CreateType
	CreateType PerType = 2
)

// PerType PerType
type PerType int64

// PermissionGroup PermissionGroup
type PermissionGroup struct {
	ID          string      `bson:"_id"`
	AppID       string      `bson:"app_id"`
	Name        string      `bson:"name"`
	CreatedBy   string      `bson:"created_by"`
	Scopes      []*ScopesVO `bson:"scopes"`
	CreatedAt   int64       `bson:"created_at"`
	Description string      `bson:"description"`
	Types       PerType     `bson:"types"`
	Pages       []*string   `bson:"pages"`
}

// PermissionGroupRepo UserGroupRepo
type PermissionGroupRepo interface {
	Create(ctx context.Context, db *mongo.Database, permissionGroup *PermissionGroup) error
	Update(ctx context.Context, db *mongo.Database, permissionGroup *PermissionGroup) error
	Delete(ctx context.Context, db *mongo.Database, id string) error
	GetByIDUserGroup(ctx context.Context, db *mongo.Database, id string) (*PermissionGroup, error)
	GetListUserGroup(ctx context.Context, db *mongo.Database, appID string) ([]*PermissionGroup, error)
	GetByName(ctx context.Context, db *mongo.Database, name, id, appID string) (bool, error)
	GetBYScopeIDs(ctx context.Context, db *mongo.Database, userID, depID, appID string) (*PermissionGroup, error)
	GetByUserInfo(ctx context.Context, db *mongo.Database, userID, depID string, appID string) ([]*PermissionGroup, error)
	VisibilityByAppID(ctx context.Context, db *mongo.Database, appID string) ([]string, error)
	UpdatePagePermission(ctx context.Context, db *mongo.Database, id string, pageIDs []*string) error
	DeletePagePermissionByPageID(ctx context.Context, db *mongo.Database, appID, pageID string) error
	DeletePagePermissionByID(ctx context.Context, db *mongo.Database, groupID, pageID string) error
	AddPagePermission(ctx context.Context, db *mongo.Database, groupID, pageID string) error
	FindGroupByIDs(ctx context.Context, db *mongo.Database, groupIDs []*string) ([]*PermissionGroup, error)
	FindGroupByPageID(ctx context.Context, db *mongo.Database, pageID string) ([]*PermissionGroup, error)
	GetInitGroup(ctx context.Context, db *mongo.Database, AppID string) (*PermissionGroup, error)
}
