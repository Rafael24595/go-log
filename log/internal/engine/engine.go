package engine

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/Rafael24595/go-log/log/internal/clock"
	"github.com/Rafael24595/go-log/log/logger"
	"github.com/Rafael24595/go-log/log/record"
)

// WriteAction defines the function signature for persisting a single record.
type WriteAction func(record.Record) error

// CloseAction defines the function signature for cleanup operations when the engine stops.
type CloseAction func() error

func VoidWriteAction(record.Record) error { return nil }
func VoidCloseAction() error              { return nil }

// Engine is the core concurrent processor for log entries.
// It handles asynchronous writing via channels and manages a history of records.
type Engine struct {
	rw sync.RWMutex

	ch     chan record.Record
	errCh  chan error
	done   chan struct{}
	closed atomic.Bool

	timestamp int64
	clock     clock.Clock

	name        logger.Logger
	writeAction WriteAction
	closeAction CloseAction

	recordStore record.Store
}

// NewEngine initializes a new engine, starts the background processing loop,
// and begins watching the provided context for cancellation.
func NewEngine(
	name logger.Logger,
	opts ...Option,
) (*Engine, error) {
	cfg := makeConfig(name, opts...)
	timestamp := clock.UnixMilliClock()

	logger := &Engine{
		ch:          make(chan record.Record, cfg.buffer),
		errCh:       make(chan error, 1),
		done:        make(chan struct{}),
		clock:       clock.UnixMilliClock,
		timestamp:   timestamp,
		name:        name,
		writeAction: cfg.writeAction,
		closeAction: cfg.closeAction,
		recordStore: cfg.recordStore,
	}

	go logger.runLoop()
	go logger.watchExit(cfg.ctx)

	return logger, nil
}

// Name returns the unique identifier assigned to this engine instance.
func (e *Engine) Name() logger.Logger {
	return e.name
}

// Closed returns true if this engine instance has been shut down.
func (e *Engine) Closed() bool {
	return e.closed.Load()
}

// Custom processes a message with a specific category string.
// It automatically converts the category to uppercase for consistency.
func (e *Engine) Custom(category string, message string) record.Record {
	upperCategory := strings.ToUpper(category)
	return e.Customc(record.Category(upperCategory), message)
}

// Custom processes a message with a specific category.
func (e *Engine) Customc(category record.Category, message string) record.Record {
	return e.write(category, message)
}

// Custome is a helper that logs a custom category using an error's message.
func (e *Engine) Custome(category string, err error) record.Record {
	return e.Custom(category, err.Error())
}

// Customf logs a formatted custom message.
func (e *Engine) Customf(category string, format string, a ...any) record.Record {
	return e.Custom(category, fmt.Sprintf(format, a...))
}

// Message dispatches an informational record to the processing loop.
func (e *Engine) Message(message string) record.Record {
	return e.write(record.MESSAGE, message)
}

// Messagef dispatches a formatted informational record.
func (e *Engine) Messagef(format string, a ...any) record.Record {
	return e.Message(fmt.Sprintf(format, a...))
}

// Warning dispatches a warning record to the processing loop.
func (e *Engine) Warning(message string) record.Record {
	return e.write(record.WARNING, message)
}

// Warningf dispatches a formatted warning record.
func (e *Engine) Warningf(format string, a ...any) record.Record {
	return e.Warning(fmt.Sprintf(format, a...))
}

// Error dispatches a record based on a standard Go error.
func (e *Engine) Error(err error) record.Record {
	return e.Errors(err.Error())
}

// Errors dispatches an error record using a string message.
func (e *Engine) Errors(message string) record.Record {
	return e.write(record.ERROR, message)
}

// Errorf dispatches a formatted error record.
func (e *Engine) Errorf(format string, a ...any) record.Record {
	return e.Errors(fmt.Sprintf(format, a...))
}

// Record sends one or more pre-built records directly into the engine's channel.
// This is useful for re-logging records from a bootstrap or buffer.
func (e *Engine) Record(records ...record.Record) []record.Record {
	e.rw.RLock()
	defer e.rw.RUnlock()

	if e.closed.Load() {
		return records
	}

	for _, record := range records {
		e.ch <- record
	}
	return records
}

// Close initiates a graceful shutdown. It stops accepting new records,
// waits for the internal buffer to be processed by the runLoop, and
// finally executes the CloseAction.
func (l *Engine) Close() error {
	l.rw.Lock()

	closing := l.closed.CompareAndSwap(false, true)
	if closing {
		close(l.ch)
	}

	l.rw.Unlock()

	if !closing {
		return nil
	}

	<-l.done
	return l.closeAction()
}

// Done returns a channel that is closed when the engine has fully stopped.
func (e *Engine) Done() <-chan struct{} {
	return e.done
}

func (e *Engine) write(category record.Category, message string) record.Record {
	record := record.Record{
		Category:  category,
		Message:   message,
		Timestamp: e.clock(),
	}

	e.rw.RLock()
	defer e.rw.RUnlock()

	if e.closed.Load() {
		return record
	}

	e.ch <- record

	return record
}

func (e *Engine) runLoop() {
	for record := range e.ch {
		err := e.recordStore.Add(record)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		err = e.writeAction(record)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	close(e.done)
}

func (e *Engine) watchExit(ctx context.Context) {
	select {
	case <-ctx.Done():
		if err := e.Close(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	case <-e.done:
		return
	}
}
