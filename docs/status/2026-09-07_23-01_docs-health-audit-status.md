# Status Report: Docs-Health Audit — go-ndjson

**Snapshot:** 2026-09-07 23:01 CEST
**Scope:** This session's work only — full docs-health AUDIT (BUILD + VERIFY) of `github.com/larsartmann/go-ndjson`. No new research was done for this report.
**Verdict:** Docs set went from 7.75/10 Accuracy / 8/10 Fitness to healthy. Two must-have docs were missing entirely and the README shipped non-compiling code.

---

## a) FULLY DONE

Every item below was verified with a passing gate before being called done.

1. **Viewed ALL 14 project files** (2 source, 3 test, 6 docs, dprint.json, flake.nix, go.mod, .gitattributes) — no blind spots claimed.
2. **Quality gates established green baseline**: `nix run .#test` ✅, `.#lint` 0 issues ✅, `.#vet` ✅, `nix flake check` ✅, coverage 96.4% root / 87.5% loader / 92.3% total.
3. **Critical README fix**: both code examples used `ndjson.Read(reader, nil)`, which **does not compile** (`cannot infer T`) — proven via scratch build, fixed to `Read[Event](...)`, then re-verified by compiling and running the *exact* README code (output: `2 events`, `ErrEmpty ok`). README.md:37,83.
4. **AGENTS.md de-rotted**: removed temporal pollution ("Go 1.27 is not released yet") and an unverifiable claim ("jsonv2 in Go 1.24+"); added two gotchas I hit myself this session: gopls `[stdversion]` false positives, and the treefmt/dprint dual-formatter setup. AGENTS.md:20-33.
5. **CHANGELOG `[Unreleased]` filled** with the verified `go` directive relaxation (`1.26.4` → `1.26` since tag v0.0.1, confirmed via `git show v0.0.1:go.mod`).
6. **FEATURES.md citation precision fix** (`loader/format.go:54` → `:54-55`); all 17 feature rows independently re-verified against source line numbers — zero ghosts found.
7. **TODO_LIST.md built from scratch** — 4 items, every one verified against code (no CI: no `.github/` exists; `loader.Detect` edge cases at format.go:78/81/90; loader coverage gap; no benchmarks). Zero completed items retained.
8. **ROADMAP.md built from scratch** — 3 raw-idea themes (streaming, write API, generalized detection) + explicit non-goals; zero bounded tasks leaked in.
9. **Inline health report delivered** with visible math (Accuracy 10 − 1·1 − 0.5·2 − 0.25·1 = 7.75; Fitness 10 − 2·1 = 8, both 10 post-fix).
10. **Archive check executed**: no `docs/status/`, `docs/reviews/`, or historical snapshots exist — nothing to ANNOTATE or archive. Correct no-op, not a skipped step.

## b) PARTIALLY DONE

1. **dprint verification of the new/edited markdown** — the config exists (`dprint.json`) but the `dprint` binary is not on PATH and not wired into the flake, so table alignment in TODO_LIST.md could not be machine-checked. Works: all docs render correctly and tables are syntactically valid (5 columns throughout). Remains open: run `dprint fmt` once a runner exists. Effort: S.
2. **The "10/10 post-fix" score** — honestly: it is my own self-measurement, not an independent audit. The fixes are real and evidenced, but the score has a single verifier (me). Treat as "all named findings fixed", not "docs are perfect".

## c) NOT STARTED

All discovered this session, all correctly routed to TODO_LIST/ROADMAP (a docs audit should not silently start feature work). None blocked; none started.

1. CI workflow (no `.github/` at all) — in TODO_LIST.md.
2. `loader.Detect` edge-case hardening — in TODO_LIST.md.
3. Loader branch test coverage — in TODO_LIST.md.
4. `BenchmarkRead` — in TODO_LIST.md.
5. All ROADMAP raw ideas: streaming reader, `Write[T]`, generalized/configurable detection, `ReadContext`.
6. v0.0.2 release — the tag would now have a real `[Unreleased]` section to ship.
7. Archive/annotate workflow for this project's own future status reports (first one is being created right now — the workflow starts existing with this file).

## d) TOTALLY FUCKED UP

Radical honesty section. Nothing here blocks development, but all of it deserved to be named.

1. **The README shipped broken code since project creation.** Both examples called `Read(reader, nil)` — a guaranteed compile error for every first-time user following the quick start. Root cause: docs were written and maintained but never once *executed*. Workaround: none for a newcomer; they must read AGENTS.md (which contained the exact counter-knowledge — a split brain living in the same repo until this session).
2. **`dprint.json` is a ghost system.** Added deliberately (commit 7afa208), excluded from every automation: not in treefmt, not in devShell, binary not installed, nothing enforces it. Right now it is decoration pretending to be a formatting contract. Value: only if integrated (Question 2 below).
3. **`loader.Detect`'s core behaviors are untested accidents posing as an API.** Malformed first line → JSON (format.go:78); object with neither key → NDJSON (format.go:90); `{"version":""}` → falls through both probes → NDJSON. Zero tests pin any of these. A future refactor can silently change the "contract" and all tests stay green.
4. **Git history pollution by the auto-commit daemon** ("chore: auto-commit 2 changed file(s) (heuristic)" — commit 3cc2fd2 landed mid-session). Not authored by this session, but it obscures what actually changed.
5. **The `FormatAuto` constant is dead weight in practice**: it is only ever returned *together with an error* (format.go:62,65) — callers can never meaningfully receive it as a detection *result*. The public API surface advertises a state that cannot occur on the success path.

