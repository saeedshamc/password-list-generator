package core

import (
	"context"
	"strings"
	"unicode"

	"github.com/passgen/config"
)

// SmartFilter applies intelligent filtering to passwords
type SmartFilter struct {
	cfg *config.FilterConfig
}

// NewSmartFilter creates a new smart filter
func NewSmartFilter(cfg *config.FilterConfig) *SmartFilter {
	return &SmartFilter{cfg: cfg}
}

// Filter checks if a password passes all filter criteria
func (f *SmartFilter) Filter(password string) bool {
	if !f.cfg.Enabled {
		return true
	}

	// Check complexity
	if !f.checkComplexity(password) {
		return false
	}

	// Check character requirements
	if f.cfg.RequireLetter && !hasLetter(password) {
		return false
	}
	if f.cfg.RequireDigit && !hasDigit(password) {
		return false
	}
	if f.cfg.RequireSymbol && !hasSymbol(password) {
		return false
	}

	// Check ASCII only
	if f.cfg.OnlyASCII && !isASCIIString(password) {
		return false
	}

	// Check minimum unique characters
	if f.cfg.MinUniqueChars > 0 {
		if countUniqueChars(password) < f.cfg.MinUniqueChars {
			return false
		}
	}

	// Check excluded patterns
	for _, pattern := range f.cfg.ExcludePatterns {
		if strings.Contains(strings.ToLower(password), strings.ToLower(pattern)) {
			return false
		}
	}

	return true
}

// checkComplexity checks if password meets minimum complexity requirement
func (f *SmartFilter) checkComplexity(password string) bool {
	switch f.cfg.MinComplexity {
	case "simple":
		return true // All passwords pass
	case "medium":
		// Must have at least 2 different character types
		hasLetter := hasLetter(password)
		hasDigit := hasDigit(password)
		hasSymbol := hasSymbol(password)
		
		types := 0
		if hasLetter {
			types++
		}
		if hasDigit {
			types++
		}
		if hasSymbol {
			types++
		}
		return types >= 2
	case "complex":
		// Must have all 3 character types and minimum length
		return hasLetter(password) && hasDigit(password) && hasSymbol(password) && len(password) >= 8
	default:
		return true
	}
}

// hasLetter checks if password contains letters
func hasLetter(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

// hasDigit checks if password contains digits
func hasDigit(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

// hasSymbol checks if password contains symbols
func hasSymbol(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsSpace(r) {
			return true
		}
	}
	return false
}

// countUniqueChars counts unique characters in password
func countUniqueChars(s string) int {
	seen := make(map[rune]bool)
	for _, r := range s {
		seen[r] = true
	}
	return len(seen)
}

// isASCIIString checks if string contains only ASCII characters
func isASCIIString(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// FilterPasswords filters passwords using smart filter
func (f *SmartFilter) FilterPasswords(ctx context.Context, passwords <-chan string) <-chan string {
	output := make(chan string, 100)

	go func() {
		defer close(output)
		for password := range passwords {
			select {
			case <-ctx.Done():
				return
			default:
				if f.Filter(password) {
					output <- password
				}
			}
		}
	}()

	return output
}

