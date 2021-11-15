package clause

import "go.mongodb.org/mongo-driver/bson"

type common struct {
	alias string
	value interface{}
}

func (c *common) CSetFiled(alias string, value interface{}) {
	c.alias = alias
	c.value = value
}

func (c *common) CMongoAgg(builder Builder, tag string) {
	value, ok := c.value.(map[string]interface{})
	if ok {
		builder.AddAggVar(c.alias, bson.M{"$" + tag: "$" + value["field"].(string)})
	}

}

// Agg Agg
func (c *common) Agg(builder Builder, tag string) {

}

// Sum Sum
type Sum struct {
	common
}

// GetTag GetTag
func (sum *Sum) GetTag() string {
	return "sum"
}

// SetField SetField
func (sum *Sum) SetField(alias string, values ...interface{}) {
	var value interface{}
	if len(values) > 0 {
		value = values[0]
	}
	sum.CSetFiled(alias, value)
}

// MongoAgg MongoAgg
func (sum *Sum) MongoAgg(builder Builder) {
	sum.CMongoAgg(builder, "sum")
}

// Agg Agg
func (sum *Sum) Agg(builder Builder) {

}

// Avg Avg
type Avg struct {
	common
}

// GetTag GetTag
func (avg *Avg) GetTag() string {
	return "avg"
}

// SetField SetField
func (avg *Avg) SetField(alias string, values ...interface{}) {
	var value interface{}
	if len(values) > 0 {
		value = values[0]
	}
	avg.CSetFiled(alias, value)
}

// MongoAgg MongoAgg
func (avg *Avg) MongoAgg(builder Builder) {
	avg.CMongoAgg(builder, "avg")
}

// Agg Agg
func (avg *Avg) Agg(builder Builder) {

}

// Min Min
type Min struct {
	common
}

// GetTag GetTag
func (min *Min) GetTag() string {
	return "min"
}

// SetField SetField
func (min *Min) SetField(alias string, values ...interface{}) {
	var value interface{}
	if len(values) > 0 {
		value = values[0]
	}
	min.CSetFiled(alias, value)
}

// MongoAgg MongoAgg
func (min *Min) MongoAgg(builder Builder) {
	min.CMongoAgg(builder, "min")
}

// Agg Agg
func (min *Min) Agg(builder Builder) {

}

// Max max
type Max struct {
	common
}

// GetTag GetTag
func (max *Max) GetTag() string {
	return "max"
}

// SetField SetField
func (max *Max) SetField(alias string, values ...interface{}) {
	var value interface{}
	if len(values) > 0 {
		value = values[0]
	}
	max.CSetFiled(alias, value)
}

// MongoAgg MongoAgg
func (max *Max) MongoAgg(builder Builder) {
	max.CMongoAgg(builder, "max")
}

// Agg Agg
func (max *Max) Agg(builder Builder) {

}

// Count max
type Count struct {
	common
}

// GetTag GetTag
func (c *Count) GetTag() string {
	return "count"
}

// SetField SetField
func (c *Count) SetField(alias string, values ...interface{}) {
	var value interface{}
	if len(values) > 0 {
		value = values[0]
	}
	c.CSetFiled(alias, value)
}

// MongoAgg MongoAgg
func (c *Count) MongoAgg(builder Builder) {
	builder.AddAggVar("$count", c.alias)
}

// Agg Agg
func (c *Count) Agg(builder Builder) {

}

// Groups Groups
type Groups struct {
	value []interface{}
}

// GetTag GetTag
func (group *Groups) GetTag() string {
	return "group"
}

// SetField SetField
func (group *Groups) SetField(alias string, fieldName ...interface{}) {
	group.value = fieldName
}

//MongoAgg MongoAgg
func (group *Groups) MongoAgg(builder Builder) {
	for _, value := range group.value {
		switch value := value.(type) {
		case MONGOAggregation:
			value.MongoAgg(builder)
		default:
			continue
		}

	}
	builder.WriteQuotedAgg("$group")
}

// Agg Agg
func (group *Groups) Agg(builder Builder) {

}
