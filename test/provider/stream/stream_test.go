package stream_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-log/log/provider/stream"
)

func TestStreamLogger_Integrity(t *testing.T) {
	var buf bytes.Buffer
	totalLogs := 100

	lg, err := stream.StreamProvider{
		Writer: &buf,
		Buffer: 10,
	}.Build(t.Context())

	assert.Nil(t, err)

	for i := range totalLogs {
		lg.Message(fmt.Sprintf("log numero %d", i))
	}

	records, err := lg.Close()
	assert.Nil(t, err)

	assert.Len(t, totalLogs, records)

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	assert.Nil(t, err)
	assert.Len(t, totalLogs, lines)
}
