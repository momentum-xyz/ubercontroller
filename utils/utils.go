package utils

import (
	"reflect"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func BinID(id uuid.UUID) []byte {
	binID, err := id.MarshalBinary()
	if err != nil {
		log.Errorf("Utils: BinID: failed to marshal binary: %+v", errors.WithStack(err))
		return nil
	}
	return binID
}

// MergeStructs recursively merge optional struct with default one.
func MergeStructs[T any](opt, def *T) *T {
	if opt == nil {
		return def
	}
	if def == nil {
		return opt
	}

	var t T
	resVal := reflect.ValueOf(&t)
	optVal := reflect.ValueOf(opt)
	defVal := reflect.ValueOf(def)

	merge(resVal, optVal, defVal)

	return resVal.Interface().(*T)
}

func merge(resVal, optVal, defVal reflect.Value) {
	if optVal.Kind() == reflect.Invalid {
		resVal.Set(defVal)
		return
	}
	if defVal.Kind() == reflect.Invalid {
		resVal.Set(optVal)
		return
	}

	switch optVal.Kind() {
	case reflect.Map:
		mergeMap(resVal, optVal, defVal)
		return
	case reflect.Pointer:
		if optVal.IsNil() {
			resVal.Set(defVal)
			return
		}
		if defVal.IsNil() {
			resVal.Set(optVal)
			return
		}

		switch optVal.Elem().Kind() {
		case reflect.Struct:
			mergeStruct(resVal, optVal, defVal)
			return
		}
	}

	resVal.Set(optVal)
}

func mergeMap(resVal, optVal, defVal reflect.Value) {
	if resVal.IsNil() {
		resVal.Set(reflect.MakeMap(optVal.Type()))
	}

	var keys []reflect.Value
	for _, val := range append(optVal.MapKeys(), defVal.MapKeys()...) {
		var found bool
		for _, v := range keys {
			if val == v {
				found = true
				break
			}
		}
		if !found {
			keys = append(keys, val)
		}
	}

	for i := range keys {
		optElem := optVal.MapIndex(keys[i])
		defElem := defVal.MapIndex(keys[i])
		var resElem reflect.Value
		if optElem.Kind() != reflect.Invalid {
			resElem = reflect.New(optElem.Type()).Elem()
		} else {
			resElem = reflect.New(defElem.Type()).Elem()
		}

		merge(resElem, optElem, defElem)

		resVal.SetMapIndex(keys[i], resElem)
	}
}

func mergeStruct(resVal, optVal, defVal reflect.Value) {
	optElem := optVal.Elem()
	defElem := defVal.Elem()
	if resVal.IsNil() {
		resVal.Set(reflect.New(optElem.Type()))
	}
	resElem := resVal.Elem()

	for i := 0; i < resElem.NumField(); i++ {
		resField := resElem.Field(i)
		optField := optElem.Field(i)
		defField := defElem.Field(i)

		merge(resField, optField, defField)
	}
}
