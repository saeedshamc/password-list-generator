package core

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// HashcatRule represents a single Hashcat rule
type HashcatRule struct {
	Command string
	Param   string
}

// HashcatRulesEngine applies Hashcat rules to passwords
type HashcatRulesEngine struct {
	rules []HashcatRule
}

// NewHashcatRulesEngine creates a new Hashcat rules engine
func NewHashcatRulesEngine(rules []HashcatRule) *HashcatRulesEngine {
	return &HashcatRulesEngine{rules: rules}
}

// LoadRulesFromFile loads Hashcat rules from a file
func LoadRulesFromFile(filePath string) ([]HashcatRule, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open rules file: %w", err)
	}
	defer file.Close()

	var rules []HashcatRule
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}
		
		rule := parseRule(line)
		if rule != nil {
			rules = append(rules, *rule)
		}
	}
	
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read rules file: %w", err)
	}
	
	return rules, nil
}

// parseRule parses a single rule line
func parseRule(line string) *HashcatRule {
	if len(line) == 0 {
		return nil
	}
	
	// Simple rule: command only or command:param
	parts := strings.SplitN(line, ":", 2)
	command := parts[0]
	param := ""
	if len(parts) > 1 {
		param = parts[1]
	}
	
	return &HashcatRule{
		Command: command,
		Param:   param,
	}
}

