package json_test

import (
	"encoding/json"
	"strings"
	"testing"

	fmt_json "github.com/Rafael24595/go-log/log/format/json"
	"github.com/Rafael24595/go-log/log/model/record"
)

func TestJsonLinesFormat(t *testing.T) {
	records := []record.Record{
		{Message: "First log entry"},
		{Message: "Second log entry"},
	}

	output, err := fmt_json.JsonLineFormat.Format(records...)
	if err != nil {
		t.Fatalf("Failed to format records: %v", err)
	}

	trimmedOutput := strings.TrimSpace(output)
	lines := strings.Split(trimmedOutput, "\n")

	if len(lines) != len(records) {
		t.Errorf("Expected %d lines, but got %d", len(records), len(lines))
	}

	for i, line := range lines {
		var r record.Record
		if err := json.Unmarshal([]byte(line), &r); err != nil {
			t.Errorf("Line %d is not a valid JSON: %v", i+1, err)
		}

		if r.Message != records[i].Message {
			t.Errorf("Data corruption at line %d: expected %q, got %q",
				i+1, records[i].Message, r.Message)
		}
	}
}

func TestJsonLinesFormat_Empty(t *testing.T) {
	output, err := fmt_json.JsonLineFormat.Format()
	if err != nil {
		t.Errorf("Should not fail with zero records: %v", err)
	}
	if output != "" {
		t.Error("Expected empty string for zero records")
	}
}
