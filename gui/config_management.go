package gui

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/passgen/config"
	"github.com/passgen/core"
	"context"
)

// saveConfiguration saves the current configuration
func (a *App) saveConfiguration() {
	// Show dialog to enter profile name
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Enter profile name...")
	
	content := container.NewVBox(
		widget.NewLabel("Save Configuration Profile"),
		widget.NewLabel("Enter a name for this configuration:"),
		nameEntry,
	)
	
	dialog.ShowCustomConfirm("Save Configuration", "Save", "Cancel", content, func(save bool) {
		if save {
			profileName := strings.TrimSpace(nameEntry.Text)
			if profileName == "" {
				dialog.ShowError(fmt.Errorf("profile name cannot be empty"), a.window)
				return
			}
			
			if err := config.SaveProfile(a.config, profileName); err != nil {
				dialog.ShowError(fmt.Errorf("failed to save profile: %v", err), a.window)
			} else {
				dialog.ShowInformation("Success", fmt.Sprintf("Configuration saved as '%s'", profileName), a.window)
			}
		}
	}, a.window)
}

// loadConfiguration loads a saved configuration
func (a *App) loadConfiguration() {
	profiles, err := config.ListSavedProfiles()
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to list profiles: %v", err), a.window)
		return
	}
	
	if len(profiles) == 0 {
		dialog.ShowInformation("No Profiles", "No saved profiles found. Save a configuration first.", a.window)
		return
	}
	
	// Create list of profiles
	profileList := widget.NewList(
		func() int {
			return len(profiles)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(profiles[id])
		},
	)
	
	var selectedProfile string
	profileList.OnSelected = func(id widget.ListItemID) {
		selectedProfile = profiles[id]
	}
	
	content := container.NewVBox(
		widget.NewLabel("Load Configuration Profile"),
		widget.NewLabel("Select a profile to load:"),
		container.NewScroll(profileList),
	)
	
	dialog.ShowCustomConfirm("Load Configuration", "Load", "Cancel", content, func(load bool) {
		if load && selectedProfile != "" {
			loadedConfig, err := config.LoadProfile(selectedProfile)
			if err != nil {
				dialog.ShowError(fmt.Errorf("failed to load profile: %v", err), a.window)
				return
			}
			
			// Update current config
			a.config = loadedConfig
			
			// Refresh UI - this is a simplified approach
			// In a full implementation, we'd need to refresh all UI elements
			dialog.ShowInformation("Success", fmt.Sprintf("Configuration '%s' loaded. Please restart the application for full effect.", selectedProfile), a.window)
		}
	}, a.window)
}

// showPreview shows a preview of generated passwords
func (a *App) showPreview() {
	// Create a context for preview (with timeout)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Create mutation engine for preview
	var mutationEngine *core.MutationEngine
	if a.config.Mutation.Enabled {
		mutationEngine = core.NewMutationEngine(&a.config.Mutation)
	}
	
	// Collect preview passwords (limit to 100)
	previewPasswords := []string{}
	previewCount := 0
	maxPreview := 100
	
	// Use a channel to collect passwords
	passwordCh := make(chan string, 100)
	
	go func() {
		defer close(passwordCh)
		
		// Collect inputs
		inputCh := core.CollectInputs(ctx, a.config)
		
		// Apply mutations
		mutatedCh := make(chan string, 100)
		go func() {
			defer close(mutatedCh)
			for word := range inputCh {
				select {
				case <-ctx.Done():
					return
				default:
					if a.config.Advanced.CaseVariations {
						variations := core.ApplyCaseVariations(word)
						for _, variant := range variations {
							if previewCount >= maxPreview {
								return
							}
							mutatedCh <- variant
							previewCount++
						}
					} else {
						if previewCount >= maxPreview {
							return
						}
						mutatedCh <- word
						previewCount++
					}
					
					if a.config.Mutation.Enabled && mutationEngine != nil {
						mutations := mutationEngine.Mutate(word)
						for _, mut := range mutations {
							if mut != word && previewCount < maxPreview {
								select {
								case <-ctx.Done():
									return
								case mutatedCh <- mut:
									previewCount++
								}
							}
						}
					}
				}
			}
		}()
		
		// Collect from mutated channel
		for word := range mutatedCh {
			select {
			case <-ctx.Done():
				return
			case passwordCh <- word:
			}
		}
	}()
	
	// Collect passwords
	for word := range passwordCh {
		if len(previewPasswords) >= maxPreview {
			break
		}
		previewPasswords = append(previewPasswords, word)
	}
	
	// Display preview
	if len(previewPasswords) == 0 {
		dialog.ShowInformation("Preview", "No passwords generated. Check your input settings.", a.window)
		return
	}
	
	// Create preview content
	previewText := strings.Join(previewPasswords, "\n")
	previewEntry := widget.NewMultiLineEntry()
	previewEntry.SetText(previewText)
	previewEntry.Disable() // Read-only
	previewEntry.Wrapping = fyne.TextWrapOff
	
	scroll := container.NewScroll(previewEntry)
	scroll.SetMinSize(fyne.NewSize(600, 400))
	
	content := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("Preview (%d passwords):", len(previewPasswords))),
		scroll,
	)
	
	dialog.ShowCustom("Password Preview", "Close", content, a.window)
}

