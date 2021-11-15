package clause

// Expression expression
type Expression interface {
	GetTag() string
	Set(column string, values ...interface{}) Expressions
}

// MYSQLExpression MYSQLExpression expression builder
type MYSQLExpression interface {
	Expression
	Build(builder Builder)
}

// MONGOExpression mongo expression builder
type MONGOExpression interface {
	Expression
	MongoBuild(builder Builder)
}

// Expressions expression set
type Expressions interface {
	Expression
	MYSQLExpression
	MONGOExpression
}
