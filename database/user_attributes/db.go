package user_attributes

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
	getUserAttributesQuery                              = `SELECT * FROM user_attribute;`
	getSpaceAttributesQueryByUserId                     = `SELECT * FROM user_attribute WHERE user_id = $1;`
	updateUserAttributeValueQuery                       = `UPDATE user_attribute SET value = $4 WHERE plugin_id=$1 and attribute_name = $2 and user_id = $3;`
	updateUserAttributeOptionsQuery                     = `UPDATE user_attribute SET options = $4 WHERE plugin_id=$1 and attribute_name = $2 and user_id = $3;`
	removeUserAttributeByNameQuery                      = `DELETE FROM user_attribute WHERE attribute_name = $1;`
	removeUserAttributesByNamesQuery                    = `DELETE FROM user_attribute WHERE attribute_name IN ($1);`
	removeUserAttributesByPluginIdQuery                 = `DELETE FROM user_attribute WHERE plugin_id = $1;`
	removeUserAttributeByPluginIdAndNameQuery           = `DELETE FROM user_attribute WHERE plugin_id = $1 and attribute_name =$2;`
	removeUserAttributesByUserIdQuery                   = `DELETE FROM user_attribute WHERE user_id = $1;`
	removeUserAttributeByNameAndUserIdQuery             = `DELETE FROM user_attribute WHERE attribute_name = $1 and user_id = $2;`
	removeUserAttributesByNamesAndUserIdQuery           = `DELETE FROM user_attribute WHERE attribute_name IN ($1)  and user_id = $2;`
	removeUserAttributesByPluginIdAndUserIdQuery        = `DELETE FROM user_attribute WHERE plugin_id = $1  and user_id = $2;`
	removeUserAttributesByPluginIdAndUserIdAndNameQuery = `DELETE FROM user_attribute WHERE plugin_id = $1  and user_id = $2 and name = $3;`

	upsertUserAttributeQuery = `INSERT INTO user_attribute
									(plugin_id, user_attribute_name,user_id, value, options)
								VALUES
									($1, $2, $3, $4, $5)
								ON CONFLICT (plugin_id,attribute_name, user_id)
								DO UPDATE SET
									value = $4,options = $5;`
)

var _ database.UserAttributesDB = (*DB)(nil)

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

func (db *DB) UserAttributesGetUserAttributes(ctx context.Context) ([]*entry.UserAttribute, error) {
	var assets []*entry.UserAttribute
	if err := pgxscan.Select(ctx, db.conn, &assets, getUserAttributesQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return assets, nil
}

func (db *DB) UserAttributesGetUserAttributesByUserId(
	ctx context.Context, userId uuid.UUID,
) ([]*entry.SpaceAttribute, error) {
	var assets []*entry.SpaceAttribute
	if err := pgxscan.Select(ctx, db.conn, &assets, getSpaceAttributesQueryByUserId, userId); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return assets, nil
}

func (db *DB) UserAttributesUpsertUserAttribute(ctx context.Context, userAttribute *entry.UserAttribute) error {
	if _, err := db.conn.Exec(
		ctx, upsertUserAttributeQuery, userAttribute.PluginID, userAttribute.Name, userAttribute.UserID,
		userAttribute.Value, userAttribute.Options,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserAttributesUpsertUserAttributes(
	ctx context.Context, userAttributes []*entry.UserAttribute,
) error {
	batch := &pgx.Batch{}
	for _, userAttribute := range userAttributes {
		batch.Queue(
			upsertUserAttributeQuery, userAttribute.PluginID, userAttribute.Name, userAttribute.UserID,
			userAttribute.Value, userAttribute.Options,
		)
	}

	batchRes := db.conn.SendBatch(ctx, batch)
	defer batchRes.Close()

	var errs *multierror.Error
	for i := 0; i < batch.Len(); i++ {
		if _, err := batchRes.Exec(); err != nil {
			errs = multierror.Append(
				errs, errors.WithMessagef(err, "failed to exec db for: %v", userAttributes[i].Name),
			)
		}
	}

	return errs.ErrorOrNil()
}

func (db *DB) UserAttributesRemoveUserAttributeByName(ctx context.Context, attributeName string) error {
	if _, err := db.conn.Exec(ctx, removeUserAttributeByNameQuery, attributeName); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserAttributesRemoveUserAttributesByNames(ctx context.Context, attributeNames []string) error {
	if _, err := db.conn.Exec(ctx, removeUserAttributesByNamesQuery, attributeNames); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserAttributesRemoveUserAttributesByPluginId(ctx context.Context, pluginID uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeUserAttributesByPluginIdQuery, pluginID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserAttributesRemoveUserAttributeByPluginIdAndName(
	ctx context.Context, pluginID uuid.UUID, attributeName string,
) error {
	if _, err := db.conn.Exec(
		ctx, removeUserAttributeByPluginIdAndNameQuery, pluginID, attributeName,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserAttributesRemoveUserAttributeByUserId(
	ctx context.Context, userID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeUserAttributesByUserIdQuery, userID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserAttributesRemoveUserAttributeByNameAndUserId(
	ctx context.Context, attributeName string, userID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeUserAttributeByNameAndUserIdQuery, attributeName, userID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserAttributesRemoveUserAttributeByNamesAndUserId(
	ctx context.Context, attributeNames []string, userID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeUserAttributesByNamesAndUserIdQuery, attributeNames, userID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserAttributesRemoveUserAttributeByPluginIdAndUserId(
	ctx context.Context, pluginId uuid.UUID, userID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeUserAttributesByPluginIdAndUserIdQuery, pluginId, userID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserAttributesRemoveUserAttributeByPluginIdAndUserIdAndName(
	ctx context.Context, pluginId uuid.UUID, attributeName string, userID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeUserAttributesByPluginIdAndUserIdAndNameQuery, pluginId, userID, attributeName,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserAttributesUpdateUserAttributeOptions(
	ctx context.Context, pluginID uuid.UUID, attributeName string, userId uuid.UUID, options *entry.AttributeOptions,
) error {
	if _, err := db.conn.Exec(
		ctx, updateUserAttributeOptionsQuery, attributeName, pluginID, userId, options,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserAttributesUpdateUserAttributeValue(
	ctx context.Context, pluginID uuid.UUID, attributeName string, userId uuid.UUID, value *entry.AttributeValue,
) error {
	if _, err := db.conn.Exec(
		ctx, updateUserAttributeValueQuery, attributeName, pluginID, userId, value,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}
