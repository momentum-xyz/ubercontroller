package harvester3

import (
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sasha-s/go-deadlock"
	"go.uber.org/zap"
)

type NFTs struct {
	updates        chan any
	updatesDB      chan<- any
	output         chan<- UpdateCell
	adapter        Adapter
	logger         *zap.SugaredLogger
	block          uint64
	mu             deadlock.RWMutex
	data           map[common.Address]map[common.Address]map[common.Hash]bool
	contracts      []common.Address
	SubscribeQueue *SubscribeQueue
	DB             *DB
}

type UpsertNFTToDB struct {
	contract common.Address
	wallet   common.Address
	id       common.Hash
}

type UpsertEtherToDB struct {
	wallet common.Address
	value  *big.Int
}

type FlushNFTToDB struct {
	block uint64
}

type FlushEthersToDB struct {
	block uint64
}

type RemoveNFTFromDB struct {
	contract common.Address
	wallet   common.Address
	id       common.Hash
}

type UpdateNFTEvent struct {
	Contract common.Address
	Wallet   common.Address
	IDs      []common.Hash
}

func NewNFTs(db *pgxpool.Pool, adapter Adapter, logger *zap.SugaredLogger, output chan<- UpdateCell) *NFTs {

	updates := make(chan any)
	updatesDB := make(chan any)
	blockchainID, blockchainName, _ := adapter.GetInfo()

	return &NFTs{
		updates:        updates,
		updatesDB:      updatesDB,
		output:         output,
		adapter:        adapter,
		logger:         logger,
		block:          0,
		mu:             deadlock.RWMutex{},
		data:           make(map[common.Address]map[common.Address]map[common.Hash]bool),
		contracts:      nil,
		SubscribeQueue: NewSubscribeQueue(updates),
		DB:             NewDB(updatesDB, db, logger, blockchainID, blockchainName),
	}
}

func (n *NFTs) Run() error {

	cells, err := n.DB.loadNFTsFromDB()
	if err != nil {
		return err
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	for _, cell := range cells {
		contract := cell.Contract
		_, ok := n.data[contract]
		if !ok {
			n.data[contract] = make(map[common.Address]map[common.Hash]bool)
			n.contracts = append(n.contracts, contract)
		}
		_, ok = n.data[contract][cell.Wallet]
		if !ok {
			n.data[contract][cell.Wallet] = make(map[common.Hash]bool)
		}
		n.data[contract][cell.Wallet][cell.ItemID] = true

		n.SubscribeQueue.MarkAsLoadedFromDB(contract, cell.Wallet)
	}

	if len(cells) > 0 {
		n.block = cells[0].Block
	} else {
		block, err := n.adapter.GetLastBlockNumber()
		if err != nil {
			n.logger.Error(err)
		}

		n.block = block
		if n.block > 0 {
			n.block--
		}
	}

	n.DB.Run()
	n.adapter.RegisterNewBlockListener(n.newBlockTicker)

	go n.worker()
	n.runInitTicker()

	return nil
}

func (n *NFTs) runInitTicker() {
	ticker := time.NewTicker(300 * time.Millisecond)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				n.updates <- DoInit{}
			}
		}
	}()
}

func (n *NFTs) newBlockTicker(blockNumber uint64) {
	n.updates <- NewBlock{
		block: blockNumber,
	}
}

func (n *NFTs) worker() {
	initJobs := make([]QueueInit, 0)
	var wg sync.WaitGroup
	for {
		select {
		case update := <-n.updates:
			switch u := update.(type) {
			case QueueInit:
				n.logger.Debug("QueueInit", u.contract.Hex(), u.wallet.Hex())
				initJobs = append(initJobs, u)
			case DoInit:
				for _, j := range initJobs {
					wg.Add(1)
					go func(c common.Address, w common.Address) {
						n.logger.Debug("Init", c, w)
						ids, err := n.adapter.GetNFTBalance(&c, &w, n.block)
						if err != nil {
							n.logger.Error(err)
						}
						n.setCell(c, w, ids)
						wg.Done()
					}(j.contract, j.wallet)
				}
				wg.Wait()
				initJobs = make([]QueueInit, 0)
				n.updatesDB <- FlushNFTToDB{
					block: n.block,
				}
			case NewBlock:
				n.logger.Debug("NewBlock ", u.block)
				if u.block <= n.block {
					break
				}
				if n.contracts == nil {
					continue
				}
				adapterLogs, err := n.adapter.GetNFTLogs(n.block+1, u.block, n.contracts)
				if err != nil {
					n.logger.Error(err)
				}

				for _, l := range adapterLogs {
					log, ok := l.(*TransferNFTLog)
					if !ok {
						n.logger.Error("Log variable must has *TransferERC20Log type")
						continue
					}

					n.updateCell(log.Contract, log.From, log.TokenID, "remove")
					n.updateCell(log.Contract, log.To, log.TokenID, "add")
				}

				// Only place where we update block!
				n.block = u.block
				n.updatesDB <- FlushNFTToDB{
					block: n.block,
				}
			}
		}
	}
}

func (n *NFTs) setCell(contract common.Address, wallet common.Address, ids []common.Hash) {
	n.mu.Lock()
	defer n.mu.Unlock()
	_, ok := n.data[contract]
	if !ok {
		n.data[contract] = make(map[common.Address]map[common.Hash]bool)
		n.contracts = append(n.contracts, contract)
	}

	_, ok = n.data[contract][wallet]
	if !ok {
		n.data[contract][wallet] = make(map[common.Hash]bool)
	}

	for _, id := range ids {
		n.data[contract][wallet][id] = true

		n.updatesDB <- UpsertNFTToDB{
			contract: contract,
			wallet:   wallet,
			id:       id,
		}
	}

	n.output <- UpdateCell{
		Contract: contract,
		Wallet:   wallet,
		IDs:      ids,
	}
}

func (n *NFTs) updateCell(contract common.Address, wallet common.Address, id common.Hash, updateType string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	_, ok := n.data[contract]
	if !ok {
		return
	}
	_, ok = n.data[contract][wallet]
	if !ok {
		return
	}

	if updateType == "add" {
		n.data[contract][wallet][id] = true
		n.updatesDB <- UpsertNFTToDB{
			contract: contract,
			wallet:   wallet,
			id:       id,
		}
	} else {
		delete(n.data[contract][wallet], id)
		n.updatesDB <- RemoveNFTFromDB{
			contract: contract,
			wallet:   wallet,
			id:       id,
		}
	}

	ids := make([]common.Hash, 0)

	for tokenID := range n.data[contract][wallet] {
		ids = append(ids, tokenID)
	}

	n.output <- UpdateCell{
		Contract: contract,
		Wallet:   wallet,
		IDs:      ids,
	}
}

func (n *NFTs) AddContract(contract common.Address) error {
	return n.SubscribeQueue.AddContract(contract)
}

func (n *NFTs) AddWallet(wallet common.Address) error {
	return n.SubscribeQueue.AddWallet(wallet)
}
