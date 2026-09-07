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

### 2. Symmetric write API

The library reads NDJSON but has no story for producing it. A writer
counterpart would make this the complete NDJSON toolkit for the audit-log
ecosystem it was extracted from.

Raw ideas:

- `Write[T]` serializing a slice to newline-delimited JSON
- Streaming writer with the same blank-line and size-cap discipline the
  reader enforces

### 3. Generalized format detection

`loader.Detect` hardcodes audit-log vocabulary (`"version"` / `"event_type"`
keys, see `loader/doc.go`). Useful for its home ecosystem, but generic
consumers may want detection rules that fit their own schemas.

Raw ideas:

- Configurable probe keys or pluggable classification predicates
- Richer heuristics (sampling more than the first non-blank line)
- Decide and document a principled policy for ambiguous inputs instead of the
  current silent defaults

## Non-goals

Things we are deliberately NOT pursuing and why:

- **JSON Schema validation:** the per-line `validate` callback already gives
  callers full control; embedding a schema engine would balloon the API.
- **Third-party dependencies:** the stdlib-only, dependency-free property is
  a feature (`go.mod` has no requires).
- **CLI tooling:** this is a library; shells and pipelines already handle
  line-oriented workflows.
