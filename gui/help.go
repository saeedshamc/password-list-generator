package gui

// GetHelpText returns help text for each tab
func GetHelpText(tabName string) string {
	helpTexts := map[string]string{
		"input": `INPUT SECTION
================================================================================

NAMES:
   Enter first names, last names, usernames
   Example: john, mary, admin, root
   One name per line

KEYWORDS:
   Enter custom keywords and terms
   Example: password, secret, key, company
   For WiFi: router, wifi, network
   For Server: server, admin, root
   One keyword per line

NUMERIC RANGES:
   Define ranges of numbers to include
   Format: Start-End
   Example: 0-9999, 1000-2000
   Used in combinations (word + number)

DATES:
   Select date ranges
   Format: YYYY-MM-DD to YYYY-MM-DD
   Example: 2020-01-01 to 2024-12-31
   Generates dates in various formats

FREE TEXT:
   Enter any custom text
   Example: company name, project name, team name
   One item per line

EXTERNAL FILES:
   Load words from external text files
   One word per line
   Supports multiple files
   Click "Add File" to add wordlist files

COUNTRY NAMES:
   Select countries to include common names
   Choose single or multiple countries
   Names are used as base words for mutations
   Enable "Use All Country Names" for complete dataset
   Or select specific countries from the list
   50+ countries available with common names
   Deduplication is mandatory when multiple countries selected

ADVANCED DATES:
   Include common birthdays, holidays, and special dates
   Birthdays: Common birth date patterns (YYYY-MM-DD, DD-MM-YYYY)
   Holidays: Major holidays and celebrations
   Special Dates: Notable dates (New Year, Valentine's Day, etc.)
   Automatically generates multiple date formats

KEYBOARD PATTERNS:
   Generate common keyboard-based patterns
   Includes: qwerty, asdf, 1234, zxcv, etc.
   Horizontal, vertical, and diagonal patterns
   Useful for realistic password generation
   Common sequences like "qwerty", "asdfgh", "123456"

COMMON PASSWORD LISTS:
   Integrate well-known password lists
   Built-in options: RockYou, SecLists
   Custom files: Load your own password lists
   Enable "Apply Mutations" to mutate common passwords
   Mutations make common passwords more effective
   Automatically loads and processes large lists

MASK ATTACK:
   Generate passwords using Hashcat-style masks
   Format: ?l (lowercase), ?u (uppercase), ?d (digit), ?s (symbol)
   Example: ?l?l?l?l?d?d?d (4 letters + 3 digits)
   Supports multiple mask patterns
   Each pattern generates all possible combinations
   Very powerful for targeted attacks

TIPS:
   - For WiFi: Enter network name (SSID), router, wifi
   - For Server: Enter admin, root, server
   - Combine multiple input sources for comprehensive wordlists
   - Use country names for region-specific attacks
   - Enable keyboard patterns for realistic passwords
   - Use mask attack for systematic brute force
   - Apply mutations to common passwords for better coverage`,

		"numeric": `NUMERIC SECTION
================================================================================

ENABLE NUMERIC MODE:
   Generates numeric-only passwords
   Useful for PINs, numeric passwords
   Example: 1234, 0000, 9999

RANGE SETTINGS:
   Start Number: Starting number (e.g., 0)
   End Number: Ending number (e.g., 9999)
   Example: 0-9999 generates 10,000 passwords

FIXED LENGTH:
   Set to force all numbers to same length
   Example: Length 4 → 0000, 0001, 0002, ... 9999
   0 = variable length

LEADING ZEROS:
   Pad numbers with leading zeros
   Example: 1 → 0001 (if length is 4)
   Only works with Fixed Length

HEXADECIMAL:
   Generate hexadecimal numbers (0-9, a-f)
   Example: 0000, 0001, ... ffff
   Works with Fixed Length and Leading Zeros
   Useful for hex-based passwords and keys
   Format: 0, 1, 2, ..., 9, a, b, c, d, e, f, 10, 11, ...

RANDOMIZED:
   Generate numbers in random order
   Instead of sequential (0, 1, 2, ...)
   Uses seed for reproducibility

SHUFFLE:
   Shuffle the generated sequence
   Different from randomized (shuffles after generation)

DETERMINISTIC:
   Use seed for reproducible generation
   Same seed = same sequence
   Useful for testing and resuming

USE CASES:
   - For PIN codes: 4-8 digits
   - For simple numeric passwords
   - For combination: password123
   - For hexadecimal: MAC addresses, hash prefixes, hex keys

TIPS:
   - Use fixed length for consistent password length
   - Enable hexadecimal for hex-based systems
   - Use leading zeros for formatted numbers
   - Combine with shuffle for randomized sequences`,

		"combination": `COMBINATION SECTION
================================================================================

COMBINATION ENGINE:
   Combines words, numbers, and symbols
   Creates realistic password patterns

MAX DEPTH:
   Maximum combination depth
   Example: Depth 2 = word + number
   Depth 3 = word + number + symbol
   Higher depth = more combinations
   Warning: Higher depth increases output size exponentially

MIN/MAX LENGTH:
   Filter combinations by length
   Min Length: Minimum password length
   Max Length: Maximum password length
   Passwords outside range are filtered
   Enforced by target preset rules

PATTERNS:
   Define custom combination patterns
   Format: {word}, {number}, {symbol}
   Example: {word}{number} = password123
   Example: {word}_{word} = password_secret
   Supports complex patterns with multiple placeholders

WORD + NUMBER:
   Combine words with numbers
   Example: password123, admin2024
   Uses common numbers (0-99, years, etc.)
   Prevents combinatorial explosion

WORD + WORD:
   Combine two words
   Example: passwordsecret, adminuser
   Creates compound passwords

WORD + SYMBOL + NUMBER:
   Combine word, symbol, and number
   Example: password!123, admin@2024
   Realistic password patterns

MULTI-WORD:
   Enable multi-word combinations
   Creates longer password chains
   Example: passwordsecretadmin
   Use with caution (can create very large wordlists)

PASSWORD TEMPLATES:
   Use predefined or custom templates
   Predefined templates: Common password patterns
   Custom templates: Define your own patterns
   Templates use placeholders: {word}, {number}, {symbol}, {date}
   Example template: {word}{number}{symbol}
   Generates: password123!, admin2024@, etc.

TIPS:
   - Start with lower depth for testing
   - Use templates for consistent patterns
   - Combine with mutations for variety
   - Monitor output size with preview mode`,

		"mutation": `MUTATION SECTION
================================================================================

MUTATION ENGINE:
   Applies character substitutions to passwords
   Makes passwords more realistic and varied

LEETSPEAK:
   Replaces letters with similar-looking numbers/symbols
   Examples:
   a → @, 4
   e → 3
   i → 1, !
   o → 0
   s → $, 5
   Example: password -> p@ssw0rd, p4ssw0rd

UNICODE:
   Uses Unicode look-alike characters
   Example: a → α (Greek alpha), а (Cyrillic a)
   Warning: May cause compatibility issues with some tools

INTENSITY LEVELS:
   Low: Minimal substitutions (a→@, e→3)
   Medium: Common substitutions (default)
   Aggressive: Maximum substitutions (all variants)

HASHCAT RULES:
   Apply Hashcat rule files (.rule format)
   Enable "Hashcat Rules" checkbox
   Select a .rule file path
   Rules will be applied to all generated passwords
   Example rules: : (nothing), l (lowercase), u (uppercase), c (capitalize)
   Supports advanced rules: append, prepend, insert, delete, substitute
   Very powerful for hashcat-specific transformations

JOHN THE RIPPER RULES:
   Apply John the Ripper rule files
   Enable "John Rules" checkbox
   Select a John rules file path
   Rules will be applied to all generated passwords
   Example rules: : (nothing), l (lowercase), u (uppercase), c (capitalize)
   Supports advanced rules: reverse, duplicate, reflect, rotate, insert, overwrite
   Very powerful for John the Ripper-specific transformations

MAX MUTATIONS:
   Limits how many mutations per password
   Prevents explosion of variants
   Recommended: 1-3 mutations

TIPS:
   - Use Hashcat rules for hashcat-specific transformations
   - Use John rules for John the Ripper-specific transformations
   - Combine with leetspeak for maximum coverage
   - Test with small wordlists first to verify rules work correctly
   - Rules are applied sequentially, order matters`,

		"advanced": `ADVANCED SECTION
================================================================================

CASE VARIATIONS:
   Converts to lowercase, uppercase, and mixed case
   Example: Password, PASSWORD, PaSsWoRd
   Applies all case variations to each password
   Can significantly increase output size

PREFIXES:
   Add symbols or text at the beginning of passwords
   Example: !password, @admin, #secret
   Enter one prefix per line
   Each prefix creates a new password variant
   Common prefixes: !, @, #, $, %

SUFFIXES:
   Add symbols or text at the end of passwords
   Example: password123, admin!, secret2024
   Enter one suffix per line
   Each suffix creates a new password variant
   Common suffixes: 123, 2024, !, @, #

REPEATED CHARACTERS:
   Repeat the last character
   Example: password -> passworddd
   Creates variation by duplicating final character

WORD DUPLICATION:
   Duplicate the entire word
   Example: password -> passwordpassword
   Creates longer passwords by repeating the word

RANDOM SYMBOL INSERTION:
   Insert random symbols into passwords
   Example: password -> p@ssw0rd
   Adds unpredictability to passwords

TIPS:
   - For WiFi: Use common years (2024, 2023) and numbers (123, 1234)
   - For Server: Use common symbols (!, @, #) as prefixes/suffixes
   - Combine prefixes and suffixes for maximum coverage
   - Use case variations with other mutations for variety
   - Monitor output size as advanced features can create large wordlists`,

		"output": `OUTPUT SECTION
================================================================================

OUTPUT DIRECTORY:
   Select where to save the wordlist
   Click "Select Directory" to choose a folder
   You can also enter the path manually

FILE NAME:
   Name of the output file
   Default: wordlist.txt
   You can use any name you want
   If file splitting is enabled, files will be named:
   wordlist_part1.txt, wordlist_part2.txt, etc.

OUTPUT FORMAT:
   Select the format for your cracking tool:
   
   Plain: Standard wordlist (one password per line)
   Example: password123
   
   Hashcat: Compatible with hashcat
   Format: One password per line
   Example: password123
   
   Hydra: Compatible with Hydra
   Format: username:password
   Example: admin:password123
   Note: Enter username in the username field
   
   John: Compatible with John the Ripper
   Format: username:password or hash:password
   Example: admin:password123 or $2a$10$hash:password123
   Note: Enter username or hash in the respective fields
   
   Aircrack: Compatible with Aircrack-ng (WPA/WPA2)
   Format: One password per line
   Example: password123

COMPRESSION:
   Enable Gzip compression to reduce file size
   Output file will have .gz extension
   Example: wordlist.txt.gz
   Useful for very large wordlists to save disk space

FILE SPLITTING:
   Automatically split large files into multiple parts
   Enable "Split Files" to activate
   Set "Max File Size" to control split size (MB)
   Files will be named: wordlist_part1.txt, wordlist_part2.txt, etc.
   Useful for managing very large wordlists

SETTINGS:
   Deduplication: Remove duplicate passwords
   Essential for large wordlists
   
   Max Output Lines: Limit number of lines
   0 = unlimited
   Prevents extremely large files
   
   Max File Size: Limit file size (MB)
   0 = unlimited
   Example: 1000 = 1GB, 100 = 100MB
   Also used for file splitting when enabled

RECOMMENDATIONS:
   - For brute force: No limits, enable compression
   - For testing: Set limits for faster generation
   - For WiFi cracking: Use Aircrack format
   - For hash cracking: Use Hashcat or John format
   - For very large lists: Enable file splitting and compression`,

		"bruteforce": `PROFESSIONAL BRUTE FORCE GUIDE
================================================================================

FOR WIFI (WPA/WPA2):
   Names: Enter network name (SSID)
   Keywords: wifi, router, network
   Numeric: 8 digits (for WPA2)
   Combination: Enable name + date + number
   Suffixes: Add 123, 2024, !

   Recommended Settings:
   - Enable Word + Number combination
   - Enable Word + Symbol + Number
   - Add common years (2024, 2023, 2022) to suffixes
   - Set numeric range: 0-9999
   - Use Aircrack output format
   - Enable compression for large wordlists
   - Use file splitting for very large lists

FOR SERVER:
   Names: admin, root, user
   Keywords: server, system, pass
   Combination: Enable admin + year
   Mutation: Enable Leetspeak
   Prefixes: Add !, @, #

   Recommended Settings:
   - Enable Leetspeak mutations
   - Add common symbols to prefixes
   - Use years in suffixes
   - Enable case variations
   - Use Hashcat or John output format
   - Apply Hashcat/John rules if available

FOR WEB APPLICATIONS:
   Names: admin, user, test
   Keywords: login, password, secret
   Combination: Enable word + number
   Mutation: Enable Leetspeak (medium intensity)
   Advanced: Enable case variations

   Recommended Settings:
   - Use Hydra format for username:password
   - Enable smart filtering for complexity
   - Apply custom rules for password policies
   - Use common password lists integration

FOR MOBILE/IOT:
   Names: Device names, brand names
   Keywords: default, admin, root
   Numeric: 4-6 digit PINs
   Combination: Simple patterns (word + number)

   Recommended Settings:
   - Use numeric mode for PINs
   - Enable keyboard patterns
   - Keep combinations simple
   - Use target preset: Mobile or IoT

ADVANCED TECHNIQUES:
   - Use Hashcat rules for hashcat-specific attacks
   - Use John rules for John the Ripper attacks
   - Apply mask attack for systematic brute force
   - Integrate common password lists (RockYou, SecLists)
   - Use country names for region-specific attacks
   - Enable keyboard patterns for realistic passwords
   - Use password templates for consistent patterns

PERFORMANCE OPTIMIZATION:
   - Enable compression for large wordlists
   - Use file splitting for very large outputs
   - Set appropriate limits to prevent explosion
   - Use preview mode to test configurations
   - Monitor statistics during generation
   - Use checkpoint/resume for long-running jobs

IMPORTANT TIPS:
   - Always enable Deduplication
   - Don't set limits for full brute force
   - Use dates and years in combinations
   - Test with small limits first
   - Use external wordlists for common passwords
   - Save configurations for reuse
   - Export statistics for analysis
   - Use target presets for automatic rule enforcement`,

		"filter": `FILTER SECTION
================================================================================

SMART FILTERING:
   Advanced filtering based on password characteristics
   Helps generate higher quality wordlists

COMPLEXITY LEVEL:
   Simple: Basic passwords (letters only)
   Medium: Mixed passwords (letters + numbers)
   Complex: Strong passwords (letters + numbers + symbols)
   Select minimum complexity required

CHARACTER TYPE REQUIREMENTS:
   Require Letter: Password must contain at least one letter
   Require Digit: Password must contain at least one digit
   Require Symbol: Password must contain at least one symbol
   Use these to enforce password policies

EXCLUDE PATTERNS:
   Enter patterns to exclude from output
   One pattern per line
   Supports regex patterns
   Example: ^admin$ (exact match), .*test.* (contains "test")

MINIMUM UNIQUE CHARACTERS:
   Minimum number of unique characters required
   Helps filter out repetitive passwords
   Example: "aaa111" has only 2 unique characters

ONLY ASCII:
   Filter to ASCII-only passwords
   Useful for systems that don't support Unicode

CUSTOM RULES:
   Define custom validation rules
   Add rules with name, pattern, and description
   Rules can require specific character classes
   Example: Require at least 2 uppercase letters
   Rules support regex patterns for advanced matching

PASSWORD ANALYSIS:
   The system automatically analyzes passwords for:
   - Common patterns (password, 123456, qwerty)
   - Keyboard walks (asdf, qwerty, 1234)
   - Repeating characters (aaa, 111)
   - Sequential patterns (abc, 123)
   - Entropy calculation
   - Character type distribution

WORD FREQUENCY:
   Tracks frequency of base words in passwords
   Helps identify most common password patterns
   Useful for prioritizing wordlist generation

USE CASES:
   - For strong passwords: Use Complex + all requirements
   - For filtering weak: Exclude common patterns
   - For specific targets: Set character requirements

TIPS:
   - Use filters to match target password policies
   - Combine multiple filters for precise control
   - Test filters with preview mode first
   - Use password analysis to understand password quality
   - Monitor word frequency to optimize wordlist effectiveness`,
	}
	
	if text, ok := helpTexts[tabName]; ok {
		return text
	}
	return "No help available for this section."
}
