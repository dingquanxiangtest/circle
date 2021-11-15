package comet

import (
	"context"
	"encoding/json"
	"errors"
	"git.internal.yunify.com/qxp/molecule/internal/dorm"
	"git.internal.yunify.com/qxp/molecule/internal/listener"
	"net/http"

	"git.internal.yunify.com/qxp/misc/header2"
	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/molecule/internal/dorm/clause"
	"git.internal.yunify.com/qxp/molecule/internal/dorm/dao"
	"git.internal.yunify.com/qxp/molecule/internal/models"
	"git.internal.yunify.com/qxp/molecule/internal/service"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Comet comet
type Comet interface {
	Handle(c *gin.Context)
	HandleWithoutAuth(c *gin.Context)
	Search(c *gin.Context)
}

type comet struct {
	cm         *CMongo
	permission service.Permission
	filter     service.Filter
	observers  []listener.Observer
}

type input struct {
	dao.FindOptions

	Method     string
	Conditions *InputConditionVo `json:"conditions"`
	Entity     Entity
}

// InputConditionVo InputConditionVo
type InputConditionVo struct {
	Condition []dorm.Condition `json:"condition"`
	Tag       string           `json:"tag"`
}

// New new a comet
func New(conf *config.Config, opts ...service.Options) (Comet, error) {
	p, err := service.NewPermission(conf, opts...)
	if err != nil {
		return nil, err
	}
	filter, err := service.NewFilter(conf, opts...)
	if err != nil {
		return nil, err
	}
	c := &comet{
		cm: &CMongo{
			dc:    clause.New(),
			query: dorm.NewQuery(),
			ag:    clause.NewAg(),
		},
		permission: p,
		filter:     filter,
		observers:  make([]listener.Observer, 0),
	}
	for _, opt := range opts {
		opt(c)
	}
	// 注册监听器
	process, err := NewProcess(conf, opts...)
	if err != nil {
		return nil, err
	}
	c.AddObserve(process)
	return c, nil
}

func (comet *comet) SetMongo(client *mongo.Client, dbName string) {
	comet.cm.DB = client.Database(dbName)
}

// HandleWithoutAuth HandleWithoutAuth
func (comet *comet) HandleWithoutAuth(c *gin.Context) {
	tableName, ok := getTableName(c)
	if !ok {
		err := errors.New("invalid URI")
		logger.Logger.Errorw(err.Error(), logger.GINRequestID(c))
		Format(c, WithError(err))
		return
	}
	logger.Logger.Infow(tableName, logger.GINRequestID(c))
	profile := header2.GetProfile(c)

	in := new(input)
	if err := getEntity(c, in); err != nil {
		logger.Logger.Errorw(err.Error(), logger.GINRequestID(c))
		Format(c, WithError(err))
		return
	}
	bus := &bus{
		tableName: tableName,
		input:     in,
		profile:   &profile,
		permissionGroup: &service.GetByConditionPerGroupResp{
			DataAccessPer: nil,
		},
	}
	ctx := logger.CTXTransfer(c)
	data, err := comet.cm.Handler(ctx, bus)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.GINRequestID(c))
		Format(c, WithError(err))
		return
	}
	comet.dataModifyEvent(ctx, bus, data)
	Format(c, WithPack(data))
}

// prefix 用户权限（接口权限、默认condition）
// suffix 字段权限 数据封装

func (comet *comet) Handle(c *gin.Context) {
	in := new(input)
	if err := getEntity(c, in); err != nil {
		logger.Logger.Errorw(err.Error(), logger.GINRequestID(c))
		Format(c, WithError(err))
	}
	bus, err := comet.pre(c, in.Method, in.Entity)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.GINRequestID(c))
		Format(c, WithError(err))
		return
	}
	bus.input = in
	if bus == nil || !checkOperatePermission(bus, bus.input.Method) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	ctx := logger.CTXTransfer(c)
	data, err := comet.cm.Handler(ctx, bus)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.GINRequestID(c))
		Format(c, WithError(err))
		return
	}
	comet.dataModifyEvent(ctx, bus, data)
	comet.postfix(c, bus, data)
}

