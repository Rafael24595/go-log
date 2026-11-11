package file

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Rafael24595/go-log/log"
	"github.com/Rafael24595/go-log/log/format"
	"github.com/Rafael24595/go-log/log/format/text"
	"github.com/Rafael24595/go-log/log/internal/clock"
	"github.com/Rafael24595/go-log/log/internal/engine"
	"github.com/Rafael24595/go-log/log/internal/file"
	"github.com/Rafael24595/go-log/log/internal/paths"
	"github.com/Rafael24595/go-log/log/logger"
	"github.com/Rafael24595/go-log/log/provider/stream"
	"github.com/Rafael24595/go-log/log/record"
)

// LoggerFile is the default identifier for file-based loggers.
const LoggerFile logger.Logger = "File"

// FileProvider configures and creates a logger that persists records to a physical file.
// It automatically handles file creation, naming conventions, and resource cleanup.
type FileProvider struct {
	// Session is an optional identifier added to the filename (e.g., "user-auth").
	Session string
	// Buffer size for the underlying engine channel.
	Buffer uint
	// Path is the directory where log files will be created.
	Path string
	// Format defines the layout of the log entries (e.g., Text or JSONL).
	Format *format.Format
	// RecordStore allows for optional in-memory storage of records for later retrieval.
	RecordStore record.Store
}

// New returns a new, unconfigured FileProvider as a log.Provider interface.
func New() log.Provider {
	return FileProvider{}
}

// Build initializes a file-based logger. It generates a unique filename using
// the current timestamp and session name, ensures the file is opened with
// appropriate permissions, and leverages StreamProvider for the writing logic.
func (p FileProvider) Build(ctx context.Context) (log.Log, error) {
	timestamp := clock.UnixMilliClock()

	if p.Path == "" {
		p.Path = paths.BootstrapDir()
	}

	if p.Format == nil {
		p.Format = &text.TextFormat
	}

	name := fmt.Sprintf(
		"log-%s-%s.%s",
		p.Session,
		format.FormatMillisecondsCompact(timestamp),
		p.Format.Extension,
	)

	source := filepath.Join(p.Path, name)

	file, err := file.NewFile(source, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}

	l, err := stream.StreamProvider{
		Name:        LoggerFile,
		Buffer:      p.Buffer,
		Format:      p.Format,
		Writer:      file,
		CloseAction: makeCloseAction(file),
		RecordStore: p.RecordStore,
	}.Build(ctx)

	if err != nil {
		return nil, errors.Join(err, file.Close())
	}

	return l, nil
}

func makeCloseAction(file *file.File) engine.CloseAction {
	return func() error {
		return file.Close()
	}
}
