package stream

import (
	"context"
	"io"
	"os"

	"github.com/Rafael24595/go-log/log"
	"github.com/Rafael24595/go-log/log/format"
	"github.com/Rafael24595/go-log/log/format/text"
	"github.com/Rafael24595/go-log/log/internal/constants"
	"github.com/Rafael24595/go-log/log/internal/engine"
	"github.com/Rafael24595/go-log/log/logger"
	"github.com/Rafael24595/go-log/log/record"
)

// LoggerStream is the default identifier for stream-based loggers.
const LoggerStream logger.Logger = "Stream"

// StreamProvider configures and creates a logger that writes to an io.Writer.
// It allows customization of formatting, buffering, and output destination.
type StreamProvider struct {
	// Name is the unique identifier for this logger instance.
	Name logger.Logger
	// Buffer size for the underlying engine channel.
	Buffer uint
	// Format defines how record.Record objects are converted to strings.
	Format *format.Format
	// Writer is the destination for the logs (e.g., os.Stdout, a file, or a buffer).
	Writer io.Writer
	// CloseAction defines a custom cleanup behavior when the logger stops.
	CloseAction engine.CloseAction
	// RecordStore allows for optional in-memory storage of records for later retrieval.
	RecordStore record.Store
}

// New returns a new, unconfigured StreamProvider as a log.Provider interface.
func New() log.Provider {
	return StreamProvider{}
}

// Build validates the provider configuration and initializes the stream logger engine.
// It sets default values for Name (Stream), Buffer (from constants),
// Format (Text), and Writer (os.Stdout) if they are not provided.
func (p StreamProvider) Build(ctx context.Context) (log.Log, error) {
	if p.Name == "" {
		p.Name = LoggerStream
	}

	if p.Buffer == 0 {
		p.Buffer = constants.DefaultBufferSize
	}

	if p.Format == nil {
		p.Format = &text.TextFormat
	}

	if p.Writer == nil {
		p.Writer = os.Stdout
	}

	if p.CloseAction == nil {
		p.CloseAction = engine.VoidCloseAction
	}

	if p.RecordStore == nil {
		p.RecordStore = record.NewNoOp()
	}

	return newStreamLogger(
		ctx,
		p.Name,
		p.Buffer,
		*p.Format,
		p.Writer,
		p.CloseAction,
		p.RecordStore,
	)
}

type streamLogger struct {
	format format.Format
	writer io.Writer
}

func newStreamLogger(
	ctx context.Context,
	name logger.Logger,
	buffer uint,
	format format.Format,
	writer io.Writer,
	closeAction engine.CloseAction,
	recordStore record.Store,
) (log.Log, error) {
	instance := &streamLogger{
		format: format,
		writer: writer,
	}

	return engine.NewEngine(
		name,
		engine.WithContext(ctx),
		engine.WithBufferSize(buffer),
		engine.WithWriteAction(instance.write),
		engine.WithCloseAction(closeAction),
		engine.WithRecordStore(recordStore),
	)
}

func (l *streamLogger) write(record record.Record) error {
	data, err := l.format.Format(record)
	if err != nil {
		return err
	}

	_, err = io.WriteString(l.writer, data+"\n")
	return err
}
