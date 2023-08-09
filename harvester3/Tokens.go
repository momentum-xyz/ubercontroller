package harvester3

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sasha-s/go-deadlock"
	"go.uber.org/zap"
)

type Tokens struct {
	updates        chan any
	updatesDB      chan any
	output         chan UpdateCell
	adapter        Adapter
	logger         *zap.SugaredLogger
	block          uint64
	mu             deadlock.RWMutex
	data           map[common.Address]map[common.Address]*big.Int
	contracts      []common.Address
	SubscribeQueue *SubscribeQueue
	DB             *DB
}

type UpdateCell struct {
	Contract common.Address
	Wallet   common.Address
	Value    *big.Int
	IDs      []common.Hash
	Block    uint64
}

type QueueInit struct {
	contract common.Address
	wallet   common.Address
}

type DoInit struct {
}

type NewBlock struct {
	block uint64
}

type InsertOrUpdateToDB struct {
	contract common.Address
	wallet   common.Address
	value    *big.Int
}

type FlushToDB struct {
	block uint64
}

func NewTokens(db *pgxpool.Pool, adapter Adapter, logger *zap.SugaredLogger, output chan UpdateCell) *Tokens {

	updates := make(chan any)
	updatesDB := make(chan any)
	blockchainID, blockchainName, _ := adapter.GetInfo()

	return &Tokens{
		updates:        updates,
		updatesDB:      updatesDB,
		output:         output,
		adapter:        adapter,
		logger:         logger,
		block:          0,
		mu:             deadlock.RWMutex{},
		data:           map[common.Address]map[common.Address]*big.Int{},
		contracts:      nil,
		SubscribeQueue: NewSubscribeQueue(updates),
		DB:             NewDB(updatesDB, db, blockchainID, blockchainName),
	}
}

func (t *Tokens) Run() error {

	cells, err := t.DB.loadTokensFromDB()
	if err != nil {
		return err
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	for _, cell := range cells {
		contract := cell.Contract
		_, ok := t.data[contract]
		if !ok {
			t.data[contract] = make(map[common.Address]*big.Int)
			t.contracts = append(t.contracts, contract)
		}
		t.data[contract][cell.Wallet] = cell.Value
		t.SubscribeQueue.MarkAsLoadedFromDB(contract, cell.Wallet)
	}

	if len(cells) > 0 {
		t.block = cells[0].Block
	} else {
		block, err := t.adapter.GetLastBlockNumber()
		if err != nil {
			fmt.Println(err)
		}

		t.block = block
		if t.block > 0 {
			t.block--
		}
	}

	t.DB.Run()
	t.adapter.RegisterNewBlockListener(t.newBlockTicker)

	go t.worker()
	t.runInitTicker()

	return nil
}

func (t *Tokens) runInitTicker() {
	ticker := time.NewTicker(300 * time.Millisecond)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				t.updates <- DoInit{}
			}
		}
	}()
}

func (t *Tokens) newBlockTicker(blockNumber uint64) {
	t.updates <- NewBlock{
		block: blockNumber,
	}
}

func (t *Tokens) worker() {
	initJobs := make([]QueueInit, 0)
	var wg sync.WaitGroup
	for {
		select {
		case update := <-t.updates:
			switch u := update.(type) {
			case QueueInit:
				fmt.Println("QueueInit", u.contract.Hex(), u.wallet.Hex())
				initJobs = append(initJobs, u)
			case DoInit:
				for _, j := range initJobs {
					wg.Add(1)
					go func(c common.Address, w common.Address) {
						fmt.Println("Init", c, w)
						balance, _, err := t.adapter.GetTokenBalance(&c, &w, t.block)
						if err != nil {
							t.logger.Error(err)
						}
						t.setCell(c, w, balance)
						wg.Done()
					}(j.contract, j.wallet)
				}
				wg.Wait()
				initJobs = make([]QueueInit, 0)
				t.updatesDB <- FlushToDB{
					block: t.block,
				}
			case NewBlock:
				fmt.Println("NewBlock", u.block)
				if u.block <= t.block {
					break
				}
				adapterLogs, err := t.adapter.GetTokenLogs(t.block+1, u.block, t.contracts)
				if err != nil {
					t.logger.Error(err)
				}

				for _, l := range adapterLogs {
					log, ok := l.(*TransferERC20Log)
					if !ok {
						t.logger.Error("Log variable must has *TransferERC20Log type")
						continue
					}

					t.updateCell(log.Contract, log.From, log.Value.Neg(log.Value))
					t.updateCell(log.Contract, log.To, log.Value)
				}

				// Only place where we update block!
				t.block = u.block
				t.updatesDB <- FlushToDB{
					block: t.block,
				}
			}
		}
	}
}

func (t *Tokens) setCell(contract common.Address, wallet common.Address, value *big.Int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	_, ok := t.data[contract]
	if !ok {
		t.data[contract] = make(map[common.Address]*big.Int)
		t.contracts = append(t.contracts, contract)
	}
	t.data[contract][wallet] = value

	t.updatesDB <- InsertOrUpdateToDB{
		contract: contract,
		wallet:   wallet,
		value:    value,
	}

	t.output <- UpdateCell{
		Contract: contract,
		Wallet:   wallet,
		Value:    t.data[contract][wallet],
		Block:    t.block,
	}

	fmt.Println("setCell ", contract.Hex(), wallet.Hex(), t.block, t.data[contract][wallet].String())
}

func (t *Tokens) updateCell(contract common.Address, wallet common.Address, value *big.Int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	_, ok := t.data[contract]
	if !ok {
		return
	}
	_, ok = t.data[contract][wallet]
	if !ok {
		return
	}

	// Update only existing cells
	t.data[contract][wallet].Add(t.data[contract][wallet], value)

	t.updatesDB <- InsertOrUpdateToDB{
		contract: contract,
		wallet:   wallet,
		value:    t.data[contract][wallet],
	}

	t.output <- UpdateCell{
		Contract: contract,
		Wallet:   wallet,
		Value:    t.data[contract][wallet],
		Block:    t.block,
	}

	fmt.Println("updateCell ", contract.Hex(), wallet.Hex(), t.block, t.data[contract][wallet].String())
}

func (t *Tokens) AddContract(contract common.Address) error {
	return t.SubscribeQueue.AddContract(contract)
}

func (t *Tokens) AddWallet(wallet common.Address) error {
	return t.SubscribeQueue.AddWallet(wallet)
}
