package filters

import (
	"reflect"
)

// JSONFilter2 json field filter,inputJSON JUST match (map[string]interface,[]map[string]interface)
func JSONFilter2(inputJSON interface{}, requiredFields map[string]interface{}) {
	switch reflect.TypeOf(inputJSON).Kind() {
	case reflect.Map:

		v := reflect.ValueOf(inputJSON)
		iter := v.MapRange()
		for iter.Next() {
			if _, ok := requiredFields[iter.Key().String()]; ok {
				// TODO
				switch reflect.TypeOf(requiredFields[iter.Key().String()]).Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
				default:
					JSONFilter2(iter.Value().Interface(), requiredFields[iter.Key().String()].(map[string]interface{}))
					continue
				}
			} else {
				// TODO delete
				v.SetMapIndex(iter.Key(), reflect.Value{})
			}
		}
	case reflect.Slice, reflect.Array:
		of := reflect.ValueOf(inputJSON)
		for i := 0; i < of.Len(); i++ {
			JSONFilter2(of.Index(i).Interface(), requiredFields)
			continue
		}
	case reflect.Ptr:
		JSONFilter2(reflect.ValueOf(inputJSON).Elem().Interface(), requiredFields)
	default:

	}
}

const (
	str            = "string"
	obj            = "object"
	arr            = "array"
	num            = "number"
	dateTime       = "datetime"
	boolean        = "boolean"
	decimal        = "decimal"
	editPermission = 2
)

// Schema schema
type Schema struct {
	Title      string            `json:"title,omitempty"`
	Types      string            `json:"type,omitempty"`
	XInternal  XInternal         `json:"x-internal,omitempty"`
	Properties map[string]Schema `json:"properties,omitempty"` // type==object
	Item       *Schema           `json:"item,omitempty"`       //type==array
}

// XInternal x-internal
type XInternal struct {
	Sortable   bool    `json:"sortable"`   //排序
	Permission float64 `json:"permission"` //权限属性，第一位可不可见，第二位可不可编辑
}

// DealSchemaToFilterType 将schema处理成过滤器需要的格式
func DealSchemaToFilterType(schema Schema) map[string]interface{} {
	out := make(map[string]interface{})
	if schema.Types == obj && schema.XInternal.Permission != 0 {
		for k := range schema.Properties {
			if schema.Properties[k].XInternal.Permission != 0 {
				if schema.Properties[k].Types == obj || schema.Properties[k].Types == arr {
					out[k] = schema.Properties[k].XInternal.Permission
					res := DealSchemaToFilterType(schema.Properties[k])
					if res != nil {
						for k1, v1 := range res {
							out[k1] = v1
						}
					}
					continue
				} else {

					out[k] = schema.Properties[k].XInternal.Permission

					continue
				}
			}
		}
		return out
	}
	if schema.Types == arr && schema.XInternal.Permission != 0 {
		filterType := DealSchemaToFilterType(*schema.Item)
		if len(filterType) > 0 {
			return filterType
		}
		return nil
	}
	return nil
}

