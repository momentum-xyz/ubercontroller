package worlds

import (
	"fmt"
	"math/big"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

// @Summary Get NFT metadata.
// @Description Returns ERC721 metadata.
// @Tags nfts
// @Param nftID path string true "NFT token ID"
// @Success 200 {object} dto.WorldNFTMeta
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/nft/{nftID} [get]
func (w *Worlds) apiNFTMetaData(c *gin.Context) {
	nftID := c.Param("nftID")
	tokenID := new(big.Int)
	tokenID, ok := tokenID.SetString(nftID, 10)
	if !ok {
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request", fmt.Errorf("failed to parse nft ID"), w.log)
		return
	}
	odysseyID, err := umid.FromBytes(tokenID.FillBytes(make([]byte, 16)))
	if err != nil {
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request", fmt.Errorf("failed to parse UMID %s: %w", tokenID, err), w.log)
		return
	}

	var response *dto.WorldNFTMeta
	world, ok := w.GetWorld(odysseyID)
	if !ok {
		response = w.nftMetadataDefault(odysseyID)
	} else {
		response, err = w.nftMetadata(world)
		if err != nil {
			api.AbortRequest(c, http.StatusBadRequest, "invalid_request", fmt.Errorf("failed to get world metadata: %w", err), w.log)
			return
		}
	}
	c.JSON(http.StatusOK, response)
}

func (w *Worlds) nftMetadataDefault(worldID umid.UMID) *dto.WorldNFTMeta {
	seqID := utils.UMIDToSEQ(worldID)
	objectName := "Odyssey#" + strconv.FormatUint(seqID, 10)
	return &dto.WorldNFTMeta{
		Name:        objectName,
		Description: "",
		Image:       w.nftImageDefault(),
		ExternalURL: w.nftExternalURLDefault(),
		Attributes:  nil,
	}
}

func (w *Worlds) nftMetadata(world universe.World) (*dto.WorldNFTMeta, error) {
	attributes := make([]dto.NFTAttributes, 0)

	worldNr := utils.UMIDToSEQ(world.GetID())
	if worldNr <= 100 { // Will we get more ranges like this?
		attribute := dto.NFTAttributes{
			TraitType: "type",
			Value:     "origin",
		}
		attributes = append(attributes, attribute)
	}

	pluginID := universe.GetSystemPluginID()
	attributeID := entry.NewAttributeID(pluginID, universe.ReservedAttributes.Object.WorldAvatar.Name)
	imageValue, ok := world.GetObjectAttributes().GetValue(attributeID)
	if !ok {
		w.log.Debugf("NFT metadata: could not get image attribute for world")
		// Ignore ok return here, no image attr, use a default fallback.
	}

	var worldAvatarHash string
	if imageValue != nil {
		worldAvatarHash = utils.GetFromAnyMap(*imageValue, universe.ReservedAttributes.Object.WorldAvatar.Key, "")
	} else {
		worldAvatarHash = w.nftImageDefault()
	}

	return &dto.WorldNFTMeta{
		Name:        world.GetName(),
		Description: world.GetDescription(),
		Image:       w.nftImage(worldAvatarHash),
		ExternalURL: w.nftExternalURL(world),
		Attributes:  attributes,
	}, nil
}

func (w *Worlds) nftImageDefault() string {
	return w.nftImage("bd6563cc9fceac3e1ed6fcad752c902d") // TODO: node level config
}

func (w *Worlds) nftImage(imgHash string) string {
	return w.cfg.Settings.FrontendURL + "/api/v3/render/get/" + imgHash
}

func (w *Worlds) nftExternalURLDefault() string {
	return w.cfg.Settings.FrontendURL
}

func (w *Worlds) nftExternalURL(world universe.World) string {
	return w.cfg.Settings.FrontendURL + "/odyssey/" + world.GetID().String()
}
