package swagger

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

// TestConvert TestConvert
func TestConvert(t *testing.T) {
	data := "{\"_id\":{\"display\":false,\"readOnly\":false,\"title\":\"id\",\"type\":\"string\",\"x-component\":\"Input\",\"x-component-props\":{},\"x-index\":6,\"x-internal\":{\"isSystem\":true,\"permission\":3},\"x-mega-props\":{\"labelCol\":4}},\"creator_id\":{\"display\":false,\"readOnly\":false,\"title\":\"创建者 ID\",\"type\":\"string\",\"x-component\":\"Input\",\"x-component-props\":{},\"x-index\":10,\"x-internal\":{\"isSystem\":true,\"permission\":3},\"x-mega-props\":{\"labelCol\":4}},\"creator_name\":{\"display\":false,\"readOnly\":false,\"title\":\"创建者\",\"type\":\"string\",\"x-component\":\"Input\",\"x-component-props\":{},\"x-index\":9,\"x-internal\":{\"isSystem\":true,\"permission\":3},\"x-mega-props\":{\"labelCol\":4}},\"field_AOJsqDBt\":{\"description\":\"\",\"display\":true,\"format\":\"mobile_phone\",\"readOnly\":false,\"required\":true,\"title\":\"手机号码\",\"type\":\"string\",\"x-component\":\"Input\",\"x-component-props\":{\"defaultValue\":\"\",\"placeholder\":\"请输入\"},\"x-index\":4,\"x-internal\":{\"defaultValueFrom\":\"customized\",\"isSystem\":false,\"permission\":3,\"sortable\":false},\"x-mega-props\":{\"labelCol\":4}},\"field_G0cx2SeQ\":{\"description\":\"\",\"display\":true,\"format\":\"\",\"readOnly\":false,\"required\":true,\"title\":\"负责人\",\"type\":\"string\",\"x-component\":\"Input\",\"x-component-props\":{\"defaultValue\":\"\",\"placeholder\":\"请输入\"},\"x-index\":1,\"x-internal\":{\"defaultValueFrom\":\"customized\",\"isSystem\":false,\"permission\":3,\"sortable\":false},\"x-mega-props\":{\"labelCol\":4}},\"field_aX3VVj8x\":{\"description\":\"\",\"display\":true,\"format\":\"\",\"readOnly\":false,\"required\":true,\"title\":\"地址\",\"type\":\"string\",\"x-component\":\"Input\",\"x-component-props\":{\"defaultValue\":\"\",\"placeholder\":\"请输入\"},\"x-index\":3,\"x-internal\":{\"defaultValueFrom\":\"customized\",\"isSystem\":false,\"permission\":3,\"sortable\":false},\"x-mega-props\":{\"labelCol\":4}},\"field_brNbwe5Y\":{\"description\":\"\",\"display\":true,\"format\":\"email\",\"readOnly\":false,\"required\":true,\"title\":\"邮箱\",\"type\":\"string\",\"x-component\":\"Input\",\"x-component-props\":{\"defaultValue\":\"\",\"placeholder\":\"请输入\"},\"x-index\":2,\"x-internal\":{\"defaultValueFrom\":\"customized\",\"isSystem\":false,\"permission\":3,\"sortable\":false},\"x-mega-props\":{\"labelCol\":4}},\"field_cfj7AP1h\":{\"description\":\"\",\"display\":true,\"minimum\":0,\"readOnly\":false,\"required\":false,\"title\":\"承接项目数\",\"type\":\"number\",\"x-component\":\"NumberPicker\",\"x-component-props\":{\"placeholder\":\"请输入\",\"precision\":0,\"step\":1},\"x-index\":5,\"x-internal\":{\"defaultValueFrom\":\"customized\",\"isSystem\":false,\"permission\":3,\"sortable\":false},\"x-mega-props\":{\"labelCol\":4}},\"field_juSTcHjM\":{\"description\":\"\",\"display\":true,\"format\":\"\",\"readOnly\":false,\"required\":true,\"title\":\"供货商名称\",\"type\":\"string\",\"x-component\":\"Input\",\"x-component-props\":{\"defaultValue\":\"\",\"placeholder\":\"请输入\"},\"x-index\":0,\"x-internal\":{\"defaultValueFrom\":\"customized\",\"isSystem\":false,\"permission\":3,\"sortable\":false},\"x-mega-props\":{\"labelCol\":4}},\"modifier_id\":{\"display\":false,\"readOnly\":false,\"title\":\"修改者 ID\",\"type\":\"string\",\"x-component\":\"Input\",\"x-component-props\":{},\"x-index\":12,\"x-internal\":{\"isSystem\":true,\"permission\":3},\"x-mega-props\":{\"labelCol\":4}},\"modifier_name\":{\"display\":false,\"readOnly\":false,\"title\":\"修改者\",\"type\":\"string\",\"x-component\":\"Input\",\"x-component-props\":{},\"x-index\":11,\"x-internal\":{\"isSystem\":true,\"permission\":3},\"x-mega-props\":{\"labelCol\":4}}}"
	convert1 := make(map[string]interface{})
	err := json.Unmarshal([]byte(data), &convert1)
	if err != nil {

	}
	sw := NewSW()
	convert, _, _ := Convert1(convert1)
	swagger, err := GenSwagger(convert, "111", "111", sw)

	marshal, err := json.Marshal(swagger)

	file, err := os.OpenFile("./swagger.json", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("open file failed, err:", err)
		return
	}
	defer file.Close()

	file.Write(marshal) //写入字节切片数据

}
