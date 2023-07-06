package harvester

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/sasha-s/go-deadlock"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

/**
Table2 features:
	- Always load data from Arbitrum from block 0 on every UC start
*/

type Table2 struct {
	mu                deadlock.RWMutex
	blockNumber       uint64
	data              map[string]map[string]*big.Int
	stakesData        map[umid.UMID]map[string]map[uint8]*big.Int
	nftData           map[umid.UMID]string
	db                *pgxpool.Pool
	adapter           Adapter
	harvesterListener func(bcName string, p []*UpdateEvent, s []*StakeEvent, n []*NftEvent) error
}

func NewTable2(db *pgxpool.Pool, adapter Adapter, listener func(bcName string, p []*UpdateEvent, s []*StakeEvent, n []*NftEvent) error) *Table2 {
	return &Table2{
		blockNumber:       0,
		data:              make(map[string]map[string]*big.Int),
		stakesData:        make(map[umid.UMID]map[string]map[uint8]*big.Int),
		nftData:           make(map[umid.UMID]string),
		adapter:           adapter,
		harvesterListener: listener,
		db:                db,
	}
}

func (t *Table2) Run() {
	t.fastForward()

	t.adapter.RegisterNewBlockListener(t.listener)
}

func (t *Table2) fastForward() {
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

	logs, err := t.adapter.GetLogs(int64(t.blockNumber)+1, int64(lastBlockNumber), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	t.ProcessLogs(lastBlockNumber, logs)
}

func (t *Table2) ProcessLogs(blockNumber uint64, logs []any) {
	fmt.Printf("Block: %d \n", blockNumber)
	events := make([]*UpdateEvent, 0)
	stakeEvents := make([]*StakeEvent, 0)
	nftEvents := make([]*NftEvent, 0)

	nftLogs := make([]*TransferNFTLog, 0)
	stakeLogs := make([]*StakeLog, 0)

	for _, log := range logs {
		switch log.(type) {
		case *TransferERC20Log:
			diff := log.(*TransferERC20Log)

			_, ok := t.data[diff.Contract]
			if !ok {
				// Table2 store everything came from adapter
				t.data[diff.Contract] = make(map[string]*big.Int)
			}

			b, ok := t.data[diff.Contract][diff.From]
			if !ok {
				t.data[diff.Contract][diff.From] = big.NewInt(0)
			}
			b = t.data[diff.Contract][diff.From]
			b.Sub(b, diff.Value)
			events = append(events, &UpdateEvent{
				Wallet:   diff.From,
				Contract: diff.Contract,
				Amount:   b,
			})

			b, ok = t.data[diff.Contract][diff.To]
			if !ok {
				t.data[diff.Contract][diff.To] = big.NewInt(0)
			}
			b = t.data[diff.Contract][diff.To]
			b.Add(b, diff.Value)
			events = append(events, &UpdateEvent{
				Wallet:   diff.To,
				Contract: diff.Contract,
				Amount:   b,
			})

		case *StakeLog:
			stake := log.(*StakeLog)

			createIfEmpty(t.stakesData, stake.OdysseyID, stake.UserWallet, stake.TokenType)

			t.stakesData[stake.OdysseyID][stake.UserWallet][stake.TokenType].Add(t.stakesData[stake.OdysseyID][stake.UserWallet][stake.TokenType], stake.AmountStaked)
			stakeEvents = append(stakeEvents, &StakeEvent{
				TxHash:    stake.TxHash,
				LogIndex:  strconv.FormatUint(uint64(stake.LogIndex), 10),
				Wallet:    stake.UserWallet,
				Kind:      stake.TokenType,
				OdysseyID: stake.OdysseyID,
				Amount:    t.stakesData[stake.OdysseyID][stake.UserWallet][stake.TokenType],
			})

			stakeLogs = append(stakeLogs, stake)

		case *UnstakeLog:
			stake := log.(*UnstakeLog)

			createIfEmpty(t.stakesData, stake.OdysseyID, stake.UserWallet, stake.TokenType)

			t.stakesData[stake.OdysseyID][stake.UserWallet][stake.TokenType].Sub(t.stakesData[stake.OdysseyID][stake.UserWallet][stake.TokenType], stake.AmountUnstaked)
			stakeEvents = append(stakeEvents, &StakeEvent{
				Wallet:    stake.UserWallet,
				OdysseyID: stake.OdysseyID,
				Amount:    t.stakesData[stake.OdysseyID][stake.UserWallet][stake.TokenType],
			})
		case *RestakeLog:
			stake := log.(*RestakeLog)

			createIfEmpty(t.stakesData, stake.FromOdysseyID, stake.UserWallet, stake.TokenType)

			stakeEvents = append(stakeEvents, &StakeEvent{
				Wallet:    stake.UserWallet,
				OdysseyID: stake.FromOdysseyID,
				Amount:    t.stakesData[stake.FromOdysseyID][stake.UserWallet][stake.TokenType],
			})

			createIfEmpty(t.stakesData, stake.ToOdysseyID, stake.UserWallet, stake.TokenType)

			t.stakesData[stake.ToOdysseyID][stake.UserWallet][stake.TokenType].Add(t.stakesData[stake.ToOdysseyID][stake.UserWallet][stake.TokenType], stake.Amount)
			stakeEvents = append(stakeEvents, &StakeEvent{
				Wallet:    stake.UserWallet,
				OdysseyID: stake.ToOdysseyID,
				Amount:    t.stakesData[stake.ToOdysseyID][stake.UserWallet][stake.TokenType],
			})
		case *TransferNFTLog:
			e := log.(*TransferNFTLog)
			t.nftData[e.TokenID] = e.To
			nftEvents = append(nftEvents, &NftEvent{
				From:      e.From,
				To:        e.To,
				OdysseyID: e.TokenID,
			})
			nftLogs = append(nftLogs, e)
		}

	}

	t.blockNumber = blockNumber

	_, name, _ := t.adapter.GetInfo()
	if err := t.harvesterListener(name, events, stakeEvents, nftEvents); err != nil {
		log.Printf("Error in harvester listener: %v\n", err)
	}

	err := t.SaveToDB(events, stakeEvents, nftLogs)
	if err != nil {
		log.Fatal(err)
	}
	//t.Display()
}

func createIfEmpty(m map[umid.UMID]map[string]map[uint8]*big.Int, odysseyID umid.UMID, wallet string, tokenType uint8) {
	_, ok := m[odysseyID]
	if !ok {
		m[odysseyID] = make(map[string]map[uint8]*big.Int)
	}
	_, ok = m[odysseyID][wallet]
	if !ok {
		m[odysseyID][wallet] = make(map[uint8]*big.Int)
	}
	if m[odysseyID][wallet][tokenType] == nil {
		m[odysseyID][wallet][tokenType] = big.NewInt(0)
	}
}

func (t *Table2) listener(blockNumber uint64, diffs []*BCDiff, stakes []*BCStake) {
	t.fastForward()
	//t.mu.Lock()
	//t.ProcessDiffs(blockNumber, diffs, stakes)
	//t.mu.Unlock()
}

func (t *Table2) SaveToDB(events []*UpdateEvent, stakeEvents []*StakeEvent, nftLogs []*TransferNFTLog) (err error) {
	wallets := make([]Address, 0)
	contracts := make([]Address, 0)
	// Save balance by value to quickly unlock mutex, otherwise have to unlock util DB transaction finished
	balances := make([]*entry.Balance, 0)
	stakes := make([]*entry.Stake, 0)

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

		var comment string
		if stake.TxHash != "" {
			comment, err = t.GetLastCommentByTxHash(stake.TxHash)
			if err != nil {
				return errors.WithMessagef(err, "failed to GetLastCommentByTxHash: %s", stake.TxHash)
			}
		}

		stakes = append(stakes, &entry.Stake{
			WalletID:     HexToAddress(stake.Wallet),
			BlockchainID: blockchainUMID,
			ObjectID:     stake.OdysseyID,
			LastComment:  comment,
			Amount:       (*entry.BigInt)(stake.Amount),
			Kind:         stake.Kind,
		})
	}

	wallets = unique(wallets)

	return t.saveToDB(wallets, contracts, balances, stakes, nftLogs)
}

