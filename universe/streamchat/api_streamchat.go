package streamchat

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/pkg/errors"
)

// @Summary Get a authorization token
// @Description Request a chat authorization token for current user and given object (world or object)
// @Description The user is required to be connected to the world/object.
// @Description This also automatically joins the user as a member to the channel, so the join endpoint does not have to be called.
// @Tags chat
// @Accept json
// @Produce json
// @Param objectID path string true "World or object ID"
// @Success 200 {object} streamchat.apiChannelToken.Response
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/streamchat/{objectID}/token [post]
func (s *StreamChat) apiChannelToken(c *gin.Context) {
	object, user, err := s.getRequestContextObjects(c)
	if err != nil {
		// TODO: better error handling and response
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request", err, s.log)
		return
	}

	channel, err := s.GetChannel(c, object)
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
// @Param objectID path string true "World or object ID"
// @Success 204 ""
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/streamchat/{objectID}/join [post]
func (s *StreamChat) apiChannelJoin(c *gin.Context) {
	object, user, err := s.getRequestContextObjects(c)
	if err != nil {
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request", err, s.log)
		return
	}

	channel, err := s.GetChannel(c, object)
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
// @Param objectID path string true "World or object ID"
// @Success 204 ""
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/streamchat/{objectID}/leave [post]
func (s *StreamChat) apiChannelLeave(c *gin.Context) {
	object, user, err := s.getRequestContextObjects(c)
	if err != nil {
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request", err, s.log)
		return
	}

	channel, err := s.GetChannel(c, object)
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
func (s *StreamChat) getRequestContextObjects(c *gin.Context) (object universe.Object, user universe.User, err error) {
	objectID := c.Param("objectID")
	object, err = s.getObject(objectID)
	if err != nil {
		return
	}

	user, err = s.getUserFromContext(c, object)
	if err != nil {
		return
	}
	return

}

// Get object by UUID string.
func (s *StreamChat) getObject(id string) (universe.Object, error) {
	objectID, err := uuid.Parse(id)
	if err != nil {
		err := errors.WithMessagef(err, "Failed to parse ID %s", id)
		return nil, err
	}
	object, ok := s.node.GetObjectFromAllObjects(objectID)
	if !ok {
		err := errors.Errorf("Object not found: %s", objectID)
		return nil, err
	}
	return object, nil
}

// Resolve the user object from the request context and the given object.
func (s *StreamChat) getUserFromContext(ctx *gin.Context, object universe.Object) (universe.User, error) {

	userID, err := api.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "Streamchat: failed to get userID from context")
	}
	// The user needs to be 'in' the object (a.k.a world or object)
	user, ok := object.GetUser(userID, false)
	if !ok {
		return nil, fmt.Errorf("Streamchat: User %s not found in %s", userID, object.GetID())
	}
	return user, nil
}
