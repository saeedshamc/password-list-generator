package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// FileSplitter handles splitting output into multiple files
type FileSplitter struct {
	basePath    string
	maxFileSize int64
	currentFile *os.File
	currentSize int64
	fileIndex   int
	writer      io.Writer
}

// NewFileSplitter creates a new file splitter
func NewFileSplitter(basePath string, maxFileSize int64) *FileSplitter {
	return &FileSplitter{
		basePath:    basePath,
		maxFileSize: maxFileSize,
		fileIndex:   1,
	}
}

// Write writes data and splits to new file if needed
func (f *FileSplitter) Write(p []byte) (int, error) {
	// Check if we need to create a new file
	if f.currentFile == nil || (f.maxFileSize > 0 && f.currentSize+int64(len(p)) > f.maxFileSize) {
		if err := f.nextFile(); err != nil {
			return 0, err
		}
	}

	n, err := f.writer.Write(p)
	if err != nil {
		return n, err
	}

	f.currentSize += int64(n)
	return n, nil
}

// WriteString writes a string
func (f *FileSplitter) WriteString(s string) (int, error) {
	return f.Write([]byte(s))
}

// nextFile creates the next file in sequence
func (f *FileSplitter) nextFile() error {
	// Close current file if exists
	if f.currentFile != nil {
		if bw, ok := f.writer.(*BufferedWriter); ok {
			bw.Flush()
		}
		f.currentFile.Close()
	}

	// Generate new filename
	ext := filepath.Ext(f.basePath)
	base := f.basePath[:len(f.basePath)-len(ext)]
	newPath := fmt.Sprintf("%s_part%d%s", base, f.fileIndex, ext)

	// Create new file
	file, err := os.Create(newPath)
	if err != nil {
		return fmt.Errorf("failed to create split file: %w", err)
	}

	f.currentFile = file
	f.writer = NewBufferedWriter(file)
	f.currentSize = 0
	f.fileIndex++

	return nil
}

// Close closes the current file
func (f *FileSplitter) Close() error {
	if f.currentFile != nil {
		if bw, ok := f.writer.(*BufferedWriter); ok {
			bw.Flush()
		}
		return f.currentFile.Close()
	}
	return nil
}

// Flush flushes the current file
func (f *FileSplitter) Flush() error {
	if f.writer != nil {
		if bw, ok := f.writer.(*BufferedWriter); ok {
			return bw.Flush()
		}
	}
	return nil
}

// GetCurrentFilePath returns the path of the current file
func (f *FileSplitter) GetCurrentFilePath() string {
	if f.currentFile != nil {
		return f.currentFile.Name()
	}
	return f.basePath
}

// GetFileCount returns the number of files created
func (f *FileSplitter) GetFileCount() int {
	return f.fileIndex - 1
}

