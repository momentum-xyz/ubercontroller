package harvester

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type Harvester struct {
	clients *Callbacks
	db      *pgxpool.Pool
	bc      map[string]*Table
}

func (h *Harvester) SubscribeForWallet(bcType string, wallet, callback Callback) {
	//TODO implement me
	panic("implement me")
}

func (h *Harvester) SubscribeForWalletAndContract(bcType string, wallet string, contract string, callback Callback) error {
	table, ok := h.bc[bcType]
	if !ok {
		return errors.New("failed to find blockchain:" + bcType)
	}

	h.clients.Add(bcType, BalanceChange, callback)
	table.AddWalletContract(wallet, contract)

	return nil
}

type Address []byte

func NewHarvester(db *pgxpool.Pool) *Harvester {
	return &Harvester{
		clients: NewCallbacks(),
		db:      db,
		bc:      make(map[string]*Table),
	}
}

func (h *Harvester) Init() {

}

func (h *Harvester) OnBalanceChange() {

}

func (h *Harvester) OnNewBlock(bcType string, block *BCBlock) {
	//fmt.Printf("On new block: %+v %+v \n", block.Hash, block.Number)
	h.clients.Trigger(bcType, NewBlock, block)
}

func (h *Harvester) RegisterAdapter(adapter Adapter) error {
	_, bcType, _ := adapter.GetInfo()

	h.bc[bcType] = NewTable(h.db, adapter, h.updateHook)
	h.bc[bcType].Run()

	return nil
}

func (h *Harvester) updateHook(bcType string, updates []*UpdateEvent) {
	for update := range updates {
		h.clients.Trigger(bcType, BalanceChange, update)
	}
}

func (h *Harvester) Run() error {
	return nil
}

func (h *Harvester) Subscribe(bcType string, eventName Event, callback Callback) {
	h.clients.Add(bcType, eventName, callback)
}

func (h *Harvester) Unsubscribe(bcType string, eventName Event, callback Callback) {
	h.clients.Remove(bcType, eventName, callback)
}
