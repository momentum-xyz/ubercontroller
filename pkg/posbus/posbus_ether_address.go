package posbus

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/niubaoshu/gotiny"
)

type PBEthAddress common.Address

func (v PBEthAddress) MarshalMUS(buf []byte) int {
	b := gotiny.Marshal(&v)
	copy(buf, b)
	return len(b)
}

func (v *PBEthAddress) UnmarshalMUS(buf []byte) (int, error) {
	l := gotiny.Unmarshal(buf, v)
	return l, nil
}

func (v PBEthAddress) SizeMUS() int {
	return len(gotiny.Marshal(&v))
}
