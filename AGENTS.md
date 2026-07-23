# AGENTS.md

Concise, enduring context for AI sessions working in `github.com/larsartmann/go-ndjson`.

## What This Is

A small Go **library** (not an application) for reading newline-delimited JSON. Two packages:

- **`ndjson`** (root) — generic streaming reader. `Read[T any]` parses NDJSON into `[]T` with an optional per-line validation callback. Skips blank lines, enforces a 1 MB per-line cap, distinguishes empty input from all-blank input via distinct sentinel errors.
- **`loader`** — format detection. `Detect([]byte)` inspects the first non-blank line to classify content as a single JSON Report object (has `"version"` key) vs an NDJSON Event stream (has `"event_type"` key). Pretty-printed/multi-line JSON defaults to `FormatJSON`.

## Critical Gotcha: Go Version Mismatch

The code imports **`encoding/json/v2`**, the new stdlib JSON v2 package promoted in **Go 1.27**. However, `go.mod` declares `go 1.26.4`, so:

- `go test ./...` and `go build` **FAIL** on Go 1.26.4 with `build constraints exclude all Go files in .../encoding/json/v2`.
- To compile, you need **Go 1.27+**. Either install it or bump the `go` directive in `go.mod` to `1.27` (with `GOTOOLCHAIN=auto`, Go will then fetch the right toolchain automatically).
- Do not "fix" this by downgrading the import to `encoding/json` (v1) without understanding why v2 was chosen.

## Commands

No `flake.nix`, `justfile`, or `Makefile` exists — use raw `go` commands:

```bash
go test ./...              # run all tests (requires Go 1.27+)
go test -race ./...        # with race detector (recommended; see CONTRIBUTING.md)
go test -fuzz=FuzzRead .   # run the fuzzer (time-box it; it runs indefinitely)
golangci-lint run ./...    # lint (referenced in CONTRIBUTING.md; no config file present)
```

## Architecture & Data Flow

```
raw bytes/io.Reader
      │
      ├── loader.Detect([]byte) ──► Format (JSON | NDJSON | Auto)
      │        probes first non-blank line for "version" vs "event_type"
      │
      └── ndjson.Read[T](io.Reader, validate) ──► ([]T, error)
               bufio.Scanner (1 MB cap) → json.Unmarshal per line → validate callback → collect
```

The two packages are independent: `loader` only classifies bytes; `ndjson` only parses. Callers compose them (detect format, then read if NDJSON).

## Conventions

- **Generics over interfaces**: `Read[T any]` is the core API. Type inference works for the reader argument; pass an explicit type parameter or a non-nil `validate` to pin `T` (see `TestRead_NilValidatorSkipsValidation` for the `Read[testEvent](..., nil)` form).
- **Optional validation via callback**: `validate func(lineNum int, v T) error` — `nil` skips validation entirely. Line numbers are 1-based.
- **Sentinel errors** checked with `errors.Is`:
  - Root: `ErrEmpty` (no data at all), `ErrNoEvents` (only blank lines), `ErrOversizedLine` (> 1 MB).
  - Loader: `ErrNoContent`.
- **Constants**: `MaxLineBytes = 1 << 20` (1 MB) is the public cap; `loader` has an unexported `maxScanBytes` twin.
- **Package docs**: each package has a `doc.go` with the package comment and sentinel-error inventory.

## Testing Patterns

- **External test packages**: `ndjson_test` and `loader_test` — tests live in `*_test.go` importing the package by path, not `package ndjson`.
- **`t.Parallel()`** on every test.
- **Fuzz test** (`fuzz_test.go`): `FuzzRead` asserts the reader never panics and never returns events alongside an error. Seed corpus covers happy paths, malformed JSON, blank-only, oversized, and unknown enum values.
- **Table-driven** where appropriate (e.g. `TestFormat_String`); otherwise one focused test per behavior.
- Test fixture type `testEvent` (`event_type`/`phase` fields) and its validator (`validateTestEvent`) are defined in `reader_test.go` and reused by the fuzz test.
