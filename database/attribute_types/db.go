package attribute_types

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
	getAttributeTypesQuery = `SELECT * FROM attribute_type;`

	removeAttributeTypeByNameQuery             = `DELETE FROM attribute_type WHERE attribute_name = $1;`
	removeAttributeTypesByNamesQuery           = `DELETE FROM attribute_type WHERE attribute_name = ANY($1);`
	removeAttributeTypesByPluginIDQuery        = `DELETE FROM attribute_type WHERE plugin_id = $1;`
	removeAttributeTypeByPluginIDAndNameQuery  = `DELETE FROM attribute_type WHERE plugin_id = $1 AND attribute_name = $2;`
	removeAttributeTypesByPluginIDAndNameQuery = `DELETE FROM attribute_type WHERE plugin_id = $1 AND attribute_name = ANY($2);`

	updateAttributeTypeNameQuery        = `UPDATE attribute_type SET attribute_name = $3 WHERE plugin_id = $1 AND attribute_name = $2;`
	updateAttributeTypeDescriptionQuery = `UPDATE attribute_type SET description = $3 WHERE plugin_id = $1 AND attribute_name = $2;`
	updateAttributeTypeOptionsQuery     = `UPDATE attribute_type SET options = $3 WHERE plugin_id = $1 AND attribute_name = $2;`

	upsertAttributeTypeQuery = `INSERT INTO attribute_type
									(plugin_id, attribute_name, description, options)
								VALUES
									($1, $2, $3, $4)
								ON CONFLICT (plugin_id, attribute_name)
								DO UPDATE SET
									description = $3,options = $4;`
)

var _ database.AttributeTypesDB = (*DB)(nil)

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

func (db *DB) GetAttributeTypes(ctx context.Context) ([]*entry.AttributeType, error) {
	var attributeTypes []*entry.AttributeType
	if err := pgxscan.Select(ctx, db.conn, &attributeTypes, getAttributeTypesQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return attributeTypes, nil
}

func (db *DB) UpsertAttributeType(ctx context.Context, attributeType *entry.AttributeType) error {
	if _, err := db.conn.Exec(
		ctx, upsertAttributeTypeQuery,
		attributeType.PluginID, attributeType.Name, attributeType.Description, attributeType.Options,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpsertAttributeTypes(ctx context.Context, attributeTypes []*entry.AttributeType) error {
	batch := &pgx.Batch{}
	for _, attributeType := range attributeTypes {
		batch.Queue(
			upsertAttributeTypeQuery,
			attributeType.PluginID, attributeType.Name, attributeType.Description, attributeType.Options,
		)
	}

	batchRes := db.conn.SendBatch(ctx, batch)
	defer batchRes.Close()

	var errs *multierror.Error
	for i := 0; i < batch.Len(); i++ {
		if _, err := batchRes.Exec(); err != nil {
			errs = multierror.Append(
				errs, errors.WithMessagef(err, "failed to exec db for: %v", attributeTypes[i].Name),
			)
		}
	}

	return errs.ErrorOrNil()
}

func (db *DB) RemoveAttributeTypeByName(ctx context.Context, name string) error {
	if _, err := db.conn.Exec(ctx, removeAttributeTypeByNameQuery, name); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) RemoveAttributeTypesByNames(ctx context.Context, names []string) error {
	if _, err := db.conn.Exec(ctx, removeAttributeTypesByNamesQuery, names); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) RemoveAttributeTypesByPluginID(ctx context.Context, pluginID uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeAttributeTypesByPluginIDQuery, pluginID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) RemoveAttributeTypeByID(ctx context.Context, attributeTypeID entry.AttributeTypeID) error {
	if _, err := db.conn.Exec(
		ctx, removeAttributeTypeByPluginIDAndNameQuery, attributeTypeID.PluginID, attributeTypeID.Name,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) RemoveAttributeTypesByIDs(ctx context.Context, attributeTypeIDs []entry.AttributeTypeID) error {
	batch := &pgx.Batch{}
	for _, attributeTypeID := range attributeTypeIDs {
		batch.Queue(
			removeAttributeTypesByPluginIDAndNameQuery,
			attributeTypeID.PluginID, attributeTypeID.Name,
		)
	}

	batchRes := db.conn.SendBatch(ctx, batch)
	defer batchRes.Close()

	var errs *multierror.Error
	for i := 0; i < batch.Len(); i++ {
		if _, err := batchRes.Exec(); err != nil {
			errs = multierror.Append(
				errs, errors.WithMessagef(err, "failed to exec db for: %v", attributeTypeIDs[i].Name),
			)
		}
	}

	return errs.ErrorOrNil()
}

func (db *DB) UpdateAttributeTypeName(
	ctx context.Context, attributeTypeID entry.AttributeTypeID, name string,
) error {
	if _, err := db.conn.Exec(
		ctx, updateAttributeTypeNameQuery, attributeTypeID.PluginID, attributeTypeID.Name, name,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateAttributeTypeOptions(
	ctx context.Context, attributeTypeID entry.AttributeTypeID, options *entry.AttributeOptions,
) error {
	if _, err := db.conn.Exec(
		ctx, updateAttributeTypeOptionsQuery, attributeTypeID.PluginID, attributeTypeID.Name, options,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateAttributeTypeDescription(
	ctx context.Context, attributeTypeID entry.AttributeTypeID, description *string,
) error {
	if _, err := db.conn.Exec(
		ctx, updateAttributeTypeDescriptionQuery, attributeTypeID.PluginID, attributeTypeID.Name, description,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}
