package core

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/passgen/config"
)

// ProgressCallback is called with progress updates
type ProgressCallback func(lines int64, bytes int64)

// Generator orchestrates the password generation pipeline
type Generator struct {
	cfg            *config.Config
	mutationEngine *MutationEngine
	deduplicator   *Deduplicator
	randomizer     *RandomizationEngine
	combiner       *CombinationEngine
	validator      *PasswordValidator
	smartFilter    *SmartFilter
	customRules    *CustomRulesEngine
	hashcatRules   *HashcatRulesEngine
	johnRules     *JohnRulesEngine
	keyboardPatterns *KeyboardPatternsEngine
	commonPasswords *CommonPasswordsManager
	customCharset  *CustomCharsetEngine
	maskAttack     *MaskAttackEngine
	perfOptimizer  *PerformanceOptimizer
	statistics     *Statistics
	checkpoint     *Checkpoint
	checkpointMgr  *CheckpointManager
	parallelProc   *ParallelProcessor
	jobID          string
	paused         bool
	mu             sync.Mutex
}

func NewGenerator(cfg *config.Config) *Generator {
	// Get target constraints
	minLen, maxLen, asciiOnly := GetTargetConstraints(cfg.TargetPreset.Preset)
	if cfg.TargetPreset.MinLength > 0 {
		minLen = cfg.TargetPreset.MinLength
	}
	if cfg.TargetPreset.MaxLength > 0 {
		maxLen = cfg.TargetPreset.MaxLength
	}
	if cfg.TargetPreset.ASCIIOnly {
		asciiOnly = true
	}

	validator := NewPasswordValidator(cfg.TargetPreset.Preset, minLen, maxLen, asciiOnly)

	var customRules *CustomRulesEngine
	if len(cfg.Filter.CustomRules) > 0 {
		customRules = NewCustomRulesEngine(cfg.Filter.CustomRules)
	}

	var hashcatRules *HashcatRulesEngine
	if cfg.Mutation.HashcatRules && cfg.Mutation.HashcatRulesFile != "" {
		rules, err := LoadRulesFromFile(cfg.Mutation.HashcatRulesFile)
		if err == nil {
			hashcatRules = NewHashcatRulesEngine(rules)
		}
	}

	var johnRules *JohnRulesEngine
	if cfg.Mutation.JohnRules && cfg.Mutation.JohnRulesFile != "" {
		rules, err := LoadJohnRulesFromFile(cfg.Mutation.JohnRulesFile)
		if err == nil {
			johnRules = NewJohnRulesEngine(rules)
		}
	}

	var keyboardPatterns *KeyboardPatternsEngine
	if cfg.Input.KeyboardPatterns {
		keyboardPatterns = NewKeyboardPatternsEngine()
	}

	commonPasswords := NewCommonPasswordsManager()
	// Load built-in lists if specified
	for _, listName := range cfg.Input.CommonPasswordLists {
		commonPasswords.LoadBuiltinList(listName)
	}
	// Load custom files
	for _, filePath := range cfg.Input.CommonPasswordFiles {
		listName := filepath.Base(filePath)
		commonPasswords.LoadList(listName, filePath)
	}

	customCharset := NewCustomCharsetEngine()
	// Add custom charsets from config
	if cfg.CustomCharset.Enabled {
		for _, charsetDef := range cfg.CustomCharset.CustomCharsets {
			expanded := ExpandCharset(charsetDef.Characters)
			charset := &CustomCharset{
				Name:        charsetDef.Name,
				Description: charsetDef.Description,
				Characters:  expanded,
				Enabled:     charsetDef.Enabled,
			}
			customCharset.AddCharset(charset)
		}
	}

	var maskAttack *MaskAttackEngine
	if cfg.Input.MaskAttack && len(cfg.Input.MaskPatterns) > 0 {
		maskAttack = NewMaskAttackEngine(cfg.Input.MaskPatterns)
	}

	checkpointMgr, _ := NewCheckpointManager()
	parallelProc := NewParallelProcessor()
	perfOptimizer := NewPerformanceOptimizer()

	return &Generator{
		cfg:            cfg,
		mutationEngine: NewMutationEngine(&cfg.Mutation),
		deduplicator:   NewDeduplicator(&cfg.Deduplication),
		randomizer:     NewRandomizationEngine(&cfg.Randomization),
		combiner:       NewCombinationEngine(&cfg.Combination),
		validator:      validator,
		smartFilter:    NewSmartFilter(&cfg.Filter),
		customRules:    customRules,
		hashcatRules:   hashcatRules,
		johnRules:     johnRules,
		keyboardPatterns: keyboardPatterns,
		commonPasswords: commonPasswords,
		customCharset:  customCharset,
		maskAttack:     maskAttack,
		perfOptimizer:  perfOptimizer,
		statistics:     NewStatistics(),
		checkpointMgr:  checkpointMgr,
		checkpoint:     &Checkpoint{Config: cfg},
		parallelProc:   parallelProc,
	}
}

