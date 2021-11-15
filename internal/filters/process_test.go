package filters

import (
	"encoding/json"
	"fmt"
)

func (suite *JSONFilterSuite) TestSchemaLoseWeight() {
	var data = "{\"title\":\"固定资产领用\",\"type\":\"object\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"properties\":{\"_id\":{\"title\":\"\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"string\"},\"bianhao\":{\"title\":\"编号\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"string\"},\"jiluleixing\":{\"title\":\"记录类型\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"string\"},\"bumenfuzeren\":{\"title\":\"部门负责人\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"string\"},\"bumenfuzerenid\":{\"title\":\"部门负责人 ID\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"string\"},\"shenqingshiyou\":{\"title\":\"申请事由\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"string\"},\"shenqingmingxi\":{\"title\":\"申请明细\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"array\",\"item\":{\"title\":\"\",\"type\":\"object\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"properties\":{\"_id\":{\"title\":\"\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"string\"},\"gudingzichanpingcheng\":{\"title\":\"固定资产名称\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"string\"},\"danwei\":{\"title\":\"单位\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"string\"},\"shuliang\":{\"title\":\"数量\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"number\"},\"shiyongdidian\":{\"title\":\"使用地点\",\"x-internal\":{\"sortable\":true,\"permission\":1},\"type\":\"string\"}}}}}}"
	fmt.Println(data)
	schema := make(map[string]interface{})
	json.Unmarshal([]byte(data), &schema)
	m := make(map[string]interface{})
	SchemaLoseWeight(schema, m)
	fmt.Println("===result :",m)

}
