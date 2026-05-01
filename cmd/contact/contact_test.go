package contact

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
func newTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *ContactOptions) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	opts := &ContactOptions{
		GetClient: func() (*api.Client, error) {
			return api.NewClient("test-token", api.WithBaseURL(server.URL)), nil
		},
		GetOutputFormat: func() string { return "json" },
		GetTailnet:      func() string { return "test.example.com" },
	}
	return server, opts
}

// ---------------------------------------------------------------------------
// TestContactGet
// ---------------------------------------------------------------------------

func TestContactGet(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/contact_get.json with mock response")

	fixture := loadFixture(t, "contact_get.json")
	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/contacts" {
			t.Errorf("expected path /tailnet/test.example.com/contacts, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(fixture)
	})
}

// ---------------------------------------------------------------------------
// TestContactUpdate
// ---------------------------------------------------------------------------

func TestContactUpdate(t *testing.T) {
	t.Skip("TODO: implement contact update test with mock server")

	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/contacts/account" {
			t.Errorf("expected path /tailnet/test.example.com/contacts/account, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}

// ---------------------------------------------------------------------------
// TestContactResendVerification
// ---------------------------------------------------------------------------

func TestContactResendVerification(t *testing.T) {
	t.Skip("TODO: implement contact resend verification test with mock server")

	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/contacts/account/resend-verification-email" {
			t.Errorf("expected path /tailnet/test.example.com/contacts/account/resend-verification-email, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
}
