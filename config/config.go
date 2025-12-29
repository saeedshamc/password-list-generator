package config

import (
	"time"
)

// Config holds all configuration for the wordlist generator
type Config struct {
	TargetPreset TargetPresetConfig `json:"target_preset"`
	Input        InputConfig        `json:"input"`
	Numeric      NumericConfig      `json:"numeric"`
	Combination  CombinationConfig  `json:"combination"`
	Mutation     MutationConfig     `json:"mutation"`
	Randomization RandomizationConfig `json:"randomization"`
	Advanced     AdvancedConfig     `json:"advanced"`
	Deduplication DeduplicationConfig `json:"deduplication"`
	Filter       FilterConfig       `json:"filter"`
	Limits       LimitsConfig       `json:"limits"`
	Output       OutputConfig       `json:"output"`
	CustomCharset CustomCharsetConfig `json:"custom_charset"`
}

// CustomCharsetConfig defines custom character set options
type CustomCharsetConfig struct {
	Enabled     bool     `json:"enabled"`
	CharsetName string   `json:"charset_name"` // Built-in or custom
	CustomCharsets []CustomCharsetDef `json:"custom_charsets"`
	MinLength   int      `json:"min_length"`
	MaxLength   int      `json:"max_length"`
}

// CustomCharsetDef defines a custom character set
type CustomCharsetDef struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Characters  string `json:"characters"` // Can use patterns like "a-z0-9!@#"
	Enabled     bool   `json:"enabled"`
}

// TargetPresetConfig defines target-specific constraints
type TargetPresetConfig struct {
	Preset      string `json:"preset"`       // "none", "wifi", "ssh", "web", "mobile", "iot"
	MinLength   int    `json:"min_length"`    // Override min length
	MaxLength   int    `json:"max_length"`   // Override max length
	ASCIIOnly   bool   `json:"ascii_only"`    // Force ASCII encoding
	EnforceRules bool  `json:"enforce_rules"` // Enable strict validation
}

// InputConfig defines input sources
type InputConfig struct {
	Names          []string `json:"names"`
	Keywords       []string `json:"keywords"`
	NumericRanges  []NumericRange `json:"numeric_ranges"`
	Dates          []DateRange `json:"dates"`
	ExternalFiles  []string `json:"external_files"`
	FreeText       []string `json:"free_text"`
	SelectedCountries []string `json:"selected_countries"` // Country codes like "iran", "united_states"
	UseAllCountryNames bool   `json:"use_all_country_names"` // Use all names or custom subset
	IncludeBirthdays  bool   `json:"include_birthdays"`     // Include common birthdays
	IncludeHolidays   bool   `json:"include_holidays"`      // Include holidays
	IncludeSpecialDates bool `json:"include_special_dates"`  // Include special dates
	KeyboardPatterns  bool   `json:"keyboard_patterns"`     // Include keyboard patterns
	CommonPasswordLists []string `json:"common_password_lists"` // List names like "rockyou", "seclists"
	CommonPasswordFiles []string `json:"common_password_files"` // Custom file paths
	ApplyMutationsToCommon bool `json:"apply_mutations_to_common"` // Apply mutations to common passwords
	MaskAttack        bool     `json:"mask_attack"`           // Enable mask attack generation
	MaskPatterns      []string `json:"mask_patterns"`         // Hashcat mask patterns like "?l?l?l?d?d?d"
}

// NumericRange defines a range of numbers
type NumericRange struct {
	Start      int64  `json:"start"`
	End        int64  `json:"end"`
	FixedLength int    `json:"fixed_length"` // 0 = variable length
	LeadingZeros bool `json:"leading_zeros"`
}

// DateRange defines a range of dates
type DateRange struct {
	Start time.Time `json:"start"`
	End   time.Time   `json:"end"`
	Format string     `json:"format"` // "YYYY", "YYYYMMDD", "DDMMYYYY", etc.
}

// NumericConfig for numeric-only mode
type NumericConfig struct {
	Enabled      bool  `json:"enabled"`
	Start        int64 `json:"start"`
	End          int64 `json:"end"`
	FixedLength  int   `json:"fixed_length"`
	LeadingZeros bool  `json:"leading_zeros"`
	Randomized   bool  `json:"randomized"`
	Shuffle      bool  `json:"shuffle"`
	Hexadecimal  bool  `json:"hexadecimal"`
	Deterministic bool `json:"deterministic"`
	Seed         int64 `json:"seed"`
}

