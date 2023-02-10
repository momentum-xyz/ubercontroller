package seed

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/momentum-xyz/ubercontroller/universe"
)

func Node(ctx context.Context, node universe.Node, db database.DB) error {
	group, _ := errgroup.WithContext(ctx)

	group.Go(func() error {
		return seedPlugins(node)
	})

	group.Go(func() error {
		return seedAttributeType(node)
	})

	group.Go(func() error {
		return seedNodeAttributes(node)
	})

	group.Go(func() error {
		return seedAssets2d(node)
	})

	group.Go(func() error {
		return seedAssets3d(node)
	})

	group.Go(func() error {
		return seedUserTypes(ctx, node, db)
	})

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

	if err := node.SetOwnerID(uuid.MustParse("00000000-0000-0000-0000-000000000003"), false); err != nil {
		return errors.WithMessage(err, "failed to set owner ID")
	}

	if err := node.Save(); err != nil {
		return errors.WithMessage(err, "failed to save node")
	}

	return nil
	//return node.Save()
}
