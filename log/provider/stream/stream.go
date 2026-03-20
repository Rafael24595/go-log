package stream

import (
	"io"
	"os"

	"github.com/Rafael24595/go-log/log"
	"github.com/Rafael24595/go-log/log/format"
	"github.com/Rafael24595/go-log/log/format/text"
	"github.com/Rafael24595/go-log/log/internal/constants"
	"github.com/Rafael24595/go-log/log/internal/engine"
	"github.com/Rafael24595/go-log/log/logger"
	"github.com/Rafael24595/go-log/log/model/record"
)

const LoggerStream logger.Logger = "Stream"

type StreamProvider struct {
	Name        logger.Logger
	Buffer      int
	Format      *format.Format
	Writer      io.Writer
	CloseAction engine.CloseAction
}

func (p StreamProvider) Build() (log.Log, error) {
	if p.Name == "" {
		p.Name = LoggerStream
	}

	if p.Buffer <= 0 {
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

	return newStreamLogger(
		p.Name,
		p.Buffer,
		*p.Format,
		p.Writer,
		p.CloseAction,
	)
}

type streamLogger struct {
	format format.Format
	writer io.Writer
}

func newStreamLogger(
	name logger.Logger,
	buffer int,
	format format.Format,
	writer io.Writer,
	closeAction engine.CloseAction,
) (log.Log, error) {
	instance := &streamLogger{
		format: format,
		writer: writer,
	}

	return engine.NewEngine(
		name,
		buffer,
		instance.write,
		closeAction,
	)
}

func (l *streamLogger) write(record record.Record, _ []record.Record) error {
	data, err := l.format.Format(record)
	if err != nil {
		return err
	}

	_, err = io.WriteString(l.writer, data+"\n")

	return nil
}
