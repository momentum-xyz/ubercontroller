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
	database.PluginsDB
	database.SpaceTypesDB
	database.UserSpaceDB
	database.UserTypesDB
	database.AttributeTypesDB
	database.NodeAttributesDB
	database.SpaceAttributesDB
	database.SpaceUserAttributesDB
	database.UserAttributesDB
	database.UserUserAttributesDB
}

func (DB *DB) GetNodesDB() database.NodesDB {
	return DB.NodesDB
}

func (DB *DB) GetWorldsDB() database.WorldsDB {
	return DB.WorldsDB
}

func (DB *DB) GetSpacesDB() database.SpacesDB {
	return DB.SpacesDB
}

func (DB *DB) GetUsersDB() database.UsersDB {
	return DB.UsersDB
}

func (DB *DB) GetAssets2dDB() database.Assets2dDB {
	return DB.Assets2dDB
}

func (DB *DB) GetAssets3dDB() database.Assets3dDB {
	return DB.Assets3dDB
}

func (DB *DB) GetPluginsDB() database.PluginsDB {
	return DB.PluginsDB
}

func (DB *DB) GetSpaceTypesDB() database.SpaceTypesDB {
	return DB.SpaceTypesDB
}

func (DB *DB) GetUserTypesDB() database.UserTypesDB {
	return DB.UserTypesDB
}

func (DB *DB) GetAttributeTypesDB() database.AttributeTypesDB {
	return DB.AttributeTypesDB
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
	plugins database.PluginsDB,
	userSpace database.UserSpaceDB,
	spaceTypes database.SpaceTypesDB,
	userTypes database.UserTypesDB,
	attributeTypes database.AttributeTypesDB,
	nodeAttributes database.NodeAttributesDB,
	spaceAttributes database.SpaceAttributesDB,
	spaceUserAttributes database.SpaceUserAttributesDB,
	userAttributes database.UserAttributesDB,
	userUserAttributes database.UserUserAttributesDB,
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
		PluginsDB:             plugins,
		SpaceTypesDB:          spaceTypes,
		UserSpaceDB:           userSpace,
		UserTypesDB:           userTypes,
		AttributeTypesDB:      attributeTypes,
		NodeAttributesDB:      nodeAttributes,
		SpaceAttributesDB:     spaceAttributes,
		SpaceUserAttributesDB: spaceUserAttributes,
		UserAttributesDB:      userAttributes,
		UserUserAttributesDB:  userUserAttributes,
	}
}
