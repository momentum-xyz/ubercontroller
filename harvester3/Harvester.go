package harvester3

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sasha-s/go-deadlock"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/config"
)

type Harvester struct {
	tokens  *Tokens
	adapter Adapter
	mu      deadlock.RWMutex
	logger  *zap.SugaredLogger
	pool    *pgxpool.Pool
	input   chan UpdateCell
	//outputs map[common.Address]map[common.Address][]chan any
	outputs []chan any
}

func NewHarvester(cfg *config.Arbitrum3, pool *pgxpool.Pool, adapter Adapter, logger *zap.SugaredLogger) *Harvester {
	input := make(chan UpdateCell)
	//a := arbitrum_nova_adapter3.NewArbitrumNovaAdapter(cfg, logger)
	//a.Run()

	return &Harvester{
		tokens:  NewTokens(pool, adapter, logger, input),
		adapter: adapter,
		logger:  logger,
		pool:    pool,
		input:   input,
		mu:      deadlock.RWMutex{},
	}
}

func (h *Harvester) Run() error {
	//h.adapter.Run()
	err := h.tokens.Run()
	if err != nil {
		return err
	}

	go h.worker()

	return nil
}

func (h *Harvester) AddTokenContract(contract common.Address) error {
	return h.tokens.AddContract(contract)
}

func (h *Harvester) AddWallet(wallet common.Address) error {
	return h.tokens.AddWallet(wallet)
}

func (h *Harvester) worker() {
	for {
		in := <-h.input
		h.mu.RLock()
		for _, c := range h.outputs {
			c <- in
		}
		h.mu.RUnlock()
	}
}

func (h *Harvester) SubscribeForToken(tokenContract common.Address, wallet common.Address) (chan any, error) {
	err := h.AddTokenContract(tokenContract)
	if err != nil {
		return nil, err
	}

	err = h.AddWallet(wallet)
	if err != nil {
		return nil, err
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	c := make(chan any)
	h.outputs = append(h.outputs, c)

	return c, nil

}
