package core

import (
	"context"
	"fmt"
	"strings"
)

// MaskAttackEngine generates passwords from Hashcat mask patterns
type MaskAttackEngine struct {
	masks []string
}

// NewMaskAttackEngine creates a new mask attack engine
func NewMaskAttackEngine(masks []string) *MaskAttackEngine {
	return &MaskAttackEngine{masks: masks}
}

// ExpandMask expands a Hashcat mask pattern to passwords
// ?l = lowercase, ?u = uppercase, ?d = digit, ?s = symbol, ?a = all, ?b = binary
func ExpandMask(mask string, maxCombinations int64) <-chan string {
	output := make(chan string, 100)

	go func() {
		defer close(output)
		expandMaskRecursive(mask, "", 0, output, maxCombinations)
	}()

	return output
}

// expandMaskRecursive recursively expands mask pattern
func expandMaskRecursive(mask, current string, depth int, output chan<- string, maxCombinations int64) {
	if maxCombinations > 0 && depth >= int(maxCombinations) {
		return
	}

	if len(mask) == 0 {
		if len(current) > 0 {
			output <- current
		}
		return
	}

	// Check for mask pattern
	if len(mask) >= 2 && mask[0] == '?' {
		var charset string
		switch mask[1] {
		case 'l': // lowercase
			charset = "abcdefghijklmnopqrstuvwxyz"
		case 'u': // uppercase
			charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		case 'd': // digit
			charset = "0123456789"
		case 's': // symbol
			charset = " !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
		case 'a': // all printable ASCII
			charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
		case 'b': // binary (0-255)
			// For binary, generate byte values
			for i := 0; i < 256; i++ {
				expandMaskRecursive(mask[2:], current+string(byte(i)), depth+1, output, maxCombinations)
			}
			return
		case 'h': // hex lowercase
			charset = "0123456789abcdef"
		case 'H': // hex uppercase
			charset = "0123456789ABCDEF"
		default:
			// Unknown mask, treat as literal
			expandMaskRecursive(mask[1:], current+string(mask[0]), depth, output, maxCombinations)
			return
		}

		// Expand each character in charset
		for _, char := range charset {
			expandMaskRecursive(mask[2:], current+string(char), depth+1, output, maxCombinations)
		}
	} else {
		// Literal character
		expandMaskRecursive(mask[1:], current+string(mask[0]), depth, output, maxCombinations)
	}
}

// GenerateFromMasks generates passwords from multiple masks
func (m *MaskAttackEngine) GenerateFromMasks(ctx context.Context) <-chan string {
	output := make(chan string, 100)

	go func() {
		defer close(output)
		seen := make(map[string]bool)

		for _, mask := range m.masks {
			select {
			case <-ctx.Done():
				return
			default:
				// Limit combinations to prevent explosion
				maxCombinations := int64(1000000) // 1M max per mask
				for password := range ExpandMask(mask, maxCombinations) {
					select {
					case <-ctx.Done():
						return
					default:
						// Deduplicate
						if !seen[password] {
							seen[password] = true
							output <- password
						}
					}
				}
			}
		}
	}()

	return output
}

// ValidateMask validates a mask pattern
func ValidateMask(mask string) error {
	if len(mask) == 0 {
		return fmt.Errorf("mask cannot be empty")
	}

	// Check for valid mask characters
	validChars := "?l?u?d?s?a?b?h?H"
	i := 0
	for i < len(mask) {
		if mask[i] == '?' {
			if i+1 >= len(mask) {
				return fmt.Errorf("incomplete mask pattern at position %d", i)
			}
			char := mask[i+1]
			if !strings.ContainsRune(validChars, rune(char)) && char != '?' {
				return fmt.Errorf("invalid mask character '%c' at position %d", char, i+1)
			}
			i += 2
		} else {
			i++
		}
	}

	return nil
}

// EstimateMaskSize estimates the number of passwords for a mask
func EstimateMaskSize(mask string) int64 {
	size := int64(1)

	i := 0
	for i < len(mask) {
		if mask[i] == '?' && i+1 < len(mask) {
			switch mask[i+1] {
			case 'l':
				size *= 26
			case 'u':
				size *= 26
			case 'd':
				size *= 10
			case 's':
				size *= 33 // Approximate
			case 'a':
				size *= 95 // All printable ASCII
			case 'b':
				size *= 256
			case 'h', 'H':
				size *= 16
			}
			i += 2
		} else {
			i++
		}

		// Prevent overflow
		if size < 0 {
			return -1
		}
	}

	return size
}

