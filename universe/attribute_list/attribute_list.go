package attribute_list

import (
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
	"sync/atomic"
)

type AttributeData struct {
	id               entry.AttributeID
	options          atomic.Pointer[entry.AttributeOptions]
	effectiveOptions *entry.AttributeOptions
	value            atomic.Pointer[entry.AttributeValue]
	attribute        universe.Attribute
}

type AttributeList[indexType comparable] struct {
	id   indexType
	data map[indexType]AttributeData
}

func (a *AttributeData) GetOptions() *entry.AttributeOptions {

}

func (a *AttributeData) SetOptions(modifyFn modify.Fn[entry.AttributeOptions], updateDB bool) error {

}

func (a *AttributeData) GetValue() *string {

}

func (a *AttributeData) SetValue(modifyFn modify.Fn[string], updateDB bool) error {

}

func (a *AttributeData) GetEntry() *entry.Attribute {

}

func (a *AttributeData) LoadFromEntry(entry *entry.Attribute) error {

}

func (a *AttributeData) GetEffectiveOptions() *entry.AttributeOptions {
	if a.effectiveOptions == nil {
		a.effectiveOptions = utils.MergePTRs(a.options.Load(), a.attribute.GetOptions())
	}
	return a.effectiveOptions
}
