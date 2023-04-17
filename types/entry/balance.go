package entry

import (
	"database/sql/driver"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type Balance struct {
	WalletID                 []byte    `db:"wallet_id"`
	ContractID               []byte    `db:"contract_id"`
	BlockchainID             umid.UMID `db:"blockchain_id"`
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
	// TODO Find better way

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

		var s string

		r, _ := regexp.Compile("^[0-9]*e[0-9]*$")
		match := r.MatchString(t)

		if match {
			// If match string example: 265e17
			ss := strings.Split(t, "e")
			zerosCount, err := strconv.ParseInt(ss[1], 10, 8)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("failed to strconv.ParseInt: %v", ss[1])
			}
			s = ss[0] + strings.Repeat("0", int(zerosCount))
		} else {
			s = t
		}
		_, ok := (*big.Int)(b).SetString(s, 10)
		if !ok {
			return fmt.Errorf("failed to load value to string: %v", value)
		}
	default:
		return fmt.Errorf("Could not scan type %T into BigInt", t)
	}

	return nil
}
