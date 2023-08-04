package harvester3

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sasha-s/go-deadlock"
	"go.uber.org/zap"
)

type TokenMatrix struct {
	mu          deadlock.RWMutex
	data        map[common.Address]map[common.Address]*TokenCell
	delta       map[common.Address]map[common.Address]map[int64]*big.Int
	wallets     map[common.Address]bool
	contracts   []common.Address
	adapter     Adapter
	jobs        chan *TokenCell
	blockNumber int64
	startBlock  int64
	logger      *zap.SugaredLogger
}

func NewTokenMatrix(adapter Adapter, logger *zap.SugaredLogger) *TokenMatrix {
	return &TokenMatrix{
		data:    map[common.Address]map[common.Address]*TokenCell{},
		delta:   make(map[common.Address]map[common.Address]map[int64]*big.Int),
		wallets: make(map[common.Address]bool),
		jobs:    make(chan *TokenCell),
		adapter: adapter,
		logger:  logger,
	}
}

func (m *TokenMatrix) Run() error {
	err := m.LoadFromDB()
	if err != nil {
		return err
	}

	block, err := m.adapter.GetLastBlockNumber()
	if err != nil {
		fmt.Println(err)
	}
	m.startBlock = int64(block)

	m.adapter.RegisterNewBlockListener(m.newBlockTicker)

	go m.runWorker(1)

	return nil
}

func (m *TokenMatrix) newBlockTicker(blockNumber uint64) {
	fmt.Println("New block", blockNumber)

	fromBlock := m.blockNumber
	if fromBlock == 0 {
		fromBlock = m.startBlock
	}
	fromBlock += 1

	toBlock := int64(blockNumber)

	logs, err := m.adapter.GetTokenLogs(fromBlock, toBlock, m.contracts)
	if err != nil {
		m.logger.Error(err)
	}

	m.processLogs(logs)

	m.blockNumber = int64(blockNumber)
}

func (m *TokenMatrix) processLogs(logs []any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, l := range logs {
		log, ok := l.(*TransferERC20Log)
		if !ok {
			m.logger.Error("Log variable must has *TransferERC20Log type")
			continue
		}

		if _, ok := m.data[log.Contract]; !ok {
			m.logger.Error("Got log for contract which is not in the list")
			continue
		}

		if _, ok := m.data[log.Contract][log.To]; !ok {
			// This wallet not subscribed
			continue
		}

		if _, ok := m.data[log.Contract][log.To]; ok {
			m.updateCellOrDelta(log.Contract, log.To, log.Block, log.Value)
		}
		if _, ok := m.data[log.Contract][log.From]; ok {
			m.updateCellOrDelta(log.Contract, log.To, log.Block, log.Value.Neg(log.Value))
		}

		////
	}
}

func (m *TokenMatrix) updateCellOrDelta(contract common.Address, wallet common.Address, block int64, value *big.Int) {
	v := m.data[contract][wallet]
	if v.isInit {
		// If cell initialised apply delta immediately
		if v.block < block {
			v.block = block
			v.value.Add(v.value, value)
		}
	} else {
		if _, ok := m.delta[contract]; !ok {
			m.delta[contract] = make(map[common.Address]map[int64]*big.Int)
		}
		if _, ok := m.delta[contract][wallet]; !ok {
			m.delta[contract] = make(map[common.Address]map[int64]*big.Int)
		}
		m.delta[contract][wallet] = make(map[int64]*big.Int)
		m.delta[contract][wallet][block] = value
	}
}

func (m *TokenMatrix) LoadFromDB() error {
	//TODO Load from DB

	return nil
}

func (m *TokenMatrix) AddContract(contract common.Address) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.data[contract]; ok {
		// contract already subscribed
		return nil
	}

	m.contracts = append(m.contracts, contract)

	m.data[contract] = make(map[common.Address]*TokenCell)
	for w, _ := range m.wallets {
		m.data[contract][w] = NewTokenCell(contract, w)

		go m.initCellFromBC(m.data[contract][w])
	}

	return nil
}

func (m *TokenMatrix) AddWallet(wallet common.Address) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.wallets[wallet]; ok {
		// Wallet already subscribed
		return nil
	}

	m.wallets[wallet] = true

	for c, _ := range m.data {
		m.data[c][wallet] = NewTokenCell(c, wallet)

		go m.initCellFromBC(m.data[c][wallet])
	}

	return nil
}

func (m *TokenMatrix) initCellFromBC(c *TokenCell) {
	fmt.Println("initCellFromBC", *c)
	m.jobs <- c
}

func (m *TokenMatrix) runWorker(workerID int) {
	for cell := range m.jobs {
		m.doInitialiseCell(cell)
	}
}

func (m *TokenMatrix) doInitialiseCell(cell *TokenCell) {
	block := m.blockNumber
	if block == 0 {
		block = m.startBlock
	}
	balance, block, err := m.adapter.GetTokenBalance(&cell.contract, &cell.wallet, uint64(block))
	if err != nil {
		m.logger.Error(err)
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	cell.isInit = true
	cell.value = balance
	cell.initBlock = block
	cell.block = block

	defer func() {
		fmt.Println("Initialized " + cell.contract.Hex() + " " + cell.wallet.Hex() + " " + strconv.Itoa(int(cell.block)) + " =" + cell.value.String())
	}()

	_, ok := m.delta[cell.contract]
	if !ok {
		return
	}

	_, ok = m.delta[cell.contract][cell.wallet]
	if !ok {
		return
	}

	bmap := m.delta[cell.contract][cell.wallet]
	for logBlock, logValue := range bmap {
		if logBlock > block {
			// Change by that block already counted
			cell.value.Add(cell.value, logValue)
			cell.block = logBlock
		}

		delete(bmap, logBlock)
	}
}
