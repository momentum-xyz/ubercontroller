package db

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/momentum-xyz/ubercontroller/database"
)

var _ database.DB = (*DB)(nil)

type DB struct {
	conn *pgxpool.Pool
	database.CommonDB
	database.NodesDB
	database.WorldsDB
	database.SpacesDB
	database.UsersDB
	database.Assets2dDB
	database.Assets3dDB
	database.SpaceTypesDB
	database.UserTypesDB
	database.PluginsDB
	database.AttributeTypesDB
	database.NodeAttributesDB
	database.SpaceAttributesDB
	database.SpaceUserAttributesDB
	database.UserAttributesDB
	database.UserUserAttributesDB
}

func NewDB(
	conn *pgxpool.Pool,
	common database.CommonDB,
	nodes database.NodesDB,
	worlds database.WorldsDB,
	spaces database.SpacesDB,
	users database.UsersDB,
	assets2d database.Assets2dDB,
	assets3d database.Assets3dDB,
	spaceTypes database.SpaceTypesDB,
	userTypes database.UserTypesDB,
	attributeTypes database.AttributeTypesDB,
	nodeAttributes database.NodeAttributesDB,
	spaceAttributes database.SpaceAttributesDB,
	spaceUserAttributes database.SpaceUserAttributesDB,
	userAttributes database.UserAttributesDB,
	userUserAttributes database.UserUserAttributesDB,
	plugins database.PluginsDB,
) *DB {
	return &DB{
		conn:                  conn,
		CommonDB:              common,
		NodesDB:               nodes,
		WorldsDB:              worlds,
		SpacesDB:              spaces,
		UsersDB:               users,
		Assets2dDB:            assets2d,
		Assets3dDB:            assets3d,
		SpaceTypesDB:          spaceTypes,
		UserTypesDB:           userTypes,
		AttributeTypesDB:      attributeTypes,
		NodeAttributesDB:      nodeAttributes,
		SpaceAttributesDB:     spaceAttributes,
		SpaceUserAttributesDB: spaceUserAttributes,
		UserAttributesDB:      userAttributes,
		UserUserAttributesDB:  userUserAttributes,
		PluginsDB:             plugins,
	}
}
