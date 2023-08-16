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

type Ethers struct {
	updates         chan any
	updatesDB       chan any
	output          chan UpdateCell
	adapter         Adapter
	logger          *zap.SugaredLogger
	block           uint64
	mu              deadlock.RWMutex
	data            map[common.Address]map[common.Address]*big.Int
	contracts       []common.Address
	wallets         map[common.Address]bool
	SubscribeQueue  *SubscribeQueue
	DB              *DB
	DefaultContract common.Address
}

func NewEthers(db *pgxpool.Pool, adapter Adapter, logger *zap.SugaredLogger, output chan UpdateCell) *Ethers {

	updates := make(chan any)
	updatesDB := make(chan any)
	blockchainID, blockchainName, _ := adapter.GetInfo()

	return &Ethers{
		updates:         updates,
		updatesDB:       updatesDB,
		output:          output,
		adapter:         adapter,
		logger:          logger,
		block:           0,
		mu:              deadlock.RWMutex{},
		data:            map[common.Address]map[common.Address]*big.Int{},
		contracts:       nil,
		wallets:         make(map[common.Address]bool),
		SubscribeQueue:  NewSubscribeQueue(updates),
		DB:              NewDB(updatesDB, db, blockchainID, blockchainName),
		DefaultContract: common.HexToAddress("0x0000000000000000000000000000000000000001"),
	}
}

func (e *Ethers) Run() error {
	err := e.SubscribeQueue.AddContract(e.DefaultContract)
	if err != nil {
		return err
	}

	cells, err := e.DB.loadEthersFromDB()
	if err != nil {
		return err
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	for _, cell := range cells {
		contract := cell.Contract
		_, ok := e.data[contract]
		if !ok {
			e.data[contract] = make(map[common.Address]*big.Int)
			e.contracts = append(e.contracts, contract)
		}
		e.data[contract][cell.Wallet] = cell.Value
		e.wallets[cell.Wallet] = true
		e.SubscribeQueue.MarkAsLoadedFromDB(contract, cell.Wallet)
	}

	if len(cells) > 0 {
		e.block = cells[0].Block
	} else {
		block, err := e.adapter.GetLastBlockNumber()
		if err != nil {
			e.logger.Error(err)
		}

		e.block = block
		if e.block > 0 {
			e.block--
		}
	}

	e.DB.Run()
	e.adapter.RegisterNewBlockListener(e.newBlockTicker)

	go e.worker()
	e.runInitTicker()

	return nil
}

func (e *Ethers) runInitTicker() {
	ticker := time.NewTicker(300 * time.Millisecond)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				e.updates <- DoInit{}
			}
		}
	}()
}

func (e *Ethers) newBlockTicker(blockNumber uint64) {
	e.updates <- NewBlock{
		block: blockNumber,
	}
}

func (e *Ethers) worker() {
	initJobs := make([]QueueInit, 0)
	var wg sync.WaitGroup
	for {
		select {
		case update := <-e.updates:
			switch u := update.(type) {
			case QueueInit:
				fmt.Println("QueueInit", u.contract.Hex(), u.wallet.Hex())
				initJobs = append(initJobs, u)
			case DoInit:
				for _, j := range initJobs {
					wg.Add(1)
					go func(c common.Address, w common.Address) {
						fmt.Println("Init", c, w)
						balance, err := e.adapter.GetEtherBalance(&w, e.block)
						if err != nil {
							e.logger.Error(err)
						}
						e.setCell(c, w, balance)
						wg.Done()
					}(j.contract, j.wallet)
				}
				wg.Wait()
				initJobs = make([]QueueInit, 0)
				e.updatesDB <- FlushEthersToDB{
					block: e.block,
				}
			case NewBlock:
				fmt.Println("NewBlock", u.block)
				if u.block <= e.block {
					break
				}
				adapterLogs, err := e.adapter.GetEtherLogs(e.block+1, u.block, e.wallets)
				if err != nil {
					e.logger.Error(err)
				}

				for _, log := range adapterLogs {
					e.updateCell(e.DefaultContract, log.Wallet, log.Delta)
				}

				// Only place where we update block!
				e.block = u.block
				e.updatesDB <- FlushEthersToDB{
					block: e.block,
				}
			}
		}
	}
}

func (e *Ethers) setCell(contract common.Address, wallet common.Address, value *big.Int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	_, ok := e.data[contract]
	if !ok {
		e.data[contract] = make(map[common.Address]*big.Int)
		e.contracts = append(e.contracts, contract)
	}
	e.data[contract][wallet] = value
	e.wallets[wallet] = true

	e.updatesDB <- UpsertEtherToDB{
		wallet: wallet,
		value:  value,
	}

	e.output <- UpdateCell{
		Contract: contract,
		Wallet:   wallet,
		Value:    e.data[contract][wallet],
		Block:    e.block,
	}

	fmt.Println("setCell ", contract.Hex(), wallet.Hex(), e.block, e.data[contract][wallet].String())
}

func (e *Ethers) updateCell(contract common.Address, wallet common.Address, value *big.Int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	_, ok := e.data[contract]
	if !ok {
		return
	}
	_, ok = e.data[contract][wallet]
	if !ok {
		return
	}

	// Update only existing cells
	e.data[contract][wallet].Add(e.data[contract][wallet], value)

	e.updatesDB <- UpsertEtherToDB{
		wallet: wallet,
		value:  e.data[contract][wallet],
	}

	e.output <- UpdateCell{
		Contract: contract,
		Wallet:   wallet,
		Value:    e.data[contract][wallet],
		Block:    e.block,
	}

	fmt.Println("updateCell ", contract.Hex(), wallet.Hex(), e.block, e.data[contract][wallet].String())
}

//func (e *Ethers) AddContract(contract common.Address) error {
//	return e.SubscribeQueue.AddContract(contract)
//}

func (e *Ethers) AddWallet(wallet common.Address) error {
	return e.SubscribeQueue.AddWallet(wallet)
}
