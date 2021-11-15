package models

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

// PageType PageType
type PageType int

// PageStatus PageStatus
type PageStatus int

const (
	// HTMLFile Html文件上传
	HTMLFile Type = 1

	// UsingStatus 应用中状态
	UsingStatus PageStatus = 1

	// NotUsingStatus 未应用状态
	NotUsingStatus PageStatus = 0
)

// CustomPage CustomPage
type CustomPage struct {
	// ID 主键
	ID string `bson:"_id"`
	// FileURL 文件路径
	FileURL string `bson:"file_url"`
	// FileSize 文件大小
	FileSize string `bson:"file_size"`
	// AppID 应用id
	AppID string `bson:"app_id"`
	// CreatedBy 创建人ID
	CreatedBy string `bson:"created_by"`
	// CreatedName 创建人名字
	CreatedName string `bson:"created_name"`
	// UpdatedBy 修改人id
	UpdatedBy string `bson:"updated_by"`
	// UpdatedName 修改人姓名
	UpdatedName string `bson:"updated_name"`
	// CreatedAt 创建时间
	CreatedAt int64 `bson:"created_at"`
	// UpdatedAt 修改时间
	UpdatedAt int64 `bson:"updated_at"`
}

// CustomPageRepo CustomPage[存储服务]
type CustomPageRepo interface {

	// CreateCustomPage 创建自定义页面
	CreateCustomPage(ctx context.Context, db *mongo.Database, custom *CustomPage) error

	// DeleteCustomPage 根据ID删除自定义页面
	DeleteCustomPage(ctx context.Context, db *mongo.Database, id string) error

	// UpdateCustomPage 更新自定义页面
	UpdateCustomPage(ctx context.Context, db *mongo.Database, custom *CustomPage) error

	// FindOne 根据ID获取自定义页面
	FindOne(ctx context.Context, db *mongo.Database, id string) (*CustomPage, error)
}
