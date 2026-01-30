# Documentation Index

Start at the root [README.md](../README.md) for installation + quick examples, then use this index to go deeper.

## Start here

- **Project overview**: [Project Overview](./project-overview.md)
- **Architecture & code structure**: [Architecture](./architecture.md)
- **Local development**: [Development Guide](./development-guide.md)
- **Where things live**: [Source Tree Analysis](./source-tree-analysis.md)

## References

- **API quirks / gotchas**: [API Notes](./API_NOTES.md)
- **Go standards**: [GO Standards](./GO_STANDARDS.md)
- **Testing patterns**: [Testing Guide](./TESTING.md)
- **OpenAPI snapshot (offline)**: [OpenAPI README](./openapi/README.md)

## Contribution docs

- **Repo-level contributing**: [CONTRIBUTING.md](../CONTRIBUTING.md)
- **Docs contributing**: [docs/CONTRIBUTING.md](./CONTRIBUTING.md)

## Quick facts

- **Language**: Go 1.24+
- **CLI framework**: `ffcli` (peterbourgon/ff/v3)
- **Entry point**: `main.go` → `cmd.RootCommand()`
- **Key packages**:
  - `cmd/`: root command + shared wiring
  - `internal/cli/`: command implementations
  - `internal/asc/`: App Store Connect API client
  - `internal/auth/`, `internal/config/`: auth/config plumbing

---

**Last Updated:** 2026-01-30
