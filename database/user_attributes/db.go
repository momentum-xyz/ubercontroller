package user_attributes

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
	getUserAttributesQuery                                 = `SELECT * FROM user_attribute;`
	getUserAttributesQueryByUserID                         = `SELECT * FROM user_attribute WHERE user_id = $1;`
	getUserAttributeByPluginIDAndNameAndUserIDQuery        = `SELECT * FROM user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND user_id = $3;`
	getUserAttributeValueByPluginIDAndNameAndUserIDQuery   = `SELECT value FROM user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND user_id = $3;`
	getUserAttributeOptionsByPluginIDAndNameAndUserIDQuery = `SELECT options FROM user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND user_id = $3;`

	removeUserAttributeByNameQuery                      = `DELETE FROM user_attribute WHERE attribute_name = $1;`
	removeUserAttributesByNamesQuery                    = `DELETE FROM user_attribute WHERE attribute_name IN ($1);`
	removeUserAttributesByPluginIDQuery                 = `DELETE FROM user_attribute WHERE plugin_id = $1;`
	removeUserAttributeByPluginIDAndNameQuery           = `DELETE FROM user_attribute WHERE plugin_id = $1 AND attribute_name = $2;`
	removeUserAttributesByUserIDQuery                   = `DELETE FROM user_attribute WHERE user_id = $1;`
	removeUserAttributeByNameAndUserIDQuery             = `DELETE FROM user_attribute WHERE attribute_name = $1 AND user_id = $2;`
	removeUserAttributesByNamesAndUserIDQuery           = `DELETE FROM user_attribute WHERE attribute_name IN ($1)  AND user_id = $2;`
	removeUserAttributesByPluginIDAndUserIDQuery        = `DELETE FROM user_attribute WHERE plugin_id = $1 AND user_id = $2;`
	removeUserAttributesByPluginIDAndNameAndUserIDQuery = `DELETE FROM user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND user_id = $3;`

	updateUserAttributeValueQuery   = `UPDATE user_attribute SET value = $4 WHERE plugin_id = $1 AND attribute_name = $2 AND user_id = $3;`
	updateUserAttributeOptionsQuery = `UPDATE user_attribute SET options = $4 WHERE plugin_id = $1 AND attribute_name = $2 AND user_id = $3;`

	upsertUserAttributeQuery = `INSERT INTO user_attribute
									(plugin_id, attribute_name, user_id, value, options)
								VALUES
									($1, $2, $3, $4, $5)
								ON CONFLICT (plugin_id, attribute_name, user_id)
								DO UPDATE SET
									value = $4, options = $5;`
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

func (db *DB) UserAttributesGetUserAttributes(ctx context.Context) ([]*entry.UserAttribute, error) {
	var attributes []*entry.UserAttribute
	if err := pgxscan.Select(ctx, db.conn, &attributes, getUserAttributesQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return attributes, nil
}

func (db *DB) UserAttributesGetUserAttributesByUserID(
	ctx context.Context, userID uuid.UUID,
) ([]*entry.UserAttribute, error) {
	var attributes []*entry.UserAttribute
	if err := pgxscan.Select(ctx, db.conn, &attributes, getUserAttributesQueryByUserID, userID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return attributes, nil
}

func (db *DB) UserAttributesGetUserAttributeByID(
	ctx context.Context, userAttributeID entry.UserAttributeID,
) (*entry.UserAttribute, error) {
	var attribute entry.UserAttribute

	if err := pgxscan.Get(
		ctx, db.conn, &attribute, getUserAttributeByPluginIDAndNameAndUserIDQuery,
		userAttributeID.PluginID, userAttributeID.Name, userAttributeID.UserID,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}

	return &attribute, nil
}

func (db *DB) UserAttributesGetUserAttributeValueByID(
	ctx context.Context, userAttributeID entry.UserAttributeID,
) (*entry.AttributeValue, error) {
	var value entry.AttributeValue

	if err := pgxscan.Get(
		ctx, db.conn, &value, getUserAttributeValueByPluginIDAndNameAndUserIDQuery,
		userAttributeID.PluginID, userAttributeID.Name, userAttributeID.UserID,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}

	return &value, nil
}

func (db *DB) UserAttributesGetUserAttributeOptionsByID(
	ctx context.Context, userAttributeID entry.UserAttributeID,
) (*entry.AttributeOptions, error) {
	var options entry.AttributeOptions

	if err := pgxscan.Get(
		ctx, db.conn, &options, getUserAttributeOptionsByPluginIDAndNameAndUserIDQuery,
		userAttributeID.PluginID, userAttributeID.Name, userAttributeID.UserID,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}

	return &options, nil
}

// TODO: we really need to think about it
func (db *DB) UserAttributesUpsertUserAttribute(
	ctx context.Context, userAttribute *entry.UserAttribute, modifyFn modify.Fn[entry.AttributePayload],
) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	attribute, err := db.UserAttributesGetUserAttributeByID(ctx, userAttribute.UserAttributeID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return errors.WithMessage(err, "failed to query db")
		}
	}

	if err != nil {
		userAttribute.AttributePayload = modifyFn((*entry.AttributePayload)(nil))
	} else {
		userAttribute.AttributePayload = modifyFn(attribute.AttributePayload)
	}

	if _, err := db.conn.Exec(
		ctx, upsertUserAttributeQuery,
		userAttribute.PluginID, userAttribute.Name, userAttribute.UserID, userAttribute.Value, userAttribute.Options,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}

	return nil
}

