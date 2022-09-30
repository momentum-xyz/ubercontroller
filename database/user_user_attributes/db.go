package user_user_attributes

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
	getUserUserAttributesQuery                              = `SELECT * FROM user_user_attribute;`
	getUserUserAttributesQueryBySourceUserId                = `SELECT * FROM user_user_attribute where source_user_id=$1;`
	getUserUserAttributesQueryByTargetUserId                = `SELECT * FROM user_user_attribute where target_user_id=$1;`
	getUserUserAttributesQueryBySourceUserIdAndTargetUserId = `SELECT * FROM user_user_attribute where source_user_id=$1 and target_user_id=$2;`
	updateUserUserAttributeValueQuery                       = `UPDATE user_user_attribute SET value = $5 WHERE plugin_id=$1 and attribute_name = $2 and source_user_id = $3 and target_user_id = $4;`
	updateUserUserAttributeOptionsQuery                     = `UPDATE user_user_attribute SET options = $5 WHERE plugin_id=$1 and attribute_name = $2 and source_user_id = $3 and target_user_id = $4;`
	removeUserUserAttributeByNameQuery                      = `DELETE FROM user_user_attribute WHERE attribute_name = $1;`
	removeUserUserAttributesByNamesQuery                    = `DELETE FROM user_user_attribute WHERE attribute_name IN ($1);`
	removeUserUserAttributesByPluginIdQuery                 = `DELETE FROM user_user_attribute WHERE plugin_id = $1;`
	removeUserUserAttributeByPluginIdAndNameQuery           = `DELETE FROM user_user_attribute WHERE plugin_id = $1 and attribute_name =$2;`
	removeUserUserAttributesBySourceUserIdQuery             = `DELETE FROM user_user_attribute WHERE source_user_id = $1;`
	removeUserUserAttributeByNameAndSourceUserIdQuery       = `DELETE FROM user_user_attribute WHERE attribute_name = $1 and source_user_id = $2;`
	removeUserUserAttributesByNamesAndSourceUserIdQuery     = `DELETE FROM user_user_attribute WHERE attribute_name IN ($1)  and source_user_id = $2;`

	removeUserUserAttributesByTargetUserIdQuery         = `DELETE FROM user_user_attribute WHERE target_user_id = $1;`
	removeUserUserAttributeByNameAndTargetUserIdQuery   = `DELETE FROM user_user_attribute WHERE attribute_name = $1 and target_user_id = $2;`
	removeUserUserAttributesByNamesAndTargetUserIdQuery = `DELETE FROM user_user_attribute WHERE attribute_name IN ($1)  and target_user_id = $2;`

	removeUserUserAttributesBySourceUserIdAndTargetUserIdQuery         = `DELETE FROM user_user_attribute WHERE source_user_id = $1 and target_user_id = $2;`
	removeUserUserAttributeByNameAndSourceUserIdAndTargetUserIdQuery   = `DELETE FROM user_user_attribute WHERE attribute_name = $1 and source_user_id = $2 and target_user_id = $3;`
	removeUserUserAttributesByNamesAndSourceUserIdAndTargetUserIdQuery = `DELETE FROM user_user_attribute WHERE attribute_name IN ($1)  and source_user_id = $2 and target_user_id = $3;`

	removeUserUserAttributesByPluginIdAndSourceUserIdQuery        = `DELETE FROM user_user_attribute WHERE plugin_id = $1  and source_user_id = $2;`
	removeUserUserAttributesByPluginIdAndSourceUserIdAndNameQuery = `DELETE FROM user_user_attribute WHERE plugin_id = $1  and source_user_id = $2 and name = $3;`

	removeUserUserAttributesByPluginIdAndTargetUserIdQuery        = `DELETE FROM user_user_attribute WHERE plugin_id = $1  and target_user_id = $2;`
	removeUserUserAttributesByPluginIdAndTargetUserIdAndNameQuery = `DELETE FROM user_user_attribute WHERE plugin_id = $1  and target_user_id = $2 and name = $3;`

	removeUserUserAttributesByPluginIdAndSourceUserIdAndTargetUserIdQuery        = `DELETE FROM user_user_attribute WHERE plugin_id = $1  and source_user_id = $2  and target_user_id = $3;`
	removeUserUserAttributesByPluginIdAndSourceUserIdAndTargetUserIdAndNameQuery = `DELETE FROM user_user_attribute WHERE plugin_id = $1  and source_user_id = $2  and target_user_id = $3 and name = $4;`

	upsertUserUserAttributeQuery = `INSERT INTO user_user_attribute
									(plugin_id, user_user_attribute_name,source_user_id, target_user_id, value, options)
								VALUES
									($1, $2, $3, $4, $5, $6)
								ON CONFLICT (plugin_id,attribute_name, source_user_id, target_user_id)
								DO UPDATE SET
									value = $5,options = $6;`
)

