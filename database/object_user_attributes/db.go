package object_user_attributes

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

const (
	getObjectUserAttributesQuery           = `SELECT * FROM object_user_attribute;`
	getObjectUserAttributeByIDQuery        = `SELECT * FROM object_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND object_id = $3 AND user_id = $4;`
	getObjectUserAttributePayloadByIDQuery = `SELECT value, options FROM object_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND object_id = $3 AND user_id = $4;`
	getObjectUserAttributeValueByIDQuery   = `SELECT value FROM object_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND object_id = $3 AND user_id = $4;`
	getObjectUserAttributeOptionsByIDQuery = `
		SELECT COALESCE(a_type.options, '{}'::jsonb) || COALESCE(oua.options, '{}'::jsonb) as options
		FROM attribute_type as a_type
		LEFT JOIN object_user_attribute as oua ON
  			oua.plugin_id=a_type.plugin_id AND
  			oua.attribute_name=a_type.attribute_name AND
  			oua.object_id = $3 AND oua.user_id = $4
  		WHERE a_type.plugin_id = $1 AND a_type.attribute_name = $2
		;`
	getObjectUserAttributesByObjectIDQuery          = `SELECT * FROM object_user_attribute WHERE object_id = $1;`
	getObjectUserAttributesByUserIDQuery            = `SELECT * FROM object_user_attribute WHERE user_id = $1;`
	getObjectUserAttributesByObjectIDAndUserIDQuery = `SELECT * FROM object_user_attribute WHERE object_id = $1 AND user_id = $2;`
	getObjectUserAttributesByObjectAttributeIDQuery = `SELECT * FROM object_user_attribute WHERE plugin_id = $1 AND attribute_name = $2 AND object_id = $3;`

	getObjectUserAttributesCountQuery                       = `SELECT COUNT(*) FROM object_user_attribute WHERE value IS NOT NULL;`
	getObjectUserAttributesCountByObjectIDQuery             = `SELECT COUNT(*) FROM object_user_attribute WHERE value IS NOT NULL AND object_id = $1 AND attribute_name = $2;`
	getObjectUserAttributesCountByObjectIDAndUpdatedAtQuery = `SELECT COUNT(*) FROM object_user_attribute WHERE value IS NOT NULL AND object_id = $1 AND attribute_name = $2 AND updated_at >= $3;`

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
		objectUserAttributeID.ObjectID, objectUserAttributeID.UserID,
	).
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
		objectUserAttributeID.ObjectID, objectUserAttributeID.UserID,
	).
		Scan(&options); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &options, nil
}

