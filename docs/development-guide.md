# Development Guide

## Prerequisites

- **Go 1.24+** - Required for building and running the application
- **Git** - For version control
- **Make** - For build automation (optional, can use `go` commands directly)

## Local Development Setup

### 1. Clone the Repository

```bash
git clone https://github.com/rudrankriyam/App-Store-Connect-CLI.git
cd App-Store-Connect-CLI
```

### 2. Install Dependencies

```bash
go mod download
go mod tidy
```

Or using Make:

```bash
make deps
```

### 3. Build the Binary

```bash
make build
```

This creates the `asc` binary in the project root.

Or directly:

```bash
go build -o asc .
```

### 4. Run Locally

```bash
./asc --help
```

## Development Commands

### Build

```bash
make build              # Standard build
make build-debug        # Build with debug symbols
make build-all          # Build for multiple platforms
```

### Testing

```bash
make test               # Run all tests
make test-coverage      # Run tests with coverage report
make test-integration   # Run integration tests (requires API credentials)
```

### Code Quality

```bash
make lint               # Lint code (uses golangci-lint if installed, else go vet)
make format             # Format code (gofmt + gofumpt)
```

### Development Workflow

```bash
make dev                # Format, lint, test, and build
```

## Integration Tests

Integration tests require App Store Connect API credentials and are skipped by default.

### Setup

Set environment variables:

```bash
export ASC_KEY_ID="YOUR_KEY_ID"
export ASC_ISSUER_ID="YOUR_ISSUER_ID"
export ASC_PRIVATE_KEY_PATH="/path/to/AuthKey.p8"
export ASC_APP_ID="YOUR_APP_ID"
```

Or use base64-encoded key:

```bash
export ASC_PRIVATE_KEY_B64="BASE64_KEY"
```

### Run Integration Tests

```bash
make test-integration
```

Or directly:

```bash
go test -tags=integration -v ./internal/asc -run Integration
```

## Local API Testing

Test API calls locally with real credentials:

```bash
export ASC_KEY_ID="YOUR_KEY_ID"
export ASC_ISSUER_ID="YOUR_ISSUER_ID"
export ASC_PRIVATE_KEY_PATH="/path/to/AuthKey.p8"
export ASC_APP_ID="YOUR_APP_ID"

./asc feedback --app "$ASC_APP_ID"
./asc crashes --app "$ASC_APP_ID"
./asc reviews --app "$ASC_APP_ID"
```

## Credential Storage

Credentials are stored in priority order:

1. **Environment variables** (highest priority)
2. **System keychain** (macOS/Windows/Linux)
3. **Config file** (`~/.asc/config.json` or `./.asc/config.json`)

### Login

```bash
asc auth login \
  --name "MyApp" \
  --key-id "ABC123" \
  --issuer-id "DEF456" \
  --private-key /path/to/AuthKey.p8
```

### Check Status

```bash
asc auth status
asc auth status --verbose
asc auth status --validate
```

### Diagnose Issues

```bash
asc auth doctor
asc auth doctor --output json
asc auth doctor --fix --confirm
```

**Important:** Never commit credentials or `.p8` files to version control.

## Project Structure

- **`main.go`** - Entry point
- **`cmd/`** - Root command setup
- **`internal/cli/`** - Command implementations
- **`internal/asc/`** - API client
- **`internal/auth/`** - Authentication
- **`internal/config/`** - Configuration

## Adding a New Command

1. Create command directory: `internal/cli/<command-name>/`
2. Implement command function: `func <Command>Command() *ffcli.Command`
3. Register in: `internal/cli/registry/registry.go`
4. Add tests: `*_test.go` files
5. Update README.md with usage examples

## Code Standards

See `docs/GO_STANDARDS.md` for detailed coding standards.

Key principles:
- Explicit flags (no short flags like `-a`, use `--app`)
- JSON output by default (minified)
- No interactive prompts
- Clean exit codes (0 success, 1 error)

## Testing Patterns

See `docs/TESTING.md` for testing guidelines.

- Unit tests: `*_test.go` files alongside source
- Integration tests: Tagged with `integration`, require API credentials
- Test patterns: Use standard Go testing patterns

## Build System

The Makefile provides common development tasks:

- `build` - Build binary
- `test` - Run tests
- `lint` - Lint code
- `format` - Format code
- `deps` - Install dependencies
- `clean` - Clean build artifacts
- `install` - Install binary to system
- `update-openapi` - Update OpenAPI paths index

## Pull Request Guidelines

- Keep PRs small and focused
- Add or update tests for new behavior
- Update `README.md` if behavior or scope changes
- Avoid committing credentials or `.p8` files
- Follow code standards in `docs/GO_STANDARDS.md`

## Security

If you find a security issue:
- Report responsibly via private issue
- Or contact the maintainer directly
- Do not commit credentials or sensitive data

## Common Development Tasks

### Update Dependencies

```bash
make update-deps
```

Or:

```bash
go get -u ./...
go mod tidy
```

### Update OpenAPI Index

```bash
make update-openapi
```

This runs `scripts/update-openapi-index.py` to update the API paths index.

### Install Binary

```bash
make install
```

Installs to `/usr/local/bin` by default. Override with `INSTALL_PREFIX`:

```bash
make install INSTALL_PREFIX=/custom/path
```

### Clean Build Artifacts

```bash
make clean
```

Removes:
- Binary files (`asc`, `asc-debug`)
- Build directories (`build/`, `dist/`)
- Coverage files (`coverage.out`, `coverage.html`)
