package core

import (
	"context"
	"fmt"
	"time"
)

// AdvancedDateGenerator generates common dates for password generation
type AdvancedDateGenerator struct {
	includeBirthdays bool
	includeHolidays  bool
	includeSpecial   bool
	yearRange        []int
}

// NewAdvancedDateGenerator creates a new advanced date generator
func NewAdvancedDateGenerator(includeBirthdays, includeHolidays, includeSpecial bool, startYear, endYear int) *AdvancedDateGenerator {
	// Generate year range
	years := []int{}
	for year := startYear; year <= endYear; year++ {
		years = append(years, year)
	}

	return &AdvancedDateGenerator{
		includeBirthdays: includeBirthdays,
		includeHolidays:   includeHolidays,
		includeSpecial:    includeSpecial,
		yearRange:         years,
	}
}

// Generate generates date strings
func (a *AdvancedDateGenerator) Generate(ctx context.Context) <-chan string {
	output := make(chan string, 100)

	go func() {
		defer close(output)

		// Common birthdays (1980-2010)
		if a.includeBirthdays {
			for year := 1980; year <= 2010; year++ {
				select {
				case <-ctx.Done():
					return
				default:
					// Common date formats
					formats := []string{
						fmt.Sprintf("%d0101", year), // YYYYMMDD
						fmt.Sprintf("%d1231", year), // YYYYMMDD
						fmt.Sprintf("01%02d%d", 1, year), // DDMMYYYY
						fmt.Sprintf("31%02d%d", 12, year), // DDMMYYYY
					}
					for _, date := range formats {
						output <- date
					}
				}
			}
		}

		// Common holidays and special dates
		if a.includeHolidays {
			holidays := []string{
				"0101", "0214", "0308", "0401", "0501", "0601",
				"0704", "0815", "0911", "1025", "1111", "1225",
				"1231",
			}
			for _, year := range a.yearRange {
				select {
				case <-ctx.Done():
					return
				default:
					for _, holiday := range holidays {
						// YYYYMMDD format
						output <- fmt.Sprintf("%d%s", year, holiday)
						// DDMMYYYY format
						output <- fmt.Sprintf("%s%d", holiday, year)
					}
				}
			}
		}

		// Special dates (common patterns)
		if a.includeSpecial {
			specialDates := []string{
				"20200101", "20201231", "20210101", "20211231",
				"20220101", "20221231", "20230101", "20231231",
				"20240101", "20241231",
				"01012020", "31122020", "01012021", "31122021",
				"01012022", "31122022", "01012023", "31122023",
				"01012024", "31122024",
			}
			for _, date := range specialDates {
				select {
				case <-ctx.Done():
					return
				default:
					output <- date
				}
			}
		}

		// Common years as dates
		for _, year := range a.yearRange {
			select {
			case <-ctx.Done():
				return
			default:
				// Just the year
				output <- fmt.Sprintf("%d", year)
				// Year with common suffixes
				output <- fmt.Sprintf("%d01", year)
				output <- fmt.Sprintf("%d12", year)
			}
		}
	}()

	return output
}

// GetCommonBirthdays returns common birthday years
func GetCommonBirthdays() []int {
	return []int{1980, 1985, 1990, 1995, 2000, 2005, 2010}
}

// GetCommonHolidays returns common holiday dates
func GetCommonHolidays() []string {
	return []string{
		"0101", // New Year
		"0214", // Valentine's
		"0308", // Women's Day
		"0401", // April Fool's
		"0501", // Labor Day
		"0601", // Children's Day
		"0704", // Independence Day (US)
		"0815", // Various holidays
		"0911", // 9/11
		"1025", // Various
		"1111", // Singles Day
		"1225", // Christmas
		"1231", // New Year's Eve
	}
}

// FormatDate formats a date in various formats
func FormatDate(date time.Time, format string) string {
	switch format {
	case "YYYY":
		return fmt.Sprintf("%d", date.Year())
	case "YYYYMMDD":
		return fmt.Sprintf("%d%02d%02d", date.Year(), date.Month(), date.Day())
	case "DDMMYYYY":
		return fmt.Sprintf("%02d%02d%d", date.Day(), date.Month(), date.Year())
	case "MMDDYYYY":
		return fmt.Sprintf("%02d%02d%d", date.Month(), date.Day(), date.Year())
	case "YYYYMM":
		return fmt.Sprintf("%d%02d", date.Year(), date.Month())
	case "MMYYYY":
		return fmt.Sprintf("%02d%d", date.Month(), date.Year())
	default:
		return date.Format("20060102")
	}
}

