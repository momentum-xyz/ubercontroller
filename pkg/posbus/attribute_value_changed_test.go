package posbus_test

import (
	"testing"

	"github.com/goccy/go-reflect"
	"github.com/momentum-xyz/ubercontroller/pkg/posbus"
)

func TestAttributeValueChangedMarshalling(t *testing.T) {
	subTests := []struct {
		name string
		in   posbus.AttributeValueChanged
	}{
		{
			name: "empty",
			in:   posbus.AttributeValueChanged{},
		},
		{
			name: "flat",
			in: posbus.AttributeValueChanged{
				ChangeType: "attribute_changed",
				Topic:      "foo",
				Data: posbus.AttributeValueChangedData{
					AttributeName: "bar",
					Value: &posbus.StringAnyMap{
						"baz": map[string]any{"qux": "quux"},
					},
				},
			},
		},
	}

	for _, subTest := range subTests {
		t.Run("Roundtrip marshalling "+subTest.name, func(t *testing.T) {
			in := subTest.in
			buf := make([]byte, in.SizeMUS())
			in.MarshalMUS(buf)
			out := posbus.AttributeValueChanged{}
			out.UnmarshalMUS(buf)
			if !reflect.DeepEqual(in, out) {
				t.Fatalf("%+v != %+v", in, out)
			}
		})
	}
}
