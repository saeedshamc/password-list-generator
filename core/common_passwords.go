package core

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CommonPasswordsManager manages common password lists
type CommonPasswordsManager struct {
	lists map[string][]string
}

// NewCommonPasswordsManager creates a new common passwords manager
func NewCommonPasswordsManager() *CommonPasswordsManager {
	return &CommonPasswordsManager{
		lists: make(map[string][]string),
	}
}

// LoadList loads a password list from file
func (c *CommonPasswordsManager) LoadList(name, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var passwords []string
	scanner := bufio.NewScanner(file)
	
	// Limit to first 100k passwords to avoid memory issues
	maxPasswords := 100000
	count := 0
	
	for scanner.Scan() && count < maxPasswords {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			passwords = append(passwords, line)
			count++
		}
	}
	
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	
	c.lists[name] = passwords
	return nil
}

// GetList returns a password list by name
func (c *CommonPasswordsManager) GetList(name string) ([]string, bool) {
	list, ok := c.lists[name]
	return list, ok
}

// GenerateFromList generates passwords from a loaded list
func (c *CommonPasswordsManager) GenerateFromList(ctx context.Context, listName string) <-chan string {
	output := make(chan string, 100)
	
	go func() {
		defer close(output)
		
		list, ok := c.GetList(listName)
		if !ok {
			return
		}
		
		for _, password := range list {
			select {
			case <-ctx.Done():
				return
			case output <- password:
			}
		}
	}()
	
	return output
}

// GetBuiltinLists returns list of built-in common password lists
func GetBuiltinLists() []string {
	return []string{
		"rockyou",
		"seclists",
		"10-million-password-list",
		"common-passwords",
	}
}

// LoadBuiltinList attempts to load a built-in list from common locations
func (c *CommonPasswordsManager) LoadBuiltinList(name string) error {
	// Common locations for password lists
	commonLocations := []string{
		"/usr/share/wordlists",
		"/usr/share/seclists",
		"/opt/wordlists",
		"~/.wordlists",
		"./wordlists",
		"./data",
	}
	
	// Map list names to common filenames
	filenameMap := map[string][]string{
		"rockyou": {
			"rockyou.txt",
			"rockyou-50.txt",
			"rockyou-75.txt",
		},
		"seclists": {
			"Passwords/Common-Credentials/10-million-password-list-top-1000000.txt",
			"Passwords/Common-Credentials/best1050.txt",
			"Passwords/Common-Credentials/top-passwords-shortlist.txt",
		},
		"10-million-password-list": {
			"10-million-password-list-top-1000000.txt",
			"10-million-password-list-top-100000.txt",
		},
		"common-passwords": {
			"common-passwords.txt",
			"common.txt",
			"passwords.txt",
		},
	}
	
	filenames, ok := filenameMap[name]
	if !ok {
		return fmt.Errorf("unknown built-in list: %s", name)
	}
	
	// Try to find the file
	for _, location := range commonLocations {
		// Expand ~
		if strings.HasPrefix(location, "~") {
			home, err := os.UserHomeDir()
			if err != nil {
				continue
			}
			location = filepath.Join(home, location[2:])
		}
		
		for _, filename := range filenames {
			fullPath := filepath.Join(location, filename)
			if _, err := os.Stat(fullPath); err == nil {
				// File exists, try to load it
				if err := c.LoadList(name, fullPath); err == nil {
					return nil
				}
			}
		}
	}
	
	return fmt.Errorf("could not find built-in list: %s", name)
}

// GenerateCommonPasswords generates common passwords with variations
func (c *CommonPasswordsManager) GenerateCommonPasswords(ctx context.Context, listNames []string, applyMutations bool) <-chan string {
	output := make(chan string, 100)
	
	go func() {
		defer close(output)
		
		seen := make(map[string]bool)
		
		for _, listName := range listNames {
			list, ok := c.GetList(listName)
			if !ok {
				continue
			}
			
			for _, password := range list {
				select {
				case <-ctx.Done():
					return
				default:
					// Deduplicate
					if seen[password] {
						continue
					}
					seen[password] = true
					
					output <- password
					
					// Apply common mutations if enabled
					if applyMutations {
						// Case variations
						output <- strings.ToUpper(password)
						output <- strings.ToLower(password)
						if len(password) > 0 {
							output <- strings.ToUpper(string(password[0])) + strings.ToLower(password[1:])
						}
						
						// Add common suffixes
						commonSuffixes := []string{"123", "1234", "!", "!@#", "2024", "2023"}
						for _, suffix := range commonSuffixes {
							output <- password + suffix
						}
						
						// Add common prefixes
						commonPrefixes := []string{"!", "@", "#", "$"}
						for _, prefix := range commonPrefixes {
							output <- prefix + password
						}
					}
				}
			}
		}
	}()
	
	return output
}

