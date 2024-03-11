package db

import (
	"reflect"
	"strings"
)

type Condition struct {
	Eq    map[string]interface{}
	NotEq map[string]interface{}
}
type FieldsAndPointers struct {
	Fields      []string
	Pointers    map[string]interface{}
	FieldsTypes map[string][]string
}

// Хочет указатель на структру
func GetFieldsAndPointers(s interface{}) FieldsAndPointers {
	val := reflect.ValueOf(s).Elem()
	res := FieldsAndPointers{
		Fields:      make([]string, 0),
		Pointers:    make(map[string]interface{}),
		FieldsTypes: make(map[string][]string),
	}
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i).Tag.Get("db")
		res.Fields = append(res.Fields, field)
		res.Pointers[field] = val.Field(i).Addr().Interface()
		tp := val.Type().Field(i).Tag.Get("db_type")
		tps := strings.Split(tp, ",")
		if len(tps) > 1 {
			for _, t := range tps {
				res.FieldsTypes[field] = append(res.FieldsTypes[field], t)
			}
		}
		if len(tps) == 1 {
			res.FieldsTypes[field] = append(res.FieldsTypes[field], tps[0])
		}
	}

	return res
}
