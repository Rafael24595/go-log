package engine

import (
	"context"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-log/log/internal/constants"
	"github.com/Rafael24595/go-log/log/logger"
	"github.com/Rafael24595/go-log/log/record"
)

const testLogger logger.Logger = "TestEngine"

func TestDefaultConfig(t *testing.T) {
	cfg := defaultConfig(testLogger)

	assert.NotNil(t, cfg.ctx)
	assert.True(t, cfg.name == testLogger)
	assert.True(t, cfg.buffer == constants.DefaultBufferSize)
	assert.NotNil(t, cfg.writeAction)
	assert.NotNil(t, cfg.closeAction)
	assert.NotNil(t, cfg.recordStore)
}

func TestWithContextOption(t *testing.T) {
	cfg := defaultConfig(testLogger)
	customCtx := t.Context()

	opt := WithContext(customCtx)
	opt(&cfg)
	assert.True(t, cfg.ctx == customCtx)

	//nolint:staticcheck
	optNil := WithContext(nil)
	optNil(&cfg)
	assert.True(t, cfg.ctx == customCtx)
}

func TestWithBufferSizeOption(t *testing.T) {
	cfg := defaultConfig(testLogger)

	opt := WithBufferSize(42)
	opt(&cfg)

	assert.True(t, cfg.buffer == 42)
}

func TestWithWriteActionOption(t *testing.T) {
	cfg := defaultConfig(testLogger)
	var called bool
	customAction := func(r record.Record) error {
		called = true
		return nil
	}

	opt := WithWriteAction(customAction)
	opt(&cfg)
	_ = cfg.writeAction(record.Record{})
	assert.True(t, called)

	optNil := WithWriteAction(nil)
	optNil(&cfg)
	assert.NotNil(t, cfg.writeAction)
}

func TestWithCloseActionOption(t *testing.T) {
	cfg := defaultConfig(testLogger)
	var called bool
	customAction := func() error {
		called = true
		return nil
	}

	opt := WithCloseAction(customAction)
	opt(&cfg)
	_ = cfg.closeAction()
	assert.True(t, called)

	optNil := WithCloseAction(nil)
	optNil(&cfg)
	assert.NotNil(t, cfg.closeAction)
}

func TestWithRecordStoreOption(t *testing.T) {
	cfg := defaultConfig(testLogger)
	mockStore := record.NewNoOp()

	opt := WithRecordStore(mockStore)
	opt(&cfg)
	assert.True(t, cfg.recordStore == mockStore)

	optNil := WithRecordStore(nil)
	optNil(&cfg)
	assert.NotNil(t, cfg.recordStore)
}

func TestMakeConfigIntegration(t *testing.T) {
	ctx := context.TODO()
	
	cfg := makeConfig(testLogger,
		WithBufferSize(100),
		WithContext(ctx),
	)

	assert.True(t, cfg.name == testLogger)
	assert.True(t, cfg.buffer == 100)
	assert.True(t, cfg.ctx == ctx)
}