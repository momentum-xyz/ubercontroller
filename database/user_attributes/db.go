package user_attributes

import (
	"context"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"sync"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

const (
	getUserAttributesQuery           = `SELECT * FROM user_attribute;`
	getUserAttributeByIDQuery        = `SELECT * FROM user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND user_id = $3;`
	getUserAttributesByUserIDQuery   = `SELECT * FROM user_attribute WHERE user_id = $1;`
	getUserAttributePayloadByIDQuery = `SELECT value, options FROM user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND user_id = $3;`
	getUserAttributeValueByIDQuery   = `SELECT value FROM user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND user_id = $3;`
	getUserAttributeOptionsByIDQuery = `SELECT options FROM user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND user_id = $3;`

	getUserAttributesCountQuery = `SELECT COUNT(*) FROM user_attribute;`

	upsertUserAttributeQuery = `INSERT INTO user_attribute
									(plugin_id, attribute_name, user_id, value, options)
								VALUES
									($1, $2, $3, $4, $5)
								ON CONFLICT (plugin_id, attribute_name, user_id)
								DO UPDATE SET
									value = $4, options = $5;`

	updateUserAttributeValueQuery   = `UPDATE user_attribute SET value = $4 WHERE plugin_id = $1 AND attribute_name = $2 AND user_id = $3;`
	updateUserAttributeOptionsQuery = `UPDATE user_attribute SET options = $4 WHERE plugin_id = $1 AND attribute_name = $2 AND user_id = $3;`

	removeUserAttributeByIDQuery                 = `DELETE FROM user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND user_id = $3;`
	removeUserAttributesByNameQuery              = `DELETE FROM user_attribute WHERE attribute_name = $1;`
	removeUserAttributesByNamesQuery             = `DELETE FROM user_attribute WHERE attribute_name = ANY($1);`
	removeUserAttributesByPluginIDQuery          = `DELETE FROM user_attribute WHERE plugin_id = $1;`
	removeUserAttributesByAttributeIDQuery       = `DELETE FROM user_attribute WHERE plugin_id = $1 AND attribute_name = $2;`
	removeUserAttributesByUserIDQuery            = `DELETE FROM user_attribute WHERE user_id = $1;`
	removeUserAttributesByNameAndUserIDQuery     = `DELETE FROM user_attribute WHERE attribute_name = $1 AND user_id = $2;`
	removeUserAttributesByNamesAndUserIDQuery    = `DELETE FROM user_attribute WHERE attribute_name = ANY($1)  AND user_id = $2;`
	removeUserAttributesByPluginIDAndUserIDQuery = `DELETE FROM user_attribute WHERE plugin_id = $1 AND user_id = $2;`
)

var _ database.UserAttributesDB = (*DB)(nil)

type DB struct {
	conn   *pgxpool.Pool
	common database.CommonDB
	mu     sync.Mutex // TODO: think how to change this
}

func NewDB(conn *pgxpool.Pool, commonDB database.CommonDB) *DB {
	return &DB{
		conn:   conn,
		common: commonDB,
	}
}

func (db *DB) GetUserAttributes(ctx context.Context) ([]*entry.UserAttribute, error) {
	var attributes []*entry.UserAttribute
	if err := pgxscan.Select(ctx, db.conn, &attributes, getUserAttributesQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return attributes, nil
}

func (db *DB) GetUserAttributesByUserID(
	ctx context.Context, userID umid.UMID,
) ([]*entry.UserAttribute, error) {
	var attributes []*entry.UserAttribute
	if err := pgxscan.Select(ctx, db.conn, &attributes, getUserAttributesByUserIDQuery, userID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return attributes, nil
}

func (db *DB) GetUserAttributeByID(
	ctx context.Context, userAttributeID entry.UserAttributeID,
) (*entry.UserAttribute, error) {
	var attribute entry.UserAttribute

	if err := pgxscan.Get(
		ctx, db.conn, &attribute, getUserAttributeByIDQuery,
		userAttributeID.PluginID, userAttributeID.Name, userAttributeID.UserID,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}

	return &attribute, nil
}

func (db *DB) GetUserAttributePayloadByID(
	ctx context.Context, userAttributeID entry.UserAttributeID,
) (*entry.AttributePayload, error) {
	var payload entry.AttributePayload
	if err := pgxscan.Get(
		ctx, db.conn, &payload,
		getUserAttributePayloadByIDQuery, userAttributeID.PluginID, userAttributeID.Name, userAttributeID.UserID,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &payload, nil
}

func (db *DB) GetUserAttributeValueByID(
	ctx context.Context, userAttributeID entry.UserAttributeID,
) (*entry.AttributeValue, error) {
	var value entry.AttributeValue
	err := db.conn.QueryRow(
		ctx,
		getUserAttributeValueByIDQuery,
		userAttributeID.PluginID,
		userAttributeID.Name,
		userAttributeID.UserID,
	).Scan(&value)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}

	return &value, nil
}

func (db *DB) GetUserAttributeOptionsByID(
	ctx context.Context, userAttributeID entry.UserAttributeID,
) (*entry.AttributeOptions, error) {
	var options entry.AttributeOptions

	err := db.conn.QueryRow(
		ctx,
		getUserAttributeOptionsByIDQuery,
		userAttributeID.PluginID,
		userAttributeID.Name,
		userAttributeID.UserID,
	).Scan(&options)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}

	return &options, nil
}

func (db *DB) GetUserAttributesCount(ctx context.Context) (int64, error) {
	var count int64
	if err := db.conn.QueryRow(ctx, getUserAttributesCountQuery).
		Scan(&count); err != nil {
		return 0, errors.WithMessage(err, "failed to query db")
	}
	return count, nil
}

func (db *DB) UpsertUserAttribute(
	ctx context.Context, userAttributeID entry.UserAttributeID, modifyFn modify.Fn[entry.AttributePayload],
) (*entry.AttributePayload, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	payload, err := db.GetUserAttributePayloadByID(ctx, userAttributeID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.WithMessage(err, "failed to get attribute payload by umid")
		}
	}

	payload, err = modifyFn(payload)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify attribute payload")
	}

	var value *entry.AttributeValue
	var options *entry.AttributeOptions
	if payload != nil {
		value = payload.Value
		options = payload.Options
	}

	if _, err := db.conn.Exec(
		ctx, upsertUserAttributeQuery,
		userAttributeID.PluginID, userAttributeID.Name, userAttributeID.UserID, value, options,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to exec db")
	}

	return payload, nil
}

