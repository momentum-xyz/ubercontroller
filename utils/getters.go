package utils

import (
	"log"

	"github.com/goccy/go-reflect"

	"github.com/pkg/errors"
)

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
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		if rVal.IsNil() {
			return defaultValue
		}
	}

	v, ok := val.(V)
	if ok {
		return v
	}

	// TODO: return the error!
	log.Printf(
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
