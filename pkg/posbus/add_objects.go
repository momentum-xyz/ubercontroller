package posbus

import (
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type ObjectDefinition struct {
	ID               umid.UMID             `json:"id"`
	ParentID         umid.UMID             `json:"parent_id"`
	AssetType        umid.UMID             `json:"asset_type"`
	AssetFormat      dto.Asset3dType       `json:"asset_format"` // TODO: Rename AssetType to AssetID, so Type can be used for this.
	Name             string                `json:"name"`
	Transform        cmath.ObjectTransform `json:"transform"`
	IsEditable       bool                  `json:"is_editable"`
	TetheredToParent bool                  `json:"tethered_to_parent"`
	ShowOnMiniMap    bool                  `json:"show_on_minimap"`
	//InfoUI           umid.UMID
}

type AddObjects struct {
	Objects []ObjectDefinition `json:"objects"`
}

func init() {
	registerMessage(&AddObjects{})
}

func (g *AddObjects) Type() MsgType {
	return 0x2452A9C1
}
