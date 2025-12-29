package core

import (
	"context"
	"fmt"
	"strings"
)

// KeyboardLayout represents a keyboard layout
type KeyboardLayout struct {
	Name     string
	Rows     []string // Keyboard rows
	Adjacent map[rune][]rune // Adjacent keys
}

// KeyboardPatternsEngine generates keyboard-based patterns
type KeyboardPatternsEngine struct {
	layouts map[string]*KeyboardLayout
}

// NewKeyboardPatternsEngine creates a new keyboard patterns engine
func NewKeyboardPatternsEngine() *KeyboardPatternsEngine {
	engine := &KeyboardPatternsEngine{
		layouts: make(map[string]*KeyboardLayout),
	}
	
	// Initialize QWERTY layout
	engine.initQWERTY()
	
	return engine
}

// initQWERTY initializes QWERTY keyboard layout
func (k *KeyboardPatternsEngine) initQWERTY() {
	qwerty := &KeyboardLayout{
		Name: "QWERTY",
		Rows: []string{
			"qwertyuiop",
			"asdfghjkl",
			"zxcvbnm",
			"1234567890",
		},
		Adjacent: make(map[rune][]rune),
	}
	
	// Build adjacent key map
	adjacentMap := map[rune][]rune{
		'q': {'w'}, 'w': {'q', 'e'}, 'e': {'w', 'r'}, 'r': {'e', 't'}, 't': {'r', 'y'}, 'y': {'t', 'u'}, 'u': {'y', 'i'}, 'i': {'u', 'o'}, 'o': {'i', 'p'}, 'p': {'o'},
		'a': {'s'}, 's': {'a', 'd'}, 'd': {'s', 'f'}, 'f': {'d', 'g'}, 'g': {'f', 'h'}, 'h': {'g', 'j'}, 'j': {'h', 'k'}, 'k': {'j', 'l'}, 'l': {'k'},
		'z': {'x'}, 'x': {'z', 'c'}, 'c': {'x', 'v'}, 'v': {'c', 'b'}, 'b': {'v', 'n'}, 'n': {'b', 'm'}, 'm': {'n'},
		'1': {'2'}, '2': {'1', '3'}, '3': {'2', '4'}, '4': {'3', '5'}, '5': {'4', '6'}, '6': {'5', '7'}, '7': {'6', '8'}, '8': {'7', '9'}, '9': {'8', '0'}, '0': {'9'},
	}
	
	qwerty.Adjacent = adjacentMap
	k.layouts["qwerty"] = qwerty
}

// GenerateHorizontalPatterns generates horizontal keyboard patterns
func (k *KeyboardPatternsEngine) GenerateHorizontalPatterns(ctx context.Context, minLen, maxLen int) <-chan string {
	output := make(chan string, 100)
	
	go func() {
		defer close(output)
		
		qwerty := k.layouts["qwerty"]
		if qwerty == nil {
			return
		}
		
		// Generate patterns from each row
		for _, row := range qwerty.Rows {
			select {
			case <-ctx.Done():
				return
			default:
				// Generate all substrings of valid length
				for start := 0; start < len(row); start++ {
					for end := start + minLen; end <= len(row) && end-start <= maxLen; end++ {
						pattern := row[start:end]
						if len(pattern) >= minLen && len(pattern) <= maxLen {
							select {
							case <-ctx.Done():
								return
							case output <- pattern:
							}
							// Also generate reversed
							reversed := reverseString(pattern)
							select {
							case <-ctx.Done():
								return
							case output <- reversed:
							}
						}
					}
				}
			}
		}
	}()
	
	return output
}

