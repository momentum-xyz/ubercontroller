package entry

import "time"

type Node struct {
	*Object
}

type NodeAttributeID struct {
	AttributeID
}

type NodeAttribute struct {
	NodeAttributeID
	*AttributePayload
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func NewNodeAttribute(nodeAttributeID NodeAttributeID, payload *AttributePayload) *NodeAttribute {
	return &NodeAttribute{
		NodeAttributeID:  nodeAttributeID,
		AttributePayload: payload,
	}
}

func NewNodeAttributeID(attributeID AttributeID) NodeAttributeID {
	return NodeAttributeID{
		AttributeID: attributeID,
	}
}

//
//type AttributeValue struct {
//	//HashSalt           string `db:"hash_salt" json:"hash_salt"`
//	//MainDomain         string `db:"main_domain" json:"main_domain"`
//	//WorldCreatorsGroup string `db:"world_creators_group" json:"world_creators_group"`
//	//DefaultWorld       string `db:"default_world" json:"default_world"`
//}
