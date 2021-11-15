package dorm

import (
	"encoding/json"
	"fmt"
	"testing"

	"git.internal.yunify.com/qxp/molecule/internal/dorm/clause"
)

// TestConvert TestConvert
func TestConvert(t *testing.T) {
	data := "{\"query\":{\"match\":{\"xxx\":\"testName\"},\"term\":{\"xxx\":\"testName\"},\"bool\":{\"must\":[{\"match\":{\"name\":\"testName\"}},{\"term\":{\"name\":\"testName\"}},{\"bool\":{\"must\":[{\"term\":{\"name\":\"xxx\"}},{\"match\":{\"name\":\"xxx\"}}]}}]}}}"
	convert := make(map[string]interface{})
	json.Unmarshal([]byte(data), &convert)
	query := NewQuery()
	clauses := clause.New()
	c := convert["query"]
	boolsss, ok := c.(map[string]interface{})
	if ok {
		exprs := make([]clause.Expressions, 0)
		for key, value := range boolsss {
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

		}
		fmt.Println(link)

	}

}
