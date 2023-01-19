package space_user_attributes

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
	getSpaceUserAttributesQuery                     = `SELECT * FROM space_user_attribute;`
	getSpaceUserAttributeByIDQuery                  = `SELECT * FROM space_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND space_id = $3 AND user_id = $4;`
	getSpaceUserAttributePayloadQuery               = `SELECT value, options FROM space_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND space_id = $3 AND user_id = $4;`
	getSpaceUserAttributeValueByIDQuery             = `SELECT value FROM space_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND space_id = $3 AND user_id = $4;`
	getSpaceUserAttributeOptionsByIDQuery           = `SELECT options FROM space_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND space_id = $3 AND user_id = $4;`
	getSpaceUserAttributesBySpaceIDQuery            = `SELECT * FROM space_user_attribute WHERE space_id = $1;`
	getSpaceUserAttributesByUserIDQuery             = `SELECT * FROM space_user_attribute WHERE user_id = $1;`
	getSpaceUserAttributesBySpaceIDAndUserIDQuery   = `SELECT * FROM space_user_attribute WHERE space_id = $1 AND user_id = $2;`
	getObjectUserAttributesByObjectAttributeIDQuery = `SELECT * FROM object_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND object_id = $3;`

	getSpaceUserAttributesCountQuery = `SELECT COUNT(*) FROM space_user_attribute;`

	upsertObjectUserAttributeQuery = `INSERT INTO object_user_attribute
											(plugin_id, attribute_name, object_id, user_id, value, options)
										VALUES
											($1, $2, $3, $4, $5, $6)
										ON CONFLICT (plugin_id, attribute_name, object_id, user_id)
										DO UPDATE SET
											value = $5,options = $6;`

	updateSpaceUserAttributeValueQuery   = `UPDATE space_user_attribute SET value = $5 WHERE plugin_id = $1 AND attribute_name = $2 AND space_id = $3 AND user_id = $4;`
	updateSpaceUserAttributeOptionsQuery = `UPDATE space_user_attribute SET options = $5 WHERE plugin_id = $1 AND attribute_name = $2 AND space_id = $3 AND user_id = $4;`

	removeObjectUserAttributeByIDQuery               = `DELETE FROM object_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND object_id = $3 AND user_id = $4;`
	removeObjectUserAttributesByNameQuery            = `DELETE FROM object_user_attribute WHERE attribute_name = $1;`
	removeSpaceUserAttributesByNamesQuery            = `DELETE FROM space_user_attribute WHERE attribute_name = ANY($1);`
	removeSpaceUserAttributesByPluginIDQuery         = `DELETE FROM space_user_attribute WHERE plugin_id = $1;`
	removeObjectUserAttributesByAttributeIDQuery     = `DELETE FROM object_user_attribute WHERE plugin_id = $1 AND attribute_name = $2;`
	removeSpaceUserAttributesBySpaceIDQuery          = `DELETE FROM space_user_attribute WHERE space_id = $1;`
	removeObjectUserAttributesByNameAndObjectIDQuery = `DELETE FROM object_user_attribute WHERE attribute_name = $1 AND object_id = $2;`
	removeSpaceUserAttributesByNamesAndSpaceIDQuery  = `DELETE FROM space_user_attribute WHERE attribute_name = ANY($1) AND space_id = $2;`

	removeSpaceUserAttributesByUserIDQuery         = `DELETE FROM space_user_attribute WHERE user_id = $1;`
	removeObjectUserAttributesByNameAndUserIDQuery = `DELETE FROM object_user_attribute WHERE attribute_name = $1 AND user_id = $2;`
	removeSpaceUserAttributesByNamesAndUserIDQuery = `DELETE FROM space_user_attribute WHERE attribute_name = ANY($1) AND user_id = $2;`

	removeSpaceUserAttributesBySpaceIDAndUserIDQuery         = `DELETE FROM space_user_attribute WHERE space_id = $1 AND user_id = $2;`
	removeSpaceUserAttributesByNameAndObjectIDAndUserIDQuery = `DELETE FROM object_user_attribute WHERE attribute_name = $1 AND object_id = $2 AND user_id = $3;`
	removeSpaceUserAttributesByNamesAndSpaceIDAndUserIDQuery = `DELETE FROM space_user_attribute WHERE attribute_name = ANY($1) AND space_id = $2 AND user_id = $3;`

	removeSpaceUserAttributesByPluginIDAndSpaceIDQuery = `DELETE FROM space_user_attribute WHERE plugin_id = $1 AND space_id = $2;`
	removeObjectUserAttributesByObjectAttributeIDQuery = `DELETE FROM object_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND object_id = $3;`

	removeSpaceUserAttributesByPluginIDAndUserIDQuery           = `DELETE FROM space_user_attribute WHERE plugin_id = $1 AND user_id = $2;`
	removeObjectUserAttributesByUserAttributeIDQuery            = `DELETE FROM object_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND user_id = $3;`
	removeSpaceUserAttributesByPluginIDAndSpaceIDAndUserIDQuery = `DELETE FROM space_user_attribute WHERE plugin_id = $1 AND space_id = $2 AND user_id = $3;`
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
	if err := pgxscan.Select(ctx, db.conn, &attributes, getSpaceUserAttributesQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return attributes, nil
}

