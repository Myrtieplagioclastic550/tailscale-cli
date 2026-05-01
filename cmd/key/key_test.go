package key

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
func newTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, KeyOptions) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	opts := KeyOptions{
		GetClient: func() (*api.Client, error) {
			return api.NewClient("test-token", api.WithBaseURL(server.URL)), nil
		},
		GetOutputFormat: func() string { return "json" },
		GetTailnet:      func() string { return "test.example.com" },
	}
	return server, opts
}

// ---------------------------------------------------------------------------
// TestKeyList
// ---------------------------------------------------------------------------

func TestKeyList(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/key_list.json with mock response")

	fixture := loadFixture(t, "key_list.json")
	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/keys" {
			t.Errorf("expected path /tailnet/test.example.com/keys, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(fixture)
	})
}

// ---------------------------------------------------------------------------
// TestKeyCreate
// ---------------------------------------------------------------------------

func TestKeyCreate(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/key_create.json with mock response")

	fixture := loadFixture(t, "key_create.json")
	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/keys" {
			t.Errorf("expected path /tailnet/test.example.com/keys, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(fixture)
	})
}

// ---------------------------------------------------------------------------
// TestKeyGet
// ---------------------------------------------------------------------------

func TestKeyGet(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/key_get.json with mock response")

	fixture := loadFixture(t, "key_get.json")
	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/keys/kTEST1234CNTRL" {
			t.Errorf("expected path /tailnet/test.example.com/keys/kTEST1234CNTRL, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(fixture)
	})
}

// ---------------------------------------------------------------------------
// TestKeyDelete
// ---------------------------------------------------------------------------

func TestKeyDelete(t *testing.T) {
	t.Skip("TODO: implement key delete test with mock server")

	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/keys/kTEST1234CNTRL" {
			t.Errorf("expected path /tailnet/test.example.com/keys/kTEST1234CNTRL, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
}
