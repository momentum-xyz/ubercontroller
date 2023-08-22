package attributes

import (
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

// ModifyFnFromAttributeValue creates a modify function for (api) attribute input.
func ModifyFn(value map[string]any) modify.Fn[entry.AttributePayload] {
	return func(current *entry.AttributePayload) (*entry.AttributePayload, error) {
		// TODO: pointer type alias for map does not make any sense, remove it.
		var v entry.AttributeValue = value
		if current == nil {
			return entry.NewAttributePayload(&v, nil), nil
		}
		current.Value = &v
		return current, nil
	}
}
