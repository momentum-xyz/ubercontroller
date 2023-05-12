package harvester

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/config"
)

func SubscribeAllWallets(ctx context.Context, harv *Harvester, cfg *config.Config, pool *pgxpool.Pool) error {
	sql := `SELECT value -> 'wallet' AS wallets
						FROM user_attribute
						WHERE true
						  AND plugin_id = '86DC3AE7-9F3D-42CB-85A3-A71ABC3C3CB8'
						  AND attribute_name = 'wallet'`

	rows, err := pool.Query(ctx, sql)
	if err != nil {
		return err
	}

	allWallets := make([]string, 0)

	for rows.Next() {
		var wallets []string

		if err := rows.Scan(&wallets); err != nil {
			fmt.Println(allWallets)
			return errors.WithMessage(err, "failed to scan rows from user_attribute table")
		}

		allWallets = append(allWallets, wallets...)
	}

	handlerVar := handler
	handlerPointer := &handlerVar

	for _, w := range allWallets {
		if len(w) != 42 {
			// Skip Polkadot address in Hex format
			continue
		}

		if w[0:2] != "0x" {
			// Skip Polkadot addresses
			continue
		}

		err := harv.SubscribeForWalletAndContract(ArbitrumNova, w, cfg.Arbitrum.MOMTokenAddress, handlerPointer)
		if err != nil {
			return errors.WithMessage(err, "failed to subscribe wallet/contract")
		}
	}

	return nil
}

func handler(p any) {

}
