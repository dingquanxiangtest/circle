package comet

import (
	"errors"
	"git.internal.yunify.com/qxp/misc/header2"
	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/molecule/internal/filters"
	"git.internal.yunify.com/qxp/molecule/internal/models"
	"git.internal.yunify.com/qxp/molecule/internal/service"
	"github.com/gin-gonic/gin"
)

type comet1 struct {
	cm         *CMongo
	permission service.Permission
	filter     service.Filter
	schema     service.Kernel
}

func (c1 *comet1) SchemaHandle(c *gin.Context, bus *bus) (Pack, error) {
	ctx := logger.CTXTransfer(c)
	req := service.GetTableReq{
		TableID: bus.tableName,
	}
	r, err := c1.schema.GetSchemaByTableID(ctx, &req)
	if err != nil {
		return nil, err
	}
	sm := &r.Schema
	err = c1.filter.SchemaFilter(ctx, sm, bus.filter)
	if err != nil {
		return nil, err
	}
	resp := &Schemas{
		ID:      r.ID,
		TableID: r.TableID,
		Schema:  &r.Schema,
		Config:  &r.Config,
	}
	return resp, nil

}

func (c1 *comet1) DataHandle(c *gin.Context, bus *bus) (Pack, error) {
	ctx := logger.CTXTransfer(c)
	return c1.cm.Handler(ctx, bus)
}

func (c1 *comet1) Search(c *gin.Context, bus *bus) (Pack, error) {
	ctx := logger.CTXTransfer(c)
	data, err := c1.cm.Search(ctx, bus)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c1 *comet1) Create(c *gin.Context, bus *bus) (Pack, error) {
	ctx := logger.CTXTransfer(c)
	data, err := c1.cm.Create(ctx, bus)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c1 *comet1) Update(c *gin.Context, bus *bus) (Pack, error) {
	ctx := logger.CTXTransfer(c)
	data, err := c1.cm.Update(ctx, bus)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c1 *comet1) Delete(c *gin.Context, bus *bus) (Pack, error) {
	ctx := logger.CTXTransfer(c)
	data, err := c1.cm.Delete(ctx, bus)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// CheckURL CheckURL
func CheckURL(c *gin.Context, bus *bus) (err error) {
	tableName, ok := getTableName(c)
	if !ok {
		return errors.New("invalid URI")
	}
	appID, ok := getAppID(c)
	if !ok {
		return errors.New("invalid URI")
	}
	bus.tableName = tableName
	bus.AppID = appID
	return nil
}

// commonPre all
func commonPre(ctx *gin.Context, b *bus) error {
	err := CheckURL(ctx, b) // 校验参数
	in1 := new(input1)
	if err := getEntity1(ctx, in1); err != nil {
		return err
	}
	b.input1 = in1
	p := header2.GetProfile(ctx)
	b.profile = &p
	if err != nil {
		return err
	}
	return nil

}

// Postfix Postfix
func (c1 *comet1) Postfix(c *gin.Context, bus *bus, pack Pack, opts ...FilterOption) error {
	for _, opt := range opts {
		err := opt(pack.Get(), bus.filter)
		if err != nil {
			return err
		}
	}
	Format(c, WithPack(pack))
	return nil
}

type c1 interface {
}

// SearchData SearchData
func SearchData(ce c1) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var (
			pack Pack = &Body{}
			err  error
		)
		b := &bus{}
		err = commonPre(ctx, b)
		if err != nil {
			logger.Logger.Errorw(err.Error(), logger.GINRequestID(ctx))
			Format(ctx, WithError(err))
			return
		}
		err = pre(ctx, b, ce, "find", b.input1.Entity, CheckOperate(), CheckData())
		if err != nil {
			logger.Logger.Errorw(err.Error(), logger.GINRequestID(ctx))
			Format(ctx, WithError(err))
			return
		}
		if handle, ok := ce.(Search); ok {
			pack, err = handle.Search(ctx, b)
			if err != nil {
				logger.Logger.Errorw(err.Error(), logger.GINRequestID(ctx))
				Format(ctx, WithError(err))
				return
			}
		}
		post(ctx, b, pack, ce, JSONFilter())
	}
}

// CreateData CreateData
func CreateData(ce c1) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var (
			pack Pack = &Body{}
			err  error
		)
		b := &bus{}
		err = commonPre(ctx, b)
		if err != nil {
			logger.Logger.Errorw(err.Error(), logger.GINRequestID(ctx))
			Format(ctx, WithError(err))
			return
		}
		err = pre(ctx, b, ce, "create", b.input1.Entity, CheckOperate(), CheckData())
		if err != nil {
			logger.Logger.Errorw(err.Error(), logger.GINRequestID(ctx))
			Format(ctx, WithError(err))
			return
		}
		if handle, ok := ce.(Create); ok {
			pack, err = handle.Create(ctx, b)
			if err != nil {
				logger.Logger.Errorw(err.Error(), logger.GINRequestID(ctx))
				Format(ctx, WithError(err))
				return
			}
		}
		post(ctx, b, pack, ce, JSONFilter())
	}
}

