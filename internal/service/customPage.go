package service

import (
	"context"
	"git.internal.yunify.com/qxp/misc/error2"
	"git.internal.yunify.com/qxp/misc/time2"
	"git.internal.yunify.com/qxp/molecule/internal/models"
	repo "git.internal.yunify.com/qxp/molecule/internal/models/mongo"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/code"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"
	"go.mongodb.org/mongo-driver/mongo"
)

// CustomPage CustomPage
type CustomPage interface {

	// CreateCustom Create CustomPage
	CreateCustom(ctx context.Context, req *CreateCustomPageReq) (*CreateCustomPageResp, error)

	// UpdateCustomPage Update customPage information
	UpdateCustomPage(ctx context.Context, req *UpdateCustomPageReq) (*UpdateCustomPageResp, error)

	// DeletePageMenuByMenuID Removes the association between the custom page and the menu
	DeletePageMenuByMenuID(ctx context.Context, req *DeletePageMenuByMenuIDReq) (*DeletePageMenuByMenuIDResp, error)

	// GetByMenuID Get the custom page information by menu id
	GetByMenuID(ctx context.Context, req *GetByMenuIDReq) (*GetByMenuIDResp, error)
}

type customPage struct {
	mongodb             *mongo.Database
	customPageRepo      models.CustomPageRepo
	permissionGroupRepo models.PermissionGroupRepo
}