func (db *DB) GetObjectUserAttributeByID(
	ctx context.Context, spaceUserAttributeID entry.ObjectUserAttributeID,
) (*entry.ObjectUserAttribute, error) {
	var attribute entry.ObjectUserAttribute
	if err := pgxscan.Get(
		ctx, db.conn, &attribute, getSpaceUserAttributeByIDQuery,
		spaceUserAttributeID.PluginID, spaceUserAttributeID.Name, spaceUserAttributeID.ObjectID, spaceUserAttributeID.UserID,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &attribute, nil
}

func (db *DB) GetObjectUserAttributePayloadByID(
	ctx context.Context, spaceUserAttributeID entry.ObjectUserAttributeID,
) (*entry.AttributePayload, error) {
	var payload entry.AttributePayload
	if err := pgxscan.Get(ctx, db.conn, &payload, getSpaceUserAttributePayloadQuery,
		spaceUserAttributeID.PluginID, spaceUserAttributeID.Name,
		spaceUserAttributeID.ObjectID, spaceUserAttributeID.UserID,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &payload, nil
}

func (db *DB) GetObjectUserAttributeValueByID(
	ctx context.Context, spaceUserAttributeID entry.ObjectUserAttributeID,
) (*entry.AttributeValue, error) {
	var value entry.AttributeValue
	err := db.conn.QueryRow(ctx,
		getSpaceUserAttributeValueByIDQuery,
		spaceUserAttributeID.PluginID,
		spaceUserAttributeID.Name,
		spaceUserAttributeID.ObjectID,
		spaceUserAttributeID.UserID).Scan(&value)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &value, nil
}

func (db *DB) GetObjectUserAttributeOptionsByID(
	ctx context.Context, spaceUserAttributeID entry.ObjectUserAttributeID,
) (*entry.AttributeOptions, error) {
	var options entry.AttributeOptions
	err := db.conn.QueryRow(ctx,
		getSpaceUserAttributeOptionsByIDQuery,
		spaceUserAttributeID.PluginID,
		spaceUserAttributeID.Name,
		spaceUserAttributeID.ObjectID,
		spaceUserAttributeID.UserID).Scan(&options)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &options, nil
}

func (db *DB) GetObjectUserAttributesByObjectID(
	ctx context.Context, spaceID uuid.UUID,
) ([]*entry.ObjectUserAttribute, error) {
	var attributes []*entry.ObjectUserAttribute
	if err := pgxscan.Select(ctx, db.conn, &attributes, getSpaceUserAttributesBySpaceIDQuery, spaceID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return attributes, nil
}

func (db *DB) GetObjectUserAttributesByUserID(
	ctx context.Context, userID uuid.UUID,
) ([]*entry.ObjectUserAttribute, error) {
	var attributes []*entry.ObjectUserAttribute
	if err := pgxscan.Select(ctx, db.conn, &attributes, getSpaceUserAttributesByUserIDQuery, userID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return attributes, nil
}

func (db *DB) GetObjectUserAttributesByObjectIDAndUserID(
	ctx context.Context, spaceID uuid.UUID, userID uuid.UUID,
) ([]*entry.ObjectUserAttribute, error) {
	var attributes []*entry.ObjectUserAttribute
	if err := pgxscan.Select(
		ctx, db.conn, &attributes, getSpaceUserAttributesBySpaceIDAndUserIDQuery, spaceID, userID,
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
	if err := db.conn.QueryRow(ctx, getSpaceUserAttributesCountQuery).
		Scan(&count); err != nil {
		return 0, errors.WithMessage(err, "failed to query db")
	}
	return count, nil
}

func (db *DB) UpsertObjectUserAttribute(
	ctx context.Context, spaceUserAttributeID entry.ObjectUserAttributeID, modifyFn modify.Fn[entry.AttributePayload],
) (*entry.AttributePayload, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	payload, err := db.GetObjectUserAttributePayloadByID(ctx, spaceUserAttributeID)
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
		ctx, upsertObjectUserAttributeQuery, spaceUserAttributeID.PluginID, spaceUserAttributeID.Name,
		spaceUserAttributeID.ObjectID, spaceUserAttributeID.UserID,
		value, options,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to exec db")
	}

	return payload, nil
}

func (db *DB) RemoveObjectUserAttributeByName(ctx context.Context, name string) error {
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
	res, err := db.conn.Exec(ctx, removeSpaceUserAttributesByNamesQuery, names)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error {
	res, err := db.conn.Exec(ctx, removeSpaceUserAttributesByPluginIDQuery, pluginID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributeByAttributeID(
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

func (db *DB) RemoveObjectUserAttributeByObjectID(
	ctx context.Context, spaceID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeSpaceUserAttributesBySpaceIDQuery, spaceID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributeByNameAndObjectID(
	ctx context.Context, name string, spaceID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeObjectUserAttributesByNameAndObjectIDQuery, name, spaceID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributeByNamesAndObjectID(
	ctx context.Context, names []string, spaceID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeSpaceUserAttributesByNamesAndSpaceIDQuery, names, spaceID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributeByUserID(
	ctx context.Context, userID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeSpaceUserAttributesByUserIDQuery, userID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributeByNameAndUserID(
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

func (db *DB) RemoveObjectUserAttributeByNamesAndUserID(
	ctx context.Context, names []string, userID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeSpaceUserAttributesByNamesAndUserIDQuery, names, userID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributeByObjectIDAndUserID(
	ctx context.Context, spaceID uuid.UUID, userID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeSpaceUserAttributesBySpaceIDAndUserIDQuery, spaceID, userID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributeByNameAndObjectIDAndUserID(
	ctx context.Context, name string, spaceID uuid.UUID, userID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeSpaceUserAttributesByNameAndObjectIDAndUserIDQuery, name, spaceID, userID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributeByNamesAndObjectIDAndUserID(
	ctx context.Context, names []string, spaceID uuid.UUID, userID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeSpaceUserAttributesByNamesAndSpaceIDAndUserIDQuery, names, spaceID, userID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributeByPluginIDAndObjectID(
	ctx context.Context, pluginID uuid.UUID, spaceID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeSpaceUserAttributesByPluginIDAndSpaceIDQuery, pluginID, spaceID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributeByObjectAttributeID(
	ctx context.Context, spaceAttributeID entry.ObjectAttributeID,
) error {
	res, err := db.conn.Exec(
		ctx, removeObjectUserAttributesByObjectAttributeIDQuery,
		spaceAttributeID.PluginID, spaceAttributeID.Name, spaceAttributeID.ObjectID,
	)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributeByPluginIDAndUserID(
	ctx context.Context, pluginID uuid.UUID, userID uuid.UUID,
) error {
	res, err := db.conn.Exec(
		ctx, removeSpaceUserAttributesByPluginIDAndUserIDQuery, pluginID, userID,
	)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributeByUserAttributeID(
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
func (db *DB) RemoveObjectUserAttributeByPluginIDAndObjectIDAndUserID(
	ctx context.Context, pluginID uuid.UUID, spaceID uuid.UUID, userID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeSpaceUserAttributesByPluginIDAndSpaceIDAndUserIDQuery, pluginID, spaceID, userID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveObjectUserAttributeByID(
	ctx context.Context, spaceUserAttributeID entry.ObjectUserAttributeID,
) error {
	res, err := db.conn.Exec(
		ctx, removeObjectUserAttributeByIDQuery,
		spaceUserAttributeID.PluginID, spaceUserAttributeID.Name, spaceUserAttributeID.ObjectID, spaceUserAttributeID.UserID,
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
	ctx context.Context, spaceUserAttributeID entry.ObjectUserAttributeID, modifyFn modify.Fn[entry.AttributeValue],
) (*entry.AttributeValue, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	value, err := db.GetObjectUserAttributeValueByID(ctx, spaceUserAttributeID)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get value by id")
	}

	value, err = modifyFn(value)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify value")
	}

	if _, err := db.conn.Exec(
		ctx, updateSpaceUserAttributeValueQuery,
		spaceUserAttributeID.PluginID, spaceUserAttributeID.Name, spaceUserAttributeID.ObjectID, spaceUserAttributeID.UserID,
		value,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to exec db")
	}

	return value, nil
}

func (db *DB) UpdateObjectUserAttributeOptions(
	ctx context.Context, spaceUserAttributeID entry.ObjectUserAttributeID, modifyFn modify.Fn[entry.AttributeOptions],
) (*entry.AttributeOptions, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	options, err := db.GetObjectUserAttributeOptionsByID(ctx, spaceUserAttributeID)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get options by id")
	}

	options, err = modifyFn(options)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify options")
	}

	if _, err := db.conn.Exec(
		ctx, updateSpaceUserAttributeOptionsQuery,
		spaceUserAttributeID.PluginID, spaceUserAttributeID.Name, spaceUserAttributeID.ObjectID, spaceUserAttributeID.UserID,
		options,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to exec db")
	}

	return options, nil
}
