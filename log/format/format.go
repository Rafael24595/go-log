package format

import "github.com/Rafael24595/go-log/log/model/record"

// Format defines how one or more log records are serialized into a string.
// It includes the file extension associated with the specific format.
type Format struct {
	// Extension returns the file suffix (json, log, etc.)
	Extension string
	// Format converts records into a string
	Format    func(records ...record.Record) (string, error)
}
