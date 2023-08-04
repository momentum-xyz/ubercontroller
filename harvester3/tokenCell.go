package harvester3

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type TokenCell struct {
	contract common.Address
	wallet   common.Address

	block uint64
	value *big.Int

	isInit    bool
	initBlock uint64
	initValue *big.Int
}

func NewTokenCell(contract common.Address, wallet common.Address) *TokenCell {
	return &TokenCell{
		contract: contract,
		wallet:   wallet,
		isInit:   false,
		block:    0,
		value:    big.NewInt(0),
	}
}

func (c *TokenCell) initFromBC() {

}
