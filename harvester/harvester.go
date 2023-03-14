package harvester

import (
	"encoding/hex"
	"math/big"

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

func (h *Harvester) SubscribeForWallet(bcType string, wallet, callback Callback) {
	//TODO implement me
	panic("implement me")
}

func (h *Harvester) SubscribeForWalletAndContract(bcType string, wallet []byte, contract []byte, callback Callback) error {
	bc, ok := h.bc[bcType]
	if !ok {
		return errors.New("failed to find blockchain:" + bcType)
	}

	//event := hexutil.Encode(wallet) + "_" + hexutil.Encode(contract)
	h.clients.Add(bcType, BalanceChange, callback)

	bc.SubscribeForWalletAndContract(wallet, contract)
	bc.SaveBalancesToDB()
	bc.LoadBalancesFromDB()

	return nil
}

type BCBlock struct {
	Hash   string
	Number uint64
}

type BCAdapterAPI interface {
	RegisterBCAdapter(uuid uuid.UUID, bcType string, rpcURL string, bcAdapter BCAdapter) error
	OnNewBlock(bcType string, block *BCBlock)
}

type HarvesterAPI interface {
	OnBalanceChange()
	Subscribe(bcType string, eventName HarvesterEvent, callback Callback)
	Unsubscribe(bcType string, eventName HarvesterEvent, callback Callback)
	SubscribeForWallet(bcType string, wallet, callback Callback)
	SubscribeForWalletAndContract(bcType string, wallet []byte, contract []byte, callback Callback) error
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

func (h *Harvester) RegisterBCAdapter(uuid uuid.UUID, bcType string, rpcURL string, adapter BCAdapter) error {
	h.bc[bcType] = NewBlockchain(h.db, adapter, uuid, bcType, rpcURL, h.updateHook)
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

func (h *Harvester) Subscribe(bcType string, eventName HarvesterEvent, callback Callback) {
	h.clients.Add(bcType, eventName, callback)
}

func (h *Harvester) Unsubscribe(bcType string, eventName HarvesterEvent, callback Callback) {
	h.clients.Remove(bcType, eventName, callback)
}

func HexToAddress(s string) []byte {
	b, err := hex.DecodeString(s[2:])
	if err != nil {
		panic(err)
	}
	return b
}
