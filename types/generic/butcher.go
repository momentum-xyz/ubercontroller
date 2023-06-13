package generic

import (
	"context"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type BatchFn[T any] func(batch []T) error
type BatchItemFn[T any] func(item T) error

type Butcher[T any] struct {
	data []T
}

func NewButcher[T any](data []T) *Butcher[T] {
	return &Butcher[T]{
		data: data,
	}
}

func (b *Butcher[T]) HandleItems(ctx context.Context, batchSize int, itemHandler BatchItemFn[T]) error {
	handler := func(batch []T) error {
		group, _ := errgroup.WithContext(ctx)
		for i := range batch {
			item := batch[i]

			group.Go(func() error {
				if err := itemHandler(item); err != nil {
					return errors.WithMessage(err, "failed to handle item")
				}
				return nil
			})
		}
		return group.Wait()
	}

	return b.HandleBatchesSync(batchSize, handler)
}

func (b *Butcher[T]) HandleBatchesSync(batchSize int, batchHandler BatchFn[T]) error {
	data := b.data
	for len(data) > 0 {
		batch := data
		if len(data) > batchSize {
			batch = data[:batchSize]
			data = data[batchSize:]
		} else {
			data = nil
		}

		if err := batchHandler(batch); err != nil {
			return errors.WithMessage(err, "failed to handle batch")
		}
	}

	return nil
}

func (b *Butcher[T]) HandleBatchesAsync(ctx context.Context, batchSize int, batchHandler BatchFn[T]) error {
	data := b.data
	group, _ := errgroup.WithContext(ctx)
	for len(data) > 0 {
		batch := data
		if len(data) > batchSize {
			batch = data[:batchSize]
			data = data[batchSize:]
		} else {
			data = nil
		}

		group.Go(func() error {
			if err := batchHandler(batch); err != nil {
				return errors.WithMessage(err, "failed to handle batch")
			}
			return nil
		})
	}
	return group.Wait()
}

func (b *Butcher[T]) Len() int {
	return len(b.data)
}
