package harvester3

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/sasha-s/go-deadlock"
)

type SubscribeQueue struct {
	mu           deadlock.RWMutex
	wallets      map[common.Address]bool
	contracts    map[common.Address]bool
	updates      chan any
	loadedFromDB map[common.Address]map[common.Address]bool
}

func NewSubscribeQueue(updates chan any) *SubscribeQueue {
	return &SubscribeQueue{
		mu:           deadlock.RWMutex{},
		wallets:      make(map[common.Address]bool),
		contracts:    make(map[common.Address]bool),
		updates:      updates,
		loadedFromDB: make(map[common.Address]map[common.Address]bool),
	}
}

func (q *SubscribeQueue) AddWallet(wallet common.Address) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if _, ok := q.wallets[wallet]; ok {
		// Wallet already subscribed
		return nil
	}

	q.wallets[wallet] = true

	for c, _ := range q.contracts {
		if !q.isLoadedFromDB(c, wallet) {

			q.updates <- QueueInit{
				contract: c,
				wallet:   wallet,
			}
		}
	}

	return nil
}

func (q *SubscribeQueue) AddContract(contract common.Address) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if _, ok := q.contracts[contract]; ok {
		// contract already subscribed
		return nil
	}

	q.contracts[contract] = true

	for w, _ := range q.wallets {
		if !q.isLoadedFromDB(contract, w) {
			q.updates <- QueueInit{
				contract: contract,
				wallet:   w,
			}
		}
	}

	return nil
}

func (q *SubscribeQueue) MarkAsLoadedFromDB(contract, wallet common.Address) {
	q.mu.Lock()
	defer q.mu.Unlock()
	_, ok := q.loadedFromDB[contract]
	if !ok {
		q.loadedFromDB[contract] = make(map[common.Address]bool)
	}
	q.loadedFromDB[contract][wallet] = true
}

func (q *SubscribeQueue) isLoadedFromDB(contract, wallet common.Address) bool {
	//fmt.Println("isLoadedFromDB", contract.Hex(), wallet.Hex())
	_, ok := q.loadedFromDB[contract]
	if ok {
		_, ok := q.loadedFromDB[contract][wallet]
		if ok {
			return true
		}
	}

	return false
}
