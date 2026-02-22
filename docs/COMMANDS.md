# Command Reference Guide

This page is the navigation layer for command discovery.
For authoritative flag behavior, always use:

```bash
asc --help
asc <command> --help
asc <command> <subcommand> --help
```

## Usage Pattern

```bash
asc <command> <subcommand> [flags]
```

## Global Flags

Common root flags:

- `--profile` - use a named authentication profile
- `--strict-auth` - fail when credentials are resolved from multiple sources
- `--debug` - debug logging to stderr
- `--api-debug` - HTTP debug logging to stderr (redacted)
- `--report` - CI report format (for example `junit`)
- `--report-file` - path to write report output
- `--retry-log` - enable retry logging
- `--version` - print version and exit

## Command Families

### Getting Started

- `auth` - manage API key authentication
- `doctor` - diagnose auth configuration
- `init` - generate helper docs in a repository
- `docs` - access embedded guides
- `install-skills` - install the asc skill pack

### Analytics and Finance

- `analytics`
- `insights`
- `finance`
- `performance`
- `feedback`
- `crashes`

### App Management

- `apps`
- `app-setup`
- `app-tags`
- `app-info`
- `app-infos`
- `versions`
- `localizations`
- `screenshots`
- `video-previews`
- `background-assets`
- `product-pages`
- `routing-coverage`
- `pricing`
- `pre-orders`
- `categories`
- `age-rating`
- `accessibility`
- `encryption`
- `eula`
- `agreements`
- `app-clips`
- `android-ios-mapping`
- `marketplace`
- `alternative-distribution`
- `nominations`
- `game-center`

### TestFlight and Builds

- `testflight`
- `builds`
- `build-bundles`
- `pre-release-versions`
- `build-localizations`
- `beta-app-localizations`
- `beta-build-localizations`
- `sandbox`

### Review and Release

- `review`
- `reviews`
- `submit`
- `validate`
- `publish`

### Monetization

- `iap`
- `app-events`
- `subscriptions`
- `offer-codes`
- `win-back-offers`
- `promoted-purchases`

### Signing

- `signing`
- `bundle-ids`
- `certificates`
- `profiles`
- `merchant-ids`
- `pass-type-ids`
- `notarization`

### Team and Access

- `account`
- `users`
- `actors`
- `devices`

### Automation

- `webhooks`
- `xcode-cloud`
- `notify`
- `migrate`

### Utility

- `version`
- `completion`
- `diff`
- `status`
- `release-notes`
- `workflow`
- `metadata`

## Scripting Tips

- JSON output is minified by default and optimized for machine parsing.
- Use `--output table` or `--output markdown` for human-readable output.
- Use `--paginate` on list commands to fetch all pages automatically.
- Use `--limit` and `--next` for manual pagination control.
- Prefer explicit flags and deterministic outputs in CI scripts.

## High-Signal Examples

```bash
# List apps
asc apps list --output table

# Upload a build
asc builds upload --app "123456789" --file "/path/to/MyApp.ipa"

# Validate and submit an App Store version
asc validate --app "123456789" --version "1.2.3"
asc submit --app "123456789" --version "1.2.3"

# Run a local automation workflow
asc workflow run --file .asc/workflow.json --workflow release
```

## Related Documentation

- [../README.md](../README.md) - onboarding and common workflows
- [API_NOTES.md](API_NOTES.md) - API-specific behavior and caveats
- [TESTING.md](TESTING.md) - test strategy and patterns
- [CONTRIBUTING.md](CONTRIBUTING.md) - contribution and dev workflow

