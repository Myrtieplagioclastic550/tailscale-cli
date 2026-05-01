package posture

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
func newTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *PostureOptions) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	opts := &PostureOptions{
		GetClient: func() (*api.Client, error) {
			return api.NewClient("test-token", api.WithBaseURL(server.URL)), nil
		},
		GetOutputFormat: func() string { return "json" },
		GetTailnet:      func() string { return "test.example.com" },
	}
	return server, opts
}

// ---------------------------------------------------------------------------
// TestPostureList
// ---------------------------------------------------------------------------

func TestPostureList(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/posture_list.json with mock response")

	fixture := loadFixture(t, "posture_list.json")
	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/posture/integrations" {
			t.Errorf("expected path /tailnet/test.example.com/posture/integrations, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(fixture)
	})
}

// ---------------------------------------------------------------------------
// TestPostureCreate
// ---------------------------------------------------------------------------

func TestPostureCreate(t *testing.T) {
	t.Skip("TODO: implement posture create test with mock server")

	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/posture/integrations" {
			t.Errorf("expected path /tailnet/test.example.com/posture/integrations, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}

// ---------------------------------------------------------------------------
// TestPostureGet
// ---------------------------------------------------------------------------

func TestPostureGet(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/posture_get.json with mock response")

	fixture := loadFixture(t, "posture_get.json")
	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/posture/integrations/12345" {
			t.Errorf("expected path /posture/integrations/12345, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(fixture)
	})
}

// ---------------------------------------------------------------------------
// TestPostureUpdate
// ---------------------------------------------------------------------------

func TestPostureUpdate(t *testing.T) {
	t.Skip("TODO: implement posture update test with mock server")

	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/posture/integrations/12345" {
			t.Errorf("expected path /posture/integrations/12345, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}

// ---------------------------------------------------------------------------
// TestPostureDelete
// ---------------------------------------------------------------------------

func TestPostureDelete(t *testing.T) {
	t.Skip("TODO: implement posture delete test with mock server")

	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/posture/integrations/12345" {
			t.Errorf("expected path /posture/integrations/12345, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
}
