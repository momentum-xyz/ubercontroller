package harvester2

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"github.com/pkg/errors"
	"github.com/sasha-s/go-deadlock"
	"math/big"
	"sync"
)

type Wallet Address
type Contract Address

type Matrix struct {
	mu          deadlock.RWMutex
	blockNumber uint64
	tokenMatrix map[Contract]map[Wallet]*TokenCell
	nftMatrix   map[Contract]map[Wallet]*NFTCell
	stakeMatrix map[Contract]map[Wallet]*StakeCell
	db          *pgxpool.Pool
	adapter     Adapter

	wallets   map[Address]bool
	contracts map[Address]bool
	//harvesterListener func(bcName string, p []*UpdateEvent, s []*StakeEvent)

	tokenListeners map[Contract]map[string]TokenListener
	nftListeners   map[Contract]map[string]NFTListener
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

var ZeroAddress = common.Address{}

func NewMatrix(db *pgxpool.Pool, adapter Adapter) *Matrix {
	return &Matrix{
		blockNumber: 0,
		tokenMatrix: make(map[Contract]map[Wallet]*TokenCell),
		stakeMatrix: make(map[Contract]map[Wallet]*StakeCell),
		nftMatrix:   make(map[Contract]map[Wallet]*NFTCell),
		adapter:     adapter,
		//harvesterListener: listener,
		db: db,

		wallets:   make(map[Address]bool),
		contracts: make(map[Address]bool),

		tokenListeners: make(map[Contract]map[string]TokenListener),
		nftListeners:   make(map[Contract]map[string]NFTListener),
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

	m.adapter.RegisterNewBlockListener(m.newBlockTicker)
}

func (m *Matrix) newBlockTicker(blockNumber uint64) {
	m.fastForward()
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

	m.ProcessLogs(m.blockNumber, logs)
	if wg != nil {
		wg.Done()
	}
}

func (m *Matrix) fillStakeMatrixCell(block uint64, contract Contract, wallet Wallet, wg *sync.WaitGroup) {
	fmt.Println("fillStakeMatrixCell")

	//m.mu.Lock()
	//defer m.mu.Unlock()

	if m.stakeMatrix[contract][wallet].isInit {
		return
	}

	stakesMap, err := m.adapter.GetStakeBalance(int64(block), (*common.Address)(&wallet), (*common.Address)(&contract))
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

	m.saveAllToDB()

	if wg != nil {
		wg.Done()
	}
}

func (m *Matrix) fillNFTMatrixCell(block uint64, contract Contract, wallet Wallet, wg *sync.WaitGroup) {
	fmt.Println("fillNFTMatrixCell")

	//m.mu.Lock()
	//defer m.mu.Unlock()

	if m.nftMatrix[contract][wallet].isInit {
		return
	}

	nfts, err := m.adapter.GetNFTBalance(int64(block), (*common.Address)(&wallet), (*common.Address)(&contract))
	if err != nil {
		fmt.Println("ERROR: fillNFTMatrixCell: Failed to get NFTs balance")
	}
	if _, ok := m.nftMatrix[contract]; !ok {
		m.nftMatrix[contract] = make(map[Wallet]*NFTCell)
	}

	for _, nft := range nfts {
		m.nftMatrix[contract][wallet].value[nft] += 1
	}
	m.nftMatrix[contract][wallet].isInit = true

	m.saveAllToDB()

	if wg != nil {
		wg.Done()
	}
}

func (m *Matrix) fillTokenMatrixCell(block uint64, contract Contract, wallet Wallet, wg *sync.WaitGroup) {
	//m.mu.Lock()
	//defer m.mu.Unlock()

	if m.tokenMatrix[contract][wallet].isInit {
		return
	}

	fmt.Println("fillTokenMatrixCell")
	b, err := m.adapter.GetBalance((*common.Address)(&wallet), (*common.Address)(&contract), block)
	if err != nil {
		fmt.Println("ERROR: fillTokenMatrixCell: Failed to get token balance")
	}

	if b == nil {
		b = big.NewInt(0)
	}

	m.tokenMatrix[contract][wallet].value.Add(m.tokenMatrix[contract][wallet].value, b)

	m.tokenMatrix[contract][wallet].isInit = true

	m.saveAllToDB()

	if wg != nil {
		wg.Done()
	}
}

func (m *Matrix) fillMissingData(wgMain *sync.WaitGroup) {

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

	wg.Wait()

	if wgMain != nil {
		wgMain.Done()
	}
}

func (m *Matrix) getNFTEntriesFromCell(cell *NFTCell, add []*entry.NFT, remove []*entry.NFT, l *TransferNFTLog, blockChainID umid.UMID) {
	for id, v := range cell.value {
		if v == 1 {
			add = append(add, &entry.NFT{
				WalletID:     l.From.Bytes(),
				BlockchainID: blockChainID,
				ObjectID:     id,
				ContractID:   l.Contract.Bytes(),
			})
		} else {
			remove = append(remove, &entry.NFT{
				WalletID:     l.From.Bytes(),
				BlockchainID: blockChainID,
				ObjectID:     id,
				ContractID:   l.Contract.Bytes(),
			})
		}
	}
}

func (m *Matrix) ProcessLogs(blockNumber uint64, logs []any) {

	walletEntries := make([]*entry.Wallet, 0)
	balanceEntries := make([]*entry.Balance, 0)
	nftEntriesAdd := make([]*entry.NFT, 0)
	nftEntriesRemove := make([]*entry.NFT, 0)
	stakeEntries := make([]*entry.Stake, 0)
	contractsEntries := make([]*entry.Contract, 0)

	_ = walletEntries
	//_ = stakeEntries
	_ = contractsEntries

	updatedWallets := make(map[common.Address]bool)
	updatedContracts := make(map[common.Address]bool)

	blockChainID, _, _ := m.adapter.GetInfo()

	for _, log := range logs {

		switch l := log.(type) {
		case *TransferERC20Log:
			if _, ok := m.wallets[(Address)(l.From)]; ok {
				cell := m.tokenMatrix[(Contract)(l.Contract)][(Wallet)(l.From)]
				cell.value.Sub(cell.value, l.Value)
				updatedWallets[l.From] = true
				updatedContracts[l.Contract] = true
				balanceEntries = append(balanceEntries, &entry.Balance{
					WalletID:                 l.From.Bytes(),
					ContractID:               l.Contract.Bytes(),
					BlockchainID:             blockChainID,
					LastProcessedBlockNumber: m.blockNumber,
					Balance:                  (*entry.BigInt)(cell.value),
				})
			}

			if _, ok := m.wallets[(Address)(l.To)]; ok {
				cell := m.tokenMatrix[(Contract)(l.Contract)][(Wallet)(l.To)]
				cell.value.Add(cell.value, l.Value)
				updatedWallets[l.From] = true
				updatedContracts[l.Contract] = true
				balanceEntries = append(balanceEntries, &entry.Balance{
					WalletID:                 l.To.Bytes(),
					ContractID:               l.Contract.Bytes(),
					BlockchainID:             blockChainID,
					LastProcessedBlockNumber: m.blockNumber,
					Balance:                  (*entry.BigInt)(cell.value),
				})
			}

		case *TransferNFTLog:
			if _, ok := m.wallets[(Address)(l.From)]; ok {
				cell := m.nftMatrix[(Contract)(l.Contract)][(Wallet)(l.From)]
				cell.value[l.TokenID] -= 1
				updatedWallets[l.From] = true
				updatedContracts[l.Contract] = true
				m.getNFTEntriesFromCell(cell, nftEntriesAdd, nftEntriesRemove, l, blockChainID)
			}

			if _, ok := m.wallets[(Address)(l.To)]; ok {
				cell := m.nftMatrix[(Contract)(l.Contract)][(Wallet)(l.From)]
				cell.value[l.TokenID] += 1
				updatedWallets[l.From] = true
				updatedContracts[l.Contract] = true

				m.getNFTEntriesFromCell(cell, nftEntriesAdd, nftEntriesRemove, l, blockChainID)
			}

		case *StakeLog:
			if _, ok := m.wallets[(Address)(l.UserWallet)]; ok {

				cell := m.stakeMatrix[(Contract)(l.Contract)][(Wallet)(l.UserWallet)]
				if _, ok := cell.Stakes[l.OdysseyID]; !ok {
					cell.Stakes[l.OdysseyID] = &Stake{
						TotalAmount:    big.NewInt(0),
						TotalDADAmount: big.NewInt(0),
						TotalMOMAmount: big.NewInt(0),
					}
				}
				if l.TokenType == 0 {
					cell.Stakes[l.OdysseyID].TotalAmount = l.TotalStaked
					cell.Stakes[l.OdysseyID].TotalMOMAmount.Add(cell.Stakes[l.OdysseyID].TotalMOMAmount, l.AmountStaked)
				} else {
					cell.Stakes[l.OdysseyID].TotalAmount = l.TotalStaked
					cell.Stakes[l.OdysseyID].TotalDADAmount.Add(cell.Stakes[l.OdysseyID].TotalDADAmount, l.AmountStaked)
				}
				updatedWallets[l.UserWallet] = true
				updatedContracts[l.Contract] = true

				stakeEntries = append(stakeEntries, &entry.Stake{
					WalletID:     l.UserWallet.Bytes(),
					BlockchainID: blockChainID,
					ObjectID:     l.OdysseyID,
					LastComment:  "",
					Amount:       (*entry.BigInt)(l.TotalStaked),
				})
			}

		case *UnstakeLog:
			if _, ok := m.wallets[(Address)(l.UserWallet)]; ok {
				cell := m.stakeMatrix[(Contract)(l.Contract)][(Wallet)(l.UserWallet)]
				cell.Stakes[l.OdysseyID].TotalAmount = l.TotalStaked

				updatedWallets[l.UserWallet] = true
				updatedContracts[l.Contract] = true

				stakeEntries = append(stakeEntries, &entry.Stake{
					WalletID:     l.UserWallet.Bytes(),
					BlockchainID: blockChainID,
					ObjectID:     l.OdysseyID,
					LastComment:  "",
					Amount:       (*entry.BigInt)(l.TotalStaked),
				})
			}
		case *RestakeLog:
			//todo
		}

	}

	for w := range updatedWallets {
		walletEntries = append(walletEntries, &entry.Wallet{
			WalletID:     w.Bytes(),
			BlockchainID: blockChainID,
		})
	}

	for c := range updatedContracts {
		contractsEntries = append(contractsEntries, &entry.Contract{
			ContractID: c.Bytes(),
			Name:       "",
		})
	}

	m.blockNumber = blockNumber

	err := m.saveUpdateToDB(walletEntries, contractsEntries, nftEntriesAdd, nftEntriesRemove, balanceEntries, stakeEntries)
	if err != nil {
		err = errors.WithMessage(err, "failed to saveUpdateToDB")
		fmt.Println(err)
	}

}

func (m *Matrix) saveUpdateToDB(wallets []*entry.Wallet,
	contracts []*entry.Contract,
	addNFTs []*entry.NFT,
	removeNFTs []*entry.NFT,
	balances []*entry.Balance,
	stakes []*entry.Stake) error {

	blockchainUMID, name, rpcURL := m.adapter.GetInfo()

	tx, err := m.db.BeginTx(context.TODO(), pgx.TxOptions{})
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
		LastProcessedBlockNumber: m.blockNumber,
		BlockchainName:           (string)(name),
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
		_, err = tx.Exec(context.Background(), sql, w.WalletID, blockchainUMID)
		if err != nil {
			err = errors.WithMessage(err, "failed to insert wallet to DB")
			return err
		}
	}

	sql = `INSERT INTO contract (contract_id, name)
			VALUES ($1, $2)
			ON CONFLICT (contract_id) DO NOTHING`
	for _, c := range contracts {
		_, err = tx.Exec(context.TODO(), sql, c.ContractID, "")
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

	sql = `INSERT INTO stake (wallet_id, blockchain_id, object_id, amount, last_comment, updated_at, created_at)
			VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
			ON CONFLICT (blockchain_id, wallet_id, object_id)
				DO UPDATE SET updated_at   = NOW(),
				              last_comment = $5,
							  amount       = $4`

	for _, s := range stakes {
		_, err = tx.Exec(context.TODO(), sql,
			s.WalletID, blockchainUMID, s.ObjectID, s.Amount, s.LastComment)
		if err != nil {
			err = errors.WithMessage(err, "failed to insert stakes to DB")
			return err
		}
	}

	sql = `INSERT INTO nft (wallet_id, blockchain_id, object_id, contract_id, created_at, updated_at)	
			VALUES ($1, $2, $3, $4, NOW(), NOW())
			ON CONFLICT (wallet_id, contract_id, blockchain_id, object_id) DO UPDATE SET updated_at=NOW()`

	deleteSQL := `DELETE FROM nft WHERE object_id = $1`

	for _, nft := range removeNFTs {
		_, err = tx.Exec(context.TODO(), deleteSQL, nft.ObjectID)
		if err != nil {
			err = errors.WithMessage(err, "failed to delete NFT from DB")
			return err
		}
	}

	for _, nft := range addNFTs {
		_, err = tx.Exec(context.TODO(), sql, nft.WalletID, blockchainUMID, nft.ObjectID, nft.ContractID)
		if err != nil {
			err = errors.WithMessage(err, "failed to insert NFT to DB")
			return err
		}
	}

	return nil

}

func (m *Matrix) saveAllToDB() {
	//m.mu.Lock()
	//defer m.mu.Unlock()

	contracts := make([]*entry.Contract, 0)
	wallets := make([]*entry.Wallet, 0)
	addNFTs := make([]*entry.NFT, 0)
	removeNFTs := make([]*entry.NFT, 0)
	balances := make([]*entry.Balance, 0)
	stakes := make([]*entry.Stake, 0)

	_ = contracts
	_ = wallets
	_ = addNFTs
	_ = removeNFTs
	_ = balances
	_ = stakes

	blockChainID, _, _ := m.adapter.GetInfo()

	for c := range m.contracts {
		contracts = append(contracts, &entry.Contract{
			ContractID: (common.Address)(c).Bytes(),
			Name:       "",
		})
	}

	for w := range m.wallets {
		wallets = append(wallets, &entry.Wallet{
			WalletID:     (common.Address)(w).Bytes(),
			BlockchainID: blockChainID,
		})
	}

	for c, val := range m.nftMatrix {
		for w, cell := range val {
			if !cell.isInit {
				continue
			}
			for id, v := range cell.value {
				e := &entry.NFT{
					WalletID:     (common.Address)(w).Bytes(),
					BlockchainID: blockChainID,
					ObjectID:     id,
					ContractID:   (common.Address)(c).Bytes(),
				}
				if v == 1 {
					addNFTs = append(addNFTs, e)
				} else {
					removeNFTs = append(removeNFTs, e)
				}
			}
		}
	}

	for c, val := range m.tokenMatrix {
		for w, cell := range val {
			if !cell.isInit {
				continue
			}
			balances = append(balances, &entry.Balance{
				WalletID:                 (common.Address)(w).Bytes(),
				ContractID:               (common.Address)(c).Bytes(),
				BlockchainID:             blockChainID,
				LastProcessedBlockNumber: m.blockNumber,
				Balance:                  (*entry.BigInt)(cell.value),
			})
		}
	}

	for _, val := range m.stakeMatrix {
		for w, cell := range val {
			if !cell.isInit {
				continue
			}
			for id, stake := range cell.Stakes {
				stakes = append(stakes, &entry.Stake{
					WalletID:     (common.Address)(w).Bytes(),
					BlockchainID: blockChainID,
					ObjectID:     id,
					LastComment:  "",
					Amount:       (*entry.BigInt)(stake.TotalAmount),
				})
			}
		}
	}

	err := m.saveUpdateToDB(wallets, contracts, addNFTs, removeNFTs, balances, stakes)
	if err != nil {
		err = errors.WithMessage(err, "failed to saveUpdateToDB")
		fmt.Println(err)
	}
}

func (m *Matrix) AddWallet(wallet Address) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.addWallet(wallet, nil)
}

func (m *Matrix) addWallet(wallet Address, wg *sync.WaitGroup) error {
	if _, ok := m.wallets[wallet]; ok {
		// Wallet already subscribed
		if wg != nil {
			wg.Done()
		}
		return nil
	}

	w := (Wallet)(wallet)
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

	go m.fillMissingData(wg)

	return nil
}

func (m *Matrix) AddNFTContract(contract Address) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.addNFTContract(contract, nil)
}

