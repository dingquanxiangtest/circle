package models

// Condition condition
type Condition struct {
	Key   string        `json:"key"`
	OP    string        `json:"op"`
	Value []interface{} `json:"value"`
}

// ScopesVO ScopesVO
type ScopesVO struct {
	Type int16  `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ConditionVO ConditionVO
type ConditionVO struct {
	Arr []*Condition `json:"arr"`
	Tag string       `json:"tag"`
}
