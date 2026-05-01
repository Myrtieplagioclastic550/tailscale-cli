package output

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// mockDevice returns a single fictitious device as map[string]interface{}.
func mockDevice() map[string]interface{} {
	return map[string]interface{}{
		"nodeId":     "nTEST1234CNTRL",
		"hostname":   "mock-server",
		"os":         "linux",
		"authorized": true,
	}
}

// mockDevices returns a slice of two fictitious devices.
func mockDevices() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"nodeId":     "nTEST1234CNTRL",
			"hostname":   "mock-server",
			"os":         "linux",
			"authorized": true,
		},
		{
			"nodeId":     "nTEST5678CNTRL",
			"hostname":   "mock-desktop",
			"os":         "windows",
			"authorized": false,
		},
	}
}

// captureStdout captures everything written to os.Stdout during the execution of fn.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = oldStdout

	captured, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read captured output: %v", err)
	}
	return string(captured)
}

// --- PrintJSON tests ---

func TestPrintJSON_SingleObject(t *testing.T) {
	device := mockDevice()

	out := captureStdout(t, func() {
		err := PrintJSON(device)
		if err != nil {
			t.Fatalf("PrintJSON returned error: %v", err)
		}
	})

	// Verify it is valid JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v\nOutput:\n%s", err, out)
	}

	// Verify indentation (should contain 2-space indent)
	if !strings.Contains(out, "  \"nodeId\"") {
		t.Error("expected indented JSON output with 2-space indent")
	}

	// Verify values
	if parsed["hostname"] != "mock-server" {
		t.Errorf("expected hostname 'mock-server', got %v", parsed["hostname"])
	}
	if parsed["os"] != "linux" {
		t.Errorf("expected os 'linux', got %v", parsed["os"])
	}
}

func TestPrintJSON_Array(t *testing.T) {
	devices := mockDevices()

	out := captureStdout(t, func() {
		err := PrintJSON(devices)
		if err != nil {
			t.Fatalf("PrintJSON returned error: %v", err)
		}
	})

	// Verify it is valid JSON array
	var parsed []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("output is not valid JSON array: %v\nOutput:\n%s", err, out)
	}

	if len(parsed) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(parsed))
	}

	if parsed[0]["hostname"] != "mock-server" {
		t.Errorf("expected first hostname 'mock-server', got %v", parsed[0]["hostname"])
	}
	if parsed[1]["hostname"] != "mock-desktop" {
		t.Errorf("expected second hostname 'mock-desktop', got %v", parsed[1]["hostname"])
	}
}

// --- CSV tests ---

func TestPrintCSV_WithColumns(t *testing.T) {
	devices := mockDevices()
	columns := []string{"hostname", "os"}

	out := captureStdout(t, func() {
		err := Print("csv", devices, columns)
		if err != nil {
			t.Fatalf("Print csv returned error: %v", err)
		}
	})

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines (1 header + 2 rows), got %d:\n%s", len(lines), out)
	}

	// Verify data rows contain expected values
	if !strings.Contains(lines[1], "mock-server") {
		t.Errorf("expected first row to contain 'mock-server', got: %s", lines[1])
	}
	if !strings.Contains(lines[1], "linux") {
		t.Errorf("expected first row to contain 'linux', got: %s", lines[1])
	}
	if !strings.Contains(lines[2], "mock-desktop") {
		t.Errorf("expected second row to contain 'mock-desktop', got: %s", lines[2])
	}
}

func TestPrintCSV_Headers(t *testing.T) {
	device := mockDevice()
	columns := []string{"nodeId", "hostname", "os", "authorized"}

	out := captureStdout(t, func() {
		err := Print("csv", device, columns)
		if err != nil {
			t.Fatalf("Print csv returned error: %v", err)
		}
	})

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) < 1 {
		t.Fatal("expected at least 1 line of output")
	}

	header := lines[0]
	expectedHeader := "nodeId,hostname,os,authorized"
	if header != expectedHeader {
		t.Errorf("expected header %q, got %q", expectedHeader, header)
	}
}

// --- Unsupported format test ---

func TestPrint_UnsupportedFormat(t *testing.T) {
	device := mockDevice()

	err := Print("xml", device, []string{"hostname"})
	if err == nil {
		t.Fatal("expected error for unsupported format 'xml', got nil")
	}
	if !strings.Contains(err.Error(), "unsupported output format") {
		t.Errorf("expected error message to contain 'unsupported output format', got: %v", err)
	}
}

// --- toRows tests ---

func TestToRows_SingleMap(t *testing.T) {
	device := mockDevice()

	rows, err := toRows(device)
	if err != nil {
		t.Fatalf("toRows returned error: %v", err)
	}

	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	if rows[0]["hostname"] != "mock-server" {
		t.Errorf("expected hostname 'mock-server', got %v", rows[0]["hostname"])
	}
}

func TestToRows_SliceOfMaps(t *testing.T) {
	devices := mockDevices()

	rows, err := toRows(devices)
	if err != nil {
		t.Fatalf("toRows returned error: %v", err)
	}

	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	if rows[0]["nodeId"] != "nTEST1234CNTRL" {
		t.Errorf("expected first nodeId 'nTEST1234CNTRL', got %v", rows[0]["nodeId"])
	}
	if rows[1]["nodeId"] != "nTEST5678CNTRL" {
		t.Errorf("expected second nodeId 'nTEST5678CNTRL', got %v", rows[1]["nodeId"])
	}
}

func TestToRows_SliceOfInterfaces(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{
			"nodeId":   "nTEST1234CNTRL",
			"hostname": "mock-server",
		},
		map[string]interface{}{
			"nodeId":   "nTEST5678CNTRL",
			"hostname": "mock-desktop",
		},
	}

	rows, err := toRows(data)
	if err != nil {
		t.Fatalf("toRows returned error: %v", err)
	}

	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	if rows[0]["hostname"] != "mock-server" {
		t.Errorf("expected first hostname 'mock-server', got %v", rows[0]["hostname"])
	}
	if rows[1]["hostname"] != "mock-desktop" {
		t.Errorf("expected second hostname 'mock-desktop', got %v", rows[1]["hostname"])
	}
}

func TestToRows_UnsupportedType(t *testing.T) {
	_, err := toRows("this is a string")
	if err == nil {
		t.Fatal("expected error for unsupported type string, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported data type") {
		t.Errorf("expected error to contain 'unsupported data type', got: %v", err)
	}
}

// --- YAML test ---

func TestPrintYAML(t *testing.T) {
	device := mockDevice()

	out := captureStdout(t, func() {
		err := Print("yaml", device, nil)
		if err != nil {
			t.Fatalf("Print yaml returned error: %v", err)
		}
	})

	// Verify it is valid YAML
	var parsed map[string]interface{}
	if err := yaml.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("output is not valid YAML: %v\nOutput:\n%s", err, out)
	}

	if parsed["hostname"] != "mock-server" {
		t.Errorf("expected hostname 'mock-server', got %v", parsed["hostname"])
	}
	if parsed["os"] != "linux" {
		t.Errorf("expected os 'linux', got %v", parsed["os"])
	}
}
