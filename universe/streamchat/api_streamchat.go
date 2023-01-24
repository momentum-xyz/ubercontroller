package streamchat

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
	"github.com/pkg/errors"
)

// @Summary Get a authorization token
// @Description Request a chat authorization token for current user and given object (world or space)
// @Description The user is required to be connected to the world/object.
// @Description This also automatically joins the user as a member to the channel, so the join endpoint does not have to be called.
// @Tags chat
// @Accept json
// @Produce json
// @Param spaceID path string true "World or object ID"
// @Success 200 {object} streamchat.apiChannelToken.Response
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/streamchat/{spaceID}/token [post]
func (s *StreamChat) apiChannelToken(c *gin.Context) {
	space, user, err := s.getRequestContextObjects(c)
	if err != nil {
		// TODO: better error handling and response
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request", err, s.log)
		return
	}

	channel, err := s.GetChannel(c, space)
	if err != nil {
		err = errors.WithMessage(err, "Streamchat: failed to get channel")
		api.AbortRequest(c, http.StatusInternalServerError, "get_channel_failed", err, s.log)
		return
	}

	token, err := s.GetToken(c, user)
	if err != nil {
		err = errors.WithMessage(err, "Streamchat: failed to get token")
		api.AbortRequest(c, http.StatusInternalServerError, "get_token_failed", err, s.log)
		return
	}

	if err := s.MakeMember(c, channel, user); err != nil {
		err = errors.WithMessage(err, "Streamchat: failed to make user a member of channel")
		api.AbortRequest(c, http.StatusInternalServerError, "channel_member_failed", err, s.log)
	}

	type Response struct {
		Channel     string `json:"channel"`
		ChannelType string `json:"channel_type"`
		Token       string `json:"token"`
	}
	response := &Response{
		Channel:     channel.ID,
		ChannelType: MomentumChannelType,
		Token:       token,
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Join a chat channel.
// @Description Join the chat channel (as a member) for the given world or object ID.
// @Tags chat
// @Accept json
// @Produce json
// @Param spaceID path string true "World or object ID"
// @Success 204 ""
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/streamchat/{spaceID}/join [post]
func (s *StreamChat) apiChannelJoin(c *gin.Context) {
	space, user, err := s.getRequestContextObjects(c)
	if err != nil {
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request", err, s.log)
		return
	}

	channel, err := s.GetChannel(c, space)
	if err != nil {
		err = errors.WithMessage(err, "Streamchat: failed to get channel")
		api.AbortRequest(c, http.StatusInternalServerError, "get_channel_failed", err, s.log)
		return
	}

	if err := s.MakeMember(c, channel, user); err != nil {
		api.AbortRequest(c, http.StatusInternalServerError, "channel_member_failed", err, s.log)
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Leave a chat channel.
// @Description Leave the chat channel (as a member) for the given world or object ID.
// @Tags chat
// @Accept json
// @Produce json
// @Param spaceID path string true "World or object ID"
// @Success 204 ""
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/streamchat/{spaceID}/leave [post]
func (s *StreamChat) apiChannelLeave(c *gin.Context) {
	space, user, err := s.getRequestContextObjects(c)
	if err != nil {
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request", err, s.log)
		return
	}

	channel, err := s.GetChannel(c, space)
	if err != nil {
		err = errors.WithMessage(err, "Streamchat: failed to get channel")
		api.AbortRequest(c, http.StatusInternalServerError, "get_channel_failed", err, s.log)
		return
	}

	if err := s.RemoveMember(c, channel, user); err != nil {
		api.AbortRequest(c, http.StatusInternalServerError, "channel_member_failed", err, s.log)
		return
	}
	c.Status(http.StatusNoContent)

}

// Get the common objects for these api requests
// TODO: put these in the actual context in shared middleware?
func (s *StreamChat) getRequestContextObjects(c *gin.Context) (space universe.Object, user universe.User, err error) {
	spaceID := c.Param("spaceID")
	space, err = s.getSpace(spaceID)
	if err != nil {
		return
	}

	user, err = s.getUserFromContext(c, space)
	if err != nil {
		return
	}
	return

}

// Get space by UUID string.
func (s *StreamChat) getSpace(id string) (universe.Object, error) {
	spaceID, err := uuid.Parse(id)
	if err != nil {
		err := errors.WithMessagef(err, "Failed to parse ID %s", id)
		return nil, err
	}
	space, ok := s.node.GetObjectFromAllObjects(spaceID)
	if !ok {
		err := errors.Errorf("Object not found: %s", spaceID)
		return nil, err
	}
	return space, nil
}

// Resolve the user object from the request context and the given space.
func (s *StreamChat) getUserFromContext(ctx *gin.Context, space universe.Object) (universe.User, error) {

	userID, err := api.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "Streamchat: failed to get userID from context")
	}
	// The user needs to be 'in' the space (a.k.a world or object)
	user, ok := space.GetUser(userID, false)
	if !ok {
		return nil, fmt.Errorf("Streamchat: User %s not found in %s", userID, space.GetID())
	}
	return user, nil
}
