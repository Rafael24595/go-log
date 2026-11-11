package json_test

import (
	"encoding/json"
	"strings"
	"testing"

	fmt_json "github.com/Rafael24595/go-log/log/format/json"

	"github.com/Rafael24595/go-assert/assert"
	"github.com/Rafael24595/go-assert/require"
	"github.com/Rafael24595/go-log/log/record"
)

func TestJsonLinesFormat(t *testing.T) {
	records := []record.Record{
		{Message: "First log entry"},
		{Message: "Second log entry"},
	}

	output, err := fmt_json.JsonLineFormat.Format(records...)
	assert.Nil(t, err)

	lines := strings.Split(output, "\n")

	require.Size(t, len(records)+1, lines)
	assert.Equal(t, lines[len(lines)-1], "")

	for i, rec := range records {
		var r record.Record

		err := json.Unmarshal([]byte(lines[i]), &r)
		assert.Nil(t, err)

		assert.Equal(t, r.Message, rec.Message)
	}
}

func TestJsonLinesFormat_Empty(t *testing.T) {
	output, err := fmt_json.JsonLineFormat.Format()
	assert.Nil(t, err)

	assert.Equal(t, "", output)
}
