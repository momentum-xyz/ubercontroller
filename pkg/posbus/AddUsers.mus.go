// Code generated by musgen. DO NOT EDIT.

package posbus

import "github.com/ymz-ncnk/muserrs"

// MarshalMUS fills buf with the MUS encoding of v.
func (v AddUsers) MarshalMUS(buf []byte) int {
	i := 0
	{
		length := len(v.Users)
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
		for _, el := range v.Users {
			{
				si := el.MarshalMUS(buf[i:])
				i += si
			}
		}
	}
	return i
}

// UnmarshalMUS parses the MUS-encoded buf, and sets the result to *v.
func (v *AddUsers) UnmarshalMUS(buf []byte) (int, error) {
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
		v.Users = make([]UserData, length)
		for j := 0; j < length; j++ {
			{
				var sv UserData
				si := 0
				si, err = sv.UnmarshalMUS(buf[i:])
				if err == nil {
					v.Users[j] = sv
					i += si
				}
			}
			if err != nil {
				err = muserrs.NewSliceError(j, err)
				break
			}
		}
	}
	if err != nil {
		return i, muserrs.NewFieldError("Users", err)
	}
	return i, err
}

// SizeMUS returns the size of the MUS-encoded v.
func (v AddUsers) SizeMUS() int {
	size := 0
	{
		length := len(v.Users)
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
		for _, el := range v.Users {
			{
				ss := el.SizeMUS()
				size += ss
			}
		}
	}
	return size
}
