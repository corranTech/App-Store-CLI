# Foundation Models Error Fixtures

This folder is reserved for sanitized real-world ASC error payloads captured from Foundation Models failure repros.

## Capture Targets

Add one fixture per scenario:

- `400-foundation-models-validation-missing-required-field.json`
- `409-foundation-models-state-invalid-transition.json`
- `409-foundation-models-relationship-invalid-reference.json`
- `403-foundation-models-permission-insufficient-role.json`

## Capture Requirements

- Use a throwaway app/project for repros
- Save the exact non-2xx ASC response body JSON (sanitized only)
- Keep schema/wording unchanged except redaction placeholders
- Do not invent or synthesize payloads for this folder

Once fixtures are added, add parser and command-level regression tests that load them.
