package clause

import "errors"

// Aggregation Aggregation
type Aggregation interface {
	GetTag() string
	SetField(alias string, value ...interface{})
}

// MYSQLAggregation MYSQLAggregation
type MYSQLAggregation interface {
	Aggregation
	Agg(builder Builder)
}

// MONGOAggregation MONGOAggregation
type MONGOAggregation interface {
	Aggregation
	MongoAgg(builder Builder)
}

// Aggregations Aggregations
type Aggregations interface {
	MYSQLAggregation
	MONGOAggregation
}

var (
	// ErrNoAggregation no expression
	ErrNoAggregation = errors.New("no Aggregation like this")
)

var ags = []Aggregations{
	&Sum{},
	&Groups{},
	&Avg{},
	&Min{},
	&Max{},
	&Count{},
}

// Aggregate Aggregate
type Aggregate struct {
	aggregations map[string]Aggregations
}

// NewAg new Ag
func NewAg() *Aggregate {
	aggregate := &Aggregate{
		aggregations: make(map[string]Aggregations, len(ags)),
	}
	for _, ag := range ags {
		aggregate.aggregations[ag.GetTag()] = ag
	}
	return aggregate
}

// GetAg GetAg
func (a *Aggregate) GetAg(op, alias string, value ...interface{}) (Aggregations, error) {
	ag, ok := a.aggregations[op]
	if !ok {
		return nil, ErrNoAggregation
	}
	ag.SetField(alias, value...)
	return ag, nil
}
