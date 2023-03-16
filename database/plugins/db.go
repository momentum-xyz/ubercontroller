package plugins

import (
	"context"
	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/entry"
)

const (
	getPluginsQuery = `SELECT * FROM plugin;`

	upsertPluginQuery = `INSERT INTO plugin
								(plugin_id, meta, options, created_at, updated_at)
							VALUES
								($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
							ON CONFLICT (plugin_id)
							DO UPDATE SET
								meta = $2, options = $3, updated_at = CURRENT_TIMESTAMP;`

	updatePluginMetaQuery    = `UPDATE plugin SET meta = $2, updated_at = CURRENT_TIMESTAMP WHERE plugin_id = $1;`
	updatePluginOptionsQuery = `UPDATE plugin SET options = $2, updated_at = CURRENT_TIMESTAMP WHERE plugin_id = $1;`

	removePluginByIDQuery   = `DELETE FROM plugin WHERE plugin_id = $1;`
	removePluginsByIDsQuery = `DELETE FROM plugin WHERE plugin_id = ANY($1);`
)

var _ database.PluginsDB = (*DB)(nil)

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

func (db *DB) GetPlugins(ctx context.Context) ([]*entry.Plugin, error) {
	var plugins []*entry.Plugin
	if err := pgxscan.Select(ctx, db.conn, &plugins, getPluginsQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return plugins, nil
}

func (db *DB) UpsertPlugin(ctx context.Context, plugin *entry.Plugin) error {
	if _, err := db.conn.Exec(
		ctx, upsertPluginQuery, plugin.PluginID, plugin.Meta, plugin.Options,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpsertPlugins(ctx context.Context, plugins []*entry.Plugin) error {
	batch := &pgx.Batch{}
	for _, plugin := range plugins {
		batch.Queue(
			upsertPluginQuery, plugin.PluginID, plugin.Meta, plugin.Options,
		)
	}

	batchRes := db.conn.SendBatch(ctx, batch)
	defer batchRes.Close()

	var errs *multierror.Error
	for i := 0; i < batch.Len(); i++ {
		if _, err := batchRes.Exec(); err != nil {
			errs = multierror.Append(
				errs, errors.WithMessagef(err, "failed to exec db for: %s", plugins[i].PluginID),
			)
		}
	}

	return errs.ErrorOrNil()
}

func (db *DB) RemovePluginByID(ctx context.Context, PluginID umid.UMID) error {
	if _, err := db.conn.Exec(ctx, removePluginByIDQuery, PluginID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) RemovePluginsByIDs(ctx context.Context, PluginIDs []umid.UMID) error {
	if _, err := db.conn.Exec(ctx, removePluginsByIDsQuery, PluginIDs); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdatePluginMeta(ctx context.Context, pluginID umid.UMID, meta entry.PluginMeta) error {
	if _, err := db.conn.Exec(ctx, updatePluginMetaQuery, pluginID, meta); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdatePluginOptions(
	ctx context.Context, pluginID umid.UMID, options *entry.PluginOptions,
) error {
	if _, err := db.conn.Exec(ctx, updatePluginOptionsQuery, pluginID, options); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}
