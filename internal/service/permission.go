package service

import (
	"context"
	"time"

	"git.internal.yunify.com/qxp/misc/time2"
	"git.internal.yunify.com/qxp/molecule/internal/filters"
	"git.internal.yunify.com/qxp/molecule/pkg/client"

	"git.internal.yunify.com/qxp/misc/error2"
	"git.internal.yunify.com/qxp/misc/id2"
	logger "git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/misc/redis2"
	"git.internal.yunify.com/qxp/molecule/internal/models"
	repo "git.internal.yunify.com/qxp/molecule/internal/models/mongo"
	"git.internal.yunify.com/qxp/molecule/internal/models/redis"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/code"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"

	redisc "github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

// Permission permission
type Permission interface {
	CreatePermissionGroup(ctx context.Context, req *CreatePermissionGroupReq) (*CreatePermissionGroupResp, error)

	UpdatePermissionGroup(ctx context.Context, req *UpdatePermissionGroupReq) (*UpdatePermissionGroupResp, error)

	DeletePermissionGroup(ctx context.Context, req *DeletePermissionGroupReq) (*DeletePermissionGroupResp, error)

	GetByIDPermissionGroup(ctx context.Context, req *GetByIDPermissionGroupReq) (*GetByIDPermissionGroupResp, error)

	GetListPermissionGroup(ctx context.Context, req *GetListPermissionGroupReq) (*GetListPermissionGroupResp, error)

	GetByConditionPerGroup(ctx context.Context, req *GetByConditionPerGroupReq) (*GetByConditionPerGroupResp, error)

	GetFormsPerGroup(ctx context.Context, req *GetFormsPerGroupReq) (*GetFormsPerGroupResp, error)

	GetDataAccessPermission(ctx context.Context, req *GetDataAccessPermissionReq) (*GetDataAccessPermissionResp, error)

	GetOperatePermission(ctx context.Context, req *GetOperatePermissionReq) (*GetOperatePermissionResp, error)

	GetGroupPerByUserInfo(ctx context.Context, req *GetGroupPerByUserInfoReq) (*GetGroupPerByUserInfoResp, error)

	GetPerSelect(ctx context.Context, req *GetPerSelectReq) (*GetPerSelectResp, error)

	VisibilityApp(ctx context.Context, req *VisibilityAppReq) (*VisibilityAppResp, error)

	SaveForm(ctx context.Context, req *SaveFormReq) (*SaveFormResp, error)

	DeleteForm(ctx context.Context, req *DeleteFormReq) (*DeleteFormResp, error)

	GetForm(ctx context.Context, req *GetFormReq) (*GetFormResp, error)

	GetOperate(ctx context.Context, req *GetOperateReq) (*GetOperateResp, error)

	SaveUserPerMatch(c context.Context, req *SaveUserPerMatchReq) (*SaveUserPerMatchResp, error)

	UpdatePerName(c context.Context, req *UpdatePerNameReq) (*UpdatePerNameResp, error)

	UpdatePagePermission(ctx context.Context, req *UpdatePagePermissionReq) (*UpdatePagePermissionResp, error)

	GetGroupPage(ctx context.Context, req *GetGroupPageReq) (*GetGroupPageResp, error)

	ModifyPagePer(ctx context.Context, req *ModifyPagePerReq) (*ModifyPagePerResp, error)

	GetPerGroupByMenu(ctx context.Context, req *GetPerGroupByMenuReq) (*GetPerGroupByMenuResp, error)
}

const (
	lockPermission   = "lockPermission"
	lockPerMatch     = "lockPerMatch"
	lockTimeout      = time.Duration(30) * time.Second                   // 30秒
	perTime          = time.Hour * time.Duration(12) * time.Duration(30) // 30天
	timeSleep        = time.Millisecond * 500                            // 0.5 秒
	notAuthority     = -1
	removePermission = 0 //页面没有权限
	addPermission    = 1 //页面有权限
)

type permission struct {
	mongodb                  *mongo.Database
	redisClient              *redisc.ClusterClient
	permissionGroupRepo      models.PermissionGroupRepo
	dataAccessPermissionRepo models.DataAccessPermissionRepo
	operatePermissionRepo    models.OperatePermissionRepo
	permissionRepo           models.PermissionRepo
	groupFormRepo            models.GroupFormRepo
	customPage               CustomPage
	filter                   Filter
	appClient                client.AppCenter
}

// SaveUserPerMatchReq SaveUserPerMatchReq
type SaveUserPerMatchReq struct {
	UserID     string `json:"userID"`
	AppID      string `json:"appID"`
	PerGroupID string `json:"perGroupID"`
}

// SaveUserPerMatchResp SaveUserPerMatchResp
type SaveUserPerMatchResp struct {
}

// SaveUserPerMatch SaveUserPerMatch
func (per *permission) SaveUserPerMatch(ctx context.Context, req *SaveUserPerMatchReq) (*SaveUserPerMatchResp, error) {
	matchPer := &models.PermissionMatch{
		UserID:     req.UserID,
		PerGroupID: req.PerGroupID,
		AppID:      req.AppID,
	}
	err := per.permissionRepo.CreatePerMatch(ctx, matchPer)
	if err != nil {
		return nil, err
	}
	return &SaveUserPerMatchResp{}, nil
}

// GetOperateReq GetOperateReq
type GetOperateReq struct {
	FormID string `json:"formID"`
	UserID string `json:"userID"`
	DepID  string `json:"depID"`
	AppID  string `json:"appID"`
}

// GetOperateResp GetOperateResp
type GetOperateResp struct {
	Authority int64 `json:"authority"`
}

func (per *permission) GetOperate(ctx context.Context, req *GetOperateReq) (*GetOperateResp, error) {
	match, err := per.getPerMatch(ctx, req.UserID, req.DepID, req.AppID)
	if err != nil {
		return nil, err
	}
	if match == nil {
		return nil, nil
	}
	operate, err := per.operatePermissionRepo.Get(ctx, per.mongodb, req.FormID, match.PerGroupID)
	if err != nil {
		return nil, err
	}
	if operate != nil {
		return &GetOperateResp{
			Authority: operate.Authority,
		}, nil
	}

	return &GetOperateResp{
		Authority: models.OPUpdate | models.OPCreate | models.OPDelete | models.OPRead,
	}, nil

}

// GetFormReq GetFormReq
type GetFormReq struct {
	PerGroupID string `json:"perGroupID"`
}

// GetFormResp GetFormResp
type GetFormResp struct {
	FormArr []*FormVo `json:"formArr"`
}

// FormVo FormVo
type FormVo struct {
	ID        string `json:"id"`
	Authority int64  `json:"authority"`
}

func (per *permission) GetForm(ctx context.Context, req *GetFormReq) (*GetFormResp, error) {
	arr, err := per.groupFormRepo.GetByGroupID(ctx, per.mongodb, req.PerGroupID)
	if err != nil {
		return nil, err
	}
	resp := &GetFormResp{
		FormArr: make([]*FormVo, len(arr)),
	}
	for index, value := range arr {
		// 根据formID ，PerGroupID ，得到权限
		vo := &FormVo{
			ID: value.FormID,
		}
		operate, err := per.operatePermissionRepo.Get(ctx, per.mongodb, value.FormID, value.PerGroupID)
		if err != nil {
			return nil, err
		}
		if operate != nil {
			vo.Authority = operate.Authority
		}
		resp.FormArr[index] = vo
	}
	return resp, nil

}

// GetPerSelectReq GetPerSelectReq
type GetPerSelectReq struct {
	UserID string `json:"userID"`
	DepID  string `json:"depID"`
	AppID  string `json:"appID"`
}

// GetPerSelectResp GetPerSelectResp
type GetPerSelectResp struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (per *permission) GetPerSelect(c context.Context, req *GetPerSelectReq) (*GetPerSelectResp, error) {
	match, err := per.getPerMatch(c, req.UserID, req.DepID, req.AppID)
	if err != nil {
		return nil, err
	}
	if match == nil {
		return nil, nil
	}
	perGroup, err := per.permissionGroupRepo.GetByIDUserGroup(c, per.mongodb, match.PerGroupID)
	if err != nil {
		return nil, err
	}
	if perGroup != nil {
		return &GetPerSelectResp{
			ID:   match.PerGroupID,
			Name: perGroup.Name,
		}, nil
	}
	return nil, nil

}

// DeleteFormReq DeleteFormReq
type DeleteFormReq struct {
	PerGroupID string `json:"perGroupID"`
	FormID     string `json:"formID"`
}

// DeleteFormResp DeleteFormResp
type DeleteFormResp struct {
}

func (per *permission) DeleteForm(ctx context.Context, req *DeleteFormReq) (*DeleteFormResp, error) {
	// 删除映射关系
	err := per.groupFormRepo.DeleteByPerIDAndFormID(ctx, per.mongodb, req.PerGroupID, req.FormID)
	if err != nil {
		return nil, err
	}
	err = per.operatePermissionRepo.Delete(ctx, per.mongodb, req.FormID, req.PerGroupID)
	if err != nil {
		return nil, err
	}
	err = per.dataAccessPermissionRepo.Delete(ctx, per.mongodb, req.FormID, req.PerGroupID)
	if err != nil {
		return nil, err
	}
	dRep := &DELFilterReq{
		PerGroupID: req.PerGroupID,
		FormID:     req.FormID,
	}
	err = per.filter.Delete(ctx, dRep)
	if err != nil {
		return nil, err
	}
	return nil, nil

}

// NewPermission NewPermission
func NewPermission(conf *config.Config, opts ...Options) (Permission, error) {
	redisClient, err := redis2.NewClient(conf.Redis)
	if err != nil {
		return nil, err
	}
	filter, err := NewFilter(conf, opts...)
	if err != nil {
		return nil, err
	}
	customPage, err := NewCustomPage(conf, opts...)
	if err != nil {
		return nil, err
	}
	u := &permission{
		permissionGroupRepo:      repo.NewPermissionGroupRepo(),
		dataAccessPermissionRepo: repo.NewDataAccessPermissionRepo(),
		operatePermissionRepo:    repo.NewOperatePermissionRepo(),
		groupFormRepo:            repo.NewGroupFormRepo(),
		customPage:               customPage,
		filter:                   filter,
		permissionRepo:           redis.NewPermissionRepo(redisClient),
		appClient:                client.NewAppCenter(conf.InternalNet),
	}
	for _, opt := range opts {
		opt(u)
	}
	return u, nil
}

func (per *permission) SetMongo(client *mongo.Client, dbName string) {
	per.mongodb = client.Database(dbName)
}

// CreatePermissionGroupReq req
type CreatePermissionGroupReq struct {
	AppID       string         `json:"appID"`
	Name        string         `json:"name"`
	CreatedBy   string         `json:"createdBy"`
	Description string         `json:"description"`
	Types       models.PerType `json:"types"`
}

// CreatePermissionGroupResp resp
type CreatePermissionGroupResp struct {
	ID string `json:"id"`
}

// CreatePermissionGroup CreatePermissionGroup
func (per *permission) CreatePermissionGroup(ctx context.Context, req *CreatePermissionGroupReq) (*CreatePermissionGroupResp, error) {
	// app_id name and id get group per
	exist, err := per.permissionGroupRepo.GetByName(ctx, per.mongodb, req.Name, "", req.AppID)
	if err != nil {
		return nil, err
	}
	if exist {
		// 返回用户组名已经存在
		return nil, error2.NewError(code.ErrExistGroupNameState)
	}
	perGroup := &models.PermissionGroup{
		ID:          id2.GenID(),
		Name:        req.Name,
		CreatedBy:   req.CreatedBy,
		Description: req.Description,
		AppID:       req.AppID,
		CreatedAt:   time2.NowUnix(),
		Pages:       make([]*string, 0),
	}
	perGroup.Types = req.Types
	if req.Types == 0 {
		perGroup.Types = models.CreateType
	}

	err = per.permissionGroupRepo.Create(ctx, per.mongodb, perGroup)
	if err != nil {
		return nil, err
	}

	return &CreatePermissionGroupResp{ID: perGroup.ID}, nil
}

// SaveFormReq SaveFormReq
type SaveFormReq struct {
	FormID     string                         `json:"formID"`
	FormName   string                         `json:"formName"`
	PerGroupID string                         `json:"perGroupID"`
	Authority  int64                          `json:"authority"`
	Conditions map[string]*models.ConditionVO `json:"conditions"`
	Schema     filters.Schema                 `json:"schema"`
}

// SaveFormResp SaveFormResp
type SaveFormResp struct {
}

func (per *permission) SaveForm(ctx context.Context, req *SaveFormReq) (*SaveFormResp, error) {
	// 1、 查询映射关系
	perForm, err := per.groupFormRepo.GetByPerAndForm(ctx, per.mongodb, req.PerGroupID, req.FormID)
	if err != nil {
		return nil, err
	}
	if perForm == nil { // 创建映射关系
		perGroupForm := &models.GroupForm{
			ID:         id2.GenID(),
			FormID:     req.FormID,
			PerGroupID: req.PerGroupID,
			FormName:   req.FormName,
		}
		err := per.groupFormRepo.Create(ctx, per.mongodb, perGroupForm)
		if err != nil {
			return nil, err
		}
	}
	// 2、保存权限
	err = per.saveOperate(ctx, req)
	if err != nil {
		return nil, err
	}
	// 3、保存数据过滤权限
	err = per.saveFilter(ctx, req)
	if err != nil {
		return nil, err
	}
	// 4. 保存数据访问权限
	err = per.saveCondition(ctx, req)
	if err != nil {
		return nil, err
	}
	// 5、删除缓存
	permission := &models.Permission{
		PerGroupID: req.PerGroupID,
		FormID:     req.FormID,
	}
	err = per.permissionRepo.Delete(ctx, permission)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.GenRequestID(ctx))
	}
	return &SaveFormResp{}, nil

}

