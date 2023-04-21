package posbus_test

import (
	"testing"

	"github.com/goccy/go-reflect"
	"github.com/momentum-xyz/ubercontroller/pkg/posbus"
)

func TestStringAnyMapMarshalling(t *testing.T) {
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
				"foo": "bar",
			},
		},
		{
			name: "nested",
			in: posbus.StringAnyMap{
				"foo": map[string]any{"bar": "baz"},
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
