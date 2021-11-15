package comet

import (
	"encoding/json"
	"errors"
	"net/http"

	"git.internal.yunify.com/qxp/misc/header2"
	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/molecule/internal/models"
	"git.internal.yunify.com/qxp/molecule/internal/service"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Form form
type Form interface {
	Handle(c *gin.Context)
}

type table struct {
	cm         *CMongo
	permission service.Permission
	filter     service.Filter
	schema     service.Kernel
}

// SchemaResp SchemaResp
type SchemaResp struct {
	ID      string `json:"id"`
	TableID string `json:"tableID"`
	Schema interface{} `json:"schema"`
	Config interface{} `json:"config"`
}

// NewTable new a schema
func NewTable(conf *config.Config, opts ...service.Options) (Form, error) {
	p, err := service.NewPermission(conf, opts...)
	if err != nil {
		return nil, err
	}
	f, err := service.NewFilter(conf, opts...)
	if err != nil {
		return nil, err
	}
	t, err := service.NewKernel(conf, opts...)
	if err != nil {
		return nil, err
	}
	c := &table{
		cm:         &CMongo{},
		permission: p,
		filter:     f,
		schema:     t,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

func (t *table) SetMongo(client *mongo.Client, dbName string) {
	t.cm.DB = client.Database(dbName)
}

func (t *table) Handle(c *gin.Context) {
	tableName, ok := getTableName(c)

	if !ok {
		err := errors.New("invalid URI")
		logger.Logger.Errorw(err.Error(), logger.GINRequestID(c))
		Format(c, WithError(err))
		return
	}
	appID, ok := getAppID(c)
	if !ok {
		err := errors.New("invalid URI")
		logger.Logger.Errorw(err.Error(), logger.GINRequestID(c))
		Format(c, WithError(err))
		return
	}
	profile := header2.GetProfile(c)

	perGroupReq := &service.GetByConditionPerGroupReq{
		UserID: profile.UserID,
		DepID:  profile.DepartmentID,
		FormID: tableName,
		AppID:  appID,
	}
	ctx := logger.CTXTransfer(c)
	perGroupResp, err := t.permission.GetByConditionPerGroup(ctx, perGroupReq)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.GINRequestID(c))
		Format(c, WithError(err))
		return
	}
	if perGroupResp == nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	filerReq := &service.GetFilterReq{
		FormID:     tableName,
		PerGroupID: perGroupResp.ID,
	}
	filter, err := t.filter.InnerGetJSONFilter(ctx, filerReq)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.GINRequestID(c))
		Format(c, WithError(err))
		return
	}

	bus := &bus{
		tableName: tableName,
		//permissionID: perGroupResp.ID,
		filter: filter,
	}
	t.postfix(c, bus)
}

func (t *table) postfix(c *gin.Context, bus *bus) {
	ctx := logger.CTXTransfer(c)
	req := service.GetTableReq{
		TableID: bus.tableName,
	}
	r, err := t.schema.GetSchemaByTableID(ctx, &req)
	if err != nil {
		Format(c, WithError(err))
		return
	}
	s, err := json.Marshal(r)
	if err != nil {
		Format(c, WithError(err))
		return
	}

	table := models.Table{}
	err = json.Unmarshal(s, &table)
	if err != nil {
		Format(c, WithError(err))
		return
	}

	sm := &table.Schema
	// filter schema
	err = t.filter.SchemaFilter(ctx, sm, bus.filter)
	if err != nil {
		Format(c, WithError(err))
		return
	}
	resp := SchemaResp{
		ID:         table.ID,
		TableID:    table.TableID,
		Schema:     &table.Schema,
		Config:     &table.Config,
	}
	Format(c, WithSchema(&resp))
}
