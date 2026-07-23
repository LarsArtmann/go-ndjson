package ndjson_test

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-ndjson"
)

// FuzzRead fuzzes the NDJSON reader with arbitrary input. The reader must
// never panic and must always return either valid events or a non-nil error.
func FuzzRead(f *testing.F) {
	seeds := []string{
		`{"event_type":"start","phase":"before"}`,
		`{"event_type":"start","phase":"before"}` + "\n" +
			`{"event_type":"end","phase":"after"}`,
		"",
		"   \n\t\n  \n",
		"not json at all",
		"{broken",
		`{"event_type":"bogus_type","phase":"before"}`,
		`{"event_type":"start","phase":"bogus_phase"}`,
		"  " + `{"event_type":"end","phase":"after"}` + "  \n",
		strings.Repeat(`{"event_type":"start","phase":"before"}`, 100000),
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		events, err := ndjson.Read(strings.NewReader(input), validateTestEvent)
		if err != nil {
			if len(events) > 0 {
				t.Errorf("Read returned %d events alongside error %v", len(events), err)
			}

			return
		}

		for i, evt := range events {
			if evt.EventType != "" && !isValidEventType(evt.EventType) {
				t.Errorf("event %d has unknown event_type %q", i, evt.EventType)
			}

			if evt.Phase != "" && !isValidPhase(evt.Phase) {
				t.Errorf("event %d has unknown phase %q", i, evt.Phase)
			}
		}
	})
}

func isValidEventType(s string) bool {
	return s == "start" || s == "end"
}

func isValidPhase(s string) bool {
	return s == "before" || s == "after"
}