// showStatisticsDialog displays generation statistics
func (a *App) showStatisticsDialog(totalPasswords, totalBytes int64, byLength map[int]int64, byType map[string]int64) {
	statsText := fmt.Sprintf("Generation Statistics\n")
	statsText += fmt.Sprintf("===================\n\n")
	statsText += fmt.Sprintf("Total Passwords: %d\n", totalPasswords)
	statsText += fmt.Sprintf("Total Size: %d bytes (%.2f MB)\n\n", totalBytes, float64(totalBytes)/(1024*1024))
	
	statsText += fmt.Sprintf("By Length:\n")
	for length := 1; length <= 128; length++ {
		if count, ok := byLength[length]; ok && count > 0 {
			statsText += fmt.Sprintf("  %d chars: %d passwords\n", length, count)
		}
	}
	
	statsText += fmt.Sprintf("\nBy Type:\n")
	for ptype, count := range byType {
		if count > 0 {
			statsText += fmt.Sprintf("  %s: %d passwords\n", ptype, count)
		}
	}
	
	statsEntry := widget.NewMultiLineEntry()
	statsEntry.SetText(statsText)
	statsEntry.Disable()
	statsEntry.Wrapping = fyne.TextWrapOff
	
	scroll := container.NewScroll(statsEntry)
	scroll.SetMinSize(fyne.NewSize(500, 400))
	
	// Add export buttons
	exportJSONBtn := widget.NewButton("Export JSON", func() {
		if a.generator == nil {
			dialog.ShowError(fmt.Errorf("no statistics available"), a.window)
			return
		}
		
		stats := a.generator.GetStatistics()
		if stats == nil {
			dialog.ShowError(fmt.Errorf("no statistics available"), a.window)
			return
		}
		
		// Show file save dialog
		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil || writer == nil {
				return
			}
			defer writer.Close()
			
			filePath := writer.URI().Path()
			if err := core.ExportStatisticsJSON(stats, filePath); err != nil {
				dialog.ShowError(fmt.Errorf("failed to export: %v", err), a.window)
			} else {
				dialog.ShowInformation("Success", fmt.Sprintf("Statistics exported to %s", filePath), a.window)
			}
		}, a.window)
	})
	
	exportHTMLBtn := widget.NewButton("Export HTML", func() {
		if a.generator == nil {
			dialog.ShowError(fmt.Errorf("no statistics available"), a.window)
			return
		}
		
		stats := a.generator.GetStatistics()
		if stats == nil {
			dialog.ShowError(fmt.Errorf("no statistics available"), a.window)
			return
		}
		
		// Show file save dialog
		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil || writer == nil {
				return
			}
			defer writer.Close()
			
			filePath := writer.URI().Path()
			if err := core.ExportStatisticsHTML(stats, filePath); err != nil {
				dialog.ShowError(fmt.Errorf("failed to export: %v", err), a.window)
			} else {
				dialog.ShowInformation("Success", fmt.Sprintf("Statistics exported to %s", filePath), a.window)
			}
		}, a.window)
	})
	
	exportBox := container.NewHBox(exportJSONBtn, exportHTMLBtn)
	
	content := container.NewVBox(
		scroll,
		container.NewPadded(exportBox),
	)
	
	dialog.ShowCustom("Generation Statistics", "Close", content, a.window)
}