// Generate runs the complete generation pipeline
func (g *Generator) Generate(ctx context.Context, progress ProgressCallback) error {
	// Open output file (append if resuming)
	var writer io.Writer
	var file *os.File
	var compressWriter *CompressionWriter
	var err error
	var linesWritten int64
	var bytesWritten int64
	
	// Check if compression is enabled
	if g.cfg.Output.Compress {
		compressWriter, err = NewCompressionWriter(g.cfg.Output.FilePath)
		if err != nil {
			return fmt.Errorf("failed to create compression writer: %w", err)
		}
		defer compressWriter.Close()
		writer = compressWriter
	} else {
		// Check if resuming from checkpoint
		if g.checkpoint != nil && g.checkpoint.LinesWritten > 0 {
			file, err = os.OpenFile(g.cfg.Output.FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			linesWritten = g.checkpoint.LinesWritten
			bytesWritten = g.checkpoint.BytesWritten
		} else {
			file, err = os.Create(g.cfg.Output.FilePath)
		}
		
		if err != nil {
			return fmt.Errorf("failed to create/open output file: %w", err)
		}
		defer file.Close()
		writer = file
	}
	
	// Use file splitter if enabled, otherwise use buffered writer
	var finalWriter io.Writer
	var fileSplitter *FileSplitter
	
	if g.cfg.Output.SplitFiles && g.cfg.Output.MaxFileSize > 0 {
		fileSplitter = NewFileSplitter(g.cfg.Output.FilePath, g.cfg.Output.MaxFileSize)
		finalWriter = fileSplitter
		defer fileSplitter.Close()
	} else {
		bufferedWriter := NewBufferedWriter(writer)
		finalWriter = bufferedWriter
		defer bufferedWriter.Flush()
	}
	
	// Check if numeric-only mode
	if g.cfg.Numeric.Enabled {
		return g.generateNumeric(ctx, finalWriter, &linesWritten, &bytesWritten, progress)
	}
	
	// Normal pipeline: input → normalize → combine → mutate → randomize → deduplicate → output
	inputCh := CollectInputs(ctx, g.cfg)
	
	// Add keyboard patterns if enabled
	if g.keyboardPatterns != nil && g.cfg.Input.KeyboardPatterns {
		keyboardCh := g.keyboardPatterns.GenerateCommonPatterns(ctx)
		// Merge keyboard patterns with input
		mergedCh := make(chan string, 100)
		go func() {
			defer close(mergedCh)
			done := make(chan bool, 2)
			
			// Forward input
			go func() {
				defer func() { done <- true }()
				for word := range inputCh {
					select {
					case <-ctx.Done():
						return
					case mergedCh <- word:
					}
				}
			}()
			
			// Forward keyboard patterns
			go func() {
				defer func() { done <- true }()
				for pattern := range keyboardCh {
					select {
					case <-ctx.Done():
						return
					case mergedCh <- pattern:
					}
				}
			}()
			
			// Wait for both
			<-done
			<-done
		}()
		inputCh = mergedCh
	}
	
	// Apply mutations
	mutatedCh := make(chan string, 100)
	go func() {
		defer close(mutatedCh)
		for word := range inputCh {
			select {
			case <-ctx.Done():
				return
			default:
				// Apply case variations if enabled
				if g.cfg.Advanced.CaseVariations {
					variations := ApplyCaseVariations(word)
					for _, variant := range variations {
						mutatedCh <- variant
					}
				} else {
					mutatedCh <- word
				}
				
				// Apply mutations
				if g.cfg.Mutation.Enabled {
					mutations := g.mutationEngine.Mutate(word)
					for _, mut := range mutations {
						if mut != word {
							select {
							case <-ctx.Done():
								return
							case mutatedCh <- mut:
							}
						}
					}
				}
			}
		}
	}()
	
	// Apply combinations
	combinedCh := g.combiner.Combine(ctx, mutatedCh)
	
	// Apply randomization
	randomizedCh := g.randomizer.Randomize(ctx, combinedCh)
	
	// Apply advanced features
	advancedCh := g.applyAdvanced(ctx, randomizedCh)
	
	// Apply Hashcat rules (if enabled)
	var hashcatCh <-chan string
	if g.hashcatRules != nil && g.cfg.Mutation.HashcatRules {
		hashcatCh = g.hashcatRules.ProcessPasswords(ctx, advancedCh)
	} else {
		hashcatCh = advancedCh
	}
	
	// Apply validation (if enabled)
	var validatedCh <-chan string
	if g.cfg.TargetPreset.EnforceRules {
		validatedCh = FilterPasswords(g.validator, hashcatCh)
	} else {
		validatedCh = hashcatCh
	}
	
	// Apply smart filtering (if enabled)
	var filteredCh <-chan string
	if g.cfg.Filter.Enabled {
		filteredCh = g.smartFilter.FilterPasswords(ctx, validatedCh)
	} else {
		filteredCh = validatedCh
	}
	
	// Apply custom rules (if enabled)
	var rulesCh <-chan string
	if g.customRules != nil && len(g.cfg.Filter.CustomRules) > 0 {
		// Check if any custom rules are enabled
		hasEnabledRules := false
		for _, rule := range g.cfg.Filter.CustomRules {
			if rule.Enabled {
				hasEnabledRules = true
				break
			}
		}
		if hasEnabledRules {
			rulesCh = g.customRules.FilterPasswords(ctx, filteredCh)
		} else {
			rulesCh = filteredCh
		}
	} else {
		rulesCh = filteredCh
	}
	
	// Apply custom charset filter (if enabled)
	var charsetCh <-chan string
	if g.customCharset != nil && g.cfg.CustomCharset.Enabled && g.cfg.CustomCharset.CharsetName != "" {
		charsetCh = g.customCharset.FilterByCharset(ctx, rulesCh, g.cfg.CustomCharset.CharsetName)
	} else {
		charsetCh = rulesCh
	}
	
	// Deduplicate
	finalCh := g.deduplicator.Deduplicate(ctx, charsetCh)
	
	// Write to file
	for word := range finalCh {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Check limits
			if g.cfg.Limits.MaxOutputLines > 0 && linesWritten >= g.cfg.Limits.MaxOutputLines {
				return nil
			}
			if g.cfg.Limits.MaxFileSize > 0 && bytesWritten >= g.cfg.Limits.MaxFileSize {
				return nil
			}
			
			line := g.formatOutput(word) + "\n"
			var n int
			var err error
			
			// Write to appropriate writer
			if fileSplitter != nil {
				n, err = fileSplitter.WriteString(line)
			} else if bw, ok := finalWriter.(*BufferedWriter); ok {
				n, err = bw.WriteString(line)
			} else {
				n, err = fmt.Fprintf(finalWriter, "%s", line)
			}
			
			if err != nil {
				return fmt.Errorf("failed to write: %w", err)
			}
			
			atomic.AddInt64(&linesWritten, 1)
			atomic.AddInt64(&bytesWritten, int64(n))
			
			// Update checkpoint periodically (every 1000 lines)
			if g.checkpointMgr != nil && g.jobID != "" && linesWritten%1000 == 0 {
				g.mu.Lock()
				g.checkpoint.LinesWritten = linesWritten
				g.checkpoint.BytesWritten = bytesWritten
				g.checkpoint.LastPassword = word
				g.checkpointMgr.SaveCheckpoint(g.jobID, g.checkpoint)
				g.mu.Unlock()
			}
			
			// Check for pause
			g.mu.Lock()
			paused := g.paused
			g.mu.Unlock()
			if paused {
				// Save checkpoint before pausing
				if g.checkpointMgr != nil && g.jobID != "" {
					g.mu.Lock()
					g.checkpoint.LinesWritten = linesWritten
					g.checkpoint.BytesWritten = bytesWritten
					g.checkpoint.LastPassword = word
					g.checkpointMgr.SaveCheckpoint(g.jobID, g.checkpoint)
					g.mu.Unlock()
				}
				return fmt.Errorf("generation paused")
			}
			
			// Update statistics
			if g.statistics != nil {
				g.statistics.AddPassword(word)
			}
			
			if progress != nil {
				progress(linesWritten, bytesWritten)
			}
		}
	}
	
	// Save final checkpoint
	if g.checkpointMgr != nil && g.jobID != "" {
		g.mu.Lock()
		g.checkpoint.LinesWritten = linesWritten
		g.checkpoint.BytesWritten = bytesWritten
		g.checkpointMgr.SaveCheckpoint(g.jobID, g.checkpoint)
		g.mu.Unlock()
	}
	
	return nil
}

