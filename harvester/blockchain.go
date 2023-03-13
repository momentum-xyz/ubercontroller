package harvester

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/sasha-s/go-deadlock"

	"github.com/momentum-xyz/ubercontroller/types/entry"
)

const getBlockchainByID = `SELECT * FROM blockchain WHERE blockchain_id=$1`
const insertOrUpdate = `INSERT INTO blockchain (blockchain_id, last_processed_block_number, blockchain_name, rpc_url, updated_at)
VALUES ($1, $2, $3, $4, NOW())
ON CONFLICT (blockchain_id) DO UPDATE SET last_processed_block_number=$2,
                                          blockchain_name=$3,
                                          rpc_url=$4,
                                          updated_at=NOW();`

type BlockChain struct {
	uuid                     uuid.UUID
	name                     string
	lastProcessedBlockNumber uint64
	rpcURL                   string
	db                       *pgxpool.Pool
	adapter                  BCAdapter
	m                        map[string]map[string]*entry.Balance
	mu                       deadlock.RWMutex
	onUpdateWalletContract   UpdateWalletContractHook
}

type UpdateWalletContractHook func(bcType string, wallet string, contract string, blockNumber uint64, balance *big.Int)

func NewBlockchain(db *pgxpool.Pool, adapter BCAdapter, uuid uuid.UUID, name string, rpcURL string, onUpdateWalletContract UpdateWalletContractHook) *BlockChain {
	return &BlockChain{
		uuid:                   uuid,
		name:                   name,
		rpcURL:                 rpcURL,
		db:                     db,
		adapter:                adapter,
		m:                      make(map[string]map[string]*entry.Balance),
		mu:                     deadlock.RWMutex{},
		onUpdateWalletContract: onUpdateWalletContract,
	}
}

func (b *BlockChain) ToEntry() *entry.Blockchain {
	return &entry.Blockchain{
		BlockchainID:             b.uuid,
		LastProcessedBlockNumber: b.lastProcessedBlockNumber,
		BlockchainName:           b.name,
		RPCURL:                   b.rpcURL,
	}
}

func (b *BlockChain) SubscribeForWalletAndContract(wallet []byte, contract []byte) {
	walletStr := hexutil.Encode(wallet)
	contractStr := hexutil.Encode(contract)

	b.mu.Lock()
	defer b.mu.Unlock()

	_, ok := b.m[walletStr]
	if !ok {
		b.m[walletStr] = make(map[string]*entry.Balance)
	}

	_, ok = b.m[walletStr][contractStr]
	if ok {
		// Already subscribed
		fmt.Printf("Already subscribed wallet:%s, contract:%s \n", walletStr, contractStr)
		return
	}

	b.m[walletStr][contractStr] = &entry.Balance{
		WalletID:                 wallet,
		ContractID:               contract,
		BlockchainID:             b.uuid,
		LastProcessedBlockNumber: 0,
		Balance:                  0,
	}

	b.getBalanceFromBC(walletStr, contractStr)
}

func (b *BlockChain) update(walletStr string, contractStr string, blockNumber uint64, balance *big.Int) {
	b.mu.Lock()
	defer b.mu.Unlock()

}

func (b *BlockChain) getBalanceFromBC(walletStr string, contractStr string) (*big.Int, uint64, error) {
	n, err := b.adapter.GetLastBlockNumber()
	if err != nil {
		return nil, 0, err
	}
	balance, err := b.adapter.GetBalance(walletStr, contractStr, n)
	return balance, n, nil
}

func (b *BlockChain) SaveBalancesToDB() (err error) {
	wallets := make([]Address, 0)
	contracts := make([]Address, 0)
	// Save balance by value to quickly unlock mutex, otherwise have to unlock util DB transaction finished
	balances := make([]entry.Balance, 0)

	b.mu.RLock()
	for wallet, value := range b.m {
		wallets = append(wallets, HexToAddress(wallet))
		for contract, balance := range value {
			contracts = append(contracts, HexToAddress(contract))
			balances = append(balances, *balance)
		}
	}
	b.mu.RUnlock()

	fmt.Println(wallets)
	fmt.Println(contracts)
	fmt.Println(balances)

	tx, err := b.db.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return errors.WithMessage(err, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.TODO())
		} else {
			tx.Commit(context.TODO())
		}
	}()

	sql := `INSERT INTO wallet (wallet_id, blockchain_id)
			VALUES ($1::bytea, $2)
			ON CONFLICT (blockchain_id, wallet_id) DO NOTHING `
	for _, w := range wallets {
		_, err = tx.Exec(context.Background(), sql, w, b.uuid)
		if err != nil {
			err = errors.WithMessage(err, "failed to insert wallet to DB")
			return err
		}
	}

	sql = `INSERT INTO contract (contract_id, name)
			VALUES ($1, $2)
			ON CONFLICT (contract_id) DO NOTHING`
	for _, c := range contracts {
		_, err = tx.Exec(context.TODO(), sql, c, "")
		if err != nil {
			err = errors.WithMessage(err, "failed to insert contract to DB")
			return err
		}
	}

	sql = `INSERT INTO balance (wallet_id, contract_id, blockchain_id, balance, last_processed_block_number)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (wallet_id, contract_id, blockchain_id)
				DO UPDATE SET balance                     = $4,
							  last_processed_block_number = $5`

	for _, b := range balances {
		_, err = tx.Exec(context.TODO(), sql,
			b.WalletID, b.ContractID, b.BlockchainID, b.Balance, b.LastProcessedBlockNumber)
		if err != nil {
			err = errors.WithMessage(err, "failed to insert balance to DB")
			return err
		}
	}

	return nil
}

func (b *BlockChain) LoadFromDB() error {
	var vals []entry.Blockchain
	if err := pgxscan.Select(context.Background(), b.db, &vals, getBlockchainByID, b.uuid); err != nil {

		return errors.WithMessage(err, "failed to select from db")
	}

	fmt.Println(vals)
	if len(vals) != 0 {
		e := vals[0]
		b.name = e.BlockchainName
		b.rpcURL = e.RPCURL
		b.lastProcessedBlockNumber = e.LastProcessedBlockNumber
	} else {
		return b.SaveToDB()
	}

	return nil
}

func (b *BlockChain) SaveToDB() error {
	val := b.ToEntry()
	_, err := b.db.Exec(context.Background(), insertOrUpdate,
		val.BlockchainID, val.LastProcessedBlockNumber, val.BlockchainName, val.RPCURL)
	if err != nil {
		return errors.WithMessage(err, "failed to exec DB query")
	}

	return nil
}
