package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	t.Parallel()

	type T1 struct {
		S []int
	}

	type T2 struct {
		I *int
		M map[string]map[string]string
		T *T1
	}

	opt := T2{
		M: map[string]map[string]string{
			"o1": {
				"o2": "opt2",
			},
		},
		T: &T1{
			S: []int{3, 4, 5},
		},
	}

	def := T2{
		I: GetPtr(32),
		M: map[string]map[string]string{
			"d1": {
				"d2": "def2",
			},
			"o1": {
				"d3": "def3",
			},
		},
		T: &T1{
			S: []int{1, 2},
		},
	}

	exp := &T2{
		I: GetPtr(32),
		M: map[string]map[string]string{
			"d1": {
				"d2": "def2",
			},
			"o1": {
				"o2": "opt2",
				"d3": "def3",
			},
		},
		T: &T1{
			S: []int{3, 4, 5},
		},
	}

	res := Merge(&opt, &def)
	assert.Equal(t, exp, res)
}
