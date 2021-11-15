package comet

import (
	"context"
	"errors"
	"git.internal.yunify.com/qxp/misc/error2"
	"git.internal.yunify.com/qxp/misc/header2"
	"git.internal.yunify.com/qxp/molecule/internal/dorm"
	"git.internal.yunify.com/qxp/molecule/internal/dorm/clause"
	"git.internal.yunify.com/qxp/molecule/internal/dorm/dao"
	"git.internal.yunify.com/qxp/molecule/internal/models"
	"git.internal.yunify.com/qxp/molecule/internal/service"
	ce "git.internal.yunify.com/qxp/molecule/pkg/misc/code"
	"go.mongodb.org/mongo-driver/mongo"
)

// CMongo mongo handler
type CMongo struct {
	dc *clause.Clause

	DB *mongo.Database

	query *dorm.Query

	ag *clause.Aggregate
}

type bus struct {
	tableName       string
	AppID           string
	input           *input
	input1          *input1
	profile         *header2.Profile
	permissionGroup *service.GetByConditionPerGroupResp
	filter          map[string]interface{}
}

// DSLQuery DSLQuery
type DSLQuery map[string]interface{}

// DSLAgg DSLAgg
type DSLAgg map[string]interface{}

func aggregateCondition(src *models.ConditionVO) *InputConditionVo {
	vo := &InputConditionVo{
		Tag:       src.Tag,
		Condition: make([]dorm.Condition, 0),
	}
	for _, value := range src.Arr {
		v := dorm.Condition{
			Value: value.Value,
			Op:    value.OP,
			Key:   value.Key,
		}
		vo.Condition = append(vo.Condition, v)
	}
	return vo
}

// Handler handling the API corresponding to JSON schema
func (cm *CMongo) Handler(ctx context.Context, bus *bus) (Pack, error) {
	var (
		pack Pack = &Body{}
		err  error
	)
	var fail = func(l int64, cs []dorm.Condition) int64 {
		if len(cs) == 0 {
			return 0
		}
		return int64(len(cs[0].Value)) - l
	}
	var dm dao.Dao
	dm = &dao.MONGO{C: cm.DB.Collection(bus.tableName)}
	builder := &clause.MONGO{}

	dataAccess := make([]*InputConditionVo, 0)

	method := bus.input.Method
	if bus.input.Method == "findOne" {
		method = "find"
	}
	//
	if bus.permissionGroup != nil && len(bus.permissionGroup.DataAccessPer) > 0 {
		filterCondition := bus.permissionGroup.DataAccessPer[method]
		if filterCondition != nil {
			dataAccess = append(dataAccess, aggregateCondition(filterCondition))
		}
		if bus.input.Method == "delete" || bus.input.Method == "update" {
			filterCondition = bus.permissionGroup.DataAccessPer["find"]
			if filterCondition != nil {
				dataAccess = append(dataAccess, aggregateCondition(filterCondition))
			}
		}
	}

	condition, err := cm.Condition(dataAccess, bus.input.Conditions)
	if err != nil {
		return nil, err
	}
	if condition != nil {
		condition.MongoBuild(builder)
	}
	switch bus.input.Method {
	case "find":
		paging := &Paging{}
		paging.Data, err = dm.Find(ctx, builder, bus.input.FindOptions)
		if err != nil {
			break
		}
		paging.Total, err = dm.Count(ctx, builder)
		pack = paging
	case "findOne":
		body := &Body{}
		body.Data, err = dm.FindOne(ctx, builder)
		pack = body
	case "create":
		bus.input.Entity = defaultField(bus.input.Entity,
			WithID(),
			WithCreated(bus.profile))
		err = dm.Insert(ctx, bus.input.Entity)
		if err != nil {
			break
		}
		body := &Body{
			Data: bus.input.Entity,
		}
		pack = body
	case "update":
		if condition == nil {
			return nil, error2.NewError(ce.InvalidCondition)
		}
		bus.input.Entity = defaultField(bus.input.Entity,
			WithUpdated(bus.profile),
		)
		number, err := dm.Update(ctx,
			builder,
			bus.input.Entity)
		if err != nil {
			break
		}
		ud, err := dm.Find(ctx, builder, bus.input.FindOptions)
		if err != nil {
			break
		}
		res := filterUpdateData(bus, ud)
		body := &Body{
			Number: fail(number, bus.input.Conditions.Condition),
			Data:   res,
		}
		pack = body
	case "delete":
		if condition == nil {
			return nil, error2.NewError(ce.InvalidCondition)
		}
		dd, err := dm.Find(ctx, builder, bus.input.FindOptions)
		if err != nil {
			break
		}
		number, err := dm.Delete(ctx, builder)
		if err != nil {
			break
		}
		body := &Body{
			Number: fail(number, bus.input.Conditions.Condition),
			Data:   dd,
		}
		pack = body
	default:
		return nil, errors.New("invalid method")
	}
	return pack, err
}

// Condition parameter to condition
func (cm *CMongo) Condition(src []*InputConditionVo, inputCondition *InputConditionVo) (clause.Expressions, error) {
	var fc = func(conditions *InputConditionVo) (clause.Expressions, error) {
		exprs, err := dorm.Converts(cm.dc, conditions.Condition...)
		if err != nil {
			return nil, err
		}
		return dorm.Link(cm.dc, dorm.GetOP(conditions.Tag), exprs...)
	}
	var (
		exprInput, exprDefault clause.Expressions
		err                    error
	)
	srcExpr := make([]clause.Expressions, 0)
	if src != nil {
		for _, condition := range src {
			exprDefault, err = fc(condition)
			if err != nil {
				return nil, err
			}
			if exprDefault != nil {
				srcExpr = append(srcExpr, exprDefault)
			}
		}
	}
	if inputCondition != nil && len(inputCondition.Condition) != 0 {
		exprInput, err = fc(inputCondition)
		if err != nil {
			return nil, err
		}
	}
	switch {
	case exprInput == nil && len(srcExpr) == 0:
		return nil, nil
	case exprInput != nil && len(srcExpr) > 0:
		srcExpr = append(srcExpr, exprInput)
		return dorm.Link(cm.dc, dorm.AND, srcExpr...)
	default:
		if exprInput != nil {
			return exprInput, nil
		}
		// 说明长度大于0
		return dorm.Link(cm.dc, dorm.AND, srcExpr...)
	}
}

