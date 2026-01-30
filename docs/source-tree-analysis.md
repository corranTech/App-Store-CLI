# Source Tree Analysis

## Directory Structure

```
App-Store-Connect-CLI/
├── main.go                    # Application entry point
│
├── cmd/                       # Root command infrastructure
│   ├── root.go               # Root command definition and setup
│   ├── shared.go             # Shared utilities for commands
│   └── errors.go             # Error handling types
│
├── internal/                  # Internal application code
│   │
│   ├── cli/                   # CLI command implementations (177 Go files)
│   │   ├── registry/         # Command registry and registration
│   │   │   └── registry.go   # Subcommand registration (50+ commands)
│   │   │
│   │   ├── shared/           # Shared CLI utilities
│   │   │
│   │   ├── auth/             # Authentication commands (9 files)
│   │   │   └── auth.go       # Login, logout, status, doctor
│   │   │
│   │   ├── apps/             # App management commands
│   │   │   ├── apps.go       # List apps, get app details
│   │   │   ├── app_setup.go  # Post-create app setup automation
│   │   │   └── app_tags.go   # App tags management
│   │   │
│   │   ├── builds/           # Build management commands
│   │   │   └── builds_commands.go  # List, info, expire, upload, groups
│   │   │
│   │   ├── testflight/       # TestFlight operations
│   │   │   ├── testflight.go # TestFlight app management
│   │   │   ├── beta_groups.go # Beta group management
│   │   │   └── beta_testers.go # Beta tester management
│   │   │
│   │   ├── reviews/           # App Store reviews
│   │   │   ├── reviews.go    # List reviews, filter by stars/territory
│   │   │   └── review.go     # Respond to reviews, manage responses
│   │   │
│   │   ├── analytics/        # Analytics & sales reports
│   │   │   └── analytics.go  # Sales reports, analytics requests
│   │   │
│   │   ├── finance/          # Finance reports
│   │   │   └── finance.go    # Financial reports, region codes
│   │   │
│   │   ├── xcodecloud/       # Xcode Cloud workflows
│   │   │   ├── xcode_cloud.go # Workflow management
│   │   │   ├── xcode_cloud_workflows.go
│   │   │   ├── xcode_cloud_build_runs.go
│   │   │   └── ...           # Additional Xcode Cloud commands
│   │   │
│   │   ├── gamecenter/       # Game Center features
│   │   │   └── ...           # Achievements, leaderboards, sets
│   │   │
│   │   ├── subscriptions/    # In-App Purchase subscriptions
│   │   │   └── subscriptions.go
│   │   │
│   │   ├── versions/         # App Store versions
│   │   │   └── versions.go   # List, get, attach build, release
│   │   │
│   │   ├── localizations/    # Localization management
│   │   │   └── localizations.go
│   │   │
│   │   ├── devices/          # Device management
│   │   │   └── devices.go    # List, register, update devices
│   │   │
│   │   ├── sandbox/          # Sandbox testers
│   │   │   └── sandbox.go    # Sandbox tester management
│   │   │
│   │   ├── migrate/           # Fastlane compatibility
│   │   │   └── migrate.go    # Import/export metadata
│   │   │
│   │   ├── submit/           # App Store submission
│   │   │   └── submit.go     # Submit builds for review
│   │   │
│   │   ├── categories/       # App Store categories
│   │   │   └── categories.go
│   │   │
│   │   └── ...               # 30+ additional command groups
│   │
│   ├── asc/                   # App Store Connect API client (186 Go files)
│   │   └── [API endpoint implementations]
│   │       # Structured request/response handling
│   │       # Authentication integration
│   │       # Pagination support
│   │       # Error handling
│   │
│   ├── auth/                  # Authentication utilities (9 files)
│   │   └── [Keychain, config, JWT generation]
│   │
│   ├── config/                # Configuration management (2 files)
│   │   ├── config.go         # Config loading and parsing
│   │   └── config_test.go    # Configuration tests
│   │
│   └── itunes/                # iTunes Store API client (3 files)
│       ├── client.go
│       ├── client_test.go
│       └── storefronts.go
│
├── docs/                      # Project documentation
│   ├── API_NOTES.md          # API quirks and notes
│   ├── CONTRIBUTING.md        # Contribution guidelines
│   ├── GO_STANDARDS.md       # Go coding standards
│   ├── TESTING.md            # Testing patterns
│   └── openapi/              # OpenAPI schema snapshots
│       ├── latest.json       # Latest OpenAPI schema
│       ├── paths.txt         # Quick index of API paths
│       └── README.md         # Update instructions
│
├── scripts/                   # Utility scripts
│   └── update-openapi-index.py  # Updates OpenAPI paths index
│
├── .github/workflows/         # CI/CD pipelines
│   ├── main-branch.yml       # Main branch workflow
│   ├── pr-checks.yml         # PR validation
│   └── release.yml           # Release automation
│
├── Makefile                   # Build system
├── go.mod                     # Go module definition
├── go.sum                     # Dependency checksums
├── README.md                  # Project README
├── CONTRIBUTING.md            # Contributing guidelines
├── LICENSE                    # MIT License
└── .gitignore                 # Git ignore patterns
```

## Critical Directories

### Entry Points
- **`main.go`** - Application entry point, calls `cmd.RootCommand()`
- **`cmd/root.go`** - Root command setup, registers all subcommands

### Command Infrastructure
- **`internal/cli/registry/`** - Central command registration
- **`internal/cli/shared/`** - Shared utilities for commands

### API Integration
- **`internal/asc/`** - App Store Connect API client (186 files)
- **`internal/itunes/`** - iTunes Store API client

### Authentication & Configuration
- **`internal/auth/`** - Authentication utilities (keychain, JWT)
- **`internal/config/`** - Configuration management

## Command Organization

Commands are organized by domain:
- **Authentication:** `auth/`
- **App Management:** `apps/`
- **Builds:** `builds/`
- **TestFlight:** `testflight/`
- **App Store:** `reviews/`, `versions/`, `categories/`
- **Analytics:** `analytics/`, `finance/`
- **Xcode Cloud:** `xcodecloud/`
- **Game Center:** `gamecenter/`
- **Subscriptions:** `subscriptions/`
- **Devices:** `devices/`
- **Sandbox:** `sandbox/`
- **Migration:** `migrate/`
- **Submission:** `submit/`

## File Patterns

- **Command files:** `*_command.go` or `*.go` in command directories
- **Test files:** `*_test.go` alongside source files
- **Shared utilities:** `shared*.go` or `shared/` directory
- **API clients:** Organized by API endpoint/resource

## Integration Points

- **Command → API Client:** Commands call `internal/asc/` for API operations
- **Command → Auth:** Commands use `internal/auth/` for credential management
- **Command → Config:** Commands use `internal/config/` for configuration
- **Root → Registry:** Root command uses `internal/cli/registry/` for subcommand registration
