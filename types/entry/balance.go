package entry

import (
	"database/sql/driver"
	"fmt"
	"math/big"
	"strings"

	"github.com/google/uuid"
)

type Balance struct {
	WalletID                 []byte    `db:"wallet_id"`
	ContractID               []byte    `db:"contract_id"`
	BlockchainID             uuid.UUID `db:"blockchain_id"`
	LastProcessedBlockNumber uint64    `db:"last_processed_block_number"`
	Balance                  *BigInt   `db:"balance"` //TODO should be big.Int
}

type BigInt big.Int

func (b *BigInt) Value() (driver.Value, error) {
	if b != nil {
		return (*big.Int)(b).String(), nil
	}
	return nil, nil
}

func (b *BigInt) Scan(value interface{}) error {
	if value == nil {
		b = nil
	}

	switch t := value.(type) {
	case []uint8:
		_, ok := (*big.Int)(b).SetString(string(value.([]uint8)), 10)
		if !ok {
			return fmt.Errorf("failed to load value to []uint8: %v", value)
		}
	case string:
		s := strings.Replace(t, "e0", "", 1) // TODO Find better way
		_, ok := (*big.Int)(b).SetString(s, 10)
		if !ok {
			return fmt.Errorf("failed to load value to string: %v", value)
		}
	default:
		return fmt.Errorf("Could not scan type %T into BigInt", t)
	}

	return nil
}
