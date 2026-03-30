package multi

import (
	"context"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-log/log"
	log_test "github.com/Rafael24595/go-log/test"
)

func TestMultiProvider_BuildSuccess(t *testing.T) {
	ctx := context.Background()

	p1 := log_test.MockProvider{
		Name: "M1",
	}

	p2 := log_test.MockProvider{
		Name: "M2",
	}

	m := New(p1, p2)

	l, err := m.Build(ctx)

	assert.Nil(t, err)
	assert.Equal(t, "M1+M2", l.Name())
}

func TestMultiProvider_BuildFailture(t *testing.T) {
	ctx := context.Background()

	p1 := log_test.MockProvider{
		Name: "M1",
	}

	p2 := log_test.MockProvider{
		Name: "M2",
		Fail: true,
	}

	m := New(p1, p2)

	_, err := m.Build(ctx)
	assert.NotNil(t, err)
}

func TestMultiLogger_MultiplexingReachesAll(t *testing.T) {
	m1 := &log_test.MockLogger{
		Logger: "L1",
	}

	m2 := &log_test.MockLogger{
		Logger: "L2",
	}

	ml := &multiLogger{
		loggers: []log.Log{m1, m2},
	}

	msg := "hello world"
	ml.Message(msg)

	assert.Len(t, 1, m1.History)
	assert.Equal(t, msg, m1.History[0].Message)

	assert.Len(t, 1, m2.History)
	assert.Equal(t, msg, m2.History[0].Message)
}

func TestMultiLogger_MultiplexingClosesAll(t *testing.T) {
	m1 := &log_test.MockLogger{
		Logger: "L1",
	}

	m2 := &log_test.MockLogger{
		Logger: "L2",
	}

	ml := &multiLogger{
		loggers: []log.Log{m1, m2},
	}

	ml.Close()

	assert.True(t, m1.Closed())
	assert.True(t, m2.Closed())
}
