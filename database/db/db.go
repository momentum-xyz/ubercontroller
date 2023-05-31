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
	database.ActivitiesDB
	database.ObjectsDB
	database.UsersDB
	database.Assets2dDB
	database.Assets3dDB
	database.PluginsDB
	database.ObjectTypesDB
	database.UserObjectsDB
	database.UserTypesDB
	database.AttributeTypesDB
	database.NodeAttributesDB
	database.ObjectAttributesDB
	database.ObjectUserAttributesDB
	database.UserAttributesDB
	database.UserUserAttributesDB
	database.StakesDB
	database.NFTsDB
}

func NewDB(
	conn *pgxpool.Pool,
	common database.CommonDB,
	nodes database.NodesDB,
	worlds database.WorldsDB,
	objects database.ObjectsDB,
	activities database.ActivitiesDB,
	users database.UsersDB,
	assets2d database.Assets2dDB,
	assets3d database.Assets3dDB,
	plugins database.PluginsDB,
	userObjects database.UserObjectsDB,
	objectTypes database.ObjectTypesDB,
	userTypes database.UserTypesDB,
	attributeTypes database.AttributeTypesDB,
	nodeAttributes database.NodeAttributesDB,
	objectAttributes database.ObjectAttributesDB,
	objectUserAttributes database.ObjectUserAttributesDB,
	userAttributes database.UserAttributesDB,
	userUserAttributes database.UserUserAttributesDB,
	stakesDB database.StakesDB,
	nftsDB database.NFTsDB,
) *DB {
	return &DB{
		conn:                   conn,
		CommonDB:               common,
		NodesDB:                nodes,
		WorldsDB:               worlds,
		ActivitiesDB:           activities,
		ObjectsDB:              objects,
		UsersDB:                users,
		Assets2dDB:             assets2d,
		Assets3dDB:             assets3d,
		PluginsDB:              plugins,
		ObjectTypesDB:          objectTypes,
		UserObjectsDB:          userObjects,
		UserTypesDB:            userTypes,
		AttributeTypesDB:       attributeTypes,
		NodeAttributesDB:       nodeAttributes,
		ObjectAttributesDB:     objectAttributes,
		ObjectUserAttributesDB: objectUserAttributes,
		UserAttributesDB:       userAttributes,
		UserUserAttributesDB:   userUserAttributes,
		StakesDB:               stakesDB,
		NFTsDB:                 nftsDB,
	}
}

func (DB *DB) GetCommonDB() database.CommonDB {
	return DB.CommonDB
}

func (DB *DB) GetNodesDB() database.NodesDB {
	return DB.NodesDB
}

func (DB *DB) GetWorldsDB() database.WorldsDB {
	return DB.WorldsDB
}

func (DB *DB) GetActivitiesDB() database.ActivitiesDB {
	return DB.ActivitiesDB
}

func (DB *DB) GetObjectsDB() database.ObjectsDB {
	return DB.ObjectsDB
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

func (DB *DB) GetObjectTypesDB() database.ObjectTypesDB {
	return DB.ObjectTypesDB
}

func (DB *DB) GetUserTypesDB() database.UserTypesDB {
	return DB.UserTypesDB
}

func (DB *DB) GetAttributeTypesDB() database.AttributeTypesDB {
	return DB.AttributeTypesDB
}

func (DB *DB) GetNodeAttributesDB() database.NodeAttributesDB {
	return DB.NodeAttributesDB
}

func (DB *DB) GetObjectAttributesDB() database.ObjectAttributesDB {
	return DB.ObjectAttributesDB
}

func (DB *DB) GetObjectUserAttributesDB() database.ObjectUserAttributesDB {
	return DB.ObjectUserAttributesDB
}

func (DB *DB) GetUserAttributesDB() database.UserAttributesDB {
	return DB.UserAttributesDB
}

func (DB *DB) GetUserUserAttributesDB() database.UserUserAttributesDB {
	return DB.UserUserAttributesDB
}

func (DB *DB) GetUserObjectsDB() database.UserObjectsDB {
	return DB.UserObjectsDB
}

func (DB *DB) GetStakesDB() database.StakesDB {
	return DB.StakesDB
}

func (DB *DB) GetNFTsDB() database.NFTsDB {
	return DB.NFTsDB
}
