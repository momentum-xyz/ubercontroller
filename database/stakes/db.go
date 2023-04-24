package stakes

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

const (
	getStakesByWalletID = `SELECT object.object_id,
       object_attribute.value ->> 'name' AS name,
       stake.wallet_id,
       stake.blockchain_id,
       stake.amount,
       stake.last_comment,
       stake.updated_at
		FROM stake
        JOIN object USING (object_id)
    	JOIN object_attribute USING (object_id)
		WHERE attribute_name = 'name'
  		AND wallet_id = $1`
	getStakesByObjectID    = `SELECT * FROM stake WHERE object_id = $1`
	getStakesByLatestStake = `SELECT last_comment FROM stake ORDER BY created_at DESC LIMIT 1;`
)

var _ database.StakesDB = (*DB)(nil)

type DB struct {
	conn *pgxpool.Pool
}

func NewDB(conn *pgxpool.Pool) *DB {
	return &DB{
		conn: conn,
	}
}

func (db *DB) GetStakesByWalletID(ctx context.Context, walletID []byte) ([]*map[string]any, error) {
	stakes := make([]*map[string]any, 0)

	rows, err := db.conn.Query(ctx, getStakesByWalletID, walletID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var objectID umid.UMID
		var name string
		var walletID common.Address
		var blockchainID umid.UMID
		//var amount entry.BigInt // TODO use BigInt need to update SQL schema
		var amount int64
		var lastComment string
		var updatedAt time.Time

		if err := rows.Scan(&objectID, &name, &walletID, &blockchainID, &amount, &lastComment, &updatedAt); err != nil {
			return nil, errors.WithMessage(err, "failed to scan rows from table")
		}

		item := make(map[string]any)

		item["object_id"] = objectID
		item["name"] = name
		item["wallet_id"] = walletID
		item["blockchain_id"] = blockchainID
		item["amount"] = amount
		item["lastComment"] = lastComment
		item["updatedAt"] = updatedAt

		stakes = append(stakes, &item)
	}

	return stakes, nil
}

func (db *DB) GetStakesByWorldID(ctx context.Context, worldID umid.UMID) ([]*entry.Stake, error) {
	var stakes []*entry.Stake

	if err := pgxscan.Select(ctx, db.conn, &stakes, getStakesByObjectID, worldID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return stakes, nil
}

func (db *DB) GetStakeByLatestStake(ctx context.Context) (*string, error) {
	var stake *string

	if err := pgxscan.Get(ctx, db.conn, &stake, getStakesByLatestStake); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return stake, nil
}
