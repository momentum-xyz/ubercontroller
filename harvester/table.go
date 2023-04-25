package harvester

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/sasha-s/go-deadlock"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type Table struct {
	mu                deadlock.RWMutex
	blockNumber       uint64
	data              map[string]map[string]*big.Int
	stakesData        map[umid.UMID]map[string]*big.Int
	db                *pgxpool.Pool
	adapter           Adapter
	harvesterListener func(bcName string, p []*UpdateEvent, s []*StakeEvent)
}

func NewTable(db *pgxpool.Pool, adapter Adapter, listener func(bcName string, p []*UpdateEvent, s []*StakeEvent)) *Table {
	return &Table{
		blockNumber:       0,
		data:              make(map[string]map[string]*big.Int),
		stakesData:        make(map[umid.UMID]map[string]*big.Int),
		adapter:           adapter,
		harvesterListener: listener,
		db:                db,
	}
}

func (t *Table) Run() {
	err := t.LoadFromDB()
	if err != nil {
		fmt.Println(err)
	}

	t.fastForward()

	t.adapter.RegisterNewBlockListener(t.listener)
}

func (t *Table) fastForward() {
	t.mu.Lock()
	defer t.mu.Unlock()

	lastBlockNumber, err := t.adapter.GetLastBlockNumber()
	fmt.Printf("Fast Forward. From: %d to: %d\n", t.blockNumber, lastBlockNumber)
	if err != nil {
		fmt.Println(err)
		return
	}

	if t.blockNumber >= lastBlockNumber {
		// Table already processed latest BC block
		return
	}

	contracts := make([]common.Address, 0)

	//if t.blockNumber == 0 {
	//	// No blocks processed
	//	// Initialisation should be done using GetBalance for tokens
	//	// But for stakes we will use fastForward
	//	return
	//}

	for contract := range t.data {
		contracts = append(contracts, common.HexToAddress(contract))
	}

	fmt.Println("Doing Fast Forward")

	//if len(contracts) == 0 {
	//	return
	//}

	diffs, stakes, err := t.adapter.GetTransferLogs(int64(t.blockNumber)+1, int64(lastBlockNumber), contracts)
	if err != nil {
		fmt.Println(err)
		return
	}

	t.ProcessDiffs(lastBlockNumber, diffs, stakes)
}

func (t *Table) ProcessDiffs(blockNumber uint64, diffs []*BCDiff, stakes []*BCStake) {
	fmt.Printf("Block: %d \n", blockNumber)
	events := make([]*UpdateEvent, 0)
	stakeEvents := make([]*StakeEvent, 0)

	for _, diff := range diffs {
		_, ok := t.data[diff.Token]
		if !ok {
			// No such contract
			continue
		}
		b, ok := t.data[diff.Token][diff.From]
		if ok && b != nil { // if
			// From wallet found
			b.Sub(b, diff.Amount)
			events = append(events, &UpdateEvent{
				Wallet:   diff.From,
				Contract: diff.Token,
				Amount:   b, // TODO ask should we clone here by value
			})
		}
		b, ok = t.data[diff.Token][diff.To]
		if ok && b != nil {
			// To wallet found
			b.Add(b, diff.Amount)
			events = append(events, &UpdateEvent{
				Wallet:   diff.To,
				Contract: diff.Token,
				Amount:   b,
			})
		}
	}

	for _, stake := range stakes {
		_, ok := t.stakesData[stake.OdysseyID]
		if !ok {
			t.stakesData[stake.OdysseyID] = make(map[string]*big.Int)
		}

		t.stakesData[stake.OdysseyID][stake.From] = stake.TotalAmount
		stakeEvents = append(stakeEvents, &StakeEvent{
			Wallet:    stake.From,
			OdysseyID: stake.OdysseyID,
			Amount:    stake.TotalAmount,
		})
	}

	t.blockNumber = blockNumber

	_, name, _ := t.adapter.GetInfo()
	t.harvesterListener(name, events, stakeEvents)

	err := t.SaveToDB(events, stakeEvents)
	if err != nil {
		log.Fatal(err)
	}
	t.Display()
}

