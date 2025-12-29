package core

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/passgen/config"
)

// CombinationEngine generates password combinations
type CombinationEngine struct {
	cfg *config.CombinationConfig
	rng *rand.Rand
}

func NewCombinationEngine(cfg *config.CombinationConfig) *CombinationEngine {
	engine := &CombinationEngine{cfg: cfg}
	engine.rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	return engine
}

// Combine generates combinations from input words
func (c *CombinationEngine) Combine(ctx context.Context, words <-chan string) <-chan string {
	output := make(chan string, 100)
	
	go func() {
		defer close(output)
		
		if !c.cfg.Enabled {
			// Pass through unchanged
			for word := range words {
				select {
				case <-ctx.Done():
					return
				case output <- word:
				}
			}
			return
		}
		
		// Collect words into memory for combination
		wordList := []string{}
		for word := range words {
			wordList = append(wordList, word)
		}
		
		if len(wordList) == 0 {
			return
		}
		
		// Generate combinations
		c.generateCombinations(ctx, output, wordList, []string{}, 0)
	}()
	
	return output
}

func (c *CombinationEngine) generateCombinations(ctx context.Context, output chan<- string, words []string, current []string, depth int) {
	if depth >= c.cfg.MaxDepth {
		return
	}
	
	// Generate word + number combinations (RULE-BASED, not brute force)
	if c.cfg.WordNumber {
		for _, word := range words {
			select {
			case <-ctx.Done():
				return
			default:
				// Human behavior: prefer common numbers and years
				commonNumbers := []int{0, 1, 12, 123, 1234, 12345, 123456}
				years := []int{2024, 2023, 2022, 2021, 2020, 2019, 2018, 2017, 2016, 2015}
				
				// Add common numbers (realistic human behavior)
				for _, num := range commonNumbers {
					combo := fmt.Sprintf("%s%d", word, num)
					if c.isValidLength(combo) {
						output <- combo
					}
				}
				
				// Add years (very common in passwords)
				for _, year := range years {
					combo := fmt.Sprintf("%s%d", word, year)
					if c.isValidLength(combo) {
						output <- combo
					}
				}
				
				// Limited range: 0-999 (not 0-9999 to prevent explosion)
				for i := 0; i <= 999; i++ {
					combo := fmt.Sprintf("%s%d", word, i)
					if c.isValidLength(combo) {
						output <- combo
					}
					// Limit to prevent explosion
					if i >= 99 {
						break
					}
				}
			}
		}
	}
	
	// Generate word + word combinations
	if c.cfg.WordWord && len(words) > 1 {
		for i, w1 := range words {
			for j, w2 := range words {
				if i != j {
					select {
					case <-ctx.Done():
						return
					default:
						combos := []string{
							w1 + w2,
							w1 + "_" + w2,
							w1 + "-" + w2,
						}
						for _, combo := range combos {
							if c.isValidLength(combo) {
								output <- combo
							}
						}
					}
				}
			}
		}
	}
	
	// Generate word + symbol + number
	if c.cfg.WordSymbolNumber {
		symbols := []string{"!", "@", "#", "$", "%", "&", "*", "-", "_"}
		// Common years for brute force
		years := []int{2024, 2023, 2022, 2021, 2020, 2019, 2018}
		// Common numbers
		commonNumbers := []int{123, 1234, 12345, 123456, 1, 12, 123, 0, 00, 000, 0000}
		
		for _, word := range words {
			select {
			case <-ctx.Done():
				return
			default:
				// Word + symbol + year
				for _, sym := range symbols {
					for _, year := range years {
						combo := fmt.Sprintf("%s%s%d", word, sym, year)
						if c.isValidLength(combo) {
							output <- combo
						}
					}
				}
				
				// Word + symbol + common numbers
				for _, sym := range symbols {
					for _, num := range commonNumbers {
						combo := fmt.Sprintf("%s%s%d", word, sym, num)
						if c.isValidLength(combo) {
							output <- combo
						}
					}
				}
				
				// Word + symbol + limited range (0-99 to prevent explosion)
				for _, sym := range symbols {
					for i := 0; i <= 99; i++ {
						combo := fmt.Sprintf("%s%s%d", word, sym, i)
						if c.isValidLength(combo) {
							output <- combo
						}
					}
				}
			}
		}
	}
	
	// Apply patterns
	for _, pattern := range c.cfg.Patterns {
		if !pattern.Enabled {
			continue
		}
		c.applyPattern(ctx, output, words, pattern.Template)
	}
	
	// Multi-word chains
	if c.cfg.MultiWord && len(words) >= 2 {
		c.generateMultiWord(ctx, output, words, []string{}, 0)
	}
}

func (c *CombinationEngine) applyPattern(ctx context.Context, output chan<- string, words []string, template string) {
	// Simple pattern replacement
	// {word} -> word, {number} -> number, {symbol} -> symbol
	symbols := []string{"!", "@", "#", "$", "%", "&", "*", "-", "_"}
	// Extended numbers for brute force
	numbers := []string{"0", "1", "12", "123", "1234", "12345", "123456", 
		"2024", "2023", "2022", "2021", "2020", "2019", "2018",
		"00", "000", "0000", "00000", "000000"}
	
	for _, word := range words {
		select {
		case <-ctx.Done():
			return
		default:
			result := template
			result = strings.ReplaceAll(result, "{word}", word)
			
			// Replace {number} with various numbers
			if strings.Contains(result, "{number}") {
				for _, num := range numbers {
					combo := strings.ReplaceAll(result, "{number}", num)
					if c.isValidLength(combo) {
						output <- combo
					}
				}
			} else if strings.Contains(result, "{symbol}") {
				for _, sym := range symbols {
					combo := strings.ReplaceAll(result, "{symbol}", sym)
					if c.isValidLength(combo) {
						output <- combo
					}
				}
			} else {
				if c.isValidLength(result) {
					output <- result
				}
			}
		}
	}
}

func (c *CombinationEngine) generateMultiWord(ctx context.Context, output chan<- string, words []string, current []string, depth int) {
	if depth >= c.cfg.MaxDepth || len(current) >= 3 {
		if len(current) > 0 {
			combo := strings.Join(current, "")
			if c.isValidLength(combo) {
				select {
				case <-ctx.Done():
					return
				case output <- combo:
				}
			}
		}
		return
	}
	
	for _, word := range words {
		select {
		case <-ctx.Done():
			return
		default:
			newCurrent := append([]string{}, current...)
			newCurrent = append(newCurrent, word)
			c.generateMultiWord(ctx, output, words, newCurrent, depth+1)
		}
	}
}

func (c *CombinationEngine) isValidLength(s string) bool {
	length := len(s)
	return length >= c.cfg.MinLength && length <= c.cfg.MaxLength
}

