package posbus

import "github.com/momentum-xyz/ubercontroller/utils/umid"

type SetWorld struct {
	ID              umid.UMID `json:"id"`
	Name            string    `json:"name"`
	Avatar          umid.UMID `json:"avatar"`
	Owner           umid.UMID `json:"owner"`
	Avatar3DAssetID umid.UMID `json:"avatar_3d_asset_id"`
}

func init() {
	registerMessage(&SetWorld{})
}

func (g *SetWorld) Type() MsgType {
	return 0xCCDF2E49
}
