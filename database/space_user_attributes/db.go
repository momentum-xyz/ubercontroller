package space_user_attributes

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/entry"
)

const (
	getSpaceUserAttributesQuery                     = `SELECT * FROM space_user_attribute;`
	getSpaceUserAttributesQueryBySpaceId            = `SELECT * FROM space_user_attribute where space_id=$1;`
	getSpaceUserAttributesQueryByUserId             = `SELECT * FROM space_user_attribute where user_id=$1;`
	getSpaceUserAttributesQueryBySpaceIdAndUserId   = `SELECT * FROM space_user_attribute where space_id=$1 and user_id=$2;`
	updateSpaceUserAttributeValueQuery              = `UPDATE space_user_attribute SET value = $5 WHERE plugin_id=$1 and attribute_name = $2 and space_id = $3 and user_id = $4;`
	updateSpaceUserAttributeOptionsQuery            = `UPDATE space_user_attribute SET options = $5 WHERE plugin_id=$1 and attribute_name = $2 and space_id = $3 and user_id = $4;`
	removeSpaceUserAttributeByNameQuery             = `DELETE FROM space_user_attribute WHERE attribute_name = $1;`
	removeSpaceUserAttributesByNamesQuery           = `DELETE FROM space_user_attribute WHERE attribute_name IN ($1);`
	removeSpaceUserAttributesByPluginIdQuery        = `DELETE FROM space_user_attribute WHERE plugin_id = $1;`
	removeSpaceUserAttributeByPluginIdAndNameQuery  = `DELETE FROM space_user_attribute WHERE plugin_id = $1 and attribute_name =$2;`
	removeSpaceUserAttributesBySpaceIdQuery         = `DELETE FROM space_user_attribute WHERE space_id = $1;`
	removeSpaceUserAttributeByNameAndSpaceIdQuery   = `DELETE FROM space_user_attribute WHERE attribute_name = $1 and space_id = $2;`
	removeSpaceUserAttributesByNamesAndSpaceIdQuery = `DELETE FROM space_user_attribute WHERE attribute_name IN ($1)  and space_id = $2;`

	removeSpaceUserAttributesByUserIdQuery         = `DELETE FROM space_user_attribute WHERE user_id = $1;`
	removeSpaceUserAttributeByNameAndUserIdQuery   = `DELETE FROM space_user_attribute WHERE attribute_name = $1 and user_id = $2;`
	removeSpaceUserAttributesByNamesAndUserIdQuery = `DELETE FROM space_user_attribute WHERE attribute_name IN ($1)  and user_id = $2;`

	removeSpaceUserAttributesBySpaceIdAndUserIdQuery         = `DELETE FROM space_user_attribute WHERE space_id = $1 and user_id = $2;`
	removeSpaceUserAttributeByNameAndSpaceIdAndUserIdQuery   = `DELETE FROM space_user_attribute WHERE attribute_name = $1 and space_id = $2 and user_id = $3;`
	removeSpaceUserAttributesByNamesAndSpaceIdAndUserIdQuery = `DELETE FROM space_user_attribute WHERE attribute_name IN ($1)  and space_id = $2 and user_id = $3;`

	removeSpaceUserAttributesByPluginIdAndSpaceIdQuery        = `DELETE FROM space_user_attribute WHERE plugin_id = $1  and space_id = $2;`
	removeSpaceUserAttributesByPluginIdAndSpaceIdAndNameQuery = `DELETE FROM space_user_attribute WHERE plugin_id = $1  and space_id = $2 and name = $3;`

	removeSpaceUserAttributesByPluginIdAndUserIdQuery        = `DELETE FROM space_user_attribute WHERE plugin_id = $1  and user_id = $2;`
	removeSpaceUserAttributesByPluginIdAndUserIdAndNameQuery = `DELETE FROM space_user_attribute WHERE plugin_id = $1  and user_id = $2 and name = $3;`

	removeSpaceUserAttributesByPluginIdAndSpaceIdAndUserIdQuery        = `DELETE FROM space_user_attribute WHERE plugin_id = $1  and space_id = $2  and user_id = $3;`
	removeSpaceUserAttributesByPluginIdAndSpaceIdAndUserIdAndNameQuery = `DELETE FROM space_user_attribute WHERE plugin_id = $1  and space_id = $2  and user_id = $3 and name = $4;`

	upsertSpaceUserAttributeQuery = `INSERT INTO space_user_attribute
									(plugin_id, space_user_attribute_name,space_id, user_id, value, options)
								VALUES
									($1, $2, $3, $4, $5, $6)
								ON CONFLICT (plugin_id,attribute_name, space_id, user_id)
								DO UPDATE SET
									value = $5,options = $6;`
)