// UpdateData UpdateData
func UpdateData(ce c1) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var (
			pack Pack = &Body{}
			err  error
		)
		b := &bus{}
		err = commonPre(ctx, b)
		if err != nil {
			logger.Logger.Errorw(err.Error(), logger.GINRequestID(ctx))
			Format(ctx, WithError(err))
			return
		}
		err = pre(ctx, b, ce, "update", b.input1.Entity, CheckOperate(), CheckData())
		if err != nil {
			logger.Logger.Errorw(err.Error(), logger.GINRequestID(ctx))
			Format(ctx, WithError(err))
			return
		}
		if handle, ok := ce.(Update); ok {
			pack, err = handle.Update(ctx, b)
			if err != nil {
				logger.Logger.Errorw(err.Error(), logger.GINRequestID(ctx))
				Format(ctx, WithError(err))
				return
			}
		}
		post(ctx, b, pack, ce, JSONFilter())

	}

}

// DeleteData  DeleteData
func DeleteData(ce c1) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var (
			pack Pack = &Body{}
			err  error
		)
		b := &bus{}
		err = commonPre(ctx, b)
		if err != nil {
			logger.Logger.Errorw(err.Error(), logger.GINRequestID(ctx))
			Format(ctx, WithError(err))
			return
		}
		err = pre(ctx, b, ce, "delete", b.input1.Entity, CheckOperate(), CheckData())
		if err != nil {
			logger.Logger.Errorw(err.Error(), logger.GINRequestID(ctx))
			Format(ctx, WithError(err))
			return
		}
		if handle, ok := ce.(Delete); ok {
			pack, err = handle.Delete(ctx, b)
			if err != nil {
				logger.Logger.Errorw(err.Error(), logger.GINRequestID(ctx))
				Format(ctx, WithError(err))
				return
			}
		}
		post(ctx, b, pack, ce, JSONFilter())

	}

}

// HandleData HandleData
func HandleData(ce c1) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var (
			pack Pack = &Body{}
			err  error
		)
		b := &bus{}
		err = CheckURL(ctx, b) // 校验参数
		if err != nil {
			logger.Logger.Errorw(err.Error(), logger.GINRequestID(ctx))
			Format(ctx, WithError(err))
			return
		}
		// 接收参数
		in := new(input)
		if err := getEntity(ctx, in); err != nil {
			logger.Logger.Errorw(err.Error(), logger.GINRequestID(ctx))
			Format(ctx, WithError(err))
			return
		}
		b.input = in

		err = pre(ctx, b, ce, b.input.Method, b.input.Entity, CheckOperate(), CheckData())
		if err != nil {
			logger.Logger.Errorw(err.Error(), logger.GINRequestID(ctx))
			Format(ctx, WithError(err))
			return
		}

		if handle, ok := ce.(DataHandle); ok {
			pack, err = handle.DataHandle(ctx, b)
			if err != nil {
				logger.Logger.Errorw(err.Error(), logger.GINRequestID(ctx))
				Format(ctx, WithError(err))
				return
			}

		}

		post(ctx, b, pack, ce, JSONFilter())
	}
}

// HandleSchema HandleSchema
func HandleSchema(ce c1) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var (
			pack Pack = &Body{}
			err  error
		)
		b := &bus{}
		err = pre(ctx, b, ce, "", nil)
		if err != nil {
			logger.Logger.Errorw(err.Error(), logger.GINRequestID(ctx))
			Format(ctx, WithError(err))
			return
		}
		if handle, ok := ce.(SchemaHandle); ok {
			pack, err = handle.SchemaHandle(ctx, b)
			if err != nil {
				logger.Logger.Errorw(err.Error(), logger.GINRequestID(ctx))
				Format(ctx, WithError(err))
				return
			}

		}
		post(ctx, b, pack, ce, SchemaFilters())

	}
}

func pre(ctx *gin.Context, b *bus, ce c1, method string, entity Entity, opts ...PreOption) (error error) {
	if pre, ok := ce.(Pre); ok {
		err := pre.Pre(ctx, b, method, entity, opts...)
		if err != nil {
			return err
		}
	}
	return nil
}

func post(ctx *gin.Context, b *bus, pack Pack, ce c1, opts ...FilterOption) {
	if post, ok := ce.(Post); ok {
		post.Postfix(ctx, b, pack, opts...)
	}
}

//FilterOption FilterOption
type FilterOption func(data interface{}, filter map[string]interface{}) error

// SchemaFilters SchemaFilters
func SchemaFilters() FilterOption {
	return func(data interface{}, filter map[string]interface{}) error {
		if data == nil {
			return nil
		}
		if filter == nil {
			return nil
		}
		filters.SchemaFilterToNewSchema2(data, filter)
		return nil
	}
}

// JSONFilter JSONFilter
func JSONFilter() FilterOption {
	return func(data interface{}, filter map[string]interface{}) error {
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
}

// PreOption PreOption
type PreOption func(method string, bus *bus, entity Entity) bool

// CheckOperate CheckOperate
func CheckOperate() PreOption {
	return func(method string, bus *bus, entity Entity) bool {
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
}

// CheckData  CheckData
func CheckData() PreOption {
	return func(method string, bus *bus, entity Entity) bool {
		switch method {
		case "create", "update", "update#pull", "update#push":
			if entity == nil {
				return false
			}
			if bus.filter == nil {
				return true
			}
			flag := filters.FilterCheckData(entity, bus.filter)
			return flag
		}
		return true
	}
}
