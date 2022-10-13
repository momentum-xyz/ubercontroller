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

// TODO: change "node_attributes" table to "node_attribute"

const (
	getNodeAttributesQuery                   = `SELECT * FROM node_attributes;`
	getNodeAttributeByPluginIDAndName        = `SELECT * FROM node_attributes WHERE plugin_id = $1 AND attribute_name = $2;`
	getNodeAttributeValueByPluginIDAndName   = `SELECT value FROM node_attributes WHERE plugin_id = $1 AND attribute_name = $2;`
	getNodeAttributeOptionsByPluginIDAndName = `SELECT options FROM node_attributes WHERE plugin_id = $1 AND attribute_name = $2;`

	updateNodeAttributeValueQuery             = `UPDATE node_attributes SET value = $3 WHERE plugin_id = $1 AND attribute_name = $2;`
	updateNodeAttributeOptionsQuery           = `UPDATE node_attributes SET options = $3 WHERE plugin_id = $1 AND attribute_name = $2;`
	removeNodeAttributeByNameQuery            = `DELETE FROM node_attributes WHERE attribute_name = $1;`
	removeNodeAttributesByNamesQuery          = `DELETE FROM node_attributes WHERE attribute_name IN ($1);`
	removeNodeAttributesByPluginIdQuery       = `DELETE FROM node_attributes WHERE plugin_id = $1;`
	removeNodeAttributeByPluginIdAndNameQuery = `DELETE FROM node_attributes WHERE plugin_id = $1 AND attribute_name = $2;`

	upsertNodeAttributeQuery = `INSERT INTO node_attributes
									(plugin_id, attribute_name, value, options)
								VALUES
									($1, $2, $3, $4)
								ON CONFLICT (plugin_id, attribute_name)
								DO UPDATE SET
									value = $3, options = $4;`
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

func (db *DB) NodeAttributesGetNodeAttributes(ctx context.Context) ([]*entry.NodeAttribute, error) {
	var assets []*entry.NodeAttribute
	if err := pgxscan.Select(ctx, db.conn, &assets, getNodeAttributesQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return assets, nil
}

func (db *DB) NodeAttributesGetNodeAttributeByPluginIDAndName(
	ctx context.Context, pluginID uuid.UUID, attributeName string,
) (*entry.NodeAttribute, error) {
	var attr entry.NodeAttribute
	if err := pgxscan.Get(
		ctx, db.conn, &attr, getNodeAttributeByPluginIDAndName,
		pluginID, attributeName,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &attr, nil
}

func (db *DB) NodeAttributesGetNodeAttributeValueByPluginIDAndName(
	ctx context.Context, pluginID uuid.UUID, attributeName string,
) (*entry.AttributeValue, error) {
	var value entry.AttributeValue
	if err := pgxscan.Get(
		ctx, db.conn, &value, getNodeAttributeValueByPluginIDAndName,
		pluginID, attributeName,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &value, nil
}

func (db *DB) NodeAttributesGetNodeAttributeOptionsByPluginIDAndName(
	ctx context.Context, pluginID uuid.UUID, attributeName string,
) (*entry.AttributeOptions, error) {
	var options entry.AttributeOptions
	if err := pgxscan.Get(
		ctx, db.conn, &options, getNodeAttributeOptionsByPluginIDAndName,
		pluginID, attributeName,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &options, nil
}

func (db *DB) NodeAttributesUpsertNodeAttribute(ctx context.Context, nodeAttribute *entry.NodeAttribute) error {
	if _, err := db.conn.Exec(
		ctx, upsertNodeAttributeQuery,
		nodeAttribute.PluginID, nodeAttribute.Name, nodeAttribute.Value, nodeAttribute.Options,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) NodeAttributesUpsertNodeAttributes(ctx context.Context, nodeAttributes []*entry.NodeAttribute) error {
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
				errors.WithMessagef(err, "failed to exec db for: %s", nodeAttributes[i].Name),
			)
		}
	}

	return errs.ErrorOrNil()
}

func (db *DB) NodeAttributesUpdateNodeAttributeValue(
	ctx context.Context, pluginID uuid.UUID, attributeName string, value *entry.AttributeValue,
) error {
	if _, err := db.conn.Exec(
		ctx, updateNodeAttributeValueQuery,
		pluginID, attributeName, value,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) NodeAttributesUpdateNodeAttributeOptions(
	ctx context.Context, pluginID uuid.UUID, attributeName string, options *entry.AttributeOptions,
) error {
	if _, err := db.conn.Exec(
		ctx, updateNodeAttributeOptionsQuery,
		pluginID, attributeName, options,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) NodeAttributesRemoveNodeAttributeByName(ctx context.Context, attributeName string) error {
	if _, err := db.conn.Exec(ctx, removeNodeAttributeByNameQuery, attributeName); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) NodeAttributesRemoveNodeAttributesByNames(ctx context.Context, attributeNames []string) error {
	if _, err := db.conn.Exec(ctx, removeNodeAttributesByNamesQuery, attributeNames); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) NodeAttributesRemoveNodeAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeNodeAttributesByPluginIdQuery, pluginID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) NodeAttributesRemoveNodeAttributeByPluginIDAndName(
	ctx context.Context, pluginID uuid.UUID, attributeName string,
) error {
	if _, err := db.conn.Exec(
		ctx, removeNodeAttributeByPluginIdAndNameQuery, pluginID, attributeName,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}
