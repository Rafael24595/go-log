package file

import (
	"os"
	"path/filepath"
	"sync"
)

const maxRetries = 2

// File is a thread-safe wrapper around a standard os.File.
// It provides advanced features like automatic directory creation,
// file health checks (stat), and retry logic for writes.
type File struct {
	mu   sync.Mutex
	file *os.File
	name string
	flag int
	perm os.FileMode
}

// NewFile creates a new File instance. It automatically creates the
// required directory tree and opens the file with the specified flags.
func NewFile(name string, flag int, perm os.FileMode) (*File, error) {
	file, err := openFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return &File{
		file: file,
		name: name,
		flag: flag,
		perm: perm,
	}, nil
}

// Write persists bytes to the file. It ensures the file exists,
// handles retries, and calls Sync() to flush data to physical storage.
func (f *File) Write(bytes []byte) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if err := f.stat(); err != nil {
		return 0, err
	}

	return len(bytes), f.retry(func() error {
		_, err := f.file.Write(bytes)
		if err != nil {
			return err
		}
		return f.file.Sync()
	})
}

// Close gracefully releases the file descriptor.
func (f *File) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.file != nil {
		err := f.file.Close()
		f.file = nil
		return err
	}

	return nil
}

func (f *File) reopen() error {
	if f.file != nil {
		f.file.Close()
	}

	file, err := openFile(f.name, f.flag, f.perm)
	if err != nil {
		return err
	}

	f.file = file
	return nil
}

func (f *File) stat() error {
	_, err := os.Stat(f.name)
	if err == nil {
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}

	return f.reopen()
}

func (f *File) retry(fn func() error) error {
	var err error
	for range maxRetries {
		if err := fn(); err == nil {
			return nil
		}

		if err := f.reopen(); err != nil {
			return err
		}
	}

	return err
}

func openFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	dir := filepath.Dir(name)

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return os.OpenFile(name, flag, perm)
}
