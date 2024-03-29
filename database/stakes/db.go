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
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

const (
	getJoinedStakesByWalletID = `SELECT object.object_id,
       object_attribute.value ->> 'name' AS name,
       stake.wallet_id,
       stake.blockchain_id,
       stake.amount,
       stake.last_comment,
       stake.updated_at,
       stake.kind
		FROM stake
        JOIN object USING (object_id)
    	JOIN object_attribute USING (object_id)
		WHERE attribute_name = 'name'
  		AND wallet_id = $1`
	getStakesByObjectID = `SELECT * FROM stake WHERE object_id = $1`
	getStakesByWalletID = `SELECT stake.* FROM stake
				  INNER JOIN object USING (object_id)
				  WHERE wallet_id = $1`
	getStakesWithCount     = `SELECT wallet_id, COUNT(*) AS count FROM stake GROUP BY wallet_id ORDER BY count DESC;`
	getStakesByLatestStake = `SELECT last_comment FROM stake ORDER BY created_at DESC LIMIT 1;`
	getWalletsInfoQuery    = `SELECT wallet_id, contract_id, balance, blockchain_name, updated_at
					FROM balance
							 JOIN blockchain USING (blockchain_id)
					WHERE wallet_id = ANY ($1);`
	insertIntoPendingStakes = `INSERT INTO pending_stake (transaction_id,
									   object_id,
									   wallet_id,
									   blockchain_id,
									   amount,
									   comment,
									   kind,
									   updated_at,
									   created_at)
								VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())`
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

	if err := pgxscan.Select(ctx, db.conn, &stakes, getStakesByWalletID, utils.HexToAddress(walletID)); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return stakes, nil
}

func (db *DB) GetStakesWithCount(ctx context.Context) ([]*entry.Stake, error) {
	var stakes []*entry.Stake

	if err := pgxscan.Select(ctx, db.conn, &stakes, getStakesWithCount); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return stakes, nil
}

func (db *DB) GetWalletsInfo(ctx context.Context, walletIDs [][]byte) ([]*dto.WalletInfo, error) {
	wallets := make([]*dto.WalletInfo, 0)

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

		item := dto.WalletInfo{
			WalletID:       walletID.Hex(),
			ContractID:     contractID.Hex(),
			Balance:        (*big.Int)(&balance).String(),
			BlockchainName: blockchainName,
			Reward:         "0",
			Transferable:   "0",
			Staked:         "0",
			Unbonding:      "0",
			UpdatedAt:      updatedAt,
		}

		wallets = append(wallets, &item)
	}

	return wallets, nil
}

func (db *DB) GetStakes(ctx context.Context, walletID []byte) ([]*dto.Stake, error) {
	stakes := make([]*dto.Stake, 0)

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
		var kind int

		if err := rows.Scan(&objectID, &name, &walletID, &blockchainID, &amount, &lastComment, &updatedAt, &kind); err != nil {
			return nil, errors.WithMessage(err, "failed to scan rows from table")
		}

		item := dto.Stake{
			ObjectID:     objectID,
			Name:         name,
			WalletID:     walletID.Hex(),
			BlockchainID: blockchainID,
			Amount:       (*big.Int)(&amount).String(),
			Reward:       "0",
			LastComment:  lastComment,
			UpdatedAt:    updatedAt,
			Kind:         kind,
		}

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
	var comment []*string

	if err := pgxscan.Select(ctx, db.conn, &comment, getStakesByLatestStake); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	if len(comment) > 0 {
		return comment[0], nil
	}
	return nil, nil
}

func (db *DB) InsertIntoPendingStakes(ctx context.Context, transactionID []byte,
	objectID umid.UMID,
	walletID []byte,
	blockchainID umid.UMID,
	amount *big.Int,
	comment string,
	kind uint8) error {

	a := (*entry.BigInt)(amount)

	if _, err := db.conn.Exec(
		ctx, insertIntoPendingStakes,
		transactionID, objectID, walletID, blockchainID, a, comment, kind,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}
