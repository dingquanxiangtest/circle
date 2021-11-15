package service

import (
	"context"
	"git.internal.yunify.com/qxp/misc/id2"
	"git.internal.yunify.com/qxp/molecule/internal/models"
	repo "git.internal.yunify.com/qxp/molecule/internal/models/mongo"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"
	"go.mongodb.org/mongo-driver/mongo"
)

// SubTable 业务逻辑
type SubTable interface {
	CreateSubTable(ctx context.Context, req *CreateSubTableReq) (*CreateSubTableResp, error)
	GetSubTableByTableID(ctx context.Context, req *GetSubTableReq) (*GetSubTableResp, error)
	GetSubTableByCondition(ctx context.Context, req *GetSubTableReq) (*models.SubTable, error)
	GetSubTableByType(ctx context.Context,req *GetSubTableReq)([]*models.SubTable, error)
	UpdateSubTable(ctx context.Context, req *UpdateSubTableReq) (*UpdateSubTableResp, error)
	DeleteSubTable(ctx context.Context, req *DeleteSubTableReq) (*DeleteSubTableResp, error)
}

type subTable struct {
	db *mongo.Database
	subTableRepo models.SubTableRepo
}

// NewSubTable new a subTable
func NewSubTable(conf *config.Config, opts ...Options) (SubTable, error) {
	subTableRepo := repo.NewSubTableRepo()
	m := &subTable{
		subTableRepo: subTableRepo,
	}

	for _, opt := range opts {
		opt(m)
	}
	return m, nil
}

func (s *subTable) SetMongo(client *mongo.Client, dbName string) {
	s.db = client.Database(dbName)
}

// CreateSubTableReq CreateSubTableReq
type CreateSubTableReq struct {
	AppID           string      `json:"appID"`
	TableID         string      `json:"tableID"`
	SubTableID      string      `json:"subTableID"`
	FieldName       string      `json:"fieldName"`
	SubTableType    string      `json:"subTableType"`
	Filter          []string    `json:"filter"`
	AggCondition    *models.AggregationCondition `json:"aggCondition"`
}



// CreateSubTableResp CreateSubTableResp
type CreateSubTableResp struct {}

func (s *subTable)CreateSubTable(ctx context.Context, req *CreateSubTableReq) (*CreateSubTableResp, error)  {

	subTable := &models.SubTable{
		ID: id2.GenID(),
		AppID: req.AppID,
		TableID: req.TableID,
		SubTableID: req.SubTableID,
		FieldName: req.FieldName,
		Filter: req.Filter,
		SubTableType: req.SubTableType,
		AggCondition: req.AggCondition,
	}
	err := s.subTableRepo.Create(ctx,s.db,subTable)
	return nil,err
}

// GetSubTableReq GetSubTableReq
type GetSubTableReq struct {
	TableID         string  `json:"tableID"`
	SubTableID      string  `json:"subTableID"`
	FieldName       string   `json:"fieldName"`
	SubTableType    string   `json:"subTableType"`
}

// GetSubTableResp GetSubTableResp
type  GetSubTableResp struct {
	SubTables     []*models.SubTable  `json:"subTables"`
}

func (s *subTable)GetSubTableByTableID(ctx context.Context, req *GetSubTableReq) (*GetSubTableResp, error){
	subTable := &models.SubTable{
		TableID: req.TableID,
	}
	r,err := s.subTableRepo.GetByID(ctx,s.db,subTable)
	if err != nil{
		return nil, err
	}
	resp := &GetSubTableResp{
		SubTables: r,
	}
	return resp,nil
}

func (s *subTable)GetSubTableByCondition(ctx context.Context, req *GetSubTableReq) (*models.SubTable, error){
	subTable := &models.SubTable{
		TableID: req.TableID,
		SubTableID: req.SubTableID,
		FieldName: req.FieldName,
	}
	r,err := s.subTableRepo.GetByCondition(ctx,s.db,subTable)
	if err != nil{
		return nil, err
	}

	return r,nil
}

// UpdateSubTableReq UpdateSubTableReq
type UpdateSubTableReq struct {
	TableID         string  `json:"tableID"`
	SubTableID      string `json:"subTableID"`
	FieldName       string `json:"fieldName"`
}

// UpdateSubTableResp UpdateSubTableResp
type UpdateSubTableResp struct {}

func (s *subTable)UpdateSubTable(ctx context.Context, req *UpdateSubTableReq) (*UpdateSubTableResp, error) {
	subTable := &models.SubTable{
		TableID: req.TableID,
		SubTableID: req.SubTableID,
		FieldName: req.FieldName,
	}
	err := s.subTableRepo.Update(ctx,s.db,subTable)
	return nil, err
}

// DeleteSubTableReq DeleteSubTableReq
type DeleteSubTableReq struct {
	TableID         string  `json:"tableID"`
	SubTableID      string `json:"subTableID"`
}

// DeleteSubTableResp DeleteSubTableResp
type DeleteSubTableResp struct {}

func (s *subTable)DeleteSubTable(ctx context.Context, req *DeleteSubTableReq) (*DeleteSubTableResp, error) {
	subTable := &models.SubTable{
		TableID: req.TableID,
		SubTableID: req.SubTableID,
	}
	err := s.subTableRepo.Delete(ctx,s.db,subTable)
	return nil,err
}

func (s *subTable)GetSubTableByType(ctx context.Context,req *GetSubTableReq)([]*models.SubTable, error) {
	subTable := &models.SubTable{
		TableID      : req.TableID,
		SubTableType : req.SubTableType,
	}
	r,err := s.subTableRepo.GetSubTableByType(ctx,s.db,subTable)
	if err != nil {
		return nil, err
	}
	return r, nil
}
