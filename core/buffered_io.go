package core

import (
	"bufio"
	"io"
)

// BufferedWriter wraps an io.Writer with buffering
type BufferedWriter struct {
	writer *bufio.Writer
}

// NewBufferedWriter creates a new buffered writer
func NewBufferedWriter(w io.Writer) *BufferedWriter {
	return &BufferedWriter{
		writer: bufio.NewWriterSize(w, 64*1024), // 64KB buffer
	}
}

// Write writes data to the buffered writer
func (b *BufferedWriter) Write(p []byte) (int, error) {
	return b.writer.Write(p)
}

// WriteString writes a string to the buffered writer
func (b *BufferedWriter) WriteString(s string) (int, error) {
	return b.writer.WriteString(s)
}

// Flush flushes the buffer
func (b *BufferedWriter) Flush() error {
	return b.writer.Flush()
}

