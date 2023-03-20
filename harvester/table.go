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
	fmt.Printf("Block: %d \n", block.Number)
	t.mu.Lock()
	for _, diff := range diffs {
		_, ok := t.data[diff.Token]
		if !ok {
			// No such contract
			continue
		}
		b, ok := t.data[diff.Token][diff.From]
		if ok {
			// From wallet found
			b.Sub(b, diff.Amount)
		}
		b, ok = t.data[diff.Token][diff.To]
		if ok {
			// To wallet found
			b.Add(b, diff.Amount)
		}
	}
	t.mu.Unlock()
	t.Display()
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

func (t *Table) Display() {
	for token, wallets := range t.data {
		for wallet, balance := range wallets {
			fmt.Printf("%+v %+v %+v \n", token, wallet, balance.String())
		}
	}
}
