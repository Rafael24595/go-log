package provider

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-log/log/format/json"
	"github.com/Rafael24595/go-log/log/provider/file"
)

// build writes one message through a FileProvider and returns the only file it created.
func build(t *testing.T, p file.FileProvider) string {
	t.Helper()

	p.Path = filepath.Join(t.TempDir(), "logs") // a missing folder must be created
	l, err := p.Build(t.Context())
	assert.Nil(t, err)

	l.Message("hello")
	assert.Nil(t, l.Close())

	files, err := filepath.Glob(filepath.Join(p.Path, "*"))
	assert.Nil(t, err)
	assert.Len(t, 1, files)
	return files[0]
}

func TestFileProvider_DefaultFormatIsText(t *testing.T) {
	path := build(t, file.FileProvider{Session: "s"})

	assert.True(t, strings.HasPrefix(filepath.Base(path), "log-s-"))
	assert.Equal(t, ".log", filepath.Ext(path))

	data, err := os.ReadFile(path)
	assert.Nil(t, err)
	assert.Contains(t, string(data), "[MESSAGE]: hello\n")
}

func TestFileProvider_ExtensionFollowsTheFormat(t *testing.T) {
	format := json.JsonLineFormat
	path := build(t, file.FileProvider{Session: "s", Format: &format})

	assert.Equal(t, ".jsonl", filepath.Ext(path))

	data, err := os.ReadFile(path)
	assert.Nil(t, err)
	assert.Contains(t, string(data), `"message":"hello"`)
	assert.False(t, strings.Contains(string(data), "\n\n"), "one line per record")
}

func TestFileProvider_UnusableFolderFails(t *testing.T) {
	blocker := filepath.Join(t.TempDir(), "file")
	assert.Nil(t, os.WriteFile(blocker, []byte("x"), 0o644))

	_, err := file.FileProvider{Path: filepath.Join(blocker, "logs")}.Build(t.Context())
	assert.NotNil(t, err, "a folder below a regular file cannot be created")
}
