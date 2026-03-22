package stream_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/Rafael24595/go-log/log/provider/stream"
)

func TestStreamLogger_Integrity(t *testing.T) {
	var buf bytes.Buffer
	totalLogs := 100

	lg, err := stream.StreamProvider{
		Writer: &buf,
		Buffer: 10,
	}.Build(t.Context())

	if err != nil {
		t.Fatalf("failed to build logger: %v", err)
	}

	for i := range totalLogs {
		lg.Message(fmt.Sprintf("log numero %d", i))
	}

	records, err := lg.Close()
	if err != nil {
		t.Errorf("error closing logger: %v", err)
	}

	if len(records) != totalLogs {
		t.Errorf("expected %d records, got %d", totalLogs, len(records))
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != totalLogs {
		t.Errorf("expected %d lines in buffer, got %d", totalLogs, len(lines))
	}
}
