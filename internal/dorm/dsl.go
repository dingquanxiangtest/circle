package dorm

import (
	"errors"
	"git.internal.yunify.com/qxp/molecule/internal/dorm/clause"
	"reflect"
)

// DSL DSL
type DSL interface {
	ConvertExper(convert interface{}, clause *clause.Clause, query *Query) (clause.Expressions, error)
	GetTag() string
}

// Simple simple  is  match/term/terms
type Simple struct {
}

// SConvertExper SConvertExper
func (s *Simple) SConvertExper(convert interface{}, clause *clause.Clause, op string) (clause.Expressions, error) {
	c, ok := convert.(map[string]interface{})
	if !ok {
		return nil, ErrConvert
	}
	if len(c) != 1 {
		return nil, ErrConvert
	}
	for key, value := range c {
		return clause.GetExpression(op, key, value)

	}
	return nil, ErrConvert
}

var converts = []DSL{
	&Match{},
	&Term{},
	&Terms{},
	&Bool{},
	&Must{},
	&MustNot{},
	&Should{},
	&Range{},
}

// Query Query
type Query struct {
	converts map[string]DSL
}

// NewQuery new Query
func NewQuery() *Query {
	c := &Query{
		converts: make(map[string]DSL, len(converts)),
	}
	for _, converts := range converts {
		c.converts[converts.GetTag()] = converts
	}
	return c
}

// Complex must must_not should
type Complex struct {
}

//CPConvert CPConvert
func (cp *Complex) CPConvert(convert interface{}, clauses *clause.Clause, query *Query) ([]clause.Expressions, error) {
	c, ok := convert.([]interface{})
	if !ok {
		return nil, ErrConvert
	}
	expers := make([]clause.Expressions, 0)
	for _, value := range c {
		v1, ok := value.(map[string]interface{})
		if !ok {
			return nil, ErrConvert
		}
		for key, v2 := range v1 {
			dslItem, err := query.GetItem(key)
			if err != nil {
				return nil, err
			}
			exper, err := dslItem.ConvertExper(v2, clauses, query)
			if err != nil {
				return nil, err
			}
			expers = append(expers, exper)
		}
	}
	return expers, nil
}

// GetTag GetTag
func (cp *Complex) GetTag() string {
	return "complex"
}

// Bool Bool
type Bool struct {
}

// GetTag GetTag
func (b *Bool) GetTag() string {
	return "bool"
}

// ConvertExper ConvertExper
func (b *Bool) ConvertExper(convert interface{}, clauses *clause.Clause, query *Query) (clause.Expressions, error) {
	c, ok := convert.(map[string]interface{})
	if !ok {
		return nil, ErrConvert
	}
	exprs := make([]clause.Expressions, 0)
	for key, value := range c { // 遍历bool
		expression, err := query.GetItem(key)
		if err != nil {
			return nil, err
		}
		exper, err := expression.ConvertExper(value, clauses, query)

		exprs = append(exprs, exper)

	}
	link, err := Link(clauses, AND, exprs...)
	if err != nil {
		return nil, err
	}
	return link, nil
}

var (
	// ErrNoDSLItem no ErrNoDSLItem
	ErrNoDSLItem = errors.New("no ErrNoDSLItem like this")
	// ErrConvert ErrConvert
	ErrConvert = errors.New("convert  fail  ")
)

// GetItem get item with tag
func (query *Query) GetItem(tag string) (DSL, error) {
	convert, ok := query.converts[tag]
	if !ok {
		return nil, ErrNoDSLItem
	}
	return convert, nil
}

// Match Match
type Match struct {
	Simple
}

// ConvertExper ConvertExper
func (m *Match) ConvertExper(convert interface{}, clause *clause.Clause, query *Query) (clause.Expressions, error) {
	exper, err := m.SConvertExper(convert, clause, "like")
	if err != nil {
		return nil, err
	}
	return exper, nil
}

