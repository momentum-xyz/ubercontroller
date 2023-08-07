package harvester3

import (
	"context"
	"fmt"

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
	for {
		select {
		case update := <-db.updates:
			switch u := update.(type) {
			case InsertOrUpdateToDB:
				fmt.Println("InsertOrUpdateToDB")
				queue = append(queue, u)
			case FlushToDB:
				fmt.Println("FlushToDB")
				if len(queue) == 0 {
					continue
				}
				err := db.flush(queue, u.block)
				if err != nil {
					fmt.Println(err)
				}

				queue = make([]InsertOrUpdateToDB, 0)
			}
		}

	}
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
