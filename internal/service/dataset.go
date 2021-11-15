package service

import (
	"context"
	"git.internal.yunify.com/qxp/misc/error2"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/code"

	"git.internal.yunify.com/qxp/misc/id2"
	"git.internal.yunify.com/qxp/misc/time2"
	"git.internal.yunify.com/qxp/molecule/internal/models"
	repo "git.internal.yunify.com/qxp/molecule/internal/models/mongo"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"
	"go.mongodb.org/mongo-driver/mongo"
)

// DataSet DataSet
type DataSet interface {
	CreateDataSet(c context.Context, req *CreateDataSetReq) (*CreateDataSetResp, error)
	GetDataSet(c context.Context, req *GetDataSetReq) (*GetDataSetResp, error)
	UpdateDataSet(c context.Context, req *UpdateDataSetReq) (*UpdateDataSetResp, error)
	GetByConditionSet(c context.Context, req *GetByConditionSetReq) (*GetByConditionSetResp, error)
	DeleteDataSet(c context.Context, req *DeleteDataSetReq) (*DeleteDataSetResp, error)
}

type dataset struct {
	mongodb     *mongo.Database
	datasetRepo models.DataSetRepo
}

// NewDataSet NewDataSet
func NewDataSet(conf *config.Config, opts ...Options) (DataSet, error) {
	u := &dataset{
		datasetRepo: repo.NewDataSetRepo(),
	}
	for _, opt := range opts {
		opt(u)
	}
	return u, nil
}

// SetMongo SetMongo
func (per *dataset) SetMongo(client *mongo.Client, dbName string) {
	per.mongodb = client.Database(dbName)
}

// CreateDataSetReq CreateDataSetReq
type CreateDataSetReq struct {
	Name    string `json:"name" binding:"max=100"`
	Tag     string `json:"tag"  binding:"max=100"`
	Type    int64  `json:"type"`
	Content string `json:"content"`
}

// CreateDataSetResp CreateDataSetResp
type CreateDataSetResp struct {
	ID string `json:"id"`
}

// CreateDataSet CreateDataSet
func (per *dataset) CreateDataSet(c context.Context, req *CreateDataSetReq) (*CreateDataSetResp, error) {
	exist, err := per.datasetRepo.GetByName(c, per.mongodb, req.Name, "")
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, error2.NewError(code.ErrExistDataSetNameState)
	}
	dataset := &models.DataSet{
		ID:        id2.GenID(),
		Name:      req.Name,
		Tag:       req.Tag,
		Type:      req.Type,
		Content:   req.Content,
		CreatedAt: time2.NowUnix(),
	}
	err = per.datasetRepo.Insert(c, per.mongodb, dataset)
	if err != nil {
		return nil, err
	}
	return &CreateDataSetResp{
		ID: dataset.ID,
	}, nil
}

// GetDataSetReq GetDataSetReq
type GetDataSetReq struct {
	ID string `json:"id"`
}

// GetDataSetResp GetDataSetResp
type GetDataSetResp struct {
	ID      string `json:"id"`
	Name    string `json:"name" binding:"max=100"`
	Tag     string `json:"tag"  binding:"max=100"`
	Type    int64  `json:"type"`
	Content string `json:"content"`
}

// GetDataSet GetDataSet
func (per *dataset) GetDataSet(c context.Context, req *GetDataSetReq) (*GetDataSetResp, error) {
	dataset, err := per.datasetRepo.GetByID(c, per.mongodb, req.ID)
	if err != nil {
		return nil, err
	}
	if dataset == nil {
		return nil, error2.NewError(code.ErrNODataSetNameState)
	}
	resp := &GetDataSetResp{
		ID:      dataset.ID,
		Name:    dataset.Name,
		Tag:     dataset.Tag,
		Type:    dataset.Type,
		Content: dataset.Content,
	}
	return resp, nil
}

// UpdateDataSetReq UpdateDataSetReq
type UpdateDataSetReq struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Tag     string `json:"tag"`
	Type    int64  `json:"type"`
	Content string `json:"content"`
}

// UpdateDataSetResp UpdateDataSetResp
type UpdateDataSetResp struct {
}

// UpdateDataSet UpdateDataSet
func (per *dataset) UpdateDataSet(c context.Context, req *UpdateDataSetReq) (*UpdateDataSetResp, error) {
	exist, err := per.datasetRepo.GetByName(c, per.mongodb, req.Name, req.ID)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, error2.NewError(code.ErrExistDataSetNameState)
	}
	dataset := &models.DataSet{
		ID:      req.ID,
		Name:    req.Name,
		Tag:     req.Tag,
		Type:    req.Type,
		Content: req.Content,
	}
	err = per.datasetRepo.Update(c, per.mongodb, dataset)
	if err != nil {
		return nil, err
	}
	return &UpdateDataSetResp{}, nil
}

// GetByConditionSetReq GetByConditionSetReq
type GetByConditionSetReq struct {
	Name  string `json:"name"`
	Tag   string `json:"tag"`
	Types int64  `json:"type"`
}

// GetByConditionSetResp GetByConditionSetResp
type GetByConditionSetResp struct {
	List []*DataSetVo `json:"list"`
}

// DataSetVo DataSetVo
type DataSetVo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Tag     string `json:"tag"`
	Type    int64  `json:"type"`
	Content string `json:"content"`
}

// GetByConditionSet GetByConditionSet
func (per *dataset) GetByConditionSet(c context.Context, req *GetByConditionSetReq) (*GetByConditionSetResp, error) {
	arr, err := per.datasetRepo.GetByCondition(c, per.mongodb, req.Tag, req.Name, req.Types)
	if err != nil {
		return nil, err
	}
	resp := &GetByConditionSetResp{
		List: make([]*DataSetVo, len(arr)),
	}
	for index, value := range arr {
		resp.List[index] = new(DataSetVo)
		cloneDataSet(value, resp.List[index])
	}
	return resp, nil

}
func cloneDataSet(src *models.DataSet, dst *DataSetVo) {
	dst.ID = src.ID
	dst.Name = src.Name
	dst.Tag = src.Tag
	dst.Type = src.Type
	dst.Content = src.Content

}

// DeleteDataSetReq DeleteDataSetReq
type DeleteDataSetReq struct {
	ID string `json:"id"`
}

// DeleteDataSetResp DeleteDataSetResp
type DeleteDataSetResp struct {
}

// DeleteDataSet DeleteDataSet
func (per *dataset) DeleteDataSet(c context.Context, req *DeleteDataSetReq) (*DeleteDataSetResp, error) {
	err := per.datasetRepo.Delete(c, per.mongodb, req.ID)
	if err != nil {
		return nil, err
	}
	return &DeleteDataSetResp{}, nil
}
