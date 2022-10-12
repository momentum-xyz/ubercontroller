package entry

type Node struct {
	*Space
}

type NodeAttributeID struct {
	AttributeID
}

type NodeAttribute struct {
	NodeAttributeID
	AttributePayload
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
