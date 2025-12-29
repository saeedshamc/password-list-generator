package core

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
)

// JohnRule represents a John the Ripper rule
type JohnRule struct {
	Command string
	Param   string
}

// JohnRulesEngine applies John the Ripper rules to passwords
type JohnRulesEngine struct {
	rules []JohnRule
}

// NewJohnRulesEngine creates a new John rules engine
func NewJohnRulesEngine(rules []JohnRule) *JohnRulesEngine {
	return &JohnRulesEngine{rules: rules}
}

// LoadRulesFromFile loads John rules from a file
func LoadJohnRulesFromFile(filePath string) ([]JohnRule, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open rules file: %w", err)
	}
	defer file.Close()

	var rules []JohnRule
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "[") {
			continue
		}

		rule := parseJohnRule(line)
		if rule != nil {
			rules = append(rules, *rule)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read rules file: %w", err)
	}

	return rules, nil
}

// parseJohnRule parses a John rule line
func parseJohnRule(line string) *JohnRule {
	if len(line) == 0 {
		return nil
	}

	// John rules format: command:param or just command
	parts := strings.SplitN(line, ":", 2)
	command := parts[0]
	param := ""
	if len(parts) > 1 {
		param = parts[1]
	}

	return &JohnRule{
		Command: command,
		Param:   param,
	}
}

// ApplyRule applies a John rule to a password
func (j *JohnRulesEngine) ApplyRule(password string, rule JohnRule) string {
	result := password

	// Apply rule based on command (simplified John rules)
	switch rule.Command {
	case ":", "nothing":
		return result
	case "l", "lowercase":
		return strings.ToLower(result)
	case "u", "uppercase":
		return strings.ToUpper(result)
	case "c", "capitalize":
		if len(result) > 0 {
			return strings.ToUpper(string(result[0])) + strings.ToLower(result[1:])
		}
		return result
	case "C", "invertcase":
		runes := []rune(result)
		for i, r := range runes {
			if r >= 'a' && r <= 'z' {
				runes[i] = r - 32
			} else if r >= 'A' && r <= 'Z' {
				runes[i] = r + 32
			}
		}
		return string(runes)
	case "r", "reverse":
		runes := []rune(result)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes)
	case "d", "duplicate":
		return result + result
	case "f", "reflect":
		return result + reverseString(result)
	case "{", "rotateleft":
		if len(result) > 0 {
			return result[1:] + string(result[0])
		}
		return result
	case "}", "rotateright":
		if len(result) > 0 {
			return string(result[len(result)-1]) + result[:len(result)-1]
		}
		return result
	case "$", "append":
		if len(rule.Param) > 0 {
			return result + rule.Param
		}
		return result
	case "^", "prepend":
		if len(rule.Param) > 0 {
			return rule.Param + result
		}
		return result
	case "[", "deletefirst":
		if len(result) > 0 {
			return result[1:]
		}
		return result
	case "]", "deletelast":
		if len(result) > 0 {
			return result[:len(result)-1]
		}
		return result
	case "D", "deleteat":
		// Format: D:position
		if len(rule.Param) > 0 {
			pos := parsePosition(rule.Param)
			if pos >= 0 && pos < len(result) {
				runes := []rune(result)
				return string(append(runes[:pos], runes[pos+1:]...))
			}
		}
		return result
	case "x", "extract":
		// Format: x:start:length
		parts := strings.Split(rule.Param, ":")
		if len(parts) >= 2 {
			start := parsePosition(parts[0])
			length := parsePosition(parts[1])
			if start >= 0 && start+length <= len(result) {
				return result[start : start+length]
			}
		}
		return result
	case "i", "insert":
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
	case "o", "overwrite":
		// Format: o:position:character
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
	case "'", "increment":
		// Format: ':position
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
	case "s", "substitute":
		// Format: s:old:new
		parts := strings.Split(rule.Param, ":")
		if len(parts) >= 2 {
			old := parts[0]
			new := parts[1]
			return strings.ReplaceAll(result, old, new)
		}
		return result
	case "@", "purge":
		if len(rule.Param) > 0 {
			return strings.ReplaceAll(result, rule.Param, "")
		}
		return result
	case "z", "duplicatefirst":
		// Format: z:count
		if len(rule.Param) > 0 {
			n := parsePosition(rule.Param)
			if n > 0 && n <= len(result) {
				return result[:n] + result
			}
		}
		return result
	case "Z", "duplicatelast":
		// Format: Z:count
		if len(rule.Param) > 0 {
			n := parsePosition(rule.Param)
			if n > 0 && n <= len(result) {
				return result + result[len(result)-n:]
			}
		}
		return result
	case "q", "duplicateall":
		var builder strings.Builder
		for _, r := range result {
			builder.WriteRune(r)
			builder.WriteRune(r)
		}
		return builder.String()
	default:
		return result
	}
}

// ApplyRules applies all rules to a password
func (j *JohnRulesEngine) ApplyRules(password string) []string {
	results := []string{password}

	for _, rule := range j.rules {
		newResults := []string{}
		for _, pwd := range results {
			applied := j.ApplyRule(pwd, rule)
			if applied != "" {
				newResults = append(newResults, applied)
			}
		}
		results = newResults
	}

	return results
}

// ProcessPasswords processes passwords through rules
func (j *JohnRulesEngine) ProcessPasswords(ctx context.Context, passwords <-chan string) <-chan string {
	output := make(chan string, 100)

	go func() {
		defer close(output)
		for password := range passwords {
			select {
			case <-ctx.Done():
				return
			default:
				variants := j.ApplyRules(password)
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