func (db *DB) UserAttributesRemoveUserAttributeByName(ctx context.Context, name string) error {
	if _, err := db.conn.Exec(ctx, removeUserAttributeByNameQuery, name); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserAttributesRemoveUserAttributesByNames(ctx context.Context, names []string) error {
	if _, err := db.conn.Exec(ctx, removeUserAttributesByNamesQuery, names); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserAttributesRemoveUserAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeUserAttributesByPluginIDQuery, pluginID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserAttributesRemoveUserAttributeByAttributeID(
	ctx context.Context, attributeID entry.AttributeID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeUserAttributeByPluginIDAndNameQuery, attributeID.PluginID, attributeID.Name,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserAttributesRemoveUserAttributeByUserID(ctx context.Context, userID uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeUserAttributesByUserIDQuery, userID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserAttributesRemoveUserAttributeByNameAndUserID(
	ctx context.Context, name string, userID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeUserAttributeByNameAndUserIDQuery, name, userID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserAttributesRemoveUserAttributeByNamesAndUserID(
	ctx context.Context, names []string, userID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeUserAttributesByNamesAndUserIDQuery, names, userID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserAttributesRemoveUserAttributeByPluginIDAndUserID(
	ctx context.Context, pluginID uuid.UUID, userID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeUserAttributesByPluginIDAndUserIDQuery, pluginID, userID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserAttributesRemoveUserAttributeByID(
	ctx context.Context, userAttributeID entry.UserAttributeID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeUserAttributesByPluginIDAndNameAndUserIDQuery,
		userAttributeID.PluginID, userAttributeID.Name, userAttributeID.UserID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserAttributesUpdateUserAttributeValue(
	ctx context.Context, userAttributeID entry.UserAttributeID,
	modifyFn modify.Fn[entry.AttributeValue],
) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	value, err := db.UserAttributesGetUserAttributeValueByID(ctx, userAttributeID)
	if err != nil {
		return errors.WithMessage(err, "failed to query db")
	}

	value = modifyFn(value)

	if _, err := db.conn.Exec(
		ctx, updateUserAttributeValueQuery, userAttributeID.PluginID, userAttributeID.Name, userAttributeID.UserID, value,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}

	return nil
}

func (db *DB) UserAttributesUpdateUserAttributeOptions(
	ctx context.Context, userAttributeID entry.UserAttributeID, modifyFn modify.Fn[entry.AttributeOptions],
) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	options, err := db.UserAttributesGetUserAttributeOptionsByID(ctx, userAttributeID)
	if err != nil {
		return errors.WithMessage(err, "failed to query db")
	}

	options = modifyFn(options)

	if _, err := db.conn.Exec(
		ctx, updateUserAttributeOptionsQuery, userAttributeID.PluginID, userAttributeID.Name, userAttributeID.UserID, options,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}

	return nil
}