func (per *permission) saveOperate(ctx context.Context, req *SaveFormReq) error {
	// 查询数据库
	operate, err := per.operatePermissionRepo.Get(ctx, per.mongodb, req.FormID, req.PerGroupID)
	if err != nil {
		return err
	}
	if operate == nil { // 创建
		operate = &models.OperatePermission{
			ID:         id2.GenID(),
			FormID:     req.FormID,
			PerGroupID: req.PerGroupID,
			Authority:  req.Authority,
		}
		err = per.operatePermissionRepo.Create(ctx, per.mongodb, operate)
		if err != nil {
			return err
		}
		return nil
	}
	operate.Authority = req.Authority
	err = per.operatePermissionRepo.Update(ctx, per.mongodb, operate)
	if err != nil {
		return err
	}
	return nil
}

func (per *permission) saveFilter(ctx context.Context, req *SaveFormReq) error {
	SaveFilterReq := &SaveFilterReq{
		PerGroupID: req.PerGroupID,
		FormID:     req.FormID,
		Schema:     req.Schema,
	}
	_, err := per.filter.SaveJSONFilter(ctx, SaveFilterReq)
	if err != nil {
		return err
	}
	return nil
}
func (per *permission) saveCondition(ctx context.Context, req *SaveFormReq) error {
	dataAccess, err := per.dataAccessPermissionRepo.Get(ctx, per.mongodb, req.FormID, req.PerGroupID)
	if err != nil {
		return err
	}
	if dataAccess == nil {
		dataAccess = &models.DataAccessPermission{
			FormID:     req.FormID,
			PerGroupID: req.PerGroupID,
			ID:         id2.GenID(),
			Conditions: req.Conditions,
		}
		err = per.dataAccessPermissionRepo.Create(ctx, per.mongodb, dataAccess)
		if err != nil {
			return err
		}
		return nil
	}
	dataAccess.Conditions = req.Conditions
	err = per.dataAccessPermissionRepo.Update(ctx, per.mongodb, dataAccess)
	if err != nil {
		return err
	}
	return nil
}

