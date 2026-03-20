package json

import (
	"encoding/json"
	"strings"

	"github.com/Rafael24595/go-log/log/format"
	"github.com/Rafael24595/go-log/log/model/record"
)

var JsonFormat = format.Format{
	Extension: "jsonl",
	Format: func(records ...record.Record) (string, error) {
		if len(records) == 0 {
			return "", nil
		}

		var buf strings.Builder
		encoder := json.NewEncoder(&buf)

		for _, r := range records {
			if err := encoder.Encode(r); err != nil {
				return "", err
			}
		}

		return buf.String(), nil
	},
}
