package core

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

// CompressionWriter handles compressed output
type CompressionWriter struct {
	file   *os.File
	writer *gzip.Writer
}

// NewCompressionWriter creates a new compression writer
func NewCompressionWriter(filePath string) (*CompressionWriter, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	writer := gzip.NewWriter(file)

	return &CompressionWriter{
		file:   file,
		writer: writer,
	}, nil
}

// Write writes data to compressed file
func (c *CompressionWriter) Write(p []byte) (int, error) {
	return c.writer.Write(p)
}

// WriteString writes a string to compressed file
func (c *CompressionWriter) WriteString(s string) (int, error) {
	return c.writer.Write([]byte(s))
}

// Close closes the compression writer
func (c *CompressionWriter) Close() error {
	if err := c.writer.Close(); err != nil {
		c.file.Close()
		return err
	}
	return c.file.Close()
}

// Flush flushes the compression writer
func (c *CompressionWriter) Flush() error {
	return c.writer.Flush()
}

// IsCompressed checks if a file is gzip compressed
func IsCompressed(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	// Check magic number for gzip
	magic := make([]byte, 2)
	n, err := file.Read(magic)
	if err != nil || n < 2 {
		return false
	}

	return magic[0] == 0x1f && magic[1] == 0x8b
}

// DecompressFile decompresses a gzip file
func DecompressFile(inputPath, outputPath string) error {
	input, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer input.Close()

	reader, err := gzip.NewReader(input)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer reader.Close()

	output, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer output.Close()

	_, err = io.Copy(output, reader)
	return err
}

