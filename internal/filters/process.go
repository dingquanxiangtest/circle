package filters

import (
	"reflect"
	"sort"
)

// SchemaLoseWeight 简化返回值
func SchemaLoseWeight(oldSchema interface{}, rest map[string]interface{}) {
	switch t := reflect.TypeOf(oldSchema);t.Kind() {
	case reflect.Map:
		v := reflect.ValueOf(oldSchema)
		iter := v.MapRange()
		for iter.Next() {
			if iter.Key().String() == "type" {
				if iter.Value().Elem().String() == arr {
					SchemaLoseWeight(v.MapIndex(reflect.ValueOf("items")).Interface(), rest)
				}
				if iter.Value().Elem().String() == obj {
					fieldFill(v.MapIndex(reflect.ValueOf("properties")).Interface(), rest)
				}
			}
		}
	case reflect.Slice, reflect.Array:
		of := reflect.ValueOf(oldSchema)
		for i := 0; i < of.Len(); i++ {
			SchemaLoseWeight(of.Index(i).Interface(), rest)
		}
	case reflect.Ptr:
		SchemaLoseWeight(reflect.ValueOf(oldSchema).Elem(), rest)
	default:
	}
}

func fieldFill(oldSchema interface{}, rest map[string]interface{})  {
	switch reflect.TypeOf(oldSchema).Kind() {
	case reflect.Map:
		v := reflect.ValueOf(oldSchema)
		iter := v.MapRange()
		for iter.Next() {
			val := iter.Value().Elem()
			switch val.Kind() {
			case reflect.Map:
				vm := val.Interface().(map[string]interface{})

				if tp,ok := vm["type"];ok{
					if tp == arr {
						rc := make(map[string]interface{})
						if title,ok := vm["title"];ok{
							rc["title"] = title
							rc["type"] = tp
						}
						rest[iter.Key().String()] = rc
						if item,ok := vm["items"];ok{
							SchemaLoseWeight(item,rc)
						}
					}else {
						if title,ok := vm["title"];ok{
							rc := make(map[string]interface{})
							rc["title"] = title
							rc["type"] = tp
							rest[iter.Key().String()] = rc
						}
					}
				}
			}
		}
	}
}
// SchemaFilter SchemaFilter
func SchemaFilter(oldSchema interface{},filter *[]string) {
	if len(*filter) == 0 {
		return
	}
	switch t := reflect.TypeOf(oldSchema);t.Kind() {
	case reflect.Map:
		v := reflect.ValueOf(oldSchema)
		iter := v.MapRange()
		for iter.Next() {
			if iter.Key().String() == "type" {
				if iter.Value().Elem().String() == arr {
					SchemaFilter(v.MapIndex(reflect.ValueOf("item")).Interface(),filter)
				}
				if iter.Value().Elem().String() == obj {
					DoFilter(v.MapIndex(reflect.ValueOf("properties")).Interface().(map[string]interface{}),filter)
				}
			}
		}
	case reflect.Slice, reflect.Array:
		of := reflect.ValueOf(oldSchema)
		for i := 0; i < of.Len(); i++ {
			SchemaFilter(of.Index(i).Interface(),filter)
		}
	case reflect.Ptr:
		SchemaFilter(reflect.ValueOf(oldSchema).Elem(),filter)
	default:
	}
}

// DoFilter filter map
func DoFilter(schema map[string]interface{},filter *[]string)  {
	if len(*filter) == 0 {
		return
	}
	sort.Strings(*filter)
	v := reflect.ValueOf(schema)
	iter := v.MapRange()
	for iter.Next() {
		if !in(iter.Key().String(),filter){
			delete(schema,iter.Key().String())
		}
	}
}

func in(target string, filter *[]string) bool {
	index := sort.SearchStrings(*filter, target)
	if index < len(*filter) && (*filter)[index] == target {
		return true
	}
	return false
}

