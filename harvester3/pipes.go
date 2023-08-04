package harvester3

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sasha-s/go-deadlock"
	"go.uber.org/zap"
)

type Pipes struct {
	adapter Adapter
	logger  *zap.SugaredLogger
	block   uint64

	reqBalance chan ReqBalance
	resBalance chan ResBalance
	reqLogs    chan ReqLogs
	resLogs    chan ResLogs

	mu        deadlock.RWMutex
	data      map[common.Address]map[common.Address]*TokenCell
	wallets   map[common.Address]bool
	contracts []common.Address
}

type ReqBalance struct {
	contract common.Address
	wallet   common.Address
	block    uint64
}

type ResBalance struct {
	contract common.Address
	wallet   common.Address
	value    *big.Int
	block    uint64
}

type ReqLogs struct {
	contracts []common.Address
	fromBlock uint64
	toBlock   uint64
}

type ResLogs struct {
	logs map[common.Address]map[common.Address]map[uint64]*big.Int
}

func NewPipes(adapter Adapter, logger *zap.SugaredLogger) *Pipes {
	return &Pipes{
		adapter: adapter,
		logger:  logger,
		block:   0,
		data:    map[common.Address]map[common.Address]*TokenCell{},
		wallets: make(map[common.Address]bool),

		reqBalance: make(chan ReqBalance),
		resBalance: make(chan ResBalance),
		reqLogs:    make(chan ReqLogs),
		resLogs:    make(chan ResLogs),
	}
}

func (p *Pipes) Run() {
	block, err := p.adapter.GetLastBlockNumber()
	if err != nil {
		p.logger.Error(err)
	}

	p.block = block
	if p.block > 0 {
		p.block--
	}

	go p.balanceWorker() // can be many
	go p.logsWorker()    // should be one only
	go p.mainWorker()

	p.adapter.RegisterNewBlockListener(p.newBlockTicker)
}

func (p *Pipes) balanceWorker() {
	for req := range p.reqBalance {
		balance, block, err := p.adapter.GetTokenBalance(&req.contract, &req.wallet, req.block)
		if err != nil {
			p.logger.Error(err)
		}
		p.resBalance <- ResBalance{
			contract: req.contract,
			wallet:   req.wallet,
			value:    balance,
			block:    block,
		}
	}
}

func (p *Pipes) logsWorker() {
	for req := range p.reqLogs {
		fmt.Println("logsWorker", req.fromBlock, req.toBlock)
		adapterLogs, err := p.adapter.GetTokenLogs(req.fromBlock, req.toBlock, req.contracts)
		if err != nil {
			p.logger.Error(err)
		}

		logs := make(map[common.Address]map[common.Address]map[uint64]*big.Int)

		for _, l := range adapterLogs {
			log, ok := l.(*TransferERC20Log)
			if !ok {
				p.logger.Error("Log variable must has *TransferERC20Log type")
				continue
			}

			if _, ok := logs[log.Contract]; !ok {
				logs[log.Contract] = make(map[common.Address]map[uint64]*big.Int)
			}

			if _, ok := logs[log.Contract][log.From]; !ok {
				logs[log.Contract][log.From] = make(map[uint64]*big.Int)
			}

			if _, ok := logs[log.Contract][log.To]; !ok {
				logs[log.Contract][log.To] = make(map[uint64]*big.Int)
			}

			logs[log.Contract][log.From][log.Block] = big.NewInt(0).Neg(log.Value)
			logs[log.Contract][log.To][log.Block] = log.Value
		}

		p.resLogs <- ResLogs{logs: logs}
	}
}

func (p *Pipes) mainWorker() {
	for {
		select {
		case res := <-p.resBalance:
			fmt.Printf("resBalance %v %v %v %v \n", res.contract, res.wallet, res.block, res.value.String())
		case res := <-p.resLogs:
			fmt.Println("resLogs", res)

			for c := range res.logs {
				for w := range res.logs[c] {
					for b := range res.logs[c][w] {
						fmt.Printf("%v %v %v %v\n", c.Hex(), w.Hex(), b, res.logs[c][w][b].String())
					}
				}
			}
		}
		///
	}
}

func (p *Pipes) newBlockTicker(blockNumber uint64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if blockNumber > p.block {

		fromBlock := p.block + 1
		p.block = blockNumber

		p.reqLogs <- ReqLogs{
			contracts: nil,
			fromBlock: fromBlock,
			toBlock:   blockNumber,
		}
	}
}

func (p *Pipes) AddContract(contract common.Address) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.data[contract]; ok {
		// contract already subscribed
		return nil
	}

	p.contracts = append(p.contracts, contract)

	p.data[contract] = make(map[common.Address]*TokenCell)
	for w, _ := range p.wallets {
		p.data[contract][w] = NewTokenCell(contract, w)
		p.reqBalance <- ReqBalance{
			contract: contract,
			wallet:   w,
			block:    p.block,
		}
	}

	return nil
}

func (p *Pipes) AddWallet(wallet common.Address) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.wallets[wallet]; ok {
		// Wallet already subscribed
		return nil
	}

	p.wallets[wallet] = true

	for c, _ := range p.data {
		p.data[c][wallet] = NewTokenCell(c, wallet)

		p.reqBalance <- ReqBalance{
			contract: c,
			wallet:   wallet,
			block:    p.block,
		}
	}

	return nil
}