## e) WHAT WE SHOULD IMPROVE

1. **Docs audits must compile README examples systematically.** This session it was an ad-hoc idea — and it caught the only Critical finding. Make "every code block in README compiles and runs" a standing step, not luck.
2. **One formatting system, enforced.** Two formatter configs (treefmt for go/nix, dprint for md/json/yaml) where only one is wired anywhere is a drift machine waiting to happen. Pick one owner per file type and enforce it in `nix flake check`.
3. **Pin `Detect`'s fallback policy with explicit decision + tests** before any hardening refactor — otherwise the hardening PR becomes a silent behavior change.
4. **doc.go vs code nuance**: reader.go:24-26 says each line "must be a single JSON-encoded object", but the code happily parses arrays/primitives (`json.Unmarshal` into `T`). Either tighten the doc or embrace and document the generality.
5. **Harvest discipline**: this report's section (f) must be pulled into TODO_LIST/ROADMAP by the next docs-health HARVEST run, or it dies entombed in this timestamped file (the standing #1 failure mode).

## f) Top things to get done next (up to 50, brainstorm — mostly ROADMAP fuel)

Impact: Critical / High / Med / Low. Effort: S (<30min) / M (30min-2h) / L (>2h).

**Ship & infrastructure**

| #  | Task                                                                 | Impact | Effort | Category    | Home        |
|----|----------------------------------------------------------------------|--------|--------|-------------|-------------|
| 1  | Add GitHub Actions CI running `nix flake check` + `nix run .#test`    | High   | M      | Quality     | TODO_LIST ✓ |
| 2  | Cut v0.0.2 (Unreleased now non-empty; relaxes toolchain constraint)   | Med    | S      | Release     | TODO_LIST   |
| 3  | Wire dprint into flake (treefmt program or devShell app)              | Med    | S      | Cleanup     | TODO_LIST   |
| 4  | Add govulncheck job to CI (app already exists: `nix run .#vulncheck`) | Med    | S      | Quality     | TODO_LIST   |
| 5  | Add `nix run .#lint` + `.#vet` to CI matrix, not just test            | Med    | S      | Quality     | TODO_LIST   |
| 6  | Verify pkg.go.dev renders v0.0.2 after tag (proxy propagation)        | Low    | S      | Release     | TODO_LIST   |
| 7  | Add README CI badge once CI exists                                    | Low    | S      | Docs        | TODO_LIST   |
| 8  | Pin golangci-lint linter set in `.golangci.yml` (currently defaults)  | Low    | S      | Quality     | TODO_LIST   |
| 9  | Dependabot/Renovate or documented cadence for flake.lock bumps        | Low    | S      | Cleanup     | TODO_LIST   |
| 10 | Add GitHub issue templates (bug/feature)                              | Low    | S      | Docs        | TODO_LIST   |

**Correctness & API contract**

| #  | Task                                                                               | Impact | Effort | Category | Home        |
|----|-------------------------------------------------------------------------------------|--------|--------|----------|-------------|
| 11 | Switch `Detect` probes from non-empty-value to key-presence checking                | High   | M      | Bug      | TODO_LIST ✓ |
| 12 | Decide + test the neither-key default (currently silent NDJSON, format.go:90)       | High   | S      | Bug      | TODO_LIST ✓ |
| 13 | Decide + test malformed-first-line → JSON fallback (format.go:78)                   | High   | S      | Bug      | TODO_LIST ✓ |
| 14 | Test/document numeric `"version"` behavior (json/v2 unmarshal error → JSON branch)  | Med    | S      | Bug      | TODO_LIST   |
| 15 | Fix doc.go:24-26 "object" wording vs actual any-JSON-value behavior                 | Med    | S      | Docs     | TODO_LIST   |
| 16 | Wrap validate-callback errors with line number when caller omits it (reader.go:58-60)| Med    | S      | Feature  | TODO_LIST   |
| 17 | Reconsider `FormatAuto` in the success-path API (only ever returned alongside error)| Low    | S      | Cleanup  | TODO_LIST   |
| 18 | Distinguish parse errors from scan errors via sentinel wrapping                     | Low    | M      | Feature  | ROADMAP     |
| 19 | Test CRLF + blank-line + no-trailing-newline combination matrix                     | Low    | S      | Quality  | TODO_LIST   |

