---
name: asc-release-successor
description: Replace legacy pilot/deliver-style workflows with asc CLI commands for CI and agent-driven releases.
---

# ASC Release Successor

Use this skill when a request mentions `pilot`, `deliver`, metadata/screenshot migration, or release automation.

## Workflow

1. Confirm command surface from help:
   - `asc publish --help`
   - `asc migrate --help`
   - `asc localizations --help`
   - `asc assets --help`
2. Resolve required IDs:
   - app: `asc apps`
   - builds: `asc builds list --app "APP_ID" --sort -uploadedDate --limit 5`
   - TestFlight groups: `asc testflight beta-groups list --app "APP_ID"`
3. Run the target flow:
   - TestFlight distribution: `asc publish testflight ...`
   - App Store submission: `asc publish appstore ... --submit --confirm`
   - Metadata migration: `asc migrate ...`
4. Add validation output:
   - Prefer JSON output for CI.
   - Use `--output table` or `--output markdown` only for human-readable reports.

## Workflow Mapping

| Legacy concept | ASC command |
|---|---|
| `pilot upload + distribute` | `asc publish testflight --app "APP_ID" --ipa "./app.ipa" --group "External Testers" --wait --notify` |
| `deliver upload + submit` | `asc publish appstore --app "APP_ID" --ipa "./app.ipa" --version "1.0.0" --wait --submit --confirm` |
| metadata upload | `asc localizations upload --version "VERSION_ID" --path "./localizations"` |
| metadata download | `asc localizations download --version "VERSION_ID" --path "./localizations"` |
| screenshot upload | `asc assets screenshots upload --version-localization "LOC_ID" --path "./screenshots" --device-type "IPHONE_65"` |
| preview upload | `asc assets previews upload --version-localization "LOC_ID" --path "./previews" --device-type "IPHONE_65"` |

## Guardrails

- Always use explicit long flags.
- Never assume command signatures; verify with `--help`.
- For destructive actions, require `--confirm`.
- For bulk listings, use `--paginate`.
- For `publish testflight --group`, IDs and names are accepted; if duplicate names exist, use explicit group IDs.
