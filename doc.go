// Package ndjson reads newline-delimited JSON.
//
// NDJSON (Newline-Delimited JSON) is a streaming-friendly format where each
// line is a complete JSON object. The [Read] function parses NDJSON into a
// slice of any type T, with optional per-line validation.
//
// # Sentinel errors
//
//   - [ErrEmpty]: input contains no data
//   - [ErrNoEvents]: all lines were blank
//   - [ErrOversizedLine]: a line exceeds 1 MB
package ndjson