// GenerateVerticalPatterns generates vertical keyboard patterns
func (k *KeyboardPatternsEngine) GenerateVerticalPatterns(ctx context.Context, minLen, maxLen int) <-chan string {
	output := make(chan string, 100)
	
	go func() {
		defer close(output)
		
		qwerty := k.layouts["qwerty"]
		if qwerty == nil {
			return
		}
		
		// Generate vertical patterns (columns)
		// Simplified: use common vertical patterns
		verticalPatterns := []string{
			"qaz", "wsx", "edc", "rfv", "tgb", "yhn", "ujm", "ik", "ol", "p",
			"1qaz", "2wsx", "3edc", "4rfv", "5tgb", "6yhn", "7ujm", "8ik", "9ol", "0p",
			"zaq", "xsw", "cde", "vfr", "bgt", "nhy", "mju", "ki", "lo", "p",
		}
		
		for _, pattern := range verticalPatterns {
			select {
			case <-ctx.Done():
				return
			default:
				if len(pattern) >= minLen && len(pattern) <= maxLen {
					output <- pattern
					// Also reversed
					output <- reverseString(pattern)
				}
			}
		}
	}()
	
	return output
}

// GenerateDiagonalPatterns generates diagonal keyboard patterns
func (k *KeyboardPatternsEngine) GenerateDiagonalPatterns(ctx context.Context, minLen, maxLen int) <-chan string {
	output := make(chan string, 100)
	
	go func() {
		defer close(output)
		
		// Common diagonal patterns
		diagonalPatterns := []string{
			"qwe", "asd", "zxc", "wer", "sdf", "xcv",
			"qaz", "wsx", "edc", "rfv", "tgb",
			"1qaz", "2wsx", "3edc", "4rfv", "5tgb",
		}
		
		for _, pattern := range diagonalPatterns {
			select {
			case <-ctx.Done():
				return
			default:
				if len(pattern) >= minLen && len(pattern) <= maxLen {
					output <- pattern
					output <- reverseString(pattern)
				}
			}
		}
	}()
	
	return output
}

// GenerateCommonPatterns generates common keyboard patterns
func (k *KeyboardPatternsEngine) GenerateCommonPatterns(ctx context.Context) <-chan string {
	output := make(chan string, 100)
	
	go func() {
		defer close(output)
		
		// Very common keyboard patterns
		commonPatterns := []string{
			"qwerty", "qwertyuiop", "asdf", "asdfgh", "asdfghjkl",
			"zxcv", "zxcvbnm", "qaz", "wsx", "edc", "rfv", "tgb",
			"1234", "12345", "123456", "1234567", "12345678", "123456789", "1234567890",
			"qwerty123", "asdf123", "qwertyuiop123",
			"1qaz2wsx", "1qaz2wsx3edc", "1qaz2wsx3edc4rfv",
			"qwerty!@#", "asdf!@#", "qwerty123!@#",
		}
		
		for _, pattern := range commonPatterns {
			select {
			case <-ctx.Done():
				return
			default:
				output <- pattern
				// Variations
				output <- strings.ToUpper(pattern)
				output <- strings.ToLower(pattern)
				// With numbers
				for i := 0; i <= 9; i++ {
					output <- pattern + fmt.Sprintf("%d", i)
					output <- fmt.Sprintf("%d", i) + pattern
				}
			}
		}
	}()
	
	return output
}

// ProcessPasswords adds keyboard patterns to passwords
func (k *KeyboardPatternsEngine) ProcessPasswords(ctx context.Context, passwords <-chan string) <-chan string {
	output := make(chan string, 100)
	
	go func() {
		defer close(output)
		
		// Generate common patterns
		patternCh := k.GenerateCommonPatterns(ctx)
		
		// Send patterns
		for pattern := range patternCh {
			select {
			case <-ctx.Done():
				return
			case output <- pattern:
			}
		}
		
		// Also combine with input passwords
		for password := range passwords {
			select {
			case <-ctx.Done():
				return
			default:
				output <- password
				// Add keyboard patterns as suffixes/prefixes
				commonSuffixes := []string{"qwerty", "asdf", "1234", "!@#$"}
				for _, suffix := range commonSuffixes {
					output <- password + suffix
					output <- suffix + password
				}
			}
		}
	}()
	
	return output
}

