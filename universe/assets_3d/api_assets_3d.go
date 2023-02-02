package assets_3d

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
	"github.com/momentum-xyz/ubercontroller/universe/common/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/merge"
)

// @Summary Get 3d assets
// @Schemes
// @Description Returns a filtered list of 3d assets
// @Tags assets3d
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param query query assets_3d.apiGetAssets3d.InQuery true "query params"
// @Success 200 {array} dto.Asset3d
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/assets-3d/{space_id} [get]
func (a *Assets3d) apiGetAssets3d(c *gin.Context) {
	type InQuery struct {
		Category string `form:"category" json:"category"`
	}
	var inQuery InQuery

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err = errors.WithMessage(err, "Assets3d: apiGetAssets3d: failed to bind query parameters")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, a.log)
		return
	}

	var a3dMap map[uuid.UUID]universe.Asset3d
	predicateFn := func(asset3dID uuid.UUID, asset3d universe.Asset3d) bool {
		var category string
		meta := asset3d.GetMeta()

		if meta == nil {
			return false
		}
		category = utils.GetFromAnyMap(*meta, "category", "")
		return category == inQuery.Category
	}

	if inQuery.Category == "" {
		a3dMap = a.GetAssets3d()
	} else {
		a3dMap = a.FilterAssets3d(predicateFn)
	}

	assets := make([]*dto.Asset3d, 0, len(a3dMap))

	for i := range a3dMap {
		asset := a3dMap[i].GetEntry()

		assetDTO := &dto.Asset3d{
			ID:        asset.Asset3dID.String(),
			Meta:      asset.Meta,
			CreatedAt: asset.CreatedAt.Format(time.RFC3339),
			UpdatedAt: asset.UpdatedAt.Format(time.RFC3339),
		}
		assets = append(assets, assetDTO)
	}

	c.JSON(http.StatusOK, assets)
}

