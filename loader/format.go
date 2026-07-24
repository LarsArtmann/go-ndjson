package loader

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

// maxScanBytes is the maximum scan buffer for format detection (1 MB).
const maxScanBytes = 1 << 20

// ErrNoContent is returned when no non-blank content is found for detection.
var ErrNoContent = errors.New("no content found for format detection")

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
// A JSON report is identified by a top-level "version" key; an NDJSON event
// by an "event_type" key. Multi-line JSON (pretty-printed) that cannot be
// parsed as a single-line object defaults to FormatJSON.
func Detect(data []byte) (Format, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(make([]byte, 0, maxScanBytes), maxScanBytes)

	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) > 0 {
			return detectLineFormat(line), nil
		}
	}

	err := scanner.Err()
	if err != nil {
		return FormatAuto, fmt.Errorf("scan for format detection: %w", err)
	}

	return FormatAuto, ErrNoContent
}

// detectLineFormat inspects a single JSON line for Report vs Event keys.
func detectLineFormat(line []byte) Format {
	var probe struct {
		Version   string `json:"version"`
		EventType string `json:"event_type"`
	}

	err := json.Unmarshal(line, &probe)
	if err != nil {
		// Not valid single-line JSON — probably a multi-line JSON Report.
		return FormatJSON
	}

	if probe.Version != "" {
		return FormatJSON
	}

	if probe.EventType != "" {
		return FormatNDJSON
	}

	// Default: single-line object without version or event_type.
	return FormatNDJSON
}
