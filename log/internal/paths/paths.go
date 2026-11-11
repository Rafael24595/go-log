package paths

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// BootstrapDirEnv overrides the folder where records logged before any
// logger was configured are saved when the program exits.
const BootstrapDirEnv = "GO_LOG_DIR"

// BootstrapDir returns the folder for the unsigned bootstrap log:
// $GO_LOG_DIR when set; otherwise go-log/<program> inside the
// user's cache folder (%LocalAppData% on Windows, ~/.cache on Linux,
// ~/Library/Caches on macOS), or inside the temporary folder if there is none.
// It never depends on the working directory.
func BootstrapDir() string {
	if dir := os.Getenv(BootstrapDirEnv); dir != "" {
		return dir
	}

	base, err := os.UserCacheDir()
	if err != nil {
		base = os.TempDir()
	}

	return filepath.Join(base, "go-log", programName())
}

// programName is the running executable's name without its extension,
// so records from different programs do not mix.
func programName() string {
	file, err := os.Executable()
	if err != nil {
		return "unknown"
	}

	name := filepath.Base(file)
	
	ext := filepath.Ext(name)
	if runtime.GOOS == "windows" && strings.EqualFold(ext, ".exe") {
		name = strings.TrimSuffix(name, filepath.Ext(name))
	}

	name = strings.TrimSuffix(name, ext)
	return strings.ReplaceAll(name, ".", "_")
}
