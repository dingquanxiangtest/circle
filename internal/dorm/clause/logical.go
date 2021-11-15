package clause

type logical struct {
	Vars []interface{}
}

// LSet set value
func (l *logical) LSet(values ...interface{}) {
	l.Vars = values
}

// LBuild build SQL
func (l *logical) LBuild(builder Builder, tag string) {
	var i int
	builder.WriteByte('(')
	defer builder.WriteByte(')')
	for _, vars := range l.Vars {
		switch vars := vars.(type) {
		case MYSQLExpression:
			if i > 0 {
				builder.WriteByte(' ')
				builder.WriteString(tag)
				builder.WriteByte(' ')
			}
			vars.Build(builder)
			i++
		default:
			continue
		}
	}
}

// LMongoBuild build mongo bson
func (l *logical) LMongoBuild(builder Builder, tag string) {
	vars := make([]interface{}, 0, len(l.Vars))
	for _, value := range l.Vars {
		switch value := value.(type) {
		case MONGOExpression:
			value.MongoBuild(builder)
			vars = append(vars, builder.(*MONGO).Vars)
		default:
			continue
		}
	}
	builder.WriteString("$" + tag)
	builder.AddVar(vars)
}

// AND and
type AND struct {
	logical
}

// GetTag get tag
func (and AND) GetTag() string {
	return "and"
}

// Set set value
func (and AND) Set(column string, values ...interface{}) Expressions {
	and.LSet(values...)
	return and
}

// Build build SQL
func (and AND) Build(builder Builder) {
	and.LBuild(builder, and.GetTag())
}

// MongoBuild build mongo bson
func (and AND) MongoBuild(builder Builder) {
	and.LMongoBuild(builder, and.GetTag())
}

// OR or
type OR struct {
	logical
}

// GetTag get tag
func (or OR) GetTag() string {
	return "or"
}

// Set set value
func (or OR) Set(column string, values ...interface{}) Expressions {
	or.LSet(values...)
	return or
}

// Build build SQL
func (or OR) Build(builder Builder) {
	or.LBuild(builder, or.GetTag())
}

// MongoBuild build mongo bson
func (or OR) MongoBuild(builder Builder) {
	or.LMongoBuild(builder, or.GetTag())
}
