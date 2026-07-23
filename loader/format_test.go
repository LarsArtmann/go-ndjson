package loader_test

import (
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
	if err == nil {
		t.Fatal("expected error for empty input, got nil")
	}
}

func TestDetect_OnlyBlankLines(t *testing.T) {
	t.Parallel()

	_, err := loader.Detect([]byte("\n\n\n"))
	if err == nil {
		t.Fatal("expected error for blank-only input, got nil")
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

func TestFormat_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		format loader.Format
		want   string
	}{
		{loader.FormatAuto, "auto"},
		{loader.FormatJSON, "json"},
		{loader.FormatNDJSON, "ndjson"},
	}

	for _, tt := range tests {
		got := tt.format.String()
		if got != tt.want {
			t.Errorf("Format(%d).String() = %q, want %q", tt.format, got, tt.want)
		}
	}
}
