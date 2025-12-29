package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/passgen/config"
	"github.com/passgen/core"
	"github.com/passgen/gui"
)

func main() {
	// CLI flags
	cliMode := flag.Bool("cli", false, "Run in CLI mode")
	configFile := flag.String("config", "", "Configuration file path")
	flag.Parse()

	if *cliMode {
		runCLI(*configFile)
	} else {
		runGUI()
	}
}

func runGUI() {
	app := gui.NewApp()
	app.Run()
}

func runCLI(configPath string) {
	cfg := config.DefaultConfig()
	
	// TODO: Load config from file if provided
	_ = configPath
	
	// Example CLI usage
	fmt.Println("Password Wordlist Generator - CLI Mode")
	fmt.Println("=====================================")
	
	// Create generator
	generator := core.NewGenerator(cfg)
	
	// Estimate output
	lines, bytes := generator.EstimateOutputSize()
	fmt.Printf("Estimated output: %d lines, %d bytes\n", lines, bytes)
	
	// Generate
	fmt.Println("Generating wordlist...")
	err := generator.Generate(nil, func(lines int64, bytes int64) {
		if lines%10000 == 0 {
			fmt.Printf("Progress: %d lines, %d bytes\n", lines, bytes)
		}
	})
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("Generation complete!")
}

