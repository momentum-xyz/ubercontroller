package utils

type SetFn[T any] func(current *T) *T

func SetNil[T any]() SetFn[T] {
	return func(current *T) *T {
		return nil
	}
}

func SetWithReplace[T any](new *T) SetFn[T] {
	return func(current *T) *T {
		return new
	}
}

func SetWithMerge[T any](new *T) SetFn[T] {
	return func(current *T) *T {
		return MergeStructs(new, current)
	}
}
