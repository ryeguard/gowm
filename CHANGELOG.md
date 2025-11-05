# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Unified `gowm` binary with CLI and MCP server functionality
  - CLI commands: `gowm get-weather`, `gowm version`
  - MCP server: `gowm mcp` (stdio and HTTP modes)
  - Single binary simplifies installation and distribution

- Automated release infrastructure
  - GoReleaser configuration for cross-platform builds
  - Support for Linux, macOS, Windows (amd64 and arm64)
  - GitHub Actions workflow for automated releases on version tags
  - Automated changelog generation from commits
  - Binary checksums for verification

- Installation options
  - One-line installation script for macOS/Linux
  - Direct binary downloads from GitHub Releases
  - Go install support: `go install github.com/ryeguard/gowm/cmd/gowm@latest`
  - Comprehensive installation documentation

### Changed
- **BREAKING**: MCP server now accessed via `gowm mcp` subcommand
  - Previously: separate `gowm-mcp` binary
  - Now: unified `gowm` binary with `mcp` command
  - Claude Desktop config: use `"command": "/path/to/gowm"` with `"args": ["mcp"]`

### Benefits
- Single binary to install and manage
- Smaller download size (shared dependencies compiled once)
- Easier feature discovery via `gowm --help`
- Unified versioning across CLI and MCP
- Follows standard CLI tool patterns (docker, kubectl, gh)

[Unreleased]: https://github.com/ryeguard/gowm/commits/main

