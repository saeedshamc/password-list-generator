package core

import (
	"math"
	"strings"
	"unicode"
)

// PasswordAnalysis provides detailed analysis of a password
type PasswordAnalysis struct {
	Length        int
	HasUpper      bool
	HasLower      bool
	HasDigit      bool
	HasSymbol     bool
	UniqueChars   int
	CharTypes     int
	Entropy       float64
	CommonPattern bool
	KeyboardWalk  bool
	Repeating     bool
	Sequential    bool
}

// AnalyzePassword performs comprehensive password analysis
func AnalyzePassword(password string) PasswordAnalysis {
	if len(password) == 0 {
		return PasswordAnalysis{}
	}

	analysis := PasswordAnalysis{
		Length: len(password),
	}

	charSet := make(map[rune]bool)
	for _, r := range password {
		charSet[r] = true
		if unicode.IsUpper(r) {
			analysis.HasUpper = true
		}
		if unicode.IsLower(r) {
			analysis.HasLower = true
		}
		if unicode.IsDigit(r) {
			analysis.HasDigit = true
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsSpace(r) {
			analysis.HasSymbol = true
		}
	}

	analysis.UniqueChars = len(charSet)

	// Count character types
	types := 0
	if analysis.HasUpper {
		types++
	}
	if analysis.HasLower {
		types++
	}
	if analysis.HasDigit {
		types++
	}
	if analysis.HasSymbol {
		types++
	}
	analysis.CharTypes = types

	// Calculate entropy (using same function from password_strength.go)
	analysis.Entropy = calculateEntropyForAnalysis(password, charSet)

	// Check for common patterns
	analysis.CommonPattern = isCommonPattern(password)
	analysis.KeyboardWalk = isKeyboardWalk(password)
	analysis.Repeating = hasRepeatingPattern(password)
	analysis.Sequential = hasSequentialPattern(password)

	return analysis
}

// calculateEntropyForAnalysis calculates Shannon entropy for analysis
func calculateEntropyForAnalysis(password string, charSet map[rune]bool) float64 {
	if len(password) == 0 {
		return 0
	}

	// Count character frequencies
	freq := make(map[rune]int)
	for _, r := range password {
		freq[r]++
	}

	entropy := 0.0
	length := float64(len(password))

	for _, count := range freq {
		probability := float64(count) / length
		if probability > 0 {
			entropy -= probability * math.Log2(probability)
		}
	}

	// Multiply by length for total entropy
	return entropy * length
}

// isCommonPattern checks for common password patterns
func isCommonPattern(password string) bool {
	lower := strings.ToLower(password)
	common := []string{
		"password", "123456", "qwerty", "abc123", "admin",
		"letmein", "welcome", "monkey", "dragon", "master",
	}

	for _, pattern := range common {
		if strings.Contains(lower, pattern) {
			return true
		}
	}

	return false
}

// isKeyboardWalk checks for keyboard walking patterns
func isKeyboardWalk(password string) bool {
	lower := strings.ToLower(password)
	walks := []string{
		"qwerty", "asdf", "zxcv", "qaz", "wsx", "edc",
		"1234", "5678", "0987", "6543",
	}

	for _, walk := range walks {
		if strings.Contains(lower, walk) {
			return true
		}
	}

	return false
}

// hasRepeatingPattern checks for repeating characters
func hasRepeatingPattern(password string) bool {
	if len(password) < 3 {
		return false
	}

	for i := 0; i < len(password)-2; i++ {
		if password[i] == password[i+1] && password[i+1] == password[i+2] {
			return true
		}
	}

	return false
}

// hasSequentialPattern checks for sequential patterns
func hasSequentialPattern(password string) bool {
	if len(password) < 3 {
		return false
	}

	for i := 0; i < len(password)-2; i++ {
		r1 := rune(password[i])
		r2 := rune(password[i+1])
		r3 := rune(password[i+2])

		// Check ascending
		if r2 == r1+1 && r3 == r2+1 {
			return true
		}
		// Check descending
		if r2 == r1-1 && r3 == r2-1 {
			return true
		}
	}

	return false
}

// FilterByAnalysis filters passwords based on analysis criteria
func FilterByAnalysis(passwords <-chan string, minCharTypes, minUniqueChars int, excludeCommon, excludeKeyboardWalk, excludeRepeating, excludeSequential bool) <-chan string {
	output := make(chan string, 100)

	go func() {
		defer close(output)
		for password := range passwords {
			analysis := AnalyzePassword(password)

			// Check minimum character types
			if analysis.CharTypes < minCharTypes {
				continue
			}

			// Check minimum unique characters
			if analysis.UniqueChars < minUniqueChars {
				continue
			}

			// Exclude common patterns
			if excludeCommon && analysis.CommonPattern {
				continue
			}

			// Exclude keyboard walks
			if excludeKeyboardWalk && analysis.KeyboardWalk {
				continue
			}

			// Exclude repeating patterns
			if excludeRepeating && analysis.Repeating {
				continue
			}

			// Exclude sequential patterns
			if excludeSequential && analysis.Sequential {
				continue
			}

			output <- password
		}
	}()

	return output
}

