package gui

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/passgen/config"
	"github.com/passgen/core"
)

// App represents the GUI application
type App struct {
	fyneApp    fyne.App
	window     fyne.Window
	config     *config.Config
	generator  *core.Generator
	cancelFunc context.CancelFunc
}

// NewApp creates a new GUI application
func NewApp() *App {
	fyneApp := app.NewWithID("com.passgen.app")
	
	// Set Persian theme for better Persian font support
	fyneApp.Settings().SetTheme(NewPersianTheme())
	
	window := fyneApp.NewWindow("Password Wordlist Generator")
	window.Resize(fyne.NewSize(1000, 750))
	window.CenterOnScreen()
	
	app := &App{
		fyneApp: fyneApp,
		window:  window,
		config:  config.DefaultConfig(),
	}
	
	app.setupUI()
	return app
}

// Run starts the GUI application
func (a *App) Run() {
	a.window.ShowAndRun()
}

func (a *App) setupUI() {
	// Create tabs for different sections
	tabs := container.NewAppTabs()
	tabs.SetTabLocation(container.TabLocationTop)
	
	// Add tabs using the Items property
	inputTab := &container.TabItem{Text: "Input", Content: a.createInputTab()}
	numericTab := &container.TabItem{Text: "Numeric", Content: a.createNumericTab()}
	combinationTab := &container.TabItem{Text: "Combination", Content: a.createCombinationTab()}
	mutationTab := &container.TabItem{Text: "Mutation", Content: a.createMutationTab()}
	advancedTab := &container.TabItem{Text: "Advanced", Content: a.createAdvancedTab()}
	outputTab := &container.TabItem{Text: "Output", Content: a.createOutputTab()}
	
	tabs.Items = []*container.TabItem{inputTab, numericTab, combinationTab, mutationTab, advancedTab, outputTab}
	
	// Progress bar and controls
	progressBar := widget.NewProgressBar()
	progressBar.Hide()
	
	statusLabel := widget.NewLabel("Ready")
	
	// Generate button
	generateBtn := widget.NewButton("Generate Wordlist", func() {
		a.startGeneration(progressBar, statusLabel)
	})
	
	// Preview button
	previewBtn := widget.NewButton("Preview", func() {
		a.showPreview()
	})
	
	// Cancel button
	cancelBtn := widget.NewButton("Cancel", func() {
		if a.cancelFunc != nil {
			a.cancelFunc()
		}
	})
	cancelBtn.Disable()
	
	// Save/Load config buttons
	saveConfigBtn := widget.NewButton("💾 Save Config", func() {
		a.saveConfiguration()
	})
	
	loadConfigBtn := widget.NewButton("📂 Load Config", func() {
		a.loadConfiguration()
	})
	
	// Main layout with better spacing
	buttonBox := container.NewHBox(
		generateBtn,
		previewBtn,
		widget.NewSeparator(),
		saveConfigBtn,
		loadConfigBtn,
		widget.NewSeparator(),
		cancelBtn,
	)
	buttonBox.Resize(fyne.NewSize(300, 40))
	
	bottomPanel := container.NewBorder(
		nil, nil, nil, nil,
		container.NewVBox(
			container.NewPadded(buttonBox),
			container.NewPadded(progressBar),
			container.NewPadded(statusLabel),
		),
	)
	
	content := container.NewBorder(
		nil, bottomPanel, nil, nil,
		tabs,
	)
	
	a.window.SetContent(content)
}

func (a *App) startGeneration(progressBar *widget.ProgressBar, statusLabel *widget.Label) {
	// Create context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	a.cancelFunc = cancel
	
	// Update UI
	progressBar.Show()
	progressBar.SetValue(0)
	statusLabel.SetText("Generating...")
	
	// Create generator
	a.generator = core.NewGenerator(a.config)
	
		// Estimate size
		estLines, _ := a.generator.EstimateOutputSize()
		if estLines > a.config.Limits.WarnThreshold {
			statusLabel.SetText(fmt.Sprintf("Warning: Estimated %d lines. This may take a while.", estLines))
		}
	
	// Start generation in goroutine
	go func() {
		var lastUpdate time.Time
		var lastLines int64
		
		err := a.generator.Generate(ctx, func(lines int64, bytes int64) {
			// Update progress every 100ms
			if time.Since(lastUpdate) > 100*time.Millisecond {
				// Calculate progress (rough estimate)
				if a.config.Limits.MaxOutputLines > 0 {
					progress := float64(lines) / float64(a.config.Limits.MaxOutputLines)
					if progress > 1.0 {
						progress = 1.0
					}
					progressBar.SetValue(progress)
				}
				
				// Update status with statistics
				rate := float64(lines-lastLines) / time.Since(lastUpdate).Seconds()
				stats := a.generator.GetStatistics()
				if stats != nil {
					totalPasswords, totalBytes, _, _ := stats.GetStats()
					statusLabel.SetText(fmt.Sprintf("Generated: %d lines (%.0f lines/sec) | Total: %d passwords, %.2f MB", 
						lines, rate, totalPasswords, float64(totalBytes)/(1024*1024)))
				} else {
					statusLabel.SetText(fmt.Sprintf("Generated: %d lines (%.0f lines/sec)", lines, rate))
				}
				
				lastUpdate = time.Now()
				lastLines = lines
			}
		})
		
		if err != nil {
			if err == context.Canceled {
				statusLabel.SetText("Generation cancelled")
			} else {
				statusLabel.SetText(fmt.Sprintf("Error: %v", err))
			}
		} else {
			// Show final statistics
			stats := a.generator.GetStatistics()
			if stats != nil {
				totalPasswords, totalBytes, byLength, byType := stats.GetStats()
				// Show statistics dialog
				a.showStatisticsDialog(totalPasswords, totalBytes, byLength, byType)
			}
			statusLabel.SetText("Generation complete!")
			progressBar.SetValue(1.0)
		}
		
		progressBar.Hide()
		a.cancelFunc = nil
	}()
}