// NewCustomPage new a customPage service
func NewCustomPage(conf *config.Config, opts ...Options) (CustomPage, error) {
	customPageRepo := repo.NewCustomPageRepo()
	groupRepo := repo.NewPermissionGroupRepo()
	c := &customPage{
		customPageRepo:      customPageRepo,
		permissionGroupRepo: groupRepo,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

// SetMongo SetMongo
func (cus *customPage) SetMongo(client *mongo.Client, dbName string) {
	cus.mongodb = client.Database(dbName)
}

// DeletePageMenuByMenuIDReq delete relation customPage and menu by menu ID request params
type DeletePageMenuByMenuIDReq struct {
	MenuID string `json:"menuId"`
	AppID  string
}

// DeletePageMenuByMenuIDResp delete relation customPage and menu by menu ID return params
type DeletePageMenuByMenuIDResp struct {
}

/*
DeletePageMenuByMenuID :
	EXPLAINS: Removes the association between the custom page and the menu
	STEPS:
		1. get the association between the custom page and the menu by menuID
		2. delete the association.
		3. remove the menu permission
		4. get the number of menus associated with the current custom page
		5. determine whether the current page is associated with other menus.
			if there is no, modify current custom page status to "NotUsingStatus"
*/
func (cus *customPage) DeletePageMenuByMenuID(ctx context.Context, req *DeletePageMenuByMenuIDReq) (*DeletePageMenuByMenuIDResp, error) {

	// remove the custom page.
	err := cus.customPageRepo.DeleteCustomPage(ctx, cus.mongodb, req.MenuID)
	if err != nil {
		return nil, err
	}
	// remove the menu permission.
	err = cus.permissionGroupRepo.DeletePagePermissionByPageID(ctx, cus.mongodb, req.AppID, req.MenuID)
	if err != nil {
		return nil, err
	}
	return &DeletePageMenuByMenuIDResp{}, nil
}

// CreateCustomPageReq Insert customPage request params
type CreateCustomPageReq struct {
	FileURL  string `json:"fileUrl" binding:"required"`
	FileSize string `json:"fileSize" binding:"required"`
	MenuID   string `json:"menuId" binding:"required"`
	AppID    string `json:"appID"`
	UserID   string `json:"userID"`
	UserName string `json:"userName"`
}

// CreateCustomPageResp Insert customPage return params
type CreateCustomPageResp struct {
	ID          string `json:"id"`
	FileURL     string `json:"fileUrl"`
	FileSize    string `json:"fileSize"`
	CreatedName string `json:"createdBy"`
	UpdatedName string `json:"updatedBy"`
	CreatedAt   int64  `json:"createdAt"`
	UpdateAt    int64  `json:"updatedAt"`
}

/*
CreateCustom :
	EXPLAINS: Create a custom page and save to mongo
*/
func (cus *customPage) CreateCustom(ctx context.Context, req *CreateCustomPageReq) (*CreateCustomPageResp, error) {
	custom := &models.CustomPage{
		ID:          req.MenuID,
		FileURL:     req.FileURL,
		FileSize:    req.FileSize,
		AppID:       req.AppID,
		CreatedBy:   req.UserID,
		CreatedName: req.UserName,
		UpdatedBy:   req.UserID,
		UpdatedName: req.UserName,
		CreatedAt:   time2.NowUnix(),
		UpdatedAt:   time2.NowUnix(),
	}
	err := cus.customPageRepo.CreateCustomPage(ctx, cus.mongodb, custom)
	if err != nil {
		return nil, err
	}
	return &CreateCustomPageResp{
		ID:          custom.ID,
		FileURL:     custom.FileURL,
		FileSize:    custom.FileSize,
		CreatedAt:   custom.CreatedAt,
		CreatedName: custom.CreatedName,
		UpdatedName: custom.UpdatedName,
		UpdateAt:    custom.UpdatedAt,
	}, nil
}

// UpdateCustomPageReq Update customPage information request params
type UpdateCustomPageReq struct {
	ID       string `json:"id" binding:"required"`
	FileURL  string `json:"fileUrl" binding:"required"`
	FileSize string `json:"fileSize" binding:"required"`
	UserID   string `json:"userID"`
	UserName string `json:"userName"`
}

// UpdateCustomPageResp Update customPage information return params
type UpdateCustomPageResp struct {
	ID          string `json:"id"`
	FileURL     string `json:"fileUrl"`
	FileSize    string `json:"fileSize"`
	CreatedName string `json:"createdBy"`
	UpdatedName string `json:"updatedBy"`
	CreatedAt   int64  `json:"createdAt"`
	UpdateAt    int64  `json:"updatedAt"`
}

/*
UpdateCustomPage :
	EXPLAINS: Create the information of custom page and save to mongo
*/
func (cus *customPage) UpdateCustomPage(ctx context.Context, req *UpdateCustomPageReq) (*UpdateCustomPageResp, error) {
	custom := &models.CustomPage{
		ID:          req.ID,
		FileURL:     req.FileURL,
		FileSize:    req.FileSize,
		UpdatedBy:   req.UserID,
		UpdatedName: req.UserName,
		UpdatedAt:   time2.NowUnix(),
	}
	err := cus.customPageRepo.UpdateCustomPage(ctx, cus.mongodb, custom)
	if err != nil {
		return nil, err
	}
	page, err := cus.customPageRepo.FindOne(ctx, cus.mongodb, custom.ID)
	if err != nil {
		return nil, err
	}
	if page == nil {
		return nil, error2.NewError(code.ErrNODataSetNameState)
	}
	return &UpdateCustomPageResp{
		ID:          page.ID,
		FileURL:     page.FileURL,
		FileSize:    page.FileSize,
		CreatedAt:   page.CreatedAt,
		CreatedName: page.CreatedName,
		UpdatedName: page.UpdatedName,
		UpdateAt:    page.UpdatedAt,
	}, nil
}

// GetByMenuIDReq get Custom page information by menu id request params.
type GetByMenuIDReq struct {
	MenuID string `json:"menuID"`
}

// GetByMenuIDResp get Custom page information by menu id return params.
type GetByMenuIDResp struct {
	ID          string `json:"id"`
	FileURL     string `json:"fileUrl"`
	FileSize    string `json:"fileSize"`
	CreatedName string `json:"createdBy"`
	UpdatedName string `json:"updatedBy"`
	CreatedAt   int64  `json:"createdAt"`
	UpdateAt    int64  `json:"updatedAt"`
}

/*
GetByMenuID
	EXPLAINS:Get the custom page information by menu id
*/
func (cus *customPage) GetByMenuID(ctx context.Context, req *GetByMenuIDReq) (*GetByMenuIDResp, error) {
	page, err := cus.customPageRepo.FindOne(ctx, cus.mongodb, req.MenuID)
	if err != nil {
		return nil, err
	}
	if page == nil {
		return nil, error2.NewError(code.ErrNODataSetNameState)
	}
	return &GetByMenuIDResp{
		ID:          page.ID,
		FileURL:     page.FileURL,
		FileSize:    page.FileSize,
		CreatedAt:   page.CreatedAt,
		CreatedName: page.CreatedName,
		UpdatedName: page.UpdatedName,
		UpdateAt:    page.UpdatedAt,
	}, nil
}
