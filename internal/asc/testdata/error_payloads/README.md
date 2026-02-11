# API Error Payload Fixtures

Sanitized real-world App Store Connect API error bodies used for parser and CLI error-surfacing regression tests.

## Layout

- `foundation-models/` - Foundation Models repro scenarios

## Rules

- Store only sanitized JSON response bodies (raw ASC error object)
- Never include credentials, tokens, headers, cookies, or private user/team data
- Keep real error structure and wording intact
- Name files using: `<status>-<area>-<scenario>.json`

See `docs/API_ERROR_FIXTURES.md` for intake and maintenance workflow.