func (g *Generator) generateNumeric(ctx context.Context, writer io.Writer, linesWritten *int64, bytesWritten *int64, progress ProgressCallback) error {
	gen := NewNumericGenerator(&g.cfg.Numeric)
	
	for word := range gen.Generate(ctx) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Check limits
			if g.cfg.Limits.MaxOutputLines > 0 && *linesWritten >= g.cfg.Limits.MaxOutputLines {
				return nil
			}
			if g.cfg.Limits.MaxFileSize > 0 && *bytesWritten >= g.cfg.Limits.MaxFileSize {
				return nil
			}
			
			line := g.formatOutput(word) + "\n"
			var n int
			var err error
			
			// Write to appropriate writer
			if fs, ok := writer.(*FileSplitter); ok {
				n, err = fs.WriteString(line)
			} else if bw, ok := writer.(*BufferedWriter); ok {
				n, err = bw.WriteString(line)
			} else {
				n, err = fmt.Fprintf(writer, "%s", line)
			}
			
			if err != nil {
				return fmt.Errorf("failed to write: %w", err)
			}
			
			atomic.AddInt64(linesWritten, 1)
			atomic.AddInt64(bytesWritten, int64(n))
			
			// Update checkpoint periodically
			if g.checkpointMgr != nil && g.jobID != "" && *linesWritten%1000 == 0 {
				g.mu.Lock()
				g.checkpoint.LinesWritten = *linesWritten
				g.checkpoint.BytesWritten = *bytesWritten
				g.checkpoint.LastPassword = word
				g.checkpointMgr.SaveCheckpoint(g.jobID, g.checkpoint)
				g.mu.Unlock()
			}
			
			// Check for pause
			g.mu.Lock()
			paused := g.paused
			g.mu.Unlock()
			if paused {
				if g.checkpointMgr != nil && g.jobID != "" {
					g.mu.Lock()
					g.checkpoint.LinesWritten = *linesWritten
					g.checkpoint.BytesWritten = *bytesWritten
					g.checkpoint.LastPassword = word
					g.checkpointMgr.SaveCheckpoint(g.jobID, g.checkpoint)
					g.mu.Unlock()
				}
				return fmt.Errorf("generation paused")
			}
			
			// Update statistics
			if g.statistics != nil {
				g.statistics.AddPassword(word)
			}
			
			if progress != nil {
				progress(*linesWritten, *bytesWritten)
			}
		}
	}
	
	// Save final checkpoint
	if g.checkpointMgr != nil && g.jobID != "" {
		g.mu.Lock()
		g.checkpoint.LinesWritten = *linesWritten
		g.checkpoint.BytesWritten = *bytesWritten
		g.checkpointMgr.SaveCheckpoint(g.jobID, g.checkpoint)
		g.mu.Unlock()
	}
	
	return nil
}

