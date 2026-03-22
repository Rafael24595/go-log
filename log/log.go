package log

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/Rafael24595/go-log/log/logger"
	"github.com/Rafael24595/go-log/log/model/record"
)

var (
	log  Log
	once sync.Once
)

func init() {
	logg, err := newBootstrapLogger()
	if err != nil {
		panic(err.Error())
	}
	log = logg
}

// Provider defines the interface for types that can create a concrete Log instance.
type Provider interface {
	Build(context.Context) (Log, error)
}

// DefaultFromProvider initializes the global logger using a given provider.
// The provided context controls the lifetime of the underlying logging engine.
func DefaultFromProvider(ctx context.Context, provider Provider) error {
	target, err := provider.Build(ctx)
	if err != nil {
		return err
	}
	return DefaultFromLog(target)
}

// DefaultFromLog handles the transition from the initial Bootstrap logger
// to a permanent target logger. It ensures all previous logs are flushed.
func DefaultFromLog(target Log) error {
	if target == nil {
		return errors.New("nil logger")
	}

	var init bool
	var err error

	once.Do(func() {
		if b, ok := log.(Bootstrap); ok {
			err = b.Flush(target)
			if err != nil {
				return
			}
		}

		log = target
		log.Messagef("Logging is configured to use the %s instance.", target.Name())

		init = true
	})

	if !init {
		return errors.New("logger already initialized")
	}

	return err
}

// OnClose triggers a clean shutdown of the global logger, ensuring
// all buffered records are written and resources are released.
func OnClose() error {
	_, err := log.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error closing logger: %v\n", err)
	}
	return err
}

// Bootstrap defines a logger capable of transferring its records
// to another logger once the system is fully initialized.
type Bootstrap interface {
	Log
	// Flush transfers collected records to the target Log.
	Flush(Log) error
}

// Log defines the core logging capabilities and engine management.
type Log interface {
	// Name returns the identifier of the current logger instance (e.g., "File", "Stream").
	Name() logger.Logger
	// Records returns a slice of all log entries collected by this instance so far.
	Records() []record.Record
	// Custom logs a message with a user-defined category string.
	Custom(string, string) record.Record
	// Custome logs an error associated with a specific custom category.
	Custome(string, error) record.Record
	// Customf logs a formatted message with a user-defined category.
	Customf(string, string, ...any) record.Record
	// Message logs an informational message using the default category.
	Message(string) record.Record
	// Messagef logs a formatted informational message.
	Messagef(string, ...any) record.Record
	// Warning logs a message that highlights a potential issue or important state.
	Warning(string) record.Record
	// Warningf logs a formatted warning message.
	Warningf(string, ...any) record.Record
	// Error logs a standard Go error object.
	Error(error) record.Record
	// Errors logs an error message provided as a string.
	Errors(string) record.Record
	// Errorf logs a formatted error message.
	Errorf(string, ...any) record.Record
	// Record allows manual insertion of one or more pre-built Record objects.
	Record(...record.Record) []record.Record
	// Close gracefully shuts down the logger engine, flushing any pending data.
	Close() ([]record.Record, error)
}

// Name returns the identifier of the current global logger instance.
func Name() logger.Logger {
	return log.Name()
}

// Records returns a slice of all log entries collected by the global logger.
func Records() []record.Record {
	return log.Records()
}

// Custom logs a message with a user-defined category string using the global logger.
func Custom(category string, message string) {
	log.Custom(category, message)
}

// Custome logs an error associated with a specific custom category using the global logger.
func Custome(category string, err error) {
	log.Custome(category, err)
}

// Customf logs a formatted message with a user-defined category using the global logger.
func Customf(category string, format string, args ...any) {
	log.Customf(category, format, args...)
}

// Message logs an informational message using the default category.
func Message(message string) {
	log.Message(message)
}

// Messagef logs a formatted informational message using the global logger.
func Messagef(format string, args ...any) {
	log.Messagef(format, args...)
}

// Warning logs a message highlighting a potential issue to the global logger.
func Warning(message string) {
	log.Warning(message)
}

// Warningf logs a formatted warning message using the global logger.
func Warningf(format string, args ...any) {
	log.Warningf(format, args...)
}

// Error logs a standard Go error object to the global logger.
func Error(err error) {
	log.Error(err)
}

// Errors logs an error message string to the global logger.
func Errors(message string) {
	log.Errors(message)
}

// Errorf logs a formatted error message using the global logger.
func Errorf(format string, args ...any) {
	log.Errorf(format, args...)
}
