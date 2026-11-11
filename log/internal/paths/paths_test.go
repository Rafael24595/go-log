package paths

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Rafael24595/go-assert/assert"
)

func TestBootstrapDir_EnvironmentOverride(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(BootstrapDirEnv, dir)

	assert.Equal(t, dir, BootstrapDir())
}

func TestBootstrapDir_DefaultIsPerProgramAndOutsideTheWorkingDirectory(t *testing.T) {
	t.Setenv(BootstrapDirEnv, "")

	got := BootstrapDir()

	assert.True(t, filepath.IsAbs(got))
	assert.True(t, strings.HasSuffix(got, filepath.Join("go-log", programName())))

	wd, err := os.Getwd()
	assert.Nil(t, err)

	assert.False(t, strings.HasPrefix(got, wd))
}

func TestProgramName_HasNoExtension(t *testing.T) {
	name := programName()

	assert.NotEqual(t, "", name)
	assert.Equal(t, "", filepath.Ext(name))
}
