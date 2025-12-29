# Password Wordlist Generator

A high-performance, professional-grade password wordlist generator written in Go with a modern GUI interface. Designed for security researchers, penetration testers, and authorized security assessments.

## 🚀 Features

### Core Capabilities

- **Multiple Input Sources**: Names, keywords, numeric ranges, dates, external files, free text, and country-based name datasets (50+ countries)
- **Numeric-Only Mode**: Generate sequential or randomized numeric passwords with configurable formatting (decimal and hexadecimal)
- **Advanced Combination Engine**: Generate passwords using various patterns (word+number, word+word, word+symbol+number, multi-word chains)
- **Powerful Mutation Engine**: Leetspeak and Unicode look-alike substitutions with configurable intensity levels
- **Hashcat & John Rules Support**: Apply Hashcat (.rule) and John the Ripper rule files for advanced transformations
- **Keyboard Pattern Generation**: Generate common keyboard-based password patterns (qwerty, asdf, etc.)
- **Common Password Lists Integration**: Built-in support for RockYou, SecLists, and custom password lists
- **Mask Attack Support**: Generate passwords using Hashcat-style mask patterns (?l, ?u, ?d, ?s)
- **Password Templates**: Pre-defined and custom templates for consistent password patterns
- **Smart Filtering**: Advanced filtering based on complexity, character types, patterns, and custom rules
- **Password Analysis**: Automatic analysis of passwords (entropy, patterns, keyboard walks, etc.)
- **Word Frequency Analysis**: Track and analyze frequency of base words in generated passwords
- **Target-Aware Rules**: Enforce target-specific constraints (WiFi, SSH, Web, Mobile, IoT)
- **Memory-Efficient Deduplication**: Streaming-safe deduplication that scales to very large outputs
- **Output Formatting**: Support for multiple output formats (Plain, Hashcat, Hydra, John, Aircrack)
- **Compression Support**: Gzip compression for output files to save disk space
- **File Splitting**: Automatically split large files into multiple parts
- **Checkpoint/Resume**: Save and resume long-running generation jobs
- **Batch Processing**: Manage and execute multiple generation jobs concurrently
- **Multi-threaded Performance**: Parallel processing with worker pools for optimal performance
- **Statistics & Reporting**: Real-time statistics and exportable reports (JSON, HTML)
- **Preview Mode**: Generate and preview a sample of passwords before full generation
- **Save/Load Configuration**: Save and load configuration profiles for reuse

### Advanced Features

- **Custom Character Sets**: Define and apply custom character sets for filtering
- **Advanced Date Generation**: Generate common birthdays, holidays, and special dates
- **Priority Queue**: Output passwords by priority based on strength and analysis
- **Performance Optimization**: Automatic optimization of channels and buffers
- **Buffered I/O**: Optimized file writing for improved performance

