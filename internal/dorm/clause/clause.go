package clause

import (
	"errors"
)

// Writer write
type Writer interface {
	WriteByte(byte) error
	WriteString(string) (int, error)
}

// QueryBuilder db condition builder
type QueryBuilder interface {
	Writer
	WriteQuoted(field string)
	AddVar(interface{})
}

// AggBuilder AggBuilder
type AggBuilder interface {
	WriteQuotedAgg(field string)
	AddAggVar(key string, value interface{})
}

//Builder Builder
type Builder interface {
	QueryBuilder
	AggBuilder
}

var (
	// ErrNoExpression no expression
	ErrNoExpression = errors.New("no expression like this")
)

var expressions = []Expressions{
	IN{},
	LIKE{},
	EQUAL{},
	LT{},
	LTE{},
	GT{},
	GTE{},

	AND{},
	OR{},

	BETWEEN{},

	FullSubset{},
	Intersection{},
}

// Clause expressions set
type Clause struct {
	Expressions map[string]Expressions
}

// New new a clause
func New() *Clause {
	c := &Clause{
		Expressions: make(map[string]Expressions, len(expressions)),
	}
	for _, expr := range expressions {
		c.Expressions[expr.GetTag()] = expr
	}

	return c
}

// GetExpression get expression with op
func (c *Clause) GetExpression(op string, column string, values ...interface{}) (Expressions, error) {
	expr, ok := c.Expressions[op]
	if !ok {
		return nil, ErrNoExpression
	}
	expr = expr.Set(column, values...)
	return expr, nil
}
