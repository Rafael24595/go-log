package log_test

import (
	"context"
	"errors"

	"github.com/Rafael24595/go-log/log"
	"github.com/Rafael24595/go-log/log/logger"
	"github.com/Rafael24595/go-log/log/model/record"
)

type MockProvider struct {
	Name logger.Logger
	Fail bool
}

func (p MockProvider) Build(ctx context.Context) (log.Log, error) {
	if p.Fail {
		return nil, errors.New("forced failure")
	}
	return &MockLogger{Logger: p.Name}, nil
}

type MockLogger struct {
	Logger  logger.Logger
	History []record.Record
	Exit    bool
}

func (m *MockLogger) Name() logger.Logger {
	return m.Logger
}

func (m *MockLogger) Closed() bool {
	return m.Exit
}

func (m *MockLogger) Records() []record.Record {
	return m.History
}

func (m *MockLogger) Custom(cat string, msg string) record.Record {
	return m.Customc(record.Category(cat), msg)
}

func (m *MockLogger) Customc(cat record.Category, msg string) record.Record {
	r := record.Record{Category: cat, Message: msg}
	m.History = append(m.History, r)
	return r
}

func (m *MockLogger) Custome(cat string, err error) record.Record {
	return m.Custom(cat, err.Error())
}

func (m *MockLogger) Customf(cat string, f string, a ...any) record.Record {
	return m.Custom(cat, "fmt")
}

func (m *MockLogger) Message(msg string) record.Record {
	return m.Custom("MESSAGE", msg)
}

func (m *MockLogger) Messagef(f string, a ...any) record.Record {
	return m.Message("fmt")
}

func (m *MockLogger) Warning(msg string) record.Record {
	return m.Custom("WARNING", msg)
}

func (m *MockLogger) Warningf(f string, a ...any) record.Record {
	return m.Warning("fmt")
}

func (m *MockLogger) Error(err error) record.Record {
	return m.Custom("ERROR", err.Error())
}

func (m *MockLogger) Errors(msg string) record.Record {
	return m.Custom("ERROR", msg)
}

func (m *MockLogger) Errorf(f string, a ...any) record.Record {
	return m.Errors("fmt")
}

func (m *MockLogger) Record(recs ...record.Record) []record.Record {
	m.History = append(m.History, recs...)
	return recs
}

func (m *MockLogger) Close() ([]record.Record, error) {
	m.Exit = true
	return m.History, nil
}