func (m *Matrix) addNFTContract(contract Address, wg *sync.WaitGroup) error {
	if _, ok := m.contracts[contract]; ok {
		// Contract already subscribed
		return nil
	}

	c := (Contract)(contract)
	m.nftMatrix[c] = make(map[Wallet]*NFTCell)

	for wallet := range m.wallets {
		w := (Wallet)(wallet)
		// All new cells require initial fill
		m.nftMatrix[c][w] = &NFTCell{
			isInit: false,
			value:  make(map[umid.UMID]int8),
		}
	}

	m.contracts[contract] = true

	go m.fillMissingData(wg)

	return nil
}

func (m *Matrix) AddTokenContract(contract Address) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.addTokenContract(contract, nil)
}

func (m *Matrix) addTokenContract(contract Address, wg *sync.WaitGroup) error {
	if _, ok := m.contracts[contract]; ok {
		// Contract already subscribed
		if wg != nil {
			wg.Done()
		}
		return nil
	}

	c := (Contract)(contract)
	m.tokenMatrix[c] = make(map[Wallet]*TokenCell)

	for wallet := range m.wallets {
		w := (Wallet)(wallet)
		m.tokenMatrix[c][w] = &TokenCell{
			isInit: false,
			value:  big.NewInt(0),
		}
	}

	m.contracts[contract] = true

	go m.fillMissingData(wg)

	return nil
}

