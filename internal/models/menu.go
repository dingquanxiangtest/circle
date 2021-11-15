package models

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// Type Type
type Type int

const (
	// None None
	None Type = -1
	// MenuType 0:表示页面
	MenuType Type = 0
	// GroupType 1:表示分组
	GroupType Type = 1
	// CustomType 2:表示自定义页面
	CustomType Type = 2
)

// BindType BindType
type BindType int

const (
	// Unbound 未绑定
	Unbound BindType = (1 + iota) * 10
	// Bound 绑定
	Bound BindType = (1 + iota) * 10
)

// Menu menu
type Menu struct {
	// ID 主键
	ID string `bson:"_id"`

	// AppID 应用id
	AppID string `bson:"app_id"`

	// Name 菜单名称
	Name string `bson:"name"`

	// Icon 图标
	Icon string `bson:"icon"`

	// Sort 排序
	Sort int `bson:"sort"`

	// MenuType 菜单类型：0代表菜单，1代表分组
	MenuType Type `bson:"menu_type"`

	// Describe 描述
	Describe string `bson:"describe"`

	// GroupID 分组id
	GroupID string `bson:"group_id"`

	// BindingState 页面绑定状态
	BindingState BindType `bson:"binding_state"`
}

// MenuRepo menu[存储服务]
type MenuRepo interface {

	// 查找当前应用所在组中是否有这个菜单
	FindSameMenuName(ctx context.Context, db *mongo.Database, menu *Menu) (string, error)

	// 找到当前组中最大的Sort
	FindMaxSortFromGroup(ctx context.Context, db *mongo.Database, menu *Menu) (int, error)

	// 插入菜单
	InsertMenu(ctx context.Context, db *mongo.Database, menu *Menu) error

	// 删除菜单
	DeleteMenuFromGroup(ctx context.Context, db *mongo.Database, menu *Menu) error

	// 更新当前组中的所有sort
	UpdateSortFromGroup(ctx context.Context, db *mongo.Database, menu *Menu) error

	// 更新菜单
	UpdateMenuByID(ctx context.Context, db *mongo.Database, menu *Menu) error

	//从组中查出所有数据
	FindAllFromGroup(ctx context.Context, db *mongo.Database, menu *Menu) ([]*Menu, error)

	// 删除分组
	DeleteGroupByID(ctx context.Context, db *mongo.Database, id string) error

	// 查询所有分组+菜单
	ListAllGroup(ctx context.Context, db *mongo.Database, menu *Menu) ([]*Menu, error)

	// 从范围中查询所有菜单【同一分组中】
	FindAllFromRange(ctx context.Context, db *mongo.Database, fromSort, toSort int, appID, groupID string) ([]*Menu, error)

	// 通过id更新sort
	UpdateSortByID(ctx context.Context, db *mongo.Database, menu *Menu) error

	// 通过id批量修改sort增加或减少number
	BatchUpdateSortByID(ctx context.Context, db *mongo.Database, number int, ids []string) error

	// 从范围中查询所有菜单【不同分组中】
	FindAllBySortAndGroup(ctx context.Context, db *mongo.Database, menu *Menu) ([]*Menu, error)

	// 批量查找菜单
	BatchFindMenus(ctx context.Context, db *mongo.Database, ids []string) ([]*Menu, error)

	// 根据ID修改页面类型
	ModifyMenuTypeByID(ctx context.Context, db *mongo.Database, id string, menuType Type) error

	// 根据ID获取菜单
	FindByID(ctx context.Context, db *mongo.Database, id string) (*Menu, error)

	// 根据ID更新绑定状态
	UpdateBindingStateByID(ctx context.Context, db *mongo.Database, menu *Menu) error
}