// @Summary Add 3d assets
// @Schemes
// @Description Creates 3d assets with the given input
// @Tags assets3d
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param body body assets_3d.apiAddAssets3d.InBody true "body params"
// @Success 200 {object} nil
// @Failure 400	{object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/assets-3d/{space_id} [post]
func (a *Assets3d) apiAddAssets3d(c *gin.Context) {
	type InBody struct {
		Assets3dIDs []string `form:"assets3d_ids[]" binding:"required"`
	}
	var inBody InBody

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Assets3d: apiAddAssets3d: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, a.log)
		return
	}

	addAssets3d := make([]universe.Asset3d, 0, len(inBody.Assets3dIDs))
	for i := range inBody.Assets3dIDs {
		assetID, err := uuid.Parse(inBody.Assets3dIDs[i])
		if err != nil {
			err = errors.WithMessage(err, "Assets3d: apiAddAssets3d: failed to parse uuid")
			api.AbortRequest(c, http.StatusInternalServerError, "failed_to_parse_uuid", err, a.log)
			return
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

// @Summary Uploads a 3d asset
// @Schemes
// @Description Uploads a 3d asset to the media manager
// @Tags assets3d
// @Accept multipart/form-data
// @Produce json
// @Param space_id path string true "Space ID"
// @Success 202 {object} dto.Asset3d
// @Failure 400	{object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/assets-3d/{space_id}/upload [post]
// TODO: swag doc for multipart, it does not get *multipart.FileHeader
func (a *Assets3d) apiUploadAsset3d(c *gin.Context) {
	type InBody struct {
		File *multipart.FileHeader `form:"asset"`
		Name string                `form:"name"`
	}
	var request InBody
	if err := c.ShouldBind(&request); err != nil {
		err := errors.WithMessage(err, "Assets3d: apiUploadAsset3d: failed to read request")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_read", err, a.log)
		return
	}
	assetFile := request.File

	openedFile, err := assetFile.Open()
	if err != nil {
		err := errors.WithMessage(err, "Assets3d: apiUploadAsset3d: failed to open file")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_open", err, a.log)
		return
	}

	defer openedFile.Close()

	req, err := http.NewRequest("POST", a.cfg.Common.RenderInternalURL+"/addasset", openedFile)
	if err != nil {
		err := errors.WithMessage(err, "Assets3d: apiUploadAsset3d: failed to create post request")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_create_request", err, a.log)
		return
	}

	req.Header.Set("Content-Type", assetFile.Header.Get("Content-Type"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		err := errors.WithMessage(err, "Assets3d: apiUploadAsset3d: failed to post data to media-manager")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_post_request", err, a.log)
		return
	}

	defer resp.Body.Close()

	response := dto.HashResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		err := errors.WithMessage(err, "Assets3d: apiUploadAsset3d: failed to decode json into response")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_decode", err, a.log)
		return
	}

	assetID, err := uuid.Parse(response.Hash)
	if err != nil {
		err := errors.WithMessage(err, "Assets3d: apiUploadAsset3d: failed to parse hash to uuid")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_parse_hash", err, a.log)
		return
	}

	newAsset, err := a.CreateAsset3d(assetID)
	if err != nil {
		err = errors.WithMessage(err, "Assets3d: apiUploadAsset3d: failed to create asset3d from input")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_create_asset3d", err, a.log)
		return
	}

	var name string
	if request.Name != "" {
		name = request.Name
	} else {
		fileName := assetFile.Filename
		name = fileName[:len(fileName)-len(filepath.Ext(fileName))] // utility function?
	}
	meta := entry.Asset3dMeta{
		"type":     dto.GLTFAsset3dType,
		"category": "custom",
		"name":     name,
	}

	if err := newAsset.SetMeta(&meta, true); err != nil {
		err = errors.WithMessage(err, "Assets3d: apiUploadAsset3d: failed to set meta")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_set_meta", err, a.log)
		return
	}

	if err := a.AddAsset3d(newAsset, true); err != nil {
		err = errors.WithMessage(err, "Assets3d: apiUploadAsset3d: failed to add assets3d")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_add_asset3d", err, a.log)
		return
	}

	out := dto.Asset3d{
		ID:   newAsset.GetID().String(),
		Meta: newAsset.GetMeta(),
	}

	c.JSON(http.StatusAccepted, out)
}

// @Summary Delete 3d assets
// @Schemes
// @Description Deletes 3d assets by list of ids
// @Tags assets3d
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param body body assets_3d.apiRemoveAssets3dByIDs.InBody true "body params"
// @Success 200 {object} nil
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/assets-3d/{space_id} [delete]
func (a *Assets3d) apiRemoveAssets3dByIDs(c *gin.Context) {
	type InBody struct {
		Assets3dIDs []string `form:"assets3d_ids[]" binding:"required"`
	}
	var inBody InBody

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Assets3d: apiRemoveAssets3dByIDs: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, a.log)
		return
	}

	uids := make([]uuid.UUID, 0, len(inBody.Assets3dIDs))
	for i := range inBody.Assets3dIDs {
		uid, err := uuid.Parse(inBody.Assets3dIDs[i])
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

// @Summary Get 3d assets options
// @Schemes
// @Description Returns list of 3d assets options
// @Tags assets3d
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param query query assets_3d.apiGetAssets3dOptions.InQuery true "query params"
// @Success 200 {object} dto.Assets3dOptions
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/{space_id}/assets-3d/options [get]
func (a *Assets3d) apiGetAssets3dOptions(c *gin.Context) {
	type InQuery struct {
		Assets3dIDs []string `form:"assets3d_ids[]" binding:"required"`
	}
	var inQuery InQuery

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
			err = errors.Errorf("Assets3d: apiGetAsset3dOptions: failed to get asset3d")
			api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_asset3d", err, a.log)
			return
		}

		out[asset3dID] = gotAsset3d.GetOptions()
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Get 3d assets meta
// @Schemes
// @Description Returns a list of 3d assets meta
// @Tags assets3d
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param query query assets_3d.apiGetAssets3dMeta.InQuery true "query params"
// @Success 200 {object} dto.Assets3dMeta
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/assets-3d/{space_id}/meta [get]
func (a *Assets3d) apiGetAssets3dMeta(c *gin.Context) {
	type InQuery struct {
		Assets3dIDs []string `form:"assets3d_ids[]" binding:"required"`
	}
	var inQuery InQuery

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
			err = errors.Errorf("Assets3d: apiGetAsset3dMeta: failed to get asset3d")
			api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_asset3d", err, a.log)
			return
		}

		out[asset3dID] = gotAsset3d.GetMeta()
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Delete a 3d asset by its id
// @Schemes
// @Description Deletes 3d asset by its id
// @Tags assets3d
// @Accept json
// @Param space_id path string true "Space ID"
// @Param asset3d_id path string true "Asset 3D ID"
// @Produce json
// @Success 200 {object} nil
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/assets-3d/{space_id}/{asset3d_id} [delete]
func (a *Assets3d) apiRemoveAsset3dByID(c *gin.Context) {
	uid, err := uuid.Parse(c.Param("asset3dID"))
	if err != nil {
		err = errors.WithMessage(err, "Assets3d: apiRemoveAsset3dByID: failed to parse uuid")
		api.AbortRequest(c, http.StatusInternalServerError, "invalid_uuid_parse", err, a.log)
		return
	}

	if err := a.RemoveAsset3dByID(uid, true); err != nil {
		err = errors.WithMessage(err, "Assets3d: apiRemoveAsset3dByID: failed to remove asset3d")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_remove_asset3d_by_id", err, a.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Update 3d asset meta by its id
// @Schemes
// @Description Update 3d asset meta by its id
// @Tags assets3d
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param asset3d_id path string true "Asset 3D ID"
// @Param body body assets_3d.apiUpdateAsset3dByID.InBody true "body params"
// @Success 200 {object} dto.Asset3d
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/assets-3d/{space_id}/{asset3d_id} [patch]
func (a *Assets3d) apiUpdateAsset3dByID(c *gin.Context) {
	type InBody struct {
		Meta entry.Asset3dMeta `json:"meta" binding:"required"`
	}
	var inBody InBody
	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Assets3d: apiUpdateAsset3dByID: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, a.log)
		return
	}

	asset3dID, err := uuid.Parse(c.Param("asset3dID"))
	if err != nil {
		err = errors.WithMessage(err, "Assets3d: apiUpdateAsset3dByID: failed to parse uuid")
		api.AbortRequest(c, http.StatusInternalServerError, "invalid_uuid_parse", err, a.log)
		return
	}

	asset3d, ok := a.GetAssets3d()[asset3dID]
	if !ok {
		err = errors.WithMessagef(err, "Assets3d: apiUpdateAsset3dByID: asset3d not found: %s", asset3dID)
		api.AbortRequest(c, http.StatusNotFound, "not_found", err, a.log)
		return
	}

	oldMeta := asset3d.GetMeta()
	newMeta, err := merge.Auto[entry.Asset3dMeta](&inBody.Meta, oldMeta)
	if err != nil {
		err = errors.WithMessagef(err, "Assets3d: apiUpdateAsset3dByID: failed to merge meta")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, a.log)
		return
	}

	if err := asset3d.SetMeta(newMeta, true); err != nil {
		err = errors.WithMessagef(err, "Assets3d: apiUpdateAsset3dByID: failed to set meta")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, a.log)
		return
	}

	c.JSON(http.StatusOK, asset3d.GetEntry())
}
