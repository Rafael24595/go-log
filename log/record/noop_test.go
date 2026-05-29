package record

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestNoOpStore_Behavior(t *testing.T) {
	store := NewNoOp()

	err := store.Add(Record{Message: "Ghost Log"})
	assert.Nil(t, err)

	records := store.All()
	assert.Len(t, 0, records)
}
