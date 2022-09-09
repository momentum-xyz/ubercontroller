package modify

import "github.com/momentum-xyz/ubercontroller/utils"

type Fn[T any] func(current *T) *T

func SetNil[T any]() Fn[T] {
	return func(current *T) *T {
		return nil
	}
}

func ReplaceWith[T any](new *T) Fn[T] {
	return func(current *T) *T {
		return new
	}
}

func MergeWith[T any](new *T) Fn[T] {
	return func(current *T) *T {
		return utils.MergeStructs(new, current)
	}
}
