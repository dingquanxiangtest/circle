package comet

import (
	"errors"
	"git.internal.yunify.com/qxp/misc/header2"
	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/molecule/internal/dorm"
	"git.internal.yunify.com/qxp/molecule/internal/dorm/clause"
	"git.internal.yunify.com/qxp/molecule/internal/service"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Plugs  Plugs
type Plugs interface {
	SchemaHandle
	DataHandle
	Search
}

type auth struct {
	comet1
}

// NewAuth NewAuth
func NewAuth(conf *config.Config, opts ...service.Options) (Plugs, error) {
	p, err := service.NewPermission(conf, opts...)
	if err != nil {
		return nil, err
	}
	filter, err := service.NewFilter(conf, opts...)
	if err != nil {
		return nil, err
	}
	t, err := service.NewKernel(conf, opts...)
	if err != nil {
		return nil, err
	}
	a := &auth{
		comet1: comet1{
			cm: &CMongo{
				dc:    clause.New(),
				query: dorm.NewQuery(),
				ag:    clause.NewAg(),
			},
			permission: p,
			filter:     filter,
			schema:     t,
		},
	}

	for _, opt := range opts {
		opt(a)
	}
	return a, nil
}

// SetMongo SetMongo
func (c1 *auth) SetMongo(client *mongo.Client, dbName string) {
	c1.cm.DB = client.Database(dbName)
}

// Pre Pre
func (c1 *auth) Pre(c *gin.Context, bus *bus, method string, entity Entity, opts ...PreOption) error {
	profile := header2.GetProfile(c)
	ctx := logger.CTXTransfer(c)
	// 1、 得到用户组ID
	perGroupReq := &service.GetByConditionPerGroupReq{
		UserID: profile.UserID,
		DepID:  profile.DepartmentID,
		FormID: bus.tableName,
		AppID:  bus.AppID,
	}

	perGroupResp, err := c1.permission.GetByConditionPerGroup(ctx, perGroupReq)
	if err != nil {
		return err
	}
	if perGroupResp == nil || perGroupResp.ID == "" {
		return errors.New("no permission")
	}
	filterReq := &service.GetFilterReq{
		FormID:     bus.tableName,
		PerGroupID: perGroupResp.ID,
	}
	filter, err := c1.filter.InnerGetJSONFilter(ctx, filterReq)
	if err != nil {
		return err
	}
	bus.permissionGroup = perGroupResp
	bus.profile = &profile
	bus.filter = filter

	err = after(bus, entity, method, opts...)
	if err != nil {
		return err
	}
	return nil
}

func after(bus *bus, entity Entity, method string, opts ...PreOption) error {
	for _, opt := range opts {
		if !opt(method, bus, entity) {
			return errors.New("no permission")
		}
	}
	return nil
}
