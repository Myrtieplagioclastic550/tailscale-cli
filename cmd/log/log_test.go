package log

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/dimer47/tailscale-cli/internal/api"
)

// fixturesDir returns the path to the testdata/fixtures directory.
func fixturesDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..", "testdata", "fixtures")
}

// loadFixture reads a fixture file. Returns empty bytes if file is empty or contains only "{}".
func loadFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(fixturesDir(), name))
	if err != nil {
		t.Fatalf("failed to load fixture %s: %v", name, err)
	}
	return data
}

// newTestServer creates a mock HTTP server that serves fixture responses.
func newTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *LogOptions) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	opts := &LogOptions{
		GetClient: func() (*api.Client, error) {
			return api.NewClient("test-token", api.WithBaseURL(server.URL)), nil
		},
		GetOutputFormat: func() string { return "json" },
		GetTailnet:      func() string { return "test.example.com" },
	}
	return server, opts
}

// ---------------------------------------------------------------------------
// TestLogAuditList
// ---------------------------------------------------------------------------

func TestLogAuditList(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/log_audit.json with mock response")

	fixture := loadFixture(t, "log_audit.json")
	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/logging/configuration" {
			t.Errorf("expected path /tailnet/test.example.com/logging/configuration, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(fixture)
	})
}

// ---------------------------------------------------------------------------
// TestLogNetworkList
// ---------------------------------------------------------------------------

func TestLogNetworkList(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/log_network.json with mock response")

	fixture := loadFixture(t, "log_network.json")
	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/logging/network" {
			t.Errorf("expected path /tailnet/test.example.com/logging/network, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(fixture)
	})
}

// ---------------------------------------------------------------------------
// TestLogStreamStatus
// ---------------------------------------------------------------------------

func TestLogStreamStatus(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/log_stream_status.json with mock response")

	fixture := loadFixture(t, "log_stream_status.json")
	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/logging/configuration/stream/status" {
			t.Errorf("expected path /tailnet/test.example.com/logging/configuration/stream/status, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(fixture)
	})
}

// ---------------------------------------------------------------------------
// TestLogStreamGet
// ---------------------------------------------------------------------------

func TestLogStreamGet(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/log_stream_config.json with mock response")

	fixture := loadFixture(t, "log_stream_config.json")
	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/logging/configuration/stream" {
			t.Errorf("expected path /tailnet/test.example.com/logging/configuration/stream, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(fixture)
	})
}

// ---------------------------------------------------------------------------
// TestLogStreamSet
// ---------------------------------------------------------------------------

func TestLogStreamSet(t *testing.T) {
	t.Skip("TODO: implement log stream set test with mock server")

	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/logging/configuration/stream" {
			t.Errorf("expected path /tailnet/test.example.com/logging/configuration/stream, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
}

// ---------------------------------------------------------------------------
// TestLogStreamDisable
// ---------------------------------------------------------------------------

func TestLogStreamDisable(t *testing.T) {
	t.Skip("TODO: implement log stream disable test with mock server")

	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/logging/configuration/stream" {
			t.Errorf("expected path /tailnet/test.example.com/logging/configuration/stream, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
}

// ---------------------------------------------------------------------------
// TestLogAwsIdCreate
// ---------------------------------------------------------------------------

func TestLogAwsIdCreate(t *testing.T) {
	t.Skip("TODO: implement log aws-id create test with mock server")

	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/aws-external-id" {
			t.Errorf("expected path /tailnet/test.example.com/aws-external-id, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}

// ---------------------------------------------------------------------------
// TestLogAwsIdValidate
// ---------------------------------------------------------------------------

func TestLogAwsIdValidate(t *testing.T) {
	t.Skip("TODO: implement log aws-id validate test with mock server")

	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/aws-external-id/ext-12345/validate-aws-trust-policy" {
			t.Errorf("expected path /tailnet/test.example.com/aws-external-id/ext-12345/validate-aws-trust-policy, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}
