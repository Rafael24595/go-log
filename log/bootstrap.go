package log

import (
	"fmt"
	"sync/atomic"

	"github.com/Rafael24595/go-log/log/format"
	"github.com/Rafael24595/go-log/log/format/json"
	"github.com/Rafael24595/go-log/log/internal/clock"
	"github.com/Rafael24595/go-log/log/internal/constants"
	"github.com/Rafael24595/go-log/log/internal/engine"
	"github.com/Rafael24595/go-log/log/internal/file"
	"github.com/Rafael24595/go-log/log/logger"
	"github.com/Rafael24595/go-log/log/record"
)

const loggerBootstrap logger.Logger = "Bootstrap"

type bootstrapLogger struct {
	Log
	flushed *atomic.Bool
	store   record.Store
}

func newBootstrapLogger() (Bootstrap, error) {
	flushed := &atomic.Bool{}
	store := record.NewMemory()

	engine, err := engine.NewEngine(
		loggerBootstrap,
		engine.WithRecordStore(store),
		engine.WithCloseAction(
			makeCloseAction(flushed, store),
		),
	)

	if err != nil {
		return nil, err
	}

	return &bootstrapLogger{
		Log:     engine,
		flushed: flushed,
		store:   store,
	}, nil
}

func (l *bootstrapLogger) Flush(target Log) error {
	if l.flushed.Swap(true) {
		return nil
	}

	err := l.Close()
	if err != nil {
		return err
	}

	records := l.store.Drain()
	if len(records) > 0 {
		target.Record(records...)
	}

	return nil
}

func makeCloseAction(flushed *atomic.Bool, store record.Store) engine.CloseAction {
	timestamp := clock.UnixMilliClock()
	json := json.JsonLineFormat

	return func() error {
		records := store.All()
		if flushed.Load() || len(records) == 0 {
			return nil
		}

		data, err := json.Format(records...)
		if err != nil {
			return err
		}

		name := fmt.Sprintf("log-unsigned-%s", format.FormatMillisecondsCompact(timestamp))
		path := fmt.Sprintf("%s/%s.%s", constants.DefaultPath, name, json.Extension)

		return file.WriteFileSafe(path, string(data))
	}
}
