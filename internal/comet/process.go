package comet

import (
	"context"
	"git.internal.yunify.com/qxp/misc/header2"
	"git.internal.yunify.com/qxp/molecule/internal/dorm/clause"
	"git.internal.yunify.com/qxp/molecule/internal/dorm/dao"
	"git.internal.yunify.com/qxp/molecule/internal/filters"
	"git.internal.yunify.com/qxp/molecule/internal/listener"
	"git.internal.yunify.com/qxp/molecule/internal/models"
	repo "git.internal.yunify.com/qxp/molecule/internal/models/mongo"
	"git.internal.yunify.com/qxp/molecule/internal/service"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

// Process trigger handler
type Process interface {
	// GetTableByID query schema
	GetTableByID(ctx context.Context, req *service.GetTableReq) (*ProcessTableResp, error)
	// GetRawTableByID query schema
	GetRawTableByID(ctx context.Context, req *ProcessTableReq) (*ProcessTableResp, error)
	// GetData BatchGetData
	GetData(ctx context.Context, req *GetDataReq) (*CreateProcessResp, error)
	// BatchGetData BatchGetData
	BatchGetData(ctx context.Context, req *BatchGetDataReq) (*BatchGetDataResp, error)
	// UpdateProcessData CreateProcessData
	UpdateProcessData(ctx context.Context, c *gin.Context, req *CreateProcessDataReq) (*UpdateProcessResp, error)
	// CreateProcessData CreateProcessData
	CreateProcessData(ctx context.Context, c *gin.Context, req *CreateProcessDataReq) (*CreateProcessResp, error)
	// GetSubTable CreateProcessData
	GetSubTable(ctx context.Context, req *service.GetSubTableReq) (*models.SubTable, error)
}

type processServer struct {
	cm        *CMongo
	tableRepo models.TableRepo
	subTable  service.SubTable
	observers []listener.Observer
}

// NewProcessServer new
func NewProcessServer(conf *config.Config, opts ...service.Options) (Process, error) {
	s, err := service.NewSubTable(conf, opts...)
	if err != nil {
		return nil, err
	}
	p := &processServer{
		cm:        &CMongo{dc: clause.New()},
		tableRepo: repo.NewTableRepo(),
		subTable:  s,
		observers: make([]listener.Observer, 0),
	}
	for _, opt := range opts {
		opt(p)
	}
	process, err := NewProcess(conf, opts...)
	if err != nil {
		return nil, err
	}
	p.AddObserve(process)
	return p, nil
}
func (p *processServer) SetMongo(client *mongo.Client, dbName string) {
	p.cm.DB = client.Database(dbName)
}

// CreateProcessResp CreateProcessResp
type CreateProcessResp struct {
	Entity interface{} `json:"entity"`
}

// UpdateProcessResp UpdateProcessResp
type UpdateProcessResp struct{}

// ProcessTableReq ProcessTableReq
type ProcessTableReq struct {
	TableID string `json:"tableID"`
}

func (p *processServer) GetSubTable(ctx context.Context, req *service.GetSubTableReq) (*models.SubTable, error) {
	r, err := p.subTable.GetSubTableByCondition(ctx, req)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// ProcessTableResp ProcessTableResp
type ProcessTableResp struct {
	Schema map[string]interface{} `json:"schema"`
}

func (p *processServer) GetTableByID(ctx context.Context, req *service.GetTableReq) (*ProcessTableResp, error) {
	table := &models.Table{
		TableID: req.TableID,
	}
	r, err := p.tableRepo.GetByID(ctx, p.cm.DB, table)
	if r == nil || err != nil {
		return nil, err
	}

	rest := make(map[string]interface{})
	filters.SchemaLoseWeight(r.Schema, rest)

	resp := &ProcessTableResp{
		Schema: rest,
	}
	return resp, nil
}

func (p *processServer) GetRawTableByID(ctx context.Context, req *ProcessTableReq) (*ProcessTableResp, error) {
	table := &models.Table{
		TableID: req.TableID,
	}
	r, err := p.tableRepo.GetByID(ctx, p.cm.DB, table)
	if r == nil || err != nil {
		return nil, err
	}
	resp := &ProcessTableResp{
		Schema: r.Schema,
	}
	return resp, nil
}

// BatchGetDataReq BatchGetDataReq
type BatchGetDataReq struct {
	TableID string `json:"tableID"`
	Input   *input `json:"input"`
}

// BatchGetDataResp BatchGetDataResp
type BatchGetDataResp struct {
	Entity interface{} `json:"entities"`
}

// GetDataReq GetDataReq
type GetDataReq struct {
	TableID string `json:"tableID"`
	ID      string `json:"id"`
	Input   *input `json:"input"`
}

func (p *processServer) BatchGetData(ctx context.Context, req *BatchGetDataReq) (*BatchGetDataResp, error) {
	if req == nil {
		return nil, nil
	}
	bus := packBus(req.TableID, nil, req.Input)
	data, err := p.cm.Handler(ctx, bus)
	if err != nil {
		return nil, err
	}
	da := reflect.ValueOf(data.Get())
	d := da.Interface().([]dao.Data)
	return &BatchGetDataResp{
		Entity: d,
	}, err
}

// GetDataResp GetDataResp
type GetDataResp []interface{}

func (p *processServer) GetData(ctx context.Context, req *GetDataReq) (*CreateProcessResp, error) {
	if req == nil {
		return nil, nil
	}
	bus := packBus(req.TableID, nil, req.Input)
	data, err := p.cm.Handler(ctx, bus)
	if err != nil {
		return nil, err
	}
	da := reflect.ValueOf(data.Get())
	d := map[string]interface{}(da.Interface().(dao.Data))
	return &CreateProcessResp{
		Entity: d,
	}, err
}

// CreateProcessDataReq CreateProcessDataReq
type CreateProcessDataReq struct {
	ID           string `json:"id"`
	TableID      string `json:"tableID"`
	Entity       Entity `json:"entity"`
	EventTrigger bool   `json:"eventTrigger"`
	Input        *input `json:"input"`
}

func (p *processServer) UpdateProcessData(ctx context.Context, c *gin.Context, req *CreateProcessDataReq) (*UpdateProcessResp, error) {
	profile := header2.GetProfile(c)
	bus := packBus(req.TableID, &profile, req.Input)
	r, err := p.cm.Handler(ctx, bus)
	if err == nil {
		if req.EventTrigger {
			d := eventDataPack(bus, r.Get())
			if d != nil {
				p.Notify(ctx, d)
			}
		}
	}
	return nil, err
}

func (p *processServer) CreateProcessData(ctx context.Context, c *gin.Context, req *CreateProcessDataReq) (*CreateProcessResp, error) {
	profile := header2.GetProfile(c)
	bus := packBus(req.TableID, &profile, req.Input)
	r, err := p.cm.Handler(ctx, bus)
	if err != nil {
		return nil, err
	}
	if req.EventTrigger {
		d := eventDataPack(bus, r.Get())
		if d != nil {
			p.Notify(ctx, d)
		}
	}
	return &CreateProcessResp{
		Entity: r.Get(),
	}, err
}

func packBus(tableName string, profile *header2.Profile, in *input) *bus {
	return &bus{
		tableName: tableName,
		input:     in,
		profile:   profile,
		permissionGroup: &service.GetByConditionPerGroupResp{
			DataAccessPer: nil,
		},
	}
}
