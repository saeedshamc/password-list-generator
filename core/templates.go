package core

import (
	"context"
	"fmt"
	"strings"
)

// PasswordTemplate represents a password generation template
type PasswordTemplate struct {
	Name        string
	Description string
	Template    string // e.g., "{name}{year}", "{keyword}!{number}"
}

// Predefined templates
var PredefinedTemplates = []PasswordTemplate{
	{"Name + Year", "Combines name with year", "{name}{year}"},
	{"Name + Number", "Combines name with number", "{name}{number}"},
	{"Keyword + Year", "Combines keyword with year", "{keyword}{year}"},
	{"Keyword + Symbol + Number", "Combines keyword with symbol and number", "{keyword}{symbol}{number}"},
	{"Name + Symbol + Year", "Combines name with symbol and year", "{name}{symbol}{year}"},
	{"Name + Number + Symbol", "Combines name with number and symbol", "{name}{number}{symbol}"},
	{"Keyword + Number + Symbol", "Combines keyword with number and symbol", "{keyword}{number}{symbol}"},
	{"Name + Name", "Combines two names", "{name}{name}"},
	{"Name + Keyword", "Combines name with keyword", "{name}{keyword}"},
	{"Year + Name", "Year followed by name", "{year}{name}"},
	{"Number + Name", "Number followed by name", "{number}{name}"},
	{"Name + Date", "Name with date", "{name}{date}"},
}

// TemplateEngine generates passwords from templates
type TemplateEngine struct {
	templates []PasswordTemplate
	names     []string
	keywords  []string
	years     []int
	numbers   []int
	symbols   []string
	dates     []string
}

// NewTemplateEngine creates a new template engine
func NewTemplateEngine(templates []PasswordTemplate, names, keywords []string) *TemplateEngine {
	// Generate common years (2020-2024)
	years := []int{2024, 2023, 2022, 2021, 2020, 2019, 2018, 2017, 2016, 2015}
	
	// Generate common numbers
	numbers := []int{0, 1, 12, 123, 1234, 12345, 123456}
	
	// Common symbols
	symbols := []string{"!", "@", "#", "$", "%", "&", "*", "-", "_"}
	
	// Common dates (simplified)
	dates := []string{"0101", "1231", "2024", "2023"}
	
	return &TemplateEngine{
		templates: templates,
		names:     names,
		keywords:  keywords,
		years:     years,
		numbers:   numbers,
		symbols:   symbols,
		dates:     dates,
	}
}

// Generate generates passwords from templates
func (t *TemplateEngine) Generate(ctx context.Context) <-chan string {
	output := make(chan string, 100)
	
	go func() {
		defer close(output)
		
		for _, template := range t.templates {
			select {
			case <-ctx.Done():
				return
			default:
				passwords := t.expandTemplate(template.Template)
				for _, password := range passwords {
					select {
					case <-ctx.Done():
						return
					case output <- password:
					}
				}
			}
		}
	}()
	
	return output
}

// expandTemplate expands a template into passwords
func (t *TemplateEngine) expandTemplate(template string) []string {
	results := []string{}
	
	// Replace placeholders
	if strings.Contains(template, "{name}") {
		for _, name := range t.names {
			result := strings.ReplaceAll(template, "{name}", name)
			results = append(results, t.expandRemainingPlaceholders(result)...)
		}
	} else if strings.Contains(template, "{keyword}") {
		for _, keyword := range t.keywords {
			result := strings.ReplaceAll(template, "{keyword}", keyword)
			results = append(results, t.expandRemainingPlaceholders(result)...)
		}
	} else {
		// No name or keyword, expand directly
		results = append(results, t.expandRemainingPlaceholders(template)...)
	}
	
	return results
}

// expandRemainingPlaceholders expands remaining placeholders
func (t *TemplateEngine) expandRemainingPlaceholders(template string) []string {
	results := []string{template}
	
	// Expand {year}
	if strings.Contains(template, "{year}") {
		newResults := []string{}
		for _, result := range results {
			for _, year := range t.years {
				newResults = append(newResults, strings.ReplaceAll(result, "{year}", fmt.Sprintf("%d", year)))
			}
		}
		results = newResults
	}
	
	// Expand {number}
	if strings.Contains(template, "{number}") {
		newResults := []string{}
		for _, result := range results {
			for _, number := range t.numbers {
				newResults = append(newResults, strings.ReplaceAll(result, "{number}", fmt.Sprintf("%d", number)))
			}
		}
		results = newResults
	}
	
	// Expand {symbol}
	if strings.Contains(template, "{symbol}") {
		newResults := []string{}
		for _, result := range results {
			for _, symbol := range t.symbols {
				newResults = append(newResults, strings.ReplaceAll(result, "{symbol}", symbol))
			}
		}
		results = newResults
	}
	
	// Expand {date}
	if strings.Contains(template, "{date}") {
		newResults := []string{}
		for _, result := range results {
			for _, date := range t.dates {
				newResults = append(newResults, strings.ReplaceAll(result, "{date}", date))
			}
		}
		results = newResults
	}
	
	return results
}

// GetPredefinedTemplates returns all predefined templates
func GetPredefinedTemplates() []PasswordTemplate {
	return PredefinedTemplates
}

