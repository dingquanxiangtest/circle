package comet

import (
	"git.internal.yunify.com/qxp/molecule/internal/dorm/dao"
	"reflect"

	"git.internal.yunify.com/qxp/misc/header2"
	"git.internal.yunify.com/qxp/misc/id2"
	"git.internal.yunify.com/qxp/misc/time2"
)

// Entity interface of Entity
type Entity interface{}

// EntityOpt entity options
type EntityOpt func(e defaultFieldMap)

type defaultFieldMap map[string]interface{}

const (
	_id           = "_id"
	_createdAt    = "created_at"
	_creatorID    = "creator_id"
	_creatorName  = "creator_name"
	_updatedAt    = "updated_at"
	_modifierID   = "modifier_id"
	_modifierName = "modifier_name"
)

// WithID default field with id
func WithID() EntityOpt {
	return func(d defaultFieldMap) {
		d[_id] = id2.GenID()
	}
}

// WithCreated default field with created_at、creator_id and creator_name
func WithCreated(profile *header2.Profile) EntityOpt {
	return func(d defaultFieldMap) {
		d[_createdAt] = time2.Now()
		d[_creatorID] = profile.UserID
		d[_creatorName] = profile.UserName
	}
}

// WithUpdated default field with updated_at、modifier_id and modifier_name
func WithUpdated(profile *header2.Profile) EntityOpt {
	return func(d defaultFieldMap) {
		d[_updatedAt] = time2.Now()
		d[_modifierID] = profile.UserID
		d[_modifierName] = profile.UserName
	}
}

func defaultFieldWithDep(e Entity, dep int, opts ...EntityOpt) Entity {
	if e == nil {
		return e
	}
	value := reflect.ValueOf(e)
	switch _t := reflect.TypeOf(e); _t.Kind() {
	case reflect.Ptr:
		return defaultFieldWithDep(value.Elem(), dep, opts...)
	case reflect.Array, reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			if !value.Index(i).CanInterface() {
				continue
			}
			val := defaultFieldWithDep(value.Index(i).Interface(), dep, opts...)
			value.Index(i).Set(reflect.ValueOf(val))
		}
	case reflect.Map:
		if dep > 0 {
			dep--
			iter := value.MapRange()
			for iter.Next() {
				if !iter.Value().CanInterface() {
					continue
				}
				val := defaultFieldWithDep(iter.Value().Interface(), dep, opts...)
				value.SetMapIndex(reflect.ValueOf(iter.Key().String()), reflect.ValueOf(val))
			}
			return e
		}
		defaultFieldMap := make(map[string]interface{})
		for _, opt := range opts {
			opt(defaultFieldMap)
		}
		for key, val := range defaultFieldMap {
			value.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
		}
	default:
		return e
	}

	return e
}

func defaultField(e Entity, opts ...EntityOpt) Entity {
	return defaultFieldWithDep(e, 0, opts...)
}

func dataTrans(data interface{})[]interface{}  {
	entity := make([]interface{}, 0)
	switch reflect.TypeOf(data).Kind() {
	case reflect.Map:
		v := reflect.ValueOf(data)
		entity = append(entity,v.Interface())
	case reflect.Slice, reflect.Array:
		v := reflect.ValueOf(data)
		for i := 0; i < v.Len(); i++ {
			entity = append(entity,v.Index(i).Interface())
		}
	}
	return entity
}

func filterUpdateData(bus *bus,data []dao.Data) []map[string]interface{} {
	res :=  make([]map[string]interface{},0)
	if en, ok := bus.input.Entity.(map[string]interface{}); ok {
		for _,d := range data{
			dm := map[string]interface{}(d)
			if v,ok := dm[_id];ok{
				en[_id] = v
			}
			res = append(res,en)
		}
	}
	return res
}