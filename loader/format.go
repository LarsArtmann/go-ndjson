package loader

import (
	"bufio"
	"bytes"
	"encoding/json/v2"
	"encoding/json/jsontext"
	"errors"
	"fmt"
)

// maxScanBytes is the maximum scan buffer for format detection (1 MB).
const maxScanBytes = 1 << 20

// ErrNoContent is returned when no non-blank content is found for detection.
var ErrNoContent = errors.New("no content found for format detection")

// ErrUnknownFormat is returned when the first non-blank line cannot be
// classified as either format: it is not a JSON object, or it is an object
// carrying neither the "version" nor the "event_type" key.
var ErrUnknownFormat = errors.New("cannot detect format")

// Format identifies the serialization format of an audit log file.
type Format int

const (
	// FormatAuto auto-detects JSON vs NDJSON by inspecting the first line.
	FormatAuto Format = iota
	// FormatJSON is a single JSON Report object (contains "version" key).
	FormatJSON
	// FormatNDJSON is newline-delimited Event objects (contains "event_type" key).
	FormatNDJSON
)

// String returns the human-readable format name.
func (f Format) String() string {
	switch f {
	case FormatAuto:
		return "auto"
	case FormatJSON:
		return "json"
	case FormatNDJSON:
		return "ndjson"
	default:
		return "unknown"
	}
}

// Detect inspects raw bytes to determine whether they contain a JSON report
// or NDJSON events by checking the first non-blank line.
//
// Key presence, not key value, decides: any top-level "version" key marks a
// JSON report, any "event_type" key marks NDJSON events, so `{"version":""}`
// is a report. When both keys are present, "event_type" wins, because events
// may carry their own "version" field while reports never carry "event_type".
//
// An unparseable first line is classified as a multi-line (pretty-printed)
// JSON report only if it starts with "{". Anything else — text, arrays,
// scalars, or an object with neither signature key — fails with an error
// wrapping [ErrUnknownFormat]. Only the first non-blank line is probed; later
// lines are never examined.
func Detect(data []byte) (Format, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(make([]byte, 0, maxScanBytes), maxScanBytes)

	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) > 0 {
			return detectLineFormat(line)
		}
	}

	err := scanner.Err()
	if err != nil {
		return FormatAuto, fmt.Errorf("scan for format detection: %w", err)
	}

	return FormatAuto, ErrNoContent
}

// detectLineFormat inspects a single JSON line for Report vs Event keys.
func detectLineFormat(line []byte) (Format, error) {
	var fields map[string]jsontext.Value

	err := json.Unmarshal(line, &fields)
	if err != nil {
		if bytes.HasPrefix(line, []byte("{")) {
			// Incomplete first line of a multi-line (pretty-printed) JSON report.
			return FormatJSON, nil
		}

		return FormatAuto, fmt.Errorf("%w: first line is not a JSON object: %q", ErrUnknownFormat, preview(line))
	}

	if _, ok := fields["event_type"]; ok {
		return FormatNDJSON, nil
	}

	if _, ok := fields["version"]; ok {
		return FormatJSON, nil
	}

	return FormatAuto, fmt.Errorf("%w: first line has neither an %q nor a %q key: %q", ErrUnknownFormat, "event_type", "version", preview(line))
}

// preview returns data truncated to a bounded length for inclusion in error
// messages, so a hostile or oversized line cannot bloat the error.
func preview(data []byte) string {
	const maxPreviewBytes = 64
	if len(data) <= maxPreviewBytes {
		return string(data)
	}

	return string(data[:maxPreviewBytes]) + "..."
}
