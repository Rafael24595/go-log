package engine

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Rafael24595/go-log/log/model/record"

	"github.com/Rafael24595/go-assert/assert/test"

)

func TestEngine_Concurrency(t *testing.T) {
	var buf bytes.Buffer

	eng, _ := NewEngine(
		t.Context(),
		"InternalTest",
		10,
		func(r record.Record, _ []record.Record) error {
			buf.Write([]byte(r.Message))
			return nil
		},
		VoidCloseAction,
	)

	var wg sync.WaitGroup
	goroutines := 50
	logsPerRoutine := 100
	totalExpected := goroutines * logsPerRoutine

	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := range logsPerRoutine {
				eng.Message(fmt.Sprintf("Routine %d - Log %d", id, i))
			}
		}(g)
	}

	wg.Wait()
	records, _ := eng.Close()

	assert.Len(t, totalExpected, records)
}

func TestEngine_RecordsImmutability(t *testing.T) {
	eng, _ := NewEngine(
		t.Context(),
		"InternalTest",
		10,
		VoidWriteAction,
		VoidCloseAction,
	)

	eng.Message("Log 1")
	time.Sleep(1 * time.Millisecond)

	history := eng.Records()

	eng.Message("Log 2")
	time.Sleep(1 * time.Millisecond)

	assert.Len(t, 1, history)
}

func TestEngine_CloseIdempotency(t *testing.T) {
	counter := 0
	closeAction := func([]record.Record) error {
		counter++
		return nil
	}

	eng, _ := NewEngine(
		t.Context(),
		"InternalTest",
		10,
		VoidWriteAction,
		closeAction,
	)

	eng.Close()
	eng.Close()
	eng.Close()

	assert.Equal(t, 1, counter)
}

func TestEngine_ShutdownMechanism(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	eng, _ := NewEngine(
		ctx,
		"InternalTest",
		10,
		VoidWriteAction,
		VoidCloseAction,
	)

	cancel()

	assert.WillClose(t, eng.done, 100 * time.Millisecond)
}
