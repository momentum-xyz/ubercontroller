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
	getSpaceUserAttributesQuery                            = `SELECT * FROM space_user_attribute;`
	getSpaceUserAttributeByIDQuery                         = `SELECT * FROM space_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND space_id = $3 AND user_id = $4;`
	getSpaceUserAttributePayloadQuery                      = `SELECT value, options FROM space_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND space_id = $3 AND user_id = $4;`
	getSpaceUserAttributeValueByIDQuery                    = `SELECT value FROM space_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND space_id = $3 AND user_id = $4;`
	getSpaceUserAttributeOptionsByIDQuery                  = `SELECT options FROM space_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND space_id = $3 AND user_id = $4;`
	getSpaceUserAttributesBySpaceIDQuery                   = `SELECT * FROM space_user_attribute WHERE space_id = $1;`
	getSpaceUserAttributesByUserIDQuery                    = `SELECT * FROM space_user_attribute WHERE user_id = $1;`
	getSpaceUserAttributesBySpaceIDAndUserIDQuery          = `SELECT * FROM space_user_attribute WHERE space_id = $1 AND user_id = $2;`
	getSpaceUserAttributesByPluginIDAndNameAndSpaceIDQuery = `SELECT * FROM space_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND space_id = $3;`

	getSpaceUserAttributesCountQuery = `SELECT COUNT(*) FROM space_user_attribute;`

	removeSpaceUserAttributeByNameQuery             = `DELETE FROM space_user_attribute WHERE attribute_name = $1;`
	removeSpaceUserAttributesByNamesQuery           = `DELETE FROM space_user_attribute WHERE attribute_name = ANY($1);`
	removeSpaceUserAttributesByPluginIDQuery        = `DELETE FROM space_user_attribute WHERE plugin_id = $1;`
	removeSpaceUserAttributeByPluginIDAndNameQuery  = `DELETE FROM space_user_attribute WHERE plugin_id = $1 AND attribute_name = $2;`
	removeSpaceUserAttributesBySpaceIDQuery         = `DELETE FROM space_user_attribute WHERE space_id = $1;`
	removeSpaceUserAttributeByNameAndSpaceIDQuery   = `DELETE FROM space_user_attribute WHERE attribute_name = $1 AND space_id = $2;`
	removeSpaceUserAttributesByNamesAndSpaceIDQuery = `DELETE FROM space_user_attribute WHERE attribute_name = ANY($1) AND space_id = $2;`

	removeSpaceUserAttributesByUserIDQuery         = `DELETE FROM space_user_attribute WHERE user_id = $1;`
	removeSpaceUserAttributeByNameAndUserIDQuery   = `DELETE FROM space_user_attribute WHERE attribute_name = $1 AND user_id = $2;`
	removeSpaceUserAttributesByNamesAndUserIDQuery = `DELETE FROM space_user_attribute WHERE attribute_name = ANY($1) AND user_id = $2;`

	removeSpaceUserAttributesBySpaceIDAndUserIDQuery         = `DELETE FROM space_user_attribute WHERE space_id = $1 AND user_id = $2;`
	removeSpaceUserAttributeByNameAndSpaceIDAndUserIDQuery   = `DELETE FROM space_user_attribute WHERE attribute_name = $1 AND space_id = $2 AND user_id = $3;`
	removeSpaceUserAttributesByNamesAndSpaceIDAndUserIDQuery = `DELETE FROM space_user_attribute WHERE attribute_name = ANY($1) AND space_id = $2 AND user_id = $3;`

	removeSpaceUserAttributesByPluginIDAndSpaceIDQuery        = `DELETE FROM space_user_attribute WHERE plugin_id = $1 AND space_id = $2;`
	removeSpaceUserAttributesByPluginIDAndNameAndSpaceIDQuery = `DELETE FROM space_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND space_id = $3;`

	removeSpaceUserAttributesByPluginIDAndUserIDQuery        = `DELETE FROM space_user_attribute WHERE plugin_id = $1 AND user_id = $2;`
	removeSpaceUserAttributesByPluginIDAndNameAndUserIDQuery = `DELETE FROM space_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND user_id = $3;`

	removeSpaceUserAttributesByPluginIDAndSpaceIDAndUserIDQuery        = `DELETE FROM space_user_attribute WHERE plugin_id = $1 AND space_id = $2 AND user_id = $3;`
	removeSpaceUserAttributesByPluginIDAndNameAndSpaceIDAndUserIDQuery = `DELETE FROM space_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND space_id = $3 AND user_id = $4;`

	updateSpaceUserAttributeValueQuery   = `UPDATE space_user_attribute SET value = $5 WHERE plugin_id = $1 AND attribute_name = $2 AND space_id = $3 AND user_id = $4;`
	updateSpaceUserAttributeOptionsQuery = `UPDATE space_user_attribute SET options = $5 WHERE plugin_id = $1 AND attribute_name = $2 AND space_id = $3 AND user_id = $4;`

	upsertSpaceUserAttributeQuery = `INSERT INTO space_user_attribute
											(plugin_id, attribute_name, space_id, user_id, value, options)
										VALUES
											($1, $2, $3, $4, $5, $6)
										ON CONFLICT (plugin_id, attribute_name, space_id, user_id)
										DO UPDATE SET
											value = $5,options = $6;`
)

