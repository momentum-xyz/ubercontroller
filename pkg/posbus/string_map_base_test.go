package posbus_test

import (
	"testing"

	"github.com/goccy/go-reflect"
	"github.com/momentum-xyz/ubercontroller/pkg/posbus"
	"github.com/momentum-xyz/ubercontroller/utils"
)

func TestStringAnyMapMarshalling(t *testing.T) {
	str := "bar"
	// num := 42
	// flo := 42.42

	var strPtr *string = nil
	var numPtr *int = nil
	var floPtr *float64 = nil

	var nestedMap = posbus.StringAnyMap{
		"bar": "baz",
		"num": 42,
		"flo": 42.42,
	}
	var mp map[string]any
	utils.MapEncode(nestedMap, &mp)

	subTests := []struct {
		name string
		in   posbus.StringAnyMap
	}{
		{
			name: "empty",
			in:   posbus.StringAnyMap{},
		},
		{
			name: "flat",
			in: posbus.StringAnyMap{
				"foo": str,
			},
		},
		{
			name: "nested",
			in: posbus.StringAnyMap{
				"foo": map[string]any{
					"bar": "baz",
					"num": 42,
					"flo": 42.42,
					"nested": map[string]any{
						"bar": "baz",
						"num": 42,
						"flo": 42.42,
					},
					"mp": mp,
				},
			},
		},
		{
			name: "bool",
			in: posbus.StringAnyMap{
				"foo": true,
			},
		},
		{
			name: "int",
			in: posbus.StringAnyMap{
				"foo": 42,
			},
		},
		{
			name: "float",
			in: posbus.StringAnyMap{
				"foo": 42.42,
			},
		},
		{
			name: "string ptr",
			in: posbus.StringAnyMap{
				"foo": &str,
			},
		},
		{
			name: "string ptr",
			in: posbus.StringAnyMap{
				"foo": strPtr,
			},
		},
		{
			name: "int ptr",
			in: posbus.StringAnyMap{
				"foo": numPtr,
			},
		},
		{
			name: "float ptr",
			in: posbus.StringAnyMap{
				"foo": floPtr,
			},
		},
		// fails, TODO fix
		// {
		// 	name: "array",
		// 	in: posbus.StringAnyMap{
		// 		"foo": []any{"bar", "baz"},
		// 	},
		// },
		{
			name: "null",
			in: posbus.StringAnyMap{
				"foo": nil,
			},
		},
	}

	for _, subTest := range subTests {
		t.Run("Roundtrip marshalling "+subTest.name, func(t *testing.T) {
			in := subTest.in
			buf := make([]byte, in.SizeMUS())
			in.MarshalMUS(buf)
			out := posbus.StringAnyMap{}
			out.UnmarshalMUS(buf)
			if !reflect.DeepEqual(in, out) {
				t.Fatalf("%+v != %+v", in, out)
			}
		})
	}
}
