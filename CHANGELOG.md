# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

## [0.0.1] - 2026-07-23

### Added

- `ndjson.Read[T]` generic function: parses newline-delimited JSON into a typed slice with optional per-line validation callback
- Blank line skipping with distinct sentinel errors: `ErrEmpty` (no data) vs `ErrNoEvents` (only blank lines)
- `ErrOversizedLine` sentinel error with 1 MB line cap (`MaxLineBytes`)
- `loader.Detect` format detection: classifies raw bytes as JSON Report (`"version"` key) vs NDJSON Event stream (`"event_type"` key)
- `loader.Format` type with `FormatAuto`, `FormatJSON`, `FormatNDJSON` constants
- Fuzz test (`FuzzRead`) ensuring the reader never panics on arbitrary input
