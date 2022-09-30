package space_attributes

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
	getSpaceAttributesQuery          = `SELECT * FROM space_attribute;`
	getSpaceAttributesQueryBySpaceId = `SELECT * FROM space_attribute where space_id=$1;`

	updateSpaceAttributeValueQuery                        = `UPDATE space_attribute SET value = $4 WHERE plugin_id=$1 and attribute_name = $2 and space_id = $3;`
	updateSpaceAttributeOptionsQuery                      = `UPDATE space_attribute SET options = $4 WHERE plugin_id=$1 and attribute_name = $2 and space_id = $3;`
	removeSpaceAttributeByNameQuery                       = `DELETE FROM space_attribute WHERE attribute_name = $1;`
	removeSpaceAttributesByNamesQuery                     = `DELETE FROM space_attribute WHERE attribute_name IN ($1);`
	removeSpaceAttributesByPluginIdQuery                  = `DELETE FROM space_attribute WHERE plugin_id = $1;`
	removeSpaceAttributeByPluginIdAndNameQuery            = `DELETE FROM space_attribute WHERE plugin_id = $1 and attribute_name =$2;`
	removeSpaceAttributesBySpaceIdQuery                   = `DELETE FROM space_attribute WHERE space_id = $1;`
	removeSpaceAttributeByNameAndSpaceIdQuery             = `DELETE FROM space_attribute WHERE attribute_name = $1 and space_id = $2;`
	removeSpaceAttributesByNamesAndSpaceIdQuery           = `DELETE FROM space_attribute WHERE attribute_name IN ($1)  and space_id = $2;`
	removeSpaceAttributesByPluginIdAndSpaceIdQuery        = `DELETE FROM space_attribute WHERE plugin_id = $1  and space_id = $2;`
	removeSpaceAttributesByPluginIdAndSpaceIdAndNameQuery = `DELETE FROM space_attribute WHERE plugin_id = $1  and space_id = $2 and name = $3;`

	upsertSpaceAttributeQuery = `INSERT INTO space_attribute
									(plugin_id, space_attribute_name,space_id, value, options)
								VALUES
									($1, $2, $3, $4, $5)
								ON CONFLICT (plugin_id,attribute_name, space_id)
								DO UPDATE SET
									value = $4,options = $5;`
)

var _ database.SpaceAttributesDB = (*DB)(nil)

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

func (db *DB) SpaceAttributesGetSpaceAttributes(ctx context.Context) ([]*entry.SpaceAttribute, error) {
	var assets []*entry.SpaceAttribute
	if err := pgxscan.Select(ctx, db.conn, &assets, getSpaceAttributesQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return assets, nil
}

func (db *DB) SpaceAttributesGetSpaceAttributesBySpaceId(
	ctx context.Context, spaceid uuid.UUID,
) ([]*entry.SpaceAttribute, error) {
	var assets []*entry.SpaceAttribute
	if err := pgxscan.Select(ctx, db.conn, &assets, getSpaceAttributesQueryBySpaceId, spaceid); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return assets, nil
}

func (db *DB) SpaceAttributesUpsertSpaceAttribute(ctx context.Context, spaceAttribute *entry.SpaceAttribute) error {
	if _, err := db.conn.Exec(
		ctx, upsertSpaceAttributeQuery, spaceAttribute.PluginID, spaceAttribute.Name, spaceAttribute.SpaceID,
		spaceAttribute.Value, spaceAttribute.Options,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceAttributesUpsertSpaceAttributes(
	ctx context.Context, spaceAttributes []*entry.SpaceAttribute,
) error {
	batch := &pgx.Batch{}
	for _, spaceAttribute := range spaceAttributes {
		batch.Queue(
			upsertSpaceAttributeQuery, spaceAttribute.PluginID, spaceAttribute.Name, spaceAttribute.SpaceID,
			spaceAttribute.Value, spaceAttribute.Options,
		)
	}

	batchRes := db.conn.SendBatch(ctx, batch)
	defer batchRes.Close()

	var errs *multierror.Error
	for i := 0; i < batch.Len(); i++ {
		if _, err := batchRes.Exec(); err != nil {
			errs = multierror.Append(
				errs, errors.WithMessagef(err, "failed to exec db for: %v", spaceAttributes[i].Name),
			)
		}
	}

	return errs.ErrorOrNil()
}

func (db *DB) SpaceAttributesRemoveSpaceAttributeByName(ctx context.Context, attributeName string) error {
	if _, err := db.conn.Exec(ctx, removeSpaceAttributeByNameQuery, attributeName); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceAttributesRemoveSpaceAttributesByNames(ctx context.Context, attributeNames []string) error {
	if _, err := db.conn.Exec(ctx, removeSpaceAttributesByNamesQuery, attributeNames); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceAttributesRemoveSpaceAttributesByPluginId(ctx context.Context, pluginID uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeSpaceAttributesByPluginIdQuery, pluginID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceAttributesRemoveSpaceAttributeByPluginIdAndName(
	ctx context.Context, pluginID uuid.UUID, attributeName string,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceAttributeByPluginIdAndNameQuery, pluginID, attributeName,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceAttributesRemoveSpaceAttributeBySpaceId(
	ctx context.Context, spaceID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceAttributesBySpaceIdQuery, spaceID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceAttributesRemoveSpaceAttributeByNameAndSpaceId(
	ctx context.Context, attributeName string, spaceID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceAttributeByNameAndSpaceIdQuery, attributeName, spaceID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceAttributesRemoveSpaceAttributeByNamesAndSpaceId(
	ctx context.Context, attributeNames []string, spaceID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceAttributesByNamesAndSpaceIdQuery, attributeNames, spaceID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceAttributesRemoveSpaceAttributeByPluginIdAndSpaceId(
	ctx context.Context, pluginId uuid.UUID, spaceID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceAttributesByPluginIdAndSpaceIdQuery, pluginId, spaceID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceAttributesRemoveSpaceAttributeByPluginIdAndSpaceIdAndName(
	ctx context.Context, pluginId uuid.UUID, attributeName string, spaceID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceAttributesByPluginIdAndSpaceIdAndNameQuery, pluginId, spaceID, attributeName,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceAttributesUpdateSpaceAttributeOptions(
	ctx context.Context, pluginID uuid.UUID, attributeName string, spaceId uuid.UUID, options *entry.AttributeOptions,
) error {
	if _, err := db.conn.Exec(
		ctx, updateSpaceAttributeOptionsQuery, attributeName, pluginID, spaceId, options,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceAttributesUpdateSpaceAttributeValue(
	ctx context.Context, pluginID uuid.UUID, attributeName string, spaceId uuid.UUID, value *entry.AttributeValue,
) error {
	if _, err := db.conn.Exec(
		ctx, updateSpaceAttributeValueQuery, attributeName, pluginID, spaceId, value,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}
