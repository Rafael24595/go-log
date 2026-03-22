package file

import (
	"io"
	"os"
	"path/filepath"
)

// ReadFile opens a file, reads its entire content, and closes it.
// It automatically creates the directory path if it doesn't exist.
func ReadFile(filePath string) ([]byte, error) {
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return make([]byte, 0), err
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	result, readErr := io.ReadAll(file)
	err = file.Close()
	if err != nil {
		return make([]byte, 0), err
	}

	return result, readErr
}

// WriteFile writes a string to a file, truncating it if it already exists.
// It ensures the directory structure exists before writing.
func WriteFile(filePath, content string) error {
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	_, errWrite := file.Write([]byte(content))
	err = file.Close()
	if err != nil {
		return err
	}

	if errWrite != nil {
		return errWrite
	}

	return nil
}

// WriteFileSafe writes content to a temporary file first, then renames it 
// to the target path. This ensures atomicity: the target file is either 
// fully updated or remains unchanged if a crash occurs.
func WriteFileSafe(filePath, content string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp(dir, "temp-*.tmp")
	if err != nil {
		return err
	}

	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		tmpFile.Close()
		return err
	}

	if err := tmpFile.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmpPath, filePath); err != nil {
		return err
	}

	return tmpFile.Sync()
}

// AppendFileSafe appends content to a file, creating it if it doesn't exist.
// It uses Sync() to ensure the data is physically persisted to the disk.
func AppendFileSafe(filePath, content string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(content))
	if err != nil {
		return err
	}

	return f.Sync()
}