func (db *DB) GetObjectUserAttributesByObjectID(
	ctx context.Context, objectID umid.UMID,
) ([]*entry.ObjectUserAttribute, error) {
	var attributes []*entry.ObjectUserAttribute
	if err := pgxscan.Select(ctx, db.conn, &attributes, getObjectUserAttributesByObjectIDQuery, objectID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return attributes, nil
}

func (db *DB) GetObjectUserAttributesByUserID(
	ctx context.Context, userID umid.UMID,
) ([]*entry.ObjectUserAttribute, error) {
	var attributes []*entry.ObjectUserAttribute
	if err := pgxscan.Select(ctx, db.conn, &attributes, getObjectUserAttributesByUserIDQuery, userID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return attributes, nil
}

func (db *DB) GetObjectUserAttributesByObjectIDAndUserID(
	ctx context.Context, objectID umid.UMID, userID umid.UMID,
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

func (db *DB) GetObjectUserAttributesCountByObjectID(ctx context.Context, objectID umid.UMID, attributeName string, sinceTime *time.Time) (uint64, error) {
	var count uint64
	if sinceTime != nil {
		if err := db.conn.QueryRow(ctx, getObjectUserAttributesCountByObjectIDAndUpdatedAtQuery, objectID, attributeName, *sinceTime).
			Scan(&count); err != nil {
			return 0, errors.WithMessage(err, "failed to query db")
		}
	} else {
		if err := db.conn.QueryRow(ctx, getObjectUserAttributesCountByObjectIDQuery, objectID, attributeName).
			Scan(&count); err != nil {
			return 0, errors.WithMessage(err, "failed to query db")
		}
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
			return nil, errors.WithMessage(err, "failed to get attribute payload by umid")
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

func (db *DB) RemoveObjectUserAttributesByPluginID(ctx context.Context, pluginID umid.UMID) error {
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
	ctx context.Context, objectID umid.UMID,
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
	ctx context.Context, name string, objectID umid.UMID,
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
	ctx context.Context, names []string, objectID umid.UMID,
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
	ctx context.Context, userID umid.UMID,
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
	ctx context.Context, name string, userID umid.UMID,
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
	ctx context.Context, names []string, userID umid.UMID,
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
	ctx context.Context, objectID umid.UMID, userID umid.UMID,
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
	ctx context.Context, name string, objectID umid.UMID, userID umid.UMID,
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
	ctx context.Context, names []string, objectID umid.UMID, userID umid.UMID,
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
	ctx context.Context, pluginID umid.UMID, objectID umid.UMID,
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
	ctx context.Context, pluginID umid.UMID, userID umid.UMID,
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
	ctx context.Context, pluginID umid.UMID, objectID umid.UMID, userID umid.UMID,
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
		return nil, errors.WithMessage(err, "failed to get value by umid")
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
		return nil, errors.WithMessage(err, "failed to get options by umid")
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

// ValueEntries returns a combined, sorted, list of nested items inside a attribute.
// fields is the list of field name returned for each nested items, should be simple values (are seen as strings).
//
// The JSON value is expectd to be an object with key->item format.
// The key being some unique identifier and the item some nested JSON.
//
//	{
//	  "id-1": {"title": "bar", "some": {"thing": 1}},
//	  "id-2": {"title": "qux", "some": {"other": "thing"}}
//	}
//
// The result:
//
// [{"_key": "id-1", "_user": "{...}", "title": "bar"},{"_key": "id-2", "_user": "{...}", title": "qux"}]
func (db *DB) ValueEntries(
	ctx context.Context,
	attrID entry.ObjectAttributeID,
	fields []string,
	order string,
	descending bool,
	limit uint,
	offset uint,
) ([]map[string]interface{}, error) {
	// TODO: pick a sql query builder library? if we are gonna do this more often,
	// string concat starts to suck.
	sql := `SELECT entry.key as "_key",
	        jsonb_build_object(
			  'user_id', u.user_id,
	          'profile', u.profile
            ) as "_user"`
	if len(fields) > 0 {
		sql += ", item.*"
	}
	sql += `
	FROM object_user_attribute as attr
		 INNER JOIN "user" AS u USING (user_id),
         jsonb_each(attr.value) as entry`
	if len(fields) > 0 {
		sqlRecordTypes := []string{}
		for _, field := range fields {
			ident := pgx.Identifier([]string{field}).Sanitize()
			sqlRecordTypes = append(sqlRecordTypes, ident+" text") // TODO: configurables 'types' instead of text
		}
		itemType := fmt.Sprintf("item(%s)", strings.Join(sqlRecordTypes, ","))
		sql += `,
		jsonb_to_record(entry.value) AS ` + itemType
	}
	sql += `
    WHERE
      attr.plugin_id = $1 
	  AND attr.attribute_name = $2
	  AND attr.object_id = $3
	`
	if order != "" {
		// TODO: handle the case when no fields gives (value->'foo' should be used)
		// But that requires dynamic query params and we don't have named params support.
		orderBy := pgx.Identifier([]string{"item", order}).Sanitize()
		sql += "ORDER BY " + orderBy
		if descending {
			sql += "DESC\n"
		} else {
			sql += "\n"
		}
	}
	if limit > 0 {
		sql += fmt.Sprintf("LIMIT %d ", limit)
	}
	if offset > 0 {
		sql += fmt.Sprintf("OFFSET %d", offset)
	}
	sql += ";"
	var itemList []map[string]interface{}
	qArgs := []interface{}{attrID.PluginID, attrID.Name, attrID.ObjectID}
	if err := pgxscan.Select(ctx, db.conn, &itemList, sql, qArgs...); err != nil {
		return nil, err
	}
	return itemList, nil

}

// ValueEntriesCount return the number of nested entries in the JSON value.
// This assumes an JSON object in the value field.
// Each entry (key) in this object is counted.
func (db *DB) ValueEntriesCount(
	ctx context.Context,
	objectAttributeID entry.ObjectAttributeID,
) (uint, error) {
	var count uint
	sql := `SELECT count(*) FROM (
		SELECT jsonb_object_keys(value) FROM object_user_attribute
			WHERE plugin_id = $1 AND attribute_name = $2 AND object_id = $3
	) oua;`
	if err := db.conn.QueryRow(ctx, sql,
		objectAttributeID.PluginID, objectAttributeID.Name, objectAttributeID.ObjectID,
	).Scan(&count); err != nil {
		return 0, errors.WithMessage(err, "failed to query db")
	}
	return count, nil
}
