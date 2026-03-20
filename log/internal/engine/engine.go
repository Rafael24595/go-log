package engine

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/Rafael24595/go-log/log/internal/clock"
	"github.com/Rafael24595/go-log/log/logger"
	"github.com/Rafael24595/go-log/log/model/record"
)

type WriteAction func(record.Record, []record.Record) error
type CloseAction func([]record.Record) error

func VoidWriteAction(record.Record, []record.Record) error { return nil }
func VoidCloseAction([]record.Record) error                { return nil }

type Engine struct {
	mu sync.RWMutex

	ch     chan record.Record
	errCh  chan error
	done   chan struct{}
	closed atomic.Bool

	timestamp int64
	clock     clock.Clock

	name        logger.Logger
	writeAction WriteAction
	closeAction CloseAction
	records     []record.Record
}

func NewEngine(
	name logger.Logger,
	buffer int,
	writeAction WriteAction,
	closeAction CloseAction,
) (*Engine, error) {
	timestamp := clock.UnixMilliClock()

	logger := &Engine{
		ch:          make(chan record.Record, buffer),
		errCh:       make(chan error, 1),
		done:        make(chan struct{}),
		clock:       clock.UnixMilliClock,
		timestamp:   timestamp,
		name:        name,
		writeAction: writeAction,
		closeAction: closeAction,
		records:     make([]record.Record, 0),
	}

	go logger.runLoop()

	return logger, nil
}

func (l *Engine) Name() logger.Logger {
	return l.name
}

func (l *Engine) Records() []record.Record {
	l.mu.RLock()
	defer l.mu.RUnlock()

	out := make([]record.Record, len(l.records))
	copy(out, l.records)

	return out
}

func (l *Engine) Custom(category string, message string) record.Record {
	upperCategory := strings.ToUpper(category)
	return l.write(record.Category(upperCategory), message)
}

func (l *Engine) Custome(category string, err error) record.Record {
	return l.Custom(category, err.Error())
}

func (l *Engine) Customf(category string, format string, a ...any) record.Record {
	return l.Custom(category, fmt.Sprintf(format, a...))
}

func (l *Engine) Message(message string) record.Record {
	return l.write(record.MESSAGE, message)
}

func (l *Engine) Messagef(format string, a ...any) record.Record {
	return l.Message(fmt.Sprintf(format, a...))
}

func (l *Engine) Warning(message string) record.Record {
	return l.write(record.WARNING, message)
}

func (l *Engine) Warningf(format string, a ...any) record.Record {
	return l.Warning(fmt.Sprintf(format, a...))
}

func (l *Engine) Error(err error) record.Record {
	return l.Errors(err.Error())
}

func (l *Engine) Errors(message string) record.Record {
	return l.write(record.ERROR, message)
}

func (l *Engine) Errorf(format string, a ...any) record.Record {
	return l.Errors(fmt.Sprintf(format, a...))
}

func (l *Engine) Record(records ...record.Record) []record.Record {
	for _, record := range records {
		l.ch <- record
	}
	return records
}

func (l *Engine) Close() ([]record.Record, error) {
	var err error
	if l.closed.CompareAndSwap(false, true) {
		close(l.ch)
		<-l.done
		
		err = l.closeAction(l.records)
	}

	records := l.records
	l.records = nil

	return records, err
}

func (l *Engine) write(category record.Category, message string) record.Record {
	record := record.Record{
		Category:  category,
		Message:   message,
		Timestamp: l.clock(),
	}

	if l.closed.Load() {
		return record
	}

	l.ch <- record

	return record
}

func (l *Engine) runLoop() {
	for record := range l.ch {
		l.mu.Lock()
		l.records = append(l.records, record)
		l.mu.Unlock()

		//TODO: Manage write error.
		_ = l.writeAction(record, l.records)
	}

	close(l.done)
}
