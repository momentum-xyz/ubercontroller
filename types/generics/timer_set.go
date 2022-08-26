package generics

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/momentum-xyz/controller/logger"
)

var log = logger.L()

type TimerFunc[T comparable] func(key T) error

type TimerSet[T comparable] struct {
	timers *SyncMap[T, Unique[context.CancelFunc]]
}

func NewTimerSet[T comparable]() *TimerSet[T] {
	return &TimerSet[T]{
		timers: NewSyncMap[T, Unique[context.CancelFunc]](),
	}
}

func (t *TimerSet[T]) Set(key T, delay time.Duration, fn TimerFunc[T]) {
	t.timers.Mu.Lock()
	defer t.timers.Mu.Unlock()

	stopFn, ok := t.timers.Data[key]
	if ok {
		stopFn.Value()()
	}

	ctx, cancel := context.WithTimeout(context.Background(), delay)
	stopFn = NewUnique(cancel)
	t.timers.Data[key] = stopFn

	go func() {
		defer func() {
			cancel()

			t.timers.Mu.Lock()
			defer t.timers.Mu.Unlock()

			if stopFn1, ok := t.timers.Data[key]; ok && stopFn1.Equals(stopFn) {
				delete(t.timers.Data, key)
			}
		}()

		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				if err := fn(key); err != nil {
					log.Warn(errors.WithMessagef(err, "TimerSet: Set: function call failed: %+v", key))
				}
			}
		}
	}()
}

func (t *TimerSet[T]) Stop(key T) {
	if stopFn, ok := t.timers.Load(key); ok {
		stopFn.Value()()
	}
}
