# Project Overview

**Project Name:** App Store Connect CLI (ASC)  
**Type:** CLI (Command-Line Interface)  
**Language:** Go 1.24+  
**Architecture:** Monolith  
**Repository Type:** Single cohesive codebase

## Executive Summary

ASC is a fast, lightweight, and AI-agent friendly CLI for App Store Connect. It provides comprehensive automation for iOS app distribution workflows, enabling developers to manage TestFlight, App Store submissions, analytics, and more from the command line.

## Technology Stack

| Category | Technology | Version | Justification |
|----------|-----------|---------|---------------|
| Language | Go | 1.24+ | Fast compilation, single binary distribution, excellent for CLI tools |
| CLI Framework | ffcli (peterbourgon/ff/v3) | 3.4.0 | Explicit flags, no interactive prompts, AI-agent friendly |
| JWT | golang-jwt/jwt/v5 | 5.3.1 | App Store Connect API authentication |
| Keychain | 99designs/keyring | 1.2.2 | Secure credential storage |
| YAML | gopkg.in/yaml.v3 | 3.0.1 | Configuration file parsing |
| Property Lists | howett.net/plist | 1.0.1 | Apple-specific data format handling |

## Architecture Pattern

**Command-Based Architecture** - The application follows a hierarchical command structure:
- Root command (`asc`) delegates to subcommands
- Each subcommand is self-contained in `internal/cli/<command>/`
- Commands use the `ffcli` framework for flag parsing and execution
- API client layer (`internal/asc/`) provides App Store Connect API integration

## Repository Structure

- **Monolith** - Single cohesive codebase
- **Entry Point:** `main.go` → `cmd.RootCommand()`
- **Command Registry:** `internal/cli/registry/registry.go` registers all subcommands
- **API Client:** `internal/asc/` contains App Store Connect API client implementations

## Key Features

- **AI-Agent Friendly:** JSON output by default, explicit flags, clean exit codes
- **No Interactive Prompts:** All operations are flag-based for automation
- **Secure Credential Storage:** System keychain with config fallback
- **Comprehensive Coverage:** TestFlight, App Store, Analytics, Finance, Game Center, Xcode Cloud, and more
- **Fast Startup:** Go binary with minimal dependencies

## Getting Started

See [Development Guide](./development-guide.md) for local setup and development instructions.
