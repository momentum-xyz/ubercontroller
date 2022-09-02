package utils

type SetFn[T any, PtrT *T] func(current PtrT) PtrT

func SetNil[T any, PtrT *T]() SetFn[T, PtrT] {
	return func(current PtrT) PtrT {
		return nil
	}
}

func SetWithReplace[T any, PtrT *T](new PtrT) SetFn[T, PtrT] {
	return func(current PtrT) PtrT {
		return new
	}
}

func SetWithMerge[T any, PtrT *T](new PtrT) SetFn[T, PtrT] {
	return func(current PtrT) PtrT {
		return MergeStructs(new, current)
	}
}
