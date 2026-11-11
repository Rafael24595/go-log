package log

import (
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"

	"github.com/Rafael24595/go-assert/assert"
	"github.com/Rafael24595/go-assert/require"
	"github.com/Rafael24595/go-log/log/internal/paths"
	"github.com/Rafael24595/go-log/log/record"
)

func TestBootstrapClose_WritesUnflushedRecordsToTheBootstrapDir(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(paths.BootstrapDirEnv, dir)

	store := record.NewMemory()

	err := store.Add(record.Record{Category: record.ERROR, Message: "lost before configuring"})
	assert.Nil(t, err)

	err = makeCloseAction(&atomic.Bool{}, store)()
	assert.Nil(t, err)

	files, err := filepath.Glob(filepath.Join(dir, "log-unsigned-*.jsonl"))
	assert.Nil(t, err)

	require.Size(t, 1, files)

	data, err := os.ReadFile(files[0])

	assert.Nil(t, err)

	assert.Inside(t, "lost before configuring", string(data))
}

func TestBootstrapClose_WritesNothingOnceFlushed(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(paths.BootstrapDirEnv, dir)

	store := record.NewMemory()

	err := store.Add(record.Record{Category: record.MESSAGE, Message: "already flushed"})
	assert.Nil(t, err)

	flushed := &atomic.Bool{}
	flushed.Store(true)

	err = makeCloseAction(flushed, store)()
	assert.Nil(t, err)

	files, _ := filepath.Glob(filepath.Join(dir, "*"))
	assert.Size(t, 0, files)
}