func (g *Generator) applyAdvanced(ctx context.Context, input <-chan string) <-chan string {
	output := make(chan string, 100)
	
	go func() {
		defer close(output)
		
		for word := range input {
			select {
			case <-ctx.Done():
				return
			default:
				output <- word
				
				// Apply prefixes
				for _, prefix := range g.cfg.Advanced.Prefixes {
					output <- prefix + word
				}
				
				// Apply suffixes
				for _, suffix := range g.cfg.Advanced.Suffixes {
					output <- word + suffix
				}
				
				// Word duplication
				if g.cfg.Advanced.WordDuplication {
					output <- word + word
				}
				
				// Repeated characters
				if g.cfg.Advanced.RepeatedChars && len(word) > 0 {
					lastChar := string(word[len(word)-1])
					output <- word + lastChar + lastChar
				}
			}
		}
	}()
	
	return output
}

// EstimateOutputSize estimates the output size (rough estimate)
func (g *Generator) EstimateOutputSize() (lines int64, bytes int64) {
	// This is a rough estimate
	// In practice, this would require more sophisticated calculation
	inputCount := int64(len(g.cfg.Input.Names) + len(g.cfg.Input.Keywords) + len(g.cfg.Input.FreeText))
	
	// Estimate combinations
	estimate := inputCount * 100 // Rough multiplier
	
	if g.cfg.Mutation.Enabled {
		estimate *= int64(g.cfg.Mutation.MaxMutations + 1)
	}
	
	if g.cfg.Combination.Enabled {
		estimate *= 10 // Rough estimate for combinations
	}
	
	// Average password length estimate
	avgLength := int64(12)
	bytes = estimate * avgLength
	
	return estimate, bytes
}

