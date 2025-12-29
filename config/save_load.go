package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SaveConfig saves the configuration to a JSON file
func SaveConfig(cfg *Config, filePath string) error {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(filePath string) (*Config, error) {
	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal JSON
	cfg := &Config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate and set defaults for missing fields
	cfg = validateAndSetDefaults(cfg)

	return cfg, nil
}

// validateAndSetDefaults ensures all fields have valid values
func validateAndSetDefaults(cfg *Config) *Config {
	// If config is nil, return default
	if cfg == nil {
		return DefaultConfig()
	}

	// Validate TargetPreset
	if cfg.TargetPreset.Preset == "" {
		cfg.TargetPreset.Preset = "none"
	}
	if cfg.TargetPreset.MinLength <= 0 {
		cfg.TargetPreset.MinLength = 4
	}
	if cfg.TargetPreset.MaxLength <= 0 {
		cfg.TargetPreset.MaxLength = 128
	}

	// Validate Combination
	if cfg.Combination.MaxDepth <= 0 {
		cfg.Combination.MaxDepth = 3
	}
	if cfg.Combination.MinLength <= 0 {
		cfg.Combination.MinLength = 4
	}
	if cfg.Combination.MaxLength <= 0 {
		cfg.Combination.MaxLength = 128
	}

	// Validate Mutation
	if cfg.Mutation.Intensity == "" {
		cfg.Mutation.Intensity = "medium"
	}
	if cfg.Mutation.MaxMutations <= 0 {
		cfg.Mutation.MaxMutations = 3
	}

	// Validate Output
	if cfg.Output.Format == "" {
		cfg.Output.Format = "plain"
	}
	if cfg.Output.FilePath == "" {
		cfg.Output.FilePath = "wordlist.txt"
	}

	// Ensure slices are not nil
	if cfg.Input.Names == nil {
		cfg.Input.Names = []string{}
	}
	if cfg.Input.Keywords == nil {
		cfg.Input.Keywords = []string{}
	}
	if cfg.Input.ExternalFiles == nil {
		cfg.Input.ExternalFiles = []string{}
	}
	if cfg.Input.FreeText == nil {
		cfg.Input.FreeText = []string{}
	}
	if cfg.Input.SelectedCountries == nil {
		cfg.Input.SelectedCountries = []string{}
	}
	if cfg.Advanced.Prefixes == nil {
		cfg.Advanced.Prefixes = []string{}
	}
	if cfg.Advanced.Suffixes == nil {
		cfg.Advanced.Suffixes = []string{}
	}

	return cfg
}

// GetConfigProfilesDir returns the directory for saved profiles
func GetConfigProfilesDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	profilesDir := filepath.Join(homeDir, ".passgen", "profiles")
	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create profiles directory: %w", err)
	}

	return profilesDir, nil
}

// ListSavedProfiles returns a list of saved profile names
func ListSavedProfiles() ([]string, error) {
	profilesDir, err := GetConfigProfilesDir()
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(profilesDir)
	if err != nil {
		return []string{}, nil // Return empty if directory doesn't exist
	}

	profiles := []string{}
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			name := file.Name()[:len(file.Name())-5] // Remove .json extension
			profiles = append(profiles, name)
		}
	}

	return profiles, nil
}

// SaveProfile saves a configuration as a named profile
func SaveProfile(cfg *Config, profileName string) error {
	profilesDir, err := GetConfigProfilesDir()
	if err != nil {
		return err
	}

	filePath := filepath.Join(profilesDir, profileName+".json")
	return SaveConfig(cfg, filePath)
}

// LoadProfile loads a configuration from a named profile
func LoadProfile(profileName string) (*Config, error) {
	profilesDir, err := GetConfigProfilesDir()
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(profilesDir, profileName+".json")
	return LoadConfig(filePath)
}

// DeleteProfile deletes a saved profile
func DeleteProfile(profileName string) error {
	profilesDir, err := GetConfigProfilesDir()
	if err != nil {
		return err
	}

	filePath := filepath.Join(profilesDir, profileName+".json")
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}

	return nil
}