// UpdatePermissionGroupReq UpdatePermissionGroupReq
type UpdatePermissionGroupReq struct {
	ID     string             `json:"id"`
	Scopes []*models.ScopesVO `json:"scopes"`
}

// UpdatePermissionGroupResp UpdatePermissionGroupResp
type UpdatePermissionGroupResp struct {
}

// UpdatePermissionGroup updatePermissionGroup 修改 用户组的成员 或者修改 用户组的权限，
func (per *permission) UpdatePermissionGroup(ctx context.Context, req *UpdatePermissionGroupReq) (*UpdatePermissionGroupResp, error) {
	perGroup, err := per.permissionGroupRepo.GetByIDUserGroup(ctx, per.mongodb, req.ID)
	if err != nil {
		return nil, err
	}
	updateUser := &models.PermissionGroup{
		ID:     req.ID,
		Scopes: req.Scopes,
	}
	err = per.permissionGroupRepo.Update(ctx, per.mongodb, updateUser)
	if err != nil {
		return nil, err
	}
	err = per.permissionRepo.DeletePerMatch(ctx, perGroup.AppID)
	if err != nil {
		logger.Logger.Errorw("delete redis error ", logger.STDRequestID(ctx), err.Error())
	}
	return &UpdatePermissionGroupResp{}, nil
}