type input1 struct {
	dao.FindOptions
	Conditions DSLQuery `json:"query"`

	Entity Entity `json:"entity"`

	Aggs DSLAgg `json:"aggs"`
}

func (comet comet) Search(c *gin.Context) {
	in1 := new(input1)
	if err := getEntity1(c, in1); err != nil {
		logger.Logger.Errorw(err.Error(), logger.GINRequestID(c))
		Format(c, WithError(err))
	}
	bus, err := comet.pre(c, "find", in1.Entity)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.GINRequestID(c))
		Format(c, WithError(err))
		return
	}
	if bus == nil || !checkOperatePermission(bus, "find") {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	bus.input1 = in1

	ctx := logger.CTXTransfer(c)
	data, err := comet.cm.Search(ctx, bus)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.GINRequestID(c))
		Format(c, WithError(err))
		return
	}
	comet.postfix(c, bus, data)
}

func getTableName(c *gin.Context) (string, bool) {
	return c.Params.Get("tableName")
}
func getAppID(c *gin.Context) (string, bool) {
	return c.Params.Get("appID")
}

func getEntity(c *gin.Context, input *input) error {
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(input)
	return err
}
func getEntity1(c *gin.Context, input *input1) error {
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(input)
	return err
}

const notAuthority = -1

func checkOperatePermission(bus *bus, method string) bool {
	if bus.permissionGroup.Authority == notAuthority {
		return true
	}
	var op int64
	switch method {
	case "find", "findOne":
		op = models.OPRead
	case "create":
		op = models.OPCreate
	case "update":
		op = models.OPUpdate
	case "delete":
		op = models.OPDelete
	default:
		return false
	}

	return op&bus.permissionGroup.Authority != 0
}

func (comet *comet) pre(c *gin.Context, method string, entity Entity) (*bus, error) {
	tableName, ok := getTableName(c)
	if !ok {
		return nil, errors.New("invalid URI")
	}
	appID, ok := getAppID(c)
	if !ok {
		return nil, errors.New("invalid URI")
	}
	profile := header2.GetProfile(c)

	ctx := logger.CTXTransfer(c)

	// 1、 得到用户组ID
	perGroupReq := &service.GetByConditionPerGroupReq{
		UserID: profile.UserID,
		DepID:  profile.DepartmentID,
		FormID: tableName,
		AppID:  appID,
	}
	perGroupResp, err := comet.permission.GetByConditionPerGroup(ctx, perGroupReq)
	if err != nil {
		return nil, err
	}
	if perGroupResp == nil || perGroupResp.ID == "" {
		return nil, errors.New("no permission")
	}
	filterReq := &service.GetFilterReq{
		FormID:     tableName,
		PerGroupID: perGroupResp.ID,
	}
	filter, err := comet.filter.InnerGetJSONFilter(ctx, filterReq)
	if err != nil {
		return nil, err
	}
	// 数据鉴定
	if !comet.filter.DataCheck(ctx, method, entity, filter) {
		return nil, errors.New("no permission")
	}
	return &bus{
		tableName:       tableName,
		profile:         &profile,
		permissionGroup: perGroupResp,
		filter:          filter,
	}, nil
}

func (comet *comet) postfix(c *gin.Context, bus *bus, pack Pack) {
	ctx := logger.CTXTransfer(c)
	err := comet.filter.JSONFilter(ctx, pack.Get(), bus.filter)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.GINRequestID(c))
	}
	Format(c, WithPack(pack))
}

func (comet *comet) dataModifyEvent(ctx context.Context, bus *bus, pack Pack) {
	switch bus.input.Method {
	case "create", "update", "delete":
		data := pack.Get()
		d := eventDataPack(bus, data)
		if d != nil {
			comet.Notify(ctx, d)
		}
	}
}

func eventDataPack(bus *bus, data interface{}) *DataModifyTrigger {
	entity := dataTrans(data)
	if len(entity) == 0 {
		return nil
	}
	return &DataModifyTrigger{
		TableID: bus.tableName,
		Entity:  &entity,
		Method:  bus.input.Method,
		UserID:  bus.profile.UserID,
		Topic:   commonTopic,
		Event:   dataModify,
	}
}