func (t *Table) listener(blockNumber uint64, diffs []*BCDiff, stakes []*BCStake) {
	t.fastForward()
	//t.mu.Lock()
	//t.ProcessDiffs(blockNumber, diffs, stakes)
	//t.mu.Unlock()
}

func (t *Table) SaveToDB(events []*UpdateEvent, stakeEvents []*StakeEvent) (err error) {
	wallets := make([]Address, 0)
	contracts := make([]Address, 0)
	// Save balance by value to quickly unlock mutex, otherwise have to unlock util DB transaction finished
	balances := make([]*entry.Balance, 0)
	stakeEntries := make([]*entry.Stake, 0)

	blockchainUMID, _, _ := t.adapter.GetInfo()

	for _, event := range events {
		if event.Amount == nil {
			continue
		}
		wallets = append(wallets, HexToAddress(event.Wallet))
		contracts = append(contracts, HexToAddress(event.Contract))
		balances = append(balances, &entry.Balance{
			WalletID:                 HexToAddress(event.Wallet),
			ContractID:               HexToAddress(event.Contract),
			BlockchainID:             blockchainUMID,
			LastProcessedBlockNumber: t.blockNumber,
			Balance:                  (*entry.BigInt)(event.Amount),
		})
	}

	for _, stake := range stakeEvents {
		wallets = append(wallets, HexToAddress(stake.Wallet))
		stakeEntries = append(stakeEntries, &entry.Stake{
			WalletID:     HexToAddress(stake.Wallet),
			BlockchainID: blockchainUMID,
			ObjectID:     stake.OdysseyID,
			LastComment:  string(0),
			Amount:       (*entry.BigInt)(stake.Amount),
		})
	}

	wallets = unique(wallets)

	fmt.Println(stakeEntries)

	return t.saveToDB(wallets, contracts, balances, stakeEntries)
}

func (t *Table) saveToDB(wallets []Address, contracts []Address, balances []*entry.Balance, stakeEntries []*entry.Stake) error {
	blockchainUMID, name, rpcURL := t.adapter.GetInfo()

	tx, err := t.db.BeginTx(context.Background(), pgx.TxOptions{})
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
			//fmt.Println("!!! Commit")
			e := tx.Commit(context.TODO())
			if e != nil {
				fmt.Println("???!!!")
				fmt.Println(e)
			}
		}
	}()

	sql := `INSERT INTO blockchain (blockchain_id, last_processed_block_number, blockchain_name, rpc_url, updated_at)
							VALUES ($1, $2, $3, $4, NOW())
							ON CONFLICT (blockchain_id) DO UPDATE SET last_processed_block_number=$2,
																	  blockchain_name=$3,
																	  rpc_url=$4,
																	  updated_at=NOW();`

	val := &entry.Blockchain{
		BlockchainID:             blockchainUMID,
		LastProcessedBlockNumber: t.blockNumber,
		BlockchainName:           name,
		RPCURL:                   rpcURL,
	}
	_, err = tx.Exec(context.Background(), sql,
		val.BlockchainID, val.LastProcessedBlockNumber, val.BlockchainName, val.RPCURL)
	if err != nil {
		return errors.WithMessage(err, "failed to insert or update blockchain DB query")
	}

	sql = `INSERT INTO wallet (wallet_id, blockchain_id)
			VALUES ($1::bytea, $2)
			ON CONFLICT (blockchain_id, wallet_id) DO NOTHING `
	for _, w := range wallets {
		_, err = tx.Exec(context.Background(), sql, w, blockchainUMID)
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

	sql = `INSERT INTO stake (wallet_id, blockchain_id, object_id, amount, last_comment, updated_at, created_at)
			VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
			ON CONFLICT (blockchain_id, wallet_id, object_id)
				DO UPDATE SET updated_at = NOW(),
							  amount     = $4`

	for _, s := range stakeEntries {
		_, err = tx.Exec(context.TODO(), sql,
			s.WalletID, blockchainUMID, s.ObjectID, s.Amount, "")
		if err != nil {
			err = errors.WithMessage(err, "failed to insert stakes to DB")
			return err
		}
	}

	return nil
}

