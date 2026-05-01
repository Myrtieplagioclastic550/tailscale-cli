package service

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
func newTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *ServiceOptions) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	opts := &ServiceOptions{
		GetClient: func() (*api.Client, error) {
			return api.NewClient("test-token", api.WithBaseURL(server.URL)), nil
		},
		GetOutputFormat: func() string { return "json" },
		GetTailnet:      func() string { return "test.example.com" },
	}
	return server, opts
}

// ---------------------------------------------------------------------------
// TestServiceList
// ---------------------------------------------------------------------------

func TestServiceList(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/service_list.json with mock response")

	fixture := loadFixture(t, "service_list.json")
	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/services" {
			t.Errorf("expected path /tailnet/test.example.com/services, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(fixture)
	})
}

// ---------------------------------------------------------------------------
// TestServiceGet
// ---------------------------------------------------------------------------

func TestServiceGet(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/service_get.json with mock response")

	fixture := loadFixture(t, "service_get.json")
	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/services/my-service" {
			t.Errorf("expected path /tailnet/test.example.com/services/my-service, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(fixture)
	})
}

// ---------------------------------------------------------------------------
// TestServiceCreate
// ---------------------------------------------------------------------------

func TestServiceCreate(t *testing.T) {
	t.Skip("TODO: implement service create test with mock server")

	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/services/my-service" {
			t.Errorf("expected path /tailnet/test.example.com/services/my-service, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}

// ---------------------------------------------------------------------------
// TestServiceUpdate
// ---------------------------------------------------------------------------

func TestServiceUpdate(t *testing.T) {
	t.Skip("TODO: implement service update test with mock server")

	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/services/my-service" {
			t.Errorf("expected path /tailnet/test.example.com/services/my-service, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}

// ---------------------------------------------------------------------------
// TestServiceDelete
// ---------------------------------------------------------------------------

func TestServiceDelete(t *testing.T) {
	t.Skip("TODO: implement service delete test with mock server")

	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/services/my-service" {
			t.Errorf("expected path /tailnet/test.example.com/services/my-service, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
}

// ---------------------------------------------------------------------------
// TestServiceHosts
// ---------------------------------------------------------------------------

func TestServiceHosts(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/service_hosts.json with mock response")

	fixture := loadFixture(t, "service_hosts.json")
	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/services/my-service/devices" {
			t.Errorf("expected path /tailnet/test.example.com/services/my-service/devices, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(fixture)
	})
}

// ---------------------------------------------------------------------------
// TestServiceApprove
// ---------------------------------------------------------------------------

func TestServiceApprove(t *testing.T) {
	t.Skip("TODO: implement service approve test with mock server")

	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/services/my-service/device/nTEST1234CNTRL/approved" {
			t.Errorf("expected path /tailnet/test.example.com/services/my-service/device/nTEST1234CNTRL/approved, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
}