// UpdatePerNameReq UpdatePerNameReq
type UpdatePerNameReq struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdatePerNameResp UpdatePerNameResp
type UpdatePerNameResp struct {
}

// UpdatePerName UpdatePerName
func (per *permission) UpdatePerName(ctx context.Context, req *UpdatePerNameReq) (*UpdatePerNameResp, error) {
	// 判断是否更新
	perGroup, err := per.permissionGroupRepo.GetByIDUserGroup(ctx, per.mongodb, req.ID)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		exist, err := per.permissionGroupRepo.GetByName(ctx, per.mongodb, req.Name, req.ID, perGroup.AppID)
		if err != nil {
			return nil, err
		}
		if exist {
			// 返回用户组名已经存在
			return nil, error2.NewError(code.ErrExistGroupNameState)
		}
	}
	updateUser := &models.PermissionGroup{
		Name:        req.Name,
		Description: req.Description,
		ID:          req.ID,
	}
	err = per.permissionGroupRepo.Update(ctx, per.mongodb, updateUser)
	if err != nil {
		return nil, err
	}
	return &UpdatePerNameResp{}, nil
}

// DeletePermissionGroupReq req
type DeletePermissionGroupReq struct {
	ID string `json:"id"`
}

// DeletePermissionGroupResp resp
type DeletePermissionGroupResp struct {
}

// DeletePermissionGroup DeletePermissionGroup 删除
func (per *permission) DeletePermissionGroup(ctx context.Context, req *DeletePermissionGroupReq) (*DeletePermissionGroupResp, error) {
	perGroup, err := per.permissionGroupRepo.GetByIDUserGroup(ctx, per.mongodb, req.ID)
	if err != nil {
		return nil, err
	}
	forms, err := per.groupFormRepo.GetByGroupID(ctx, per.mongodb, req.ID)
	if err != nil {
		return nil, err
	}
	err = per.permissionGroupRepo.Delete(ctx, per.mongodb, req.ID)
	if err != nil {
		return nil, err
	}
	//delete GroupForm
	err = per.groupFormRepo.DeleteByGroupID(ctx, per.mongodb, req.ID)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return nil, err
	}
	err = per.permissionRepo.DeletePerMatch(ctx, perGroup.AppID)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
	}
	for _, value := range forms {
		err = per.permissionRepo.Delete(ctx, &models.Permission{
			PerGroupID: req.ID,
			FormID:     value.FormID,
		})
		if err != nil {
			logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		}
	}
	return &DeletePermissionGroupResp{}, nil

}

// GetByIDPermissionGroupReq req
type GetByIDPermissionGroupReq struct {
	ID string `json:"id"`
}

