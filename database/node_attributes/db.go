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
	getNodeAttributesQuery = `SELECT * FROM node_attributes;`

	updateNodeAttributeValueQuery             = `UPDATE node_attributes SET value = $3 WHERE plugin_id=$1 and attribute_name = $2;`
	removeNodeAttributeByNameQuery            = `DELETE FROM node_attributes WHERE attribute_name = $1;`
	removeNodeAttributesByNamesQuery          = `DELETE FROM node_attributes WHERE attribute_name IN ($1);`
	removeNodeAttributesByPluginIdQuery       = `DELETE FROM node_attributes WHERE plugin_id = $1;`
	removeNodeAttributeByPluginIdAndNameQuery = `DELETE FROM node_attributes WHERE plugin_id = $1 and attribute_name =$2;`

	upsertNodeAttributeQuery = `INSERT INTO node_attributes
									(plugin_id, node_attribute_name,value)
								VALUES
									($1, $2, $3)
								ON CONFLICT (plugin_id,attribute_name)
								DO UPDATE SET
									value = $3;`
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

func (db *DB) NodeAttributesUpsertNodeAttribute(ctx context.Context, nodeAttribute *entry.NodeAttribute) error {
	if _, err := db.conn.Exec(
		ctx, upsertNodeAttributeQuery, nodeAttribute.PluginID, nodeAttribute.Name, nodeAttribute.Value,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) NodeAttributesUpsertNodeAttributes(ctx context.Context, nodeAttributes []*entry.NodeAttribute) error {
	batch := &pgx.Batch{}
	for _, nodeAttribute := range nodeAttributes {
		batch.Queue(
			upsertNodeAttributeQuery, nodeAttribute.PluginID, nodeAttribute.Name, nodeAttribute.Value,
		)
	}

	batchRes := db.conn.SendBatch(ctx, batch)
	defer batchRes.Close()

	var errs *multierror.Error
	for i := 0; i < batch.Len(); i++ {
		if _, err := batchRes.Exec(); err != nil {
			errs = multierror.Append(
				errs, errors.WithMessagef(err, "failed to exec db for: %v", nodeAttributes[i].Name),
			)
		}
	}

	return errs.ErrorOrNil()
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

func (db *DB) NodeAttributesUpdateNodeAttributeValue(
	ctx context.Context, pluginID uuid.UUID, attributeName string, nodeId uuid.UUID, value *entry.AttributeValue,
) error {
	if _, err := db.conn.Exec(
		ctx, updateNodeAttributeValueQuery, attributeName, pluginID, nodeId, value,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}
