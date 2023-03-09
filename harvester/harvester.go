package harvester

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type Harvester struct {
	clients  *Callbacks
	Adapters *Callbacks
	db       *pgxpool.Pool
	bc       map[string]*BlockChain
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
	RegisterBCAdapter(uuid uuid.UUID, bcType string, rpcURL string, bcAdapter BCAdapter) error
	OnNewBlock(bcType BCType, block *BCBlock)
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
		bc:      make(map[string]*BlockChain),
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

func (h *Harvester) RegisterBCAdapter(uuid uuid.UUID, bcType string, rpcURL string, adapter BCAdapter) error {
	h.bc[bcType] = NewBlockchain(h.db, adapter, uuid, bcType, rpcURL)
	if err := h.bc[bcType].LoadFromDB(); err != nil {
		return errors.WithMessage(err, "failed to load from DB")
	}

	return nil
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