// CombinationConfig defines combination rules
type CombinationConfig struct {
	Enabled         bool     `json:"enabled"`
	MaxDepth        int      `json:"max_depth"`
	MinLength       int      `json:"min_length"`
	MaxLength       int      `json:"max_length"`
	Patterns        []Pattern `json:"patterns"`
	WordWord        bool     `json:"word_word"`
	WordNumber      bool     `json:"word_number"`
	WordSymbolNumber bool    `json:"word_symbol_number"`
	MultiWord       bool     `json:"multi_word"`
	UseTemplates    bool     `json:"use_templates"`
	SelectedTemplates []string `json:"selected_templates"` // Template names
	CustomTemplates  []string  `json:"custom_templates"`    // Custom template strings
}

// Pattern defines a combination pattern
type Pattern struct {
	Template string `json:"template"` // e.g., "{word}{number}", "{word}_{number}"
	Enabled  bool   `json:"enabled"`
}

// MutationConfig defines mutation rules
type MutationConfig struct {
	Enabled          bool   `json:"enabled"`
	Leetspeak        bool   `json:"leetspeak"`
	Unicode          bool   `json:"unicode"`
	Intensity        string `json:"intensity"` // "low", "medium", "aggressive"
	MaxMutations     int    `json:"max_mutations"`
	HashcatRules     bool   `json:"hashcat_rules"` // Enable Hashcat rules
	HashcatRulesFile string `json:"hashcat_rules_file"` // Path to .rule file
	JohnRules        bool   `json:"john_rules"` // Enable John the Ripper rules
	JohnRulesFile    string `json:"john_rules_file"` // Path to John rules file
}

// RandomizationConfig defines randomization options
type RandomizationConfig struct {
	Enabled            bool   `json:"enabled"`
	ShuffleChars       bool   `json:"shuffle_chars"`
	SwapWordOrder      bool   `json:"swap_word_order"`
	InsertSymbols      bool   `json:"insert_symbols"`
	Symbols            string `json:"symbols"` // e.g., "!@#$%"
	Deterministic      bool   `json:"deterministic"`
	Seed               int64  `json:"seed"`
}

// AdvancedConfig defines advanced options
type AdvancedConfig struct {
	CaseVariations    bool     `json:"case_variations"` // lower, upper, mixed
	Prefixes          []string `json:"prefixes"`
	Suffixes          []string `json:"suffixes"`
	RepeatedChars     bool     `json:"repeated_chars"`
	WordDuplication   bool     `json:"word_duplication"`
	RandomSymbolInsert bool    `json:"random_symbol_insert"`
}

// DeduplicationConfig defines deduplication options
type DeduplicationConfig struct {
	Enabled bool `json:"enabled"`
}

// CustomRule defines a custom validation rule
type CustomRule struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Pattern     string `json:"pattern"`     // Regex pattern (optional)
	MinLength   int    `json:"min_length"`  // Minimum length (0 = no limit)
	MaxLength   int    `json:"max_length"`  // Maximum length (0 = no limit)
	Require     string `json:"require"`     // Required characters (e.g., "a-z,0-9,!@#")
	Exclude     string `json:"exclude"`     // Excluded characters or patterns
	Enabled     bool   `json:"enabled"`
}

// FilterConfig defines smart filtering options
type FilterConfig struct {
	Enabled          bool         `json:"enabled"`
	MinComplexity    string       `json:"min_complexity"`    // "simple", "medium", "complex"
	RequireLetter    bool         `json:"require_letter"`
	RequireDigit     bool         `json:"require_digit"`
	RequireSymbol    bool         `json:"require_symbol"`
	ExcludePatterns  []string     `json:"exclude_patterns"`  // Patterns to exclude (e.g., "password123")
	OnlyASCII        bool         `json:"only_ascii"`
	MinUniqueChars   int          `json:"min_unique_chars"`  // Minimum unique characters
	CustomRules      []CustomRule `json:"custom_rules"`      // Custom validation rules
}

