package node

//// @Summary Get timeline for object
//// @Schemes
//// @Description Returns a timeline for an object
//// @Tags timeline
//// @Accept json
//// @Produce json
//// @Success 200 {object} node.apiNewsFeedGetAll.Out
//// @Failure 404 {object} api.HTTPError
//// @Router /api/v4/{object_id}/timeline [get]
//func (n *Node) apiTimelineForObject(c *gin.Context) {
//	objectID, err := umid.Parse(c.Param("objectID"))
//	if err != nil {
//		err := errors.WithMessage(err, "Node: apiTimelineForObject: failed to parse object umid")
//		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
//		return
//	}
//
//	object, ok := n.GetObjectFromAllObjects(objectID)
//	if !ok {
//		err := errors.Errorf("Node: apiTimelineForObject: object not found: %s", objectID)
//		api.AbortRequest(c, http.StatusNotFound, "object_not_found", err, n.log)
//		return
//	}
//
//	object.GetObjectAttributes()
//
//	type Out struct {
//		Items []any `json:"items"`
//	}
//	out := Out{
//		Items: utils.GetFromAnyMap(*value, universe.ReservedAttributes.Object.NewsFeedItems.Key, []any(nil)),
//	}
//
//	c.JSON(http.StatusOK, out)
//}
