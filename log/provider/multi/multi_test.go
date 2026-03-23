package multi

import (
	"context"
	"testing"

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
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if l.Name() != "M1+M2" {
		t.Errorf("expected name M1+M2, got %s", l.Name())
	}
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
	if err == nil {
		t.Fatal("expected error from failed provider, got nil")
	}
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

	if len(m1.History) != 1 || m1.History[0].Message != msg {
		t.Errorf("logger 1 did not receive message")
	}
	if len(m2.History) != 1 || m2.History[0].Message != msg {
		t.Errorf("logger 2 did not receive message")
	}
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
	if !m1.Closed() || !m2.Closed() {
		t.Error("not all loggers were closed")
	}
}
