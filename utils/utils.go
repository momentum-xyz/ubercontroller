package utils

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/logger"
)

var log = logger.L()

func GetFromAny[V any](val any, defaultValue V) V {
	if val == nil {
		return defaultValue
	}

	v, ok := val.(V)
	if ok {
		return v
	}

	log.Errorf("Utils: GetFromAny: invalid value type: %+v", errors.WithStack(errors.Errorf("%T != %T", val, defaultValue)))
	return defaultValue
}

func GetFromAnyMap[K comparable, V any](amap map[K]any, key K, defaultValue V) V {
	if val, ok := amap[key]; ok {
		return GetFromAny(val, defaultValue)
	}
	return defaultValue
}

func GetPtr[T any](v T) *T {
	return &v
}

func BinID(id uuid.UUID) []byte {
	binID, err := id.MarshalBinary()
	if err != nil {
		log.Errorf("Utils: BinID: failed to marshal binary: %+v", errors.WithStack(err))
		return nil
	}
	return binID
}
