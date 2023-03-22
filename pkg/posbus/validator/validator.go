package validator

import "errors"

var ErrNegative error = errors.New("negative")

func Positive(n int) error {
	if n < 0 {
		return ErrNegative
	}
	return nil
}