// formatOutput formats a password based on the selected output format
func (g *Generator) formatOutput(password string) string {
	switch g.cfg.Output.Format {
	case "hashcat":
		// Hashcat format: plain password, one per line
		return password
	case "hydra":
		// Hydra format: username:password
		username := g.cfg.Output.Username
		if username == "" {
			username = "admin" // Default username
		}
		return username + ":" + password
	case "john":
		// John format: hash:password or username:password
		if g.cfg.Output.Hash != "" {
			return g.cfg.Output.Hash + ":" + password
		}
		username := g.cfg.Output.Username
		if username == "" {
			username = "admin" // Default username
		}
		return username + ":" + password
	case "aircrack":
		// Aircrack format: plain password, one per line (for WPA/WPA2)
		return password
	default:
		// Plain format: just the password
		return password
	}
}

// GetStatistics returns current generation statistics
func (g *Generator) GetStatistics() *Statistics {
	return g.statistics
}

// SetJobID sets the job ID for checkpoint management
func (g *Generator) SetJobID(jobID string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.jobID = jobID
	if g.checkpoint != nil {
		g.checkpoint.Config = g.cfg
	}
}

// Pause pauses the generation
func (g *Generator) Pause() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.paused = true
}

// Resume resumes the generation
func (g *Generator) Resume() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.paused = false
}

// IsPaused returns whether generation is paused
func (g *Generator) IsPaused() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.paused
}

// LoadCheckpoint loads a checkpoint for resuming
func (g *Generator) LoadCheckpoint(jobID string) error {
	if g.checkpointMgr == nil {
		return fmt.Errorf("checkpoint manager not initialized")
	}
	
	checkpoint, err := g.checkpointMgr.LoadCheckpoint(jobID)
	if err != nil {
		return err
	}
	
	g.mu.Lock()
	defer g.mu.Unlock()
	g.checkpoint = checkpoint
	g.jobID = jobID
	g.cfg = checkpoint.Config
	
	return nil
}