func (t *Table2) saveToDB(wallets []Address, contracts []Address, balances []*entry.Balance, stakes []*entry.Stake, nftLogs []*TransferNFTLog) error {
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

	sql = `INSERT INTO stake (wallet_id, blockchain_id, object_id, amount, last_comment, updated_at, created_at, kind)
			VALUES ($1, $2, $3, $4, $5, NOW(), NOW(), $6)
			ON CONFLICT (blockchain_id, wallet_id, object_id, kind)
				DO UPDATE SET updated_at   = NOW(),
				              last_comment = $5,
							  amount       = $4
							  `

	for _, s := range stakes {

		_, err = tx.Exec(context.TODO(), sql,
			s.WalletID, blockchainUMID, s.ObjectID, s.Amount, s.LastComment, s.Kind)
		if err != nil {
			err = errors.WithMessage(err, "failed to insert stakes to DB")
			return err
		}
	}

	sql = `INSERT INTO nft (wallet_id, blockchain_id, object_id, contract_id, created_at, updated_at)	
			VALUES ($1, $2, $3, $4, NOW(), NOW())
			ON CONFLICT (wallet_id, contract_id, blockchain_id, object_id) DO UPDATE SET updated_at=NOW()`

	deleteSQL := `DELETE FROM nft WHERE object_id = $1`

	for _, nft := range nftLogs {
		_, err = tx.Exec(context.TODO(), deleteSQL, nft.TokenID)
		if err != nil {
			err = errors.WithMessage(err, "failed to delete NFT from DB")
			return err
		}

		_, err = tx.Exec(context.TODO(), sql, HexToAddress(nft.To), blockchainUMID, nft.TokenID, HexToAddress(nft.Contract))
		if err != nil {
			err = errors.WithMessage(err, "failed to insert NFT to DB")
			return err
		}
	}

	return nil
}

func (t *Table2) GetLastCommentByTxHash(txHash string) (string, error) {
	sqlQuery := `SELECT comment FROM pending_stake WHERE transaction_id = $1`

	row := t.db.QueryRow(context.TODO(), sqlQuery, HexToAddress(txHash))
	var comment string
	err := row.Scan(&comment)
	if err == pgx.ErrNoRows {
		return "", nil
	}

	if err != nil {
		return "", errors.WithMessage(err, "failed to scan 'comment' column from row")
	}

	return comment, err
}

func (t *Table2) LoadFromDB() error {
	panic("not implemented")
}

func (t *Table2) AddWalletContract(wallet string, contract string) {
	panic("not implemented")
}

func (t *Table2) Display() {
	fmt.Println("Display:")
	for token, wallets := range t.data {
		for wallet, balance := range wallets {
			fmt.Printf("%+v %+v %+v \n", token, wallet, balance.String())
		}
	}
}
