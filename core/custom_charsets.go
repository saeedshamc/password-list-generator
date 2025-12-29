package core

import (
	"context"
)

// CustomCharset represents a user-defined character set
type CustomCharset struct {
	Name        string
	Description string
	Characters  string
	Enabled     bool
}

// CustomCharsetEngine manages custom character sets
type CustomCharsetEngine struct {
	charsets map[string]*CustomCharset
}

// NewCustomCharsetEngine creates a new custom charset engine
func NewCustomCharsetEngine() *CustomCharsetEngine {
	return &CustomCharsetEngine{
		charsets: make(map[string]*CustomCharset),
	}
}

// AddCharset adds a custom character set
func (c *CustomCharsetEngine) AddCharset(charset *CustomCharset) {
	c.charsets[charset.Name] = charset
}

// GetCharset returns a charset by name
func (c *CustomCharsetEngine) GetCharset(name string) (*CustomCharset, bool) {
	charset, ok := c.charsets[name]
	return charset, ok
}

// ExpandCharset expands charset shortcuts (e.g., "a-z" -> "abcdefghijklmnopqrstuvwxyz")
func ExpandCharset(pattern string) string {
	result := ""
	i := 0
	
	for i < len(pattern) {
		if i+2 < len(pattern) && pattern[i+1] == '-' {
			// Range pattern like "a-z" or "0-9"
			start := rune(pattern[i])
			end := rune(pattern[i+2])
			
			if start <= end {
				for r := start; r <= end; r++ {
					result += string(r)
				}
			}
			i += 3
		} else {
			// Single character
			result += string(pattern[i])
			i++
		}
	}
	
	return result
}

// GetBuiltinCharsets returns built-in character sets
func GetBuiltinCharsets() map[string]string {
	return map[string]string{
		"lowercase":    "abcdefghijklmnopqrstuvwxyz",
		"uppercase":    "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"digits":       "0123456789",
		"symbols":      "!@#$%^&*()_+-=[]{}|;:,.<>?/~`",
		"alphanumeric": "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
		"ascii":        " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~",
		"hex":          "0123456789abcdefABCDEF",
		"base64":       "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/",
		"url-safe":     "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_",
		"printable":    " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~",
	}
}

// GenerateFromCharset generates passwords from a character set
func (c *CustomCharsetEngine) GenerateFromCharset(ctx context.Context, charsetName string, minLen, maxLen int) <-chan string {
	output := make(chan string, 100)
	
	go func() {
		defer close(output)
		
		var charset *CustomCharset
		var ok bool
		
		// Try custom charset first
		charset, ok = c.GetCharset(charsetName)
		if !ok {
			// Try built-in
			builtin := GetBuiltinCharsets()
			if chars, exists := builtin[charsetName]; exists {
				charset = &CustomCharset{
					Name:       charsetName,
					Characters: chars,
					Enabled:   true,
				}
			} else {
				return
			}
		}
		
		if !charset.Enabled || len(charset.Characters) == 0 {
			return
		}
		
		// Generate passwords of different lengths
		for length := minLen; length <= maxLen; length++ {
			select {
			case <-ctx.Done():
				return
			default:
				// Generate combinations (limited to prevent explosion)
				c.generateCombinations(ctx, output, charset.Characters, "", length, 0, 10000)
			}
		}
	}()
	
	return output
}

// generateCombinations generates all combinations of given length (with limit)
func (c *CustomCharsetEngine) generateCombinations(ctx context.Context, output chan<- string, charset, current string, length, depth, limit int) {
	if depth >= limit {
		return
	}
	
	if len(current) == length {
		select {
		case <-ctx.Done():
			return
		case output <- current:
		}
		return
	}
	
	for _, char := range charset {
		select {
		case <-ctx.Done():
			return
		default:
			c.generateCombinations(ctx, output, charset, current+string(char), length, depth+1, limit)
		}
	}
}

// ValidateCharset validates if a password uses only characters from charset
func (c *CustomCharsetEngine) ValidateCharset(password, charsetName string) bool {
	var charset *CustomCharset
	var ok bool
	
	// Try custom charset first
	charset, ok = c.GetCharset(charsetName)
	if !ok {
		// Try built-in
		builtin := GetBuiltinCharsets()
		if chars, exists := builtin[charsetName]; exists {
			charset = &CustomCharset{
				Name:       charsetName,
				Characters: chars,
				Enabled:   true,
			}
		} else {
			return false
		}
	}
	
	// Create a map for fast lookup
	charMap := make(map[rune]bool)
	for _, r := range charset.Characters {
		charMap[r] = true
	}
	
	// Check if all characters in password are in charset
	for _, r := range password {
		if !charMap[r] {
			return false
		}
	}
	
	return true
}

// FilterByCharset filters passwords to only those using specified charset
func (c *CustomCharsetEngine) FilterByCharset(ctx context.Context, passwords <-chan string, charsetName string) <-chan string {
	output := make(chan string, 100)
	
	go func() {
		defer close(output)
		
		for password := range passwords {
			select {
			case <-ctx.Done():
				return
			default:
				if c.ValidateCharset(password, charsetName) {
					output <- password
				}
			}
		}
	}()
	
	return output
}

// GetCharsetSize returns the number of possible passwords for a charset and length
func GetCharsetSize(charset string, length int) int64 {
	if length <= 0 || len(charset) == 0 {
		return 0
	}
	
	result := int64(1)
	for i := 0; i < length; i++ {
		result *= int64(len(charset))
		// Prevent overflow
		if result < 0 {
			return -1 // Overflow
		}
	}
	
	return result
}

