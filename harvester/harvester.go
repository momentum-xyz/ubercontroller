package harvester

type Harvester struct {
	clients  *Callbacks
	Adapters *Callbacks
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
}

func NewHarvester() *Harvester {
	return &Harvester{
		clients: NewCallbacks(),
	}
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