// GetByIDPermissionGroupResp GetByIDPermissionGroupResp
type GetByIDPermissionGroupResp struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	AppID       string             `json:"appID"`
	Scopes      []*models.ScopesVO `json:"scopes"`
	Description string             `json:"description"`
}

// GetByIDPermissionGroup GetByIDPermissionGroup
func (per *permission) GetByIDPermissionGroup(ctx context.Context, req *GetByIDPermissionGroupReq) (*GetByIDPermissionGroupResp, error) {
	perGroup, err := per.permissionGroupRepo.GetByIDUserGroup(ctx, per.mongodb, req.ID)
	if err != nil {
		return nil, err
	}
	resp := &GetByIDPermissionGroupResp{
		ID:          perGroup.ID,
		Name:        perGroup.Name,
		Scopes:      perGroup.Scopes,
		AppID:       perGroup.AppID,
		Description: perGroup.Description,
	}
	return resp, nil
}

// GetListPermissionGroupReq req
type GetListPermissionGroupReq struct {
	AppID string `json:"appID"`
}

// GetListPermissionGroupResp resp
type GetListPermissionGroupResp struct {
	ListVO []*listVO `json:"list"`
}
type listVO struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	CreatedBy   string             `json:"createdBy"`
	Scopes      []*models.ScopesVO `json:"scopes"`
	Description string             `json:"description"`
	AppID       string             `json:"appID"`
	Add         bool               `json:"add"`
	Types       models.PerType     `json:"types"`
}

// GetListPermissionGroup GetListPermissionGroup
func (per *permission) GetListPermissionGroup(ctx context.Context, req *GetListPermissionGroupReq) (*GetListPermissionGroupResp, error) {
	result, err := per.permissionGroupRepo.GetListUserGroup(ctx, per.mongodb, req.AppID)
	if err != nil {
		return nil, err
	}
	resp := &GetListPermissionGroupResp{
		ListVO: make([]*listVO, len(result)),
	}
	for i, value := range result {
		ids, err := per.groupFormRepo.GetByGroupID(ctx, per.mongodb, value.ID)
		if err != nil {
			return nil, err
		}
		resp.ListVO[i] = new(listVO)
		if value.Types == models.InitType {
			resp.ListVO[i].Add = true
		} else if value.Pages != nil && len(value.Pages) > 0 {
			resp.ListVO[i].Add = true
		} else {
			resp.ListVO[i].Add = len(ids) > 0
		}

		clone(resp.ListVO[i], value)
	}
	return resp, nil
}
func clone(dst *listVO, src *models.PermissionGroup) {
	dst.ID = src.ID
	dst.Name = src.Name
	dst.Scopes = src.Scopes
	dst.Description = src.Description
	dst.AppID = src.AppID
	dst.Types = src.Types
}

// GetByConditionPerGroupReq req
type GetByConditionPerGroupReq struct {
	UserID string `json:"userID"`
	DepID  string `json:"depID"`
	FormID string `json:"formID"`
	AppID  string `json:"appId"`
}

// GetByConditionPerGroupResp resp
type GetByConditionPerGroupResp struct {
	ID            string
	AppID         string
	FormID        string
	Name          string
	CreatedBy     string
	Scopes        []*models.ScopesVO
	Sequence      int64
	Description   string
	Authority     int64
	DataAccessPer map[string]*models.ConditionVO
}

// GetByConditionPerGroup 根据条件获取 权限用户组
func (per *permission) GetByConditionPerGroup(ctx context.Context, req *GetByConditionPerGroupReq) (*GetByConditionPerGroupResp, error) {
	// 1、 根据用户id 、  应用id ，得到 ====》 对应的权限组ID，
	perMatch, err := per.getPerMatch(ctx, req.UserID, req.DepID, req.AppID)
	if err != nil {
		return nil, err
	}
	if perMatch == nil {
		return nil, nil
	}
	// 2、 根据用户权限组id ，得到 权限信息
	permission, err := per.getPerInfo(ctx, perMatch.PerGroupID, req.FormID)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		return nil, nil
	}
	resp := &GetByConditionPerGroupResp{}
	per.cloneField(permission, resp)
	return resp, nil

}
func (per *permission) cloneField(permission *models.Permission, resp *GetByConditionPerGroupResp) {
	resp.ID = permission.PerGroupID
	resp.DataAccessPer = permission.Conditions
	resp.Authority = permission.Authority
	resp.FormID = permission.FormID
}
func (per *permission) getPerInfo(ctx context.Context, perGroupID, formID string) (*models.Permission, error) {
	for i := 0; i < 5; i++ {
		// 1. 去redis 查询
		permission, err := per.permissionRepo.Get(ctx, perGroupID, formID)
		if err != nil {
			logger.Logger.Errorw(err.Error(), perGroupID, logger.STDRequestID(ctx))
			return nil, err
		}
		if permission != nil {
			return permission, nil
		}
		lock, err := per.permissionRepo.Lock(ctx, lockPermission, 1, lockTimeout) // 抢占分布式锁
		if err != nil {
			logger.Logger.Errorw(err.Error(), perGroupID, logger.STDRequestID(ctx))
			return nil, err
		}

		if !lock {
			<-time.After(timeSleep)
			continue
		}
		break
	}
	defer per.permissionRepo.UnLock(ctx, lockPermission) // 删除锁
	perGroup, err := per.permissionGroupRepo.GetByIDUserGroup(ctx, per.mongodb, perGroupID)
	if err != nil {
		logger.Logger.Errorw(err.Error(), perGroupID, logger.STDRequestID(ctx))
		return nil, err
	}
	if perGroup == nil { // 数据库为空，直接返回，
		return nil, nil
	}
	permission := &models.Permission{
		PerGroupID: perGroup.ID,
		AppID:      perGroup.AppID,
		FormID:     formID,
		Name:       perGroup.Name,
	}
	// 查询
	operatePer, err := per.operatePermissionRepo.Get(ctx, per.mongodb, formID, perGroupID)
	if err != nil {
		logger.Logger.Errorw(err.Error(), perGroupID, logger.STDRequestID(ctx))
		return nil, err
	}
	if operatePer != nil {
		permission.Authority = operatePer.Authority
	} else {
		permission.Authority = notAuthority
	}
	dataAccessPer, err := per.dataAccessPermissionRepo.Get(ctx, per.mongodb, formID, perGroupID)
	if err != nil {
		logger.Logger.Errorw(err.Error(), perGroupID, logger.STDRequestID(ctx))
		return nil, err
	}
	if dataAccessPer != nil {
		permission.Conditions = dataAccessPer.Conditions
	}
	// 加入到redis 缓存
	per.permissionRepo.Create(ctx, permission, perTime)
	return permission, nil
}

