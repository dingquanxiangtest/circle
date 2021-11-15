package filters

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type JSONFilterSuite struct {
	suite.Suite
}

func TestJSONFilter(t *testing.T) {
	suite.Run(t, new(JSONFilterSuite))
}

type TestObject struct {
	Name string        `json:"name"`
	TEST string        `json:"test"`
	To   TestObject2   `json:"cd"`
	Tos  []TestObject2 `json:"tos"`
}

type TestObject2 struct {
	Name     string `json:"name"`
	TEST     string `json:"test"`
	SelfName string `json:"selfName"`
}

func (suite *JSONFilterSuite) TestFilter2() {

	to2 := make([]TestObject2, 0)
	for i := 0; i < 5; i++ {
		odl := TestObject2{
			Name:     fmt.Sprintf("%dName", i),
			TEST:     fmt.Sprintf("%dtest", i),
			SelfName: fmt.Sprintf("%d selfName", i),
		}
		to2 = append(to2, odl)
	}

	testTos := make([]TestObject, 0)
	for i := 0; i < 5; i++ {
		odl := TestObject{
			Name: fmt.Sprintf("%dName", i),
			TEST: fmt.Sprintf("%dtest", i),
			Tos:  to2,
		}
		testTos = append(testTos, odl)
	}

	bytes, _ := json.Marshal(testTos)
	m := make([]map[string]interface{}, 0)
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		panic(err)
	}

	r2 := make(map[string]interface{})
	r2["name"] = 1
	r21 := make(map[string]interface{})
	r22 := make(map[string]interface{})

	r21["selfName"] = 1
	r2["tos"] = r21
	r22["name"] = 1
	r2["cd"] = r22
	marshal1, _ := json.Marshal(m)
	fmt.Println(string(marshal1))
	JSONFilter2(&m, r2)
	marshal, _ := json.Marshal(m)
	fmt.Println(string(marshal))

}

func (suite *JSONFilterSuite) TestDealSchemaToFilterType() {
	var data = "{\"title\":\"固定资产领用\",\"type\":\"object\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"properties\":{\"_id\":{\"title\":\"\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"string\"},\"bianhao\":{\"title\":\"编号\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"string\"},\"jiluleixing\":{\"title\":\"记录类型\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"string\"},\"bumenfuzeren\":{\"title\":\"部门负责人\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"string\"},\"bumenfuzerenid\":{\"title\":\"部门负责人 ID\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"string\"},\"shenqingshiyou\":{\"title\":\"申请事由\",\"x-internal\":{\"sortable\":true,\"permission\":0},\"type\":\"string\"},\"shenqingmingxi\":{\"title\":\"申请明细\",\"x-internal\":{\"sortable\":true,\"permission\":0},\"type\":\"array\",\"item\":{\"title\":\"\",\"type\":\"object\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"properties\":{\"_id\":{\"title\":\"\",\"x-internal\":{\"sortable\":true,\"permission\":2},\"type\":\"string\"},\"gudingzichanpingcheng\":{\"title\":\"固定资产名称\",\"x-internal\":{\"sortable\":true,\"permission\":0},\"type\":\"string\"},\"danwei\":{\"title\":\"单位\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"string\"},\"shuliang\":{\"title\":\"数量\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"number\"},\"shiyongdidian\":{\"title\":\"使用地点\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"string\"}}}}}}"
	schema := Schema{}
	err := json.Unmarshal([]byte(data), &schema)
	if err != nil {
		panic(err)
	}
	resMap := DealSchemaToFilterType(schema)
	marshal, _ := json.Marshal(resMap)
	fmt.Println(string(marshal))
	suite.NotNil(resMap)
}

func (suite *JSONFilterSuite) TestSchemaFilter2() {
	var data = "{\"title\":\"固定资产领用\",\"type\":\"object\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"properties\":{\"_id\":{\"title\":\"\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"string\"},\"bianhao\":{\"title\":\"编号\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"string\"},\"jiluleixing\":{\"title\":\"记录类型\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"string\"},\"bumenfuzeren\":{\"title\":\"部门负责人\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"string\"},\"bumenfuzerenid\":{\"title\":\"部门负责人 ID\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"string\"},\"shenqingshiyou\":{\"title\":\"申请事由\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"string\"},\"shenqingmingxi\":{\"title\":\"申请明细\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"array\",\"item\":{\"title\":\"\",\"type\":\"object\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"properties\":{\"_id\":{\"title\":\"\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"string\"},\"gudingzichanpingcheng\":{\"title\":\"固定资产名称\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"string\"},\"danwei\":{\"title\":\"单位\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"string\"},\"shuliang\":{\"title\":\"数量\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"number\"},\"shiyongdidian\":{\"title\":\"使用地点\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"string\"}}}}}}"
	var filter = "{\"_id\":1,\"bianhao\":2,\"bumenfuzeren\":3,\"bumenfuzerenid\":1,\"jiluleixing\":1}"
	//var filter = "{shenqingmingxi\":{\"_id\":1,\"danwei\":1},\"shenqingshiyou\":1}}"
	//var filter = "{\"shenqingmingxi\":{\"_id\":1,\"danwei\":1,\"gudingzichanpingcheng\":1,\"shiyongdidian\":1,\"shuliang\":1},\"shenqingshiyou\":1}"
	//var filter = "{\"shenqingmingxi\":{\"danwei\":3,\"shuliang\":4}}"

	m2 := make(map[string]interface{})
	json.Unmarshal([]byte(data), &m2)
	m := make(map[string]interface{})
	json.Unmarshal([]byte(filter), &m)
	SchemaFilterToNewSchema2(m2, m)

}

type Entity interface {
}

func (suite *JSONFilterSuite) TestCheckData() {
	var data = "{\"_id\":\"123123\",\"name1\":\"test1\",\"data1\":[{\"name2\":\"123\",\"test2\":\"123\"}],\"test3\":\"test\"}"
	f := make(map[string]interface{})
	f["name1"] = 2
	f["test1"] = 2
	f["data1"] = 2
	f["name2"] = 2
	f["test2"] = 2
	f["test3"] = 2

	var a Entity
	err := json.Unmarshal([]byte(data), &a)
	suite.Nil(err)
	flag := FilterCheckData(a, f)
	suite.True(flag)
}
