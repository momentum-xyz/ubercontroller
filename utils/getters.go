package utils

import (
	"reflect"

	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/logger"
)

var log = logger.L()

func GetPTR[T any](v T) *T {
	return &v
}

func GetFromAny[V any](val any, defaultValue V) V {
	if val == nil {
		return defaultValue
	}

	// TODO: check without reflect
	rVal := reflect.ValueOf(val)
	switch rVal.Type().Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		if rVal.IsNil() {
			return defaultValue
		}
	}

	v, ok := val.(V)
	if ok {
		return v
	}

	log.Errorf(
		"Utils: GetFromAny: invalid value type: %+v", errors.WithStack(errors.Errorf("%T != %T", val, defaultValue)),
	)

	return defaultValue
}

func GetFromAnyMap[K comparable, V any](amap map[K]any, key K, defaultValue V) V {
	if val, ok := amap[key]; ok {
		return GetFromAny(val, defaultValue)
	}
	return defaultValue
}

func GetKeyByValueFromMap[K comparable, V comparable](m map[K]V, val V) (K, bool) {
	for k, v := range m {
		if v == val {
			return k, true
		}
	}

	var k K
	return k, false
}
