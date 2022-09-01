package generic

type Unique[V any] struct {
	val V
	cmp any
}

func NewUnique[V any](val V) Unique[V] {
	return Unique[V]{
		val: val,
		cmp: new(any),
	}
}

func (u Unique[V]) Value() V {
	return u.val
}

func (u Unique[V]) Cmp() any {
	return u.cmp
}

func (u Unique[V]) Compare(cmp any) bool {
	return u.cmp == cmp
}

func (u Unique[V]) Equals(unique Unique[V]) bool {
	return u.Compare(unique.Cmp())
}
