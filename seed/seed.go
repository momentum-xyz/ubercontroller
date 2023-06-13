package seed

import (
	"fmt"

	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/momentum-xyz/ubercontroller/universe"
)

func Node(ctx types.NodeContext, node universe.Node, db database.DB) error {
	log := ctx.Logger()
	group, groupCtx := errgroup.WithContext(ctx)

	log.Debugln("Seeding node...")
	group.Go(
		func() error {
			return seedPlugins(groupCtx, node)
		},
	)

	group.Go(
		func() error {
			return seedAttributeType(groupCtx, node)
		},
	)

	group.Go(
		func() error {
			return seedNodeAttributes(groupCtx, node)
		},
	)

	group.Go(
		func() error {
			return seedAssets2d(groupCtx, node)
		},
	)

	group.Go(
		func() error {
			return seedAssets3d(groupCtx, node)
		},
	)

	group.Go(
		func() error {
			return seedUserTypes(groupCtx, node, db)
		},
	)

	group.Go(
		func() error {
			return seedMedia(groupCtx, node)
		},
	)

	if err := group.Wait(); err != nil {
		return errors.WithMessage(err, "failed to seed")
	}

	// Object Types must be seeded after assets
	if err := seedObjectTypes(node); err != nil {
		fmt.Println(err)
		return errors.WithMessage(err, "failed to seed object types")
	}

	if err := seedUsers(ctx, node, db); err != nil {
		return errors.WithMessage(err, "failed to seed users")
	}

	if err := node.SetOwnerID(umid.MustParse("00000000-0000-0000-0000-000000000003"), false); err != nil {
		return errors.WithMessage(err, "failed to set owner UMID")
	}

	if err := node.Save(); err != nil {
		return errors.WithMessage(err, "failed to save node")
	}

	return nil
	//return node.Save()
}
