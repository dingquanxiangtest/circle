package clause

import (
	"bytes"
)

// MYSQL mysql
type MYSQL struct {
	SQL  bytes.Buffer
	Vars []interface{}
}

// WriteString write string
func (m *MYSQL) WriteString(str string) (int, error) {
	return m.SQL.WriteString(str)
}

// WriteByte write byte
func (m *MYSQL) WriteByte(c byte) error {
	return m.SQL.WriteByte(c)
}

// WriteQuoted write quoted
func (m *MYSQL) WriteQuoted(field string) {
	m.WriteString(field)
}

// AddVar add var
func (m *MYSQL) AddVar(value interface{}) {
	m.Vars = append(m.Vars, value)
}

// AddAggVar AddAggVar
func (m *MYSQL) AddAggVar(key string, value interface{}) {

}

// WriteQuotedAgg WriteQuotedAgg
func (m *MYSQL) WriteQuotedAgg(field string) {

}
