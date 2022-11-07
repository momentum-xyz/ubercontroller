package assets_3d

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/api"
	"github.com/momentum-xyz/ubercontroller/universe/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils"
)

// @Summary Get list of assets 3d; can be filtered by query param "kind"
// @Schemes
// @Description Returns a filtered list of assets3d
// @Tags assets3d
// @Accept json
// @Produce json
// @Param queryParams assets_3d.apiGetAssets3d.QueryParams
// @Success 200 {object} dto.Asset3d
// @Success 400 {object} api.HTTPError
// @Router /api/v4/assets-3d [get]
func (a *Assets3d) apiGetAssets3d(c *gin.Context) {
	type QueryParams struct {
		kind string `form:"kind" json:"kind"`
	}
	queryParams := QueryParams{}

	if err := c.ShouldBind(&queryParams); err != nil {
		err = errors.WithMessage(err, "Assets3d: apiGetAssets3d: failed to bind query parameters")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, a.log)
	}

	var a3dMap map[uuid.UUID]universe.Asset3d
	predicateFn := func(asset3dID uuid.UUID, asset3d universe.Asset3d) bool {
		meta := asset3d.GetMeta()
		kind := utils.GetFromAnyMap(*meta, "kind", "")
		return kind == queryParams.kind
	}

	if queryParams.kind == "" {
		a3dMap = a.GetAssets3d()
	} else {
		a3dMap = a.FilterAssets3d(predicateFn)
	}

	assets := make([]*dto.Asset3d, 0, len(a3dMap))

	for i := range a3dMap {
		asset := a3dMap[i].GetEntry()

		assetDTO := &dto.Asset3d{
			ID:        asset.Asset3dID.String(),
			CreatedAt: asset.CreatedAt.String(),
			UpdatedAt: asset.UpdatedAt.String(),
		}
		assets = append(assets, assetDTO)
	}

	c.JSON(http.StatusOK, assets)
}

