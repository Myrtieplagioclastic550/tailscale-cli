package acl

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

// loadFixture reads a fixture file.
func loadFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(fixturesDir(), name))
	if err != nil {
		t.Fatalf("failed to load fixture %s: %v", name, err)
	}
	return data
}

// newTestServer creates a mock HTTP server that serves fixture responses.
func newTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, AclOptions) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	opts := AclOptions{
		GetClient: func() (*api.Client, error) {
			return api.NewClient("test-token", api.WithBaseURL(server.URL)), nil
		},
		GetOutputFormat: func() string { return "json" },
		GetTailnet:      func() string { return "test.example.com" },
	}
	return server, opts
}

// ---------------------------------------------------------------------------
// TestAclGet
// ---------------------------------------------------------------------------

func TestAclGet(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/acl_get.json with mock response")

	fixture := loadFixture(t, "acl_get.json")
	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/acl" {
			t.Errorf("expected path /tailnet/test.example.com/acl, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(fixture)
	})
}

// ---------------------------------------------------------------------------
// TestAclSet
// ---------------------------------------------------------------------------

func TestAclSet(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/acl_get.json with mock response for set operation")

	fixture := loadFixture(t, "acl_get.json")
	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/acl" {
			t.Errorf("expected path /tailnet/test.example.com/acl, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(fixture)
	})
}

// ---------------------------------------------------------------------------
// TestAclValidate
// ---------------------------------------------------------------------------

func TestAclValidate(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/acl_validate.json with mock response")

	fixture := loadFixture(t, "acl_validate.json")
	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/acl/validate" {
			t.Errorf("expected path /tailnet/test.example.com/acl/validate, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(fixture)
	})
}

// ---------------------------------------------------------------------------
// TestAclPreview
// ---------------------------------------------------------------------------

func TestAclPreview(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/acl_get.json with mock response for preview operation")

	fixture := loadFixture(t, "acl_get.json")
	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		expectedPrefix := "/tailnet/test.example.com/acl/preview"
		if r.URL.Path != expectedPrefix {
			t.Errorf("expected path %s, got %s", expectedPrefix, r.URL.Path)
		}
		// Verify query parameters are present
		if r.URL.Query().Get("type") == "" {
			t.Error("expected 'type' query parameter to be present")
		}
		if r.URL.Query().Get("previewFor") == "" {
			t.Error("expected 'previewFor' query parameter to be present")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(fixture)
	})
}
