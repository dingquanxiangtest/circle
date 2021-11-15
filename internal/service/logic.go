package service

import (
	"context"
	"encoding/json"

	"git.internal.yunify.com/qxp/misc/error2"
	"git.internal.yunify.com/qxp/misc/id2"
	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/misc/time2"
	"git.internal.yunify.com/qxp/molecule/internal/models"
	repo "git.internal.yunify.com/qxp/molecule/internal/models/mongo"
	"git.internal.yunify.com/qxp/molecule/internal/service/swagger"
	"git.internal.yunify.com/qxp/molecule/pkg/client"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/code"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"

	"go.mongodb.org/mongo-driver/mongo"
)

// Kernel 业务逻辑
type Kernel interface {
	// CreateSchema CreateSchema
	CreateSchema(ctx context.Context, req *CreateTableReq, options ...CreateSchemaOption) (*CreateTableResp, error)
	GetSchemaByTableID(ctx context.Context, req *GetTableReq) (*models.Table, error)
	DeleteSchema(ctx context.Context, req *DeleteTableReq) (*DeleteTableResp, error)
	CreateBlankSchema(ctx context.Context, req *CreateBlankTableReq) (*CreateBlankTableResp, error)
	SearchSchema(ctx context.Context, req *SearchSchemaReq) (*SearchSchemaResp, error)
	// CreateConfig page config
	CreateConfig(ctx context.Context, req *CreateConfigReq) (*CreateConfigResp, error)
	DeleteConfig(ctx context.Context, req *DeleteConfigReq) (*DeleteConfigResp, error)
	//GetXName GetXName
	GetXName(ctx context.Context, req *GetXNameReq) (*GetXNameResp, error)
	CheckRepeat(ctx context.Context, req *CheckRepeatReq) (*CheckRepeatResp, error)
	GetModelDataByMenu(ctx context.Context, req *GetModelDataByMenuReq) (*GetModelDataByMenuResp, error)
}

// NewKernel new kernel
func NewKernel(conf *config.Config, opts ...Options) (Kernel, error) {
	k := &kernel{
		tableRepo:    repo.NewTableRepo(),
		menuRepo:     repo.NewMenuRepo(),
		databaseRepo: repo.NewDataBaseSchemaRepo(),
		polyAPI:      client.NewPolyAPI(conf.InternalNet),
		sw:           swagger.NewSW(),
	}
	for _, opt := range opts {
		opt(k)
	}
	return k, nil
}

type kernel struct {
	tableRepo    models.TableRepo
	menuRepo     models.MenuRepo
	databaseRepo models.DataBaseSchemaRepo
	db           *mongo.Database
	polyAPI      client.PolyAPI
	sw           *swagger.Swagger
}

func (k *kernel) SetMongo(client *mongo.Client, dbName string) {
	k.db = client.Database(dbName)
}

// GetTableReq GetTableReq
type GetTableReq struct {
	TableID string `json:"tableID"`
}

// DeleteTableReq DeleteTableReq
type DeleteTableReq struct {
	TableID string `json:"tableID"`
}

// CreateTableResp CreateTableResp
type CreateTableResp interface{}

// GetTableResp GetTableResp
type GetTableResp interface{}

// DeleteTableResp DeleteTableResp
type DeleteTableResp interface{}

// CreateConfigReq CreateConfigReq
type CreateConfigReq struct {
	TableID string                 `json:"tableID"`
	Config  map[string]interface{} `json:"config"`
}

// UpdateConfigReq UpdateConfigReq
type UpdateConfigReq struct {
	ID      string                 `json:"id"`
	TableID string                 `json:"tableID"`
	Config  map[string]interface{} `json:"config"`
}

// DeleteConfigReq DeleteConfigReq
type DeleteConfigReq struct {
	ID string `json:"id"`
}

// CreateConfigResp CreateConfigResp
type CreateConfigResp interface{}

// GetConfigResp GetConfigResp
type GetConfigResp interface{}

// UpdateConfigResp UpdateConfigResp
type UpdateConfigResp interface{}

// DeleteConfigResp DeleteConfigResp
type DeleteConfigResp interface{}