func (m *Matrix) AddStakeContract(contract Address) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.contracts[contract]; ok {
		// Contract already subscribed
		return nil
	}

	c := (Contract)(contract)

	m.stakeMatrix[c] = make(map[Wallet]*StakeCell)

	for wallet := range m.wallets {
		w := (Wallet)(wallet)
		m.stakeMatrix[c][w] = &StakeCell{
			isInit: false,
			Stakes: make(map[umid.UMID]*Stake),
		}
	}

	m.contracts[contract] = true

	go m.fillMissingData(nil)
	return nil
}

func (m *Matrix) AddTokenListener(contract Address, event string, listener TokenListener) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.tokenListeners[(Contract)(contract)]; !ok {
		m.tokenListeners[(Contract)(contract)] = make(map[string]TokenListener)
	}

	m.tokenListeners[(Contract)(contract)][event] = listener

	if _, ok := m.contracts[contract]; !ok {
		go func() {
			err := m.AddTokenContract(contract)
			if err != nil {
				fmt.Println(errors.WithMessage(err, "failed to AddTokenContract"))
			}
		}()
	}

	return nil
}

func (m *Matrix) AddNFTListener(contract Address, event string, listener NFTListener) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.nftListeners[(Contract)(contract)]; !ok {
		m.nftListeners[(Contract)(contract)] = make(map[string]NFTListener)
	}

	m.nftListeners[(Contract)(contract)][event] = listener

	if _, ok := m.contracts[contract]; !ok {
		go func() {
			err := m.AddNFTContract(contract)
			if err != nil {
				fmt.Println(errors.WithMessage(err, "failed to AddNFTContract"))
			}
		}()
	}

	return nil
}

