package assets_3d

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/api"
	"github.com/momentum-xyz/ubercontroller/universe/api/dto"
)

func (a *Assets3d) apiGetAsset3d(c *gin.Context) {
	asset3dID, err := uuid.Parse(c.Param("asset3dID"))
	if err != nil {
		err = errors.WithMessage(err, "Assets3d: apiGetAsset3d: failed to parse asset3dID")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_set_asset3d_uuid", err, a.log)
		return
	}

	gAsset3d, ok := a.GetAsset3d(asset3dID)
	if !ok {
		err = errors.WithMessage(err, "Assets3d: apiGetAsset3d: failed to get asset3d")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_asset3d", err, a.log)
		return
	}

	out := gAsset3d.GetEntry()

	outDTO := dto.Asset3d{
		ID:        out.Asset3dID.String(),
		CreatedAt: out.CreatedAt.String(),
		UpdatedAt: out.UpdatedAt.String(),
	}

	c.JSON(http.StatusOK, outDTO)
}

func (a *Assets3d) apiGetAssets3d(c *gin.Context) {
	// TODO: rework this in a different method
	// or in a more generic way
	//
	// This is currently used as a poor man's filter
	// for assets with `Meta = {"kind":"skybox"}`
	//
	// example "?kind=skybox` should return "skybox"
	queryParams := struct {
		kind string `form:"kind" json:"kind"`
	}{}

	if err := c.ShouldBind(&queryParams); err != nil {
		err = errors.WithMessage(err, "Assets3d: apiGetAssets3d: failed to bind query parameters")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, a.log)
	}

	a3dMap := a.GetAssets3d()
	assets := make([]*dto.Asset3d, 0, len(a3dMap))

	if queryParams.kind == "" {
		for _, el := range a3dMap {
			asset := el.GetEntry()

			assetDTO := &dto.Asset3d{
				ID:        asset.Asset3dID.String(),
				CreatedAt: asset.CreatedAt.String(),
				UpdatedAt: asset.UpdatedAt.String(),
			}
			assets = append(assets, assetDTO)
		}
	} else {
		for _, el := range a3dMap {
			asset := el.GetEntry()
			if asset.Meta != nil {
				if (*asset.Meta)["kind"] == queryParams.kind {

					assetDTO := &dto.Asset3d{
						ID:        asset.Asset3dID.String(),
						CreatedAt: asset.CreatedAt.String(),
						UpdatedAt: asset.UpdatedAt.String(),
					}
					assets = append(assets, assetDTO)
				}
			}
		}
	}

	c.JSON(http.StatusOK, assets)
}

func (a *Assets3d) apiAddAssets3d(c *gin.Context) {
	inBody := struct {
		assets3dIDs []string `json:"assets3d"`
	}{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Assets3d: apiAddAssets3d: failed to bind json")
		api.AbortRequest(c, http.StatusInternalServerError, "invalid_request_query", err, a.log)
		return
	}

	addAssets3d := make([]universe.Asset3d, 0, len(inBody.assets3dIDs))
	for i := range inBody.assets3dIDs {
		assetID, err := uuid.Parse(inBody.assets3dIDs[i])
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

func (a *Assets3d) apiRemoveAssets3dByIDs(c *gin.Context) {
	in := struct {
		ids []string `form:"ids" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&in); err != nil {
		err = errors.WithMessage(err, "Assets3d: apiRemoveAssets3dByIDs: failed to bind json")
		api.AbortRequest(c, http.StatusInternalServerError, "invalid_request_query", err, a.log)
		return
	}

	uids := make([]uuid.UUID, 0, len(in.ids))
	for _, id := range in.ids {
		uid, err := uuid.Parse(id)
		if err != nil {
			err = errors.WithMessage(err, "Assets3d: apiRemoveAssets3dByIDs: failed to parse uuid")
			api.AbortRequest(c, http.StatusInternalServerError, "invalid_uuid_parse", err, a.log)
			return
		}
		uids = append(uids, uid)
	}

	if err := a.RemoveAssets3dByIDs(uids, true); err != nil {
		err = errors.WithMessage(err, "Assets3d: apiRemoveAssets3dByIDs: failed to remove assets3d")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_remove_assets3d_by_ids", err, a.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (a *Assets3d) apiGetAsset3dOptions(c *gin.Context) {

}

func (a *Assets3d) apiGetAsset3dMeta(c *gin.Context) {

}
