package stakes

import (
	"context"
	"math/big"
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
	getJoinedStakesByWalletID = `SELECT object.object_id,
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
	getStakesByWalletID    = `SELECT * FROM stake WHERE wallet_id = $1`
	getStakesByLatestStake = `SELECT last_comment FROM stake ORDER BY created_at DESC LIMIT 1;`
	getWalletsInfoQuery    = `SELECT wallet_id, contract_id, balance, blockchain_name, updated_at
					FROM balance
							 JOIN blockchain USING (blockchain_id)
					WHERE wallet_id = ANY ($1);`
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

func (db *DB) GetStakesByWalletID(ctx context.Context, walletID string) ([]*entry.Stake, error) {
	var stakes []*entry.Stake

	if err := pgxscan.Select(ctx, db.conn, &stakes, getStakesByWalletID, walletID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return stakes, nil
}

func (db *DB) GetWalletsInfo(ctx context.Context, walletIDs [][]byte) ([]*map[string]any, error) {
	wallets := make([]*map[string]any, 0)

	rows, err := db.conn.Query(ctx, getWalletsInfoQuery, walletIDs)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var walletID common.Address
		var contractID common.Address
		var balance entry.BigInt
		var blockchainName string
		var updatedAt time.Time

		if err := rows.Scan(&walletID, &contractID, &balance, &blockchainName, &updatedAt); err != nil {
			return nil, errors.WithMessage(err, "failed to scan rows from table")
		}

		item := make(map[string]any)

		item["wallet_id"] = walletID
		item["contract_id"] = contractID
		item["balance"] = (*big.Int)(&balance).String()
		item["blockchain_name"] = blockchainName
		item["updatedAt"] = updatedAt

		wallets = append(wallets, &item)
	}

	return wallets, nil
}

func (db *DB) GetStakes(ctx context.Context, walletID []byte) ([]*map[string]any, error) {
	stakes := make([]*map[string]any, 0)

	rows, err := db.conn.Query(ctx, getJoinedStakesByWalletID, walletID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var objectID umid.UMID
		var name string
		var walletID common.Address
		var blockchainID umid.UMID
		var amount entry.BigInt
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
		item["amount"] = (*big.Int)(&amount).String()
		item["reward"] = "0"
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
