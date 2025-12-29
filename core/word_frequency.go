package core

import (
	"context"
	"sort"
	"sync"
)

// WordFrequency tracks word frequency in generated passwords
type WordFrequency struct {
	Word      string
	Frequency int64
}

// FrequencyAnalyzer analyzes word frequency in passwords
type FrequencyAnalyzer struct {
	frequencies map[string]int64
	mu          sync.Mutex
}

// NewFrequencyAnalyzer creates a new frequency analyzer
func NewFrequencyAnalyzer() *FrequencyAnalyzer {
	return &FrequencyAnalyzer{
		frequencies: make(map[string]int64),
	}
}

// AddWord adds a word to frequency analysis
func (f *FrequencyAnalyzer) AddWord(word string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.frequencies[word]++
}

// GetTopWords returns top N most frequent words
func (f *FrequencyAnalyzer) GetTopWords(n int) []WordFrequency {
	f.mu.Lock()
	defer f.mu.Unlock()

	words := make([]WordFrequency, 0, len(f.frequencies))
	for word, freq := range f.frequencies {
		words = append(words, WordFrequency{Word: word, Frequency: freq})
	}

	// Sort by frequency (descending)
	sort.Slice(words, func(i, j int) bool {
		return words[i].Frequency > words[j].Frequency
	})

	// Return top N
	if n > len(words) {
		n = len(words)
	}
	return words[:n]
}

// GetFrequency returns frequency of a word
func (f *FrequencyAnalyzer) GetFrequency(word string) int64 {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.frequencies[word]
}

// AnalyzePasswords analyzes password patterns and extracts base words
func AnalyzePasswords(ctx context.Context, passwords <-chan string) <-chan string {
	output := make(chan string, 100)
	freqAnalyzer := NewFrequencyAnalyzer()

	go func() {
		defer close(output)

		for password := range passwords {
			select {
			case <-ctx.Done():
				return
			default:
				// Extract potential base words (remove numbers, symbols)
				baseWord := extractBaseWord(password)
				if baseWord != "" {
					freqAnalyzer.AddWord(baseWord)
				}

				output <- password
			}
		}
	}()

	return output
}

// extractBaseWord extracts the base word from a password
func extractBaseWord(password string) string {
	// Simple extraction: remove trailing numbers and symbols
	base := password
	for len(base) > 0 {
		last := base[len(base)-1]
		if (last >= '0' && last <= '9') || !((last >= 'a' && last <= 'z') || (last >= 'A' && last <= 'Z')) {
			base = base[:len(base)-1]
		} else {
			break
		}
	}
	return base
}

