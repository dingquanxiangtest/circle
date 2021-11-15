package clause

// IN Whether a value is within a set of values
type IN struct {
	Column string
	Values []interface{}
}

// GetTag get tag
func (in IN) GetTag() string {
	return "in"
}

// Set set value
func (in IN) Set(column string, values ...interface{}) Expressions {
	in.Column = column
	in.Values = values
	return in
}

// Build build SQL
func (in IN) Build(builder Builder) {
	builder.WriteQuoted(in.Column)

	switch len(in.Values) {
	case 0:
		builder.WriteString(" IN (NULL)")
	case 1:
		if _, ok := in.Values[0].([]interface{}); !ok {
			builder.WriteString(" = ?")
			builder.AddVar(in.Values[0])
			break
		}
	default:
		builder.WriteString(" IN (?)")
		builder.AddVar(in.Values)
	}
}

// MongoBuild build mongo bson
func (in IN) MongoBuild(builder Builder) {
	builder.WriteString("$in")
	switch len(in.Values) {
	case 0:
		builder.AddVar("NULL")
	default:
		builder.AddVar(in.Values)
	}

	builder.WriteQuoted(in.Column)
}

// LIKE fuzzy
type LIKE struct {
	Column string
	Values interface{}
}

// GetTag get tag
func (like LIKE) GetTag() string {
	return "like"
}

// Set set value
func (like LIKE) Set(column string, values ...interface{}) Expressions {
	var checkString = func(value interface{}) bool {
		_, ok := value.(string)
		return ok
	}
	like.Column = column
	if len(values) != 1 || !checkString(values[0]) {
		like.Values = "NULL"
	} else {
		like.Values = values[0]
	}

	return like
}

// Build build SQL
func (like LIKE) Build(builder Builder) {
	builder.WriteQuoted(like.Column)

	builder.WriteString(" LIKE ?")
	builder.AddVar("%" + like.Values.(string) + "%")
}

// MongoBuild build mongo bson
func (like LIKE) MongoBuild(builder Builder) {
	builder.WriteString("$regex")
	builder.AddVar(like.Values)
	builder.WriteQuoted(like.Column)
}

// EQUAL equal
type EQUAL struct {
	Column string
	Values []interface{}
}

// GetTag get tag
func (equal EQUAL) GetTag() string {
	return "eq"
}

// Set set value
func (equal EQUAL) Set(column string, values ...interface{}) Expressions {
	equal.Column = column
	equal.Values = values
	return equal
}

// Build build SQL
func (equal EQUAL) Build(builder Builder) {
	switch len(equal.Values) {
	case 0:
		equal.Values = append(equal.Values, "NULL")
	case 1:
		builder.WriteQuoted(equal.Column)
		builder.WriteString(" = ?")
		builder.AddVar(equal.Values[0])
	default:
		exprs := make([]interface{}, 0, len(equal.Values))
		for _, value := range equal.Values {
			var expr Expressions
			expr = LIKE{}
			expr = expr.Set(equal.Column, value)
			exprs = append(exprs, expr)
		}
		var or Expressions
		or = OR{}
		or = or.Set("", exprs...)
		or.Build(builder)
	}
}

// MongoBuild build mongo bson
func (equal EQUAL) MongoBuild(builder Builder) {
	builder.WriteQuoted(equal.Column)
	switch len(equal.Values) {
	case 0:
		equal.Values = append(equal.Values, "NULL")
	case 1:
		builder.AddVar(equal.Values[0])
	default:
		builder.AddVar(equal.Values)
	}
}

type conditionOP struct {
	Column string
	Values interface{}
}

// CSet set value
func (c *conditionOP) CSet(column string, values interface{}) {
	c.Column = column
	if values == nil {
		values = "NULL"
	}
	c.Values = values
}

// CBuild build SQL
func (c *conditionOP) CBuild(builder Builder, tag string) {
	builder.WriteQuoted(c.Column)

	builder.WriteByte(' ')
	builder.WriteString(tag)
	builder.WriteByte(' ')
	builder.WriteByte('?')
	builder.AddVar(c.Values)
}

// CMongoBuild build mongo bson
func (c *conditionOP) CMongoBuild(builder Builder, tag string) {
	builder.WriteString("$" + tag)
	builder.AddVar(c.Values)
	builder.WriteQuoted(c.Column)
}

// LT less than
type LT struct {
	conditionOP
}

// GetTag get tag
func (lt LT) GetTag() string {
	return "lt"
}

// Set set value
func (lt LT) Set(column string, values ...interface{}) Expressions {
	var value interface{}
	if len(values) > 0 {
		value = values[0]
	}
	lt.CSet(column, value)
	return lt
}

// Build build SQL
func (lt LT) Build(builder Builder) {
	lt.CBuild(builder, "<")
}

// MongoBuild build mongo bson
func (lt LT) MongoBuild(builder Builder) {
	lt.CMongoBuild(builder, "lt")
}

// LTE less than or equal
type LTE struct {
	conditionOP
}

// GetTag get tag
func (lte LTE) GetTag() string {
	return "lte"
}

// Set set value
func (lte LTE) Set(column string, values ...interface{}) Expressions {
	var value interface{}
	if len(values) > 0 {
		value = values[0]
	}
	lte.CSet(column, value)
	return lte
}

// Build build SQL
func (lte LTE) Build(builder Builder) {
	lte.CBuild(builder, "<=")
}

// MongoBuild build mongo bson
func (lte LTE) MongoBuild(builder Builder) {
	lte.CMongoBuild(builder, "lte")
}

// GT greater than
type GT struct {
	conditionOP
}

// GetTag get tag
func (gt GT) GetTag() string {
	return "gt"
}

// Set set value
func (gt GT) Set(column string, values ...interface{}) Expressions {
	var value interface{}
	if len(values) > 0 {
		value = values[0]
	}
	gt.CSet(column, value)
	return gt
}

// Build build SQL
func (gt GT) Build(builder Builder) {
	gt.CBuild(builder, ">")
}

// MongoBuild build mongo bson
func (gt GT) MongoBuild(builder Builder) {
	gt.CMongoBuild(builder, "gt")
}

// GTE greater than or equal
type GTE struct {
	conditionOP
}

// GetTag get tag
func (gte GTE) GetTag() string {
	return "gte"
}

// Set set value
func (gte GTE) Set(column string, values ...interface{}) Expressions {
	var value interface{}
	if len(values) > 0 {
		value = values[0]
	}
	gte.CSet(column, value)
	return gte
}

// Build build SQL
func (gte GTE) Build(builder Builder) {
	gte.CBuild(builder, ">=")
}

// MongoBuild build mongo bson
func (gte GTE) MongoBuild(builder Builder) {
	gte.CMongoBuild(builder, "gte")
}