// GetTag GetTag
func (m *Match) GetTag() string {
	return "match"
}

// Term Term
type Term struct {
	Simple
}

// ConvertExper ConvertExper
func (t *Term) ConvertExper(convert interface{}, clause *clause.Clause, query *Query) (clause.Expressions, error) {
	exper, err := t.SConvertExper(convert, clause, "eq")
	if err != nil {
		return nil, err
	}
	return exper, nil
}

// GetTag GetTag
func (t *Term) GetTag() string {
	return "term"
}

//Terms Terms
type Terms struct {
	Simple
}

// ConvertExper ConvertExper
func (t *Terms) ConvertExper(convert interface{}, clause *clause.Clause, query *Query) (clause.Expressions, error) {
	c, ok := convert.(map[string]interface{})
	if !ok {
		return nil, ErrConvert
	}
	if len(c) != 1 {
		return nil, ErrConvert
	}
	for key, value := range c {
		if val := reflect.ValueOf(value); val.CanInterface() {
			if v := val.Interface().([]interface{}); ok {
				return clause.GetExpression("in", key, v...)
			}
		}
		return clause.GetExpression("in", key, value)
	}
	return nil, ErrConvert
}

// GetTag GetTag
func (t *Terms) GetTag() string {
	return "terms"
}

// Must Must
type Must struct {
	Complex
}

// ConvertExper ConvertExper
func (m *Must) ConvertExper(convert interface{}, clause *clause.Clause, query *Query) (clause.Expressions, error) {
	expers, err := m.CPConvert(convert, clause, query)
	if err != nil {
		return nil, err
	}
	exper, err := Link(clause, AND, expers...)
	if err != nil {
		return nil, err
	}
	return exper, nil
}

// GetTag GetTag
func (m *Must) GetTag() string {
	return "must"
}

// MustNot MustNot
type MustNot struct {
	Complex
}

// ConvertExper ConvertExper
func (m *MustNot) ConvertExper(convert interface{}, clause *clause.Clause, query *Query) (clause.Expressions, error) {
	_, err := m.CPConvert(convert, clause, query)
	if err != nil {
		return nil, err
	}
	// TODO 拿到exper 数组，取反
	return nil, nil
}

// GetTag GetTag
func (m *MustNot) GetTag() string {
	return "must_not"
}

// Should Should
type Should struct {
	Complex
}

// ConvertExper ConvertExper
func (s *Should) ConvertExper(convert interface{}, clause *clause.Clause, query *Query) (clause.Expressions, error) {
	expers, err := s.CPConvert(convert, clause, query)
	if err != nil {
		return nil, err
	}
	exper, err := Link(clause, OR, expers...)
	if err != nil {
		return nil, err
	}
	return exper, nil
}

// GetTag GetTag
func (s *Should) GetTag() string {
	return "should"
}

// Range Range
type Range struct {
	Inequality
}

// ConvertExper ConvertExper
func (r *Range) ConvertExper(convert interface{}, clauses *clause.Clause, query *Query) (clause.Expressions, error) {
	ranges, ok := convert.(map[string]interface{})
	if !ok {
		return nil, ErrConvert
	}
	expers := make([]clause.Expressions, 0)
	for field, inequalitys := range ranges {
		iqs := inequalitys.(map[string]interface{})
		for op, value := range iqs {
			exper, err := r.InConvertExper(op, field, value, clauses)
			if err != nil {
				return nil, err
			}
			expers = append(expers, exper)
		}
	}
	exper, err := Link(clauses, AND, expers...)
	if err != nil {
		return nil, err
	}
	return exper, nil
}

// GetTag GetTag
func (r *Range) GetTag() string {
	return "range"
}

// Inequality < <= > >=
type Inequality struct {
}

// InConvertExper InConvertExper
func (iq *Inequality) InConvertExper(op, key string, value interface{}, clause *clause.Clause) (clause.Expressions, error) {
	return clause.GetExpression(op, key, value)
}
