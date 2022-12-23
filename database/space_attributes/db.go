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
	getSpaceAttributesQuery                           = `SELECT * FROM space_attribute;`
	getSpaceAttributeByIDQuery                        = `SELECT * FROM space_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND space_id = $3;`
	getSpaceAttributesQueryBySpaceIDQuery             = `SELECT * FROM space_attribute WHERE space_id = $1;`
	getSpaceAttributesByPluginIDAndAttributeNameQuery = `SELECT * FROM space_attribute WHERE plugin_id = $1 AND attribute_name = $2;`

	removeSpaceAttributeByNameQuery                       = `DELETE FROM space_attribute WHERE attribute_name = $1;`
	removeSpaceAttributesByNamesQuery                     = `DELETE FROM space_attribute WHERE attribute_name = ANY($1);`
	removeSpaceAttributesByPluginIDQuery                  = `DELETE FROM space_attribute WHERE plugin_id = $1;`
	removeSpaceAttributeByAttributeIDQuery                = `DELETE FROM space_attribute WHERE plugin_id = $1 AND attribute_name = $2;`
	removeSpaceAttributesBySpaceIDQuery                   = `DELETE FROM space_attribute WHERE space_id = $1;`
	removeSpaceAttributeByNameAndSpaceIDQuery             = `DELETE FROM space_attribute WHERE attribute_name = $1 AND space_id = $2;`
	removeSpaceAttributesByNamesAndSpaceIDQuery           = `DELETE FROM space_attribute WHERE attribute_name = ANY($1) AND space_id = $2;`
	removeSpaceAttributesByPluginIDAndSpaceIDQuery        = `DELETE FROM space_attribute WHERE plugin_id = $1 AND space_id = $2;`
	removeSpaceAttributesByPluginIDAndNameAndSpaceIDQuery = `DELETE FROM space_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND space_id = $3;`

	updateSpaceAttributeValueQuery   = `UPDATE space_attribute SET value = $4 WHERE plugin_id = $1 AND attribute_name = $2 AND space_id = $3;`
	updateSpaceAttributeOptionsQuery = `UPDATE space_attribute SET options = $4 WHERE plugin_id = $1 AND attribute_name = $2 AND space_id = $3;`

	upsertSpaceAttributeQuery = `INSERT INTO space_attribute
									(plugin_id, attribute_name, space_id, value, options)
								VALUES
									($1, $2, $3, $4, $5)
								ON CONFLICT (plugin_id, attribute_name, space_id)
								DO UPDATE SET
									value = $4, options = $5;`
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
	var spaceAttribute []*entry.SpaceAttribute
	if err := pgxscan.Select(ctx, db.conn, &spaceAttribute, getSpaceAttributesQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return spaceAttribute, nil
}

func (db *DB) SpaceAttributesGetSpaceAttributesByPluginIDAndName(ctx context.Context, pluginID uuid.UUID, attributeName string) ([]*entry.SpaceAttribute, error) {
	var spaceAttribute []*entry.SpaceAttribute
	err := pgxscan.Select(ctx, db.conn, &spaceAttribute, getSpaceAttributesByPluginIDAndAttributeNameQuery, pluginID, attributeName)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return spaceAttribute, nil
}

func (db *DB) SpaceAttributesGetSpaceAttributeByID(
	ctx context.Context, spaceAttributeID entry.SpaceAttributeID,
) (*entry.SpaceAttribute, error) {
	var spaceAttribute entry.SpaceAttribute
	if err := pgxscan.Get(
		ctx, db.conn, &spaceAttribute, getSpaceAttributeByIDQuery,
		spaceAttributeID.PluginID, spaceAttributeID.Name, spaceAttributeID.SpaceID,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &spaceAttribute, nil
}

func (db *DB) SpaceAttributesGetSpaceAttributesBySpaceID(
	ctx context.Context, spaceID uuid.UUID,
) ([]*entry.SpaceAttribute, error) {
	var spaceAttributes []*entry.SpaceAttribute
	if err := pgxscan.Select(ctx, db.conn, &spaceAttributes, getSpaceAttributesQueryBySpaceIDQuery, spaceID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return spaceAttributes, nil
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

func (db *DB) SpaceAttributesUpsertSpaceAttributes(ctx context.Context, spaceAttributes []*entry.SpaceAttribute) error {
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
				errs, errors.WithMessagef(err, "failed to exec db for: %s", spaceAttributes[i].Name),
			)
		}
	}

	return errs.ErrorOrNil()
}

func (db *DB) SpaceAttributesRemoveSpaceAttributeByName(ctx context.Context, name string) error {
	if _, err := db.conn.Exec(ctx, removeSpaceAttributeByNameQuery, name); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceAttributesRemoveSpaceAttributesByNames(ctx context.Context, names []string) error {
	if _, err := db.conn.Exec(ctx, removeSpaceAttributesByNamesQuery, names); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceAttributesRemoveSpaceAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeSpaceAttributesByPluginIDQuery, pluginID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceAttributesRemoveSpaceAttributeByAttributeID(
	ctx context.Context, attributeID entry.AttributeID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceAttributeByAttributeIDQuery, attributeID.PluginID, attributeID.Name,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceAttributesRemoveSpaceAttributeBySpaceID(
	ctx context.Context, spaceID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceAttributesBySpaceIDQuery, spaceID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceAttributesRemoveSpaceAttributeByNameAndSpaceID(
	ctx context.Context, name string, spaceID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceAttributeByNameAndSpaceIDQuery, name, spaceID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceAttributesRemoveSpaceAttributeByNamesAndSpaceID(
	ctx context.Context, names []string, spaceID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceAttributesByNamesAndSpaceIDQuery, names, spaceID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceAttributesRemoveSpaceAttributeByPluginIDAndSpaceID(
	ctx context.Context, pluginID uuid.UUID, spaceID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceAttributesByPluginIDAndSpaceIDQuery, pluginID, spaceID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceAttributesRemoveSpaceAttributeByID(
	ctx context.Context, spaceAttributeID entry.SpaceAttributeID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeSpaceAttributesByPluginIDAndNameAndSpaceIDQuery,
		spaceAttributeID.PluginID, spaceAttributeID.Name, spaceAttributeID.SpaceID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceAttributesUpdateSpaceAttributeValue(
	ctx context.Context, spaceAttributeID entry.SpaceAttributeID, value *entry.AttributeValue,
) error {
	if _, err := db.conn.Exec(
		ctx, updateSpaceAttributeValueQuery,
		spaceAttributeID.PluginID, spaceAttributeID.Name, spaceAttributeID.SpaceID, value,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceAttributesUpdateSpaceAttributeOptions(
	ctx context.Context, spaceAttributeID entry.SpaceAttributeID, options *entry.AttributeOptions,
) error {
	if _, err := db.conn.Exec(
		ctx, updateSpaceAttributeOptionsQuery,
		spaceAttributeID.PluginID, spaceAttributeID.Name, spaceAttributeID.SpaceID, options,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}
