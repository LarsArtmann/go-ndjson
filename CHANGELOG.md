# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

### Added

- GitHub Actions CI (`.github/workflows/ci.yml`) running every flake gate: `nix flake check`, build, vet, lint, test, and test with race detector
- `BenchmarkRead` and `BenchmarkDetect` benchmarks tracking per-line parse and detection overhead

### Changed

- **Breaking:** `loader.Detect` now probes key presence instead of non-empty values: any top-level `"version"` key marks a JSON report (so `{"version":""}` is a report, not NDJSON), any `"event_type"` key marks NDJSON events; when both keys are present, `"event_type"` wins because events may carry their own `"version"` field
- **Breaking:** ambiguous input no longer silently defaults: an object with neither signature key, a non-object first line (text, arrays, scalars), or blank-line-only garbage fails with the new `loader.ErrUnknownFormat` sentinel; only unparseable first lines starting with `{` remain classified as multi-line JSON reports
- Loader coverage raised to 100% of statements (was 87.5%), total coverage 98.3%, pinning the `Format.String` default branch and the scanner-error wrap with tests
- Relaxed the `go` directive from `1.26.4` to `1.26`, so any Go 1.26.x toolchain can build the module

## [0.0.1] - 2026-07-23

### Added

- `ndjson.Read[T]` generic function: parses newline-delimited JSON into a typed slice with optional per-line validation callback
- Blank line skipping with distinct sentinel errors: `ErrEmpty` (no data) vs `ErrNoEvents` (only blank lines)
- `ErrOversizedLine` sentinel error with 1 MB line cap (`MaxLineBytes`)
- `loader.Detect` format detection: classifies raw bytes as JSON Report (`"version"` key) vs NDJSON Event stream (`"event_type"` key)
- `loader.Format` type with `FormatAuto`, `FormatJSON`, `FormatNDJSON` constants
- Fuzz test (`FuzzRead`) ensuring the reader never panics on arbitrary input
