package utils

import (
	"testing"

	"github.com/barkimedes/go-deepcopy"
	"github.com/stretchr/testify/assert"

	"github.com/momentum-xyz/ubercontroller/utils/merge"
)

func TestMergePTRs(t *testing.T) {
	t.Parallel()

	type T1 struct {
		S []int
		F float64
		N *T1
	}

	type T2 struct {
		I *int
		F float64
		M map[string]any
		T *T1
	}

	opt := &T2{
		F: 1.1,
		M: map[string]any{
			"o1": map[string]any{
				"o2": "opt2",
			},
		},
		T: &T1{
			S: []int{3, 4, 5},
			N: &T1{
				F: 3.3,
			},
		},
	}

	def := &T2{
		I: GetPTR(32),
		F: 2.2,
		M: map[string]any{
			"d1": map[string]any{
				"d2": "def2",
			},
			"o1": map[string]any{
				"d3": "def3",
			},
		},
		T: &T1{
			S: []int{1, 2},
			N: &T1{
				F: 4.4,
				N: &T1{
					F: 5.5,
				},
			},
		},
	}

	logFn := func(path string, new, current, result any) (any, bool) {
		t.Logf("Handle: path: %q, res: %+v\n", path, result)
		return nil, false
	}

	newT2 := func(t2 *T2) *T2 {
		return deepcopy.MustAnything(t2).(*T2)
	}

	tests := []struct {
		name     string
		opt      *T2
		def      *T2
		triggers []merge.Fn
		exp      *T2
	}{
		{
			name: "auto merge",
			opt:  newT2(opt),
			def:  newT2(def),
			triggers: []merge.Fn{
				logFn,
			},
			exp: &T2{
				I: def.I,
				F: opt.F,
				M: map[string]any{
					"d1": map[string]any{
						"d2": "def2",
					},
					"o1": map[string]any{
						"o2": "opt2",
						"d3": "def3",
					},
				},
				T: &T1{
					S: opt.T.S,
					N: &T1{
						S: []int{},
						F: 3.3,
						N: &T1{
							S: []int{},
							F: def.T.N.N.F,
						},
					},
				},
			},
		},
		{
			name: "mutateFn: .M.o1: new object",
			opt:  newT2(opt),
			def:  newT2(def),
			triggers: []merge.Fn{
				logFn,
				func(path string, new, current, result any) (any, bool) {
					if path == ".M.o1" {
						return map[string]any{
							"MY_MAP": map[string]any{
								"MY_VAR_1": 3.2,
								"MY_VAR_2": []int{10, 11, 12},
							},
						}, true
					}
					return nil, false
				},
			},
			exp: &T2{
				I: def.I,
				F: opt.F,
				M: map[string]any{
					"d1": map[string]any{
						"d2": "def2",
					},
					"o1": map[string]any{
						"MY_MAP": map[string]any{
							"MY_VAR_1": 3.2,
							"MY_VAR_2": []int{10, 11, 12},
						},
					},
				},
				T: &T1{
					S: opt.T.S,
					N: &T1{
						S: []int{},
						F: 3.3,
						N: &T1{
							S: []int{},
							F: def.T.N.N.F,
						},
					},
				},
			},
		},
	}

	for i := range tests {
		test := tests[i]

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			res, err := merge.Auto(test.opt, test.def, test.triggers...)
			assert.NoError(t, err)
			assert.Equal(t, test.exp, res)
		})
	}
}
