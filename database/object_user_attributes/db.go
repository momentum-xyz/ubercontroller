package object_user_attributes

import (
	"context"
	"sync"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

const (
	getObjectUserAttributesQuery                    = `SELECT * FROM object_user_attribute;`
	getObjectUserAttributeByIDQuery                 = `SELECT * FROM object_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND object_id = $3 AND user_id = $4;`
	getObjectUserAttributePayloadByIDQuery          = `SELECT value, options FROM object_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND object_id = $3 AND user_id = $4;`
	getObjectUserAttributeValueByIDQuery            = `SELECT value FROM object_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND object_id = $3 AND user_id = $4;`
	getObjectUserAttributeOptionsByIDQuery          = `SELECT options FROM object_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND object_id = $3 AND user_id = $4;`
	getObjectUserAttributesByObjectIDQuery          = `SELECT * FROM object_user_attribute WHERE object_id = $1;`
	getObjectUserAttributesByUserIDQuery            = `SELECT * FROM object_user_attribute WHERE user_id = $1;`
	getObjectUserAttributesByObjectIDAndUserIDQuery = `SELECT * FROM object_user_attribute WHERE object_id = $1 AND user_id = $2;`
	getObjectUserAttributesByObjectAttributeIDQuery = `SELECT * FROM object_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND object_id = $3;`

	getObjectUserAttributesCountQuery = `SELECT COUNT(*) FROM object_user_attribute;`

	upsertObjectUserAttributeQuery = `INSERT INTO object_user_attribute
											(plugin_id, attribute_name, object_id, user_id, value, options)
										VALUES
											($1, $2, $3, $4, $5, $6)
										ON CONFLICT (plugin_id, attribute_name, object_id, user_id)
										DO UPDATE SET
											value = $5,options = $6;`

	updateObjectUserAttributeValueQuery   = `UPDATE object_user_attribute SET value = $5 WHERE plugin_id = $1 AND attribute_name = $2 AND object_id = $3 AND user_id = $4;`
	updateObjectUserAttributeOptionsQuery = `UPDATE object_user_attribute SET options = $5 WHERE plugin_id = $1 AND attribute_name = $2 AND object_id = $3 AND user_id = $4;`

	removeObjectUserAttributeByIDQuery                = `DELETE FROM object_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND object_id = $3 AND user_id = $4;`
	removeObjectUserAttributesByNameQuery             = `DELETE FROM object_user_attribute WHERE attribute_name = $1;`
	removeObjectUserAttributesByNamesQuery            = `DELETE FROM object_user_attribute WHERE attribute_name = ANY($1);`
	removeObjectUserAttributesByPluginIDQuery         = `DELETE FROM object_user_attribute WHERE plugin_id = $1;`
	removeObjectUserAttributesByAttributeIDQuery      = `DELETE FROM object_user_attribute WHERE plugin_id = $1 AND attribute_name = $2;`
	removeObjectUserAttributesByObjectIDQuery         = `DELETE FROM object_user_attribute WHERE object_id = $1;`
	removeObjectUserAttributesByNameAndObjectIDQuery  = `DELETE FROM object_user_attribute WHERE attribute_name = $1 AND object_id = $2;`
	removeObjectUserAttributesByNamesAndObjectIDQuery = `DELETE FROM object_user_attribute WHERE attribute_name = ANY($1) AND object_id = $2;`

	removeObjectUserAttributesByUserIDQuery         = `DELETE FROM object_user_attribute WHERE user_id = $1;`
	removeObjectUserAttributesByNameAndUserIDQuery  = `DELETE FROM object_user_attribute WHERE attribute_name = $1 AND user_id = $2;`
	removeObjectUserAttributesByNamesAndUserIDQuery = `DELETE FROM object_user_attribute WHERE attribute_name = ANY($1) AND user_id = $2;`

	removeObjectUserAttributesByObjectIDAndUserIDQuery         = `DELETE FROM object_user_attribute WHERE object_id = $1 AND user_id = $2;`
	removeObjectUserAttributesByNameAndObjectIDAndUserIDQuery  = `DELETE FROM object_user_attribute WHERE attribute_name = $1 AND object_id = $2 AND user_id = $3;`
	removeObjectUserAttributesByNamesAndObjectIDAndUserIDQuery = `DELETE FROM object_user_attribute WHERE attribute_name = ANY($1) AND object_id = $2 AND user_id = $3;`

	removeObjectUserAttributesByPluginIDAndObjectIDQuery = `DELETE FROM object_user_attribute WHERE plugin_id = $1 AND object_id = $2;`
	removeObjectUserAttributesByObjectAttributeIDQuery   = `DELETE FROM object_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND object_id = $3;`

	removeObjectUserAttributesByPluginIDAndUserIDQuery            = `DELETE FROM object_user_attribute WHERE plugin_id = $1 AND user_id = $2;`
	removeObjectUserAttributesByUserAttributeIDQuery              = `DELETE FROM object_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND user_id = $3;`
	removeObjectUserAttributesByPluginIDAndObjectIDAndUserIDQuery = `DELETE FROM object_user_attribute WHERE plugin_id = $1 AND object_id = $2 AND user_id = $3;`
)

var _ database.ObjectUserAttributesDB = (*DB)(nil)

type DB struct {
	conn   *pgxpool.Pool
	common database.CommonDB
	mu     sync.Mutex
}

func NewDB(conn *pgxpool.Pool, commonDB database.CommonDB) *DB {
	return &DB{
		conn:   conn,
		common: commonDB,
	}
}

func (db *DB) GetObjectUserAttributes(ctx context.Context) ([]*entry.ObjectUserAttribute, error) {
	var attributes []*entry.ObjectUserAttribute
	if err := pgxscan.Select(ctx, db.conn, &attributes, getObjectUserAttributesQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return attributes, nil
}

func (db *DB) GetObjectUserAttributeByID(
	ctx context.Context, objectUserAttributeID entry.ObjectUserAttributeID,
) (*entry.ObjectUserAttribute, error) {
	var attribute entry.ObjectUserAttribute
	if err := pgxscan.Get(
		ctx, db.conn, &attribute, getObjectUserAttributeByIDQuery,
		objectUserAttributeID.PluginID, objectUserAttributeID.Name,
		objectUserAttributeID.ObjectID, objectUserAttributeID.UserID,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &attribute, nil
}

func (db *DB) GetObjectUserAttributePayloadByID(
	ctx context.Context, objectUserAttributeID entry.ObjectUserAttributeID,
) (*entry.AttributePayload, error) {
	var payload entry.AttributePayload
	if err := pgxscan.Get(
		ctx, db.conn, &payload, getObjectUserAttributePayloadByIDQuery,
		objectUserAttributeID.PluginID, objectUserAttributeID.Name,
		objectUserAttributeID.ObjectID, objectUserAttributeID.UserID,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &payload, nil
}

func (db *DB) GetObjectUserAttributeValueByID(
	ctx context.Context, objectUserAttributeID entry.ObjectUserAttributeID,
) (*entry.AttributeValue, error) {
	var value entry.AttributeValue
	if err := db.conn.QueryRow(
		ctx, getObjectUserAttributeValueByIDQuery,
		objectUserAttributeID.PluginID, objectUserAttributeID.Name,
		objectUserAttributeID.ObjectID, objectUserAttributeID.UserID).
		Scan(&value); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &value, nil
}

func (db *DB) GetObjectUserAttributeOptionsByID(
	ctx context.Context, objectUserAttributeID entry.ObjectUserAttributeID,
) (*entry.AttributeOptions, error) {
	var options entry.AttributeOptions
	if err := db.conn.QueryRow(
		ctx, getObjectUserAttributeOptionsByIDQuery,
		objectUserAttributeID.PluginID, objectUserAttributeID.Name,
		objectUserAttributeID.ObjectID, objectUserAttributeID.UserID).
		Scan(&options); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &options, nil
}

func (db *DB) GetObjectUserAttributesByObjectID(
	ctx context.Context, objectID uuid.UUID,
) ([]*entry.ObjectUserAttribute, error) {
	var attributes []*entry.ObjectUserAttribute
	if err := pgxscan.Select(ctx, db.conn, &attributes, getObjectUserAttributesByObjectIDQuery, objectID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return attributes, nil
}

func (db *DB) GetObjectUserAttributesByUserID(
	ctx context.Context, userID uuid.UUID,
) ([]*entry.ObjectUserAttribute, error) {
	var attributes []*entry.ObjectUserAttribute
	if err := pgxscan.Select(ctx, db.conn, &attributes, getObjectUserAttributesByUserIDQuery, userID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return attributes, nil
}

func (db *DB) GetObjectUserAttributesByObjectIDAndUserID(
	ctx context.Context, objectID uuid.UUID, userID uuid.UUID,
) ([]*entry.ObjectUserAttribute, error) {
	var attributes []*entry.ObjectUserAttribute
	if err := pgxscan.Select(
		ctx, db.conn, &attributes, getObjectUserAttributesByObjectIDAndUserIDQuery, objectID, userID,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return attributes, nil
}

func (db *DB) GetObjectUserAttributesByObjectAttributeID(
	ctx context.Context, objectAttributeID entry.ObjectAttributeID,
) ([]*entry.ObjectUserAttribute, error) {
	var attributes []*entry.ObjectUserAttribute
	if err := pgxscan.Select(
		ctx, db.conn, &attributes, getObjectUserAttributesByObjectAttributeIDQuery,
		objectAttributeID.PluginID, objectAttributeID.Name, objectAttributeID.ObjectID,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return attributes, nil
}

func (db *DB) GetObjectUserAttributesCount(ctx context.Context) (int64, error) {
	var count int64
	if err := db.conn.QueryRow(ctx, getObjectUserAttributesCountQuery).
		Scan(&count); err != nil {
		return 0, errors.WithMessage(err, "failed to query db")
	}
	return count, nil
}

func (db *DB) UpsertObjectUserAttribute(
	ctx context.Context, objectUserAttributeID entry.ObjectUserAttributeID, modifyFn modify.Fn[entry.AttributePayload],
) (*entry.AttributePayload, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	payload, err := db.GetObjectUserAttributePayloadByID(ctx, objectUserAttributeID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.WithMessage(err, "failed to get attribute payload by id")
		}
	}

	payload, err = modifyFn(payload)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify payload")
	}

	var value *entry.AttributeValue
	var options *entry.AttributeOptions
	if payload != nil {
		value = payload.Value
		options = payload.Options
	}

	if _, err := db.conn.Exec(
		ctx, upsertObjectUserAttributeQuery,
		objectUserAttributeID.PluginID, objectUserAttributeID.Name,
		objectUserAttributeID.ObjectID, objectUserAttributeID.UserID,
		value, options,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to exec db")
	}

	return payload, nil
}

func (db *DB) RemoveObjectUserAttributesByName(ctx context.Context, name string) error {
	res, err := db.conn.Exec(ctx, removeObjectUserAttributesByNameQuery, name)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributesByNames(ctx context.Context, names []string) error {
	res, err := db.conn.Exec(ctx, removeObjectUserAttributesByNamesQuery, names)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error {
	res, err := db.conn.Exec(ctx, removeObjectUserAttributesByPluginIDQuery, pluginID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributesByAttributeID(
	ctx context.Context, attributeID entry.AttributeID,
) error {
	res, err := db.conn.Exec(ctx, removeObjectUserAttributesByAttributeIDQuery, attributeID.PluginID, attributeID.Name)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributesByObjectID(
	ctx context.Context, objectID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeObjectUserAttributesByObjectIDQuery, objectID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributesByNameAndObjectID(
	ctx context.Context, name string, objectID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeObjectUserAttributesByNameAndObjectIDQuery, name, objectID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributesByNamesAndObjectID(
	ctx context.Context, names []string, objectID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeObjectUserAttributesByNamesAndObjectIDQuery, names, objectID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributesByUserID(
	ctx context.Context, userID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeObjectUserAttributesByUserIDQuery, userID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributesByNameAndUserID(
	ctx context.Context, name string, userID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeObjectUserAttributesByNameAndUserIDQuery, name, userID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributesByNamesAndUserID(
	ctx context.Context, names []string, userID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeObjectUserAttributesByNamesAndUserIDQuery, names, userID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributesByObjectIDAndUserID(
	ctx context.Context, objectID uuid.UUID, userID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeObjectUserAttributesByObjectIDAndUserIDQuery, objectID, userID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributesByNameAndObjectIDAndUserID(
	ctx context.Context, name string, objectID uuid.UUID, userID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeObjectUserAttributesByNameAndObjectIDAndUserIDQuery, name, objectID, userID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributesByNamesAndObjectIDAndUserID(
	ctx context.Context, names []string, objectID uuid.UUID, userID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeObjectUserAttributesByNamesAndObjectIDAndUserIDQuery, names, objectID, userID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributesByPluginIDAndObjectID(
	ctx context.Context, pluginID uuid.UUID, objectID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeObjectUserAttributesByPluginIDAndObjectIDQuery, pluginID, objectID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributesByObjectAttributeID(
	ctx context.Context, objectAttributeID entry.ObjectAttributeID,
) error {
	res, err := db.conn.Exec(
		ctx, removeObjectUserAttributesByObjectAttributeIDQuery,
		objectAttributeID.PluginID, objectAttributeID.Name, objectAttributeID.ObjectID,
	)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributesByPluginIDAndUserID(
	ctx context.Context, pluginID uuid.UUID, userID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeObjectUserAttributesByPluginIDAndUserIDQuery, pluginID, userID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributesByUserAttributeID(
	ctx context.Context, userAttributeID entry.UserAttributeID,
) error {
	res, err := db.conn.Exec(
		ctx, removeObjectUserAttributesByUserAttributeIDQuery,
		userAttributeID.PluginID, userAttributeID.Name, userAttributeID.UserID,
	)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
func (db *DB) RemoveObjectUserAttributesByPluginIDAndObjectIDAndUserID(
	ctx context.Context, pluginID uuid.UUID, objectID uuid.UUID, userID uuid.UUID,
) error {
	res, err := db.conn.Exec(
		ctx, removeObjectUserAttributesByPluginIDAndObjectIDAndUserIDQuery,
		pluginID, objectID, userID,
	)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributeByID(
	ctx context.Context, objectUserAttributeID entry.ObjectUserAttributeID,
) error {
	res, err := db.conn.Exec(
		ctx, removeObjectUserAttributeByIDQuery,
		objectUserAttributeID.PluginID, objectUserAttributeID.Name,
		objectUserAttributeID.ObjectID, objectUserAttributeID.UserID,
	)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) UpdateObjectUserAttributeValue(
	ctx context.Context, objectUserAttributeID entry.ObjectUserAttributeID, modifyFn modify.Fn[entry.AttributeValue],
) (*entry.AttributeValue, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	value, err := db.GetObjectUserAttributeValueByID(ctx, objectUserAttributeID)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get value by id")
	}

	value, err = modifyFn(value)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify value")
	}

	if _, err := db.conn.Exec(
		ctx, updateObjectUserAttributeValueQuery,
		objectUserAttributeID.PluginID, objectUserAttributeID.Name,
		objectUserAttributeID.ObjectID, objectUserAttributeID.UserID,
		value,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to exec db")
	}

	return value, nil
}

func (db *DB) UpdateObjectUserAttributeOptions(
	ctx context.Context, objectUserAttributeID entry.ObjectUserAttributeID, modifyFn modify.Fn[entry.AttributeOptions],
) (*entry.AttributeOptions, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	options, err := db.GetObjectUserAttributeOptionsByID(ctx, objectUserAttributeID)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get options by id")
	}

	options, err = modifyFn(options)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify options")
	}

	if _, err := db.conn.Exec(
		ctx, updateObjectUserAttributeOptionsQuery,
		objectUserAttributeID.PluginID, objectUserAttributeID.Name,
		objectUserAttributeID.ObjectID, objectUserAttributeID.UserID,
		options,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to exec db")
	}

	return options, nil
}
