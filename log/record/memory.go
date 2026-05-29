package record

import (
	"sync"
)

// Memory implements a thread-safe, unbounded in-memory Store.
// It utilizes an internal read-write mutex to synchronize concurrent access,
// making it ideal for transient buffering scenarios like Bootstrap logging.
type Memory struct {
	mu      sync.RWMutex
	records []Record
}

// NewMemory initializes and returns a new thread-safe Memory storage instance.
func NewMemory() *Memory {
	return &Memory{}
}

// Add appends a single log entry to the internal slice in a thread-safe manner.
// It lazily allocates the underlying array on the first write and always returns nil.
func (s *Memory) Add(r Record) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.records = append(s.records, r)
	return nil
}

// All returns a thread-safe, immutable point-in-time snapshot of the history.
// It performs a deep copy of the internal slice to prevent external callers
// from mutating the stored log state.
func (s *Memory) All() []Record {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Record, len(s.records))
	copy(out, s.records)

	return out
}

// Drain atomically extracts all accumulated records and resets the internal storage.
// It safely transfers ownership of the historical data pointer to the caller,
// clearing the internal reference to optimize memory footprint and assist the GC.
func (m *Memory) Drain() []Record {
	m.mu.Lock()
	defer m.mu.Unlock()

	extracted := m.records
	m.records = nil 

	return extracted
}
