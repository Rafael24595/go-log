package record

import (
	"fmt"
	"sync"
	"testing"

	"github.com/Rafael24595/go-assert/assert"
	"github.com/Rafael24595/go-assert/require"
)

func TestMemoryStore_SaveAndRetrieve(t *testing.T) {
	store := NewMemory()

	err := store.Add(Record{Message: "first_log"})
	assert.Nil(t, err)

	err = store.Add(Record{Message: "second_log"})
	assert.Nil(t, err)

	records := store.All()

	require.Size(t, 2, records)

	assert.True(t, records[0].Message == "first_log")
	assert.True(t, records[1].Message == "second_log")
}

func TestMemoryStore_Drain(t *testing.T) {
	store := NewMemory()
	_ = store.Add(Record{Message: "ephemeral_log"})

	extracted := store.Drain()

	require.Size(t, 1, extracted)
	assert.True(t, extracted[0].Message == "ephemeral_log")

	assert.Size(t, 0, store.All())
}

func TestMemoryStore_Immutability(t *testing.T) {
	store := NewMemory()
	_ = store.Add(Record{Message: "legit_log"})

	store.All()[0] = Record{Message: "hacked_log"}

	assert.True(t, store.All()[0].Message == "legit_log")
}

func TestMemoryStore_Concurrency(t *testing.T) {
	var wg sync.WaitGroup
	routines := 20
	logs := 50

	store := NewMemory()

	for range routines {
		wg.Go(func() {
			for i := range logs {
				_ = store.Add(Record{Message: fmt.Sprintf("log_%d", i)})
			}
		})
	}

	wg.Wait()

	assert.Size(t, routines*logs, store.All())
}
