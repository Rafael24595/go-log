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
	"github.com/Rafael24595/go-log/log/model/record"
)

const loggerBootstrap logger.Logger = "Bootstrap"

type bootstrapLogger struct {
	Log
	flushed *atomic.Bool
}

func newBootstrapLogger() (Bootstrap, error) {
	flushed := &atomic.Bool{}

	engine, err := engine.NewEngine(
		loggerBootstrap,
		constants.DefaultBufferSize,
		engine.VoidWriteAction,
		makeCloseAction(flushed),
	)

	if err != nil {
		return nil, err
	}

	return &bootstrapLogger{
		Log:     engine,
		flushed: flushed,
	}, nil
}

func makeCloseAction(flushed *atomic.Bool) engine.CloseAction {
	timestamp := clock.UnixMilliClock()
	json := json.JsonFormat

	return func(records []record.Record) error {
		if flushed.Load() || len(records) == 0 {
			return nil
		}

		data, err := json.Format(records...)
		if err != nil {
			return err
		}

		name := fmt.Sprintf("log-unsigned-%s", format.FormatMillisecondsCompact(timestamp))
		path := fmt.Sprintf("%s/%s.%s", constants.DefaultPath, name, json.Extension)

		file.WriteFileSafe(path, string(data))

		return nil
	}
}

func (l *bootstrapLogger) Flush(target Log) error {
	if l.flushed.Swap(true) {
        return nil 
    }
	
	l.flushed.Store(true)

	records, err := l.Close()
	if err != nil {
		return err
	}

	if len(records) > 0 {
		target.Record(records...)
	}

	return nil
}
