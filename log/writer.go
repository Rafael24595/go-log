package log

import (
	"io"

	"github.com/Rafael24595/go-log/log/model/record"
)

// logWriter is a dynamic proxy that redirects writes to the active global logger.
type logWriter struct {
	category record.Category
	logger   Log
}

func newLogWiter(logger Log, category ...record.Category) io.Writer {
	categ := record.MESSAGE
	if len(category) > 0 {
		categ = record.Category(category[0])
	}

	return &logWriter{
		logger:   logger,
		category: categ,
	}
}

// Write implements the io.Writer interface. It converts byte slices to strings 
// and dispatches them to the target logger. If the global logger is closed, 
// it returns io.ErrClosedPipe to signify the stream is no longer available.
func (w *logWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	if log.Closed() {
		return 0, io.ErrClosedPipe
	}

	instance := log
	if w.logger != nil {
		instance = w.logger
	}

	instance.Customc(w.category, string(p))

	return len(p), nil
}

// NewWriterFromString creates an io.Writer for a specific logger using a string category.
// If logger is nil, it will point to the global logger.
func NewWriterFromString(logger Log, category ...string) io.Writer {
	categ := make([]record.Category, 0)
	if len(category) > 0 {
		categ = append(categ, record.Category(category[0]))
	}

	return NewWriterFromCategory(logger, categ...)
}

// NewWriterFromCategory creates an io.Writer for a specific logger using a typed record.Category.
// This is the most efficient way to create a writer for a non-global logger.
func NewWriterFromCategory(logger Log, category ...record.Category) io.Writer {
	return newLogWiter(logger, category...)
}

// WriterFromString returns an io.Writer that redirects output to the
// global logger under the specified string category (defaults to MESSAGE).
func WriterFromString(category ...string) io.Writer {
	return NewWriterFromString(nil, category...)
}

// WriterFromCategory returns an io.Writer that redirects output to the
// global logger under the specified category (defaults to MESSAGE).
func WriterFromCategory(category ...record.Category) io.Writer {
	return NewWriterFromCategory(nil, category...)
}
