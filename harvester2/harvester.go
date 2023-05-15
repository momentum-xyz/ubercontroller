package harvester2

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"sync"
)

type Harvester2 struct {
	db *pgxpool.Pool
	bc map[BCType]*Matrix
	mu sync.Mutex
}

func NewHarvester(db *pgxpool.Pool) *Harvester2 {
	return &Harvester2{
		db: db,
		bc: make(map[BCType]*Matrix),
	}
}

func (h *Harvester2) RegisterAdapter(adapter Adapter) error {
	_, bcType, _ := adapter.GetInfo()

	h.bc[bcType] = NewMatrix(h.db, adapter)
	h.bc[bcType].Run()

	return nil
}

func (h *Harvester2) AddWallet(bcType BCType, wallet *Address) error {
	if _, ok := h.bc[bcType]; !ok {
		return errors.New("failed to add wallet, adapter not registered")
	}

	return h.bc[bcType].AddWallet(wallet)
}

func (h *Harvester2) RemoveWallet(bcType BCType, wallet *Address) error {
	if _, ok := h.bc[bcType]; !ok {
		return errors.New("failed to remove wallet, adapter not registered")
	}

	return nil
}

func (h *Harvester2) AddNFTContract(bcType BCType, contract *Address) error {
	if _, ok := h.bc[bcType]; !ok {
		return errors.New("failed to add nft contract, adapter not registered")
	}

	return h.bc[bcType].AddNFTContract(contract)
}

func (h *Harvester2) RemoveNFTContract(bcType BCType, contract *Address) error {
	if _, ok := h.bc[bcType]; !ok {
		return errors.New("failed to remove nft contract, adapter not registered")
	}

	return nil
}

func (h *Harvester2) AddTokenContract(bcType BCType, contract *Address) error {
	if _, ok := h.bc[bcType]; !ok {
		return errors.New("failed to add token contract, adapter not registered")
	}

	return nil
}
func (h *Harvester2) RemoveTokenContract(bcType BCType, contract *Address) error {
	if _, ok := h.bc[bcType]; !ok {
		return errors.New("failed to remove token contract, adapter not registered")
	}

	return nil
}

func (h *Harvester2) AddStakeContract(bcType BCType, contract *Address) error {
	if _, ok := h.bc[bcType]; !ok {
		return errors.New("failed to add stake contract, adapter not registered")
	}

	return nil
}
func (h *Harvester2) RemoveStakeContract(bcType BCType, contract *Address) error {
	if _, ok := h.bc[bcType]; !ok {
		return errors.New("failed to remove stake contract, adapter not registered")
	}

	return nil
}

func (h *Harvester2) AddTokenListener(bcType BCType, contract *Address, listener TokenListener) error {
	if _, ok := h.bc[bcType]; !ok {
		return errors.New("failed to add token listener, adapter not registered")
	}

	return nil
}
func (h *Harvester2) AddNFTListener(bcType BCType, contract *Address, listener NFTListener) error {
	if _, ok := h.bc[bcType]; !ok {
		return errors.New("failed to add nft listener, adapter not registered")
	}

	return nil
}
func (h *Harvester2) AddStakeListener(bcType BCType, contract *Address, listener StakeListener) error {
	if _, ok := h.bc[bcType]; !ok {
		return errors.New("failed to add stake listener, adapter not registered")
	}

	return nil
}
