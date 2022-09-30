package entry

import "github.com/google/uuid"

type Node struct {
	*Space
}

type NodeAttribute struct {
	PluginID uuid.UUID       `db:"plugin_id"`
	Name     string          `db:"attribute_name"`
	Value    *AttributeValue `db:"value"`
}

//
//type AttributeValue struct {
//	//HashSalt           string `db:"hash_salt" json:"hash_salt"`
//	//MainDomain         string `db:"main_domain" json:"main_domain"`
//	//WorldCreatorsGroup string `db:"world_creators_group" json:"world_creators_group"`
//	//DefaultWorld       string `db:"default_world" json:"default_world"`
//}
