package engine

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/Rafael24595/go-log/log/internal/clock"
	"github.com/Rafael24595/go-log/log/logger"
	"github.com/Rafael24595/go-log/log/model/record"
)

// WriteAction defines the function signature for persisting a single record.
type WriteAction func(record.Record, []record.Record) error

// CloseAction defines the function signature for cleanup operations when the engine stops.
type CloseAction func([]record.Record) error

func VoidWriteAction(record.Record, []record.Record) error { return nil }
func VoidCloseAction([]record.Record) error                { return nil }

// Engine is the core concurrent processor for log entries.
// It handles asynchronous writing via channels and manages a history of records.
type Engine struct {
	mu  sync.RWMutex
	ctx context.Context

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

// NewEngine initializes a new engine, starts the background processing loop,
// and begins watching the provided context for cancellation.
func NewEngine(
	ctx context.Context,
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
	go logger.watchExit(ctx)

	return logger, nil
}

// Name returns the unique identifier assigned to this engine instance.
func (l *Engine) Name() logger.Logger {
	return l.name
}

// Closed returns true if this engine instance has been shut down.
func (l *Engine) Closed() bool {
	return l.closed.Load()
}

// Records returns a thread-safe copy of all records processed by the engine.
// It uses a read-lock to allow multiple concurrent readers while preventing
// data races during background writes.
func (l *Engine) Records() []record.Record {
	l.mu.RLock()
	defer l.mu.RUnlock()

	out := make([]record.Record, len(l.records))
	copy(out, l.records)

	return out
}

// Custom processes a message with a specific category string.
// It automatically converts the category to uppercase for consistency.
func (l *Engine) Custom(category string, message string) record.Record {
	upperCategory := strings.ToUpper(category)
	return l.Customc(record.Category(upperCategory), message)
}

// Custom processes a message with a specific category.
func (l *Engine) Customc(category record.Category, message string) record.Record {
	return l.write(category, message)
}

// Custome is a helper that logs a custom category using an error's message.
func (l *Engine) Custome(category string, err error) record.Record {
	return l.Custom(category, err.Error())
}

// Customf logs a formatted custom message.
func (l *Engine) Customf(category string, format string, a ...any) record.Record {
	return l.Custom(category, fmt.Sprintf(format, a...))
}

// Message dispatches an informational record to the processing loop.
func (l *Engine) Message(message string) record.Record {
	return l.write(record.MESSAGE, message)
}

// Messagef dispatches a formatted informational record.
func (l *Engine) Messagef(format string, a ...any) record.Record {
	return l.Message(fmt.Sprintf(format, a...))
}

// Warning dispatches a warning record to the processing loop.
func (l *Engine) Warning(message string) record.Record {
	return l.write(record.WARNING, message)
}

// Warningf dispatches a formatted warning record.
func (l *Engine) Warningf(format string, a ...any) record.Record {
	return l.Warning(fmt.Sprintf(format, a...))
}

// Error dispatches a record based on a standard Go error.
func (l *Engine) Error(err error) record.Record {
	return l.Errors(err.Error())
}

// Errors dispatches an error record using a string message.
func (l *Engine) Errors(message string) record.Record {
	return l.write(record.ERROR, message)
}

// Errorf dispatches a formatted error record.
func (l *Engine) Errorf(format string, a ...any) record.Record {
	return l.Errors(fmt.Sprintf(format, a...))
}

// Record sends one or more pre-built records directly into the engine's channel.
// This is useful for re-logging records from a bootstrap or buffer.
func (l *Engine) Record(records ...record.Record) []record.Record {
	for _, record := range records {
		l.ch <- record
	}
	return records
}

// Close initiates a graceful shutdown. It stops accepting new records,
// waits for the internal buffer to be processed by the runLoop, and
// finally executes the CloseAction.
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

// Done returns a channel that is closed when the engine has fully stopped.
func (l *Engine) Done() <-chan struct{} {
	return l.done
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

func (l *Engine) watchExit(ctx context.Context) {
	select {
	case <-ctx.Done():
		l.Close()
	case <-l.done:
		return
	}
}
