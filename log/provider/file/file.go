package file

import (
	"fmt"
	"os"

	"github.com/Rafael24595/go-log/log"
	"github.com/Rafael24595/go-log/log/format"
	"github.com/Rafael24595/go-log/log/format/text"
	"github.com/Rafael24595/go-log/log/internal/clock"
	"github.com/Rafael24595/go-log/log/internal/constants"
	"github.com/Rafael24595/go-log/log/internal/engine"
	"github.com/Rafael24595/go-log/log/internal/file"
	"github.com/Rafael24595/go-log/log/logger"
	"github.com/Rafael24595/go-log/log/model/record"
	"github.com/Rafael24595/go-log/log/provider/stream"
)

const LoggerFile logger.Logger = "File"

type FileProvider struct {
	Session string
	Buffer  int
	Path    string
	Format  *format.Format
}

func (p FileProvider) Build() (log.Log, error) {
	timestamp := clock.UnixMilliClock()

	if p.Path == "" {
		p.Path = constants.DefaultPath
	}

	name := fmt.Sprintf("log-%s-%s", p.Session, format.FormatMillisecondsCompact(timestamp))
	source := fmt.Sprintf("%s/%s.%s", p.Path, name, text.TextFormat.Extension)

	file, err := file.NewFile(source, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}

	return stream.StreamProvider{
		Name:        LoggerFile,
		Buffer:      p.Buffer,
		Format:      p.Format,
		Writer:      file,
		CloseAction: makeCloseAction(file),
	}.Build()
}

func makeCloseAction(file *file.File) engine.CloseAction {
	return func([]record.Record) error {
		return file.Close()
	}
}
