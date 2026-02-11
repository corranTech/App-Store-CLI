# API Error Fixtures Intake

This document defines how to collect and maintain real-world App Store Connect API error payloads for regression tests.

## Why This Exists

Apple error responses can include useful fields under flexible `meta` payloads (for example, `meta.associatedErrors`). These fields are not strictly modeled in the OpenAPI schema and can change over time.

A fixture corpus helps us catch parser and UX regressions before they reach users.

## Source of Payload Variants

Use one or more of the following:

- User reports in GitHub Issues (preferred for real production failures)
- Maintainer reproduction in throwaway ASC test apps
- Integration test runs that intentionally trigger known validation/state failures

When opening an issue, use the `API Error Payload Report` template so data is complete and sanitized.

## Required Data for Each Fixture

Collect these fields:

- `asc` version
- CLI command (exact invocation)
- HTTP status code
- Full sanitized JSON error body from ASC

Optional but useful:

- Endpoint path (if known)
- Expected CLI output behavior
- Reproduction notes

## Sanitization Rules

Before committing or sharing payloads:

- Remove/redact app IDs, version IDs, submission IDs, and other resource identifiers
- Remove/redact emails, names, free-form user text, and team/vendor identifiers
- Never include tokens, private keys, JWTs, headers, cookies, or credentials
- Keep structure and error wording intact so parser behavior remains realistic

Use placeholders like `app-<redacted>` or `version-<redacted>` instead of deleting keys.

## Fixture Storage Convention

Store fixtures under:

- `internal/asc/testdata/error_payloads/`

Naming pattern:

- `<status>-<area>-<scenario>.json`
- Example: `409-review-associated-errors-missing-age-rating.json`

Fixture content should be the raw ASC error response body JSON object, not wrapped in extra metadata.

## Foundation Models Repro Matrix

Use this matrix to intentionally trigger known failure classes for Foundation Models flows and store the exact sanitized ASC responses as fixtures.

As of this repository snapshot, there is no dedicated `foundation-models` CLI command yet. When that command or endpoint support lands, discover the exact commands via `asc --help` and keep the scenario naming below.

Target fixture folder:

- `internal/asc/testdata/error_payloads/foundation-models/`

Recommended scenario files:

- `400-foundation-models-validation-missing-required-field.json`
- `409-foundation-models-state-invalid-transition.json`
- `409-foundation-models-relationship-invalid-reference.json`
- `403-foundation-models-permission-insufficient-role.json`

For each scenario:

- Reproduce the failure in a throwaway app/project
- Capture the full non-2xx error body JSON
- Sanitize identifiers/secrets using the rules above
- Commit the sanitized payload as-is under the scenario filename
- Add parser/command tests that load the fixture and assert actionable output

## Test Coverage Expectations

For each new fixture variant:

- Add/extend parser tests in `internal/asc/client_test.go`
- Assert both parsed structure and rendered error message content
- Add command-level coverage in `internal/cli/cmdtest/` when user-facing formatting matters
- Keep malformed-shape coverage to ensure graceful fallback behavior

## Maintenance

- Prefer adding new fixtures over changing existing fixtures unless the old one is invalid
- Keep fixtures small, focused, and scenario-specific
- If a fixture required aggressive redaction, include a short note in the related test explaining any intentional placeholder content
