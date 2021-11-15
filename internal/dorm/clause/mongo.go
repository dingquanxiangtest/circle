package clause

import (
	"go.mongodb.org/mongo-driver/bson"
)

// MONGO mongo
type MONGO struct {
	Vars bson.M
	Agg  bson.M
}

// WriteString write string
func (m *MONGO) WriteString(str string) (int, error) {
	m.WriteQuoted(str)
	return len(str), nil
}

// WriteByte write byte
func (m *MONGO) WriteByte(c byte) error {
	_, err := m.WriteString(string(c))
	return err
}

// WriteQuoted write quoted
func (m *MONGO) WriteQuoted(field string) {
	m.Vars = bson.M{
		field: m.Vars,
	}
}

// AddVar add var
func (m *MONGO) AddVar(value interface{}) {
	for key := range m.Vars {
		m.Vars[key] = value
		return
	}
}

// WriteQuotedAgg WriteQuotedAgg
func (m *MONGO) WriteQuotedAgg(field string) {
	m.Agg = bson.M{
		field: m.Agg,
	}

}

// AddAggVar AddAggVar
func (m *MONGO) AddAggVar(key string, value interface{}) {
	if m.Agg == nil {
		m.Agg = bson.M{}
		if key != "$count" {
			m.Agg["_id"] = "null"
		}
	}
	m.Agg[key] = value
}
