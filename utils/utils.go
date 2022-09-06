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

// Merge recursively merge optional value with default one.
func Merge[T any](opt, def *T) *T {
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

		if optVal.Kind() == reflect.Pointer && optVal.Elem().Kind() == reflect.Struct {
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
