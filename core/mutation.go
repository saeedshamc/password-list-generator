package core

import (
	"math/rand"
	"strings"
	"time"
	"unicode"

	"github.com/passgen/config"
)

// MutationEngine applies mutations to words
type MutationEngine struct {
	cfg *config.MutationConfig
	rng *rand.Rand
}

func NewMutationEngine(cfg *config.MutationConfig) *MutationEngine {
	engine := &MutationEngine{cfg: cfg}
	if cfg.Enabled {
		engine.rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	return engine
}

// Mutate applies mutations to a word based on configuration
func (m *MutationEngine) Mutate(word string) []string {
	if !m.cfg.Enabled {
		return []string{word}
	}
	
	results := []string{word}
	
	if m.cfg.Leetspeak {
		results = append(results, m.applyLeetspeak(word)...)
	}
	
	if m.cfg.Unicode {
		results = append(results, m.applyUnicode(word)...)
	}
	
	// Limit mutations
	if len(results) > m.cfg.MaxMutations+1 {
		results = results[:m.cfg.MaxMutations+1]
	}
	
	return results
}

// applyLeetspeak applies leetspeak substitutions
func (m *MutationEngine) applyLeetspeak(word string) []string {
	results := []string{}
	
	// Leetspeak mappings
	leetMap := map[rune][]string{
		'a': {"@", "4"},
		'A': {"@", "4"},
		'e': {"3"},
		'E': {"3"},
		'i': {"1", "!"},
		'I': {"1", "!"},
		'o': {"0"},
		'O': {"0"},
		's': {"$", "5"},
		'S': {"$", "5"},
		't': {"7"},
		'T': {"7"},
		'l': {"1"},
		'L': {"1"},
	}
	
	// Apply based on intensity
	switch m.cfg.Intensity {
	case "low":
		results = append(results, m.applyLeetLow(word, leetMap)...)
	case "medium":
		results = append(results, m.applyLeetMedium(word, leetMap)...)
	case "aggressive":
		results = append(results, m.applyLeetAggressive(word, leetMap)...)
	}
	
	return results
}

func (m *MutationEngine) applyLeetLow(word string, leetMap map[rune][]string) []string {
	results := []string{}
	runes := []rune(word)
	
	// Human behavior: Replace first occurrence of most common letters (a, e, i, o, s)
	// Prefer first letter or first vowel
	priority := []rune{'a', 'e', 'i', 'o', 's', 'A', 'E', 'I', 'O', 'S'}
	
	for _, priorityRune := range priority {
		for i, r := range runes {
			if r == priorityRune {
				if subs, ok := leetMap[r]; ok {
					// Use first substitution only (most common)
					if len(subs) > 0 {
						newWord := string(runes[:i]) + subs[0] + string(runes[i+1:])
						results = append(results, newWord)
					}
					return results // Only one replacement for low intensity
				}
			}
		}
	}
	
	return results
}

func (m *MutationEngine) applyLeetMedium(word string, leetMap map[rune][]string) []string {
	results := []string{}
	runes := []rune(word)
	
	// Human behavior: Replace 1-2 most common letters (a, e, i, o, s)
	// Prefer first substitution for each letter
	priority := []rune{'a', 'e', 'i', 'o', 's', 'A', 'E', 'I', 'O', 'S'}
	maxReplace := 2
	if len(word) > 6 {
		maxReplace = 3
	}
	
	replaced := 0
	for _, priorityRune := range priority {
		if replaced >= maxReplace {
			break
		}
		for i, r := range runes {
			if r == priorityRune {
				if subs, ok := leetMap[r]; ok && len(subs) > 0 {
					// Use first substitution (most common)
					newWord := string(runes[:i]) + subs[0] + string(runes[i+1:])
					results = append(results, newWord)
					replaced++
					if replaced >= maxReplace {
						return results
					}
					break
				}
			}
		}
	}
	
	return results
}

func (m *MutationEngine) applyLeetAggressive(word string, leetMap map[rune][]string) []string {
	results := []string{}
	runes := []rune(word)
	
	// Generate multiple combinations
	var generate func([]rune, int, []string)
	generate = func(rs []rune, pos int, current []string) {
		if pos >= len(rs) {
			results = append(results, strings.Join(current, ""))
			return
		}
		
		r := rs[pos]
		if subs, ok := leetMap[r]; ok && m.rng.Float32() < 0.7 {
			// Try substitution
			for _, sub := range subs {
				newCurrent := append([]string{}, current...)
				newCurrent = append(newCurrent, sub)
				generate(rs, pos+1, newCurrent)
			}
		} else {
			// Keep original
			newCurrent := append([]string{}, current...)
			newCurrent = append(newCurrent, string(r))
			generate(rs, pos+1, newCurrent)
		}
	}
	
	// Limit aggressive generation
	if len(word) <= 8 {
		generate(runes, 0, []string{})
		if len(results) > 20 {
			results = results[:20]
		}
	} else {
		// For longer words, use medium approach
		return m.applyLeetMedium(word, leetMap)
	}
	
	return results
}

// applyUnicode applies Unicode look-alike substitutions
func (m *MutationEngine) applyUnicode(word string) []string {
	results := []string{}
	
	unicodeMap := map[rune][]rune{
		'a': {'α', 'а', 'ᴀ'},
		'A': {'Α', 'А'},
		'e': {'е', '℮'},
		'E': {'Е'},
		'o': {'ο', 'о', '૦'},
		'O': {'Ο', 'О'},
		'c': {'с'},
		'C': {'С'},
		'p': {'р'},
		'P': {'Р'},
		'x': {'х'},
		'X': {'Х'},
	}
	
	runes := []rune(word)
	for i, r := range runes {
		if subs, ok := unicodeMap[r]; ok {
			for _, sub := range subs {
				newWord := string(runes[:i]) + string(sub) + string(runes[i+1:])
				results = append(results, newWord)
			}
		}
	}
	
	return results
}

// ApplyCaseVariations generates case variations
func ApplyCaseVariations(word string) []string {
	if word == "" {
		return []string{word}
	}
	
	results := []string{
		strings.ToLower(word),
		strings.ToUpper(word),
		strings.Title(strings.ToLower(word)),
	}
	
	// Mixed case variations
	runes := []rune(word)
	if len(runes) > 0 {
		// First letter uppercase, rest lowercase
		mixed1 := string(unicode.ToUpper(runes[0])) + strings.ToLower(string(runes[1:]))
		if !contains(results, mixed1) {
			results = append(results, mixed1)
		}
		
		// Alternating case
		if len(runes) > 1 {
			alt := make([]rune, len(runes))
			for i, r := range runes {
				if i%2 == 0 {
					alt[i] = unicode.ToUpper(r)
				} else {
					alt[i] = unicode.ToLower(r)
				}
			}
			altStr := string(alt)
			if !contains(results, altStr) {
				results = append(results, altStr)
			}
		}
	}
	
	return results
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

