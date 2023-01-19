package object_attributes

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
	getObjectAttributesQuery              = `SELECT * FROM object_attribute;`
	getObjectAttributeByIDQuery           = `SELECT * FROM object_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND object_id = $3;`
	getObjectAttributesByObjectIDQuery    = `SELECT * FROM object_attribute WHERE object_id = $1;`
	getObjectAttributesByAttributeIDQuery = `SELECT * FROM object_attribute WHERE plugin_id = $1 AND attribute_name = $2;`

	upsertObjectAttributeQuery = `INSERT INTO object_attribute
    									(plugin_id, attribute_name, object_id, value, options)
									VALUES
									    ($1, $2, $3, $4, $5)
									ON CONFLICT (plugin_id, attribute_name, object_id)
									DO UPDATE SET
									    value = $4, options = $5;`

	updateObjectAttributeValueQuery   = `UPDATE object_attribute SET value = $4 WHERE plugin_id = $1 AND attribute_name = $2 AND object_id = $3;`
	updateObjectAttributeOptionsQuery = `UPDATE object_attribute SET options = $4 WHERE plugin_id = $1 AND attribute_name = $2 AND object_id = $3;`

	removeObjectAttributeByIDQuery                   = `DELETE FROM object_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND object_id = $3;`
	removeObjectAttributesByNameQuery                = `DELETE FROM object_attribute WHERE attribute_name = $1;`
	removeObjectAttributesByNamesQuery               = `DELETE FROM object_attribute WHERE attribute_name = ANY($1);`
	removeObjectAttributesByPluginIDQuery            = `DELETE FROM object_attribute WHERE plugin_id = $1;`
	removeObjectAttributesByAttributeIDQuery         = `DELETE FROM object_attribute WHERE plugin_id = $1 AND attribute_name = $2;`
	removeObjectAttributesByObjectIDQuery            = `DELETE FROM object_attribute WHERE object_id = $1;`
	removeObjectAttributesByNameAndObjectIDQuery     = `DELETE FROM object_attribute WHERE attribute_name = $1 AND object_id = $2;`
	removeObjectAttributesByNamesAndObjectIDQuery    = `DELETE FROM object_attribute WHERE attribute_name = ANY($1) AND object_id = $2;`
	removeObjectAttributesByPluginIDAndObjectIDQuery = `DELETE FROM object_attribute WHERE plugin_id = $1 AND object_id = $2;`
)

var _ database.ObjectAttributesDB = (*DB)(nil)

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

func (db *DB) GetObjectAttributes(ctx context.Context) ([]*entry.ObjectAttribute, error) {
	var attributes []*entry.ObjectAttribute
	if err := pgxscan.Select(ctx, db.conn, &attributes, getObjectAttributesQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return attributes, nil
}

func (db *DB) GetObjectAttributesByAttributeID(ctx context.Context, attributeID entry.AttributeID) ([]*entry.ObjectAttribute, error) {
	var attributes []*entry.ObjectAttribute
	if err := pgxscan.Select(
		ctx, db.conn, &attributes, getObjectAttributesByAttributeIDQuery,
		attributeID.PluginID, attributeID.Name,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return attributes, nil
}

func (db *DB) GetObjectAttributeByID(
	ctx context.Context, objectAttributeID entry.ObjectAttributeID,
) (*entry.ObjectAttribute, error) {
	var attribute entry.ObjectAttribute
	if err := pgxscan.Get(
		ctx, db.conn, &attribute, getObjectAttributeByIDQuery,
		objectAttributeID.PluginID, objectAttributeID.Name, objectAttributeID.ObjectID,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &attribute, nil
}

func (db *DB) GetObjectAttributesByObjectID(ctx context.Context, objectID uuid.UUID) ([]*entry.ObjectAttribute, error) {
	var attributes []*entry.ObjectAttribute
	if err := pgxscan.Select(ctx, db.conn, &attributes, getObjectAttributesByObjectIDQuery, objectID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return attributes, nil
}

func (db *DB) UpsertObjectAttribute(ctx context.Context, objectAttribute *entry.ObjectAttribute) error {
	if _, err := db.conn.Exec(
		ctx, upsertObjectAttributeQuery, objectAttribute.PluginID, objectAttribute.Name, objectAttribute.ObjectID,
		objectAttribute.Value, objectAttribute.Options,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpsertObjectAttributes(ctx context.Context, objectAttributes []*entry.ObjectAttribute) error {
	batch := &pgx.Batch{}
	for _, objectAttribute := range objectAttributes {
		batch.Queue(
			upsertObjectAttributeQuery, objectAttribute.PluginID, objectAttribute.Name, objectAttribute.ObjectID,
			objectAttribute.Value, objectAttribute.Options,
		)
	}

	batchRes := db.conn.SendBatch(ctx, batch)
	defer batchRes.Close()

	var errs *multierror.Error
	for i := 0; i < batch.Len(); i++ {
		if _, err := batchRes.Exec(); err != nil {
			errs = multierror.Append(
				errs, errors.WithMessagef(err, "failed to exec db for: %s", objectAttributes[i].Name),
			)
		}
	}

	return errs.ErrorOrNil()
}

func (db *DB) RemoveObjectAttributesByName(ctx context.Context, name string) error {
	if _, err := db.conn.Exec(ctx, removeObjectAttributesByNameQuery, name); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) RemoveObjectAttributesByNames(ctx context.Context, names []string) error {
	if _, err := db.conn.Exec(ctx, removeObjectAttributesByNamesQuery, names); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) RemoveObjectAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeObjectAttributesByPluginIDQuery, pluginID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) RemoveObjectAttributesByAttributeID(ctx context.Context, attributeID entry.AttributeID) error {
	if _, err := db.conn.Exec(
		ctx, removeObjectAttributesByAttributeIDQuery, attributeID.PluginID, attributeID.Name,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) RemoveObjectAttributesByObjectID(ctx context.Context, objectID uuid.UUID) error {
	if _, err := db.conn.Exec(
		ctx, removeObjectAttributesByObjectIDQuery, objectID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) RemoveObjectAttributesByNameAndObjectID(ctx context.Context, name string, objectID uuid.UUID) error {
	if _, err := db.conn.Exec(
		ctx, removeObjectAttributesByNameAndObjectIDQuery, name, objectID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) RemoveObjectAttributesByNamesAndObjectID(
	ctx context.Context, names []string, objectID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeObjectAttributesByNamesAndObjectIDQuery, names, objectID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) RemoveObjectAttributesByPluginIDAndObjectID(
	ctx context.Context, pluginID uuid.UUID, objectID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeObjectAttributesByPluginIDAndObjectIDQuery, pluginID, objectID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) RemoveObjectAttributeByID(
	ctx context.Context, objectAttributeID entry.ObjectAttributeID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeObjectAttributeByIDQuery,
		objectAttributeID.PluginID, objectAttributeID.Name, objectAttributeID.ObjectID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateObjectAttributeValue(
	ctx context.Context, objectAttributeID entry.ObjectAttributeID, value *entry.AttributeValue,
) error {
	if _, err := db.conn.Exec(
		ctx, updateObjectAttributeValueQuery,
		objectAttributeID.PluginID, objectAttributeID.Name, objectAttributeID.ObjectID, value,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateObjectAttributeOptions(
	ctx context.Context, objectAttributeID entry.ObjectAttributeID, options *entry.AttributeOptions,
) error {
	if _, err := db.conn.Exec(
		ctx, updateObjectAttributeOptionsQuery,
		objectAttributeID.PluginID, objectAttributeID.Name, objectAttributeID.ObjectID, options,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}