// LimitsConfig defines explosion control
type LimitsConfig struct {
	MaxOutputLines int64  `json:"max_output_lines"` // 0 = unlimited
	MaxFileSize    int64  `json:"max_file_size"`    // bytes (stored internally), 0 = unlimited
	MaxMutations   int    `json:"max_mutations"`
	WarnThreshold  int64  `json:"warn_threshold"`   // warn when exceeding this
}

// OutputConfig defines output options
type OutputConfig struct {
	FilePath        string `json:"file_path"`
	Format          string `json:"format"` // "plain", "hashcat", "hydra", "john", "aircrack"
	Username        string `json:"username"` // For hydra/john username:password format
	Hash            string `json:"hash"`     // For john hash:password format
	Sort            bool   `json:"sort"`
	Shuffle         bool   `json:"shuffle"`
	ProbabilityOrder bool `json:"probability_order"`
	Compress        bool   `json:"compress"` // Compress output with gzip
	SplitFiles      bool   `json:"split_files"` // Split into multiple files
	MaxFileSize     int64  `json:"max_file_size"` // Max size per file when splitting (bytes)
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		TargetPreset: TargetPresetConfig{
			Preset:       "none",
			MinLength:    4,
			MaxLength:    128,
			ASCIIOnly:    false,
			EnforceRules: true,
		},
		Input: InputConfig{
			Names:              []string{},
			Keywords:           []string{},
			NumericRanges:      []NumericRange{},
			Dates:              []DateRange{},
			ExternalFiles:      []string{},
			FreeText:           []string{},
			SelectedCountries:  []string{},
			UseAllCountryNames: true,
			IncludeBirthdays:   false,
			IncludeHolidays:    false,
			IncludeSpecialDates: false,
		},
		Numeric: NumericConfig{
			Enabled:      false,
			Start:        0,
			End:          9999,
			FixedLength:  0,
			LeadingZeros: false,
			Randomized:   false,
			Shuffle:      false,
			Hexadecimal:  false,
			Deterministic: false,
			Seed:         0,
		},
		Combination: CombinationConfig{
			Enabled:         true,
			MaxDepth:        3,
			MinLength:       4,
			MaxLength:       128,
			Patterns:        []Pattern{},
			WordWord:        true,
			WordNumber:      true,
			WordSymbolNumber: true,
			MultiWord:       false,
			UseTemplates:    false,
			SelectedTemplates: []string{},
			CustomTemplates:  []string{},
		},
		Mutation: MutationConfig{
			Enabled:         true,
			Leetspeak:       true,
			Unicode:         false,
			Intensity:       "medium",
			MaxMutations:    3,
			HashcatRules:    false,
			HashcatRulesFile: "",
		},
		Randomization: RandomizationConfig{
			Enabled:       false,
			ShuffleChars:  false,
			SwapWordOrder: false,
			InsertSymbols: false,
			Symbols:       "!@#$%^&*",
			Deterministic: false,
			Seed:          0,
		},
		Advanced: AdvancedConfig{
			CaseVariations:    true,
			Prefixes:          []string{},
			Suffixes:          []string{},
			RepeatedChars:     false,
			WordDuplication:   false,
			RandomSymbolInsert: false,
		},
		Deduplication: DeduplicationConfig{
			Enabled: true,
		},
		Filter: FilterConfig{
			Enabled:         false,
			MinComplexity:    "simple",
			RequireLetter:   false,
			RequireDigit:    false,
			RequireSymbol:   false,
			ExcludePatterns: []string{},
			OnlyASCII:       false,
			MinUniqueChars:  0,
			CustomRules:     []CustomRule{},
		},
		Limits: LimitsConfig{
			MaxOutputLines: 0,
			MaxFileSize:    0,
			MaxMutations:   10,
			WarnThreshold:  1000000,
		},
		CustomCharset: CustomCharsetConfig{
			Enabled:     false,
			CharsetName: "",
			CustomCharsets: []CustomCharsetDef{},
			MinLength:   4,
			MaxLength:   128,
		},
		Output: OutputConfig{
			FilePath:        "wordlist.txt",
			Format:          "plain",
			Username:        "",
			Hash:            "",
			Sort:            false,
			Shuffle:         false,
			ProbabilityOrder: false,
		},
	}
}

