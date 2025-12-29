package core

import (
	"context"
	"math/rand"
	"time"

	"github.com/passgen/config"
)

// RandomizationEngine applies randomization to passwords
type RandomizationEngine struct {
	cfg *config.RandomizationConfig
	rng *rand.Rand
}

func NewRandomizationEngine(cfg *config.RandomizationConfig) *RandomizationEngine {
	engine := &RandomizationEngine{cfg: cfg}
	seed := time.Now().UnixNano()
	if cfg.Deterministic && cfg.Seed != 0 {
		seed = cfg.Seed
	}
	engine.rng = rand.New(rand.NewSource(seed))
	return engine
}

// Randomize applies randomization to passwords
func (r *RandomizationEngine) Randomize(ctx context.Context, input <-chan string) <-chan string {
	output := make(chan string, 100)
	
	go func() {
		defer close(output)
		
		if !r.cfg.Enabled {
			// Pass through unchanged
			for word := range input {
				select {
				case <-ctx.Done():
					return
				case output <- word:
				}
			}
			return
		}
		
		for word := range input {
			select {
			case <-ctx.Done():
				return
			default:
				variations := r.applyRandomization(word)
				for _, variant := range variations {
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

func (r *RandomizationEngine) applyRandomization(word string) []string {
	results := []string{word}
	
	if r.cfg.ShuffleChars {
		results = append(results, r.shuffleChars(word))
	}
	
	if r.cfg.InsertSymbols && len(r.cfg.Symbols) > 0 {
		results = append(results, r.insertSymbols(word)...)
	}
	
	return results
}

func (r *RandomizationEngine) shuffleChars(word string) string {
	runes := []rune(word)
	r.rng.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})
	return string(runes)
}

func (r *RandomizationEngine) insertSymbols(word string) []string {
	if len(r.cfg.Symbols) == 0 {
		return []string{word}
	}
	
	results := []string{}
	runes := []rune(word)
	symbols := []rune(r.cfg.Symbols)
	
	// Insert symbol at random positions
	for i := 0; i < len(runes); i++ {
		if r.rng.Float32() < 0.3 { // 30% chance
			sym := symbols[r.rng.Intn(len(symbols))]
			newWord := string(runes[:i]) + string(sym) + string(runes[i:])
			results = append(results, newWord)
		}
	}
	
	// Append symbol
	if len(symbols) > 0 {
		sym := symbols[r.rng.Intn(len(symbols))]
		results = append(results, word+string(sym))
	}
	
	// Prepend symbol
	if len(symbols) > 0 {
		sym := symbols[r.rng.Intn(len(symbols))]
		results = append(results, string(sym)+word)
	}
	
	return results
}

// SwapWordOrder swaps word order in multi-word passwords
func (r *RandomizationEngine) SwapWordOrder(words []string) []string {
	if len(words) < 2 {
		return words
	}
	
	result := make([]string, len(words))
	copy(result, words)
	r.rng.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})
	return result
}

