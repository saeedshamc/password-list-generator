package core

import (
	"context"
	"regexp"
	"strings"
	"unicode"

	"github.com/passgen/config"
)

// CustomRulesEngine applies custom validation rules
type CustomRulesEngine struct {
	rules []config.CustomRule
}

// NewCustomRulesEngine creates a new custom rules engine
func NewCustomRulesEngine(rules []config.CustomRule) *CustomRulesEngine {
	return &CustomRulesEngine{rules: rules}
}

// Validate checks if password passes all enabled rules
func (e *CustomRulesEngine) Validate(password string) bool {
	for _, rule := range e.rules {
		if !rule.Enabled {
			continue
		}

		if !e.validateRule(password, rule) {
			return false
		}
	}

	return true
}

// validateRule validates a password against a single rule
func (e *CustomRulesEngine) validateRule(password string, rule config.CustomRule) bool {
	// Check length constraints
	if rule.MinLength > 0 && len(password) < rule.MinLength {
		return false
	}
	if rule.MaxLength > 0 && len(password) > rule.MaxLength {
		return false
	}

	// Check regex pattern
	if rule.Pattern != "" {
		matched, err := regexp.MatchString(rule.Pattern, password)
		if err != nil || !matched {
			return false
		}
	}

	// Check required characters
	if rule.Require != "" {
		if !e.checkRequired(password, rule.Require) {
			return false
		}
	}

	// Check excluded characters/patterns
	if rule.Exclude != "" {
		if e.checkExcluded(password, rule.Exclude) {
			return false
		}
	}

	return true
}

// checkRequired checks if password contains required character classes
func (e *CustomRulesEngine) checkRequired(password, require string) bool {
	requirements := strings.Split(require, ",")
	
	for _, req := range requirements {
		req = strings.TrimSpace(req)
		if !e.hasCharacterClass(password, req) {
			return false
		}
	}
	
	return true
}

// hasCharacterClass checks if password has characters from a class
func (e *CustomRulesEngine) hasCharacterClass(password, class string) bool {
	switch class {
	case "a-z", "lowercase":
		for _, r := range password {
			if r >= 'a' && r <= 'z' {
				return true
			}
		}
		return false
	case "A-Z", "uppercase":
		for _, r := range password {
			if r >= 'A' && r <= 'Z' {
				return true
			}
		}
		return false
	case "0-9", "digit", "digits":
		for _, r := range password {
			if r >= '0' && r <= '9' {
				return true
			}
		}
		return false
	case "!@#$%^&*", "symbol", "symbols":
		for _, r := range password {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsSpace(r) {
				return true
			}
		}
		return false
	default:
		// Check if it's a specific character set
		for _, r := range password {
			if strings.ContainsRune(class, r) {
				return true
			}
		}
		return false
	}
}

// checkExcluded checks if password contains excluded characters/patterns
func (e *CustomRulesEngine) checkExcluded(password, exclude string) bool {
	// Check if it's a pattern (contains regex-like syntax)
	if strings.Contains(exclude, "*") || strings.Contains(exclude, "+") || strings.Contains(exclude, "?") {
		matched, err := regexp.MatchString(exclude, password)
		if err == nil && matched {
			return true
		}
	}
	
	// Check for excluded characters
	for _, r := range password {
		if strings.ContainsRune(exclude, r) {
			return true
		}
	}
	
	return false
}

// FilterPasswords filters passwords using custom rules
func (e *CustomRulesEngine) FilterPasswords(ctx context.Context, passwords <-chan string) <-chan string {
	output := make(chan string, 100)

	go func() {
		defer close(output)
		for password := range passwords {
			select {
			case <-ctx.Done():
				return
			default:
				if e.Validate(password) {
					output <- password
				}
			}
		}
	}()

	return output
}

