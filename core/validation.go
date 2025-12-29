package core

import (
	"strings"
	"unicode"
)

// PasswordValidator validates passwords based on target constraints and human behavior rules
type PasswordValidator struct {
	targetPreset string
	minLength    int
	maxLength    int
	asciiOnly    bool
}

// NewPasswordValidator creates a validator for a specific target
func NewPasswordValidator(targetPreset string, minLength, maxLength int, asciiOnly bool) *PasswordValidator {
	return &PasswordValidator{
		targetPreset: targetPreset,
		minLength:    minLength,
		maxLength:    maxLength,
		asciiOnly:    asciiOnly,
	}
}

// Validate checks if a password meets all constraints
func (v *PasswordValidator) Validate(password string) bool {
	// Check length constraints
	if len(password) < v.minLength || len(password) > v.maxLength {
		return false
	}

	// Check ASCII-only constraint
	if v.asciiOnly && !isASCII(password) {
		return false
	}

	// Human behavior rules
	if !v.isRealisticHumanPassword(password) {
		return false
	}

	return true
}

// isRealisticHumanPassword checks if password follows human behavior patterns
func (v *PasswordValidator) isRealisticHumanPassword(password string) bool {
	// Rule 1: No trivial repetitions (aaaaaa, 111111, etc.)
	if hasTrivialRepetition(password) {
		return false
	}

	// Rule 2: No meaningless sequences (asdfgh, qwertyui, etc.)
	if hasMeaninglessSequence(password) {
		return false
	}

	// Rule 3: No empty or whitespace-only
	if strings.TrimSpace(password) == "" {
		return false
	}

	// Rule 4: Must have at least some variation
	if !hasVariation(password) {
		return false
	}

	return true
}

// hasTrivialRepetition checks for patterns like aaaaaa, 111111, etc.
func hasTrivialRepetition(password string) bool {
	if len(password) < 3 {
		return false
	}

	runes := []rune(password)
	
	// Check for 4+ consecutive identical characters
	consecutive := 1
	lastChar := runes[0]
	for i := 1; i < len(runes); i++ {
		if runes[i] == lastChar {
			consecutive++
			if consecutive >= 4 {
				return true
			}
		} else {
			consecutive = 1
			lastChar = runes[i]
		}
	}

	// Check for simple patterns (aaa, 111, etc.)
	if len(runes) >= 4 {
		repeated := true
		for i := 1; i < len(runes) && i < 6; i++ {
			if runes[i] != runes[0] {
				repeated = false
				break
			}
		}
		if repeated {
			return true
		}
	}

	return false
}

// hasMeaninglessSequence checks for keyboard patterns
func hasMeaninglessSequence(password string) bool {
	lower := strings.ToLower(password)
	
	// Common keyboard sequences
	sequences := []string{
		"asdfgh", "qwerty", "qwertyui", "zxcvbn",
		"123456", "654321", "abcdef", "abcdefgh",
		"password", "passw0rd", "admin123",
	}

	for _, seq := range sequences {
		if strings.Contains(lower, seq) && len(password) <= len(seq)+2 {
			return true
		}
	}

	// Check for sequential patterns (abc, 123, etc.)
	if isSequential(password) {
		return true
	}

	return false
}

// isSequential checks for sequential character patterns
func isSequential(password string) bool {
	if len(password) < 4 {
		return false
	}

	runes := []rune(password)
	sequential := 0
	lastOrd := -1

	for i, r := range runes {
		ord := int(r)
		if i > 0 {
			diff := ord - lastOrd
			if diff == 1 || diff == -1 {
				sequential++
				if sequential >= 3 {
					return true
				}
			} else {
				sequential = 0
			}
		}
		lastOrd = ord
	}

	return false
}

// hasVariation checks if password has some character variation
func hasVariation(password string) bool {
	if len(password) < 2 {
		return true // Single character is fine
	}

	runes := []rune(password)
	firstChar := runes[0]
	allSame := true

	for _, r := range runes {
		if r != firstChar {
			allSame = false
			break
		}
	}

	return !allSame
}

// isASCII checks if string contains only ASCII characters
func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// FilterPasswords filters passwords based on validation rules
func FilterPasswords(validator *PasswordValidator, passwords <-chan string) <-chan string {
	output := make(chan string, 100)

	go func() {
		defer close(output)
		for password := range passwords {
			if validator.Validate(password) {
				output <- password
			}
		}
	}()

	return output
}

// GetTargetConstraints returns constraints for a specific target preset
func GetTargetConstraints(targetPreset string) (minLength, maxLength int, asciiOnly bool) {
	switch targetPreset {
	case "wifi", "aircrack":
		// WPA/WPA2: 8-63 characters, ASCII only
		return 8, 63, true
	case "ssh", "web":
		// SSH/Web: 4-128 characters, ASCII preferred
		return 4, 128, true
	case "mobile", "pin":
		// Mobile PIN: 4-8 digits, numeric only
		return 4, 8, true
	case "iot":
		// IoT devices: 4-32 characters, ASCII only
		return 4, 32, true
	default:
		// Default: 4-128 characters, any encoding
		return 4, 128, false
	}
}

