package harvester2

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"github.com/sasha-s/go-deadlock"
	"math/big"
)

type Wallet Address
type Contract Address

type Matrix struct {
	mu          deadlock.RWMutex
	blockNumber uint64
	tokenMatrix map[*Contract]map[*Wallet]*big.Int
	nftMatrix   map[*Contract]map[*Wallet][]*umid.UMID
	stakeMatrix map[*Contract]map[*Wallet]map[*umid.UMID]*Stake
	db          *pgxpool.Pool
	adapter     Adapter

	wallets   map[*Address]bool
	contracts map[*Address]bool
	//harvesterListener func(bcName string, p []*UpdateEvent, s []*StakeEvent)
}

func NewMatrix(db *pgxpool.Pool, adapter Adapter) *Matrix {
	return &Matrix{
		blockNumber: 0,
		tokenMatrix: make(map[*Contract]map[*Wallet]*big.Int),
		stakeMatrix: make(map[*Contract]map[*Wallet]map[*umid.UMID]*Stake),
		nftMatrix:   make(map[*Contract]map[*Wallet][]*umid.UMID),
		adapter:     adapter,
		//harvesterListener: listener,
		db: db,

		wallets:   make(map[*Address]bool),
		contracts: make(map[*Address]bool),
	}
}

func (m *Matrix) Run() {
	//t.fastForward()

	m.adapter.RegisterNewBlockListener(m.listener)
}

func (m *Matrix) listener(blockNumber uint64) {
	//t.fastForward()
	//t.mu.Lock()
	//t.ProcessDiffs(blockNumber, diffs, stakes)
	//t.mu.Unlock()
}

func (m *Matrix) AddWallet(wallet *Address) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.wallets[wallet]; ok {
		// Wallet already subscribed
		return nil
	}

	w := (*Wallet)(wallet)
	for c := range m.tokenMatrix {
		m.tokenMatrix[c][w] = nil
	}
	for c := range m.nftMatrix {
		m.nftMatrix[c][w] = nil
	}
	for c := range m.stakeMatrix {
		m.stakeMatrix[c][w] = make(map[*umid.UMID]*Stake)
	}

	m.wallets[wallet] = true

	return nil
}

func (m *Matrix) AddNFTContract(contract *Address) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.contracts[contract]; ok {
		// Contract already subscribed
		return nil
	}

	c := (*Contract)(contract)
	m.nftMatrix[c] = make(map[*Wallet][]*umid.UMID)

	for wallet := range m.wallets {
		w := (*Wallet)(wallet)
		m.nftMatrix[c][w] = nil
	}

	m.contracts[contract] = true
	return nil
}

func (m *Matrix) AddTokenContract(contract *Address) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.contracts[contract]; ok {
		// Contract already subscribed
		return nil
	}

	c := (*Contract)(contract)
	m.tokenMatrix[c] = make(map[*Wallet]*big.Int)

	for wallet := range m.wallets {
		w := (*Wallet)(wallet)
		m.tokenMatrix[c][w] = nil
	}

	m.contracts[contract] = true
	return nil
}

func (m *Matrix) AddStakeContract(contract *Address) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.contracts[contract]; ok {
		// Contract already subscribed
		return nil
	}

	c := (*Contract)(contract)

	m.stakeMatrix[c] = make(map[*Wallet]map[*umid.UMID]*Stake)

	for wallet := range m.wallets {
		w := (*Wallet)(wallet)
		m.stakeMatrix[c][w] = make(map[*umid.UMID]*Stake)
	}

	m.contracts[contract] = true
	return nil
}
