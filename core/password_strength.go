package core

import (
	"math"
	"strings"
	"unicode"
)

// PasswordStrength represents password strength metrics
type PasswordStrength struct {
	Score      int     // 0-100
	Entropy    float64 // Bits of entropy
	Category   string  // "Very Weak", "Weak", "Medium", "Strong", "Very Strong"
	Length     int
	HasUpper   bool
	HasLower   bool
	HasDigit   bool
	HasSymbol  bool
	UniqueChars int
}

// CalculateStrength calculates password strength
func CalculateStrength(password string) PasswordStrength {
	if len(password) == 0 {
		return PasswordStrength{
			Score:    0,
			Entropy:  0,
			Category: "Very Weak",
		}
	}

	strength := PasswordStrength{
		Length: len(password),
	}

	// Analyze character types
	charSet := make(map[rune]bool)
	for _, r := range password {
		charSet[r] = true
		if unicode.IsUpper(r) {
			strength.HasUpper = true
		}
		if unicode.IsLower(r) {
			strength.HasLower = true
		}
		if unicode.IsDigit(r) {
			strength.HasDigit = true
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsSpace(r) {
			strength.HasSymbol = true
		}
	}
	strength.UniqueChars = len(charSet)

	// Calculate entropy
	strength.Entropy = calculateEntropy(password, charSet)

	// Calculate score (0-100)
	strength.Score = calculateScore(strength)

	// Determine category
	strength.Category = getCategory(strength.Score)

	return strength
}

// calculateEntropy calculates Shannon entropy
func calculateEntropy(password string, charSet map[rune]bool) float64 {
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

// calculateScore calculates a strength score (0-100)
func calculateScore(strength PasswordStrength) int {
	score := 0

	// Length score (0-40 points)
	if strength.Length >= 12 {
		score += 40
	} else if strength.Length >= 8 {
		score += 30
	} else if strength.Length >= 6 {
		score += 20
	} else {
		score += 10
	}

	// Character variety (0-30 points)
	variety := 0
	if strength.HasUpper {
		variety++
	}
	if strength.HasLower {
		variety++
	}
	if strength.HasDigit {
		variety++
	}
	if strength.HasSymbol {
		variety++
	}
	score += variety * 7 // Max 28 points

	// Unique characters (0-20 points)
	uniqueRatio := float64(strength.UniqueChars) / float64(strength.Length)
	if uniqueRatio >= 0.8 {
		score += 20
	} else if uniqueRatio >= 0.6 {
		score += 15
	} else if uniqueRatio >= 0.4 {
		score += 10
	} else {
		score += 5
	}

	// Entropy bonus (0-10 points)
	if strength.Entropy >= 60 {
		score += 10
	} else if strength.Entropy >= 40 {
		score += 7
	} else if strength.Entropy >= 20 {
		score += 4
	} else {
		score += 1
	}

	// Penalties
	if strength.Length < 4 {
		score = 0
	} else if strength.Length < 6 {
		score = score / 2
	}

	// Cap at 100
	if score > 100 {
		score = 100
	}

	return score
}

// getCategory returns strength category
func getCategory(score int) string {
	if score >= 80 {
		return "Very Strong"
	} else if score >= 60 {
		return "Strong"
	} else if score >= 40 {
		return "Medium"
	} else if score >= 20 {
		return "Weak"
	}
	return "Very Weak"
}

// FilterByStrength filters passwords by minimum strength
func FilterByStrength(passwords <-chan string, minScore int) <-chan string {
	output := make(chan string, 100)

	go func() {
		defer close(output)
		for password := range passwords {
			strength := CalculateStrength(password)
			if strength.Score >= minScore {
				output <- password
			}
		}
	}()

	return output
}

// IsCommonPassword checks if password is in common list
func IsCommonPassword(password string) bool {
	common := []string{
		"password", "123456", "12345678", "1234", "qwerty",
		"abc123", "monkey", "1234567", "letmein", "trustno1",
		"dragon", "baseball", "iloveyou", "master", "sunshine",
		"ashley", "bailey", "passw0rd", "shadow", "123123",
		"654321", "superman", "qazwsx", "michael", "football",
		"welcome", "jesus", "ninja", "mustang", "password1",
	}

	lower := strings.ToLower(password)
	for _, common := range common {
		if lower == common || strings.Contains(lower, common) {
			return true
		}
	}

	return false
}

