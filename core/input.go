package core

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/passgen/config"
)

// InputSource represents a source of input words
type InputSource interface {
	Words(ctx context.Context) <-chan string
}

// AdvancedDateInput wraps AdvancedDateGenerator as InputSource
type AdvancedDateInput struct {
	generator *AdvancedDateGenerator
}

func (a *AdvancedDateInput) Words(ctx context.Context) <-chan string {
	return a.generator.Generate(ctx)
}

// NameInput generates words from names
type NameInput struct {
	names []string
}

func NewNameInput(names []string) *NameInput {
	return &NameInput{names: names}
}

func (n *NameInput) Words(ctx context.Context) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for _, name := range n.names {
			select {
			case <-ctx.Done():
				return
			case ch <- strings.TrimSpace(name):
			}
		}
	}()
	return ch
}

// KeywordInput generates words from keywords
type KeywordInput struct {
	keywords []string
}

func NewKeywordInput(keywords []string) *KeywordInput {
	return &KeywordInput{keywords: keywords}
}

func (k *KeywordInput) Words(ctx context.Context) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for _, keyword := range k.keywords {
			select {
			case <-ctx.Done():
				return
			case ch <- strings.TrimSpace(keyword):
			}
		}
	}()
	return ch
}

// NumericRangeInput generates numbers from ranges
type NumericRangeInput struct {
	ranges []config.NumericRange
}

func NewNumericRangeInput(ranges []config.NumericRange) *NumericRangeInput {
	return &NumericRangeInput{ranges: ranges}
}

func (n *NumericRangeInput) Words(ctx context.Context) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for _, r := range n.ranges {
			for i := r.Start; i <= r.End; i++ {
				select {
				case <-ctx.Done():
					return
				default:
					var numStr string
					if r.FixedLength > 0 {
						numStr = fmt.Sprintf("%0*d", r.FixedLength, i)
					} else if r.LeadingZeros {
						// Determine max length from end
						maxLen := len(strconv.FormatInt(r.End, 10))
						numStr = fmt.Sprintf("%0*d", maxLen, i)
					} else {
						numStr = strconv.FormatInt(i, 10)
					}
					ch <- numStr
				}
			}
		}
	}()
	return ch
}

// DateRangeInput generates dates from ranges
type DateRangeInput struct {
	ranges []config.DateRange
}

func NewDateRangeInput(ranges []config.DateRange) *DateRangeInput {
	return &DateRangeInput{ranges: ranges}
}

func (d *DateRangeInput) Words(ctx context.Context) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for _, r := range d.ranges {
			current := r.Start
			for !current.After(r.End) {
				select {
				case <-ctx.Done():
					return
				default:
					var dateStr string
					switch r.Format {
					case "YYYY":
						dateStr = strconv.Itoa(current.Year())
					case "YYYYMMDD":
						dateStr = current.Format("20060102")
					case "DDMMYYYY":
						dateStr = current.Format("02012006")
					case "MMDDYYYY":
						dateStr = current.Format("01022006")
					case "YY":
						dateStr = current.Format("06")
					case "YYMMDD":
						dateStr = current.Format("060102")
					default:
						dateStr = current.Format("2006-01-02")
					}
					ch <- dateStr
					current = current.AddDate(0, 0, 1)
				}
			}
		}
	}()
	return ch
}

// FileInput reads words from external files
type FileInput struct {
	files []string
}

func NewFileInput(files []string) *FileInput {
	return &FileInput{files: files}
}

func (f *FileInput) Words(ctx context.Context) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for _, filepath := range f.files {
			file, err := os.Open(filepath)
			if err != nil {
				continue
			}
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				select {
				case <-ctx.Done():
					file.Close()
					return
				default:
					line := strings.TrimSpace(scanner.Text())
					if line != "" {
						ch <- line
					}
				}
			}
			file.Close()
		}
	}()
	return ch
}

// FreeTextInput generates words from free text
type FreeTextInput struct {
	texts []string
}

func NewFreeTextInput(texts []string) *FreeTextInput {
	return &FreeTextInput{texts: texts}
}

func (f *FreeTextInput) Words(ctx context.Context) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for _, text := range f.texts {
			select {
			case <-ctx.Done():
				return
			case ch <- strings.TrimSpace(text):
			}
		}
	}()
	return ch
}

// CollectInputs collects all input sources into a single channel
func CollectInputs(ctx context.Context, cfg *config.Config) <-chan string {
	ch := make(chan string)
	
	go func() {
		defer close(ch)
		
		// Collect all input sources
		sources := []InputSource{}
		
		if len(cfg.Input.Names) > 0 {
			sources = append(sources, NewNameInput(cfg.Input.Names))
		}
		if len(cfg.Input.Keywords) > 0 {
			sources = append(sources, NewKeywordInput(cfg.Input.Keywords))
		}
		if len(cfg.Input.NumericRanges) > 0 {
			sources = append(sources, NewNumericRangeInput(cfg.Input.NumericRanges))
		}
		if len(cfg.Input.Dates) > 0 {
			sources = append(sources, NewDateRangeInput(cfg.Input.Dates))
		}
		if len(cfg.Input.ExternalFiles) > 0 {
			sources = append(sources, NewFileInput(cfg.Input.ExternalFiles))
		}
		if len(cfg.Input.FreeText) > 0 {
			sources = append(sources, NewFreeTextInput(cfg.Input.FreeText))
		}
		
		// Add country names if selected
		if len(cfg.Input.SelectedCountries) > 0 {
			countryNames := GetCountryNames(cfg.Input.SelectedCountries)
			if len(countryNames) > 0 {
				sources = append(sources, NewNameInput(countryNames))
			}
		}
		
		// Add advanced dates if enabled
		if cfg.Input.IncludeBirthdays || cfg.Input.IncludeHolidays || cfg.Input.IncludeSpecialDates {
			advancedDates := NewAdvancedDateGenerator(
				cfg.Input.IncludeBirthdays,
				cfg.Input.IncludeHolidays,
				cfg.Input.IncludeSpecialDates,
				2020, 2024,
			)
			// Create a wrapper to convert channel to InputSource
			dateSource := &AdvancedDateInput{generator: advancedDates}
			sources = append(sources, dateSource)
		}
		
		// Merge all sources
		done := make(chan bool, len(sources))
		for _, source := range sources {
			go func(s InputSource) {
				for word := range s.Words(ctx) {
					select {
					case <-ctx.Done():
						return
					case ch <- word:
					}
				}
				done <- true
			}(source)
		}
		
		// Wait for all sources to complete
		for i := 0; i < len(sources); i++ {
			select {
			case <-ctx.Done():
				return
			case <-done:
			}
		}
	}()
	
	return ch
}

