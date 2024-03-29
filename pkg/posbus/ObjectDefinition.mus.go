// Code generated by musgen. DO NOT EDIT.

package posbus

import (
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"github.com/ymz-ncnk/muserrs"
)

// MarshalMUS fills buf with the MUS encoding of v.
func (v ObjectDefinition) MarshalMUS(buf []byte) int {
	i := 0
	{
		si := v.ID.MarshalMUS(buf[i:])
		i += si
	}
	{
		si := v.ParentID.MarshalMUS(buf[i:])
		i += si
	}
	{
		si := v.AssetType.MarshalMUS(buf[i:])
		i += si
	}
	{
		si := v.AssetFormat.MarshalMUS(buf[i:])
		i += si
	}
	{
		length := len(v.Name)
		{
			uv := uint64(length)
			if length < 0 {
				uv = ^(uv << 1)
			} else {
				uv = uv << 1
			}
			{
				for uv >= 0x80 {
					buf[i] = byte(uv) | 0x80
					uv >>= 7
					i++
				}
				buf[i] = byte(uv)
				i++
			}
		}
		if len(buf[i:]) < length {
			panic(muserrs.ErrSmallBuf)
		}
		i += copy(buf[i:], v.Name)
	}
	{
		si := v.Transform.MarshalMUS(buf[i:])
		i += si
	}
	{
		if v.IsEditable {
			buf[i] = 0x01
		} else {
			buf[i] = 0x00
		}
		i++
	}
	{
		if v.TetheredToParent {
			buf[i] = 0x01
		} else {
			buf[i] = 0x00
		}
		i++
	}
	{
		if v.ShowOnMiniMap {
			buf[i] = 0x01
		} else {
			buf[i] = 0x00
		}
		i++
	}
	return i
}

// UnmarshalMUS parses the MUS-encoded buf, and sets the result to *v.
func (v *ObjectDefinition) UnmarshalMUS(buf []byte) (int, error) {
	i := 0
	var err error
	{
		var sv umid.UMID
		si := 0
		si, err = sv.UnmarshalMUS(buf[i:])
		if err == nil {
			v.ID = sv
			i += si
		}
	}
	if err != nil {
		return i, muserrs.NewFieldError("ID", err)
	}
	{
		var sv umid.UMID
		si := 0
		si, err = sv.UnmarshalMUS(buf[i:])
		if err == nil {
			v.ParentID = sv
			i += si
		}
	}
	if err != nil {
		return i, muserrs.NewFieldError("ParentID", err)
	}
	{
		var sv umid.UMID
		si := 0
		si, err = sv.UnmarshalMUS(buf[i:])
		if err == nil {
			v.AssetType = sv
			i += si
		}
	}
	if err != nil {
		return i, muserrs.NewFieldError("AssetType", err)
	}
	{
		var sv dto.Asset3dType
		si := 0
		si, err = sv.UnmarshalMUS(buf[i:])
		if err == nil {
			v.AssetFormat = sv
			i += si
		}
	}
	if err != nil {
		return i, muserrs.NewFieldError("AssetFormat", err)
	}
	{
		var length int
		{
			var uv uint64
			{
				if i > len(buf)-1 {
					return i, muserrs.ErrSmallBuf
				}
				shift := 0
				done := false
				for l, b := range buf[i:] {
					if l == 9 && b > 1 {
						return i, muserrs.ErrOverflow
					}
					if b < 0x80 {
						uv = uv | uint64(b)<<shift
						done = true
						i += l + 1
						break
					}
					uv = uv | uint64(b&0x7F)<<shift
					shift += 7
				}
				if !done {
					return i, muserrs.ErrSmallBuf
				}
			}
			if uv&1 == 1 {
				uv = ^(uv >> 1)
			} else {
				uv = uv >> 1
			}
			length = int(uv)
		}
		if length < 0 {
			return i, muserrs.ErrNegativeLength
		}
		if len(buf) < i+length {
			return i, muserrs.ErrSmallBuf
		}
		v.Name = string(buf[i : i+length])
		i += length
	}
	if err != nil {
		return i, muserrs.NewFieldError("Name", err)
	}
	{
		var sv cmath.Transform
		si := 0
		si, err = sv.UnmarshalMUS(buf[i:])
		if err == nil {
			v.Transform = sv
			i += si
		}
	}
	if err != nil {
		return i, muserrs.NewFieldError("Transform", err)
	}
	{
		if i > len(buf)-1 {
			return i, muserrs.ErrSmallBuf
		}
		if buf[i] == 0x01 {
			v.IsEditable = true
			i++
		} else if buf[i] == 0x00 {
			v.IsEditable = false
			i++
		} else {
			err = muserrs.ErrWrongByte
		}
	}
	if err != nil {
		return i, muserrs.NewFieldError("IsEditable", err)
	}
	{
		if i > len(buf)-1 {
			return i, muserrs.ErrSmallBuf
		}
		if buf[i] == 0x01 {
			v.TetheredToParent = true
			i++
		} else if buf[i] == 0x00 {
			v.TetheredToParent = false
			i++
		} else {
			err = muserrs.ErrWrongByte
		}
	}
	if err != nil {
		return i, muserrs.NewFieldError("TetheredToParent", err)
	}
	{
		if i > len(buf)-1 {
			return i, muserrs.ErrSmallBuf
		}
		if buf[i] == 0x01 {
			v.ShowOnMiniMap = true
			i++
		} else if buf[i] == 0x00 {
			v.ShowOnMiniMap = false
			i++
		} else {
			err = muserrs.ErrWrongByte
		}
	}
	if err != nil {
		return i, muserrs.NewFieldError("ShowOnMiniMap", err)
	}
	return i, err
}

// SizeMUS returns the size of the MUS-encoded v.
func (v ObjectDefinition) SizeMUS() int {
	size := 0
	{
		ss := v.ID.SizeMUS()
		size += ss
	}
	{
		ss := v.ParentID.SizeMUS()
		size += ss
	}
	{
		ss := v.AssetType.SizeMUS()
		size += ss
	}
	{
		ss := v.AssetFormat.SizeMUS()
		size += ss
	}
	{
		length := len(v.Name)
		{
			uv := uint64(length<<1) ^ uint64(length>>63)
			{
				for uv >= 0x80 {
					uv >>= 7
					size++
				}
				size++
			}
		}
		size += len(v.Name)
	}
	{
		ss := v.Transform.SizeMUS()
		size += ss
	}
	{
		_ = v.IsEditable
		size++
	}
	{
		_ = v.TetheredToParent
		size++
	}
	{
		_ = v.ShowOnMiniMap
		size++
	}
	return size
}
