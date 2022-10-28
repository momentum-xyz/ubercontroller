package assets_3d

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
)

func (a *Assets3d) apiCreateAsset3d(c *gin.Context) {
	in := struct {
		Meta    *entry.Asset3dMeta    `json:"meta"`
		Options *entry.Asset3dOptions `json:"options"`
	}{}

	if err := c.ShouldBindJSON(&in); err != nil {
		err = errors.WithMessage(err, "Assets3d: apiCreateAsset3d: failed to bind json")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "",
		})
		return
	}

	asset3dID, err := uuid.Parse(c.Param("asset3dID"))
	if err != nil {
		err = errors.WithMessage(err, "Assets3d: apiCreateAsset3d: failed to parse asset3dID")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "",
		})
		return
	}

	crAsset3d, err := a.CreateAsset3d(asset3dID)
	if err != nil {
		err = errors.WithMessage(err, "Assets3d: apiCreateAsset3d: failed to create asset3d")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "",
		})
		return
	}

	err = crAsset3d.SetMeta(in.Meta, false)
	if err != nil {
		err = errors.WithMessage(err, "Assets3d: apiCreateAsset3d: failed to set asset3d meta")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "",
		})
		return
	}

	// TODO: set the asset options with a proper modify func - is it currently needed
	// as there are no such options so far?
	// modFn := func(ops *entry.Asset3dOptions) (*entry.Asset3dOptions, error) {

	// }

	// err := crAsset3d.SetOptions(modFn(in.Options), false)
	// if err != nil {
	// 	err = errors.WithMessage(err, "Assets3d: apiCreateAsset3d: failed to set asset3d options")
	// 	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
	// 		"message": "",
	// 	})
	// 	return
	// }

	out := crAsset3d.GetEntry()
	c.JSON(http.StatusOK, out)
}

func (a *Assets3d) apiGetAsset3d(c *gin.Context) {
	asset3dID, err := uuid.Parse(c.Param("asset3dID"))
	if err != nil {
		err = errors.WithMessage(err, "Assets3d: apiGetAsset3d: failed to parse asset3dID")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "",
		})
		return
	}

	gAsset3d, ok := a.GetAsset3d(asset3dID)
	if !ok {
		err = errors.WithMessage(err, "Assets3d: apiGetAsset3d: failed to get asset3d")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "",
		})
		return
	}

	outBody := gAsset3d.GetEntry()
	c.JSON(http.StatusOK, outBody)
}

func (a *Assets3d) apiGetAssets3d(c *gin.Context) {
	// TODO: add filter; example "?kind=skybox"` for assets with `Meta = {"kind":"skybox"}`
	a3dMap := a.GetAssets3d()

	assets := make([]*entry.Asset3d, len(a3dMap))
	for _, el := range a3dMap {
		asset := el.GetEntry()
		assets = append(assets, asset)
	}

	c.JSON(http.StatusOK, assets)
}

func (a *Assets3d) apiAddAsset3d(c *gin.Context) {
	in := struct {
		asset3d entry.Asset3d `json:"asset3d"`
	}{}

	if err := c.ShouldBindJSON(&in); err != nil {
		err = errors.WithMessage(err, "Assets3d: apiAddAsset3d: failed to bind json")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "",
		})
		return
	}

	newAsset, err := a.CreateAsset3d(in.asset3d.Asset3dID)
	if err != nil {
		err = errors.WithMessage(err, "Assets3d: apiAddAsset3d: failed to create asset3d from input")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "",
		})
		return
	}

	if err := newAsset.SetMeta(in.asset3d.Meta, false); err != nil {
		err = errors.WithMessage(err, "Assets3d: apiAddAsset3d: failed to set asset3d meta")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "",
		})
		return

	}

	if err := a.AddAsset3d(newAsset, false); err != nil {
		err = errors.WithMessage(err, "Assets3d: apiAddAsset3d: failed to add asset3d")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "",
		})
		return
	}

	c.JSON(http.StatusOK, in)
}

func (a *Assets3d) apiAddAssets3d(c *gin.Context) {
	in := struct {
		assets3d []entry.Asset3d `json:"assets3d"`
	}{}

	if err := c.ShouldBindJSON(&in); err != nil {
		err = errors.WithMessage(err, "Assets3d: apiAddAssets3d: failed to bind json")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "",
		})
		return
	}

	addAssets3d := make([]universe.Asset3d, len(in.assets3d))
	for _, asset3d := range in.assets3d {
		newAsset, err := a.CreateAsset3d(asset3d.Asset3dID)
		if err != nil {
			err = errors.WithMessage(err, "Assets3d: apiAddAssets3d: failed to create asset3d from input")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "",
			})
			return
		}

		if err := newAsset.SetMeta(asset3d.Meta, false); err != nil {
			err = errors.WithMessage(err, "Assets3d: apiAddAssets3d: failed to set asset3d meta")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "",
			})
			return

		}
		addAssets3d = append(addAssets3d, newAsset)
	}

	if err := a.AddAssets3d(addAssets3d, false); err != nil {
		err = errors.WithMessage(err, "Assets3d: apiAddAssets3d: failed to add assets3d")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "",
		})
		return
	}

	c.JSON(http.StatusOK, in)
}

func (a *Assets3d) apiRemoveAsset3d(c *gin.Context) {
	asset3dID, err := uuid.Parse(c.Param("asset3dID"))
	if err != nil {
		err = errors.WithMessage(err, "Assets3d: apiRemoveAsset3d: failed to parse asset3dID")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "",
		})
		return
	}

	gAsset3d, ok := a.GetAsset3d(asset3dID)
	if !ok {
		err = errors.WithMessage(err, "Assets3d: apiRemoveAsset3d: failed to get asset3d")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "",
		})
		return
	}

	if err := a.RemoveAsset3d(gAsset3d, false); err != nil {
		err = errors.WithMessage(err, "Assets3d: apiRemoveAsset3d: failed to remove asset3d")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "",
		})
		return
	}
	c.JSON(http.StatusOK, gAsset3d)
}

func (a *Assets3d) apiRemoveAssets3d(c *gin.Context) {
	//TODO
}
