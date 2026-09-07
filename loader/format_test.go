package loader_test

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/larsartmann/go-ndjson/loader"
)

func TestDetect_JSONReport(t *testing.T) {
	t.Parallel()

	data := []byte(`{"version":"0.2.0","services":[]}`)

	format, err := loader.Detect(data)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if format != loader.FormatJSON {
		t.Errorf("expected FormatJSON, got %s", format)
	}
}

func TestDetect_NDJSON(t *testing.T) {
	t.Parallel()

	data := []byte(`{"event_type":"start","phase":"before"}` + "\n" +
		`{"event_type":"end","phase":"after"}`)

	format, err := loader.Detect(data)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if format != loader.FormatNDJSON {
		t.Errorf("expected FormatNDJSON, got %s", format)
	}
}

func TestDetect_BlankLinesBeforeContent(t *testing.T) {
	t.Parallel()

	data := []byte("\n\n  \n" + `{"version":"0.1.0"}`)

	format, err := loader.Detect(data)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if format != loader.FormatJSON {
		t.Errorf("expected FormatJSON after blank lines, got %s", format)
	}
}

func TestDetect_EmptyInput(t *testing.T) {
	t.Parallel()

	_, err := loader.Detect([]byte(""))
	if !errors.Is(err, loader.ErrNoContent) {
		t.Fatalf("expected ErrNoContent for empty input, got %v", err)
	}
}

func TestDetect_OnlyBlankLines(t *testing.T) {
	t.Parallel()

	_, err := loader.Detect([]byte("\n\n\n"))
	if !errors.Is(err, loader.ErrNoContent) {
		t.Fatalf("expected ErrNoContent for blank-only input, got %v", err)
	}
}

func TestDetect_MultiLineJSON(t *testing.T) {
	t.Parallel()

	// Pretty-printed JSON — first line is just "{".
	data := []byte("{\n  \"version\": \"0.2.0\",\n  \"services\": []\n}")

	format, err := loader.Detect(data)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if format != loader.FormatJSON {
		t.Errorf("expected FormatJSON for multi-line JSON, got %s", format)
	}
}

func TestDetect_KeyPresence(t *testing.T) {
	t.Parallel()

	// The signature keys decide by presence, not by value.
	tests := []struct {
		name  string
		input string
		want  loader.Format
	}{
		{"empty version value", `{"version":""}`, loader.FormatJSON},
		{"null version value", `{"version":null}`, loader.FormatJSON},
		{"numeric version value", `{"version":1}`, loader.FormatJSON},
		{"empty event_type value", `{"event_type":""}`, loader.FormatNDJSON},
		{"null event_type value", `{"event_type":null}`, loader.FormatNDJSON},
		// "event_type" wins when both keys are present: events may carry their
		// own "version" field, reports never carry "event_type".
		{
			"event with own version field",
			`{"event_type":"start","version":"0.2.0"}`,
			loader.FormatNDJSON,
		},
	}

	for _, tt := range tests {
		format, err := loader.Detect([]byte(tt.input))
		if err != nil {
			t.Errorf("%s: Detect failed: %v", tt.name, err)
			continue
		}

		if format != tt.want {
			t.Errorf("%s: got %s, want %s", tt.name, format, tt.want)
		}
	}
}

func TestDetect_UnknownFormat(t *testing.T) {
	t.Parallel()

	// Input that matches neither format fails with ErrUnknownFormat instead of
	// silently defaulting to a guess.
	tests := []struct {
		name  string
		input string
	}{
		{"object with neither key", `{"name":"job","status":"ok"}`},
		{"empty object", `{}`},
		{"plain text", "hello world"},
		{"array", `[1,2,3]`},
		{"number", `42`},
		{"string scalar", `"hello"`},
		{"null scalar", `null`},
	}

	for _, tt := range tests {
		_, err := loader.Detect([]byte(tt.input))
		if !errors.Is(err, loader.ErrUnknownFormat) {
			t.Errorf("%s: expected ErrUnknownFormat, got %v", tt.name, err)
		}
	}
}

func TestDetect_UnparseableFirstLine(t *testing.T) {
	t.Parallel()

	// An unparseable first line is a multi-line JSON report guess only when it
	// starts with "{"; anything else is ErrUnknownFormat, not a JSON default.
	tests := []struct {
		name  string
		input string
		want  loader.Format
	}{
		{"broken line starting with brace", `{"broken"`, loader.FormatJSON},
		// Degenerate duplicate keys cannot parse downstream either; the brace
		// fallback classifies them as JSON so the reader reports the real error.
		{"duplicate version keys", `{"version":"1","version":"2"}`, loader.FormatJSON},
	}

	for _, tt := range tests {
		format, err := loader.Detect([]byte(tt.input))
		if err != nil {
			t.Errorf("%s: Detect failed: %v", tt.name, err)
			continue
		}

		if format != tt.want {
			t.Errorf("%s: got %s, want %s", tt.name, format, tt.want)
		}
	}

	_, err := loader.Detect([]byte("not json at all"))
	if !errors.Is(err, loader.ErrUnknownFormat) {
		t.Errorf("malformed non-object line: expected ErrUnknownFormat, got %v", err)
	}
}

func TestDetect_UnknownFormatErrorPreviewsLine(t *testing.T) {
	t.Parallel()

	long := `{"pad":"` + strings.Repeat("x", 200) + `"}`

	_, err := loader.Detect([]byte(long))
	if !errors.Is(err, loader.ErrUnknownFormat) {
		t.Fatalf("expected ErrUnknownFormat, got %v", err)
	}

	msg := err.Error()
	if !strings.Contains(msg, strings.Repeat("x", 50)) {
		t.Errorf("expected error to include a preview of the offending line: %v", err)
	}

	if strings.Contains(msg, strings.Repeat("x", 57)) {
		t.Errorf("expected error preview to be truncated, got full line: %v", err)
	}
}

func TestDetect_OversizedLine(t *testing.T) {
	t.Parallel()

	huge := strings.Repeat("x", 1<<20+1)

	_, err := loader.Detect([]byte(huge))
	if !errors.Is(err, bufio.ErrTooLong) {
		t.Fatalf("expected bufio.ErrTooLong for oversized line, got %v", err)
	}

	if !strings.Contains(err.Error(), "scan for format detection") {
		t.Errorf("expected scanner error to carry detection context, got: %v", err)
	}
}

func TestFormat_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		format loader.Format
		want   string
	}{
		{loader.FormatAuto, "auto"},
		{loader.FormatJSON, "json"},
		{loader.FormatNDJSON, "ndjson"},
		{loader.Format(999), "unknown"},
	}

	for _, tt := range tests {
		got := tt.format.String()
		if got != tt.want {
			t.Errorf("Format(%d).String() = %q, want %q", tt.format, got, tt.want)
		}
	}
}

// BenchmarkDetect tracks first-line probing overhead over a fixed 1k-event file.
func BenchmarkDetect(b *testing.B) {
	var buf bytes.Buffer
	for i := range 1_000 {
		fmt.Fprintf(&buf, `{"event_type":"start","phase":"before","seq":%d}`+"\n", i)
	}

	payload := buf.Bytes()

	b.SetBytes(int64(len(payload)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		format, err := loader.Detect(payload)
		if err != nil {
			b.Fatal(err)
		}

		if format != loader.FormatNDJSON {
			b.Fatalf("got %s, want ndjson", format)
		}
	}
}
