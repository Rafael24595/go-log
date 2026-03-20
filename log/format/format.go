package format

import "github.com/Rafael24595/go-log/log/model/record"

type Format struct {
	Extension string
	Format    func(records ...record.Record) (string, error)
}
