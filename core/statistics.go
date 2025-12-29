package core

import (
	"sync"
	"sync/atomic"
)

// Statistics tracks generation statistics
type Statistics struct {
	TotalPasswords   int64
	TotalBytes       int64
	PasswordsByLength map[int]int64
	PasswordsByType  map[string]int64
	mu               sync.Mutex
}

// NewStatistics creates a new statistics tracker
func NewStatistics() *Statistics {
	return &Statistics{
		PasswordsByLength: make(map[int]int64),
		PasswordsByType:   make(map[string]int64),
	}
}

// AddPassword adds a password to statistics
func (s *Statistics) AddPassword(password string) {
	atomic.AddInt64(&s.TotalPasswords, 1)
	atomic.AddInt64(&s.TotalBytes, int64(len(password)+1)) // +1 for newline
	
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Track by length
	length := len(password)
	s.PasswordsByLength[length]++
	
	// Track by type (simple classification)
	passwordType := classifyPassword(password)
	s.PasswordsByType[passwordType]++
}

// GetStats returns current statistics
func (s *Statistics) GetStats() (totalPasswords, totalBytes int64, byLength map[int]int64, byType map[string]int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Create copies of maps
	byLengthCopy := make(map[int]int64)
	byTypeCopy := make(map[string]int64)
	
	for k, v := range s.PasswordsByLength {
		byLengthCopy[k] = v
	}
	for k, v := range s.PasswordsByType {
		byTypeCopy[k] = v
	}
	
	return atomic.LoadInt64(&s.TotalPasswords), atomic.LoadInt64(&s.TotalBytes), byLengthCopy, byTypeCopy
}

// classifyPassword classifies a password by type
func classifyPassword(password string) string {
	hasLetter := false
	hasDigit := false
	hasSymbol := false
	
	for _, r := range password {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' {
			hasLetter = true
		} else if r >= '0' && r <= '9' {
			hasDigit = true
		} else {
			hasSymbol = true
		}
	}
	
	if hasLetter && hasDigit && hasSymbol {
		return "Mixed (Letter+Digit+Symbol)"
	} else if hasLetter && hasDigit {
		return "Letter+Digit"
	} else if hasLetter && hasSymbol {
		return "Letter+Symbol"
	} else if hasDigit && hasSymbol {
		return "Digit+Symbol"
	} else if hasLetter {
		return "Letters Only"
	} else if hasDigit {
		return "Digits Only"
	} else if hasSymbol {
		return "Symbols Only"
	}
	return "Unknown"
}

// Reset resets all statistics
func (s *Statistics) Reset() {
	atomic.StoreInt64(&s.TotalPasswords, 0)
	atomic.StoreInt64(&s.TotalBytes, 0)
	
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.PasswordsByLength = make(map[int]int64)
	s.PasswordsByType = make(map[string]int64)
}

