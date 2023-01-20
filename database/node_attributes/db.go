package node_attributes

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
	getNodeAttributesQuery                    = `SELECT * FROM node_attribute;`
	getNodeAttributeByAttributeIDQuery        = `SELECT * FROM node_attribute WHERE plugin_id = $1 AND attribute_name = $2;`
	getNodeAttributeValueByAttributeIDQuery   = `SELECT value FROM node_attribute WHERE plugin_id = $1 AND attribute_name = $2;`
	getNodeAttributeOptionsByAttributeIDQuery = `SELECT options FROM node_attribute WHERE plugin_id = $1 AND attribute_name = $2;`

	upsertNodeAttributeQuery = `INSERT INTO node_attribute
									(plugin_id, attribute_name, value, options)
								VALUES
									($1, $2, $3, $4)
								ON CONFLICT (plugin_id, attribute_name)
								DO UPDATE SET
									value = $3, options = $4;`

	updateNodeAttributeValueQuery   = `UPDATE node_attribute SET value = $3 WHERE plugin_id = $1 AND attribute_name = $2;`
	updateNodeAttributeOptionsQuery = `UPDATE node_attribute SET options = $3 WHERE plugin_id = $1 AND attribute_name = $2;`

	removeNodeAttributesByNameQuery       = `DELETE FROM node_attribute WHERE attribute_name = $1;`
	removeNodeAttributesByNamesQuery      = `DELETE FROM node_attribute WHERE attribute_name = ANY($1);`
	removeNodeAttributesByPluginIDQuery   = `DELETE FROM node_attribute WHERE plugin_id = $1;`
	removeNodeAttributeByAttributeIDQuery = `DELETE FROM node_attribute WHERE plugin_id = $1 AND attribute_name = $2;`
)

var _ database.NodeAttributesDB = (*DB)(nil)

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

func (db *DB) GetNodeAttributes(ctx context.Context) ([]*entry.NodeAttribute, error) {
	var assets []*entry.NodeAttribute
	if err := pgxscan.Select(ctx, db.conn, &assets, getNodeAttributesQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return assets, nil
}

func (db *DB) GetNodeAttributeByAttributeID(
	ctx context.Context, attributeID entry.AttributeID,
) (*entry.NodeAttribute, error) {
	var attr entry.NodeAttribute
	if err := pgxscan.Get(
		ctx, db.conn, &attr, getNodeAttributeByAttributeIDQuery,
		attributeID.PluginID, attributeID.Name,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &attr, nil
}

func (db *DB) GetNodeAttributeValueByAttributeID(
	ctx context.Context, attributeID entry.AttributeID,
) (*entry.AttributeValue, error) {
	var value entry.AttributeValue
	err := db.conn.QueryRow(ctx,
		getNodeAttributeValueByAttributeIDQuery,
		attributeID.PluginID,
		attributeID.Name).Scan(&value)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &value, nil
}

func (db *DB) GetNodeAttributeOptionsByAttributeID(
	ctx context.Context, attributeID entry.AttributeID,
) (*entry.AttributeOptions, error) {
	var options entry.AttributeOptions
	err := db.conn.QueryRow(ctx,
		getNodeAttributeOptionsByAttributeIDQuery,
		attributeID.PluginID,
		attributeID.Name).Scan(&options)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &options, nil
}

func (db *DB) UpsertNodeAttribute(ctx context.Context, nodeAttribute *entry.NodeAttribute) error {
	if _, err := db.conn.Exec(
		ctx, upsertNodeAttributeQuery,
		nodeAttribute.PluginID, nodeAttribute.Name, nodeAttribute.Value, nodeAttribute.Options,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpsertNodeAttributes(ctx context.Context, nodeAttributes []*entry.NodeAttribute) error {
	batch := &pgx.Batch{}
	for i := range nodeAttributes {
		batch.Queue(
			upsertNodeAttributeQuery,
			nodeAttributes[i].PluginID, nodeAttributes[i].Name, nodeAttributes[i].Value, nodeAttributes[i].Options,
		)
	}

	batchRes := db.conn.SendBatch(ctx, batch)
	defer batchRes.Close()

	var errs *multierror.Error
	for i := 0; i < batch.Len(); i++ {
		if _, err := batchRes.Exec(); err != nil {
			errs = multierror.Append(errs,
				errors.WithMessagef(err, "failed to exec db for: %+v", nodeAttributes[i].NodeAttributeID),
			)
		}
	}

	return errs.ErrorOrNil()
}

func (db *DB) RemoveNodeAttributesByName(ctx context.Context, name string) error {
	if _, err := db.conn.Exec(ctx, removeNodeAttributesByNameQuery, name); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) RemoveNodeAttributesByNames(ctx context.Context, names []string) error {
	if _, err := db.conn.Exec(ctx, removeNodeAttributesByNamesQuery, names); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) RemoveNodeAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeNodeAttributesByPluginIDQuery, pluginID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) RemoveNodeAttributeByAttributeID(ctx context.Context, attributeID entry.AttributeID) error {
	if _, err := db.conn.Exec(
		ctx, removeNodeAttributeByAttributeIDQuery, attributeID.PluginID, attributeID.Name,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateNodeAttributeValue(
	ctx context.Context, attributeID entry.AttributeID, value *entry.AttributeValue,
) error {
	if _, err := db.conn.Exec(
		ctx, updateNodeAttributeValueQuery,
		attributeID.PluginID, attributeID.Name, value,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateNodeAttributeOptions(
	ctx context.Context, attributeID entry.AttributeID, options *entry.AttributeOptions,
) error {
	if _, err := db.conn.Exec(
		ctx, updateNodeAttributeOptionsQuery,
		attributeID.PluginID, attributeID.Name, options,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}
