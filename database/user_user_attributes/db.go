package user_user_attributes

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
	getUserUserAttributesQuery                    = `SELECT * FROM user_user_attribute;`
	getUserUserAttributesQueryBySourceUserIDQuery = `SELECT * FROM user_user_attribute WHERE source_user_id = $1;`
	getUserUserAttributesQueryByTargetUserIDQuery = `SELECT * FROM user_user_attribute WHERE target_user_id = $1;`

	getUserUserAttributePayloadByIDQuery = `SELECT value, options FROM user_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND source_user_id = $3 AND target_user_id = $4;`
	getUserUserAttributeValueByIDQuery   = `SELECT value FROM user_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND source_user_id = $3 AND target_user_id = $4;`
	getUserUserAttributeOptionsByIDQuery = `SELECT options FROM user_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND source_user_id = $3 AND target_user_id = $4;`

	getUserUserAttributeByIDQuery                           = `SELECT * FROM user_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND source_user_id = $3 AND target_user_id = $4;`
	getUserUserAttributesBySourceUserIDAndTargetUserIDQuery = `SELECT * FROM user_user_attribute WHERE source_user_id = $1 AND target_user_id = $2;`

	getUserUserAttributesCountQuery = `SELECT COUNT(*) FROM user_user_attribute;`

	removeUserUserAttributeByNameQuery                  = `DELETE FROM user_user_attribute WHERE attribute_name = $1;`
	removeUserUserAttributesByNamesQuery                = `DELETE FROM user_user_attribute WHERE attribute_name = ANY($1);`
	removeUserUserAttributesByPluginIDQuery             = `DELETE FROM user_user_attribute WHERE plugin_id = $1;`
	removeUserUserAttributeByPluginIDAndNameQuery       = `DELETE FROM user_user_attribute WHERE plugin_id = $1 AND attribute_name = $2;`
	removeUserUserAttributesBySourceUserIDQuery         = `DELETE FROM user_user_attribute WHERE source_user_id = $1;`
	removeUserUserAttributeByNameAndSourceUserIDQuery   = `DELETE FROM user_user_attribute WHERE attribute_name = $1 and source_user_id = $2;`
	removeUserUserAttributesByNamesAndSourceUserIDQuery = `DELETE FROM user_user_attribute WHERE attribute_name = ANY($1) AND source_user_id = $2;`

	removeUserUserAttributesByTargetUserIDQuery         = `DELETE FROM user_user_attribute WHERE target_user_id = $1;`
	removeUserUserAttributeByNameAndTargetUserIDQuery   = `DELETE FROM user_user_attribute WHERE attribute_name = $1 AND target_user_id = $2;`
	removeUserUserAttributesByNamesAndTargetUserIDQuery = `DELETE FROM user_user_attribute WHERE attribute_name = ANY($1) AND target_user_id = $2;`

	removeUserUserAttributesBySourceUserIDAndTargetUserIDQuery         = `DELETE FROM user_user_attribute WHERE source_user_id = $1 AND target_user_id = $2;`
	removeUserUserAttributeByNameAndSourceUserIDAndTargetUserIDQuery   = `DELETE FROM user_user_attribute WHERE attribute_name = $1 AND source_user_id = $2 AND target_user_id = $3;`
	removeUserUserAttributesByNamesAndSourceUserIDAndTargetUserIDQuery = `DELETE FROM user_user_attribute WHERE attribute_name = ANY($1) AND source_user_id = $2 AND target_user_id = $3;`

	removeUserUserAttributesByPluginIDAndSourceUserIDQuery        = `DELETE FROM user_user_attribute WHERE plugin_id = $1 AND source_user_id = $2;`
	removeUserUserAttributesByPluginIDAndNameAndSourceUserIDQuery = `DELETE FROM user_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND source_user_id = $3;`

	removeUserUserAttributesByPluginIDAndTargetUserIDQuery        = `DELETE FROM user_user_attribute WHERE plugin_id = $1 AND target_user_id = $2;`
	removeUserUserAttributesByPluginIDAndNameAndTargetUserIDQuery = `DELETE FROM user_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND target_user_id = $3;`

	removeUserUserAttributesByIDQuery                                     = `DELETE FROM user_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND source_user_id = $3 AND target_user_id = $4;`
	removeUserUserAttributesByPluginIDAndSourceUserIDAndTargetUserIDQuery = `DELETE FROM user_user_attribute WHERE plugin_id = $1 AND source_user_id = $2 AND target_user_id = $3;`

	updateUserUserAttributeValueQuery   = `UPDATE user_user_attribute SET value = $5 WHERE plugin_id = $1 AND attribute_name = $2 AND source_user_id = $3 AND target_user_id = $4;`
	updateUserUserAttributeOptionsQuery = `UPDATE user_user_attribute SET options = $5 WHERE plugin_id = $1 AND attribute_name = $2 AND source_user_id = $3 AND target_user_id = $4;`

	upsertUserUserAttributeQuery = `INSERT INTO user_user_attribute
										(plugin_id, attribute_name, source_user_id, target_user_id, value, options)
									VALUES
										($1, $2, $3, $4, $5, $6)
									ON CONFLICT (plugin_id, attribute_name, source_user_id, target_user_id)
									DO UPDATE SET
									    value = $5, options = $6;`
)

