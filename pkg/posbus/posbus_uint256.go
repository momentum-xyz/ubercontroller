package posbus

import (
	"github.com/holiman/uint256"
	"github.com/niubaoshu/gotiny"
)

type PBUint256 uint256.Int

func (v PBUint256) MarshalMUS(buf []byte) int {
	b := gotiny.Marshal(&v)
	copy(buf, b)
	return len(b)
}

func (v *PBUint256) UnmarshalMUS(buf []byte) (int, error) {
	l := gotiny.Unmarshal(buf, v)
	return l, nil
}

func (v PBUint256) SizeMUS() int {
	return len(gotiny.Marshal(&v))
}
