package seed

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/momentum-xyz/ubercontroller/universe"
)

func Node(ctx context.Context, node universe.Node) error {
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
		return seedUserTypes(node)
	})

	if err := group.Wait(); err != nil {
		return errors.WithMessage(err, "failed to seed")
	}

	// Object Types must be seeded after assets
	if err := seedObjectTypes(node); err != nil {
		fmt.Println(err)
		return errors.WithMessage(err, "failed to seed object types")
	}

	return node.Save()
}