var _ database.UserUserAttributesDB = (*DB)(nil)

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

func (db *DB) UserUserAttributesGetUserUserAttributes(ctx context.Context) ([]*entry.UserUserAttribute, error) {
	var attributes []*entry.UserUserAttribute
	if err := pgxscan.Select(ctx, db.conn, &attributes, getUserUserAttributesQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return attributes, nil
}

func (db *DB) UserUserAttributesGetUserUserAttributeByID(
	ctx context.Context, userUserAttributeID entry.UserUserAttributeID,
) (*entry.UserUserAttribute, error) {
	var attribute entry.UserUserAttribute
	if err := pgxscan.Get(
		ctx, db.conn, &attribute, getUserUserAttributeByIDQuery,
		userUserAttributeID.PluginID, userUserAttributeID.Name,
		userUserAttributeID.SourceUserID, userUserAttributeID.TargetUserID,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &attribute, nil
}

func (db *DB) UserUserAttributesGetUserUserAttributePayloadByID(
	ctx context.Context, userUserAttributeID entry.UserUserAttributeID,
) (*entry.AttributePayload, error) {
	var payload entry.AttributePayload
	if err := pgxscan.Get(ctx, db.conn, &payload,
		getUserUserAttributePayloadByIDQuery,
		userUserAttributeID.PluginID,
		userUserAttributeID.Name,
		userUserAttributeID.SourceUserID,
		userUserAttributeID.TargetUserID,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &payload, nil
}

func (db *DB) UserUserAttributesGetUserUserAttributeValueByID(
	ctx context.Context, userUserAttributeID entry.UserUserAttributeID,
) (*entry.AttributeValue, error) {
	var value entry.AttributeValue
	err := db.conn.QueryRow(ctx,
		getUserUserAttributeValueByIDQuery,
		userUserAttributeID.PluginID,
		userUserAttributeID.Name,
		userUserAttributeID.SourceUserID,
		userUserAttributeID.TargetUserID).Scan(&value)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &value, nil
}

func (db *DB) UserUserAttributesGetUserUserAttributeOptionsByID(
	ctx context.Context, userUserAttributeID entry.UserUserAttributeID,
) (*entry.AttributeOptions, error) {
	var options entry.AttributeOptions
	err := db.conn.QueryRow(ctx,
		getUserUserAttributeOptionsByIDQuery,
		userUserAttributeID.PluginID,
		userUserAttributeID.Name,
		userUserAttributeID.SourceUserID,
		userUserAttributeID.TargetUserID).Scan(&options)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &options, nil
}

func (db *DB) UserUserAttributesGetUserUserAttributesBySourceUserID(
	ctx context.Context, sourceUserID uuid.UUID,
) ([]*entry.UserUserAttribute, error) {
	var attributes []*entry.UserUserAttribute
	if err := pgxscan.Select(
		ctx, db.conn, &attributes, getUserUserAttributesQueryBySourceUserIDQuery, sourceUserID,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return attributes, nil
}

func (db *DB) UserUserAttributesGetUserUserAttributesByTargetUserID(
	ctx context.Context, targetUserID uuid.UUID,
) ([]*entry.UserUserAttribute, error) {
	var attributes []*entry.UserUserAttribute
	if err := pgxscan.Select(
		ctx, db.conn, &attributes, getUserUserAttributesQueryByTargetUserIDQuery, targetUserID,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return attributes, nil
}

func (db *DB) UserUserAttributesGetUserUserAttributesBySourceUserIDAndTargetUserID(
	ctx context.Context, sourceUserID uuid.UUID, targetUserID uuid.UUID,
) ([]*entry.UserUserAttribute, error) {
	var attributes []*entry.UserUserAttribute
	if err := pgxscan.Select(
		ctx, db.conn, &attributes, getUserUserAttributesBySourceUserIDAndTargetUserIDQuery, sourceUserID, targetUserID,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return attributes, nil
}

func (db *DB) UserUserAttributesGetUserUserAttributesCount(ctx context.Context) (int64, error) {
	var count int64
	if err := db.conn.QueryRow(ctx, getUserUserAttributesCountQuery).
		Scan(&count); err != nil {
		return 0, errors.WithMessage(err, "failed to query db")
	}
	return count, nil
}

func (db *DB) UserUserAttributesUpsertUserUserAttribute(
	ctx context.Context, userUserAttributeID entry.UserUserAttributeID, modifyFn modify.Fn[entry.AttributePayload],
) (*entry.AttributePayload, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	payload, err := db.UserUserAttributesGetUserUserAttributePayloadByID(ctx, userUserAttributeID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.WithMessage(err, "failed to get attribute payload by id")
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
		ctx, upsertUserUserAttributeQuery, userUserAttributeID.PluginID, userUserAttributeID.Name,
		userUserAttributeID.SourceUserID, userUserAttributeID.TargetUserID,
		value, options,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to exec db")
	}

	return payload, nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeByName(ctx context.Context, name string) error {
	res, err := db.conn.Exec(ctx, removeUserUserAttributeByNameQuery, name)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributesByNames(ctx context.Context, names []string) error {
	res, err := db.conn.Exec(ctx, removeUserUserAttributesByNamesQuery, names)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error {
	res, err := db.conn.Exec(ctx, removeUserUserAttributesByPluginIDQuery, pluginID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeByAttributeID(ctx context.Context, attributeID entry.AttributeID) error {
	res, err := db.conn.Exec(ctx, removeUserUserAttributeByPluginIDAndNameQuery, attributeID.PluginID, attributeID.Name)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeBySourceUserID(ctx context.Context, sourceUserID uuid.UUID) error {
	res, err := db.conn.Exec(ctx, removeUserUserAttributesBySourceUserIDQuery, sourceUserID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeByNameAndSourceUserID(
	ctx context.Context, name string, sourceUserID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeUserUserAttributeByNameAndSourceUserIDQuery, name, sourceUserID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeByNamesAndSourceUserID(
	ctx context.Context, names []string, sourceUserID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeUserUserAttributesByNamesAndSourceUserIDQuery, names, sourceUserID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeByTargetUserID(ctx context.Context, targetUserID uuid.UUID) error {
	res, err := db.conn.Exec(ctx, removeUserUserAttributesByTargetUserIDQuery, targetUserID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeByNameAndTargetUserID(
	ctx context.Context, name string, targetUserID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeUserUserAttributeByNameAndTargetUserIDQuery, name, targetUserID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeByNamesAndTargetUserID(
	ctx context.Context, names []string, targetUserID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeUserUserAttributesByNamesAndTargetUserIDQuery, names, targetUserID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeBySourceUserIDAndTargetUserID(
	ctx context.Context, sourceUserID uuid.UUID, targetUserID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeUserUserAttributesBySourceUserIDAndTargetUserIDQuery, sourceUserID, targetUserID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeByNameAndSourceUserIDAndTargetUserID(
	ctx context.Context, name string, sourceUserID uuid.UUID, targetUserID uuid.UUID,
) error {
	res, err := db.conn.Exec(
		ctx, removeUserUserAttributeByNameAndSourceUserIDAndTargetUserIDQuery, name, sourceUserID, targetUserID,
	)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeByNamesAndSourceUserIDAndTargetUserID(
	ctx context.Context, names []string, sourceUserID uuid.UUID, targetUserID uuid.UUID,
) error {
	res, err := db.conn.Exec(
		ctx, removeUserUserAttributesByNamesAndSourceUserIDAndTargetUserIDQuery, names, sourceUserID, targetUserID,
	)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeByPluginIDAndSourceUserID(
	ctx context.Context, pluginID uuid.UUID, sourceUserID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeUserUserAttributesByPluginIDAndSourceUserIDQuery, pluginID, sourceUserID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeBySourceUserAttributeID(
	ctx context.Context, sourceUserAttributeID entry.UserAttributeID,
) error {
	res, err := db.conn.Exec(
		ctx, removeUserUserAttributesByPluginIDAndNameAndSourceUserIDQuery,
		sourceUserAttributeID.PluginID, sourceUserAttributeID.Name, sourceUserAttributeID.UserID,
	)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeByPluginIDAndTargetUserID(
	ctx context.Context, pluginID uuid.UUID, targetUserID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeUserUserAttributesByPluginIDAndTargetUserIDQuery, pluginID, targetUserID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeByTargetUserAttributeID(
	ctx context.Context, targetUserAttributeID entry.UserAttributeID,
) error {
	res, err := db.conn.Exec(
		ctx, removeUserUserAttributesByPluginIDAndNameAndTargetUserIDQuery,
		targetUserAttributeID.PluginID, targetUserAttributeID.Name, targetUserAttributeID.UserID,
	)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
func (db *DB) UserUserAttributesRemoveUserUserAttributeByPluginIDAndSourceUserIDAndTargetUserID(
	ctx context.Context, pluginId uuid.UUID, sourceUserID uuid.UUID, targetUserID uuid.UUID,
) error {
	res, err := db.conn.Exec(
		ctx, removeUserUserAttributesByPluginIDAndSourceUserIDAndTargetUserIDQuery, pluginId, sourceUserID, targetUserID,
	)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeByID(
	ctx context.Context, userUserAttributeID entry.UserUserAttributeID,
) error {
	res, err := db.conn.Exec(
		ctx, removeUserUserAttributesByIDQuery,
		userUserAttributeID.PluginID, userUserAttributeID.Name,
		userUserAttributeID.SourceUserID, userUserAttributeID.TargetUserID,
	)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) UserUserAttributesUpdateUserUserAttributeValue(
	ctx context.Context, userUserAttributeID entry.UserUserAttributeID, modifyFn modify.Fn[entry.AttributeValue],
) (*entry.AttributeValue, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	value, err := db.UserUserAttributesGetUserUserAttributeValueByID(ctx, userUserAttributeID)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get value by id")
	}

	value, err = modifyFn(value)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify value")
	}

	if _, err := db.conn.Exec(
		ctx, updateUserUserAttributeValueQuery,
		userUserAttributeID.PluginID, userUserAttributeID.Name,
		userUserAttributeID.SourceUserID, userUserAttributeID.TargetUserID,
		value,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to exec db")
	}

	return value, nil
}

func (db *DB) UserUserAttributesUpdateUserUserAttributeOptions(
	ctx context.Context, userUserAttributeID entry.UserUserAttributeID, modifyFn modify.Fn[entry.AttributeOptions],
) (*entry.AttributeOptions, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	options, err := db.UserUserAttributesGetUserUserAttributeOptionsByID(ctx, userUserAttributeID)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get options by id")
	}

	options, err = modifyFn(options)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify options")
	}

	if _, err := db.conn.Exec(
		ctx, updateUserUserAttributeOptionsQuery,
		userUserAttributeID.PluginID, userUserAttributeID.Name,
		userUserAttributeID.SourceUserID, userUserAttributeID.TargetUserID,
		options,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to exec db")
	}

	return options, nil
}
