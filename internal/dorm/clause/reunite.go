package clause

// BETWEEN interval closed at the right
type BETWEEN struct {
	Column string
	Values []interface{}
}

// GetTag get tag
func (b BETWEEN) GetTag() string {
	return "between"
}

// Set set value
func (b BETWEEN) Set(column string, values ...interface{}) Expressions {
	b.Column = column
	b.Values = values
	return b
}

// Build build SQL
func (b BETWEEN) Build(builder Builder) {
	b.build().Build(builder)
}

// MongoBuild build mongo bson
func (b BETWEEN) MongoBuild(builder Builder) {
	b.build().MongoBuild(builder)
}

func (b BETWEEN) build() Expressions {
	var expr Expressions
	switch len(b.Values) {
	case 0:
		expr = EQUAL{}
		expr = expr.Set(b.Column, "NULL")
	case 1:
		expr = GTE{}
		expr = expr.Set(b.Column, b.Values[0])
	default:
		var left, right Expressions
		left = GTE{}
		left = left.Set(b.Column, b.Values[0])
		right = LT{}
		right = right.Set(b.Column, b.Values[1])

		expr = AND{}
		expr = expr.Set("", left, right)
	}
	return expr
}

// Intersection Intersection
type Intersection struct {
	Column string
	Values []interface{}
}

// GetTag get tag
func (intersection Intersection) GetTag() string {
	return "intersection"
}

// Build Build
func (intersection Intersection) Build(builder Builder) {

}

// Set set value
func (intersection Intersection) Set(column string, values ...interface{}) Expressions {
	intersection.Column = column
	intersection.Values = values
	return intersection
}

// MongoBuild MongoBuild
func (intersection Intersection) MongoBuild(builder Builder) {
	builder.WriteString("$in")
	switch len(intersection.Values) {
	case 0:
		builder.AddVar("NULL")
	default:
		builder.AddVar(intersection.Values)
	}
	builder.WriteQuoted(intersection.Column)
}

//FullSubset FullSubset
type FullSubset struct {
	Column string
	Values []interface{}
}

// GetTag get tag
func (fs FullSubset) GetTag() string {
	return "fullSubset"
}

// Build Build
func (fs FullSubset) Build(builder Builder) {

}

// Set set value
func (fs FullSubset) Set(column string, values ...interface{}) Expressions {
	fs.Column = column
	fs.Values = values
	return fs
}

// MongoBuild MongoBuild
func (fs FullSubset) MongoBuild(builder Builder) {
	builder.WriteString("$all")
	switch len(fs.Values) {
	case 0:
		builder.AddVar("NULL")
	default:
		builder.AddVar(fs.Values)
	}

	builder.WriteQuoted(fs.Column)
}
