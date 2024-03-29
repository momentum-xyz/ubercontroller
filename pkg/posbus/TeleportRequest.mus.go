// Code generated by musgen. DO NOT EDIT.

package posbus

import (
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"github.com/ymz-ncnk/muserrs"
)

// MarshalMUS fills buf with the MUS encoding of v.
func (v TeleportRequest) MarshalMUS(buf []byte) int {
	i := 0
	{
		si := v.Target.MarshalMUS(buf[i:])
		i += si
	}
	return i
}

// UnmarshalMUS parses the MUS-encoded buf, and sets the result to *v.
func (v *TeleportRequest) UnmarshalMUS(buf []byte) (int, error) {
	i := 0
	var err error
	{
		var sv umid.UMID
		si := 0
		si, err = sv.UnmarshalMUS(buf[i:])
		if err == nil {
			v.Target = sv
			i += si
		}
	}
	if err != nil {
		return i, muserrs.NewFieldError("Target", err)
	}
	return i, err
}

// SizeMUS returns the size of the MUS-encoded v.
func (v TeleportRequest) SizeMUS() int {
	size := 0
	{
		ss := v.Target.SizeMUS()
		size += ss
	}
	return size
}