func (db *DB) RemoveUserAttributesByName(ctx context.Context, name string) error {
	res, err := db.conn.Exec(ctx, removeUserAttributesByNameQuery, name)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveUserAttributesByNames(ctx context.Context, names []string) error {
	res, err := db.conn.Exec(ctx, removeUserAttributesByNamesQuery, names)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveUserAttributesByPluginID(ctx context.Context, pluginID umid.UMID) error {
	res, err := db.conn.Exec(ctx, removeUserAttributesByPluginIDQuery, pluginID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveUserAttributesByAttributeID(
	ctx context.Context, attributeID entry.AttributeID,
) error {
	res, err := db.conn.Exec(ctx, removeUserAttributesByAttributeIDQuery, attributeID.PluginID, attributeID.Name)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveUserAttributesByUserID(ctx context.Context, userID umid.UMID) error {
	res, err := db.conn.Exec(ctx, removeUserAttributesByUserIDQuery, userID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveUserAttributesByNameAndUserID(
	ctx context.Context, name string, userID umid.UMID,
) error {
	res, err := db.conn.Exec(ctx, removeUserAttributesByNameAndUserIDQuery, name, userID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveUserAttributesByNamesAndUserID(
	ctx context.Context, names []string, userID umid.UMID,
) error {
	res, err := db.conn.Exec(ctx, removeUserAttributesByNamesAndUserIDQuery, names, userID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveUserAttributesByPluginIDAndUserID(
	ctx context.Context, pluginID umid.UMID, userID umid.UMID,
) error {
	res, err := db.conn.Exec(ctx, removeUserAttributesByPluginIDAndUserIDQuery, pluginID, userID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveUserAttributeByID(
	ctx context.Context, userAttributeID entry.UserAttributeID,
) error {
	res, err := db.conn.Exec(
		ctx, removeUserAttributeByIDQuery,
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

func (db *DB) UpdateUserAttributeValue(
	ctx context.Context, userAttributeID entry.UserAttributeID, modifyFn modify.Fn[entry.AttributeValue],
) (*entry.AttributeValue, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	value, err := db.GetUserAttributeValueByID(ctx, userAttributeID)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get value by umid")
	}

	value, err = modifyFn(value)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify value")
	}

	if _, err := db.conn.Exec(
		ctx, updateUserAttributeValueQuery, userAttributeID.PluginID, userAttributeID.Name, userAttributeID.UserID,
		value,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to exec db")
	}

	return value, nil
}

func (db *DB) UpdateUserAttributeOptions(
	ctx context.Context, userAttributeID entry.UserAttributeID, modifyFn modify.Fn[entry.AttributeOptions],
) (*entry.AttributeOptions, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	options, err := db.GetUserAttributeOptionsByID(ctx, userAttributeID)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get options by umid")
	}

	options, err = modifyFn(options)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify options")
	}

	if _, err := db.conn.Exec(
		ctx, updateUserAttributeOptionsQuery, userAttributeID.PluginID, userAttributeID.Name, userAttributeID.UserID,
		options,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to exec db")
	}

	return options, nil
}