## 📋 Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Configuration](#configuration)
- [Input Sources](#input-sources)
- [Numeric Generation](#numeric-generation)
- [Combination Engine](#combination-engine)
- [Mutation Engine](#mutation-engine)
- [Filtering & Validation](#filtering--validation)
- [Output Options](#output-options)
- [Target Presets](#target-presets)
- [Advanced Features](#advanced-features)
- [Performance Tips](#performance-tips)
- [Examples](#examples)
- [Architecture](#architecture)
- [Contributing](#contributing)
- [License](#license)

## 🔧 Installation

### Prerequisites

- Go 1.21 or later
- C compiler (for Fyne GUI dependencies)

### Build from Source

**Quick Build:**
```bash
./build.sh
```

**Manual Build:**
```bash
# Install dependencies
go mod download

# Build for current platform
go build -o passgen

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o passgen-linux

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o passgen.exe

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o passgen-macos
```

The build script (`build.sh`) will create a standalone binary with embedded assets.

## 🎯 Quick Start

### GUI Mode (Default)

```bash
./passgen
```

The GUI will launch with a default configuration. You can:
1. Configure input sources, combinations, mutations, etc.
2. Set output file path and format
3. Click "Generate" to start generation
4. Monitor progress in real-time
5. View statistics when complete

### CLI Mode

```bash
./passgen -cli -config config.json
```

Create a `config.json` file with your settings (see Configuration section).

## 📖 Usage

### GUI Interface

The GUI is organized into tabs:

1. **Input**: Configure input sources (names, keywords, dates, files, country names, etc.)
2. **Numeric**: Configure numeric-only password generation
3. **Combination**: Configure password combination patterns and templates
4. **Mutation**: Configure character mutations (leetspeak, Unicode, Hashcat/John rules)
5. **Advanced**: Configure prefixes, suffixes, case variations, etc.
6. **Filter**: Configure smart filtering and custom rules
7. **Output**: Configure output file, format, compression, and splitting

### Configuration Management

- **Save Config**: Save current configuration to a JSON file
- **Load Config**: Load a previously saved configuration
- **Preview**: Generate a sample of passwords to test your configuration
- **Export Stats**: Export generation statistics to JSON or HTML

### Target Presets

Select a target preset to automatically apply appropriate constraints:

- **None**: No specific constraints
- **WiFi**: ASCII only, length 8-63 (WPA/WPA2 compatible)
- **SSH**: ASCII only, length 4-128
- **Web**: ASCII only, length 4-128
- **Mobile**: Numeric preferred, length 4-8
- **IoT**: ASCII only, length 4-32

## ⚙️ Configuration

### Input Sources

```json
{
  "input": {
    "names": ["john", "admin", "root"],
    "keywords": ["password", "secret"],
    "numeric_ranges": [{"start": 0, "end": 9999}],
    "dates": [{"start": "2020-01-01", "end": "2024-12-31"}],
    "external_files": ["wordlist1.txt", "wordlist2.txt"],
    "free_text": ["company", "project"],
    "selected_countries": ["iran", "united_states"],
    "use_all_country_names": false,
    "include_birthdays": true,
    "include_holidays": true,
    "include_special_dates": true,
    "keyboard_patterns": true,
    "common_password_lists": ["rockyou", "seclists"],
    "common_password_files": ["custom_list.txt"],
    "apply_mutations_to_common": true,
    "mask_attack": false,
    "mask_patterns": ["?l?l?l?l?d?d?d"]
  }
}
```

### Numeric Generation

```json
{
  "numeric": {
    "enabled": false,
    "start": 0,
    "end": 9999,
    "fixed_length": 4,
    "leading_zeros": true,
    "randomized": false,
    "shuffle": false,
    "hexadecimal": false,
    "deterministic": false,
    "seed": 0
  }
}
```

### Combination Engine

```json
{
  "combination": {
    "enabled": true,
    "max_depth": 3,
    "min_length": 4,
    "max_length": 128,
    "patterns": [
      {"template": "{word}{number}", "enabled": true},
      {"template": "{word}_{word}", "enabled": true}
    ],
    "word_word": true,
    "word_number": true,
    "word_symbol_number": true,
    "multi_word": false,
    "use_templates": false,
    "selected_templates": [],
    "custom_templates": []
  }
}
```

### Mutation Engine

```json
{
  "mutation": {
    "enabled": true,
    "leetspeak": true,
    "unicode": false,
    "intensity": "medium",
    "max_mutations": 3,
    "hashcat_rules": false,
    "hashcat_rules_file": "",
    "john_rules": false,
    "john_rules_file": ""
  }
}
```

### Filtering

```json
{
  "filter": {
    "enabled": false,
    "min_complexity": "simple",
    "require_letter": false,
    "require_digit": false,
    "require_symbol": false,
    "exclude_patterns": ["password123", "admin"],
    "only_ascii": false,
    "min_unique_chars": 0,
    "custom_rules": [
      {
        "name": "Require 2 uppercase",
        "pattern": ".*[A-Z].*[A-Z].*",
        "description": "At least 2 uppercase letters",
        "enabled": true
      }
    ]
  }
}
```

### Output Options

```json
{
  "output": {
    "file_path": "wordlist.txt",
    "format": "plain",
    "username": "",
    "hash": "",
    "sort": false,
    "shuffle": false,
    "probability_order": false,
    "compress": false,
    "split_files": false,
    "max_file_size": 0
  }
}
```

## 📥 Input Sources

### Names
Enter names, usernames, or nicknames (one per line).

### Keywords
Enter custom keywords and terms (one per line).

### Numeric Ranges
Define ranges of numbers to include in combinations:
- Format: `{"start": 0, "end": 9999}`
- Used in combinations (e.g., `password123`)

### Dates
Select date ranges with various formats:
- Format: `{"start": "2020-01-01", "end": "2024-12-31"}`
- Generates dates in multiple formats (YYYY-MM-DD, DD-MM-YYYY, etc.)

### External Files
Load words from external text files:
- One word per line
- Supports multiple files
- Automatically processed and integrated

### Country Names
Select from 50+ countries with common names:
- Single or multiple country selection
- Names used as base words for mutations
- Deduplication mandatory when multiple countries selected

### Advanced Dates
- **Birthdays**: Common birth date patterns
- **Holidays**: Major holidays and celebrations
- **Special Dates**: Notable dates (New Year, Valentine's Day, etc.)

### Keyboard Patterns
Generate common keyboard-based patterns:
- Horizontal: `qwerty`, `asdf`
- Vertical: `147`, `258`
- Diagonal: `qaz`, `wsx`
- Common sequences: `123456`, `qwertyui`

### Common Password Lists
- **Built-in**: RockYou, SecLists
- **Custom**: Load your own password lists
- **Mutations**: Apply mutations to common passwords for better coverage

### Mask Attack
Generate passwords using Hashcat-style masks:
- `?l` = lowercase letter
- `?u` = uppercase letter
- `?d` = digit
- `?s` = symbol
- Example: `?l?l?l?l?d?d?d` generates 4 letters + 3 digits

## 🔢 Numeric Generation

### Features
- Sequential or randomized generation
- Fixed or variable length
- Leading zero support
- Hexadecimal number generation
- Shuffling capability
- Deterministic generation with seed

### Examples
- **PIN codes**: 4-8 digits with fixed length
- **Numeric passwords**: 0-9999 range
- **Hexadecimal**: For MAC addresses, hash prefixes, hex keys

## 🔗 Combination Engine

### Patterns
- **Word + Number**: `password123`, `admin2024`
- **Word + Word**: `passwordsecret`, `adminuser`
- **Word + Symbol + Number**: `password!123`, `admin@2024`
- **Multi-word**: `passwordsecretadmin`

### Templates
Use predefined or custom templates:
- Predefined: Common password patterns
- Custom: Define your own with placeholders
- Placeholders: `{word}`, `{number}`, `{symbol}`, `{date}`

### Depth Control
- **Max Depth**: Limits combination depth to prevent explosion
- **Min/Max Length**: Filter combinations by length
- **Pattern-based**: Define custom combination patterns

## 🎨 Mutation Engine

### Leetspeak
Character substitutions:
- `a → @, 4`
- `e → 3`
- `i → 1, !`
- `o → 0`
- `s → $, 5`

### Intensity Levels
- **Low**: Minimal substitutions
- **Medium**: Common substitutions (default)
- **Aggressive**: Maximum substitutions

### Hashcat Rules
Apply Hashcat rule files (`.rule` format):
- Supports all Hashcat rule commands
- Example: `:`, `l`, `u`, `c`, append, prepend, insert, delete, substitute

### John the Ripper Rules
Apply John the Ripper rule files:
- Supports all John rule commands
- Example: `:`, `l`, `u`, `c`, reverse, duplicate, reflect, rotate, insert, overwrite

## 🔍 Filtering & Validation

### Smart Filtering
- **Complexity Levels**: Simple, Medium, Complex
- **Character Requirements**: Require letter, digit, symbol
- **Exclude Patterns**: Regex-based pattern exclusion
- **Min Unique Chars**: Filter repetitive passwords
- **ASCII Only**: Filter to ASCII-only passwords

### Custom Rules
Define custom validation rules:
- Regex patterns
- Character class requirements
- Multiple rules with AND/OR logic

### Password Analysis
Automatic analysis includes:
- Common patterns detection
- Keyboard walk detection
- Repeating character detection
- Sequential pattern detection
- Entropy calculation
- Character type distribution

### Word Frequency
Track frequency of base words:
- Identify most common password patterns
- Optimize wordlist generation
- Prioritize high-frequency words

## 📤 Output Options

### Formats
- **Plain**: Standard wordlist (one password per line)
- **Hashcat**: Compatible with hashcat
- **Hydra**: `username:password` format
- **John**: `username:password` or `hash:password` format
- **Aircrack**: Compatible with Aircrack-ng (WPA/WPA2)

### Compression
- Enable Gzip compression to reduce file size
- Output file will have `.gz` extension
- Useful for very large wordlists

### File Splitting
- Automatically split large files into multiple parts
- Set maximum file size (MB)
- Files named: `wordlist_part1.txt`, `wordlist_part2.txt`, etc.

### Sorting & Shuffling
- **Sort**: Alphabetically sort passwords
- **Shuffle**: Randomly shuffle passwords
- **Probability Order**: Order by password strength/priority

## 🎯 Target Presets

### WiFi (WPA/WPA2)
- ASCII only
- Length: 8-63 characters
- Optimized for WPA/WPA2 cracking

### SSH
- ASCII only
- Length: 4-128 characters
- Optimized for SSH brute force

### Web
- ASCII only
- Length: 4-128 characters
- Optimized for web application attacks

### Mobile
- Numeric preferred
- Length: 4-8 characters
- Optimized for PIN/passcode attacks

### IoT
- ASCII only
- Length: 4-32 characters
- Optimized for IoT device attacks

## 🚀 Advanced Features

### Checkpoint/Resume
- Save generation progress periodically
- Resume from last checkpoint
- Useful for long-running jobs

### Batch Processing
- Manage multiple generation jobs
- Concurrent execution
- Job queue management

### Multi-threaded Performance
- Parallel processing with worker pools
- Automatic CPU detection and optimization
- Optimized channel buffering

### Statistics & Reporting
- Real-time generation statistics
- Export to JSON or HTML
- Password length distribution
- Character type distribution
- Total passwords and bytes generated

### Preview Mode
- Generate a sample of passwords
- Test configuration before full generation
- Quick feedback on output quality

## 💡 Performance Tips

1. **Use Limits**: Set appropriate limits to prevent explosion
2. **Enable Compression**: For very large wordlists
3. **Use File Splitting**: For managing large outputs
4. **Preview First**: Test with preview mode before full generation
5. **Save Configs**: Reuse configurations for similar attacks
6. **Monitor Statistics**: Track generation progress and quality
7. **Use Target Presets**: Automatic rule enforcement
8. **Optimize Input**: Use targeted input sources instead of broad ranges

## 📝 Examples

### WiFi Password Cracking

```json
{
  "target_preset": {"preset": "wifi"},
  "input": {
    "names": ["NetworkName"],
    "keywords": ["wifi", "router"],
    "dates": [{"start": "2020-01-01", "end": "2024-12-31"}]
  },
  "combination": {
    "word_number": true,
    "word_symbol_number": true
  },
  "advanced": {
    "suffixes": ["123", "2024", "!"]
  },
  "output": {
    "format": "aircrack",
    "compress": true
  }
}
```

### Server Brute Force

```json
{
  "target_preset": {"preset": "ssh"},
  "input": {
    "names": ["admin", "root", "user"],
    "keywords": ["server", "system"]
  },
  "mutation": {
    "leetspeak": true,
    "intensity": "medium"
  },
  "advanced": {
    "prefixes": ["!", "@", "#"],
    "suffixes": ["2024", "123"]
  },
  "output": {
    "format": "hashcat"
  }
}
```

### Custom Mask Attack

```json
{
  "input": {
    "mask_attack": true,
    "mask_patterns": ["?l?l?l?l?d?d?d", "?u?l?l?l?d?d?d?d"]
  },
  "output": {
    "format": "plain",
    "compress": false
  }
}
```

## 🏗️ Architecture

The system is designed as a pipeline:

```
Input Sources
    ↓
Normalization
    ↓
Combination Engine
    ↓
Mutation Engine (Leetspeak, Unicode, Hashcat/John Rules)
    ↓
Randomization
    ↓
Advanced Features (Prefixes, Suffixes, Case Variations)
    ↓
Validation (Target-Aware Rules)
    ↓
Smart Filtering
    ↓
Custom Rules
    ↓
Custom Character Set Filtering
    ↓
Deduplication
    ↓
Output Streaming (with Compression/Splitting)
```

### Key Design Principles

- **Core logic is independent of GUI**: Can be used programmatically
- **Streaming architecture**: Results written directly to disk, never fully loaded into memory
- **Concurrent processing**: Uses goroutines and channels for high performance
- **Modular design**: Each component is independent and testable
- **Memory efficient**: Scales to very large wordlists without memory issues

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is designed for authorized security testing and educational purposes only. Users are responsible for ensuring they have proper authorization before using this tool.

## ⚠️ Disclaimer

This tool is intended for:
- Authorized penetration testing
- Security research
- Educational purposes
- Personal security assessment (on systems you own)

**DO NOT** use this tool for unauthorized access to systems or networks. Unauthorized access is illegal and unethical.

## 📞 Support

For issues, questions, or contributions, please open an issue on the project repository.

---

**Built with ❤️ for the security community**
