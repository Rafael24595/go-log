package json

import (
	"encoding/json"

	"github.com/Rafael24595/go-log/log/format"
	"github.com/Rafael24595/go-log/log/model/record"
)

var JsonFormat = format.Format{
	Extension: "json",
	Format: func(records ...record.Record) (string, error) {
		if len(records) == 0 {
			return "", nil
		}

		data, err := json.MarshalIndent(records, "", "  ")
		if err != nil {
			return "", err
		}

		return string(data), nil
	},
}
