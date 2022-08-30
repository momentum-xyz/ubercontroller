package utils

import (
	"reflect"

	"github.com/pkg/errors"
)

// MergeStructs merge optional structure with default one.
// If optional struct is nil, return passed default one,
// otherwise return pointer to new struct with merged fields.
func MergeStructs[T any](opt, def *T) *T {
	if opt == nil {
		return def
	}

	var t T
	resVal := reflect.ValueOf(&t)
	resElem := resVal.Elem()

	optElem := reflect.ValueOf(opt).Elem()
	defElem := reflect.ValueOf(def).Elem()
	if optElem.Kind() != reflect.Struct {
		log.Errorf("Utils: MergeStructs: invalid type: %+v", errors.WithStack(errors.Errorf("%T is not a struture", opt)))
		return def
	}

	for i := 0; i < resElem.NumField(); i++ {
		resField := resElem.Field(i)

		optField := optElem.Field(i)
		if !optField.IsNil() {
			resField.Set(optField)
			continue
		}

		resField.Set(defElem.Field(i))
	}

	return resVal.Interface().(*T)
}
