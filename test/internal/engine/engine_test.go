package engine

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/Rafael24595/go-log/log/model/record"
	"github.com/Rafael24595/go-log/log/provider/stream"
)

func TestEngine_Concurrency(t *testing.T) {
	var buf bytes.Buffer
	lg, _ := stream.StreamProvider{Writer: &buf}.Build()

	var wg sync.WaitGroup
	goroutines := 50
	logsPerRoutine := 100
	totalExpected := goroutines * logsPerRoutine

	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := range logsPerRoutine {
				lg.Message(fmt.Sprintf("Routine %d - Log %d", id, i))
			}
		}(g)
	}

	wg.Wait()
	records, _ := lg.Close()

	if len(records) != totalExpected {
		t.Errorf("Race condition detected: expected %d records, got %d", totalExpected, len(records))
	}
}

func TestEngine_RecordsImmutability(t *testing.T) {
	lg, _ := stream.StreamProvider{
		Writer: io.Discard,
	}.Build()

	lg.Message("Log 1")
	time.Sleep(1 * time.Millisecond)

	history := lg.Records()

	lg.Message("Log 2")
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

	lg, _ := stream.StreamProvider{
		Writer: io.Discard,
		CloseAction: closeAction,
	}.Build()

	lg.Close()
	lg.Close()
	lg.Close()

	if counter != 1 {
		t.Errorf("El CloseAction se ejecutó %d veces, debería ser solo 1", counter)
	}
}
