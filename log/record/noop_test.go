package record

import (
	"testing"

	"github.com/Rafael24595/go-assert/assert"
)

func TestNoOpStore_Behavior(t *testing.T) {
	store := NewNoOp()

	err := store.Add(Record{Message: "Ghost Log"})
	assert.Nil(t, err)

	records := store.All()
	assert.Size(t, 0, records)
}