// Tab creation methods
// showHelpDialog shows help dialog for a tab
func (a *App) showHelpDialog(tabName string) {
	helpText := GetHelpText(tabName)
	
	// Use Label widget with monospace font for better formatting
	helpLabel := widget.NewLabel(helpText)
	helpLabel.Wrapping = fyne.TextWrapWord
	helpLabel.Alignment = fyne.TextAlignLeading
	
	// Create a container with padding for better spacing
	content := container.NewPadded(helpLabel)
	
	// Create scrollable container
	scroll := container.NewScroll(content)
	scroll.SetMinSize(fyne.NewSize(800, 600))
	
	dialog.ShowCustom("Help - "+tabName, "Close", scroll, a.window)
}

func (a *App) createInputTab() *container.Scroll {
	// Help button
	helpBtn := widget.NewButton("📖 Help", func() {
		a.showHelpDialog("input")
	})
	
	namesEntry := widget.NewMultiLineEntry()
	namesEntry.SetPlaceHolder("Enter names (one per line)\nExample: john, mary, admin")
	namesEntry.Resize(fyne.NewSize(0, 100))
	
	keywordsEntry := widget.NewMultiLineEntry()
	keywordsEntry.SetPlaceHolder("Enter keywords (one per line)\nExample: password, secret, key")
	keywordsEntry.Resize(fyne.NewSize(0, 100))
	
	freeTextEntry := widget.NewMultiLineEntry()
	freeTextEntry.SetPlaceHolder("Enter free text (one per line)\nExample: company, project, team")
	freeTextEntry.Resize(fyne.NewSize(0, 100))
	
	fileList := widget.NewList(
		func() int {
			return len(a.config.Input.ExternalFiles)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(a.config.Input.ExternalFiles[id])
		},
	)
	fileList.Resize(fyne.NewSize(0, 150))
	
	addFileBtn := widget.NewButton("📁 Add File", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				uri := reader.URI()
				if uri != nil {
					a.config.Input.ExternalFiles = append(a.config.Input.ExternalFiles, uri.Path())
					fileList.Refresh()
				}
				reader.Close()
			}
		}, a.window)
	})
	
	removeFileBtn := widget.NewButton("🗑️ Remove Last", func() {
		if len(a.config.Input.ExternalFiles) > 0 {
			a.config.Input.ExternalFiles = a.config.Input.ExternalFiles[:len(a.config.Input.ExternalFiles)-1]
			fileList.Refresh()
		}
	})
	
	updateBtn := widget.NewButton("✅ Update Input", func() {
		a.config.Input.Names = filterEmpty(strings.Split(namesEntry.Text, "\n"))
		a.config.Input.Keywords = filterEmpty(strings.Split(keywordsEntry.Text, "\n"))
		a.config.Input.FreeText = filterEmpty(strings.Split(freeTextEntry.Text, "\n"))
	})
	
	// Country selection - create checkboxes for each country
	allCountries := core.GetAllCountries()
	countryChecks := []*widget.Check{}
	countryContainer := container.NewVBox()
	
	for _, countryCode := range allCountries {
		countryCode := countryCode // Capture for closure
		displayName := core.GetCountryDisplayName(countryCode)
		
		// Check if country is selected
		isSelected := false
		for _, selected := range a.config.Input.SelectedCountries {
			if selected == countryCode {
				isSelected = true
				break
			}
		}
		
		check := widget.NewCheck(displayName, func(checked bool) {
			if checked {
				// Add country if not already selected
				found := false
				for _, selected := range a.config.Input.SelectedCountries {
					if selected == countryCode {
						found = true
						break
					}
				}
				if !found {
					a.config.Input.SelectedCountries = append(a.config.Input.SelectedCountries, countryCode)
				}
			} else {
				// Remove country
				newList := []string{}
				for _, selected := range a.config.Input.SelectedCountries {
					if selected != countryCode {
						newList = append(newList, selected)
					}
				}
				a.config.Input.SelectedCountries = newList
			}
		})
		check.SetChecked(isSelected)
		countryChecks = append(countryChecks, check)
		countryContainer.Add(check)
	}
	
	countryScroll := container.NewScroll(countryContainer)
	countryScroll.SetMinSize(fyne.NewSize(0, 200))
	
	selectAllCountriesBtn := widget.NewButton("✅ Select All", func() {
		a.config.Input.SelectedCountries = allCountries
		for _, check := range countryChecks {
			check.SetChecked(true)
		}
	})
	
	deselectAllCountriesBtn := widget.NewButton("❌ Deselect All", func() {
		a.config.Input.SelectedCountries = []string{}
		for _, check := range countryChecks {
			check.SetChecked(false)
		}
	})
	
	useAllNamesCheck := widget.NewCheck("Use All Names (or custom subset)", func(checked bool) {
		a.config.Input.UseAllCountryNames = checked
	})
	useAllNamesCheck.SetChecked(a.config.Input.UseAllCountryNames)
	
	// Advanced Date Generation
	includeBirthdays := widget.NewCheck("Include Common Birthdays (1980-2010)", func(checked bool) {
		a.config.Input.IncludeBirthdays = checked
	})
	includeBirthdays.SetChecked(a.config.Input.IncludeBirthdays)
	
	includeHolidays := widget.NewCheck("Include Holidays", func(checked bool) {
		a.config.Input.IncludeHolidays = checked
	})
	includeHolidays.SetChecked(a.config.Input.IncludeHolidays)
	
	includeSpecialDates := widget.NewCheck("Include Special Dates", func(checked bool) {
		a.config.Input.IncludeSpecialDates = checked
	})
	includeSpecialDates.SetChecked(a.config.Input.IncludeSpecialDates)
	
	content := container.NewVBox(
		container.NewHBox(helpBtn, widget.NewSeparator()),
		widget.NewCard("🌍 Country Names", "Select countries for name datasets",
			container.NewVBox(
				container.NewHBox(selectAllCountriesBtn, deselectAllCountriesBtn),
				container.NewPadded(countryScroll),
				container.NewPadded(useAllNamesCheck),
			),
		),
		widget.NewCard("📅 Advanced Dates", "Common dates for password generation",
			container.NewVBox(
				container.NewPadded(includeBirthdays),
				container.NewPadded(includeHolidays),
				container.NewPadded(includeSpecialDates),
			),
		),
		widget.NewCard("👤 Names", "First names, last names, usernames", 
			container.NewPadded(namesEntry)),
		widget.NewCard("🔑 Keywords", "Custom keywords and terms",
			container.NewPadded(keywordsEntry)),
		widget.NewCard("📝 Free Text", "Any custom text input",
			container.NewPadded(freeTextEntry)),
		widget.NewCard("📂 External Files", "Load words from text files",
			container.NewVBox(
				container.NewHBox(addFileBtn, removeFileBtn),
				container.NewPadded(fileList),
			),
		),
		container.NewPadded(updateBtn),
	)
	
	return container.NewScroll(content)
}

