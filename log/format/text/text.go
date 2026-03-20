package text

import (
	"fmt"
	"strings"

	"github.com/Rafael24595/go-log/log/format"
	"github.com/Rafael24595/go-log/log/model/record"
)

var TextFormat = format.Format{
	Extension: "log",
	Format: func(records ...record.Record) (string, error) {
		if len(records) == 0 {
			return "", nil
		}

		lines := make([]string, len(records))
		for i, r := range records {
			timestamp := format.FormatMilliseconds(r.Timestamp)
			lines[i] = fmt.Sprintf("%s - [%s]: %s", timestamp, r.Category, r.Message)
		}

		return strings.Join(lines, "\n"), nil
	},
}
