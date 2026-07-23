# go-ndjson

> Parse newline-delimited JSON into typed Go slices with optional per-line validation and format auto-detection.

## Why

NDJSON (one JSON object per line) is the standard streaming format for logs, events, and audit trails. Go's stdlib has no first-class reader for it. This library fills that gap with a generic, dependency-free API that validates each line as it parses and distinguishes empty input from all-blank input via distinct sentinel errors.

## Installation

```bash
go get github.com/larsartmann/go-ndjson
```

Requires Go 1.26+ with `GOEXPERIMENT=jsonv2` (the `encoding/json/v2` package is currently experimental).

## Quick start

```go
package main

import (
	"fmt"
	"strings"

	"github.com/larsartmann/go-ndjson"
)

type Event struct {
	EventType string `json:"event_type"`
	Phase     string `json:"phase"`
}

func main() {
	input := `{"event_type":"start","phase":"before"}
{"event_type":"end","phase":"after"}`

	events, err := ndjson.Read(strings.NewReader(input), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%d events\n", len(events)) // 2 events
}
```

## Usage

### With per-line validation

Pass a callback to validate each parsed object. The callback receives the 1-based line number:

```go
events, err := ndjson.Read(reader, func(lineNum int, e Event) error {
	if e.EventType != "start" && e.EventType != "end" {
		return fmt.Errorf("line %d: unknown event_type %q", lineNum, e.EventType)
	}
	return nil
})
```

### Format detection

The `loader` package detects whether raw bytes are a single JSON report or an NDJSON event stream:

```go
format, err := loader.Detect(data)
switch format {
case loader.FormatJSON:
	// single JSON object (has "version" key)
case loader.FormatNDJSON:
	// newline-delimited events (has "event_type" key)
}
```

### Sentinel errors

```go
events, err := ndjson.Read(reader, nil)
switch {
case errors.Is(err, ndjson.ErrEmpty):
	// input contained no data
case errors.Is(err, ndjson.ErrNoEvents):
	// all lines were blank
case errors.Is(err, ndjson.ErrOversizedLine):
	// a line exceeded 1 MB
}
```

## License

Proprietary — see [LICENSE](LICENSE).
