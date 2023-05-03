package nfts

import (
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

const (
	missingWorldsByWallet = `SELECT * FROM nft
		WHERE nft.wallet_id = $1
		AND NOT EXISTS (
			SELECT object_id from object WHERE object.object_id=nft.object_id)
	`
)

type DB struct {
	conn *pgxpool.Pool
}

func NewDB(conn *pgxpool.Pool) *DB {
	return &DB{
		conn: conn,
	}
}

func (db *DB) ListNewByWallet(ctx context.Context, wallet string) ([]*entry.NFT, error) {
	var nfts []*entry.NFT
	if err := pgxscan.Select(
		ctx, db.conn, &nfts,
		missingWorldsByWallet, utils.HexToAddress(wallet)); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return nfts, nil
}