var _ database.SpaceUserAttributesDB = (*DB)(nil)

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

func (db *DB) SpaceUserAttributesGetSpaceUserAttributes(ctx context.Context) ([]*entry.SpaceUserAttribute, error) {
	var attributes []*entry.SpaceUserAttribute
	if err := pgxscan.Select(ctx, db.conn, &attributes, getSpaceUserAttributesQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return attributes, nil
}

func (db *DB) SpaceUserAttributesGetSpaceUserAttributeByID(
	ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID,
) (*entry.SpaceUserAttribute, error) {
	var attribute entry.SpaceUserAttribute
	if err := pgxscan.Get(
		ctx, db.conn, &attribute, getSpaceUserAttributeByIDQuery,
		spaceUserAttributeID.PluginID, spaceUserAttributeID.Name, spaceUserAttributeID.SpaceID, spaceUserAttributeID.UserID,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &attribute, nil
}

func (db *DB) SpaceUserAttributesGetSpaceUserAttributePayloadByID(
	ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID,
) (*entry.AttributePayload, error) {
	var payload entry.AttributePayload
	if err := pgxscan.Get(ctx, db.conn, &payload, getSpaceUserAttributePayloadQuery,
		spaceUserAttributeID.PluginID, spaceUserAttributeID.Name,
		spaceUserAttributeID.SpaceID, spaceUserAttributeID.UserID,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &payload, nil
}

func (db *DB) SpaceUserAttributesGetSpaceUserAttributeValueByID(
	ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID,
) (*entry.AttributeValue, error) {
	var value entry.AttributeValue
	err := db.conn.QueryRow(ctx,
		getSpaceUserAttributeValueByIDQuery,
		spaceUserAttributeID.PluginID,
		spaceUserAttributeID.Name,
		spaceUserAttributeID.SpaceID,
		spaceUserAttributeID.UserID).Scan(&value)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &value, nil
}

func (db *DB) SpaceUserAttributesGetSpaceUserAttributeOptionsByID(
	ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID,
) (*entry.AttributeOptions, error) {
	var options entry.AttributeOptions
	err := db.conn.QueryRow(ctx,
		getSpaceUserAttributeOptionsByIDQuery,
		spaceUserAttributeID.PluginID,
		spaceUserAttributeID.Name,
		spaceUserAttributeID.SpaceID,
		spaceUserAttributeID.UserID).Scan(&options)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &options, nil
}

func (db *DB) SpaceUserAttributesGetSpaceUserAttributesBySpaceID(
	ctx context.Context, spaceID uuid.UUID,
) ([]*entry.SpaceUserAttribute, error) {
	var attributes []*entry.SpaceUserAttribute
	if err := pgxscan.Select(ctx, db.conn, &attributes, getSpaceUserAttributesBySpaceIDQuery, spaceID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return attributes, nil
}

func (db *DB) SpaceUserAttributesGetSpaceUserAttributesByUserID(
	ctx context.Context, userID uuid.UUID,
) ([]*entry.SpaceUserAttribute, error) {
	var attributes []*entry.SpaceUserAttribute
	if err := pgxscan.Select(ctx, db.conn, &attributes, getSpaceUserAttributesByUserIDQuery, userID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return attributes, nil
}

func (db *DB) SpaceUserAttributesGetSpaceUserAttributesBySpaceIDAndUserID(
	ctx context.Context, spaceID uuid.UUID, userID uuid.UUID,
) ([]*entry.SpaceUserAttribute, error) {
	var attributes []*entry.SpaceUserAttribute
	if err := pgxscan.Select(
		ctx, db.conn, &attributes, getSpaceUserAttributesBySpaceIDAndUserIDQuery, spaceID, userID,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return attributes, nil
}

func (db *DB) SpaceUserAttributesGetSpaceUserAttributesByPluginIDAndNameAndSpaceID(
	ctx context.Context, pluginID uuid.UUID, name string, spaceID uuid.UUID,
) ([]*entry.SpaceUserAttribute, error) {
	var attributes []*entry.SpaceUserAttribute
	if err := pgxscan.Select(ctx, db.conn, &attributes, getSpaceUserAttributesByPluginIDAndNameAndSpaceIDQuery,
		pluginID, name, spaceID,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return attributes, nil
}

func (db *DB) SpaceUserAttributesGetSpaceUserAttributesCount(ctx context.Context) (int64, error) {
	var count int64
	if err := db.conn.QueryRow(ctx, getSpaceUserAttributesCountQuery).
		Scan(&count); err != nil {
		return 0, errors.WithMessage(err, "failed to query db")
	}
	return count, nil
}

func (db *DB) SpaceUserAttributesUpsertSpaceUserAttribute(
	ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID, modifyFn modify.Fn[entry.AttributePayload],
) (*entry.AttributePayload, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	payload, err := db.SpaceUserAttributesGetSpaceUserAttributePayloadByID(ctx, spaceUserAttributeID)
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
		ctx, upsertSpaceUserAttributeQuery, spaceUserAttributeID.PluginID, spaceUserAttributeID.Name,
		spaceUserAttributeID.SpaceID, spaceUserAttributeID.UserID,
		value, options,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to exec db")
	}

	return payload, nil
}

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByName(ctx context.Context, name string) error {
	res, err := db.conn.Exec(ctx, removeSpaceUserAttributeByNameQuery, name)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributesByNames(ctx context.Context, names []string) error {
	res, err := db.conn.Exec(ctx, removeSpaceUserAttributesByNamesQuery, names)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error {
	res, err := db.conn.Exec(ctx, removeSpaceUserAttributesByPluginIDQuery, pluginID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByAttributeID(
	ctx context.Context, attributeID entry.AttributeID,
) error {
	res, err := db.conn.Exec(ctx, removeSpaceUserAttributeByPluginIDAndNameQuery, attributeID.PluginID, attributeID.Name)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeBySpaceID(
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

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByNameAndSpaceID(
	ctx context.Context, name string, spaceID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeSpaceUserAttributeByNameAndSpaceIDQuery, name, spaceID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByNamesAndSpaceID(
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

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByUserID(
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

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByNameAndUserID(
	ctx context.Context, name string, userID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeSpaceUserAttributeByNameAndUserIDQuery, name, userID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByNamesAndUserID(
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

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeBySpaceIDAndUserID(
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

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByNameAndSpaceIDAndUserID(
	ctx context.Context, name string, spaceID uuid.UUID, userID uuid.UUID,
) error {
	res, err := db.conn.Exec(ctx, removeSpaceUserAttributeByNameAndSpaceIDAndUserIDQuery, name, spaceID, userID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByNamesAndSpaceIDAndUserID(
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

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByPluginIDAndSpaceID(
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

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeBySpaceAttributeID(
	ctx context.Context, spaceAttributeID entry.SpaceAttributeID,
) error {
	res, err := db.conn.Exec(
		ctx, removeSpaceUserAttributesByPluginIDAndNameAndSpaceIDQuery,
		spaceAttributeID.PluginID, spaceAttributeID.Name, spaceAttributeID.SpaceID,
	)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByPluginIDAndUserID(
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

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByUserAttributeID(
	ctx context.Context, userAttributeID entry.UserAttributeID,
) error {
	res, err := db.conn.Exec(
		ctx, removeSpaceUserAttributesByPluginIDAndNameAndUserIDQuery,
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
func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByPluginIDAndSpaceIDAndUserID(
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

func (db *DB) SpaceUserAttributesRemoveSpaceUserAttributeByID(
	ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID,
) error {
	res, err := db.conn.Exec(
		ctx, removeSpaceUserAttributesByPluginIDAndNameAndSpaceIDAndUserIDQuery,
		spaceUserAttributeID.PluginID, spaceUserAttributeID.Name, spaceUserAttributeID.SpaceID, spaceUserAttributeID.UserID,
	)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) SpaceUserAttributesUpdateSpaceUserAttributeValue(
	ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID, modifyFn modify.Fn[entry.AttributeValue],
) (*entry.AttributeValue, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	value, err := db.SpaceUserAttributesGetSpaceUserAttributeValueByID(ctx, spaceUserAttributeID)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get value by id")
	}

	value, err = modifyFn(value)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify value")
	}

	if _, err := db.conn.Exec(
		ctx, updateSpaceUserAttributeValueQuery,
		spaceUserAttributeID.PluginID, spaceUserAttributeID.Name, spaceUserAttributeID.SpaceID, spaceUserAttributeID.UserID,
		value,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to exec db")
	}

	return value, nil
}

func (db *DB) SpaceUserAttributesUpdateSpaceUserAttributeOptions(
	ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID, modifyFn modify.Fn[entry.AttributeOptions],
) (*entry.AttributeOptions, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	options, err := db.SpaceUserAttributesGetSpaceUserAttributeOptionsByID(ctx, spaceUserAttributeID)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get options by id")
	}

	options, err = modifyFn(options)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify options")
	}

	if _, err := db.conn.Exec(
		ctx, updateSpaceUserAttributeOptionsQuery,
		spaceUserAttributeID.PluginID, spaceUserAttributeID.Name, spaceUserAttributeID.SpaceID, spaceUserAttributeID.UserID,
		options,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to exec db")
	}

	return options, nil
}
