# AGENTS.md

Non-obvious context for AI sessions working in `github.com/larsartmann/go-ndjson`.

## What This Is

A Go **library** (not an application) for reading newline-delimited JSON. Two packages, no dependencies:

- **`ndjson`** (root) — generic streaming reader. `Read[T any]` parses NDJSON into `[]T` with an optional per-line validation callback.
- **`loader`** — format detection. `Detect([]byte)` inspects the first non-blank line to classify content as JSON Report (has `"version"` key) vs NDJSON Event stream (has `"event_type"` key).

## Critical Gotcha: GOEXPERIMENT=jsonv2

The code imports **`encoding/json/v2`**, which is behind a build constraint in Go 1.26. It requires:

```bash
GOEXPERIMENT=jsonv2
```

Without it, `go build` and `go test` fail with `build constraints exclude all Go files in .../encoding/json/v2`. **This is NOT a Go version problem** — Go 1.27 is not released yet. The experiment flag unlocks the v2 package in Go 1.24+.

**The flake.nix handles this automatically** — all apps (`nix run .#test`, `nix run .#build`, etc.) and the devShell set `GOEXPERIMENT=jsonv2` for you. If running raw `go` commands outside the flake, you must `export GOEXPERIMENT=jsonv2` first.

## Commands

All commands go through `flake.nix`:

```bash
nix run .#test          # run tests (GOEXPERIMENT baked in)
nix run .#test-race     # with race detector
nix run .#build         # go build ./...
nix run .#lint          # golangci-lint run ./...
nix run .#vet           # go vet ./...
nix run .#coverage      # test + coverage report
nix run .#vulncheck     # govulncheck ./...
nix run .#clean         # trash coverage.out + clean test cache
nix flake check         # validate flake + treefmt formatting check
nix develop             # enter devShell (GOEXPERIMENT pre-set)
```

Raw `go` commands work if you set the flag: `GOEXPERIMENT=jsonv2 go test ./...`

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

- **Generics over interfaces**: `Read[T any]` is the core API. Pass an explicit type parameter when `validate` is nil (`Read[testEvent](reader, nil)`) since type inference needs the callback to pin `T`.
- **Sentinel errors** checked with `errors.Is`: `ErrEmpty`, `ErrNoEvents`, `ErrOversizedLine` (root); `ErrNoContent` (loader).
- **1 MB line cap**: `MaxLineBytes = 1 << 20` — public in root, unexported `maxScanBytes` twin in loader.
- **External test packages**: tests live in `*_test.go` with `package ndjson_test`, importing the library by path.
- **Fuzz test**: `FuzzRead` asserts the reader never panics and never returns events alongside an error.
