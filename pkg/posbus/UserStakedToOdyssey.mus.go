// Code generated by musgen. DO NOT EDIT.

package posbus

import (
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"github.com/ymz-ncnk/muserrs"
)

// MarshalMUS fills buf with the MUS encoding of v.
func (v UserStakedToOdyssey) MarshalMUS(buf []byte) int {
	i := 0
	{
		length := len(v.TransactionHash)
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
		i += copy(buf[i:], v.TransactionHash)
	}
	{
		si := v.ObjectID.MarshalMUS(buf[i:])
		i += si
	}
	{
		length := len(v.Amount)
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
		i += copy(buf[i:], v.Amount)
	}
	{
		length := len(v.Comment)
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
		i += copy(buf[i:], v.Comment)
	}
	return i
}

// UnmarshalMUS parses the MUS-encoded buf, and sets the result to *v.
func (v *UserStakedToOdyssey) UnmarshalMUS(buf []byte) (int, error) {
	i := 0
	var err error
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
		v.TransactionHash = string(buf[i : i+length])
		i += length
	}
	if err != nil {
		return i, muserrs.NewFieldError("TransactionHash", err)
	}
	{
		var sv umid.UMID
		si := 0
		si, err = sv.UnmarshalMUS(buf[i:])
		if err == nil {
			v.ObjectID = sv
			i += si
		}
	}
	if err != nil {
		return i, muserrs.NewFieldError("ObjectID", err)
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
		v.Amount = string(buf[i : i+length])
		i += length
	}
	if err != nil {
		return i, muserrs.NewFieldError("Amount", err)
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
		v.Comment = string(buf[i : i+length])
		i += length
	}
	if err != nil {
		return i, muserrs.NewFieldError("Comment", err)
	}
	return i, err
}

// SizeMUS returns the size of the MUS-encoded v.
func (v UserStakedToOdyssey) SizeMUS() int {
	size := 0
	{
		length := len(v.TransactionHash)
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
		size += len(v.TransactionHash)
	}
	{
		ss := v.ObjectID.SizeMUS()
		size += ss
	}
	{
		length := len(v.Amount)
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
		size += len(v.Amount)
	}
	{
		length := len(v.Comment)
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
		size += len(v.Comment)
	}
	return size
}