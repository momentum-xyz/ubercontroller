package assets_3d

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/merge"
)

// @Summary Get 3d assets
// @Schemes
// @Description Returns a filtered list of 3d assets
// @Tags assets3d
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Param query query assets_3d.apiGetAssets3d.InQuery true "query params"
// @Success 200 {array} dto.Asset3d
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/assets-3d/{object_id} [get]
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

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Assets3d: apiUploadAsset3d: failed to get user umid from context")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, a.log)
		return
	}

	var a3dMap map[universe.AssetUserIDPair]universe.UserAsset3d
	predicateFn := func(assetUserID universe.AssetUserIDPair, userAsset3d universe.UserAsset3d) bool {
		var category string

		if userAsset3d.IsPrivate() && userAsset3d.GetUserID() != userID {
			return false
		}

		asset3d := userAsset3d.GetAsset3d()
		if asset3d == nil {
			fmt.Printf("asset3d is nil, strange! %+v\n", userAsset3d.GetAssetUserIDPair())
			return false
		}

		meta := (*asset3d).GetMeta()
		if meta == nil {
			return false
		}

		category = utils.GetFromAnyMap(*meta, "category", "")
		return category == inQuery.Category
	}

	if inQuery.Category == "" {
		a3dMap = a.GetUserAssets3d()
	} else {
		a3dMap = a.FilterUserAssets3d(predicateFn)
	}

	assets := make([]*dto.Asset3d, 0, len(a3dMap))

	for i := range a3dMap {
		asset := a3dMap[i].GetEntry()
		baseAsset3d := a3dMap[i].GetAsset3d()
		meta := asset.Meta
		baseMeta := (*baseAsset3d).GetMeta()

		combinedMeta, err := merge.Auto(baseMeta, meta)
		if err != nil {
			err = errors.WithMessagef(err, "Assets3d: apiGetAssets3d: failed to merge meta")
			api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, a.log)
			return
		}

		assetDTO := &dto.Asset3d{
			ID:        asset.Asset3dID.String(),
			UserID:    asset.UserID.String(),
			Meta:      combinedMeta,
			Private:   asset.Private,
			CreatedAt: asset.CreatedAt.Format(time.RFC3339),
			UpdatedAt: asset.UpdatedAt.Format(time.RFC3339),
		}
		assets = append(assets, assetDTO)
	}

	c.JSON(http.StatusOK, assets)
}

// @Summary Uploads a 3d asset
// @Schemes
// @Description Uploads a 3d asset to the media manager
// @Tags assets3d
// @Accept multipart/form-data
// @Produce json
// @Param object_id path string true "Object UMID"
// @Success 202 {object} dto.Asset3d
// @Failure 400	{object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/assets-3d/{object_id}/upload [post]
// TODO: swag doc for multipart, it does not get *multipart.FileHeader
func (a *Assets3d) apiUploadAsset3d(c *gin.Context) {
	type InBody struct {
		File        *multipart.FileHeader `form:"asset"`
		Name        string                `form:"name"`
		PreviewHash *string               `form:"preview_hash"`
		Private     *bool                 `form:"is_private"`
	}
	var request InBody
	if err := c.ShouldBind(&request); err != nil {
		err := errors.WithMessage(err, "Assets3d: apiUploadAsset3d: failed to read request")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_read", err, a.log)
		return
	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Assets3d: apiUploadAsset3d: failed to get user umid from context")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, a.log)
		return
	}

	assetFile := request.File
	if assetFile == nil {
		api.AbortRequest(
			c, http.StatusBadRequest, "failed_to_open", errors.New("Assets3d: apiUploadAsset3d: no file in request"),
			a.log,
		)
		return
	}

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

	fmt.Printf("Upload to mm response: %+v\n", response)

	assetID, err := umid.Parse(response.Hash)
	if err != nil {
		err := errors.WithMessage(err, "Assets3d: apiUploadAsset3d: failed to parse hash to uuid")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_parse_hash", err, a.log)
		return
	}

	baseAsset, err, isNewInstance := a.CreateAsset3d(assetID)
	if err != nil {
		err := errors.WithMessage(err, "Assets3d: apiUploadAsset3d: failed to get or create asset3d")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_or_create_asset3d", err, a.log)
		return
	}

	isPrivate := false
	if request.Private != nil {
		isPrivate = *request.Private
	}

	newUserAsset, err := a.CreateUserAsset3d(assetID, userID, isPrivate)
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
	baseMeta := entry.Asset3dMeta{
		"type":     dto.GLTFAsset3dType,
		"category": "custom",
	}
	meta := entry.Asset3dMeta{
		"name": name,
	}

	if request.PreviewHash != nil {
		meta["preview_hash"] = request.PreviewHash
	}

	if isNewInstance {
		if err := baseAsset.SetMeta(&baseMeta, false); err != nil {
			err = errors.WithMessage(err, "Assets3d: apiUploadAsset3d: failed to set meta")
			api.AbortRequest(c, http.StatusInternalServerError, "failed_to_set_meta", err, a.log)
			return
		}

		// it's added in Create but not saved to db
		if err := a.AddAsset3d(baseAsset, true); err != nil {
			err = errors.WithMessage(err, "Assets3d: apiUploadAsset3d: failed to add assets3d")
			api.AbortRequest(c, http.StatusInternalServerError, "failed_to_add_asset3d", err, a.log)
			return
		}
	}

	if err := newUserAsset.SetMeta(&meta, true); err != nil {
		err = errors.WithMessage(err, "Assets3d: apiUploadAsset3d: failed to set meta")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_set_meta", err, a.log)
		return
	}

	// it's added in Create but not saved to db
	if err := a.AddUserAsset3d(newUserAsset, true); err != nil {
		err = errors.WithMessage(err, "Assets3d: apiUploadAsset3d: failed to add user assets3d")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_add_user_asset3d", err, a.log)
		return
	}

	combinedMeta, err := merge.Auto(&baseMeta, &meta)
	if err != nil {
		err = errors.WithMessagef(err, "Assets3d: apiGetAssets3d: failed to merge meta")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, a.log)
		return
	}

	out := dto.Asset3d{
		ID:      newUserAsset.GetAssetID().String(),
		UserID:  newUserAsset.GetUserID().String(),
		Private: newUserAsset.IsPrivate(),
		Meta:    combinedMeta,
	}

	c.JSON(http.StatusAccepted, out)
}

