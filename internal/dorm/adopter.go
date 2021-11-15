package dorm

import "git.internal.yunify.com/qxp/molecule/internal/dorm/clause"

// Entity entity
type Entity interface{}

// Condition condition from comet
type Condition struct {
	// Key column key
	Key string
	// Op in/like/eq/range
	Op string
	// Value condition value
	Value []interface{}
}

// OP database and | or
type OP string

const (
	// AND and
	AND OP = "and"
	// OR or
	OR OP = "or"
)

// GetOP GetOP
func GetOP(tag string) OP {
	if tag == "or" {
		return OR
	}
	return AND
}

// Converts convert conditions to expressions
func Converts(c *clause.Clause, conditions ...Condition) ([]clause.Expressions, error) {
	exprs := make([]clause.Expressions, 0, len(conditions))
	for _, cond := range conditions {
		expr, err := Convert(c, cond)
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, expr)
	}
	return exprs, nil
}

// Convert convert condition to expression
func Convert(c *clause.Clause, cond Condition) (clause.Expressions, error) {
	return c.GetExpression(cond.Op, cond.Key, cond.Value...)
}

// Link link expressions this op
func Link(c *clause.Clause, op OP, exprs ...clause.Expressions) (clause.Expressions, error) {
	if len(exprs) == 0 {
		return nil, nil
	}
	if len(exprs) == 1 {
		return exprs[0], nil
	}
	exprInterface := make([]interface{}, 0, len(exprs))
	for _, expr := range exprs {
		exprInterface = append(exprInterface, expr)
	}
	return c.GetExpression(
		string(op),
		"",
		exprInterface...,
	)
}

// DslToExper DslToExper
func DslToExper(clauses *clause.Clause, query *Query, condition map[string]interface{}) (clause.Expressions, error) {
	exprs := make([]clause.Expressions, 0)
	for key, value := range condition {
		expression, err := query.GetItem(key)
		if err != nil {
			break
		}
		expr, err := expression.ConvertExper(value, clauses, query)
		if err != nil {

		}
		exprs = append(exprs, expr)
	}
	link, err := Link(clauses, AND, exprs...)
	if err != nil {
		return nil, err
	}
	return link, nil
}

// DslToAgg DslToAgg
func DslToAgg(ag *clause.Aggregate, aggregations map[string]interface{}) (clause.Aggregations, error) {
	aggregates := make([]interface{}, 0)
	for alias, value := range aggregations {
		value, ok := value.(map[string]interface{})
		if !ok {
			return nil, nil
		}
		// funcs  sum avg
		for funcs, filed := range value {
			aggregate, err := ag.GetAg(funcs, alias, filed)
			if funcs == "count" {
				return aggregate, nil
			}
			if err != nil {
				return nil, err
			}
			aggregates = append(aggregates, aggregate)
		}
	}
	if len(aggregates) == 0 {
		return nil, nil
	}
	group, err := ag.GetAg("group", "", aggregates...)
	if err != nil {
		return nil, nil
	}
	return group, nil

}