// getPerMatch 查找匹配的关系
func (per *permission) getPerMatch(ctx context.Context, userID, depID, appID string) (*models.PermissionMatch, error) {
	for i := 0; i < 5; i++ {
		perMatch, err := per.permissionRepo.GetPerMatch(ctx, userID, appID)
		if err != nil {
			logger.Logger.Errorw(err.Error(), userID, logger.STDRequestID(ctx))
			return nil, err
		}
		if perMatch != nil {
			return perMatch, nil
		}
		lock, err := per.permissionRepo.Lock(ctx, lockPerMatch, 1, lockTimeout) // 抢占分布式锁
		if err != nil {
			logger.Logger.Errorw(err.Error(), userID, logger.STDRequestID(ctx))
			return nil, err
		}
		if !lock {
			<-time.After(timeSleep)
			continue
		}
		break
	}
	defer per.permissionRepo.UnLock(ctx, lockPerMatch) // 删除锁
	perGroup, err := per.permissionGroupRepo.GetBYScopeIDs(ctx, per.mongodb, userID, depID, appID)
	if err != nil {
		logger.Logger.Errorw(err.Error(), depID, logger.STDRequestID(ctx))
		return nil, err
	}
	if perGroup == nil { // 数据库有
		return nil, nil
	}
	perMatch := &models.PermissionMatch{
		UserID:     userID,
		AppID:      appID,
		PerGroupID: perGroup.ID,
	}
	err = per.permissionRepo.CreatePerMatch(ctx, perMatch)
	if err != nil {
		logger.Logger.Errorw(err.Error(), userID, logger.STDRequestID(ctx))
	}
	return perMatch, nil
}

// GetFormsPerGroupReq GetFormsPerGroupReq
type GetFormsPerGroupReq struct {
	UserID string `json:"userID"`
	AppID  string `json:"appID"`
	DepID  string `json:"depID"`
}

// GetFormsPerGroupResp GetFormsPerGroupResp
type GetFormsPerGroupResp struct {
	FormID  []string       `json:"formID"`
	PerType models.PerType `json:"perType"`
}

// GetFormsPerGroup 根据应用id ，用户id ，部门id ，返回表单数组
func (per *permission) GetFormsPerGroup(ctx context.Context, req *GetFormsPerGroupReq) (*GetFormsPerGroupResp, error) {
	// 1. 根据 获取 权限组
	match, err := per.getPerMatch(ctx, req.UserID, req.DepID, req.AppID)
	if err != nil {
		return nil, err
	}
	if match == nil {
		return nil, nil
	}
	// 2. 查询Per Form
	perGroup, err := per.permissionGroupRepo.GetByIDUserGroup(ctx, per.mongodb, match.PerGroupID)
	if err != nil {
		return nil, err
	}
	groupForms, err := per.groupFormRepo.GetByGroupID(ctx, per.mongodb, match.PerGroupID)
	if err != nil {
		return nil, err
	}

	resp := &GetFormsPerGroupResp{
		FormID:  make([]string, len(groupForms)),
		PerType: perGroup.Types,
	}
	for index, value := range groupForms {
		resp.FormID[index] = value.FormID
	}
	if perGroup.Pages != nil {
		for _, m := range perGroup.Pages {
			resp.FormID = append(resp.FormID, *m)
		}
	}
	return resp, nil
}