func (t *Table) LoadFromDB() error {
	blockchainUMID, _, _ := t.adapter.GetInfo()

	tx, err := t.db.BeginTx(context.TODO(), pgx.TxOptions{})
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

	sql := `SELECT last_processed_block_number FROM blockchain WHERE blockchain_id=$1`

	row := tx.QueryRow(context.TODO(), sql, blockchainUMID)
	b := &entry.Blockchain{}

	t.mu.Lock()
	defer t.mu.Unlock()

	if err := row.Scan(&b.LastProcessedBlockNumber); err != nil {
		if err != pgx.ErrNoRows {
			return errors.WithMessage(err, "failed to scan row from blockchain table")
		}
	}

	sql = `SELECT wallet_id, contract_id, balance.balance
			FROM balance
			WHERE blockchain_id = $1`

	rows, err := tx.Query(context.TODO(), sql, blockchainUMID)
	if err != nil {
		return err
	}

	for rows.Next() {
		var wallet common.Address
		var contract common.Address
		var balance entry.BigInt

		if err := rows.Scan(&wallet, &contract, &balance); err != nil {
			return errors.WithMessage(err, "failed to scan rows from balance table")
		}

		walletStr := strings.ToLower(wallet.Hex())
		contractStr := strings.ToLower(contract.Hex())

		_, ok := t.data[contractStr]
		if !ok {
			t.data[contractStr] = make(map[string]*big.Int)
		}
		t.data[contractStr][walletStr] = (*big.Int)(&balance)

		//fmt.Println(wallet.Hex(), contract.Hex(), (*big.Int)(&balance).String())
	}

	// If DB transaction fail block will not be updated
	t.blockNumber = b.LastProcessedBlockNumber
	return nil
}

func (t *Table) AddWalletContract(wallet string, contract string) {
	wallet = strings.ToLower(wallet)
	contract = strings.ToLower(contract)

	t.mu.Lock()
	_, ok := t.data[contract]
	if !ok {
		t.data[contract] = make(map[string]*big.Int)
	}
	_, ok = t.data[contract][wallet]
	if !ok {
		t.data[contract][wallet] = nil
		// Such wallet has not existed so need to get initial balance
		go t.syncBalance(wallet, contract)
	}
	t.mu.Unlock()
}

func (t *Table) syncBalance(wallet string, contract string) {
	wallet = strings.ToLower(wallet)
	contract = strings.ToLower(contract)

	t.mu.RLock()
	blockNumber, err := t.adapter.GetLastBlockNumber()
	t.mu.RUnlock()
	if err != nil {
		err = errors.WithMessage(err, "failed to get last block number")
		fmt.Println(err)
	}
	balance, err := t.adapter.GetBalance(wallet, contract, blockNumber)
	if err != nil {
		err = errors.WithMessagef(err, "failed to get balance: %s, %s, %d", wallet, contract, blockNumber)
		fmt.Println(err)
	}
	t.mu.Lock()
	if t.blockNumber <= blockNumber {
		t.data[contract][wallet] = balance
		t.mu.Unlock()
	} else {
		t.mu.Unlock()
		t.syncBalance(wallet, contract)
	}
}

func (t *Table) Display() {
	fmt.Println("Display:")
	for token, wallets := range t.data {
		for wallet, balance := range wallets {
			fmt.Printf("%+v %+v %+v \n", token, wallet, balance.String())
		}
	}
}

func HexToAddress(s string) []byte {
	b, err := hex.DecodeString(s[2:])
	if err != nil {
		panic(err)
	}
	return b
}

func unique(slice []Address) []Address {
	keys := make(map[string]bool)
	list := []Address{}
	for _, entry := range slice {
		entryStr := hex.EncodeToString(entry)
		if _, value := keys[entryStr]; !value {
			keys[entryStr] = true
			list = append(list, entry)
		}
	}
	return list
}
