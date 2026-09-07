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

## High Impact

| Task                                                  | Status     | Impact | Effort | Evidence                                                                                               |
| ----------------------------------------------------- | ---------- | ------ | ------ | ------------------------------------------------------------------------------------------------------ |
| Add CI workflow running `nix flake check` + `nix run .#test` | 🔴 `TODO` | High | 1h | No `.github/` directory exists; all gates already defined as flake apps (`flake.nix:88-135`) |
| Harden `loader.Detect` key probing against edge cases | 🔴 `TODO` | High | 2h | Probes test non-empty values, not key presence: `{"version":""}` falls through to NDJSON (`loader/format.go:81`); object with neither key silently defaults to NDJSON (`loader/format.go:90`); unparseable first line is misclassified as JSON (`loader/format.go:78`); only the first non-blank line is ever probed |

## Medium Impact

| Task                                     | Status     | Impact | Effort | Evidence                                                                                             |
| ---------------------------------------- | ---------- | ------ | ------ | ---------------------------------------------------------------------------------------------------- |
| Cover the untested `loader` branches     | 🔴 `TODO` | Med    | 30min  | `loader` at 87.5% coverage vs 96.4% root: `Format.String` default branch (`loader/format.go:38-39`) and scanner-error wrap (`loader/format.go:61-63`) are untested |

## Low Impact

| Task                            | Status     | Impact | Effort | Evidence                                                        |
| ------------------------------- | ---------- | ------ | ------ | --------------------------------------------------------------- |
| Add `BenchmarkRead` benchmark   | 🔴 `TODO` | Low    | 30min  | No `Benchmark*` functions exist; a parsing library should track per-line overhead |
