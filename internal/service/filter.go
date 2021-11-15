package service

import (
	"context"
	"encoding/json"
	"git.internal.yunify.com/qxp/misc/id2"
	"git.internal.yunify.com/qxp/molecule/internal/filters"
	"git.internal.yunify.com/qxp/molecule/internal/models"
	repo "git.internal.yunify.com/qxp/molecule/internal/models/mongo"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"
	"go.mongodb.org/mongo-driver/mongo"
)

// Filter 业务接口
type Filter interface {
	SaveJSONFilter(c context.Context, req *SaveFilterReq) (*SaveFilterResp, error)

	GetJSONFilter(c context.Context, req *GetFilterReq) (*GetFilterResp, error)

	Delete(c context.Context, req *DELFilterReq) error

	InnerGetJSONFilter(c context.Context, req *GetFilterReq) (map[string]interface{}, error)
	// JSONFilter 吞字段
	JSONFilter(c context.Context, data interface{}, filter map[string]interface{}) error
	// SchemaFilter 对schema过滤
	SchemaFilter(c context.Context, data interface{}, filter map[string]interface{}) error
	// DataCheck 数据检查
	DataCheck(c context.Context, method string, data interface{}, filter map[string]interface{}) bool
}

// SaveFilterResp 设置返回体
type SaveFilterResp struct {
}

// GetFilterReq 请求体
type GetFilterReq struct {
	ID         string `json:"id"`
	PerGroupID string `json:"perGroupID"`
	FormID     string `json:"formID"`
}

// GetFilterResp 请求响应体
type GetFilterResp struct {
	ID         string         `json:"id"`
	PerGroupID string         `json:"perGroupID"`
	FormID     string         `json:"formID"`
	Schema     filters.Schema `json:"schema"`
}

// DELFilterReq 删除
type DELFilterReq struct {
	PerGroupID string `json:"perGroupID"`
	FormID     string `json:"formID"`
}

// SaveFilterReq 设置请求体
type SaveFilterReq struct {
	ID         string         `json:"id"`
	PerGroupID string         `json:"perGroupID"`
	FormID     string         `json:"formID"`
	Schema     filters.Schema `json:"schema"`
}

type filter struct {
	mongo      *mongo.Database
	filterRepo models.FilterRepo
}

func (f *filter) Delete(c context.Context, req *DELFilterReq) error {
	return f.filterRepo.Delete(c, f.mongo, req.PerGroupID, req.FormID)
}

func (f *filter) DataCheck(c context.Context, method string, data interface{}, filter map[string]interface{}) bool {
	switch method {
	case "create", "update", "update#pull", "update#push":
		if data == nil {
			return false
		}
		if filter == nil {
			return true
		}
		flag := filters.FilterCheckData(data, filter)
		return flag
	}
	return true
}

func (f *filter) SchemaFilter(c context.Context, data interface{}, filter map[string]interface{}) error {
	if data == nil {
		return nil
	}
	if filter == nil {
		return nil
	}
	filters.SchemaFilterToNewSchema2(data, filter)
	return nil
}

func (f *filter) JSONFilter(c context.Context, data interface{}, filter map[string]interface{}) error {
	if data == nil {
		return nil
	}
	if filter == nil {
		var empty interface{}
		data = empty
		return nil
	}
	filters.JSONFilter2(data, filter)
	return nil
}

func (f *filter) InnerGetJSONFilter(c context.Context, req *GetFilterReq) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	one, err := f.filterRepo.GetByCondition(c, f.mongo, req.PerGroupID, req.FormID)
	if err != nil {
		return nil, err
	}
	if one != nil {
		err := json.Unmarshal([]byte(one.FieldJSON), &m)
		if err != nil {
			return nil, err
		}
		return m, nil

	}
	return nil, nil
}

func (f *filter) SaveJSONFilter(c context.Context, req *SaveFilterReq) (*SaveFilterResp, error) {
	one, err := f.filterRepo.GetByCondition(c, f.mongo, req.PerGroupID, req.FormID)
	if err != nil {
		return nil, err
	}
	ft := models.Filter{}
	ft.FormID = req.FormID
	ft.PerGroupID = req.PerGroupID
	marshal, err := json.Marshal(req.Schema)
	if err != nil {
		return nil, err
	}
	ft.WebSchema = string(marshal)
	//将schema处理成过滤器需要的格式
	filterType := filters.DealSchemaToFilterType(req.Schema)

	bytes, err := json.Marshal(filterType)
	if err != nil {
		return nil, err
	}
	ft.FieldJSON = string(bytes)
	if one == nil {
		ft.ID = id2.GenID()
		err := f.filterRepo.Insert(c, f.mongo, &ft)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	err = f.filterRepo.Update(c, f.mongo, &ft)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (f filter) GetJSONFilter(c context.Context, req *GetFilterReq) (*GetFilterResp, error) {
	//需要和mongo配合
	res, err := f.filterRepo.GetByCondition(c, f.mongo, req.PerGroupID, req.FormID)
	if err != nil {
		return nil, err
	}
	if res != nil {
		schema := filters.Schema{}
		err := json.Unmarshal([]byte(res.WebSchema), &schema)
		if err != nil {
			return nil, err
		}
		return &GetFilterResp{
			Schema: schema,
			FormID: res.FormID,
		}, nil
	}
	return nil, nil
}

// NewFilter new
func NewFilter(conf *config.Config, opts ...Options) (Filter, error) {
	f := &filter{
		filterRepo: repo.NewFilterRepo(),
	}
	for _, opt := range opts {
		opt(f)
	}
	return f, nil
}
func (f *filter) SetMongo(client *mongo.Client, dbName string) {
	f.mongo = client.Database(dbName)
}
