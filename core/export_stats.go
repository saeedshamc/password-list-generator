package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// StatisticsReport represents a statistics report
type StatisticsReport struct {
	GeneratedAt      time.Time              `json:"generated_at"`
	TotalPasswords   int64                  `json:"total_passwords"`
	TotalBytes       int64                  `json:"total_bytes"`
	TotalSizeMB      float64                `json:"total_size_mb"`
	PasswordsByLength map[int]int64         `json:"passwords_by_length"`
	PasswordsByType   map[string]int64      `json:"passwords_by_type"`
	AverageLength     float64                `json:"average_length"`
	MinLength         int                    `json:"min_length"`
	MaxLength         int                    `json:"max_length"`
}

// ExportStatisticsJSON exports statistics to JSON file
func ExportStatisticsJSON(stats *Statistics, filePath string) error {
	if stats == nil {
		return fmt.Errorf("statistics is nil")
	}

	totalPasswords, totalBytes, byLength, byType := stats.GetStats()

	// Calculate average length
	var totalLength int64
	var count int64
	var minLength, maxLength int = 999, 0

	for length, cnt := range byLength {
		totalLength += int64(length) * cnt
		count += cnt
		if length < minLength && cnt > 0 {
			minLength = length
		}
		if length > maxLength && cnt > 0 {
			maxLength = length
		}
	}

	avgLength := 0.0
	if count > 0 {
		avgLength = float64(totalLength) / float64(count)
	}

	report := StatisticsReport{
		GeneratedAt:      time.Now(),
		TotalPasswords:   totalPasswords,
		TotalBytes:       totalBytes,
		TotalSizeMB:      float64(totalBytes) / (1024 * 1024),
		PasswordsByLength: byLength,
		PasswordsByType:   byType,
		AverageLength:     avgLength,
		MinLength:         minLength,
		MaxLength:         maxLength,
	}

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write report: %w", err)
	}

	return nil
}

// ExportStatisticsHTML exports statistics to HTML file
func ExportStatisticsHTML(stats *Statistics, filePath string) error {
	if stats == nil {
		return fmt.Errorf("statistics is nil")
	}

	totalPasswords, totalBytes, byLength, byType := stats.GetStats()

	// Calculate average length
	var totalLength int64
	var count int64
	var minLength, maxLength int = 999, 0

	for length, cnt := range byLength {
		totalLength += int64(length) * cnt
		count += cnt
		if length < minLength && cnt > 0 {
			minLength = length
		}
		if length > maxLength && cnt > 0 {
			maxLength = length
		}
	}

	avgLength := 0.0
	if count > 0 {
		avgLength = float64(totalLength) / float64(count)
	}

	// Generate HTML
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<title>Password Wordlist Statistics</title>
	<style>
		body { font-family: Arial, sans-serif; margin: 20px; background: #f5f5f5; }
		.container { max-width: 1200px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
		h1 { color: #333; border-bottom: 2px solid #4CAF50; padding-bottom: 10px; }
		h2 { color: #555; margin-top: 30px; }
		.stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 20px; margin: 20px 0; }
		.stat-card { background: #f9f9f9; padding: 15px; border-radius: 5px; border-left: 4px solid #4CAF50; }
		.stat-label { font-weight: bold; color: #666; font-size: 0.9em; }
		.stat-value { font-size: 1.5em; color: #333; margin-top: 5px; }
		table { width: 100%%; border-collapse: collapse; margin: 20px 0; }
		th, td { padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }
		th { background-color: #4CAF50; color: white; }
		tr:hover { background-color: #f5f5f5; }
		.footer { margin-top: 30px; padding-top: 20px; border-top: 1px solid #ddd; color: #666; font-size: 0.9em; }
	</style>
</head>
<body>
	<div class="container">
		<h1>Password Wordlist Statistics Report</h1>
		<p><strong>Generated:</strong> %s</p>
		
		<div class="stats-grid">
			<div class="stat-card">
				<div class="stat-label">Total Passwords</div>
				<div class="stat-value">%d</div>
			</div>
			<div class="stat-card">
				<div class="stat-label">Total Size</div>
				<div class="stat-value">%.2f MB</div>
			</div>
			<div class="stat-card">
				<div class="stat-label">Average Length</div>
				<div class="stat-value">%.2f chars</div>
			</div>
			<div class="stat-card">
				<div class="stat-label">Length Range</div>
				<div class="stat-value">%d - %d chars</div>
			</div>
		</div>

		<h2>Distribution by Length</h2>
		<table>
			<thead>
				<tr>
					<th>Length</th>
					<th>Count</th>
					<th>Percentage</th>
				</tr>
			</thead>
			<tbody>`, time.Now().Format("2006-01-02 15:04:05"), totalPasswords, float64(totalBytes)/(1024*1024), avgLength, minLength, maxLength)

	// Add length distribution
	for length := minLength; length <= maxLength; length++ {
		if count, ok := byLength[length]; ok && count > 0 {
			percentage := float64(count) * 100.0 / float64(totalPasswords)
			html += fmt.Sprintf(`
				<tr>
					<td>%d</td>
					<td>%d</td>
					<td>%.2f%%</td>
				</tr>`, length, count, percentage)
		}
	}

	html += `
			</tbody>
		</table>

		<h2>Distribution by Type</h2>
		<table>
			<thead>
				<tr>
					<th>Type</th>
					<th>Count</th>
					<th>Percentage</th>
				</tr>
			</thead>
			<tbody>`

	// Add type distribution
	for ptype, count := range byType {
		if count > 0 {
			percentage := float64(count) * 100.0 / float64(totalPasswords)
			html += fmt.Sprintf(`
				<tr>
					<td>%s</td>
					<td>%d</td>
					<td>%.2f%%</td>
				</tr>`, ptype, count, percentage)
		}
	}

	html += fmt.Sprintf(`
			</tbody>
		</table>

		<div class="footer">
			<p>Report generated by Password Wordlist Generator</p>
			<p>Total bytes: %d | Total passwords: %d</p>
		</div>
	</div>
</body>
</html>`, totalBytes, totalPasswords)

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Write to file
	if err := os.WriteFile(filePath, []byte(html), 0644); err != nil {
		return fmt.Errorf("failed to write HTML report: %w", err)
	}

	return nil
}

