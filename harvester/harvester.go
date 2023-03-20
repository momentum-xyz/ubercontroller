package harvester

import (
	"math/big"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type Harvester struct {
	clients *Callbacks
	db      *pgxpool.Pool
	bc      map[string]*BlockChain
}

func (h *Harvester) SubscribeForWallet(bcType string, wallet, callback Callback) {
	//TODO implement me
	panic("implement me")
}

func (h *Harvester) SubscribeForWalletAndContract(bcType string, wallet []byte, contract []byte, callback Callback) error {
	bc, ok := h.bc[bcType]
	if !ok {
		return errors.New("failed to find blockchain:" + bcType)
	}

	h.clients.Add(bcType, BalanceChange, callback)

	bc.SubscribeForWalletAndContract(wallet, contract)
	bc.SaveBalancesToDB()
	bc.LoadBalancesFromDB()

	return nil
}

type Address []byte

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

func (h *Harvester) OnNewBlock(bcType string, block *BCBlock) {
	//fmt.Printf("On new block: %+v %+v \n", block.Hash, block.Number)
	h.clients.Trigger(bcType, NewBlock, block)
}

func (h *Harvester) RegisterAdapter(umid umid.UMID, bcType string, rpcURL string, adapter Adapter) error {
	h.bc[bcType] = NewBlockchain(h.db, adapter, umid, bcType, rpcURL, h.updateHook)
	if err := h.bc[bcType].LoadFromDB(); err != nil {
		return errors.WithMessage(err, "failed to load from DB")
	}

	return nil
}

func (h *Harvester) updateHook(bcType string, wallet string, contract string, blockNumber uint64, balance *big.Int) {
	h.clients.Trigger(bcType, BalanceChange, []any{wallet, contract, balance})
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
