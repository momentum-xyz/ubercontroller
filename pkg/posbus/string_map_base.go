package posbus

import (
	"github.com/niubaoshu/gotiny"
)

type StringAnyMap map[string]any

func (v StringAnyMap) MarshalMUS(buf []byte) int {
	b := gotiny.Marshal(&v)
	copy(buf, b)
	return len(b)
}

func (v *StringAnyMap) UnmarshalMUS(buf []byte) (int, error) {
	l := gotiny.Unmarshal(buf, v)
	return l, nil
}

func (v StringAnyMap) SizeMUS() int {
	return len(gotiny.Marshal(&v))
}