var _ database.UserUserAttributesDB = (*DB)(nil)

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

func (db *DB) UserUserAttributesGetUserUserAttributes(ctx context.Context) ([]*entry.UserUserAttribute, error) {
	var assets []*entry.UserUserAttribute
	if err := pgxscan.Select(ctx, db.conn, &assets, getUserUserAttributesQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return assets, nil
}

func (db *DB) UserUserAttributesGetUserUserAttributesBySourceUserId(
	ctx context.Context, sourceUserID uuid.UUID,
) ([]*entry.UserUserAttribute, error) {
	var assets []*entry.UserUserAttribute
	if err := pgxscan.Select(
		ctx, db.conn, &assets, getUserUserAttributesQueryBySourceUserId, sourceUserID,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return assets, nil
}

func (db *DB) UserUserAttributesGetUserUserAttributesByTargetUserId(
	ctx context.Context, targetUserId uuid.UUID,
) ([]*entry.UserUserAttribute, error) {
	var assets []*entry.UserUserAttribute
	if err := pgxscan.Select(
		ctx, db.conn, &assets, getUserUserAttributesQueryByTargetUserId, targetUserId,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return assets, nil
}

func (db *DB) UserUserAttributesGetUserUserAttributesBySourceUserIdAndTargetUserId(
	ctx context.Context, sourceUserID uuid.UUID, targetUserId uuid.UUID,
) ([]*entry.UserUserAttribute, error) {
	var assets []*entry.UserUserAttribute
	if err := pgxscan.Select(
		ctx, db.conn, &assets, getUserUserAttributesQueryBySourceUserIdAndTargetUserId, sourceUserID, targetUserId,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return assets, nil
}

func (db *DB) UserUserAttributesUpsertUserUserAttribute(
	ctx context.Context, userUserAttribute *entry.UserUserAttribute,
) error {
	if _, err := db.conn.Exec(
		ctx, upsertUserUserAttributeQuery, userUserAttribute.PluginID, userUserAttribute.Name,
		userUserAttribute.SourceUserID,
		userUserAttribute.TargetUserID, userUserAttribute.Value, userUserAttribute.Options,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserUserAttributesUpsertUserUserAttributes(
	ctx context.Context, userUserAttributes []*entry.UserUserAttribute,
) error {
	batch := &pgx.Batch{}
	for _, userUserAttribute := range userUserAttributes {
		batch.Queue(
			upsertUserUserAttributeQuery, userUserAttribute.PluginID, userUserAttribute.Name,
			userUserAttribute.SourceUserID,
			userUserAttribute.TargetUserID, userUserAttribute.Value, userUserAttribute.Options,
		)
	}

	batchRes := db.conn.SendBatch(ctx, batch)
	defer batchRes.Close()

	var errs *multierror.Error
	for i := 0; i < batch.Len(); i++ {
		if _, err := batchRes.Exec(); err != nil {
			errs = multierror.Append(
				errs, errors.WithMessagef(err, "failed to exec db for: %v", userUserAttributes[i].Name),
			)
		}
	}

	return errs.ErrorOrNil()
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeByName(ctx context.Context, attributeName string) error {
	if _, err := db.conn.Exec(ctx, removeUserUserAttributeByNameQuery, attributeName); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributesByNames(ctx context.Context, attributeNames []string) error {
	if _, err := db.conn.Exec(ctx, removeUserUserAttributesByNamesQuery, attributeNames); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributesByPluginId(ctx context.Context, pluginID uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeUserUserAttributesByPluginIdQuery, pluginID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeByPluginIdAndName(
	ctx context.Context, pluginID uuid.UUID, attributeName string,
) error {
	if _, err := db.conn.Exec(
		ctx, removeUserUserAttributeByPluginIdAndNameQuery, pluginID, attributeName,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeBySourceUserId(
	ctx context.Context, sourceUserID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeUserUserAttributesBySourceUserIdQuery, sourceUserID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeByNameAndSourceUserId(
	ctx context.Context, attributeName string, sourceUserID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeUserUserAttributeByNameAndSourceUserIdQuery, attributeName, sourceUserID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeByNamesAndSourceUserId(
	ctx context.Context, attributeNames []string, sourceUserID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeUserUserAttributesByNamesAndSourceUserIdQuery, attributeNames, sourceUserID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeByTargetUserId(
	ctx context.Context, targetUserId uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeUserUserAttributesByTargetUserIdQuery, targetUserId,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeByNameAndTargetUserId(
	ctx context.Context, attributeName string, targetUserId uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeUserUserAttributeByNameAndTargetUserIdQuery, attributeName, targetUserId,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeByNamesAndTargetUserId(
	ctx context.Context, attributeNames []string, targetUserId uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeUserUserAttributesByNamesAndTargetUserIdQuery, attributeNames, targetUserId,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeBySourceUserIdAndTargetUserId(
	ctx context.Context, sourceUserID uuid.UUID, targetUserId uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeUserUserAttributesBySourceUserIdAndTargetUserIdQuery, sourceUserID, targetUserId,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeByNameAndSourceUserIdAndTargetUserId(
	ctx context.Context, attributeName string, sourceUserID uuid.UUID, targetUserId uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeUserUserAttributeByNameAndSourceUserIdAndTargetUserIdQuery, attributeName, sourceUserID,
		targetUserId,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeByNamesAndSourceUserIdAndTargetUserId(
	ctx context.Context, attributeNames []string, sourceUserID uuid.UUID, targetUserId uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeUserUserAttributesByNamesAndSourceUserIdAndTargetUserIdQuery, attributeNames, sourceUserID,
		targetUserId,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeByPluginIdAndSourceUserId(
	ctx context.Context, pluginId uuid.UUID, sourceUserID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeUserUserAttributesByPluginIdAndSourceUserIdQuery, pluginId, sourceUserID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeByPluginIdAndSourceUserIdAndName(
	ctx context.Context, pluginId uuid.UUID, attributeName string, sourceUserID uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeUserUserAttributesByPluginIdAndSourceUserIdAndNameQuery, pluginId, sourceUserID, attributeName,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeByPluginIdAndTargetUserId(
	ctx context.Context, pluginId uuid.UUID, targetUserId uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeUserUserAttributesByPluginIdAndTargetUserIdQuery, pluginId, targetUserId,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeByPluginIdAndTargetUserIdAndName(
	ctx context.Context, pluginId uuid.UUID, attributeName string, targetUserId uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeUserUserAttributesByPluginIdAndTargetUserIdAndNameQuery, pluginId, targetUserId, attributeName,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}
func (db *DB) UserUserAttributesRemoveUserUserAttributeByPluginIdAndSourceUserIdAndTargetUserId(
	ctx context.Context, pluginId uuid.UUID, sourceUserID uuid.UUID, targetUserId uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeUserUserAttributesByPluginIdAndSourceUserIdAndTargetUserIdQuery, pluginId, sourceUserID,
		targetUserId,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserUserAttributesRemoveUserUserAttributeByPluginIdAndSourceUserIdAndTargetUserIdAndName(
	ctx context.Context, pluginId uuid.UUID, attributeName string, sourceUserID uuid.UUID, targetUserId uuid.UUID,
) error {
	if _, err := db.conn.Exec(
		ctx, removeUserUserAttributesByPluginIdAndSourceUserIdAndTargetUserIdAndNameQuery, pluginId, sourceUserID,
		targetUserId,
		attributeName,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserUserAttributesUpdateUserUserAttributeOptions(
	ctx context.Context, pluginID uuid.UUID, attributeName string, sourceUserID uuid.UUID, targetUserId uuid.UUID,
	options *entry.AttributeOptions,
) error {
	if _, err := db.conn.Exec(
		ctx, updateUserUserAttributeOptionsQuery, attributeName, pluginID, sourceUserID, targetUserId, options,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UserUserAttributesUpdateUserUserAttributeValue(
	ctx context.Context, pluginID uuid.UUID, attributeName string, sourceUserID uuid.UUID, targetUserId uuid.UUID,
	value *entry.AttributeValue,
) error {
	if _, err := db.conn.Exec(
		ctx, updateUserUserAttributeValueQuery, attributeName, pluginID, sourceUserID, targetUserId, value,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}