// 版本2 的crud

// Search Search
func (cm *CMongo) Search(ctx context.Context, bus *bus) (Pack, error) {
	var (
		pack Pack = &Body{}
		err  error
	)
	var dm dao.Dao
	dm = &dao.MONGO{C: cm.DB.Collection(bus.tableName)}

	builder := &clause.MONGO{}
	err = cm.Builders(bus, builder, "find")
	if err != nil {
		return nil, err
	}
	datas, err := dm.Find(ctx, builder, bus.input1.FindOptions)
	if err != nil {
		return nil, err
	}
	paging := &Paging{}
	if len(bus.input1.Aggs) != 0 {
		paging.Aggregations = datas
	} else {
		paging.Data = datas
		paging.Total, err = dm.Count(ctx, builder)
	}
	pack = paging
	return pack, err
}

//Builders Builders
func (cm *CMongo) Builders(bus *bus, builder clause.Builder, method string) error {
	expressions, err := cm.GetExpressions(bus, method)
	if err != nil {
		return err
	}
	if expressions != nil {
		expressions.MongoBuild(builder)
	}
	aggregate, err := cm.GetAggregation(bus.input1.Aggs)
	if err != nil {
		return nil
	}
	if aggregate != nil {
		aggregate.MongoAgg(builder)
	}
	return nil
}

// GetExpressions GetExpressions
func (cm *CMongo) GetExpressions(bus *bus, method string) (clause.Expressions, error) {
	dataAccess := make([]*InputConditionVo, 0)
	if bus.permissionGroup != nil && len(bus.permissionGroup.DataAccessPer) > 0 {
		filterCondition := bus.permissionGroup.DataAccessPer[method]
		if filterCondition != nil {
			dataAccess = append(dataAccess, aggregateCondition(filterCondition))
		}
	}
	expers := make([]clause.Expressions, 0)
	if len(bus.input1.Conditions) != 0 {
		expressions, err := dorm.DslToExper(cm.dc, cm.query, bus.input1.Conditions)
		if err != nil {
			return nil, err
		}
		expers = append(expers, expressions)
	}
	for _, condition := range dataAccess {
		exper, err := dorm.Converts(cm.dc, condition.Condition...)
		if err != nil {
			return nil, err
		}
		expers = append(expers, exper...)
	}

	expressions, err := dorm.Link(cm.dc, dorm.AND, expers...)
	if err != nil {
		return nil, err
	}
	return expressions, nil
}

// GetAggregation GetAggregation
func (cm *CMongo) GetAggregation(value map[string]interface{}) (clause.Aggregations, error) {
	if len(value) == 0 {
		return nil, nil
	}
	agg, err := dorm.DslToAgg(cm.ag, value)

	if err != nil {
		return nil, err
	}
	return agg, nil

}

// Create Create
func (cm *CMongo) Create(ctx context.Context, bus *bus) (Pack, error) {
	var (
		pack Pack = &Body{}
		err  error
	)
	bus.input1.Entity = defaultField(bus.input1.Entity,
		WithID(),
		WithCreated(bus.profile))
	var dm dao.Dao
	dm = &dao.MONGO{C: cm.DB.Collection(bus.tableName)}

	builder := &clause.MONGO{}
	err = cm.Builders(bus, builder, "create")
	if err != nil {
		return nil, err
	}
	err = dm.Insert(ctx, bus.input1.Entity)
	if err != nil {
		return nil, err
	}
	body := &Body{
		Data: bus.input.Entity,
	}
	pack = body
	return pack, err
}

// Update Update
func (cm *CMongo) Update(ctx context.Context, bus *bus) (Pack, error) {
	var (
		pack Pack = &Body{}
		err  error
	)
	bus.input1.Entity = defaultField(bus.input1.Entity,
		WithUpdated(bus.profile),
	)
	var dm dao.Dao
	dm = &dao.MONGO{C: cm.DB.Collection(bus.tableName)}
	builder := &clause.MONGO{}
	err = cm.Builders(bus, builder, "update")
	if err != nil {
		return nil, err
	}
	number, err := dm.Update(ctx,
		builder,
		bus.input.Entity)
	if err != nil {
		return nil, err
	}

	body := &Body{
		Number: number,
	}
	pack = body
	return pack, nil
}

// Delete  Delete
func (cm *CMongo) Delete(ctx context.Context, bus *bus) (Pack, error) {
	var (
		pack Pack = &Body{}
		err  error
	)
	var dm dao.Dao
	dm = &dao.MONGO{C: cm.DB.Collection(bus.tableName)}
	builder := &clause.MONGO{}
	err = cm.Builders(bus, builder, "delete")

	dd, err := dm.Find(ctx, builder, bus.input.FindOptions)
	if err != nil {
		return nil, err
	}
	number, err := dm.Delete(ctx, builder)
	if err != nil {
		return nil, err
	}
	body := &Body{
		Data:   dd,
		Number: number,
	}
	pack = body
	return pack, nil
}
