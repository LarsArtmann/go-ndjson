package ndjson

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// MaxLineBytes is the maximum allowed size for a single NDJSON line (1 MB).
const MaxLineBytes = 1 << 20

// ErrEmpty is returned when the NDJSON input contains no data.
var ErrEmpty = errors.New("ndjson input is empty")

// ErrNoEvents is returned when all lines in the input were blank.
var ErrNoEvents = errors.New("ndjson input contains no events")

// ErrOversizedLine is returned when a single NDJSON line exceeds MaxLineBytes.
var ErrOversizedLine = errors.New("ndjson line exceeds maximum size")

// Read parses line-delimited JSON objects from reader into a slice of T.
// Each line must be a single JSON-encoded object. Blank lines are skipped.
// Returns the parsed objects in order.
//
// The validate function is called for each parsed object with its 1-based line
// number. Pass nil to skip validation.
//
// Returns ErrEmpty if the input contains no data, ErrNoEvents if all lines
// were blank, or ErrOversizedLine if any line exceeds MaxLineBytes.
func Read[T any](reader io.Reader, validate func(lineNum int, v T) error) ([]T, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 0, MaxLineBytes), MaxLineBytes)

	var items []T

	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()

		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}

		var item T

		err := json.Unmarshal(line, &item)
		if err != nil {
			return nil, fmt.Errorf("ndjson line %d: %w", lineNum, err)
		}

		if validate != nil {
			err = validate(lineNum, item)
			if err != nil {
				return nil, err
			}
		}

		items = append(items, item)
	}

	err := scanner.Err()
	if err != nil {
		if errors.Is(err, bufio.ErrTooLong) {
			return nil, fmt.Errorf("%w (max %d bytes)", ErrOversizedLine, MaxLineBytes)
		}

		return nil, fmt.Errorf("scan ndjson: %w", err)
	}

	if len(items) == 0 {
		if lineNum == 0 {
			return nil, ErrEmpty
		}

		return nil, ErrNoEvents
	}

	return items, nil
}
