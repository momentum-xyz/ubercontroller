package assets_2d

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"github.com/pkg/errors"
)

// @Summary Get 2d asset
// @Description Returns a 2d asset
// @Tags assets2d
// @Security Bearer
// @Param asset2dID path string true "Asset2d UMID"
// @Success 200 {array} dto.Asset2d
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/assets-2d/{asset2dID} [get]
func (a *Assets2d) apiGetAsset2d(c *gin.Context) {
	asset2dID, err := umid.Parse(c.Param("asset2dID"))
	if err != nil {
		err := errors.WithMessage(err, "Assets2d: apiGetAsset2d: failed to parse asset 2d umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_asset_2d_id", err, a.log)
		return
	}

	asset2d, ok := a.GetAsset2d(asset2dID)
	if !ok {
		err := errors.Errorf("Assets2d: apiGetAsset2d: asset 2d not found: %s", asset2dID)
		api.AbortRequest(c, http.StatusNotFound, "asset_2d_not_found", err, a.log)
		return
	}

	out := dto.Asset2d{
		Meta:    dto.Asset2dMeta(asset2d.GetMeta()),
		Options: asset2d.GetOptions(),
	}

	c.JSON(http.StatusOK, out)
}
