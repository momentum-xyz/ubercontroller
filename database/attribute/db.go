package attribute

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
	getAttributesQuery = `SELECT * FROM attribute;`

	updateAttributeNameQuery                = `UPDATE attribute SET attribute_name = $3 WHERE attribute_name = $1 and plugin_id=$2;`
	updateAttributeOptionsQuery             = `UPDATE attribute SET options = $3 WHERE attribute_name = $1 and plugin_id=$2;`
	removeAttributeByNameQuery              = `DELETE FROM attribute WHERE attribute_name = $1;`
	removeAttributesByNamesQuery            = `DELETE FROM attribute WHERE attribute_name IN ($1);`
	removeAttributesByPluginIdQuery         = `DELETE FROM attribute WHERE plugin_id = $1;`
	removeAttributeByNameAndPluginIdQuery   = `DELETE FROM attribute WHERE attribute_name = $1 and  plugin_id=$2;`
	removeAttributesByNamesAndPluginIdQuery = `DELETE FROM attribute WHERE attribute_name IN ($1) and  plugin_id=$2;`
	updateAttributeDescriptionQuery         = `UPDATE attribute SET description = $3 WHERE attribute_name = $1 and plugin_id=$2;`

	upsertAttributeQuery = `INSERT INTO attribute
									(plugin_id, attribute_name,description, options)
								VALUES
									($1, $2, $3, $4)
								ON CONFLICT (plugin_id,attribute_name)
								DO UPDATE SET
									description = $3,options = $4;`
)

var _ database.AttributesDB = (*DB)(nil)

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

func (db *DB) AttributesGetAttributes(ctx context.Context) ([]*entry.Attribute, error) {
	var assets []*entry.Attribute
	if err := pgxscan.Select(ctx, db.conn, &assets, getAttributesQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return assets, nil
}

func (db *DB) AttributesUpsertAttribute(ctx context.Context, attribute *entry.Attribute) error {
	if _, err := db.conn.Exec(
		ctx, upsertAttributeQuery, attribute.AttributeID.PluginID, attribute.AttributeID.Name, attribute.Description,
		attribute.Options,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) AttributesUpsertAttributes(ctx context.Context, attributes []*entry.Attribute) error {
	batch := &pgx.Batch{}
	for _, attribute := range attributes {
		batch.Queue(
			upsertAttributeQuery, attribute.AttributeID.PluginID, attribute.AttributeID.Name, attribute.Description,
			attribute.Options,
		)
	}

	batchRes := db.conn.SendBatch(ctx, batch)
	defer batchRes.Close()

	var errs *multierror.Error
	for i := 0; i < batch.Len(); i++ {
		if _, err := batchRes.Exec(); err != nil {
			errs = multierror.Append(
				errs, errors.WithMessagef(err, "failed to exec db for: %v", attributes[i].AttributeID.Name),
			)
		}
	}

	return errs.ErrorOrNil()
}

func (db *DB) AttributesRemoveAttributeByName(ctx context.Context, attributeName string) error {
	if _, err := db.conn.Exec(ctx, removeAttributeByNameQuery, attributeName); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) AttributesRemoveAttributesByNames(ctx context.Context, attributeNames []string) error {
	if _, err := db.conn.Exec(ctx, removeAttributesByNamesQuery, attributeNames); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) AttributesRemoveAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeAttributesByPluginIdQuery, pluginID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) AttributesRemoveAttributeByID(ctx context.Context, attributeId entry.AttributeID) error {
	if _, err := db.conn.Exec(
		ctx, removeAttributeByNameAndPluginIdQuery, attributeId.Name, attributeId.PluginID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) AttributesRemoveAttributesByIDs(
	ctx context.Context, attributeIDs []entry.AttributeID,
) error {
	// TODO: implement this!
	//if _, err := db.conn.Exec(ctx, removeAttributesByNamesAndPluginIdQuery, attributeNames, pluginID); err != nil {
	//	return errors.WithMessage(err, "failed to exec db")
	//}
	return nil
}

func (db *DB) AttributesUpdateAttributeName(
	ctx context.Context, attributeId entry.AttributeID, newName string,
) error {
	if _, err := db.conn.Exec(
		ctx, updateAttributeNameQuery, attributeId.Name, attributeId.PluginID, newName,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) AttributesUpdateAttributeOptions(
	ctx context.Context, attributeId entry.AttributeID, attributeOptions *entry.AttributeOptions,
) error {
	if _, err := db.conn.Exec(
		ctx, updateAttributeOptionsQuery, attributeId.Name, attributeId.PluginID, attributeOptions,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) AttributesUpdateAttributeDescription(
	ctx context.Context, attributeId entry.AttributeID, description *string,
) error {
	if _, err := db.conn.Exec(
		ctx, updateAttributeDescriptionQuery, attributeId.Name, attributeId.PluginID, description,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}
