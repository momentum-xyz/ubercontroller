package modify

import (
	"github.com/momentum-xyz/ubercontroller/utils/merge"
)

type Fn[T any] func(current *T) (*T, error)

func Nop[T any]() Fn[T] {
	return func(current *T) (*T, error) {
		return current, nil
	}
}

func SetNil[T any]() Fn[T] {
	return func(current *T) (*T, error) {
		return nil, nil
	}
}

func ReplaceWith[T any](new *T) Fn[T] {
	return func(current *T) (*T, error) {
		return new, nil
	}
}

func MergeWith[T any](new *T, triggers ...merge.Fn) Fn[T] {
	return func(current *T) (*T, error) {
		return merge.Auto(new, current, triggers...)
	}
}
