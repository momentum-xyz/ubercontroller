package utils

import (
	"reflect"
)

// MergeStructs recursively merge optional structure with default one.
// If optional struct is nil, return passed default one,
// otherwise return pointer to new struct with merged fields.
func MergeStructs[T any](opt, def *T) *T {
	var merge func(resVal, optVal, defVal reflect.Value)
	merge = func(resVal, optVal, defVal reflect.Value) {
		if optVal.IsNil() {
			resVal.Set(defVal)
			return
		}

		if resVal.Kind() == reflect.Pointer && resVal.Elem().Kind() == reflect.Struct {
			resElem := resVal.Elem()
			optElem := optVal.Elem()
			defElem := defVal.Elem()
			for i := 0; i < resElem.NumField(); i++ {
				resField := resElem.Field(i)
				optField := optElem.Field(i)
				defField := defElem.Field(i)
				merge(resField, optField, defField)
			}
			return
		}

		resVal.Set(optVal)
	}

	var t T
	resVal := reflect.ValueOf(&t)
	optVal := reflect.ValueOf(opt)
	defVal := reflect.ValueOf(def)

	merge(resVal, optVal, defVal)

	return resVal.Interface().(*T)
}
