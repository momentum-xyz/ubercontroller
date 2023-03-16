package harvester

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/sasha-s/go-deadlock"
)

type Table struct {
	mu          deadlock.RWMutex
	blockNumber uint64
	data        map[string]map[string]*big.Int
	db          *pgxpool.Pool
	adapter     Adapter
}

func NewTable(db *pgxpool.Pool, adapter Adapter, listener func(p any)) *Table {
	return &Table{
		blockNumber: 0,
		data:        make(map[string]map[string]*big.Int),
		adapter:     adapter,
	}
}

func (t *Table) Run() {
	t.adapter.RegisterNewBlockListener(t.listener)
}

func (t *Table) listener(block *BCBlock, diffs []*BCDiff) {

}

func (t *Table) AddWalletContract(wallet string, contract string) {
	wallet = strings.ToLower(wallet)
	contract = strings.ToLower(contract)

	t.mu.Lock()
	_, ok := t.data[contract]
	if !ok {
		t.data[contract] = make(map[string]*big.Int)
	}
	_, ok = t.data[contract][wallet]
	if !ok {
		t.data[contract][wallet] = nil
		// Such wallet has not existed so need to get initial balance
		go t.syncBalance(wallet, contract)
	}
	t.mu.Unlock()
}

func (t *Table) syncBalance(wallet string, contract string) {
	wallet = strings.ToLower(wallet)
	contract = strings.ToLower(contract)

	t.mu.RLock()
	blockNumber, err := t.adapter.GetLastBlockNumber()
	t.mu.RUnlock()
	if err != nil {
		err = errors.WithMessage(err, "failed to get last block number")
		fmt.Println(err)
	}
	balance, err := t.adapter.GetBalance(wallet, contract, blockNumber)
	if err != nil {
		err = errors.WithMessage(err, "failed to get balance")
		fmt.Println(err)
	}
	t.mu.Lock()
	if t.blockNumber <= blockNumber {
		t.data[contract][wallet] = balance
		t.mu.Unlock()
	} else {
		t.mu.Unlock()
		t.syncBalance(wallet, contract)
	}
}