// @Summary Get 3d assets options
// @Schemes
// @Description Returns list of 3d assets options
// @Tags assets3d
// @Accept json
// @Produce json
// @Param query query assets_3d.apiGetAssets3dOptions.InQuery true "query params"
// @Success 200 {object} dto.Assets3dOptions
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/assets-3d/{object_id}/options [get]
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
		asset3dID, err := umid.Parse(inQuery.Assets3dIDs[i])
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

// @Summary Delete a 3d asset by its umid
// @Schemes
// @Description Deletes 3d asset by its umid
// @Tags assets3d
// @Accept json
// @Produce json
// @Success 200 {object} nil
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/assets-3d/{object_id}/{asset3d_id} [delete]
func (a *Assets3d) apiRemoveAsset3dByID(c *gin.Context) {
	uid, err := umid.Parse(c.Param("asset3dID"))
	if err != nil {
		err := errors.WithMessage(err, "Assets3d: apiRemoveAsset3dByID: failed to parse uuid")
		api.AbortRequest(c, http.StatusInternalServerError, "invalid_uuid_parse", err, a.log)
		return
	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Assets3d: apiRemoveAsset3dByID: failed to get user umid from context")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, a.log)
		return
	}

	assetUserID := universe.AssetUserIDPair{
		AssetID: uid,
		UserID:  userID,
	}

	removed, err := a.RemoveUserAsset3dByID(assetUserID, true)
	if err != nil {
		err := errors.WithMessage(err, "Assets3d: apiRemoveAsset3dByID: failed to remove asset3d")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_remove_asset3d_by_id", err, a.log)
		return
	}
	if !removed {
		err := errors.WithMessage(err, "Assets3d: apiRemoveAsset3dByID: failed to remove asset3d")
		api.AbortRequest(c, http.StatusNotFound, "asset3d_not_removed", err, a.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Update 3d asset meta by its umid
// @Schemes
// @Description Update 3d asset meta by its umid
// @Tags assets3d
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Param asset3d_id path string true "Asset 3D UMID"
// @Param body body assets_3d.apiUpdateAsset3dByID.InBody true "body params"
// @Success 200 {object} dto.Asset3d
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/assets-3d/{object_id}/{asset3d_id} [patch]
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

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err = errors.WithMessage(err, "Assets3d: apiUpdateAsset3dByID: failed to get user umid from context")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, a.log)
		return
	}

	asset3dID, err := umid.Parse(c.Param("asset3dID"))
	if err != nil {
		err = errors.WithMessage(err, "Assets3d: apiUpdateAsset3dByID: failed to parse uuid")
		api.AbortRequest(c, http.StatusInternalServerError, "invalid_uuid_parse", err, a.log)
		return
	}

	userAsset3d, ok := a.GetUserAsset3d(asset3dID, userID)
	if !ok {
		err = errors.WithMessagef(err, "Assets3d: apiUpdateAsset3dByID: asset3d not found: %s", asset3dID)
		api.AbortRequest(c, http.StatusNotFound, "not_found", err, a.log)
		return
	}

	oldMeta := userAsset3d.GetMeta()
	newMeta, err := merge.Auto[entry.Asset3dMeta](&inBody.Meta, oldMeta)
	if err != nil {
		err = errors.WithMessagef(err, "Assets3d: apiUpdateAsset3dByID: failed to merge meta")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, a.log)
		return
	}

	if err := userAsset3d.SetMeta(newMeta, true); err != nil {
		err = errors.WithMessagef(err, "Assets3d: apiUpdateAsset3dByID: failed to set meta")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, a.log)
		return
	}

	c.JSON(http.StatusOK, userAsset3d.GetEntry())
}
