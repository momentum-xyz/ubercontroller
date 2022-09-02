package utils

import (
	"reflect"
)

// MergeStructs recursively merge optional structure with default one
// and returns pointer to merged structure.
// If optional struct is nil, returns passed default one.
// If default struct is nil returns optional one.
func MergeStructs[T any, PtrT *T](opt, def PtrT) PtrT {
	var merge func(resVal, optVal, defVal reflect.Value)
	merge = func(resVal, optVal, defVal reflect.Value) {
		if optVal.IsNil() {
			resVal.Set(defVal)
			return
		}
		if defVal.IsNil() {
			resVal.Set(optVal)
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
