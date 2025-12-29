# Architecture Overview

## System Design

The password wordlist generator is built with a clean separation between core logic and GUI, following a pipeline architecture.

## Pipeline Architecture

```
Input Sources
    ↓
Normalization (case, trimming)
    ↓
Combination Engine (word+word, word+number, patterns)
    ↓
Mutation Engine (leetspeak, Unicode)
    ↓
Randomization (shuffling, symbol insertion)
    ↓
Advanced Features (prefixes, suffixes, duplication)
    ↓
Deduplication (hash-based, memory-efficient)
    ↓
Output Streaming (direct to disk)
```

## Core Components

### 1. Input Sources (`core/input.go`)

Multiple input source types:
- **NameInput**: Names, usernames, nicknames
- **KeywordInput**: Custom keywords
- **NumericRangeInput**: Numeric ranges with formatting options
- **DateRangeInput**: Date ranges with various formats
- **FileInput**: External wordlist files
- **FreeTextInput**: Manual text input

All sources implement the `InputSource` interface with a `Words(ctx context.Context) <-chan string` method for streaming.

### 2. Numeric Generator (`core/numeric.go`)

Generates numeric-only passwords:
- Sequential or randomized generation
- Fixed or variable length
- Leading zero support
- Shuffling capability

### 3. Combination Engine (`core/combination.go`)

Generates password combinations:
- Word + Number: `password123`
- Word + Word: `passwordsecret`, `password_secret`
- Word + Symbol + Number: `password!123`
- Pattern-based: Custom templates like `{word}{number}`
- Multi-word chains: Up to configurable depth

### 4. Mutation Engine (`core/mutation.go`)

Applies character substitutions:
- **Leetspeak**: `a→@,4`, `e→3`, `i→1,!`, `o→0`, `s→$,5`
- **Unicode**: Look-alike substitutions (a→α,а)
- **Intensity Levels**: Low, medium, aggressive
- **Case Variations**: Lower, upper, title, mixed case

### 5. Randomization Engine (`core/randomization.go`)

Applies randomization:
- Character shuffling
- Word order swapping
- Random symbol insertion
- Deterministic mode with seed support

### 6. Deduplication (`core/deduplication.go`)

Memory-efficient deduplication:
- String-based for small datasets (< 1M entries)
- Hash-based (SHA256) for large datasets
- Streaming-safe implementation
- Thread-safe with mutex protection

### 7. Generator (`core/generator.go`)

Main orchestrator:
- Coordinates the pipeline
- Manages context for cancellation
- Handles progress callbacks
- Enforces limits (max lines, file size)
- Streams output directly to disk

## Configuration (`config/config.go`)

Centralized configuration structure:
- **InputConfig**: All input sources
- **NumericConfig**: Numeric generation settings
- **CombinationConfig**: Combination rules and patterns
- **MutationConfig**: Mutation settings
- **RandomizationConfig**: Randomization options
- **AdvancedConfig**: Advanced features
- **DeduplicationConfig**: Deduplication settings
- **LimitsConfig**: Explosion control
- **OutputConfig**: Output file and options

## GUI Layer (`gui/app.go`)

Fyne-based GUI with:
- Tabbed interface for different configuration sections
- Real-time progress updates
- Cancellation support
- File browser integration
- Input validation

## Concurrency Model

- **Goroutines**: Each pipeline stage runs in its own goroutine
- **Channels**: Communication between stages via buffered channels
- **Context**: Cancellation propagation through the pipeline
- **Atomic Operations**: Thread-safe progress tracking

## Memory Management

- **Streaming**: Never loads full wordlist into memory
- **Buffered Channels**: Small buffers (100 items) for pipeline stages
- **Hash-based Deduplication**: Uses SHA256 hashes for large datasets
- **Garbage Collection**: Go's GC handles cleanup automatically

## Performance Considerations

1. **Channel Buffering**: Buffered channels prevent blocking
2. **Batch Processing**: Some operations batch for efficiency
3. **Early Termination**: Context cancellation stops generation immediately
4. **Limit Enforcement**: Hard limits prevent resource exhaustion
5. **Progress Throttling**: Progress updates throttled to 100ms intervals

## Extensibility

The architecture supports easy extension:
- **New Input Sources**: Implement `InputSource` interface
- **New Mutation Rules**: Add methods to `MutationEngine`
- **New Patterns**: Extend `CombinationEngine` with new pattern types
- **New Output Formats**: Modify `Generator.Generate()` output handling

## Error Handling

- Context cancellation for graceful shutdown
- File I/O errors propagated to caller
- Invalid configuration handled at generation start
- Progress callbacks allow UI error display

## Testing Considerations

Each component can be tested independently:
- Input sources: Test word generation
- Mutations: Test substitution rules
- Combinations: Test pattern matching
- Deduplication: Test hash collisions
- Generator: Test full pipeline with mock inputs

