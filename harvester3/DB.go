package harvester3

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type DB struct {
	updates        chan any
	db             *pgxpool.Pool
	blockchainID   umid.UMID
	blockchainName string
}

func NewDB(updates chan any, db *pgxpool.Pool, blockchainID umid.UMID, blockchainName string) *DB {
	return &DB{
		updates:        updates,
		db:             db,
		blockchainID:   blockchainID,
		blockchainName: blockchainName,
	}
}

func (db *DB) Run() {
	go db.worker()
}

func (db *DB) worker() {
	fmt.Println("DB Worker")
	queue := make([]InsertOrUpdateToDB, 0)
	nftQueue := make([]any, 0)
	for {
		select {
		case update := <-db.updates:
			switch u := update.(type) {
			case InsertOrUpdateToDB:
				fmt.Println("InsertOrUpdateToDB")
				queue = append(queue, u)
			case FlushToDB:
				//fmt.Println("FlushToDB")
				if len(queue) == 0 {
					continue
				}
				err := db.flush(queue, u.block)
				if err != nil {
					fmt.Println(err)
				}

				queue = make([]InsertOrUpdateToDB, 0)
			case UpsertNFTToDB:
				nftQueue = append(nftQueue, u)
			case FlushNFTToDB:
				if len(nftQueue) == 0 {
					continue
				}
				err := db.flushNFT(nftQueue, u.block)
				if err != nil {
					fmt.Println(err)
				}
			}
		}

	}
}

func (db *DB) flushNFT(queue []any, block uint64) error {
	tx, err := db.db.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return errors.WithMessage(err, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			fmt.Println("!!! Rollback")
			e := tx.Rollback(context.TODO())
			if e != nil {
				fmt.Println("???")
				fmt.Println(e)
			}
		} else {
			e := tx.Commit(context.TODO())
			if e != nil {
				fmt.Println("???!!!")
				fmt.Println(e)
			}
		}
	}()

	sql := `INSERT INTO harvester_blockchain (blockchain_id,
											  blockchain_name,
											  last_processed_block_for_tokens,
											  last_processed_block_for_nfts,
											  last_processed_block_for_ethers,
											  updated_at)
			VALUES ($1, $2, $3, $4, $5, NOW())
			ON CONFLICT (blockchain_id) DO UPDATE SET last_processed_block_for_nfts=$4,
													  updated_at                     = NOW()`

	_, err = tx.Exec(context.Background(), sql,
		db.blockchainID, db.blockchainName, 0, block, 0)
	if err != nil {
		return errors.WithMessage(err, "failed to insert or update blockchain DB query")
	}

	upsertSQL := `INSERT INTO harvester_nfts (wallet_id, contract_id, blockchain_id, item_id, updated_at)
				 VALUES ($1, $2, $3, $4, NOW())
				 ON CONFLICT (blockchain_id, contract_id, wallet_id, item_id)
				 	DO UPDATE SET updated_at = NOW();`

	for _, i := range queue {
		if item, ok := i.(UpsertNFTToDB); ok {
			_, err = tx.Exec(context.TODO(), upsertSQL,
				item.wallet, item.contract, db.blockchainID, item.id.Big().String())
			if err != nil {
				err = errors.WithMessage(err, "failed to insert balance to DB")
				return err
			}
		}
	}

	return nil
}

func (db *DB) flush(queue []InsertOrUpdateToDB, block uint64) error {

	tx, err := db.db.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return errors.WithMessage(err, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			fmt.Println("!!! Rollback")
			e := tx.Rollback(context.TODO())
			if e != nil {
				fmt.Println("???")
				fmt.Println(e)
			}
		} else {
			e := tx.Commit(context.TODO())
			if e != nil {
				fmt.Println("???!!!")
				fmt.Println(e)
			}
		}
	}()

	sql := `INSERT INTO harvester_blockchain (blockchain_id,
											  blockchain_name,
											  last_processed_block_for_tokens,
											  last_processed_block_for_nfts,
											  last_processed_block_for_ethers,
											  updated_at)
			VALUES ($1, $2, $3, $4, $5, NOW())
			ON CONFLICT (blockchain_id) DO UPDATE SET last_processed_block_for_tokens=$3,
													  updated_at                     = NOW()`

	_, err = tx.Exec(context.Background(), sql,
		db.blockchainID, db.blockchainName, block, 0, 0)
	if err != nil {
		return errors.WithMessage(err, "failed to insert or update blockchain DB query")
	}

	sql = `INSERT INTO harvester_tokens (wallet_id, contract_id, blockchain_id, balance, updated_at)
			VALUES ($1, $2, $3, $4, NOW())
			ON CONFLICT (blockchain_id, contract_id, wallet_id) DO UPDATE SET balance   =$4,
																			  updated_at=NOW();`
	for _, b := range queue {
		_, err = tx.Exec(context.TODO(), sql,
			b.wallet, b.contract, db.blockchainID, (*entry.BigInt)(b.value))
		if err != nil {
			err = errors.WithMessage(err, "failed to insert balance to DB")
			return err
		}
	}

	return nil
}

func (db *DB) loadNFTsFromDB() ([]TokenCell, error) {
	return nil, nil
}

func (db *DB) loadTokensFromDB() ([]TokenCell, error) {
	sql := `select harvester_tokens.contract_id,
				   harvester_tokens.wallet_id,
				   harvester_tokens.balance,
				   harvester_blockchain.last_processed_block_for_tokens
			from harvester_tokens
					 join harvester_blockchain using (blockchain_id)
			where blockchain_id = $1`

	rows, err := db.db.Query(context.Background(), sql, db.blockchainID)
	if err != nil {
		return nil, err
	}

	cells := make([]TokenCell, 0)

	for rows.Next() {
		var contractID common.Address
		var walletID common.Address
		var balance entry.BigInt
		var block uint64

		if err := rows.Scan(&contractID, &walletID, &balance, &block); err != nil {
			return nil, errors.WithMessage(err, "failed to scan rows from table")
		}

		cells = append(cells, TokenCell{
			Contract: contractID,
			Wallet:   walletID,
			Value:    (*big.Int)(&balance),
			Block:    block,
		})
	}

	return cells, nil
}
