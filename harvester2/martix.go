package harvester2

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"github.com/sasha-s/go-deadlock"
	"math/big"
	"sync"
)

type Wallet Address
type Contract Address

type Matrix struct {
	mu          deadlock.RWMutex
	blockNumber uint64
	tokenMatrix map[*Contract]map[*Wallet]*TokenCell
	nftMatrix   map[*Contract]map[*Wallet]*NFTCell
	stakeMatrix map[*Contract]map[*Wallet]*StakeCell
	db          *pgxpool.Pool
	adapter     Adapter

	wallets   map[*Address]bool
	contracts map[*Address]bool
	//harvesterListener func(bcName string, p []*UpdateEvent, s []*StakeEvent)
}

type TokenCell struct {
	isInit bool
	value  *big.Int
}

type NFTCell struct {
	isInit bool
	value  map[umid.UMID]int8
}

type StakeCell struct {
	isInit bool
	Stakes map[umid.UMID]*Stake
}

func NewMatrix(db *pgxpool.Pool, adapter Adapter) *Matrix {
	return &Matrix{
		blockNumber: 0,
		tokenMatrix: make(map[*Contract]map[*Wallet]*TokenCell),
		stakeMatrix: make(map[*Contract]map[*Wallet]*StakeCell),
		nftMatrix:   make(map[*Contract]map[*Wallet]*NFTCell),
		adapter:     adapter,
		//harvesterListener: listener,
		db: db,

		wallets:   make(map[*Address]bool),
		contracts: make(map[*Address]bool),
	}
}

func (m *Matrix) Run() {
	//t.fastForward()

	m.mu.Lock()
	defer m.mu.Unlock()

	block, err := m.adapter.GetLastBlockNumber()
	if err != nil {
		fmt.Println(err)
	}
	m.blockNumber = block

	m.adapter.RegisterNewBlockListener(m.listener)
}

func (m *Matrix) listener(blockNumber uint64) {
	//t.fastForward()
	//t.mu.Lock()
	//t.ProcessDiffs(blockNumber, diffs, stakes)
	//t.mu.Unlock()
}

func (m *Matrix) fillMissingDataForContract(contract *Address, wg *sync.WaitGroup) {
	if contract == nil {
		return
	}
	c := (common.Address)(*contract)
	// Get all logs for given contract from beginning to current block
	logs, err := m.adapter.GetLogs(int64(m.blockNumber)+1, 0, []common.Address{c})
	if err != nil {
		fmt.Println(err)
		return
	}

	m.ProcessLogs(m.blockNumber, logs, m.wallets)
	if wg != nil {
		wg.Done()
	}
}

func (m *Matrix) fillStakeMatrixCell(block uint64, contract *Contract, wallet *Wallet, wg *sync.WaitGroup) {
	fmt.Println("fillStakeMatrixCell")

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.stakeMatrix[contract][wallet].isInit {
		return
	}
	fmt.Println("fillStakeMatrixCell 2")

	stakesMap, err := m.adapter.GetStakeBalance(int64(block), (*common.Address)(wallet), (*common.Address)(contract))
	if err != nil {
		fmt.Println("ERROR: fillStakeMatrixCell: Failed to GetStakeBalance")
	}

	m.stakeMatrix[contract][wallet].isInit = true

	for id, val := range stakesMap {
		if _, ok := m.stakeMatrix[contract][wallet].Stakes[id]; !ok {
			m.stakeMatrix[contract][wallet].Stakes[id] = &Stake{
				TotalAmount:    big.NewInt(0),
				TotalMOMAmount: big.NewInt(0),
				TotalDADAmount: big.NewInt(0),
			}
		}
		s := m.stakeMatrix[contract][wallet].Stakes[id]
		s.TotalAmount.Add(s.TotalAmount, val[0])
		s.TotalMOMAmount.Add(s.TotalMOMAmount, val[1])
		s.TotalDADAmount.Add(s.TotalDADAmount, val[2])
	}

	if wg != nil {
		wg.Done()
	}
}

func (m *Matrix) fillNFTMatrixCell(block uint64, contract *Contract, wallet *Wallet, wg *sync.WaitGroup) {
	fmt.Println("fillNFTMatrixCell")

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.nftMatrix[contract][wallet].isInit {
		return
	}

	nfts, err := m.adapter.GetNFTBalance(int64(block), (*common.Address)(wallet), (*common.Address)(contract))
	if err != nil {
		fmt.Println("ERROR: fillNFTMatrixCell: Failed to get NFTs balance")
	}
	if _, ok := m.nftMatrix[contract]; !ok {
		m.nftMatrix[contract] = make(map[*Wallet]*NFTCell)
	}

	for _, nft := range nfts {
		m.nftMatrix[contract][wallet].value[nft] += 1
	}
	m.nftMatrix[contract][wallet].isInit = true

	if wg != nil {
		wg.Done()
	}
}

func (m *Matrix) fillTokenMatrixCell(block uint64, contract *Contract, wallet *Wallet, wg *sync.WaitGroup) {
	//TODO start here Add mutex and set isInit=true
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.tokenMatrix[contract][wallet].isInit {
		return
	}

	fmt.Println("fillTokenMatrixCell")
	b, err := m.adapter.GetBalance((*common.Address)(wallet), (*common.Address)(contract), block)
	if err != nil {
		fmt.Println("ERROR: fillTokenMatrixCell: Failed to get token balance")
	}

	if b == nil {
		b = big.NewInt(0)
	}

	m.tokenMatrix[contract][wallet].value.Add(m.tokenMatrix[contract][wallet].value, b)

	m.tokenMatrix[contract][wallet].isInit = true

	if wg != nil {
		wg.Done()
	}
}

