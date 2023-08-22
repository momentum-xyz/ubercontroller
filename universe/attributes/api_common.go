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

// Struct to use for API attribute objects.
type AttributeValue struct {
	QueryPluginAttribute
	Value map[string]any `json:"value"`
}

// Type for API input/output of multiple attribute types.
// For now, we only support 1 type.
type AttributeMap struct { //map[AttributeType][]AttributeValue
	UserObject []AttributeValue `json:"object_user"`
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

// Get plugin attribute definition from API URL params.
func PluginAttributeFromURL(c *gin.Context, n universe.Node) (universe.AttributeType, entry.AttributeID, error) {
	var attrID entry.AttributeID
	pluginID, err := umid.Parse(c.Param("pluginID"))
	if err != nil {
		err := fmt.Errorf("invalid plugin ID: %w", err)
		return nil, attrID, err
	}
	attrName := c.Param("attrName")
	if attrName == "" {
		err := fmt.Errorf("invalid attribute name: \"%s\"", attrName)
		return nil, attrID, err
	}
	attrType, ok := n.GetAttributeTypes().GetAttributeType(
		entry.AttributeTypeID{PluginID: pluginID, Name: attrName})
	if !ok {
		err := fmt.Errorf("attribute type not found with %s %s", pluginID, attrName)
		return nil, attrID, err
	}
	attrID = entry.NewAttributeID(pluginID, attrName)
	return attrType, attrID, nil
}