// ApplyRule applies a single rule to a password
func (h *HashcatRulesEngine) ApplyRule(password string, rule HashcatRule) string {
	result := password
	
	// Apply rule based on command
	switch rule.Command {
	case ":": // No operation
		return result
	case "l": // Lowercase all
		return strings.ToLower(result)
	case "u": // Uppercase all
		return strings.ToUpper(result)
	case "c": // Capitalize first, lowercase rest
		if len(result) > 0 {
			return strings.ToUpper(string(result[0])) + strings.ToLower(result[1:])
		}
		return result
	case "C": // Lowercase first, uppercase rest
		if len(result) > 0 {
			return strings.ToLower(string(result[0])) + strings.ToUpper(result[1:])
		}
		return result
	case "t": // Toggle case
		runes := []rune(result)
		for i, r := range runes {
			if unicode.IsLower(r) {
				runes[i] = unicode.ToUpper(r)
			} else if unicode.IsUpper(r) {
				runes[i] = unicode.ToLower(r)
			}
		}
		return string(runes)
	case "T": // Toggle case at position N
		if len(rule.Param) > 0 {
			pos := parsePosition(rule.Param)
			if pos >= 0 && pos < len(result) {
				runes := []rune(result)
				if unicode.IsLower(runes[pos]) {
					runes[pos] = unicode.ToUpper(runes[pos])
				} else if unicode.IsUpper(runes[pos]) {
					runes[pos] = unicode.ToLower(runes[pos])
				}
				return string(runes)
			}
		}
		return result
	case "r": // Reverse
		runes := []rune(result)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes)
	case "d": // Duplicate
		return result + result
	case "f": // Duplicate and reverse
		return result + reverseString(result)
	case "{": // Rotate left
		if len(result) > 0 {
			return result[1:] + string(result[0])
		}
		return result
	case "}": // Rotate right
		if len(result) > 0 {
			return string(result[len(result)-1]) + result[:len(result)-1]
		}
		return result
	case "$": // Append character
		if len(rule.Param) > 0 {
			return result + rule.Param
		}
		return result
	case "^": // Prepend character
		if len(rule.Param) > 0 {
			return rule.Param + result
		}
		return result
	case "[": // Delete first character
		if len(result) > 0 {
			return result[1:]
		}
		return result
	case "]": // Delete last character
		if len(result) > 0 {
			return result[:len(result)-1]
		}
		return result
	case "D": // Delete at position N
		if len(rule.Param) > 0 {
			pos := parsePosition(rule.Param)
			if pos >= 0 && pos < len(result) {
				runes := []rune(result)
				return string(append(runes[:pos], runes[pos+1:]...))
			}
		}
		return result
	case "x": // Extract range
		// Simplified: extract single character at position
		if len(rule.Param) > 0 {
			pos := parsePosition(rule.Param)
			if pos >= 0 && pos < len(result) {
				return string(result[pos])
			}
		}
		return result
	case "O": // Overwrite at position
		// Format: O:position:character
		parts := strings.Split(rule.Param, ":")
		if len(parts) >= 2 {
			pos := parsePosition(parts[0])
			char := parts[1]
			if pos >= 0 && pos < len(result) && len(char) > 0 {
				runes := []rune(result)
				runes[pos] = []rune(char)[0]
				return string(runes)
			}
		}
		return result
	case "i": // Insert at position
		// Format: i:position:character
		parts := strings.Split(rule.Param, ":")
		if len(parts) >= 2 {
			pos := parsePosition(parts[0])
			char := parts[1]
			if pos >= 0 && pos <= len(result) && len(char) > 0 {
				runes := []rune(result)
				newRunes := make([]rune, len(runes)+1)
				copy(newRunes, runes[:pos])
				newRunes[pos] = []rune(char)[0]
				copy(newRunes[pos+1:], runes[pos:])
				return string(newRunes)
			}
		}
		return result
	case "o": // Overwrite with character at position
		// Similar to O but different format
		if len(rule.Param) > 0 {
			pos := parsePosition(rule.Param)
			if pos >= 0 && pos < len(result) {
				// Use character from another position (simplified)
				return result
			}
		}
		return result
	case "'": // Increment character at position
		if len(rule.Param) > 0 {
			pos := parsePosition(rule.Param)
			if pos >= 0 && pos < len(result) {
				runes := []rune(result)
				if runes[pos] < 0xFFFF {
					runes[pos]++
				}
				return string(runes)
			}
		}
		return result
	case "s": // Substitute character
		// Format: s:old:new
		parts := strings.Split(rule.Param, ":")
		if len(parts) >= 2 {
			old := parts[0]
			new := parts[1]
			return strings.ReplaceAll(result, old, new)
		}
		return result
	case "@": // Purge character
		if len(rule.Param) > 0 {
			return strings.ReplaceAll(result, rule.Param, "")
		}
		return result
	case "z": // Duplicate first N characters
		if len(rule.Param) > 0 {
			n := parsePosition(rule.Param)
			if n > 0 && n <= len(result) {
				return result[:n] + result
			}
		}
		return result
	case "Z": // Duplicate last N characters
		if len(rule.Param) > 0 {
			n := parsePosition(rule.Param)
			if n > 0 && n <= len(result) {
				return result + result[len(result)-n:]
			}
		}
		return result
	case "q": // Duplicate all characters
		var builder strings.Builder
		for _, r := range result {
			builder.WriteRune(r)
			builder.WriteRune(r)
		}
		return builder.String()
	case "k": // Swap first two characters
		if len(result) >= 2 {
			runes := []rune(result)
			runes[0], runes[1] = runes[1], runes[0]
			return string(runes)
		}
		return result
	case "K": // Swap last two characters
		if len(result) >= 2 {
			runes := []rune(result)
			last := len(runes) - 1
			runes[last], runes[last-1] = runes[last-1], runes[last]
			return string(runes)
		}
		return result
	case "*": // Bitwise shift left
		// Simplified: not implemented
		return result
	case "L": // Bitwise shift right
		// Simplified: not implemented
		return result
	case "+": // ASCII increment
		if len(rule.Param) > 0 {
			pos := parsePosition(rule.Param)
			if pos >= 0 && pos < len(result) {
				runes := []rune(result)
				if runes[pos] < 0xFFFF {
					runes[pos]++
				}
				return string(runes)
			}
		}
		return result
	case "-": // ASCII decrement
		if len(rule.Param) > 0 {
			pos := parsePosition(rule.Param)
			if pos >= 0 && pos < len(result) {
				runes := []rune(result)
				if runes[pos] > 0 {
					runes[pos]--
				}
				return string(runes)
			}
		}
		return result
	case ".": // Replace with next character
		if len(rule.Param) > 0 {
			pos := parsePosition(rule.Param)
			if pos >= 0 && pos < len(result) {
				runes := []rune(result)
				if runes[pos] < 0xFFFF {
					runes[pos]++
				}
				return string(runes)
			}
		}
		return result
	case ",": // Replace with previous character
		if len(rule.Param) > 0 {
			pos := parsePosition(rule.Param)
			if pos >= 0 && pos < len(result) {
				runes := []rune(result)
				if runes[pos] > 0 {
					runes[pos]--
				}
				return string(runes)
			}
		}
		return result
	default:
		// Unknown rule, return original
		return result
	}
}

// ApplyRules applies all rules to a password
func (h *HashcatRulesEngine) ApplyRules(password string) []string {
	results := []string{password}
	
	for _, rule := range h.rules {
		newResults := []string{}
		for _, pwd := range results {
			applied := h.ApplyRule(pwd, rule)
			if applied != "" {
				newResults = append(newResults, applied)
			}
		}
		results = newResults
	}
	
	return results
}

// ProcessPasswords processes passwords through rules
func (h *HashcatRulesEngine) ProcessPasswords(ctx context.Context, passwords <-chan string) <-chan string {
	output := make(chan string, 100)
	
	go func() {
		defer close(output)
		for password := range passwords {
			select {
			case <-ctx.Done():
				return
			default:
				variants := h.ApplyRules(password)
				for _, variant := range variants {
					select {
					case <-ctx.Done():
						return
					case output <- variant:
					}
				}
			}
		}
	}()
	
	return output
}

// Helper functions
func parsePosition(s string) int {
	// Try to parse as integer
	var pos int
	fmt.Sscanf(s, "%d", &pos)
	return pos
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

