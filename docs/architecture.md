# Architecture Documentation

## Executive Summary

ASC is a Go-based CLI application that provides a command-line interface to the App Store Connect API. The architecture is designed for simplicity, maintainability, and AI-agent friendliness.

## Technology Stack

### Core Dependencies

- **Go 1.24+** - Programming language
- **ffcli (peterbourgon/ff/v3)** - CLI framework for explicit flag-based commands
- **golang-jwt/jwt/v5** - JWT token generation for API authentication
- **99designs/keyring** - Secure credential storage (system keychain)
- **gopkg.in/yaml.v3** - YAML configuration parsing
- **howett.net/plist** - Apple Property List format handling

### Architecture Pattern

**Command-Based Hierarchical Architecture**

The application follows a clear command hierarchy:

```
asc (root)
├── auth (authentication)
├── apps (app management)
├── builds (build management)
├── testflight (TestFlight operations)
├── reviews (App Store reviews)
├── analytics (analytics & sales)
├── finance (finance reports)
├── xcode-cloud (Xcode Cloud workflows)
├── game-center (Game Center features)
└── ... (50+ command groups)
```

## Source Tree

```
App-Store-Connect-CLI/
├── main.go                    # Entry point
├── cmd/                       # Root command setup
│   ├── root.go               # Root command definition
│   ├── shared.go             # Shared utilities
│   └── errors.go             # Error handling
├── internal/
│   ├── cli/                   # CLI command implementations
│   │   ├── registry/         # Command registry
│   │   ├── auth/             # Authentication commands
│   │   ├── apps/             # App management commands
│   │   ├── builds/           # Build commands
│   │   ├── testflight/       # TestFlight commands
│   │   ├── reviews/          # Review commands
│   │   ├── analytics/       # Analytics commands
│   │   ├── finance/          # Finance commands
│   │   ├── xcodecloud/       # Xcode Cloud commands
│   │   ├── gamecenter/       # Game Center commands
│   │   └── ...               # 50+ command groups
│   ├── asc/                   # App Store Connect API client
│   │   └── [186 Go files]    # API endpoint implementations
│   ├── auth/                  # Authentication utilities
│   ├── config/                # Configuration management
│   └── itunes/                # iTunes Store API client
├── docs/                      # Documentation
│   ├── API_NOTES.md
│   ├── CONTRIBUTING.md
│   ├── GO_STANDARDS.md
│   ├── TESTING.md
│   └── openapi/               # OpenAPI schema snapshots
├── scripts/                   # Utility scripts
│   └── update-openapi-index.py
├── .github/workflows/         # CI/CD pipelines
├── Makefile                   # Build system
├── go.mod                     # Go module definition
└── README.md                  # Project documentation
```

## Entry Points

### Application Entry Point

**File:** `main.go`

```go
func main() {
    os.Exit(run())
}

func run() int {
    versionInfo := fmt.Sprintf("%s (commit: %s, date: %s)", version, commit, date)
    root := cmd.RootCommand(versionInfo)
    defer cmd.CleanupTempPrivateKeys()
    
    if err := root.Parse(os.Args[1:]); err != nil {
        // Handle parsing errors
    }
    
    if err := root.Run(context.Background()); err != nil {
        // Handle execution errors
    }
    
    return 0
}
```

### Command Registration

**File:** `cmd/root.go`

The root command is created with all subcommands registered via `registry.Subcommands()`.

**File:** `internal/cli/registry/registry.go`

All commands are registered in a single array, maintaining display order and grouping.

## Command Structure

Each command follows a consistent pattern:

1. **Command Definition** - Creates `*ffcli.Command` with name, usage, help text
2. **Flag Set** - Defines command-specific flags
3. **Subcommands** - Optional nested command groups
4. **Exec Function** - Command execution logic that calls API client

### Example Command Structure

```go
func SomeCommand() *ffcli.Command {
    fs := flag.NewFlagSet("some-command", flag.ExitOnError)
    
    return &ffcli.Command{
        Name:       "some-command",
        ShortUsage: "asc some-command <subcommand> [flags]",
        ShortHelp:  "Brief description",
        LongHelp:   "Detailed description with examples",
        FlagSet:    fs,
        Subcommands: []*ffcli.Command{
            // Subcommands
        },
        Exec: func(ctx context.Context, args []string) error {
            // Execution logic
        },
    }
}
```

## API Client Layer

**Location:** `internal/asc/`

The API client layer provides:
- **186 Go files** implementing App Store Connect API endpoints
- Structured request/response handling
- Authentication integration
- Error handling and retry logic
- Pagination support

## Authentication Architecture

**Location:** `internal/auth/`

Authentication supports multiple sources (priority order):
1. Environment variables (`ASC_KEY_ID`, `ASC_ISSUER_ID`, `ASC_PRIVATE_KEY_PATH`)
2. System keychain (macOS/Windows/Linux)
3. Config file (`~/.asc/config.json` or `./.asc/config.json`)

JWT tokens are generated on-demand for API requests.

## Configuration Management

**Location:** `internal/config/`

- Profile-based configuration (multiple API keys)
- Environment variable fallback
- Config file support with restricted permissions
- Strict auth mode (fail on multiple credential sources)

## Development Workflow

1. **Build:** `make build` - Compiles binary with version info
2. **Test:** `make test` - Runs unit tests
3. **Lint:** `make lint` - Code quality checks
4. **Format:** `make format` - Code formatting
5. **Integration Tests:** `make test-integration` - Requires API credentials

## Testing Strategy

- **Unit Tests:** Standard Go testing (`go test ./...`)
- **Integration Tests:** Opt-in tests requiring API credentials (tagged with `integration`)
- **Test Patterns:** `*_test.go` files alongside source code

## Build System

**Makefile** provides:
- `build` - Standard build
- `build-all` - Multi-platform builds
- `test` - Run tests
- `test-coverage` - Coverage reports
- `lint` - Code linting
- `format` - Code formatting
- `install` - Install binary

## Output Formats

The CLI supports three output formats:
- **JSON (default)** - Minified JSON for AI agents and scripting
- **Table** - Human-readable table format (`--output table`)
- **Markdown** - Markdown format for documentation (`--output markdown`)

## Error Handling

- **Exit Codes:** 0 for success, 1 for errors
- **Error Reporting:** Structured error types (`cmd.ReportedError`)
- **Help Display:** Automatic help on invalid commands/flags

## Security Considerations

- Credentials never stored in plain text
- System keychain preferred over config files
- Private keys referenced by path, not stored
- Environment variables as fallback
- Strict auth mode prevents credential conflicts
