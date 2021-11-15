package swagger

import (
	"errors"
	"fmt"
	"git.internal.yunify.com/qxp/misc/error2"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/code"
	"reflect"
)

const (
	// SwagVersion SwagVersion
	SwagVersion = "2.0"
	// jsons jsons
	jsons = "application/json"
	// http http
	http       = "http"
	datetime   = "datetime"
	labelValue = "label-value"
	url        = "/api/v1/structor/%s/home/form/%s/%s"
	//url        = "/api/v1/structor/noAuth/%s/form/%s/%s"
)

var (
	// ErrNoSwagger no ErrNoSwagger
	ErrNoSwagger = errors.New("no swagger")
)

// Schema  Schema
type Schema map[string]interface{}

// GenSwagger GenSwagger
func GenSwagger(schema Schema, appID, tableID string, sw *Swagger) (*SwagDoc, error) {
	var getSwagMethod = func(appID, tabID string, path map[string]SwagMethods, tags ...string) error {
		for _, tag := range tags {
			item, err := sw.GetSwagger(schema, tag, appID, tableID)
			if err != nil {
				return err
			}
			methods, err := item.GenSwagMethods()
			if err != nil {
				return err
			}
			url := fmt.Sprintf(url, appID, tableID, tag)
			path[url] = methods
		}
		return nil
	}
	path := make(map[string]SwagMethods, 0)
	tags := []string{"create", "update", "delete", "search"}
	getSwagMethod(appID, tableID, path, tags...)
	doc := &SwagDoc{
		Version:      SwagVersion,
		Paths:        path,
		Schemes:      []string{http},
		EncodingsIn:  []string{jsons},
		EncodingsOut: []string{jsons},
		BasePath:     "/",
		Info: SwagInfo{
			Title:   "structor",
			Version: "last",
			Desc:    "表单引擎",
		},
		Auth: SwagValue{
			"type":  "system",
			"extra": SwagValue{},
		},
	}
	return doc, nil
}

var swaggers = []SwagItem{
	&Create{},
	&Update{},
	&Delete{},
	&Search{},
}

// Swagger Swagger
type Swagger struct {
	swaggers map[string]SwagItem
	filter   map[string]interface{}
}

// NewSW new Ag
func NewSW() *Swagger {
	f := map[string]interface{}{
		"type":       "",
		"length":     "",
		"title":      "",
		"not_null":   "",
		"properties": "",
		"items":      "",
	}
	swagger := &Swagger{
		swaggers: make(map[string]SwagItem, len(swaggers)),
		filter:   f,
	}
	for _, sw := range swaggers {
		swagger.swaggers[sw.GetTag()] = sw
	}
	return swagger
}

// GetSwagger GetSwagger
func (s *Swagger) GetSwagger(value interface{}, method, appID, tableID string) (SwagItem, error) {
	sw, ok := s.swaggers[method]
	if !ok {
		return nil, ErrNoSwagger
	}
	sw.SetField(value, appID, tableID)
	return sw, nil
}

// Convert Convert
func (s *Swagger) Convert(schema Schema) (Schema, int64, error) {
	//total = 0
	//// value 是interface
	//for key, value := range schema {
	//
	//
	//}
	return nil, 0, nil
}

// Convert1  Convert1
func Convert1(schema Schema) (s Schema, total int64, err error) {
	s = make(Schema, 0)
	total = 0
	for key, value := range schema {
		if v, ok := value.(map[string]interface{}); ok {
			temp := make(Schema, 0)
			// 1、 判断x-component  是SubTable （子表单），AssociatedData  （关联数据）直接放行 ，
			if component, ok := v["x-component"]; ok && (component == "SubTable" || component == "AssociatedRecords") {
				continue
			}
			// 2、 判断是否是布局组件
			isLayout := isLayoutComponent(value)
			if isLayout {
				if p, ok := v["properties"]; ok {
					if p1, ok := p.(map[string]interface{}); ok {
						s2, t, err := Convert1(p1)
						if err != nil {
							return nil, 0, err
						}
						for key, value := range s2 {
							s[key] = value
						}
						total = t + total
						continue
					}
				}
			}
			total = total + 1
			for k1, v1 := range v {
				switch k1 {
				case "type":
					temp[k1] = v1
					if v1 == datetime || v1 == labelValue {
						temp[k1] = "string"
					}
					if v1 == "array" {
						if _, ok := v["items"]; !ok {
							temp["items"] = SwagValue{
								"type": "string",
							}
							continue
						}
						if _, ok := v["items"].(map[string]interface{}); !ok {
							return nil, 0, error2.NewError(code.ErrItemConvert)
						}
					}
				case "length", "title", "not_null":
					temp[k1] = v1
				case "properties":
					if p, ok := v1.(map[string]interface{}); ok {
						s2, _, _ := Convert1(p)
						temp[k1] = s2
					}
				case "items":
					if item, ok := v1.(map[string]interface{}); ok {
						temp[k1] = item
					}
				default:
					continue
				}
			}
			s[key] = temp
		} else {

			return nil, 0, error2.NewError(code.ErrValueConvert)

		}
	}
	return
}

func isLayoutComponent(value interface{}) bool {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Map:
		v := reflect.ValueOf(value)
		if value := v.MapIndex(reflect.ValueOf("x-internal")); value.IsValid() {
			if value.CanInterface() {
				return isLayoutComponent(value.Interface())
			}
		}
		if value := v.MapIndex(reflect.ValueOf("isLayoutComponent")); value.IsValid() {
			if _, ok := value.Interface().(bool); ok {
				return value.Interface().(bool)
			}
		}
	default:
		return false
	}
	return false

}
