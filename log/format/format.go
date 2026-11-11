package format

import "github.com/Rafael24595/go-log/log/record"

// Format defines how one or more log records are serialized into a string.
// It includes the file extension associated with the specific format.
type Format struct {
	// Extension returns the file suffix (json, log, etc.)
	Extension string
	// Format serialises the records. The result is written as it is, so the
	// format decides the line breaks: end each record with "\n" to get one per line.
	Format func(records ...record.Record) (string, error)
}
