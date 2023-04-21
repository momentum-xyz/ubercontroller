package posbus

import (
	"github.com/niubaoshu/gotiny"
)

func init() {
	// workaround, sometimes when receiving StringAnyMap, we end up in the 'interface' branch of gotiny unmarshalling and not the map handling :/
	// Needs some more debugging. But for now avoid the panic when handling these.
	gotiny.Register("")
	gotiny.Register(map[string]any{})
}

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