func (m *Matrix) fillMissingData(wgMain *sync.WaitGroup) {

	contracts := make([]common.Address, 0)
	wallets := make(map[common.Address]bool, 0)

	wg := &sync.WaitGroup{}

	for contract, val := range m.tokenMatrix {
		for wallet, cell := range val {
			if !cell.isInit {
				wg.Add(1)
				m.fillTokenMatrixCell(m.blockNumber, contract, wallet, wg)
			}
		}
	}

	for contract, val := range m.nftMatrix {
		for wallet, cell := range val {
			if !cell.isInit {
				wg.Add(1)
				m.fillNFTMatrixCell(m.blockNumber, contract, wallet, wg)
			}
		}
	}

	for contract, val := range m.stakeMatrix {
		for wallet, cell := range val {
			if !cell.isInit {
				wg.Add(1)
				m.fillStakeMatrixCell(m.blockNumber, contract, wallet, wg)
			}
		}
	}

	fmt.Println(contracts)
	fmt.Println(wallets)

	//// Get all logs for given contract from beginning to current block
	//logs, err := m.adapter.GetLogs(int64(m.blockNumber)+1, 0, []common.Address{c})
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//m.ProcessLogs(m.blockNumber, logs)
	//if wg != nil {
	//	wg.Done()
	//}
}

func (m *Matrix) ProcessLogs(blockNumber uint64, logs []any, wallets map[*Address]bool) {

	for _, log := range logs {
		switch log.(type) {
		case *TransferERC20Log:
			//l := log.(*TransferERC20Log)
			//c := (*Contract)(&l.Contract)
			//
			//w := (*Address)(&l.From)
			//if wallets[w] {
			//
			//	m.tokenMatrix[c][(*Wallet)(w)] =
			//}
		}
	}
}

func (m *Matrix) AddWallet(wallet *Address) error {
	return m.addWallet(wallet, nil)
}

func (m *Matrix) addWallet(wallet *Address, wg *sync.WaitGroup) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.wallets[wallet]; ok {
		// Wallet already subscribed
		return nil
	}

	w := (*Wallet)(wallet)
	for c := range m.tokenMatrix {
		m.tokenMatrix[c][w] = &TokenCell{
			isInit: false,
			value:  big.NewInt(0),
		}
	}
	for c := range m.nftMatrix {
		m.nftMatrix[c][w] = &NFTCell{
			isInit: false,
			value:  map[umid.UMID]int8{},
		}
	}
	for c := range m.stakeMatrix {
		m.stakeMatrix[c][w] = &StakeCell{
			isInit: false,
			Stakes: make(map[umid.UMID]*Stake),
		}
	}

	m.wallets[wallet] = true

	go m.fillMissingData(nil)

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
	m.nftMatrix[c] = make(map[*Wallet]*NFTCell)

	for wallet := range m.wallets {
		w := (*Wallet)(wallet)
		// All new cells require initial fill
		m.nftMatrix[c][w] = &NFTCell{
			isInit: false,
			value:  make(map[umid.UMID]int8),
		}
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
	m.tokenMatrix[c] = make(map[*Wallet]*TokenCell)

	for wallet := range m.wallets {
		w := (*Wallet)(wallet)
		m.tokenMatrix[c][w] = &TokenCell{
			isInit: false,
			value:  big.NewInt(0),
		}
	}

	m.contracts[contract] = true

	go m.fillMissingData(nil)

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

	m.stakeMatrix[c] = make(map[*Wallet]*StakeCell)

	for wallet := range m.wallets {
		w := (*Wallet)(wallet)
		m.stakeMatrix[c][w] = &StakeCell{
			isInit: false,
			Stakes: make(map[umid.UMID]*Stake),
		}
	}

	m.contracts[contract] = true

	go m.fillMissingData(nil)
	return nil
}

func (m *Matrix) AddTokenListener(contract *Address, listener TokenListener) error {
	m.AddTokenContract(contract)

	return nil
}

func (m *Matrix) Display() {
	fmt.Println("Token Matrix:")
	for contract, value := range m.tokenMatrix {
		for wallet, v := range value {
			fmt.Printf("%v %v %v \n", (*common.Address)(contract).Hex(), (*common.Address)(wallet).Hex(), v.value.String())
		}
	}

	fmt.Println("NFT Matrix:")
	for contract, value := range m.nftMatrix {
		for wallet, v := range value {
			fmt.Printf("%v %v %v \n", (*common.Address)(contract).Hex(), (*common.Address)(wallet).Hex(), v.value)
		}
	}

	if len(m.stakeMatrix) == 1 {
		var stakeContract *Contract
		for c, _ := range m.stakeMatrix {
			stakeContract = c
		}
		fmt.Println("STAKE Matrix:")
		fmt.Println("Contract:", (*common.Address)(stakeContract).Hex())
		for wallet, val := range m.stakeMatrix[stakeContract] {
			for id, v := range val.Stakes {
				fmt.Println((*common.Address)(wallet).Hex(), id.String(), val.isInit, v.TotalAmount, v.TotalMOMAmount, v.TotalDADAmount)
			}
		}
	}
}
