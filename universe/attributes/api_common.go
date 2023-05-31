package attributes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

// TODO: so attributes package or api package? or api/attributes? or vice versa? :)

// Struct to use for API query/params to select a plugin attribute.
type QueryPluginAttribute struct {
	PluginID      string `form:"plugin_id" json:"plugin_id" binding:"required"`
	AttributeName string `form:"attribute_name" json:"attribute_name" binding:"required"`
}

// Get attribute definition for API query.
func PluginAttributeFromQuery(c *gin.Context, n universe.Node) (universe.AttributeType, entry.AttributeID, error) {
	var attrID entry.AttributeID
	inQuery := QueryPluginAttribute{}
	if err := c.ShouldBindQuery(&inQuery); err != nil {
		return nil, attrID, fmt.Errorf("failed to bind query: %w", err)
	}
	pluginID, err := umid.Parse(inQuery.PluginID)
	if err != nil {
		return nil, attrID, fmt.Errorf("failed to parse plugin ID: %w", err)
	}
	attrType, ok := n.GetAttributeTypes().GetAttributeType(
		entry.AttributeTypeID{PluginID: pluginID, Name: inQuery.AttributeName})
	if !ok {
		return nil, attrID, fmt.Errorf("attribute type for %+v not found", inQuery)
	}
	attrID = entry.NewAttributeID(pluginID, inQuery.AttributeName)
	return attrType, attrID, nil
}
