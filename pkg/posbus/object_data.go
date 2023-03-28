package posbus

import (
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type ObjectData struct {
	ID umid.UMID
	//Entries map[ObjectDataIndex]interface{}
	Entries map[string]string
}

func init() {
	registerMessage(ObjectData{})
}

func (l *ObjectData) GetType() MsgType {
	return 0xCACE197C
}

// // addToMaps(TypeSetObjectData, "set_object_data", ObjectData{})
//
//	func (o *ObjectData) MarshalJSON() ([]byte, error) {
//		q := make(map[string]map[string]interface{})
//		for k, v := range o.Entries {
//			t, ok := q[string(k.Kind)]
//			if !ok {
//				t = make(map[string]interface{})
//			}
//			t[k.SlotName] = v
//			q[string(k.Kind)] = t
//		}
//
//		return json.Marshal(
//			&struct {
//				ID      umid.UMID                         `json:"id"`
//				Entries map[string]map[string]interface{} `json:"entries"`
//			}{
//				ID:      o.ID,
//				Entries: q,
//			},
//		)
//	}
type ObjectDataIndex struct {
	Kind     entry.UnitySlotType
	SlotName string
}