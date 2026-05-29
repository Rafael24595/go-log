package engine

import (
	"context"

	"github.com/Rafael24595/go-log/log/internal/constants"
	"github.com/Rafael24595/go-log/log/logger"
	"github.com/Rafael24595/go-log/log/record"
)

// Option defines a functional configuration closure applied to the internal engine runtime.
type Option func(*config)

type config struct {
	ctx         context.Context
	name        logger.Logger
	buffer      uint
	writeAction WriteAction
	closeAction CloseAction
	recordStore record.Store
}

func makeConfig(name logger.Logger, opts ...Option) config {
	cfg := defaultConfig(name)
	for _, o := range opts {
		o(&cfg)
	}
	return cfg
}

func defaultConfig(name logger.Logger) config {
	return config{
		ctx:         context.Background(),
		name:        name,
		buffer:      constants.DefaultBufferSize,
		writeAction: VoidWriteAction,
		closeAction: VoidCloseAction,
		recordStore: record.NewNoOp(),
	}
}

// WithContext overrides the default background context used by the engine's processing loops.
// It guards against nil values, preserving the pre-configured context if called with nil.
func WithContext(ctx context.Context) Option {
	return func(cfg *config) {
		if ctx != nil {
			cfg.ctx = ctx
		}
	}
}

// WithBufferSize tunes the allocation capacity of the underlying asynchronous log channel.
// Pass 0 to create an unbuffered (blocking/synchronous hand-off) pipeline.
func WithBufferSize(size uint) Option {
	return func(cfg *config) {
		cfg.buffer = size
	}
}

// WithWriteAction injects the core processing callback triggered on every incoming log record.
// It ignores nil actions to prevent unhandled runtime panic dispatches during the loop.
func WithWriteAction(action WriteAction) Option {
	return func(cfg *config) {
		if action != nil {
			cfg.writeAction = action
		}
	}
}

// WithCloseAction registers a cleanup hook executed during the engine's graceful shutdown cycle.
// It ignores nil actions, maintaining the engine's non-blocking shutdown boundaries intact.
func WithCloseAction(action CloseAction) Option {
	return func(cfg *config) {
		if action != nil {
			cfg.closeAction = action
		}
	}
}

// WithRecordStore couples a historical log repository to the active background logging stream.
// If the provided store is nil, the engine retains its existing storage strategy without altering state.
func WithRecordStore(store record.Store) Option {
	return func(cfg *config) {
		if store != nil {
			cfg.recordStore = store
		}
	}
}
