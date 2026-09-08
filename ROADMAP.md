# Roadmap

> Long-term direction and raw ideas. Items here are NOT actionable tasks.
> When an idea is refined into bounded work, it moves to TODO_LIST.md.

## Themes

### 1. Streaming and scale

`Read[T]` materializes every event into a `[]T` before returning. For
unbounded streams (long-lived log tails, event firehoses), callers need a way
to process events as they arrive instead of holding them all in memory.

Raw ideas:

- Iterator/callback-based streaming read (e.g. `iter.Seq2[T, error]`) as a
  sibling to the slice-returning `Read[T]`
- Context-aware reads so callers can cancel mid-stream
- `Detect` directly from an `io.Reader` (today it requires the whole payload
  as `[]byte`, forcing callers to buffer just to classify)
- Combined Detect+Read convenience entry point for the common
  "classify then parse" flow
- Configurable `MaxLineBytes` per call (today it is a package-level constant,
  `reader.go:13`)

### 2. Symmetric write API

The library reads NDJSON but has no story for producing it. A writer
counterpart would make this the complete NDJSON toolkit for the audit-log
ecosystem it was extracted from.

Raw ideas:

- `Write[T]` serializing a slice to newline-delimited JSON
- Streaming writer with the same blank-line and size-cap discipline the
  reader enforces
- Roundtrip property test once a writer exists: `Detect(Write(x)) == NDJSON`

### 3. Generalized format detection

`loader.Detect` hardcodes audit-log vocabulary (`"version"` / `"event_type"`
keys, see `loader/doc.go`). Useful for its home ecosystem, but generic
consumers may want detection rules that fit their own schemas.

Raw ideas:

- Configurable probe keys or pluggable classification predicates
- Richer heuristics (sampling more than the first non-blank line)

The ambiguous-input policy is no longer silent: `Detect` classifies by key
presence and fails with `ErrUnknownFormat` on ambiguous input (decided
2026-09-07, see `loader/format.go` and `CHANGELOG.md`).

Watch item: track `encoding/json/v2` API evolution as Go 1.27 nears (e.g. a
returning `RawMessage` could simplify the `map[string]jsontext.Value` probe
in `loader/format.go:83`).

### 4. Performance tuning (only if profiles justify it)

`BenchmarkRead` (~56 MB/s, ~1 alloc/event) and `BenchmarkDetect` exist, but
the library has never been profiled in anger. Raw ideas to evaluate only when
a real caller reports a hotspot:

- Token-based `jsontext` probe in `Detect` (no map allocation)
- Pre-sizing the result slice in `Read` heuristically

## Non-goals

Things we are deliberately NOT pursuing and why:

- **JSON Schema validation:** the per-line `validate` callback already gives
  callers full control; embedding a schema engine would balloon the API.
- **Third-party dependencies:** the stdlib-only, dependency-free property is
  a feature (`go.mod` has no requires).
- **CLI tooling:** this is a library; shells and pipelines already handle
  line-oriented workflows.