// GetGroupPerByUserInfoReq GetGroupPerByUserInfoReq
type GetGroupPerByUserInfoReq struct {
	UserID string `json:"userID"`
	DepID  string `json:"depID"`
	AppID  string `json:"appID"`
}

// GetGroupPerByUserInfoResp perGroupArr
type GetGroupPerByUserInfoResp struct {
	PerGroupArr []*PerGroupVO `json:"perGroupArr"`
}

// PerGroupVO PerGroupVO
type PerGroupVO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GetGroupPerByUserInfo GetGroupPerByUserInfo
func (per *permission) GetGroupPerByUserInfo(ctx context.Context, req *GetGroupPerByUserInfoReq) (*GetGroupPerByUserInfoResp, error) {
	// get GetByUserInfo
	perGroupArr, err := per.permissionGroupRepo.GetByUserInfo(ctx, per.mongodb, req.UserID, req.DepID, req.AppID)
	if err != nil {
		return nil, err
	}
	resp := &GetGroupPerByUserInfoResp{
		PerGroupArr: make([]*PerGroupVO, len(perGroupArr)),
	}
	for index, v := range perGroupArr {
		vo := &PerGroupVO{
			ID:   v.ID,
			Name: v.Name,
		}
		resp.PerGroupArr[index] = vo
	}
	return resp, nil
}

// GetDataAccessPermissionReq req
type GetDataAccessPermissionReq struct {
	FormID     string `json:"formID"`
	PerGroupID string `json:"perGroupID"`
}

// GetDataAccessPermissionResp resp
type GetDataAccessPermissionResp struct {
	FormID     string                         `json:"formID"`
	PerGroupID string                         `json:"perGroupID"`
	Conditions map[string]*models.ConditionVO `json:"conditions"`
}

// GetDataAccessPermission GetDataAccessPermission
func (per *permission) GetDataAccessPermission(ctx context.Context, req *GetDataAccessPermissionReq) (*GetDataAccessPermissionResp, error) {
	dataAccessPer, err := per.dataAccessPermissionRepo.Get(ctx, per.mongodb, req.FormID, req.PerGroupID)
	if err != nil {
		return nil, err
	}
	if dataAccessPer == nil {
		return nil, nil
	}
	resp := &GetDataAccessPermissionResp{
		FormID:     dataAccessPer.FormID,
		PerGroupID: dataAccessPer.PerGroupID,
		Conditions: dataAccessPer.Conditions,
	}
	return resp, nil
}

// GetOperatePermissionReq GetOperatePermissionReq
type GetOperatePermissionReq struct {
	PerGroupID string `json:"perGroupID"`
	FormID     string `json:"formID"`
}

// GetOperatePermissionResp resp
type GetOperatePermissionResp struct {
	FormID     string `json:"formID"`
	PerGroupID string `json:"perGroupID"`
	// Authority authority 对应的权限
	Authority int64 `json:"authority"`
}

// GetOperatePermission get
func (per *permission) GetOperatePermission(ctx context.Context, req *GetOperatePermissionReq) (*GetOperatePermissionResp, error) {
	operatePer, err := per.operatePermissionRepo.Get(ctx, per.mongodb, req.FormID, req.PerGroupID)
	if err != nil {
		return nil, err
	}
	if operatePer == nil {
		return nil, nil
	}
	resp := &GetOperatePermissionResp{
		PerGroupID: operatePer.PerGroupID,
		FormID:     operatePer.FormID,
		Authority:  operatePer.Authority,
	}
	return resp, nil
}

// VisibilityAppReq VisibilityAppReq
type VisibilityAppReq struct {
	AppID string `json:"appID"`
}

// VisibilityAppResp VisibilityAppResp
type VisibilityAppResp struct {
}

// VisibilityApp VisibilityApp  // 根据部门或者用户， 判断是否在某个权限组内
func (per *permission) VisibilityApp(ctx context.Context, req *VisibilityAppReq) (*VisibilityAppResp, error) {
	scopeVo, err := per.permissionGroupRepo.VisibilityByAppID(ctx, per.mongodb, req.AppID)
	logger.Logger.Info("scopes is", scopeVo, logger.STDRequestID(ctx))
	if err != nil {
		return nil, err
	}
	_, err = per.appClient.AddAppScope(ctx, req.AppID, scopeVo)
	if err != nil {
		return nil, err
	}
	return &VisibilityAppResp{}, nil
}

// UpdatePagePermissionReq update custom page permission request params
type UpdatePagePermissionReq struct {
	GroupID string    `json:"groupId" binding:"required"`
	Pages   []*string `json:"pageIds"`
}

// UpdatePagePermissionResp update custom page permission return params
type UpdatePagePermissionResp struct {
}

/*
UpdatePagePermission :
	EXPLAINS: Update custom page permission. All selected custom page IDs are required
*/
func (per *permission) UpdatePagePermission(ctx context.Context, req *UpdatePagePermissionReq) (*UpdatePagePermissionResp, error) {
	err := per.permissionGroupRepo.UpdatePagePermission(ctx, per.mongodb, req.GroupID, req.Pages)
	if err != nil {
		return nil, err
	}
	return &UpdatePagePermissionResp{}, nil
}

