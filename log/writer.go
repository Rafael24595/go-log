package log

import (
	"io"

	"github.com/Rafael24595/go-log/log/model/record"
)

// logWriter is a dynamic proxy that redirects writes to the active global logger.
type logWriter struct {
	category record.Category
}

func newLogWiter(category ...record.Category) io.Writer {
	categ := record.MESSAGE
	if len(category) > 0 {
		categ = record.Category(category[0])
	}

	return &logWriter{
		category: categ,
	}
}

// Write implements io.Writer. it dynamically fetches the current global 'log' 
// instance to ensure it follows the transition from Bootstrap to the final Provider.
func (w *logWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	if log.Closed() {
		return 0, io.ErrClosedPipe
	}

	log.Customc(w.category, string(p))
	return len(p), nil
}

// WriterFromString returns an io.Writer that redirects output to the 
// global logger under the specified string category (defaults to MESSAGE).
func WriterFromString(category ...string) io.Writer {
	categ := make([]record.Category, 0)
	if len(category) > 0 {
		categ = append(categ, record.Category(category[0]))
	}

	return newLogWiter(categ...)
}

// WriterFromCategory returns an io.Writer that redirects output to the 
// global logger under the specified category (defaults to MESSAGE).
func WriterFromCategory(category ...record.Category) io.Writer {
	return newLogWiter(category...)
}
