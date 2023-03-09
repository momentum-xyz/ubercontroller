package harvester

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

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
}

func NewBlockchain(db *pgxpool.Pool, adapter BCAdapter, uuid uuid.UUID, name string, rpcURL string) *BlockChain {
	return &BlockChain{
		uuid:    uuid,
		name:    name,
		rpcURL:  rpcURL,
		db:      db,
		adapter: adapter,
		m:       make(map[string]map[string]*entry.Balance),
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
	walletStr := string(wallet)
	contractStr := string(contract)

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

}

func (b *BlockChain) SaveBalancesToDB() error {
	wallets := make([]Address, 0)
	contracts := make([]Address, 0)
	balances := make([]*entry.Balance, 0)

	for wallet, value := range b.m {
		wallets = append(wallets, Address(wallet))
		for contract, balance := range value {
			contracts = append(contracts, Address(contract))
			balances = append(balances, balance)
		}
	}

	fmt.Println(wallets)
	fmt.Println(contracts)
	fmt.Println(balances)

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
