package engine

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Rafael24595/go-log/log/record"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestEngine_Concurrency(t *testing.T) {
	var buf bytes.Buffer

	store := record.NewMemory()

	eng, _ := NewEngine(
		"InternalTest",
		WithContext(
			t.Context(),
		),
		WithBufferSize(10),
		WithRecordStore(store),
		WithWriteAction(
			func(r record.Record) error {
				buf.Write([]byte(r.Message))
				return nil
			},
		),
	)

	var wg sync.WaitGroup

	goroutines := 50
	logsPerRoutine := 100
	totalExpected := goroutines * logsPerRoutine

	for g := range goroutines {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := range logsPerRoutine {
				eng.Message(fmt.Sprintf("Routine %d - Log %d", id, i))
			}
		}(g)
	}

	wg.Wait()

	err := eng.Close()
	assert.Nil(t, err)

	records := store.All()
	assert.Len(t, totalExpected, records)
}

func TestEngine_RecordsImmutability(t *testing.T) {
	store := record.NewMemory()
	
	eng, _ := NewEngine(
		"InternalTest",
		WithContext(
			t.Context(),
		),
		WithBufferSize(10),
		WithRecordStore(store),
	)

	eng.Message("Log 1")
	time.Sleep(1 * time.Millisecond)

	history := store.All()

	eng.Message("Log 2")
	time.Sleep(1 * time.Millisecond)

	assert.Len(t, 1, history)
}

func TestEngine_CloseIdempotency(t *testing.T) {
	counter := 0

	eng, _ := NewEngine(
		"InternalTest",
		WithContext(
			t.Context(),
		),
		WithBufferSize(10),
		WithCloseAction(
			func() error {
				counter++
				return nil
			},
		),
	)

	assert.Nil(t, eng.Close())
	assert.Nil(t, eng.Close())
	assert.Nil(t, eng.Close())

	assert.Equal(t, 1, counter)
}

func TestEngine_ShutdownMechanism(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	eng, _ := NewEngine(
		"InternalTest",
		WithContext(ctx),
		WithBufferSize(10),
	)

	cancel()

	assert.WillClose(t, eng.done, 100*time.Millisecond)
}
