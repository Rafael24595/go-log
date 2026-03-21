package console

import (
	"context"
	"os"

	"github.com/Rafael24595/go-log/log"
	"github.com/Rafael24595/go-log/log/format"
	"github.com/Rafael24595/go-log/log/format/text"
	"github.com/Rafael24595/go-log/log/internal/constants"
	"github.com/Rafael24595/go-log/log/logger"
	"github.com/Rafael24595/go-log/log/provider/stream"
)

// ConsoleStream is the default identifier for console-based loggers.
const ConsoleStream logger.Logger = "Console"

// ConsoleProvider simplifies the creation of a logger that outputs to the standard output (stdout).
// It is a specialized wrapper around StreamProvider pre-configured for terminal usage.
type ConsoleProvider struct {
	Buffer int
	Format *format.Format
}

// New returns a new, unconfigured ConsoleProvider as a log.Provider interface.
func New() log.Provider {
	return ConsoleProvider{}
}

// Build initializes a console logger. It ensures logs are directed to os.Stdout 
// and applies default values for buffering and formatting if they are not specified.
func (p ConsoleProvider) Build(ctx context.Context) (log.Log, error) {
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
	}.Build(ctx)
}