// GetGroupPageReq get the list of page id that has permission in this group request params
type GetGroupPageReq struct {
	ID string `json:"groupId" binding:"required"`
}

// GetGroupPageResp get the list of page id that has permission in this group return params
type GetGroupPageResp struct {
	Pages []*string `json:"pages"`
}

/*
GetGroupPage :
	EXPLAINS: get the pages IDs by group ID.
*/
func (per *permission) GetGroupPage(ctx context.Context, req *GetGroupPageReq) (*GetGroupPageResp, error) {
	group, err := per.permissionGroupRepo.GetByIDUserGroup(ctx, per.mongodb, req.ID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, error2.NewError(code.ErrNODataSetNameState)
	}
	if group.Pages == nil {
		return &GetGroupPageResp{
			Pages: make([]*string, 0),
		}, nil
	}
	return &GetGroupPageResp{
		Pages: group.Pages,
	}, nil
}

// ModifyPagePerReq modify the page's permission request params
type ModifyPagePerReq struct {
	GroupID string `json:"groupId" binding:"required"`
	PageID  string `json:"pageId" binding:"required"`
	Status  int    `json:"status"`
}

// ModifyPagePerResp modify the page's permission return params
type ModifyPagePerResp struct {
}

/*
ModifyPagePer :
	EXPLAINS: modify the page's permission in the permission group
	STEPS:
		1.judge the request params status
		2.insert or remove page id to mongo
*/
func (per *permission) ModifyPagePer(ctx context.Context, req *ModifyPagePerReq) (*ModifyPagePerResp, error) {
	// TODO
	perGroup, err := per.permissionGroupRepo.GetByIDUserGroup(ctx, per.mongodb, req.GroupID)
	if err != nil {
		return nil, err
	}
	if perGroup.Pages == nil {
		switch req.Status {
		case addPermission:
			pages := make([]*string, 0)
			err := per.permissionGroupRepo.UpdatePagePermission(ctx, per.mongodb, req.GroupID, pages)
			if err != nil {
				return nil, err
			}
			return &ModifyPagePerResp{}, nil
		}
	}
	switch req.Status {
	case removePermission:
		err := per.permissionGroupRepo.DeletePagePermissionByID(ctx, per.mongodb, req.GroupID, req.PageID)
		if err != nil {
			return nil, err
		}
	case addPermission:
		err := per.permissionGroupRepo.AddPagePermission(ctx, per.mongodb, req.GroupID, req.PageID)
		if err != nil {
			return nil, err
		}
	}
	return &ModifyPagePerResp{}, nil
}

// GetPerGroupByMenuReq GetPerGroupByMenuReq
type GetPerGroupByMenuReq struct {
	MenuID   string `json:"menuId" binding:"required"`
	MenuType models.Type
	AppID    string
}

// GetPerGroupByMenuResp GetPerGroupByMenuResp
type GetPerGroupByMenuResp struct {
	PerList []*PerGroupVO `json:"perGroups"`
}

/*
GetPerGroupByMenu :
	EXPLAINS: get the permission group by menu id.
	STEPS:
		1.judge the menu type.
		2.if menu type is form page, query the permission group ids by form id.
			Then, find the permission group information by group ids.
		3.if menu type is custom page, query the permission group by page id.
			It is find the page id in the column of pages.
*/
func (per *permission) GetPerGroupByMenu(ctx context.Context, req *GetPerGroupByMenuReq) (*GetPerGroupByMenuResp, error) {
	resp := &GetPerGroupByMenuResp{
		PerList: make([]*PerGroupVO, 0),
	}
	initGroup, err := per.permissionGroupRepo.GetInitGroup(ctx, per.mongodb, req.AppID)
	if err != nil {
		return nil, err
	}
	resp.PerList = append(resp.PerList, &PerGroupVO{
		ID:   initGroup.ID,
		Name: initGroup.Name,
	})
	switch req.MenuType {
	case models.MenuType:
		groupForms, err := per.groupFormRepo.FindGroupByFormID(ctx, per.mongodb, req.MenuID)
		if err != nil {
			return nil, err
		}
		ids := make([]*string, 0)
		for _, groupForm := range groupForms {
			ids = append(ids, &groupForm.PerGroupID)
		}
		groups, err := per.permissionGroupRepo.FindGroupByIDs(ctx, per.mongodb, ids)
		if err != nil {
			return nil, err
		}
		for _, g := range groups {
			resp.PerList = append(resp.PerList, &PerGroupVO{
				ID:   g.ID,
				Name: g.Name,
			})
		}
	case models.CustomType:
		groups, err := per.permissionGroupRepo.FindGroupByPageID(ctx, per.mongodb, req.MenuID)
		if err != nil {
			return nil, err
		}

		for _, g := range groups {
			resp.PerList = append(resp.PerList, &PerGroupVO{
				ID:   g.ID,
				Name: g.Name,
			})
		}
	}
	return resp, nil
}
