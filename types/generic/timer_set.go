package generic

import (
	"context"
	"fmt"
	"time"
)

type TimerFunc[T comparable] func(key T) error

type TimerSet[T comparable] struct {
	timers *SyncMap[T, Unique[context.CancelFunc]]
}

func NewTimerSet[T comparable]() *TimerSet[T] {
	return &TimerSet[T]{
		timers: NewSyncMap[T, Unique[context.CancelFunc]](0),
	}
}

func (t *TimerSet[T]) Set(ctx context.Context, key T, delay time.Duration, fn TimerFunc[T]) {
	t.timers.Mu.Lock()
	defer t.timers.Mu.Unlock()

	stopFn, ok := t.timers.Data[key]
	if ok {
		stopFn.Value()()
	}

	ctx, cancel := context.WithTimeout(ctx, delay)
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
					//TODO: return error
					fmt.Printf("TimerSet: Set: function call failed: %+v", key)
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

func (t *TimerSet[T]) StopAll() {
	t.timers.Mu.RLock()
	defer t.timers.Mu.RUnlock()

	for _, v := range t.timers.Data {
		v.Value()()
	}
}
