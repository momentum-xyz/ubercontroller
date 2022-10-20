package merge

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/logger"
)

type Fn func(path string, old, new, result any) (any, bool)

var log = logger.L()

// Auto recursively merge optional pointer with default one.
func Auto[T any](opt, def *T, triggers ...Fn) (*T, error) {
	if opt == nil {
		return def, nil
	}
	if def == nil {
		return opt, nil
	}

	optVal := reflect.ValueOf(opt)
	defVal := reflect.ValueOf(def)

	resVal, err := merge(optVal, defVal, ".", triggers...)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to merge")
	}

	res, ok := resVal.Interface().(*T)
	if !ok {
		return nil, errors.Errorf("invalid result type: %T != %T", resVal.Interface(), def)
	}

	return res, nil
}

func merge(optVal, defVal reflect.Value, path string, triggers ...Fn) (reflect.Value, error) {
	if !optVal.IsValid() {
		resVal, err := mergeHandle(path, optVal, defVal, defVal, triggers...)
		if err != nil {
			return reflect.Value{}, errors.WithMessagef(err, "failed to handle empty opt val: %q:", path)
		}
		return resVal, nil
	}
	if !defVal.IsValid() {
		resVal, err := mergeHandle(path, optVal, defVal, optVal, triggers...)
		if err != nil {
			return reflect.Value{}, errors.WithMessagef(err, "failed to handle empty def val: %q", path)
		}
		return resVal, nil
	}

	optValKind := optVal.Kind()
	if optValKind == reflect.Map ||
		optValKind == reflect.Slice ||
		optValKind == reflect.Pointer ||
		optValKind == reflect.Interface {
		if optVal.IsZero() {
			resVal, err := mergeHandle(path, optVal, defVal, defVal, triggers...)
			if err != nil {
				return reflect.Value{}, errors.WithMessagef(err, "failed to handle zero opt val: %q", path)
			}
			return resVal, nil
		}
		if defVal.IsZero() {
			resVal, err := mergeHandle(path, optVal, defVal, optVal)
			if err != nil {
				return reflect.Value{}, errors.WithMessagef(err, "failed to handle zero def val: %q", path)
			}
			return resVal, nil
		}
	}

	switch optValKind {
	case reflect.Struct:
		resVal, err := mergeStruct(optVal, defVal, path, triggers...)
		if err != nil {
			return reflect.Value{}, errors.WithMessagef(err, "failed to merge struct: %q", path)
		}
		return resVal, nil
	case reflect.Map:
		resVal, err := mergeMap(optVal, defVal, path, triggers...)
		if err != nil {
			return reflect.Value{}, errors.WithMessagef(err, "failed to merge map: %q", path)
		}
		return resVal, nil
	case reflect.Pointer:
		res, err := merge(optVal.Elem(), defVal.Elem(), path, triggers...)
		if err != nil {
			return reflect.Value{}, errors.WithMessagef(err, "failed to merge pointer: %q", path)
		}

		resVal := reflect.New(res.Type())
		resVal.Elem().Set(res)

		return resVal, nil
	case reflect.Interface:
		res, err := merge(optVal.Elem(), defVal.Elem(), path, triggers...)
		if err != nil {
			return reflect.Value{}, errors.WithMessagef(err, "failed to merge interface: %q", path)
		}

		return res, nil
	}

	resVal, err := mergeHandle(path, optVal, defVal, optVal, triggers...)
	if err != nil {
		return reflect.Value{}, errors.WithMessagef(err, "failed to handle final opt val: %q", path)
	}

	return resVal, nil
}

func mergeMap(optVal, defVal reflect.Value, path string, triggers ...Fn) (reflect.Value, error) {
	var keys []reflect.Value
	for _, val := range append(optVal.MapKeys(), defVal.MapKeys()...) {
		var found bool
		for _, v := range keys {
			if val.Interface() == v.Interface() {
				found = true
				break
			}
		}
		if !found {
			keys = append(keys, val)
		}
	}

	resVal := reflect.MakeMap(optVal.Type())
	for i := range keys {
		optElem := optVal.MapIndex(keys[i])
		defElem := defVal.MapIndex(keys[i])

		keyPath := mergeAddPathKey(path, keys[i].Interface())

		resElem, err := merge(optElem, defElem, keyPath, triggers...)
		if err != nil {
			return reflect.Value{}, errors.WithMessagef(
				err, "failed to merge map: %q: key %+v", path, keys[i].Interface(),
			)
		}
		if !resElem.IsValid() {
			continue
		}

		if !resElem.Type().AssignableTo(resVal.Type().Elem()) {
			return reflect.Value{}, errors.Errorf(
				"failed to set map: %q: key: %+v: %+v != %+v",
				path, keys[i].Interface(), resElem.Type(), resVal.Type().Elem())
		}
		resVal.SetMapIndex(keys[i], resElem)
	}

	resVal, err := mergeHandle(path, optVal, defVal, resVal, triggers...)
	if err != nil {
		return reflect.Value{}, errors.WithMessagef(err, "failed to handle map res val: %q", path)
	}
	return resVal, nil
}

func mergeStruct(optVal, defVal reflect.Value, path string, triggers ...Fn) (reflect.Value, error) {
	resVal := reflect.New(optVal.Type()).Elem()
	for i := 0; i < resVal.NumField(); i++ {
		optField := optVal.Field(i)
		defField := defVal.Field(i)

		fieldName := resVal.Type().Field(i).Name
		fieldPath := mergeAddPathKey(path, fieldName)

		newField, err := merge(optField, defField, fieldPath, triggers...)
		if err != nil {
			return reflect.Value{}, errors.WithMessagef(err, "failed to merge struct field: %q", fieldName)
		}
		if !newField.IsValid() {
			continue
		}

		if !newField.Type().AssignableTo(resVal.Type().Field(i).Type) {
			return reflect.Value{}, errors.Errorf(
				"failed to set struct field: %q: %s: %+v != %+v",
				path, fieldName, newField.Type(), resVal.Type().Field(i).Type,
			)
		}

		resVal.Field(i).Set(newField)
	}

	resVal, err := mergeHandle(path, optVal, defVal, resVal, triggers...)
	if err != nil {
		return reflect.Value{}, errors.WithMessagef(err, "failed to handle struct res val: %q", path)
	}
	return resVal, nil
}

func mergeHandle(path string, optVal, defVal, resVal reflect.Value, triggers ...Fn) (reflect.Value, error) {
	var res any
	if resVal.IsValid() {
		res = resVal.Interface()
	}

	if len(triggers) == 0 {
		return reflect.ValueOf(res), nil
	}

	var opt any
	var def any
	if optVal.IsValid() {
		opt = optVal.Interface()
	}
	if defVal.IsValid() {
		def = defVal.Interface()
	}

	for i := range triggers {
		if val, ok := triggers[i](path, def, opt, res); ok {
			return reflect.ValueOf(val), nil
		}
	}

	return reflect.ValueOf(res), nil
}

func mergeAddPathKey(path string, key any) string {
	res := bytes.NewBufferString(path)
	if path != "." {
		res.WriteByte('.')
	}

	if _, err := fmt.Fprintf(res, "%+v", key); err != nil {
		log.Errorf(
			"Utils: mergeAddPathKey: failed to fprintf: %+v", errors.WithStack(err),
		)
		return "invalid-path-key"
	}

	return res.String()
}
