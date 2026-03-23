package multi

import (
	"context"

	"github.com/Rafael24595/go-log/log"
	"github.com/Rafael24595/go-log/log/logger"
	"github.com/Rafael24595/go-log/log/model/record"
)

// MultiProvider is a composite provider that orchestrates multiple log.Provider instances.
// It allows sending the same log record to various destinations (e.g., Console and File) simultaneously.
type MultiProvider struct {
	Providers []log.Provider
}

// New creates a new MultiProvider with the given list of child providers.
func New(providers ...log.Provider) log.Provider {
	return MultiProvider{
		Providers: providers,
	}
}

// Build initializes all child providers concurrently. It returns a multiLogger 
// that multiplexes every log call to all successfully built underlying loggers.
// If one child fails to build, the error is returned immediately.
func (p MultiProvider) Build(ctx context.Context) (log.Log, error) {
	loggers := make([]log.Log, 0, len(p.Providers))

	for _, provider := range p.Providers {
		l, err := provider.Build(ctx)
		if err != nil {
			return nil, err
		}
		loggers = append(loggers, l)
	}

	return &multiLogger{
		loggers: loggers,
	}, nil
}

type multiLogger struct {
	loggers []log.Log
}

// Name returns a combined string of all child logger names separated by "+".
// Example: "Console+File+Stream".
func (m *multiLogger) Name() logger.Logger {
	var name logger.Logger
	for i, l := range m.loggers {
		if i > 0 {
			name += "+"
		}
		name += l.Name()
	}
	return name
}

// Closed returns true only if all underlying loggers are closed.
// If at least one child logger is still active, it returns false to 
// allow continued operations on the remaining functional channels.
func (m *multiLogger) Closed() bool {
	for _, l := range m.loggers {
		if !l.Closed() {
			return true
		}
	}
	return true
}

// Records aggregates and returns all log entries collected by every 
// child logger in the orchestration.
func (m *multiLogger) Records() []record.Record {
	records := make([]record.Record, 0)
	for _, l := range m.loggers {
		records = append(records, l.Records()...)
	}
	return records
}

// Custom logs a message under a user-defined category string.
// The category is normalized to ensure consistent formatting across providers.
func (m *multiLogger) Custom(category string, message string) record.Record {
	var last record.Record
	for _, l := range m.loggers {
		last = l.Custom(category, message)
	}
	return last
}

// Customc is the high-performance version of Custom. 
// It accepts a typed record.Category, avoiding type conversions and 
// ensuring that only valid, pre-defined categories are used.
func (m *multiLogger) Customc(category record.Category, message string) record.Record {
	var last record.Record
	for _, l := range m.loggers {
		last = l.Customc(category, message)
	}
	return last
}

// Custome logs a standard Go error object under a custom category.
func (m *multiLogger) Custome(category string, err error) record.Record {
	var last record.Record
	for _, l := range m.loggers {
		last = l.Custome(category, err)
	}
	return last
}

// Customf logs a formatted message under a custom category.
func (m *multiLogger) Customf(category string, format string, a ...any) record.Record {
	var last record.Record
	for _, l := range m.loggers {
		last = l.Customf(category, format, a...)
	}
	return last
}

// Message logs a plain text informational message using the default MESSAGE category.
// It broadcasts the message to all underlying loggers in the multi-provider.
func (m *multiLogger) Message(message string) record.Record {
	var last record.Record
	for _, l := range m.loggers {
		last = l.Message(message)
	}
	return last
}

// Messagef logs a formatted informational message (printf-style).
// It is useful for dynamic messages that include variables or state information.
func (m *multiLogger) Messagef(format string, a ...any) record.Record {
	var last record.Record
	for _, l := range m.loggers {
		last = l.Messagef(format, a...)
	}
	return last
}

// Warning logs a potential issue or a non-critical anomalous state.
// This is broadcasted to all active loggers to ensure visibility across all channels.
func (m *multiLogger) Warning(message string) record.Record {
	var last record.Record
	for _, l := range m.loggers {
		last = l.Warning(message)
	}
	return last
}

// Warningf logs a formatted warning message.
// Useful for providing context, such as "Retry attempt %d of %d".
func (m *multiLogger) Warningf(format string, a ...any) record.Record {
	var last record.Record
	for _, l := range m.loggers {
		last = l.Warningf(format, a...)
	}
	return last
}

// Error logs a standard Go error object. It automatically extracts the error string.
// This is the preferred method for logging caught exceptions and failures.
func (m *multiLogger) Error(err error) record.Record {
	var last record.Record
	for _, l := range m.loggers {
		last = l.Error(err)
	}
	return last
}

// Errors logs an error represented as a raw string instead of a Go error object.
func (m *multiLogger) Errors(message string) record.Record {
	var last record.Record
	for _, l := range m.loggers {
		last = l.Errors(message)
	}
	return last
}

// Errorf logs a formatted error message, allowing developers to wrap 
// error context into a single readable entry.
func (m *multiLogger) Errorf(format string, a ...any) record.Record {
	var last record.Record
	for _, l := range m.loggers {
		last = l.Errorf(format, a...)
	}
	return last
}

// Record allows manual insertion of pre-built Record objects into 
// all underlying engines simultaneously.
func (m *multiLogger) Record(records ...record.Record) []record.Record {
	for _, l := range m.loggers {
		l.Record(records...)
	}
	return records
}

// Close gracefully shuts down all child loggers. It collects all 
// flushed records from every child and returns the last encountered 
// error (if any) during the mass-closing process.
func (m *multiLogger) Close() ([]record.Record, error) {
	var allRecords []record.Record
	var lastErr error

	for _, l := range m.loggers {
		records, err := l.Close()
		if err != nil {
			lastErr = err
		}
		allRecords = append(allRecords, records...)
	}
	return allRecords, lastErr
}