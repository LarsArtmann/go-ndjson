// Package loader detects audit log file formats (JSON report vs NDJSON events).
//
// [Detect] inspects the first non-blank line of raw bytes to determine whether
// the content is a single JSON report object (identified by a "version" key)
// or NDJSON event stream (identified by an "event_type" key).
package loader