// @Summary Add a list of assets3d
// @Schemes
// @Description Creates assets3d with the given input
// @Tags assets3d
// @Accept json
// @Produce json
// @Param request body assets_3d.apiAddAssets3d.InQuery
// @Success 200 {object} nil
// @Success 400	{object} api.HTTPError
// @Success 500 {object} api.HTTPError
// @Router /api/v4/assets3d [post]
func (a *Assets3d) apiAddAssets3d(c *gin.Context) {
	type InQuery struct {
		assets3dIDs []string `form:"assets3d_ids[]" binding:"required"`
	}
	inQuery := InQuery{}

	if err := c.ShouldBindJSON(&inQuery); err != nil {
		err = errors.WithMessage(err, "Assets3d: apiAddAssets3d: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, a.log)
		return
	}

	addAssets3d := make([]universe.Asset3d, 0, len(inQuery.assets3dIDs))
	for i := range inQuery.assets3dIDs {
		assetID, err := uuid.Parse(inQuery.assets3dIDs[i])
		if err != nil {
			err = errors.WithMessage(err, "Assets3d: apiAddAssets3d: failed to parse uuid")
			api.AbortRequest(c, http.StatusInternalServerError, "failed_to_parse_uuid", err, a.log)
		}
		newAsset, err := a.CreateAsset3d(assetID)
		if err != nil {
			err = errors.WithMessage(err, "Assets3d: apiAddAssets3d: failed to create asset3d from input")
			api.AbortRequest(c, http.StatusInternalServerError, "failed_to_create_asset3d", err, a.log)
			return
		}

		addAssets3d = append(addAssets3d, newAsset)
	}

	if err := a.AddAssets3d(addAssets3d, true); err != nil {
		err = errors.WithMessage(err, "Assets3d: apiAddAssets3d: failed to add assets3d")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_add_assets3d", err, a.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Remove assets3d by ids
// @Schemes
// @Description Removes assets3d by list of ids
// @Tags assets3d
// @Accept json
// @Produce json
// @Param request body assets_3d.apiRemoveAssets3dByIDs.InQuery
// @Success 200 {object} nil
// @Success 400 {object} api.HTTPError
// @Success 500 {object} api.HTTPError
// @Router /api/v4/assets3d [delete]
func (a *Assets3d) apiRemoveAssets3dByIDs(c *gin.Context) {
	type InQuery struct {
		assets3dIDs []string `form:"assets3d_ids[]" binding:"required"`
	}
	inQuery := InQuery{}

	if err := c.ShouldBindJSON(&inQuery); err != nil {
		err = errors.WithMessage(err, "Assets3d: apiRemoveAssets3dByIDs: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, a.log)
		return
	}

	uids := make([]uuid.UUID, 0, len(inQuery.assets3dIDs))
	for i := range inQuery.assets3dIDs {
		uid, err := uuid.Parse(inQuery.assets3dIDs[i])
		if err != nil {
			err = errors.WithMessage(err, "Assets3d: apiRemoveAssets3dByIDs: failed to parse uuid")
			api.AbortRequest(c, http.StatusInternalServerError, "invalid_uuid_parse", err, a.log)
			return
		}
		uids[i] = uid
	}

	if err := a.RemoveAssets3dByIDs(uids, true); err != nil {
		err = errors.WithMessage(err, "Assets3d: apiRemoveAssets3dByIDs: failed to remove assets3d")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_remove_assets3d_by_ids", err, a.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Get assets3d options
// @Schemes
// @Description Returns list of options of assets3d
// @Tags assets3d
// @Accept json
// @Produce json
// @Param request body assets_3d.apiGetAssets3dOptions.InQuery
// @Success 200 {object} dto.Assets3dOptions
// @Success 400 {object} api.HTTPError
// @Success 500 {object} api.HTTPError
// @Router /api/v4/assets3d/options [get]
func (a *Assets3d) apiGetAssets3dOptions(c *gin.Context) {
	type InQuery struct {
		Assets3dIDs []string `form:"assets3d_ids[]" binding:"required"`
	}
	inQuery := InQuery{}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Assets3d: apiGetAssets3dOptions: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, a.log)
		return
	}

	out := make(dto.Assets3dOptions, len(inQuery.Assets3dIDs))

	for i := range inQuery.Assets3dIDs {
		asset3dID, err := uuid.Parse(inQuery.Assets3dIDs[i])
		if err != nil {
			err := errors.WithMessagef(err, "Assets3d: apiGetAssets3dOptions: failed to parse uuid")
			api.AbortRequest(c, http.StatusBadRequest, "invalid_asset3d_uuid", err, a.log)
			return
		}

		gotAsset3d, ok := a.GetAsset3d(asset3dID)
		if !ok {
			err = errors.WithMessage(err, "Assets3d: apiGetAsset3dOptions: failed to get asset3d")
			api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_asset3d", err, a.log)
			return
		}

		out[asset3dID] = gotAsset3d.GetOptions()
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Get assets3d meta
// @Schemes
// @Description Returns a list of meta of assets3d
// @Tags assets3d
// @Accept json
// @Produce json
// @Param request body assets_3d.apiGetAssets3dMeta.InQuery
// @Success 200 {object} dto.Assets3dMeta
// @Success 400 {object} api.HTTPError
// @Success 500 {object} api.HTTPError
// @Router /api/v4/assets3d/meta [get]
func (a *Assets3d) apiGetAssets3dMeta(c *gin.Context) {
	type InQuery struct {
		Assets3dIDs []string `form:"assets3d_ids[]" binding:"required"`
	}
	inQuery := InQuery{}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Assets3d: apiGetAssets3dMeta: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, a.log)
		return
	}

	out := make(dto.Assets3dMeta, len(inQuery.Assets3dIDs))

	for i := range inQuery.Assets3dIDs {
		asset3dID, err := uuid.Parse(inQuery.Assets3dIDs[i])
		if err != nil {
			err := errors.WithMessagef(err, "Assets3d: apiGetAssets3dMeta: failed to parse uuid")
			api.AbortRequest(c, http.StatusBadRequest, "invalid_asset3d_uuid", err, a.log)
			return
		}

		gotAsset3d, ok := a.GetAsset3d(asset3dID)
		if !ok {
			err = errors.WithMessage(err, "Assets3d: apiGetAsset3dMeta: failed to get asset3d")
			api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_asset3d", err, a.log)
			return
		}

		out[asset3dID] = gotAsset3d.GetMeta()
	}

	c.JSON(http.StatusOK, out)
}
