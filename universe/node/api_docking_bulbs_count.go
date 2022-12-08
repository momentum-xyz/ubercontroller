package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe/common/api"
)

// @Summary Get count of docking bulbs in the space
// @Schemes
// @Description Returns count of docking bulbs in the space
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Success 202 {object} node.apiGetSpaceDockingBulbsCount.Out
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/docking-bulbs-count [get]
func (n *Node) apiGetSpaceDockingBulbsCount(c *gin.Context) {
	type Out struct {
		Count int `json:"count"`
	}

	out := Out{
		Count: 0,
	}

	types := n.GetSpaceTypes().GetSpaceTypes()
	var bulbTypeID uuid.UUID
	for _, v := range types {
		if v.GetName() == "Docking bulb" {
			bulbTypeID = v.GetID()
			break
		}
	}

	if bulbTypeID == uuid.Nil {
		err := errors.New("Node: apiGetSpaceDockingBulbsCount: failed to find spaceType by name 'Docking bulb'")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceDockingBulbsCount: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	space, ok := n.GetSpaceFromAllSpaces(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiGetSpaceDockingBulbsCount: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	spaces := space.GetSpaces(true)
	for _, v := range spaces {
		st := v.GetSpaceType()
		if st.GetID() == bulbTypeID {
			out.Count++
		}
	}

	c.JSON(http.StatusOK, out)
}
