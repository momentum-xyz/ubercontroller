// WARNING:
// Most of the code was generated using ChatGPT
// and can have some invalid test cases

package generic

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/types"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	cfg := &config.Config{} //TODO: proper test config getter
	ctx, err := types.NewNodeContext(
		context.Background(),
		zap.NewExample().Sugar(),
		cfg,
	)
	if err != nil {
		fmt.Printf("Failed to create context: %s", err)
		os.Exit(1)
	}
	if err := Initialize(ctx); err != nil {
		fmt.Printf("Failed to initialize: %s", err)
		os.Exit(1)
	}

	code := m.Run()

	os.Exit(code)
}

func TestButcher_HandleItems(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	b := NewButcher(data)

	var sum int
	var mu sync.Mutex
	itemHandler := func(item int) error {
		mu.Lock()
		defer mu.Unlock()

		sum += item
		return nil
	}

	batchSize := 2
	err := b.HandleItems(batchSize, itemHandler)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := 15
	if sum != expected {
		t.Errorf("Expected handled data to be %d, got %d", expected, sum)
	}
}

func TestButcher_HandleItems_Error(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	b := NewButcher(data)

	itemHandler := func(item int) error {
		return errors.New("Error handling item")
	}

	batchSize := 2
	err := b.HandleItems(batchSize, itemHandler)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestButcher_HandleBatchesSync(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	b := NewButcher(data)

	var handledData [][]int
	batchHandler := func(batch []int) error {
		handledData = append(handledData, batch)
		return nil
	}

	batchSize := 2
	err := b.HandleBatchesSync(batchSize, batchHandler)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := [][]int{{1, 2}, {3, 4}, {5}}
	if !reflect.DeepEqual(expected, handledData) {
		t.Errorf("Expected handled data to be %v, got %v", expected, handledData)
	}
}

func TestButcher_HandleBatchesSync_Error(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	b := NewButcher(data)

	batchHandler := func(batch []int) error {
		return errors.New("Error handling batch")
	}

	batchSize := 2
	err := b.HandleBatchesSync(batchSize, batchHandler)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestButcher_HandleBatchesAsync(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	b := NewButcher(data)

	var sum int
	var mu sync.Mutex
	batchHandler := func(batch []int) error {
		mu.Lock()
		defer mu.Unlock()

		for _, n := range batch {
			sum += n
		}
		return nil
	}

	batchSize := 2
	err := b.HandleBatchesAsync(batchSize, batchHandler)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := 15
	if sum != expected {
		t.Errorf("Expected handled data to be %d, got %d", expected, sum)
	}
}

func TestButcher_HandleBatchesAsync_Error(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	b := NewButcher(data)

	batchHandler := func(batch []int) error {
		return errors.New("Error handling batch")
	}

	batchSize := 2
	err := b.HandleBatchesAsync(batchSize, batchHandler)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
