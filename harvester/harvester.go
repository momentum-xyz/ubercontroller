package harvester

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type Harvester struct {
	clients  *Callbacks
	Adapters *Callbacks
	db       *pgxpool.Pool
}

func (h *Harvester) SubscribeForWallet(bcType BCType, wallet, callback Callback) {
	//TODO implement me
	panic("implement me")
}

func (h *Harvester) SubscribeForWalletAndContract(bcType BCType, wallet, callback Callback) {
	//TODO implement me
	panic("implement me")
}

type BCBlock struct {
	Hash   string
	Number uint64
}

type BCAdapterAPI interface {
	RegisterBCAdapter(bcType BCType, bcAdapter BCAdapter)
	OnNewBlock(bcType BCType, block *BCBlock)
}

type BCAdapter interface {
}

type HarvesterAPI interface {
	OnBalanceChange()
	Subscribe(bcType BCType, eventName HarvesterEvent, callback Callback)
	Unsubscribe(bcType BCType, eventName HarvesterEvent, callback Callback)
	SubscribeForWallet(bcType BCType, wallet, callback Callback)
	SubscribeForWalletAndContract(bcType BCType, wallet, callback Callback)
}

func NewHarvester(db *pgxpool.Pool) *Harvester {
	return &Harvester{
		clients: NewCallbacks(),
		db:      db,
	}
}

func (h *Harvester) Init() {

}

func (h *Harvester) OnBalanceChange() {

}

func (h *Harvester) OnNewBlock(bcType BCType, block *BCBlock) {
	//fmt.Printf("On new block: %+v %+v \n", block.Hash, block.Number)
	h.clients.Trigger(bcType, NewBlock, block)
}

func (h *Harvester) RegisterBCAdapter(bcType BCType, bcAdapter BCAdapter) {
}

func (h *Harvester) Run() error {
	return nil
}

func (h *Harvester) Subscribe(bcType BCType, eventName HarvesterEvent, callback Callback) {
	h.clients.Add(bcType, eventName, callback)
}

func (h *Harvester) Unsubscribe(bcType BCType, eventName HarvesterEvent, callback Callback) {
	h.clients.Remove(bcType, eventName, callback)
}
