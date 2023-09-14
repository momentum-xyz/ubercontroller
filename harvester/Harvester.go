package harvester

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sasha-s/go-deadlock"
	"go.uber.org/zap"
)

type Harvester struct {
	tokens      *Tokens
	nfts        *NFTs
	ethers      *Ethers
	adapter     Adapter
	mu          deadlock.RWMutex
	logger      *zap.SugaredLogger
	pool        *pgxpool.Pool
	updateCells chan UpdateCell
	outputs     []chan any
}

type TokenCell struct {
	Contract common.Address
	Wallet   common.Address
	Value    *big.Int
	Block    uint64
}

type NFTCell struct {
	Contract common.Address
	Wallet   common.Address
	ItemID   common.Hash
	Block    uint64
}

func NewHarvester(updateCells chan UpdateCell, pool *pgxpool.Pool, adapter Adapter, logger *zap.SugaredLogger) *Harvester {
	return &Harvester{
		tokens:      NewTokens(pool, adapter, logger, updateCells),
		nfts:        NewNFTs(pool, adapter, logger, updateCells),
		ethers:      NewEthers(pool, adapter, logger, updateCells),
		adapter:     adapter,
		logger:      logger,
		pool:        pool,
		updateCells: updateCells,
		mu:          deadlock.RWMutex{},
	}
}

func (h *Harvester) Run() error {
	err := h.tokens.Run()
	if err != nil {
		return err
	}
	err = h.nfts.Run()
	if err != nil {
		return err
	}
	err = h.ethers.Run()
	if err != nil {
		return err
	}

	return nil
}

func (h *Harvester) AddTokenContract(tokenContract common.Address) error {
	return h.tokens.AddContract(tokenContract)
}

func (h *Harvester) AddNFTContract(contract common.Address) error {
	return h.nfts.AddContract(contract)
}

func (h *Harvester) AddWallet(wallet common.Address) error {
	err := h.tokens.AddWallet(wallet)
	if err != nil {
		return err
	}
	err = h.nfts.AddWallet(wallet)
	if err != nil {
		return err
	}

	return h.ethers.AddWallet(wallet)
}
