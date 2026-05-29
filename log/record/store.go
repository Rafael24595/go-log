package record

// Store defines the contract for managing and retaining log record history.
// It decouples the core logging engine from the specific storage mechanism,
// enabling strategies ranging from zero-allocation operations to in-memory buffering.
type Store interface {
	// Add appends a single log entry to the storage engine.
	// Returns an error if the record cannot be persisted or if the store is bounded.Add(r Record) error
	Add(r Record) error
	// All returns a point-in-time snapshot of all accumulated records.
	// This operation is non-destructive and preserves the current state of the store.All() []Record
	All() []Record
	// Drain extracts all accumulated records and atomically clears the store.
	// It transfers ownership of the historical data to the caller, resetting the
	// internal state to optimize memory footprint and prevent memory leaks.
	Drain() []Record
}