func filterEmpty(items []string) []string {
	var result []string
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func (a *App) createNumericTab() *container.Scroll {
	helpBtn := widget.NewButton("📖 Help", func() {
		a.showHelpDialog("numeric")
	})
	
	enabled := widget.NewCheck("🔢 Enable Numeric Mode", func(checked bool) {
		a.config.Numeric.Enabled = checked
	})
	enabled.SetChecked(a.config.Numeric.Enabled)
	
	startEntry := widget.NewEntry()
	startEntry.SetText(strconv.FormatInt(a.config.Numeric.Start, 10))
	startEntry.SetPlaceHolder("0")
	startEntry.OnChanged = func(s string) {
		if val, err := strconv.ParseInt(s, 10, 64); err == nil {
			a.config.Numeric.Start = val
		}
	}
	
	endEntry := widget.NewEntry()
	endEntry.SetText(strconv.FormatInt(a.config.Numeric.End, 10))
	endEntry.SetPlaceHolder("9999")
	endEntry.OnChanged = func(s string) {
		if val, err := strconv.ParseInt(s, 10, 64); err == nil {
			a.config.Numeric.End = val
		}
	}
	
	fixedLengthEntry := widget.NewEntry()
	fixedLengthEntry.SetText(strconv.Itoa(a.config.Numeric.FixedLength))
	fixedLengthEntry.SetPlaceHolder("0 = variable length")
	fixedLengthEntry.OnChanged = func(s string) {
		if val, err := strconv.Atoi(s); err == nil {
			a.config.Numeric.FixedLength = val
		}
	}
	
	leadingZeros := widget.NewCheck("0️⃣ Leading Zeros", func(checked bool) {
		a.config.Numeric.LeadingZeros = checked
	})
	leadingZeros.SetChecked(a.config.Numeric.LeadingZeros)
	
	randomized := widget.NewCheck("🎲 Randomized", func(checked bool) {
		a.config.Numeric.Randomized = checked
	})
	randomized.SetChecked(a.config.Numeric.Randomized)
	
	shuffle := widget.NewCheck("🔀 Shuffle", func(checked bool) {
		a.config.Numeric.Shuffle = checked
	})
	shuffle.SetChecked(a.config.Numeric.Shuffle)
	
	hexadecimal := widget.NewCheck("🔢 Hexadecimal", func(checked bool) {
		a.config.Numeric.Hexadecimal = checked
	})
	hexadecimal.SetChecked(a.config.Numeric.Hexadecimal)
	
	content := container.NewVBox(
		container.NewHBox(helpBtn, widget.NewSeparator()),
		container.NewPadded(enabled),
		widget.NewCard("📊 Range Settings", "Define the numeric range",
			container.NewPadded(
				widget.NewForm(
					widget.NewFormItem("Start Number", startEntry),
					widget.NewFormItem("End Number", endEntry),
					widget.NewFormItem("Fixed Length (0=variable)", fixedLengthEntry),
				),
			),
		),
		widget.NewCard("⚙️ Options", "Generation options",
			container.NewVBox(
				container.NewPadded(leadingZeros),
				container.NewPadded(randomized),
				container.NewPadded(shuffle),
				container.NewPadded(hexadecimal),
			),
		),
	)
	
	return container.NewScroll(content)
}

func (a *App) createCombinationTab() *container.Scroll {
	helpBtn := widget.NewButton("📖 Help", func() {
		a.showHelpDialog("combination")
	})
	
	// Password Templates section
	useTemplates := widget.NewCheck("📋 Use Password Templates", func(checked bool) {
		a.config.Combination.UseTemplates = checked
	})
	useTemplates.SetChecked(a.config.Combination.UseTemplates)
	
	// Get predefined templates
	predefinedTemplates := core.GetPredefinedTemplates()
	templateList := widget.NewList(
		func() int {
			return len(predefinedTemplates)
		},
		func() fyne.CanvasObject {
			return widget.NewCheck("", nil)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			check := obj.(*widget.Check)
			template := predefinedTemplates[id]
			check.SetText(fmt.Sprintf("%s: %s", template.Name, template.Template))
			
			// Check if template is selected
			isSelected := false
			for _, selected := range a.config.Combination.SelectedTemplates {
				if selected == template.Name {
					isSelected = true
					break
				}
			}
			check.SetChecked(isSelected)
			
			// Update on change
			check.OnChanged = func(checked bool) {
				if checked {
					// Add template if not already selected
					found := false
					for _, selected := range a.config.Combination.SelectedTemplates {
						if selected == template.Name {
							found = true
							break
						}
					}
					if !found {
						a.config.Combination.SelectedTemplates = append(a.config.Combination.SelectedTemplates, template.Name)
					}
				} else {
					// Remove template
					newList := []string{}
					for _, selected := range a.config.Combination.SelectedTemplates {
						if selected != template.Name {
							newList = append(newList, selected)
						}
					}
					a.config.Combination.SelectedTemplates = newList
				}
			}
		},
	)
	templateList.Resize(fyne.NewSize(0, 150))
	
	// Custom templates entry
	customTemplatesEntry := widget.NewMultiLineEntry()
	customTemplatesEntry.SetPlaceHolder("Enter custom templates (one per line)\nExample: {name}{year}, {keyword}!{number}")
	customTemplatesEntry.Resize(fyne.NewSize(0, 100))
	customTemplatesEntry.OnChanged = func(s string) {
		a.config.Combination.CustomTemplates = filterEmpty(strings.Split(s, "\n"))
	}
	
	enabled := widget.NewCheck("🔗 Enable Combinations", func(checked bool) {
		a.config.Combination.Enabled = checked
	})
	enabled.SetChecked(a.config.Combination.Enabled)
	
	maxDepthEntry := widget.NewEntry()
	maxDepthEntry.SetText(strconv.Itoa(a.config.Combination.MaxDepth))
	maxDepthEntry.SetPlaceHolder("3")
	maxDepthEntry.OnChanged = func(s string) {
		if val, err := strconv.Atoi(s); err == nil {
			a.config.Combination.MaxDepth = val
		}
	}
	
	minLengthEntry := widget.NewEntry()
	minLengthEntry.SetText(strconv.Itoa(a.config.Combination.MinLength))
	minLengthEntry.SetPlaceHolder("4")
	minLengthEntry.OnChanged = func(s string) {
		if val, err := strconv.Atoi(s); err == nil {
			a.config.Combination.MinLength = val
		}
	}
	
	maxLengthEntry := widget.NewEntry()
	maxLengthEntry.SetText(strconv.Itoa(a.config.Combination.MaxLength))
	maxLengthEntry.SetPlaceHolder("128")
	maxLengthEntry.OnChanged = func(s string) {
		if val, err := strconv.Atoi(s); err == nil {
			a.config.Combination.MaxLength = val
		}
	}
	
	wordWord := widget.NewCheck("📝 Word + Word", func(checked bool) {
		a.config.Combination.WordWord = checked
	})
	wordWord.SetChecked(a.config.Combination.WordWord)
	
	wordNumber := widget.NewCheck("🔢 Word + Number", func(checked bool) {
		a.config.Combination.WordNumber = checked
	})
	wordNumber.SetChecked(a.config.Combination.WordNumber)
	
	wordSymbolNumber := widget.NewCheck("✨ Word + Symbol + Number", func(checked bool) {
		a.config.Combination.WordSymbolNumber = checked
	})
	wordSymbolNumber.SetChecked(a.config.Combination.WordSymbolNumber)
	
	content := container.NewVBox(
		container.NewHBox(helpBtn, widget.NewSeparator()),
		container.NewPadded(enabled),
		widget.NewCard("📏 Length Settings", "Password length constraints",
			container.NewPadded(
				widget.NewForm(
					widget.NewFormItem("Max Depth", maxDepthEntry),
					widget.NewFormItem("Min Length", minLengthEntry),
					widget.NewFormItem("Max Length", maxLengthEntry),
				),
			),
		),
		widget.NewCard("🔀 Combination Types", "Select combination patterns",
			container.NewVBox(
				container.NewPadded(wordWord),
				container.NewPadded(wordNumber),
				container.NewPadded(wordSymbolNumber),
			),
		),
	)
	
	return container.NewScroll(content)
}

func (a *App) createFilterTab() *container.Scroll {
	helpBtn := widget.NewButton("📖 Help", func() {
		a.showHelpDialog("filter")
	})
	
	enabled := widget.NewCheck("✅ Enable Smart Filtering", func(checked bool) {
		a.config.Filter.Enabled = checked
	})
	enabled.SetChecked(a.config.Filter.Enabled)
	
	// Complexity selection
	complexityOptions := []string{"simple", "medium", "complex"}
	complexitySelect := widget.NewSelect(complexityOptions, func(selected string) {
		a.config.Filter.MinComplexity = selected
	})
	complexitySelect.SetSelected(a.config.Filter.MinComplexity)
	
	requireLetter := widget.NewCheck("Require Letter", func(checked bool) {
		a.config.Filter.RequireLetter = checked
	})
	requireLetter.SetChecked(a.config.Filter.RequireLetter)
	
	requireDigit := widget.NewCheck("Require Digit", func(checked bool) {
		a.config.Filter.RequireDigit = checked
	})
	requireDigit.SetChecked(a.config.Filter.RequireDigit)
	
	requireSymbol := widget.NewCheck("Require Symbol", func(checked bool) {
		a.config.Filter.RequireSymbol = checked
	})
	requireSymbol.SetChecked(a.config.Filter.RequireSymbol)
	
	onlyASCII := widget.NewCheck("ASCII Only", func(checked bool) {
		a.config.Filter.OnlyASCII = checked
	})
	onlyASCII.SetChecked(a.config.Filter.OnlyASCII)
	
	minUniqueCharsEntry := widget.NewEntry()
	minUniqueCharsEntry.SetText(strconv.Itoa(a.config.Filter.MinUniqueChars))
	minUniqueCharsEntry.SetPlaceHolder("0 = no limit")
	minUniqueCharsEntry.OnChanged = func(s string) {
		if val, err := strconv.Atoi(s); err == nil {
			a.config.Filter.MinUniqueChars = val
		}
	}
	
	excludePatternsEntry := widget.NewMultiLineEntry()
	excludePatternsEntry.SetPlaceHolder("Enter patterns to exclude (one per line)\nExample: password123, admin")
	excludePatternsEntry.Resize(fyne.NewSize(0, 100))
	excludePatternsEntry.OnChanged = func(s string) {
		a.config.Filter.ExcludePatterns = filterEmpty(strings.Split(s, "\n"))
	}
	
	content := container.NewVBox(
		container.NewHBox(helpBtn, widget.NewSeparator()),
		container.NewPadded(enabled),
		widget.NewCard("⚙️ Complexity Settings", "Password complexity requirements",
			container.NewVBox(
				container.NewPadded(widget.NewLabel("Minimum Complexity:")),
				container.NewPadded(complexitySelect),
				container.NewPadded(requireLetter),
				container.NewPadded(requireDigit),
				container.NewPadded(requireSymbol),
				container.NewPadded(onlyASCII),
				widget.NewForm(
					widget.NewFormItem("Min Unique Chars (0=no limit)", minUniqueCharsEntry),
				),
			),
		),
		widget.NewCard("🚫 Exclude Patterns", "Patterns to exclude from wordlist",
			container.NewPadded(excludePatternsEntry),
		),
	)
	
	return container.NewScroll(content)
}

func (a *App) createMutationTab() *container.Scroll {
	helpBtn := widget.NewButton("📖 Help", func() {
		a.showHelpDialog("mutation")
	})
	
	enabled := widget.NewCheck("🔄 Enable Mutations", func(checked bool) {
		a.config.Mutation.Enabled = checked
	})
	enabled.SetChecked(a.config.Mutation.Enabled)
	
	leetspeak := widget.NewCheck("💬 Leetspeak (a→@, e→3, i→1)", func(checked bool) {
		a.config.Mutation.Leetspeak = checked
	})
	leetspeak.SetChecked(a.config.Mutation.Leetspeak)
	
	unicode := widget.NewCheck("🌐 Unicode Substitutions", func(checked bool) {
		a.config.Mutation.Unicode = checked
	})
	unicode.SetChecked(a.config.Mutation.Unicode)
	
	intensitySelect := widget.NewSelect([]string{"low", "medium", "aggressive"}, func(s string) {
		a.config.Mutation.Intensity = s
	})
	intensitySelect.SetSelected(a.config.Mutation.Intensity)
	
	maxMutationsEntry := widget.NewEntry()
	maxMutationsEntry.SetText(strconv.Itoa(a.config.Mutation.MaxMutations))
	maxMutationsEntry.SetPlaceHolder("3")
	maxMutationsEntry.OnChanged = func(s string) {
		if val, err := strconv.Atoi(s); err == nil {
			a.config.Mutation.MaxMutations = val
		}
	}
	
	content := container.NewVBox(
		container.NewHBox(helpBtn, widget.NewSeparator()),
		container.NewPadded(enabled),
		widget.NewCard("🎨 Mutation Types", "Character substitution options",
			container.NewVBox(
				container.NewPadded(leetspeak),
				container.NewPadded(unicode),
			),
		),
		widget.NewCard("⚙️ Settings", "Mutation intensity and limits",
			container.NewPadded(
				widget.NewForm(
					widget.NewFormItem("Intensity", intensitySelect),
					widget.NewFormItem("Max Mutations", maxMutationsEntry),
				),
			),
		),
	)
	
	return container.NewScroll(content)
}

func (a *App) createAdvancedTab() *container.Scroll {
	helpBtn := widget.NewButton("📖 Help", func() {
		a.showHelpDialog("advanced")
	})
	
	caseVariations := widget.NewCheck("🔤 Case Variations (lower/upper/mixed)", func(checked bool) {
		a.config.Advanced.CaseVariations = checked
	})
	caseVariations.SetChecked(a.config.Advanced.CaseVariations)
	
	prefixesEntry := widget.NewMultiLineEntry()
	prefixesEntry.SetPlaceHolder("Enter prefixes (one per line)\nExample: !, @, #")
	prefixesEntry.Resize(fyne.NewSize(0, 100))
	prefixesEntry.OnChanged = func(s string) {
		a.config.Advanced.Prefixes = filterEmpty(strings.Split(s, "\n"))
	}
	
	suffixesEntry := widget.NewMultiLineEntry()
	suffixesEntry.SetPlaceHolder("Enter suffixes (one per line)\nExample: 123, !, 2024")
	suffixesEntry.Resize(fyne.NewSize(0, 100))
	suffixesEntry.OnChanged = func(s string) {
		a.config.Advanced.Suffixes = filterEmpty(strings.Split(s, "\n"))
	}
	
	repeatedChars := widget.NewCheck("🔁 Repeated Characters", func(checked bool) {
		a.config.Advanced.RepeatedChars = checked
	})
	repeatedChars.SetChecked(a.config.Advanced.RepeatedChars)
	
	wordDuplication := widget.NewCheck("📋 Word Duplication", func(checked bool) {
		a.config.Advanced.WordDuplication = checked
	})
	wordDuplication.SetChecked(a.config.Advanced.WordDuplication)
	
	content := container.NewVBox(
		container.NewHBox(helpBtn, widget.NewSeparator()),
		container.NewPadded(caseVariations),
		widget.NewCard("🔝 Prefixes", "Add prefixes to passwords",
			container.NewPadded(prefixesEntry)),
		widget.NewCard("🔚 Suffixes", "Add suffixes to passwords",
			container.NewPadded(suffixesEntry)),
		widget.NewCard("⚙️ Advanced Options", "Additional generation options",
			container.NewVBox(
				container.NewPadded(repeatedChars),
				container.NewPadded(wordDuplication),
			),
		),
	)
	
	return container.NewScroll(content)
}

func (a *App) createOutputTab() *container.Scroll {
	helpBtn := widget.NewButton("📖 Help", func() {
		a.showHelpDialog("output")
	})
	
	bruteforceBtn := widget.NewButton("🎯 Brute Force Guide", func() {
		a.showHelpDialog("bruteforce")
	})
	
	// Target Preset selection
	presetOptions := []string{"None", "WiFi (WPA/WPA2)", "SSH", "Web", "Mobile (PIN)", "IoT"}
	presetSelect := widget.NewSelect(presetOptions, func(selected string) {
		// Map UI selection to config preset
		switch selected {
		case "WiFi (WPA/WPA2)":
			a.config.TargetPreset.Preset = "wifi"
			a.config.TargetPreset.MinLength = 8
			a.config.TargetPreset.MaxLength = 63
			a.config.TargetPreset.ASCIIOnly = true
			a.config.TargetPreset.EnforceRules = true
		case "SSH":
			a.config.TargetPreset.Preset = "ssh"
			a.config.TargetPreset.MinLength = 4
			a.config.TargetPreset.MaxLength = 128
			a.config.TargetPreset.ASCIIOnly = true
			a.config.TargetPreset.EnforceRules = true
		case "Web":
			a.config.TargetPreset.Preset = "web"
			a.config.TargetPreset.MinLength = 4
			a.config.TargetPreset.MaxLength = 128
			a.config.TargetPreset.ASCIIOnly = true
			a.config.TargetPreset.EnforceRules = true
		case "Mobile (PIN)":
			a.config.TargetPreset.Preset = "mobile"
			a.config.TargetPreset.MinLength = 4
			a.config.TargetPreset.MaxLength = 8
			a.config.TargetPreset.ASCIIOnly = true
			a.config.TargetPreset.EnforceRules = true
		case "IoT":
			a.config.TargetPreset.Preset = "iot"
			a.config.TargetPreset.MinLength = 4
			a.config.TargetPreset.MaxLength = 32
			a.config.TargetPreset.ASCIIOnly = true
			a.config.TargetPreset.EnforceRules = true
		default:
			a.config.TargetPreset.Preset = "none"
			a.config.TargetPreset.MinLength = 4
			a.config.TargetPreset.MaxLength = 128
			a.config.TargetPreset.ASCIIOnly = false
			a.config.TargetPreset.EnforceRules = true
		}
	})
	
	// Set current preset
	currentPreset := "None"
	switch a.config.TargetPreset.Preset {
	case "wifi", "aircrack":
		currentPreset = "WiFi (WPA/WPA2)"
	case "ssh":
		currentPreset = "SSH"
	case "web":
		currentPreset = "Web"
	case "mobile", "pin":
		currentPreset = "Mobile (PIN)"
	case "iot":
		currentPreset = "IoT"
	}
	presetSelect.SetSelected(currentPreset)
	
	enforceRules := widget.NewCheck("✅ Enforce Validation Rules", func(checked bool) {
		a.config.TargetPreset.EnforceRules = checked
	})
	enforceRules.SetChecked(a.config.TargetPreset.EnforceRules)
	
	// Output directory selection
	outputDirEntry := widget.NewEntry()
	outputDirEntry.SetText("./")
	outputDirEntry.SetPlaceHolder("Select output directory...")
	outputDirEntry.OnChanged = func(s string) {
		// Directory will be set when user selects it
	}
	
	selectDirBtn := widget.NewButton("📁 Select Directory", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				dirPath := uri.Path()
				outputDirEntry.SetText(dirPath)
				// Update config with full path
				if a.config.Output.FilePath == "" || !strings.HasPrefix(a.config.Output.FilePath, dirPath) {
					a.config.Output.FilePath = dirPath + "/wordlist.txt"
				}
			}
		}, a.window)
	})
	
	// Output file name
	fileNameEntry := widget.NewEntry()
	fileNameEntry.SetText("wordlist.txt")
	fileNameEntry.SetPlaceHolder("wordlist.txt")
	fileNameEntry.OnChanged = func(s string) {
		if outputDirEntry.Text != "" {
			a.config.Output.FilePath = outputDirEntry.Text + "/" + s
		} else {
			a.config.Output.FilePath = s
		}
	}
	
	// Full path display
	fullPathLabel := widget.NewLabel("")
	fullPathLabel.Wrapping = fyne.TextWrapWord
	updateFullPath := func() {
		if outputDirEntry.Text != "" && fileNameEntry.Text != "" {
			fullPath := outputDirEntry.Text + "/" + fileNameEntry.Text
			a.config.Output.FilePath = fullPath
			fullPathLabel.SetText("📄 Full path: " + fullPath)
		} else {
			fullPathLabel.SetText("📄 Full path: " + fileNameEntry.Text)
			a.config.Output.FilePath = fileNameEntry.Text
		}
	}
	
	fileNameEntry.OnChanged = func(s string) {
		updateFullPath()
	}
	outputDirEntry.OnChanged = func(s string) {
		updateFullPath()
	}
	updateFullPath()
	
	deduplication := widget.NewCheck("✅ Enable Deduplication", func(checked bool) {
		a.config.Deduplication.Enabled = checked
	})
	deduplication.SetChecked(a.config.Deduplication.Enabled)
	
	maxLinesEntry := widget.NewEntry()
	maxLinesEntry.SetText(strconv.FormatInt(a.config.Limits.MaxOutputLines, 10))
	maxLinesEntry.SetPlaceHolder("0 = unlimited")
	maxLinesEntry.OnChanged = func(s string) {
		if val, err := strconv.ParseInt(s, 10, 64); err == nil {
			a.config.Limits.MaxOutputLines = val
		}
	}
	
	maxFileSizeEntry := widget.NewEntry()
	// Convert bytes to MB for display
	maxFileSizeMB := a.config.Limits.MaxFileSize / (1024 * 1024)
	maxFileSizeEntry.SetText(strconv.FormatInt(maxFileSizeMB, 10))
	maxFileSizeEntry.SetPlaceHolder("0 = unlimited (MB)")
	maxFileSizeEntry.OnChanged = func(s string) {
		if val, err := strconv.ParseInt(s, 10, 64); err == nil {
			// Convert MB to bytes
			a.config.Limits.MaxFileSize = val * 1024 * 1024
		}
	}
	
	// Username field (for hydra/john)
	usernameEntry := widget.NewEntry()
	usernameEntry.SetText(a.config.Output.Username)
	usernameEntry.SetPlaceHolder("Enter username (for hydra/john format)")
	usernameEntry.OnChanged = func(s string) {
		a.config.Output.Username = s
	}
	if a.config.Output.Format != "hydra" && a.config.Output.Format != "john" {
		usernameEntry.Hide()
	}
	
	// Hash field (for john)
	hashEntry := widget.NewEntry()
	hashEntry.SetText(a.config.Output.Hash)
	hashEntry.SetPlaceHolder("Enter hash (for john hash:password format)")
	hashEntry.OnChanged = func(s string) {
		a.config.Output.Hash = s
	}
	if a.config.Output.Format != "john" {
		hashEntry.Hide()
	}
	
	// Output format selection
	formatOptions := []string{"Plain", "Hashcat", "Hydra", "John", "Aircrack"}
	formatSelect := widget.NewSelect(formatOptions, func(selected string) {
		// Convert to lowercase for config
		a.config.Output.Format = strings.ToLower(selected)
		// Show/hide username/hash fields based on format
		if selected == "Hydra" || selected == "John" {
			usernameEntry.Show()
			if selected == "John" {
				hashEntry.Show()
			} else {
				hashEntry.Hide()
			}
		} else {
			usernameEntry.Hide()
			hashEntry.Hide()
		}
	})
	// Set default format
	if a.config.Output.Format == "" {
		a.config.Output.Format = "plain"
	}
	formatSelect.SetSelected(strings.Title(a.config.Output.Format))
	
	content := container.NewVBox(
		container.NewHBox(helpBtn, bruteforceBtn, widget.NewSeparator()),
		widget.NewCard("🎯 Target Preset", "Select target type for optimized generation",
			container.NewVBox(
				container.NewPadded(widget.NewLabel("Target Type:")),
				container.NewPadded(presetSelect),
				container.NewPadded(enforceRules),
			),
		),
		widget.NewCard("📂 Output Directory", "Select where to save the wordlist",
			container.NewVBox(
				container.NewHBox(outputDirEntry, selectDirBtn),
				container.NewPadded(fullPathLabel),
			),
		),
		widget.NewCard("📝 File Name", "Name of the output file",
			container.NewPadded(fileNameEntry),
		),
		widget.NewCard("🔧 Output Format", "Select format for cracking tools",
			container.NewVBox(
				container.NewPadded(widget.NewLabel("Format:")),
				container.NewPadded(formatSelect),
				container.NewPadded(usernameEntry),
				container.NewPadded(hashEntry),
			),
		),
		widget.NewCard("⚙️ Settings", "Generation settings",
			container.NewVBox(
				container.NewPadded(deduplication),
				widget.NewForm(
					widget.NewFormItem("Max Output Lines (0=unlimited)", maxLinesEntry),
					widget.NewFormItem("Max File Size (MB, 0=unlimited)", maxFileSizeEntry),
				),
			),
		),
	)
	
	return container.NewScroll(content)
}

