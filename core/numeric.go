package core

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/passgen/config"
)

// NumericGenerator generates numeric-only passwords
type NumericGenerator struct {
	cfg *config.NumericConfig
	rng *rand.Rand
}

func NewNumericGenerator(cfg *config.NumericConfig) *NumericGenerator {
	gen := &NumericGenerator{cfg: cfg}
	if cfg.Randomized || cfg.Shuffle {
		seed := time.Now().UnixNano()
		if cfg.Deterministic && cfg.Seed != 0 {
			seed = cfg.Seed
		}
		gen.rng = rand.New(rand.NewSource(seed))
	}
	return gen
}

func (n *NumericGenerator) Generate(ctx context.Context) <-chan string {
	ch := make(chan string, 100)
	
	go func() {
		defer close(ch)
		
		if !n.cfg.Enabled {
			return
		}
		
		// Collect all numbers first if we need to shuffle
		var numbers []string
		if n.cfg.Shuffle {
			for i := n.cfg.Start; i <= n.cfg.End; i++ {
				numbers = append(numbers, n.formatNumber(i))
			}
			// Shuffle
			n.rng.Shuffle(len(numbers), func(i, j int) {
				numbers[i], numbers[j] = numbers[j], numbers[i]
			})
			// Output shuffled
			for _, num := range numbers {
				select {
				case <-ctx.Done():
					return
				case ch <- num:
				}
			}
		} else if n.cfg.Randomized {
			// Generate random numbers in range (limited to range size)
			rangeSize := n.cfg.End - n.cfg.Start + 1
			if rangeSize > 1000000 {
				rangeSize = 1000000 // Limit to prevent excessive generation
			}
			// Generate each number in range once, in random order
			numbers := make([]int64, rangeSize)
			for i := int64(0); i < rangeSize; i++ {
				numbers[i] = n.cfg.Start + i
			}
			n.rng.Shuffle(len(numbers), func(i, j int) {
				numbers[i], numbers[j] = numbers[j], numbers[i]
			})
			for _, num := range numbers {
				select {
				case <-ctx.Done():
					return
				case ch <- n.formatNumber(num):
				}
			}
		} else {
			// Sequential generation
			for i := n.cfg.Start; i <= n.cfg.End; i++ {
				select {
				case <-ctx.Done():
					return
				case ch <- n.formatNumber(i):
				}
			}
		}
	}()
	
	return ch
}

func (n *NumericGenerator) formatNumber(num int64) string {
	if n.cfg.Hexadecimal {
		// Format as hexadecimal
		hexStr := strconv.FormatInt(num, 16)
		if n.cfg.FixedLength > 0 {
			// Pad with zeros to fixed length
			padding := n.cfg.FixedLength - len(hexStr)
			if padding > 0 {
				hexStr = strings.Repeat("0", padding) + hexStr
			}
			return hexStr
		}
		if n.cfg.LeadingZeros {
			// Pad to match end number length
			maxHex := strconv.FormatInt(n.cfg.End, 16)
			maxLen := len(maxHex)
			if len(hexStr) < maxLen {
				padding := maxLen - len(hexStr)
				hexStr = strings.Repeat("0", padding) + hexStr
			}
			return hexStr
		}
		return hexStr
	}
	
	// Decimal formatting
	if n.cfg.FixedLength > 0 {
		return fmt.Sprintf("%0*d", n.cfg.FixedLength, num)
	}
	if n.cfg.LeadingZeros {
		maxLen := len(strconv.FormatInt(n.cfg.End, 10))
		return fmt.Sprintf("%0*d", maxLen, num)
	}
	return strconv.FormatInt(num, 10)
}

