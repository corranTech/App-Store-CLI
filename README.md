# Unofficial App Store Connect CLI

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-yellow?style=for-the-badge" alt="License">
  <img src="https://img.shields.io/badge/Homebrew-compatible-blue?style=for-the-badge" alt="Homebrew">
</p>

A **fast**, **lightweight**, and **AI-agent friendly** CLI for App Store Connect. Built for automation: explicit flags, JSON-first output, and clean exit codes.

## Why ASC?

| Problem | Solution |
|---------|----------|
| Manual App Store Connect work | Automate everything from the CLI |
| Slow, heavy tooling | Single Go binary, fast startup |
| Hard to script/agent | Default JSON output, explicit flags, no prompts |

## Table of Contents

- [Quick Start](#quick-start)
- [Discover Commands](#discover-commands)
- [Automation Features](#automation-features)
- [Wait / Completion / Progress](#wait--completion--progress)
- [Docs](#docs)
- [Security](#security)
- [Contributing](#contributing)
- [License](#license)

## Quick Start

### Install

```bash
# Homebrew (recommended)
brew tap rudrankriyam/tap
brew install rudrankriyam/tap/asc

# Install script (macOS/Linux)
curl -fsSL https://raw.githubusercontent.com/rudrankriyam/App-Store-Connect-CLI/main/install.sh | bash

# Or build from source
git clone https://github.com/rudrankriyam/App-Store-Connect-CLI.git
cd App-Store-Connect-CLI
make build
./asc --help
```

### Authenticate

Generate keys in App Store Connect:
[App Store Connect API Keys](https://appstoreconnect.apple.com/access/integrations/api)

```bash
# Save an API key profile (recommended)
asc auth login \
  --name "MyApp" \
  --key-id "ABC123" \
  --issuer-id "DEF456" \
  --private-key /path/to/AuthKey.p8

# Check auth status
asc auth status --verbose

# Use a profile for a single command
asc --profile "MyApp" apps list
```

## Discover Commands

This CLI is designed to be self-documenting:

```bash
asc --help
asc builds --help
asc builds list --help
```

## Automation Features

- **JSON-first output**: Minified JSON by default. Use `--pretty` while debugging.
- **Human output when you want it**: `--output table` or `--output markdown`.
- **Pagination**: `--paginate` will fetch all pages automatically for list-style commands.
- **Retries (GET/HEAD)**: Built-in exponential backoff with jitter; configurable via env/config.
  - Enable retry logs to stderr with `--retry-log` or `ASC_RETRY_LOG=1`.
- **Timeouts**:
  - Requests: `ASC_TIMEOUT` / `ASC_TIMEOUT_SECONDS`
  - Uploads: `ASC_UPLOAD_TIMEOUT` / `ASC_UPLOAD_TIMEOUT_SECONDS`

## Wait / Completion / Progress

Some workflows support **polling-based waiting** so you can treat a command as “complete” only when the remote operation finishes.

### Publish workflows (upload + optional wait)

```bash
# Upload + distribute to TestFlight, and wait for processing
asc publish testflight --app "APP_ID" --ipa "./app.ipa" --group "GROUP_ID" --wait

# Control polling + overall timeout
asc publish testflight --app "APP_ID" --ipa "./app.ipa" --group "GROUP_ID" --wait --poll-interval 30s --timeout 45m
```

### Xcode Cloud (trigger + wait)

```bash
# Trigger and wait for completion (non-zero exit code on failure)
asc xcode-cloud run --app "APP_ID" --workflow "CI" --branch "main" --wait

# Wait for an existing build run to finish
asc xcode-cloud status --run-id "BUILD_RUN_ID" --wait --poll-interval 30s --timeout 1h
```

Notes:
- Waiting is implemented via polling; the **final status is returned in your selected output format**.
- When a waited build completes unsuccessfully, the command returns a **non-zero exit code** (good for CI).

## Docs

- **Start here**: `docs/index.md`
- **Local dev**: `docs/development-guide.md`
- **Architecture**: `docs/architecture.md`
- **API quirks**: `docs/API_NOTES.md`
- **Go standards**: `docs/GO_STANDARDS.md`
- **Testing**: `docs/TESTING.md`
- **Offline OpenAPI snapshot**: `docs/openapi/README.md`

## Security

- Credentials are stored in the **system keychain** when available, with a config fallback.
- Private key **content** is never committed; use paths or env values.
- Use `--strict-auth` / `ASC_STRICT_AUTH=1` to fail if credentials resolve from mixed sources.

## Contributing

Please read `CONTRIBUTING.md` before opening a PR.

## License

MIT License — see `LICENSE`.
