# CLI Command Patterns

Keep CLI commands explicit, predictable, and easy to extend. This guide documents
the planned approach for reducing duplication while preserving clarity.

## Principles

- Keep explicit flags (`--app`, `--output`, `--pretty`, `--paginate`).
- Avoid interactive prompts; use `--confirm` for destructive actions.
- Validate required flags early and return `flag.ErrHelp`.
- Keep custom validation inline when it is domain-specific.

## Planned Refactor Phase: Command Builders + Reusable Groups

Phase 4 (planned): Introduce a light command builder plus reusable subcommand
groups for recurring patterns (localizations, images, releases, etc.).

Goals:

- Reduce boilerplate in list/get/create/update/delete commands.
- Keep output handling and common flag registration consistent.
- Avoid over-generalized CRUD interfaces that hide behavior.

Planned shape:

- Add `internal/cli/shared/command_builder.go` with a small builder for:
  - Common flags (`--id`, `--output`, `--pretty`, `--limit`, `--next`, `--paginate`)
  - Standard validation helpers
  - Output formatting
- Add reusable group builders for shared subcommand families (e.g., localizations).
- Keep domain-specific validation in the command implementation.

Expected impact:

- 30–40% fewer lines in heavily repeated command families.
- Faster, safer addition of new subcommands without losing clarity.