var _ database.SpaceUserAttributesDB = (*DB)(nil)

type DB struct {
	conn   *pgxpool.Pool
	common database.CommonDB
}

func NewDB(conn *pgxpool.Pool, commonDB database.CommonDB) *DB {
	return &DB{
		conn:   conn,
		common: commonDB,
	}
}

func (db *DB) SpaceUserAttributesGetSpaceUserAttributes(ctx context.Context) ([]*entry.SpaceUserAttribute, error) {
	var assets []*entry.SpaceUserAttribute
	if err := pgxscan.Select(ctx, db.conn, &assets, getSpaceUserAttributesQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return assets, nil
}

func (db *DB) SpaceUserAttributesGetSpaceUserAttributesBySpaceId(
	ctx context.Context, spaceId uuid.UUID,
) ([]*entry.SpaceUserAttribute, error) {
	var assets []*entry.SpaceUserAttribute
	if err := pgxscan.Select(ctx, db.conn, &assets, getSpaceUserAttributesQueryBySpaceId, spaceId); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return assets, nil
}

func (db *DB) SpaceUserAttributesGetSpaceUserAttributesByUserId(
	ctx context.Context, userId uuid.UUID,
) ([]*entry.SpaceUserAttribute, error) {
	var assets []*entry.SpaceUserAttribute
	if err := pgxscan.Select(ctx, db.conn, &assets, getSpaceUserAttributesQueryByUserId, userId); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return assets, nil
}

func (db *DB) SpaceUserAttributesGetSpaceUserAttributesBySpaceIdAndUserId(
	ctx context.Context, spaceId uuid.UUID, userId uuid.UUID,
) ([]*entry.SpaceUserAttribute, error) {
	var assets []*entry.SpaceUserAttribute
	if err := pgxscan.Select(
		ctx, db.conn, &assets, getSpaceUserAttributesQueryBySpaceIdAndUserId, spaceId, userId,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return assets, nil
}

func (db *DB) SpaceUserAttributesUpsertSpaceUserAttribute(
	ctx context.Context, spaceUserAttribute *entry.SpaceUserAttribute,
) error {
	if _, err := db.conn.Exec(
		ctx, upsertSpaceUserAttributeQuery, spaceUserAttribute.PluginID, spaceUserAttribute.Name,
		spaceUserAttribute.SpaceID,
		spaceUserAttribute.UserID, spaceUserAttribute.Value, spaceUserAttribute.Options,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceUserAttributesUpsertSpaceUserAttributes(
	ctx context.Context, spaceUserAttributes []*entry.SpaceUserAttribute,
) error {
	batch := &pgx.Batch{}
	for _, spaceUserAttribute := range spaceUserAttributes {
		batch.Queue(
			upsertSpaceUserAttributeQuery, spaceUserAttribute.PluginID, spaceUserAttribute.Name,
			spaceUserAttribute.SpaceID,
			spaceUserAttribute.UserID, spaceUserAttribute.Value, spaceUserAttribute.Options,
		)
	}

	batchRes := db.conn.SendBatch(ctx, batch)
	defer batchRes.Close()

	var errs *multierror.Error
	for i := 0; i < batch.Len(); i++ {
		if _, err := batchRes.Exec(); err != nil {
			errs = multierror.Append(
				errs, errors.WithMessagef(err, "failed to exec db for: %v", spaceUserAttributes[i].Name),
			)
		}
	}

	return errs.ErrorOrNil()
}

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByName(ctx context.Context, attributeName string) error {
	if _, err := db.conn.Exec(ctx, removeSpaceUserAttributeByNameQuery, attributeName); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributesByNames(ctx context.Context, attributeNames []string) error {
	if _, err := db.conn.Exec(ctx, removeSpaceUserAttributesByNamesQuery, attributeNames); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributesByPluginId(ctx context.Context, pluginID uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeSpaceUserAttributesByPluginIdQuery, pluginID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByPluginIdAndName(
	ctx context.Context, pluginID uuid.UUID, attributeName string,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceUserAttributeByPluginIdAndNameQuery, pluginID, attributeName,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeBySpaceId(
	ctx context.Context, spaceID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceUserAttributesBySpaceIdQuery, spaceID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByNameAndSpaceId(
	ctx context.Context, attributeName string, spaceID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceUserAttributeByNameAndSpaceIdQuery, attributeName, spaceID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByNamesAndSpaceId(
	ctx context.Context, attributeNames []string, spaceID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceUserAttributesByNamesAndSpaceIdQuery, attributeNames, spaceID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByUserId(
	ctx context.Context, userId uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceUserAttributesByUserIdQuery, userId,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByNameAndUserId(
	ctx context.Context, attributeName string, userId uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceUserAttributeByNameAndUserIdQuery, attributeName, userId,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByNamesAndUserId(
	ctx context.Context, attributeNames []string, userId uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceUserAttributesByNamesAndUserIdQuery, attributeNames, userId,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeBySpaceIdAndUserId(
	ctx context.Context, spaceID uuid.UUID, userId uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceUserAttributesBySpaceIdAndUserIdQuery, spaceID, userId,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByNameAndSpaceIdAndUserId(
	ctx context.Context, attributeName string, spaceID uuid.UUID, userId uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceUserAttributeByNameAndSpaceIdAndUserIdQuery, attributeName, spaceID, userId,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByNamesAndSpaceIdAndUserId(
	ctx context.Context, attributeNames []string, spaceID uuid.UUID, userId uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceUserAttributesByNamesAndSpaceIdAndUserIdQuery, attributeNames, spaceID, userId,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByPluginIdAndSpaceId(
	ctx context.Context, pluginId uuid.UUID, spaceID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceUserAttributesByPluginIdAndSpaceIdQuery, pluginId, spaceID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByPluginIdAndSpaceIdAndName(
	ctx context.Context, pluginId uuid.UUID, attributeName string, spaceID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceUserAttributesByPluginIdAndSpaceIdAndNameQuery, pluginId, spaceID, attributeName,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByPluginIdAndUserId(
	ctx context.Context, pluginId uuid.UUID, userId uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceUserAttributesByPluginIdAndUserIdQuery, pluginId, userId,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByPluginIdAndUserIdAndName(
	ctx context.Context, pluginId uuid.UUID, attributeName string, userId uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceUserAttributesByPluginIdAndUserIdAndNameQuery, pluginId, userId, attributeName,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}
func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByPluginIdAndSpaceIdAndUserId(
	ctx context.Context, pluginId uuid.UUID, spaceID uuid.UUID, userId uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceUserAttributesByPluginIdAndSpaceIdAndUserIdQuery, pluginId, spaceID, userId,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByPluginIdAndSpaceIdAndUserIdAndName(
	ctx context.Context, pluginId uuid.UUID, attributeName string, spaceID uuid.UUID, userId uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceUserAttributesByPluginIdAndSpaceIdAndUserIdAndNameQuery, pluginId, spaceID, userId,
		attributeName,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceUserAttributesUpdateSpaceUserAttributeOptions(
	ctx context.Context, pluginID uuid.UUID, attributeName string, spaceId uuid.UUID, userId uuid.UUID,
	options *entry.AttributeOptions,
) error {
	if _, err := db.conn.Exec(
		ctx, updateSpaceUserAttributeOptionsQuery, attributeName, pluginID, spaceId, userId, options,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceUserAttributesUpdateSpaceUserAttributeValue(
	ctx context.Context, pluginID uuid.UUID, attributeName string, spaceId uuid.UUID, userId uuid.UUID,
	value *entry.AttributeValue,
) error {
	if _, err := db.conn.Exec(
		ctx, updateSpaceUserAttributeValueQuery, attributeName, pluginID, spaceId, userId, value,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}
