package console

import (
	"os"

	"github.com/Rafael24595/go-log/log"
	"github.com/Rafael24595/go-log/log/format"
	"github.com/Rafael24595/go-log/log/format/text"
	"github.com/Rafael24595/go-log/log/internal/constants"
	"github.com/Rafael24595/go-log/log/logger"
	"github.com/Rafael24595/go-log/log/provider/stream"
)

const ConsoleStream logger.Logger = "Console"

type ConsoleProvider struct {
	Buffer int
	Format *format.Format
}

func (p ConsoleProvider) Build() (log.Log, error) {
	if p.Buffer <= 0 {
		p.Buffer = constants.DefaultBufferSize
	}

	if p.Format == nil {
		p.Format = &text.TextFormat
	}

	return stream.StreamProvider{
		Name:   ConsoleStream,
		Buffer: p.Buffer,
		Format: p.Format,
		Writer: os.Stdout,
	}.Build()
}
