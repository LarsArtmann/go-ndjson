# TODO List

> Short-term, actionable, bounded work items, verified against the actual code.
> For long-term vision and unrefined ideas, use ROADMAP.md.
> Items are ranked by impact. Status is verified, not assumed.

## Status legend

| Status           | Meaning                                                     |
| ---------------- | ----------------------------------------------------------- |
| 🔴 `TODO`        | Not started. Needs doing.                                   |
| 🟡 `IN_PROGRESS` | Actively being worked on.                                   |
| 🔵 `BLOCKED`     | Cannot proceed, external dependency or decision needed.     |
| 🟢 `DONE`        | Completed. Remove from this list and log in `CHANGELOG.md`. |

Harvested 2026-09-08 from `docs/status/2026-09-07_23-01_docs-health-audit-status.md`
and `docs/status/2026-09-07_23-34_todo-implementation-status.md` (both archived
after harvest). Every item re-verified against the code at harvest time.

## Items

| #   | Task                                                                                           | Impact | Effort | Status            | Evidence                                                                                                   |
| --- | ---------------------------------------------------------------------------------------------- | ------ | ------ | ----------------- | ---------------------------------------------------------------------------------------------------------- |
| 1   | Push and verify first green CI run; confirm README badge renders                               | High   | S      | 🟡 `IN_PROGRESS`  | First run failed at lint (`unknown GOEXPERIMENT jsonv2`); root cause fixed in `flake.nix:109-113,120-123`    |
| 2   | Cut v0.0.2 (`[Unreleased]` carries breaking `Detect` changes); verify pkg.go.dev after tag     | High   | M      | 🔴 `TODO`         | Consumer-compat checked 2026-09-08: `samber-do-auditlog/loader.go:81` wraps Detect errors, `go-workflow-auditlog` uses `Read` only |
| 3   | Shrink scanner initial buffer 1 MB → 64 KB in `Read` + `Detect`; verify with benchmarks        | High   | S      | 🔴 `TODO`         | `reader.go:35`, `loader/format.go:64`; `BenchmarkDetect` reports ~1.05 MB/op for one ~45-byte line          |
| 4   | Make `preview()` truncation rune-safe (no mid-rune slice mojibake)                             | Med    | S      | 🔴 `TODO`         | `loader/format.go:116-123` slices at byte 64                                                                |
| 5   | Add `FuzzDetect` seeded with the new edge cases (null/numeric values, both keys, non-objects)  | Med    | S      | 🔴 `TODO`         | Only `FuzzRead` exists (`fuzz_test.go`); parse path changed in `f5fe7a6`                                    |
| 6   | Wrap validate-callback errors with line number when the caller omits it                        | Med    | S      | 🔴 `TODO`         | `reader.go:56-60` returns the callback error unwrapped                                                      |
| 7   | Fix "must be a single JSON-encoded object" wording vs actual any-JSON-value behavior           | Med    | S      | 🔴 `TODO`         | `reader.go:24-26` vs `reader.go:51` (arrays/primitives parse fine)                                          |
| 8   | Add `example_test.go` godoc examples (Go library convention)                                   | Med    | S      | 🔴 `TODO`         | No `example_test.go` exists                                                                                 |
| 9   | Systematize README-example compilation (flake app or CI step), not ad-hoc per session          | Med    | S      | 🔴 `TODO`         | Two sessions caught real breakage this way; both depended on someone remembering                            |
| 10  | Record benchmark baselines + benchstat discipline for perf-touching PRs                        | Med    | S      | 🔴 `TODO`         | `BenchmarkRead`/`BenchmarkDetect` results currently live only in terminal scrollback                        |
| 11  | Add govulncheck job to CI (app already exists: `nix run .#vulncheck`)                          | Med    | S      | 🔴 `TODO`         | `.github/workflows/ci.yml` has no vulncheck step                                                            |
| 12  | `dprint.json`: wire into flake (treefmt program or app + CI) or delete it                      | Med    | S      | 🔵 `BLOCKED`      | User decision (tooling preference); ghost config flagged in 3 consecutive sessions                           |
| 13  | CI platform matrix (macOS/arm64) or document linux-only support                                | Low    | S      | 🔵 `BLOCKED`      | User decision; `nix flake check --all-systems` warns about omitted systems                                  |
| 14  | Decide duplicate-key first-line policy consciously (JSON guess vs `ErrUnknownFormat`)          | Low    | S      | 🔴 `TODO`         | `loader/format.go:87-89` brace-fallback currently classifies duplicate-key lines as JSON                    |
| 15  | Sentinel-based scan-error wrapping parity: loader wraps raw `bufio.ErrTooLong`, root maps it   | Low    | M      | 🔴 `TODO`         | `loader/format.go:75` vs `reader.go:68-70`                                                                  |
| 16  | Reconsider `FormatAuto` in the success-path API (only ever returned alongside an error)        | Low    | S      | 🔴 `TODO`         | `loader/format.go:75,78,92,107`                                                                             |
| 17  | CRLF + blank-line + no-trailing-newline combination matrix tests                               | Low    | S      | 🔴 `TODO`         | `TestRead_CarriageReturns` and `TestRead_NoTrailingNewline` exist separately, no matrix                     |
| 18  | Add `actionlint` as flake app + CI step so workflow changes are gated like code                | Low    | S      | 🔴 `TODO`         | `ci.yml` added 2026-09-07; only run ad-hoc so far                                                           |
| 19  | Pin golangci-lint linter set in `.golangci.yml` (currently defaults)                           | Low    | S      | 🔴 `TODO`         | No `.golangci.yml` exists                                                                                   |
| 20  | Dependabot/Renovate or documented cadence for `flake.lock` bumps                               | Low    | S      | 🔴 `TODO`         | `flake.lock` bumped manually so far                                                                         |
| 21  | Add GitHub issue templates (bug/feature)                                                       | Low    | S      | 🔴 `TODO`         | No `.github/ISSUE_TEMPLATE/` exists                                                                         |
| 22  | Create `docs/DOMAIN_LANGUAGE.md` (Report vs Event, `version`, `event_type`)                    | Low    | S      | 🔴 `TODO`         | Terms live in `loader/doc.go` only                                                                          |
| 23  | CONTRIBUTING.md: add release process section                                                   | Low    | S      | 🔴 `TODO`         | `CONTRIBUTING.md` has no release section                                                                    |
| 24  | README: error-handling philosophy section + state the 1 MB line cap in quick-start prose       | Low    | S      | 🔴 `TODO`         | Cap currently only in `README.md` sentinel-errors example                                                   |
| 25  | CHANGELOG: compare-links footer + decide `### Breaking` vs bold-marker convention              | Low    | S      | 🔴 `TODO`         | `CHANGELOG.md` has no link footer; breaking entries currently use bold markers                              |
| 26  | AGENTS.md: revisit gopls-gotcha note after Go 1.27 ships (it will become stale)                | Low    | S      | 🔵 `BLOCKED`      | Time-gated on Go 1.27 release                                                                               |
| 27  | `-cpu` benchmark sweep + benchstat CI regression check                                         | Low    | M      | 🔴 `TODO`         | Follows #10                                                                                                 |
