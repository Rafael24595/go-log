package engine

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Rafael24595/go-log/log/model/record"
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

	if len(records) != totalExpected {
		t.Errorf("Race condition detected: expected %d records, got %d", totalExpected, len(records))
	}
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

	if len(history) != 1 {
		t.Errorf("The original slice was modified. 1 expected but %d got", len(history))
	}
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

	if counter != 1 {
		t.Errorf("El CloseAction se ejecutó %d veces, debería ser solo 1", counter)
	}
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

	select {
	case <-eng.done:
	case <-time.After(100 * time.Millisecond):
		t.Error("Engine failed to close internal 'done' channel")
	}
}