**Testing**

| #  | Task                                                         | Impact | Effort | Category | Home        |
|----|--------------------------------------------------------------|--------|--------|----------|-------------|
| 20 | Cover `Format.String` default branch (format.go:38-39)       | Med    | S      | Quality  | TODO_LIST ✓ |
| 21 | Cover loader scanner-error wrap (format.go:61-63)            | Med    | S      | Quality  | TODO_LIST ✓ |
| 22 | Add `FuzzDetect` (fuzzing covers only `Read` today)          | Med    | S      | Quality  | TODO_LIST   |
| 23 | Add `BenchmarkRead` + `BenchmarkDetect`                      | Low    | S      | Quality  | TODO_LIST ✓ |
| 24 | Roundtrip property test: `Detect(Write(x)) == NDJSON`        | Low    | S      | Quality  | ROADMAP     |
| 25 | `example_test.go` godoc examples (Go library convention)     | Med    | S      | Docs     | TODO_LIST   |

**Docs**

| #  | Task                                                                 | Impact | Effort | Category | Home      |
|----|----------------------------------------------------------------------|--------|--------|----------|-----------|
| 26 | Create `docs/DOMAIN_LANGUAGE.md` (Report vs Event, `version`, `event_type`) | Low | S   | Docs     | TODO_LIST |
| 27 | CONTRIBUTING.md: add release process section                          | Low    | S      | Docs     | TODO_LIST |
| 28 | README: "Error handling philosophy" short section (sentinels + wrapping) | Low | S      | Docs     | TODO_LIST |
| 29 | README: state the 1 MB line cap in the quick-start prose (currently only in errors section) | Low | S | Docs  | TODO_LIST |
| 30 | CHANGELOG: add keep-a-changelog compare links footer                  | Low    | S      | Docs     | TODO_LIST |
| 31 | AGENTS.md: revisit gopls gotcha after Go 1.27 ships (will become stale) | Low  | S      | Docs     | TODO_LIST |
| 32 | Make "compile every README example" a standing audit step (skill/AGENTS note) | Med | S   | Process  | TODO_LIST |

**Design ideas (ROADMAP fuel — unrefined by design)**

| #  | Task                                                                 | Impact | Effort | Category | Home      |
|----|----------------------------------------------------------------------|--------|--------|----------|-----------|
| 33 | Streaming read via `iter.Seq2[T, error]` (don't materialize `[]T`)   | High   | L      | Feature  | ROADMAP ✓ |
| 34 | `Write[T]` NDJSON writer counterpart                                 | Med    | M      | Feature  | ROADMAP ✓ |
| 35 | `Detect` from `io.Reader` without full buffering                     | Med    | M      | Feature  | ROADMAP ✓ |
| 36 | Configurable/pluggable detection probe keys (decouple audit-log vocab)| Med   | L      | Feature  | ROADMAP ✓ |
| 37 | `ReadContext` for cancellation mid-stream                            | Low    | M      | Feature  | ROADMAP ✓ |
| 38 | Combined Detect+Read convenience entry point                         | Low    | M      | Feature  | ROADMAP   |
| 39 | Configurable `MaxLineBytes` per call (currently package const only)  | Low    | S      | Feature  | ROADMAP   |
| 40 | Sampling more than first non-blank line for detection confidence     | Low    | M      | Feature  | ROADMAP ✓ |

(\✓ = already placed in TODO_LIST.md / ROADMAP.md this session. Items without ✓ are new leads from this report and need HARVEST routing before they count as tracked.)

**Deliberately not listed** (considered, rejected): JSON Schema integration (non-goal), CLI tooling (non-goal), third-party deps (non-goal), test framework additions (stdlib testing is a strength here), coverage-chasing microtests (92.3% with the named branch gaps handled in #20-21 is healthy).

## g) Top 3 questions I can NOT figure out myself

1. **Should `loader` stay audit-log-specific or become generic?** `Detect` hardcodes `"version"` / `"event_type"` (loader/doc.go says "audit log file formats") inside an otherwise generic `go-ndjson` library. Keeping it domain-flavored is fine — but it decides the entire hardening scope (TODO #11-14) and whether ROADMAP theme 3 lives or dies. What was the intent when it was extracted from auditlog-core?
2. **dprint.json: integrate or delete?** I can wire it into the flake (enforced like treefmt) or remove it. Ghost configs rot; I won't guess which you want (Question flagged: it's your tooling preference, not derivable from the repo).
3. **Are there downstream consumers of `Detect` whose current edge-case behavior must be preserved?** If auditlog-core (or anything else) already depends on "malformed → JSON" or "neither-key → NDJSON", hardening is a compat-constrained refactor; if not, we can fix the defaults properly while still pre-1.0. I cannot grep your consumers' intentions from this repo.

---

*First status report for this project — no prior baseline. All evidence from session 2026-09-07. Next docs-health HARVEST run should route section (f).*
