package json_test

import (
	"encoding/json"
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	fmt_json "github.com/Rafael24595/go-log/log/format/json"
	"github.com/Rafael24595/go-log/log/model/record"
)

func TestJsonLinesFormat(t *testing.T) {
	records := []record.Record{
		{Message: "First log entry"},
		{Message: "Second log entry"},
	}

	output, err := fmt_json.JsonLineFormat.Format(records...)
	assert.Nil(t, err)

	trimmedOutput := strings.TrimSpace(output)
	lines := strings.Split(trimmedOutput, "\n")

	assert.Len(t, len(lines), records)

	for i, line := range lines {
		var r record.Record

		err := json.Unmarshal([]byte(line), &r)
		assert.Nil(t, err)

		assert.Equal(t, r.Message, records[i].Message)
	}
}

func TestJsonLinesFormat_Empty(t *testing.T) {
	output, err := fmt_json.JsonLineFormat.Format()
	assert.Nil(t, err)

	assert.Equal(t, "", output)
}
