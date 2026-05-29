package record

// NoOp implements a lock-free, zero-allocation Store that discards all records.
// It is designed as a high-performance placeholder for production environments
// where historical tracking is disabled, ensuring zero memory consumption.
type NoOp struct{}

// NewNoOp initializes and returns a new NoOp storage instance.
func NewNoOp() *NoOp {
	return &NoOp{}
}

// Add satisfies the Store interface by silently discarding the record.
// It performs no operations and always returns nil.
func (n *NoOp) Add(r Record) error {
	return nil
}

// All satisfies the Store interface by returning an empty slice allocation.
// This operation is non-destructive and safe for concurrent use.
func (n *NoOp) All() []Record {
	return make([]Record, 0)
}

// Drain satisfies the Store interface by returning an empty slice.
// It performs no state mutation since no records are ever retained.
func (n *NoOp) Drain() []Record {
	return n.All()
}
