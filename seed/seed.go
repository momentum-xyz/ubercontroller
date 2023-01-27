package seed

import (
	"context"

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

	if err := group.Wait(); err != nil {
		return errors.WithMessage(err, "failed to seed plugins or attribute types")
	}

	return node.Save()
}