// CreateBlankTableReq CreateBlankTableReq
type CreateBlankTableReq struct {
}

// CreateBlankTableResp CreateBlankTableResp
type CreateBlankTableResp struct {
	TableID string `json:"tableID"`
}

//CreateSchemaOption CreateSchemaOption
type CreateSchemaOption func(ctx context.Context, req *OptionReq) (interface{}, error)

// GenSwag 转换
func GenSwag(k Kernel) CreateSchemaOption {
	return func(ctx context.Context, req *OptionReq) (interface{}, error) {
		if k2, ok := k.(*kernel); ok {
			genSwagger, err := swagger.GenSwagger(req.Schema, req.AppID, req.TableID, k2.sw)
			if err != nil {
				return nil, err
			}
			swagger, err := json.Marshal(genSwagger)
			if err != nil {
				return nil, err
			}
			content := "form"
			if req.Source == models.ModelSource {
				content = "custom"
			}
			_, err = k2.polyAPI.RegSwagger(ctx, "structor", string(swagger), req.AppID, content)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	}
}

// OptionReq OptionReq
type OptionReq struct {
	AppID       string
	TableID     string
	Total       int64
	Schema      map[string]interface{}
	Title       string
	Source      models.SourceType
	UserID      string
	UserName    string
	Description string
	ISCreate    bool
}

// ConvertSchema ConvertSchema
func ConvertSchema(k Kernel) CreateSchemaOption {
	return func(ctx context.Context, req *OptionReq) (interface{}, error) {
		if k2, ok := k.(*kernel); ok {
			one, err := k2.databaseRepo.GetByTableID(ctx, k2.db, req.TableID)
			if err != nil {
				return nil, err
			}
			table := &models.DataBaseSchema{
				Title:       req.Title,
				TableID:     req.TableID,
				Schema:      req.Schema,
				FieldLen:    req.Total,
				UpdatedAt:   time2.NowUnix(),
				Description: req.Description,
			}
			if one == nil { // 新增
				req.ISCreate = true
				table.ID = id2.GenID()
				table.Source = req.Source
				table.CreatedAt = time2.NowUnix()
				table.AppID = req.AppID
				table.CreatorName = req.UserName
				table.CreatorID = req.UserID
				err := k2.databaseRepo.Create(ctx, k2.db, table)
				if err != nil {
					return nil, err
				}
				return nil, nil
			}
			// 修改
			table.TableID = req.TableID
			table.EditorID = req.UserID
			table.EditorName = req.UserName
			err = k2.databaseRepo.Update(ctx, k2.db, table)
			if err != nil {
				return nil, err
			}
		} else {
			error2.NewError(code.ErrConvert)
		}
		return nil, nil

	}
}

// CreateTableReq CreateTableReq
type CreateTableReq struct {
	AppID    string
	TableID  string                 `json:"tableID"`
	Schema   map[string]interface{} `json:"schema"`
	UserID   string                 `json:"user_id"`
	UserName string                 `json:"user_name"`
	Source   models.SourceType      `json:"source"` // source 1 是表单驱动，2是模型驱动
}

// CreateSchema 相当于create 和  update
func (k *kernel) CreateSchema(ctx context.Context, req *CreateTableReq, options ...CreateSchemaOption) (resp *CreateTableResp, err error) {
	defer func() {
		if err == nil {
			k.after(ctx, req, options...)
			if err != nil {
				logger.Logger.Error(logger.STDRequestID(ctx), "after is", err.Error())
			}
		}
	}()
	t, err := k.tableRepo.GetByID(ctx, k.db, &models.Table{
		TableID: req.TableID,
	})
	if err != nil {
		return nil, err
	}
	table := &models.Table{
		ID:      id2.GenID(),
		TableID: req.TableID,
		Schema:  req.Schema,
	}

	if t == nil { // 新增
		if req.Source == models.ModelSource { // 手工创建的表单 ,需要判重 ，以及 改变tableID
			//req.TableID = req.AppID + "_" + req.TableID
			dataSchema, err := k.databaseRepo.GetByCondition(ctx, k.db, req.AppID, req.TableID, "")
			if err != nil {
				return nil, err
			}
			if dataSchema != nil {
				return nil, error2.NewError(code.ErrRepeatTableID)
			}
			table.TableID = req.TableID
		}
		err = k.tableRepo.Create(ctx, k.db, table)
		if err != nil {
			return nil, err
		}
		// update page binding status
		if err = k.menuRepo.UpdateBindingStateByID(ctx, k.db, &models.Menu{
			ID:           req.TableID,
			BindingState: models.Bound,
		}); err != nil {
			return nil, err
		}
	}
	err = k.tableRepo.Update(ctx, k.db, table)
	return nil, err
}

func (k *kernel) after(ctx context.Context, req *CreateTableReq, options ...CreateSchemaOption) error {
	if c, ok := req.Schema["properties"].(map[string]interface{}); ok {
		convert, total, err := swagger.Convert1(c)
		if err != nil {
			return err
		}
		if req.Source == 0 {
			req.Source = models.FormSource
		}
		description, _ := req.Schema["description"].(string)
		optionReq := &OptionReq{
			Title:       req.Schema["title"].(string),
			Description: description,
			UserName:    req.UserName,
			UserID:      req.UserID,
			Source:      req.Source,
			Schema:      convert,
			Total:       total,
			AppID:       req.AppID,
			TableID:     req.TableID,
		}
		for _, options := range options {
			_, err := options(ctx, optionReq)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// GetXNameReq GetXNameReq
type GetXNameReq struct {
	TableID string `json:"TableID"`
	Action  string `json:"action"`
	AppID   string `json:"AppID"`
}

// GetXNameResp GetXNameResp
type GetXNameResp struct {
	Name string `json:"name"`
}

// GetXName GetXName
func (k *kernel) GetXName(ctx context.Context, req *GetXNameReq) (*GetXNameResp, error) {
	table, err := k.databaseRepo.GetByTableID(ctx, k.db, req.TableID)
	if err != nil {
		return nil, err
	}
	content := "form"
	if table.Source == models.ModelSource {
		content = "custom"
	}
	resp := &GetXNameResp{
		Name: swagger.GenXName(req.AppID, req.TableID, req.Action, content),
	}

	return resp, nil
}

// GetSchemaByTableID  GetSchemaByTableID
func (k *kernel) GetSchemaByTableID(ctx context.Context, req *GetTableReq) (*models.Table, error) {
	table := &models.Table{
		TableID: req.TableID,
	}
	r, err := k.tableRepo.GetByID(ctx, k.db, table)
	return r, err
}

// DeleteSchema DeleteSchema
func (k *kernel) DeleteSchema(ctx context.Context, req *DeleteTableReq) (*DeleteTableResp, error) {

	err := k.tableRepo.Delete(ctx, k.db, req.TableID)
	if err != nil {
		return nil, err
	}
	err = k.databaseRepo.Delete(ctx, k.db, req.TableID)
	return nil, err
}

// SearchSchemaReq SearchSchemaReq
type SearchSchemaReq struct {
	Title  string            `json:"title"`
	AppID  string            `json:"appID"`
	Page   int64             `json:"page"`
	Size   int64             `json:"size"`
	Source models.SourceType `json:"source"`
}

// SearchSchemaResp SearchSchemaResp
type SearchSchemaResp struct {
	List  []*tableVo `json:"list"`
	Total int64      `json:"total"`
}

// tableVo tableVo
type tableVo struct {
	ID          string            `json:"id"`
	TableID     string            `json:"tableID"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	FieldLen    int64             `json:"fieldLen"`
	Source      models.SourceType `json:"source"`
	CreatedAt   int64             `json:"createdAt"`
	UpdatedAt   int64             `json:"updatedAt"`
	Editor      string            `json:"editor"`
	CreatorName string            `json:"creatorName"`
}

func (k *kernel) SearchSchema(ctx context.Context, req *SearchSchemaReq) (*SearchSchemaResp, error) {
	tables, total, err := k.databaseRepo.Search(ctx, k.db, req.AppID, req.Title, req.Source, req.Size, req.Page)
	if err != nil {
		return nil, err
	}
	resp := &SearchSchemaResp{
		List: make([]*tableVo, len(tables)),
	}
	for index, v := range tables {
		vo := &tableVo{
			ID:          v.ID,
			TableID:     v.TableID,
			Title:       v.Title,
			Description: v.Description,
			Source:      v.Source,
			FieldLen:    v.FieldLen,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
			Editor:      v.EditorName,
			CreatorName: v.CreatorName,
		}
		resp.List[index] = vo
	}
	resp.Total = total
	return resp, nil
}
func (k *kernel) CreateConfig(ctx context.Context, req *CreateConfigReq) (*CreateConfigResp, error) {
	table := &models.Table{
		TableID: req.TableID,
		Config:  req.Config,
	}
	err := k.tableRepo.UpdateConfig(ctx, k.db, table)
	return nil, err
}

// CreateBlankSchema 相当于创建
func (k *kernel) CreateBlankSchema(ctx context.Context, req *CreateBlankTableReq) (*CreateBlankTableResp, error) {
	table := &models.Table{
		ID:      id2.GenID(),
		TableID: id2.GenUpperID(),
	}
	err := k.tableRepo.Create(ctx, k.db, table)
	if err != nil {
		return nil, err
	}
	resp := &CreateBlankTableResp{
		TableID: table.TableID,
	}
	return resp, nil
}

func (k *kernel) DeleteConfig(ctx context.Context, req *DeleteConfigReq) (*DeleteConfigResp, error) {
	table := &models.Table{
		ID:     req.ID,
		Config: nil,
	}
	err := k.tableRepo.DeleteConfig(ctx, k.db, table)
	return nil, err
}

// CheckRepeatReq CheckRepeatReq
type CheckRepeatReq struct {
	AppID    string `json:"appID"`
	TableID  string `json:"tableID"`
	Title    string `json:"title"`
	IsModify bool   `json:"isModify"`
}

// CheckRepeatResp CheckRepeatResp
type CheckRepeatResp struct {
}

func (k *kernel) CheckRepeat(ctx context.Context, req *CheckRepeatReq) (*CheckRepeatResp, error) {
	//req.TableID = req.AppID + "_" + req.TableID
	// name 如果等于空，检测tableID
	var tableSchema *models.DataBaseSchema
	// 先检测名字
	tableSchema, err := k.databaseRepo.GetByCondition(ctx, k.db, req.AppID, "", req.Title)
	if err != nil {
		return nil, err
	}
	if tableSchema != nil && tableSchema.TableID != req.TableID {
		return nil, error2.NewError(code.ErrRepeatTableTitle)
	}
	if req.IsModify {
		return &CheckRepeatResp{}, nil
	}
	tableSchema, err = k.databaseRepo.GetByCondition(ctx, k.db, req.AppID, req.TableID, "")
	if err != nil {
		return nil, err
	}
	if tableSchema != nil {
		return nil, error2.NewError(code.ErrRepeatTableID)
	}
	return &CheckRepeatResp{}, nil
}

// GetModelDataByMenuReq GetModelDataByMenuReq
type GetModelDataByMenuReq struct {
	TableID string `json:"menuId" binding:"required"`
}

// GetModelDataByMenuResp GetModelDataByMenuResp
type GetModelDataByMenuResp struct {
	TableID     string `json:"tableID"`
	FieldLen    int64  `json:"fieldLen"`
	CreatorName string `json:"createdBy"`
	CreatedAt   int64  `json:"createdAt"`
	EditorName  string `json:"updatedBy"`
	UpdatedAt   int64  `json:"updatedAt"`
}

func (k *kernel) GetModelDataByMenu(ctx context.Context, req *GetModelDataByMenuReq) (*GetModelDataByMenuResp, error) {
	modelData, err := k.databaseRepo.GetByTableID(ctx, k.db, req.TableID)
	if err != nil {
		return nil, err
	}
	return &GetModelDataByMenuResp{
		TableID:     modelData.TableID,
		FieldLen:    modelData.FieldLen,
		CreatorName: modelData.CreatorName,
		EditorName:  modelData.EditorName,
		UpdatedAt:   modelData.UpdatedAt,
		CreatedAt:   modelData.CreatedAt,
	}, nil
}