func (m *Matrix) Display() {
	fmt.Println("")
	fmt.Println("Token Matrix:")
	for contract, value := range m.tokenMatrix {
		for wallet, v := range value {
			fmt.Printf("%v %v %v \n", (common.Address)(contract).Hex(), (common.Address)(wallet).Hex(), v.value.String())
		}
	}

	fmt.Println("NFT Matrix:")
	for contract, value := range m.nftMatrix {
		for wallet, v := range value {
			fmt.Printf("%v %v %v \n", (common.Address)(contract).Hex(), (common.Address)(wallet).Hex(), v.value)
		}
	}

	if len(m.stakeMatrix) == 1 {
		var stakeContract Contract
		for c, _ := range m.stakeMatrix {
			stakeContract = c
		}
		fmt.Println("STAKE Matrix:")
		fmt.Println("Contract:", (common.Address)(stakeContract).Hex())
		for wallet, val := range m.stakeMatrix[stakeContract] {
			for id, v := range val.Stakes {
				fmt.Println((common.Address)(wallet).Hex(), id.String(), val.isInit, v.TotalAmount, v.TotalMOMAmount, v.TotalDADAmount)
			}
		}
	}
}

func (m *Matrix) fastForward() {
	m.mu.Lock()
	defer m.mu.Unlock()

	lastBlockNumber, err := m.adapter.GetLastBlockNumber()
	fmt.Printf("Fast Forward. From: %d to: %d\n", m.blockNumber, lastBlockNumber)
	if err != nil {
		fmt.Println(err)
		return
	}

	if m.blockNumber >= lastBlockNumber {
		// Matrix already processed latest BC block
		return
	}

	if len(m.contracts) == 0 {
		return
	}

	fmt.Println("Doing Fast Forward")

	logs, err := m.adapter.GetLogs(int64(m.blockNumber)+1, int64(lastBlockNumber), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	m.ProcessLogs(lastBlockNumber, logs)
}

func (m *Matrix) GetWalletTokenData(contract Address, wallet Address) (*TokenData, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	err := m.addTokenContract(contract, wg)
	if err != nil {
		return nil, err
	}
	wg.Wait()

	wg.Add(1)
	err = m.addWallet(wallet, wg)
	if err != nil {
		return nil, err
	}

	wg.Wait()

	cell := m.tokenMatrix[(Contract)(contract)][(Wallet)(wallet)]

	return &TokenData{
		Wallet:      &wallet,
		Contract:    &contract,
		TotalAmount: cell.value,
	}, nil
}

func (m *Matrix) GetWalletNFTData(contract Address, wallet Address) (*NFTData, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	err := m.addNFTContract(contract, wg)
	if err != nil {
		return nil, err
	}
	wg.Wait()

	wg.Add(1)
	err = m.addWallet(wallet, wg)
	if err != nil {
		return nil, err
	}

	wg.Wait()

	cell := m.nftMatrix[(Contract)(contract)][(Wallet)(wallet)]

	data := &NFTData{
		Wallet:   &wallet,
		Contract: &contract,
		TokenIDs: make([]umid.UMID, 0),
	}

	for id, v := range cell.value {
		if v == 1 {
			data.TokenIDs = append(data.TokenIDs, id)
		}
	}

	return data, nil
}