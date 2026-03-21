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

func (m *multiLogger) Records() []record.Record {
	records := make([]record.Record, 0)
	for _, l := range m.loggers {
		records = append(records, l.Records()...)
	}
	return records
}

func (m *multiLogger) Custom(category string, message string) record.Record {
	var last record.Record
	for _, l := range m.loggers {
		last = l.Custom(category, message)
	}
	return last
}

func (m *multiLogger) Custome(category string, err error) record.Record {
	var last record.Record
	for _, l := range m.loggers {
		last = l.Custome(category, err)
	}
	return last
}

func (m *multiLogger) Customf(category string, format string, a ...any) record.Record {
	var last record.Record
	for _, l := range m.loggers {
		last = l.Customf(category, format, a...)
	}
	return last
}

func (m *multiLogger) Message(message string) record.Record {
	var last record.Record
	for _, l := range m.loggers {
		last = l.Message(message)
	}
	return last
}

func (m *multiLogger) Messagef(format string, a ...any) record.Record {
	var last record.Record
	for _, l := range m.loggers {
		last = l.Messagef(format, a...)
	}
	return last
}

func (m *multiLogger) Warning(message string) record.Record {
	var last record.Record
	for _, l := range m.loggers {
		last = l.Warning(message)
	}
	return last
}

func (m *multiLogger) Warningf(format string, a ...any) record.Record {
	var last record.Record
	for _, l := range m.loggers {
		last = l.Warningf(format, a...)
	}
	return last
}

func (m *multiLogger) Error(err error) record.Record {
	var last record.Record
	for _, l := range m.loggers {
		last = l.Error(err)
	}
	return last
}

func (m *multiLogger) Errors(message string) record.Record {
	var last record.Record
	for _, l := range m.loggers {
		last = l.Errors(message)
	}
	return last
}

func (m *multiLogger) Errorf(format string, a ...any) record.Record {
	var last record.Record
	for _, l := range m.loggers {
		last = l.Errorf(format, a...)
	}
	return last
}

func (m *multiLogger) Record(records ...record.Record) []record.Record {
	for _, l := range m.loggers {
		l.Record(records...)
	}
	return records
}

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