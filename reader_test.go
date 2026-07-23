package ndjson_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/larsartmann/go-ndjson"
)

type testEvent struct {
	EventType string `json:"event_type"`
	Phase     string `json:"phase"`
}

func validateTestEvent(lineNum int, evt testEvent) error {
	knownEvents := map[string]bool{"start": true, "end": true}
	knownPhases := map[string]bool{"before": true, "after": true}

	if evt.EventType != "" && !knownEvents[evt.EventType] {
		return fmt.Errorf("line %d: unknown event_type: %q", lineNum, evt.EventType)
	}

	if evt.Phase != "" && !knownPhases[evt.Phase] {
		return fmt.Errorf("line %d: unknown phase: %q", lineNum, evt.Phase)
	}

	return nil
}

func TestRead_SingleEvent(t *testing.T) {
	t.Parallel()

	input := `{"event_type":"start","phase":"before"}`

	events, err := ndjson.Read(strings.NewReader(input), validateTestEvent)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	if events[0].EventType != "start" {
		t.Errorf("expected EventType %q, got %q", "start", events[0].EventType)
	}
}

func TestRead_MultipleEvents(t *testing.T) {
	t.Parallel()

	input := `{"event_type":"start","phase":"before"}` + "\n" +
		`{"event_type":"end","phase":"after"}`

	events, err := ndjson.Read(strings.NewReader(input), validateTestEvent)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}

	if events[0].EventType != "start" {
		t.Errorf("event 0: expected EventType %q, got %q", "start", events[0].EventType)
	}

	if events[1].EventType != "end" {
		t.Errorf("event 1: expected EventType %q, got %q", "end", events[1].EventType)
	}
}

func TestRead_BlankLinesSkipped(t *testing.T) {
	t.Parallel()

	input := "\n\n" + `{"event_type":"start","phase":"before"}` + "\n\n\n"

	events, err := ndjson.Read(strings.NewReader(input), validateTestEvent)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 event after blank lines, got %d", len(events))
	}
}

func TestRead_EmptyInput(t *testing.T) {
	t.Parallel()

	_, err := ndjson.Read(strings.NewReader(""), validateTestEvent)
	if !errors.Is(err, ndjson.ErrEmpty) {
		t.Fatalf("expected ErrEmpty, got %v", err)
	}
}

func TestRead_OnlyBlankLines(t *testing.T) {
	t.Parallel()

	_, err := ndjson.Read(strings.NewReader("\n\n\n"), validateTestEvent)
	if !errors.Is(err, ndjson.ErrNoEvents) {
		t.Fatalf("expected ErrNoEvents, got %v", err)
	}
}

func TestRead_MalformedJSON(t *testing.T) {
	t.Parallel()

	input := `{"event_type":"start","phase":"before"}` + "\n" + "not json"

	_, err := ndjson.Read(strings.NewReader(input), validateTestEvent)
	if err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}

	if !strings.Contains(err.Error(), "line 2") {
		t.Errorf("expected error to reference line 2, got: %v", err)
	}
}

func TestRead_OversizedLine(t *testing.T) {
	t.Parallel()

	huge := strings.Repeat(`{"event_type":"start","phase":"before"}`, 100000)

	_, err := ndjson.Read(strings.NewReader(huge), validateTestEvent)
	if !errors.Is(err, ndjson.ErrOversizedLine) {
		t.Fatalf("expected ErrOversizedLine, got %v", err)
	}
}

func TestRead_ValidationRejectsUnknownEventType(t *testing.T) {
	t.Parallel()

	input := `{"event_type":"bogus","phase":"before"}`

	_, err := ndjson.Read(strings.NewReader(input), validateTestEvent)
	if err == nil {
		t.Fatal("expected validation error for unknown event_type, got nil")
	}

	if !strings.Contains(err.Error(), "bogus") {
		t.Errorf("expected error to mention %q, got: %v", "bogus", err)
	}
}

func TestRead_ValidationRejectsUnknownPhase(t *testing.T) {
	t.Parallel()

	input := `{"event_type":"start","phase":"bogus"}`

	_, err := ndjson.Read(strings.NewReader(input), validateTestEvent)
	if err == nil {
		t.Fatal("expected validation error for unknown phase, got nil")
	}

	if !strings.Contains(err.Error(), "bogus") {
		t.Errorf("expected error to mention %q, got: %v", "bogus", err)
	}
}

func TestRead_NilValidatorSkipsValidation(t *testing.T) {
	t.Parallel()

	// Unknown event_type would fail validation, but nil skips it.
	input := `{"event_type":"bogus","phase":"before"}`

	events, err := ndjson.Read[testEvent](strings.NewReader(input), nil)
	if err != nil {
		t.Fatalf("Read with nil validate failed: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
}

func TestRead_CarriageReturns(t *testing.T) {
	t.Parallel()

	input := "{\"event_type\":\"start\",\"phase\":\"before\"}\r\n"

	events, err := ndjson.Read(strings.NewReader(input), validateTestEvent)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 event with CR, got %d", len(events))
	}
}

func TestRead_NoTrailingNewline(t *testing.T) {
	t.Parallel()

	input := `{"event_type":"start","phase":"before"}`

	events, err := ndjson.Read(strings.NewReader(input), validateTestEvent)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 event without trailing newline, got %d", len(events))
	}
}