// SchemaFilterToNewSchema2 将全量schema过滤成不同权限组需要的
func SchemaFilterToNewSchema2(oldSchema interface{}, filter map[string]interface{}) {
	switch reflect.TypeOf(oldSchema).Kind() {
	case reflect.Map:
		v := reflect.ValueOf(oldSchema)
		if value := v.MapIndex(reflect.ValueOf("type")); value.IsValid() {
			if value.Elem().String() == arr {
				if itemValue := v.MapIndex(reflect.ValueOf("item")); itemValue.IsValid() {
					SchemaFilterToNewSchema2(itemValue.Interface(), filter)
				}
			}
			if value.Elem().String() == obj {
				if propertiesValue := v.MapIndex(reflect.ValueOf("properties")); propertiesValue.IsValid() {
					schemaFilter2(propertiesValue.Interface(), filter)
				}
			}
		}
	case reflect.Slice, reflect.Array:
		of := reflect.ValueOf(oldSchema)
		for i := 0; i < of.Len(); i++ {
			if value := of.Index(i); value.IsValid() {
				SchemaFilterToNewSchema2(of.Index(i).Interface(), filter)
			}
			continue
		}
	case reflect.Ptr:
		if reflect.ValueOf(oldSchema).IsValid() {
			SchemaFilterToNewSchema2(reflect.ValueOf(oldSchema).Elem().Interface(), filter)
		}
	default:

	}
}
func schemaFilter2(oldSchema interface{}, filter map[string]interface{}) {
	switch reflect.TypeOf(oldSchema).Kind() {
	case reflect.Map:
		v := reflect.ValueOf(oldSchema)
		iter := v.MapRange()
		for iter.Next() {
			if _, ok := filter[iter.Key().String()]; ok {
				switch reflect.TypeOf(filter[iter.Key().String()]).Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
					if iter.Value().IsValid() {
						schemaUpdatePermission2(iter.Value().Interface(), filter[iter.Key().String()])
					}
					continue
				default:
					if iter.Value().IsValid() {
						SchemaFilterToNewSchema2(iter.Value().Interface(), filter)
					}
					continue
				}
			} else {
				if !isLayoutComponent(iter.Value().Interface(), filter) {
					// TODO delete
					v.SetMapIndex(iter.Key(), reflect.Value{})
				}

			}
		}

	default:

	}
}
func isLayoutComponent(oldSchema interface{}, filter map[string]interface{}) bool {
	switch reflect.TypeOf(oldSchema).Kind() {
	case reflect.Map:
		v := reflect.ValueOf(oldSchema)
		if value := v.MapIndex(reflect.ValueOf("x-internal")); value.IsValid() {
			if isLayoutComponent(value.Interface(), filter) {
				if propertiesValue := v.MapIndex(reflect.ValueOf("properties")); propertiesValue.IsValid() {
					schemaFilter2(propertiesValue.Interface(), filter)
				}
				return true
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

func schemaUpdatePermission2(oldSchema interface{}, permission interface{}) {
	switch reflect.TypeOf(oldSchema).Kind() {
	case reflect.Map:
		v := reflect.ValueOf(oldSchema)

		if value := v.MapIndex(reflect.ValueOf("x-internal")); value.IsValid() {
			schemaUpdatePermission2(value.Interface(), permission)
		}
		if value := v.MapIndex(reflect.ValueOf("permission")); value.IsValid() {
			v.SetMapIndex(reflect.ValueOf("permission"), reflect.ValueOf(permission))
		}

	default:
	}
}

const (
	_id           = "_id"
	_createdAt    = "created_at"
	_creatorID    = "creator_id"
	_creatorName  = "creator_name"
	_updatedAt    = "updated_at"
	_modifierID   = "modifier_id"
	_modifierName = "modifier_name"
)

// FilterCheckData 提交数据时检查数据权限
func FilterCheckData(data interface{}, filter map[string]interface{}) bool {
	switch reflect.TypeOf(data).Kind() {
	case reflect.Map:
		v := reflect.ValueOf(data)
		iter := v.MapRange()
		for iter.Next() {
			switch iter.Key().String() {
			case _id, _createdAt, _creatorID, _creatorName, _updatedAt, _modifierID, _modifierName:
				v.SetMapIndex(iter.Key(), reflect.Value{})
				continue
			default:
				if _, ok := filter[iter.Key().String()]; !ok {
					return false
				}
				switch reflect.TypeOf(filter[iter.Key().String()]).Kind() {
				case reflect.Int8:
					if (filter[iter.Key().String()].(int8) & editPermission) == 0 {
						return false
					}
				case reflect.Int:
					if (filter[iter.Key().String()].(int) & editPermission) == 0 {
						return false
					}
				case reflect.Int16:
					if (filter[iter.Key().String()].(int16) & editPermission) == 0 {
						return false
					}
				case reflect.Int32:
					if (filter[iter.Key().String()].(int32) & editPermission) == 0 {
						return false
					}
				case reflect.Int64:
					if (filter[iter.Key().String()].(int64) & editPermission) == 0 {
						return false
					}
				case reflect.Float32:
					if (int64(filter[iter.Key().String()].(float32)) & editPermission) == 0 {
						return false
					}
				case reflect.Float64:
					if (int64(filter[iter.Key().String()].(float64)) & editPermission) == 0 {
						return false
					}
				default:
					flag := FilterCheckData(iter.Value().Interface(), filter)

					if !flag {
						return false
					}
					continue
				}
			}
		}
		return true
	case reflect.Array, reflect.Slice:
		of := reflect.ValueOf(data)
		for i := 0; i < of.Len(); i++ {
			flag := FilterCheckData(of.Index(i).Interface(), filter)

			if !flag {
				return false
			}
			continue
		}
		return true
	default:
		return false
	}
}
