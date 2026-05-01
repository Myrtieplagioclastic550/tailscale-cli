package settings

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/dimer47/tailscale-cli/internal/api"
	"github.com/spf13/cobra"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func fixturesDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..", "testdata", "fixtures")
}

func loadFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(fixturesDir(), name))
	if err != nil {
		t.Fatalf("failed to load fixture %s: %v", name, err)
	}
	return data
}

// makeOpts returns a SettingsOptions wired to the given httptest server.
func makeOpts(srv *httptest.Server) SettingsOptions {
	return SettingsOptions{
		GetClient: func() (*api.Client, error) {
			return api.NewClient("tskey-test-token", api.WithBaseURL(srv.URL)), nil
		},
		GetOutputFormat: func() string { return "json" },
		GetTailnet:      func() string { return "test-tailnet" },
	}
}

// execCmd builds and runs the settings cobra command with the provided arguments.
func execCmd(opts SettingsOptions, args []string) error {
	root := &cobra.Command{Use: "tailscale"}
	root.AddCommand(NewCmdSettings(opts))
	root.SetArgs(args)
	root.SilenceUsage = true
	root.SilenceErrors = true
	return root.Execute()
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestSettingsGet(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/settings_get.json with mock response")

	fixture := loadFixture(t, "settings_get.json")

	var mu sync.Mutex
	var gotPath string
	var gotMethod string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		gotPath = r.URL.Path
		gotMethod = r.Method
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(fixture)
	}))
	defer srv.Close()

	opts := makeOpts(srv)
	err := execCmd(opts, []string{"settings", "get"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if gotMethod != http.MethodGet {
		t.Errorf("expected method GET, got %s", gotMethod)
	}
	expectedPath := "/tailnet/test-tailnet/settings"
	if gotPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, gotPath)
	}
}

func TestSettingsUpdate(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/settings_get.json with mock response")

	fixture := loadFixture(t, "settings_get.json")

	var mu sync.Mutex
	var gotPath string
	var gotMethod string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		gotPath = r.URL.Path
		gotMethod = r.Method
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(fixture)
	}))
	defer srv.Close()

	opts := makeOpts(srv)
	err := execCmd(opts, []string{"settings", "update", "--https"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if gotMethod != http.MethodPatch {
		t.Errorf("expected method PATCH, got %s", gotMethod)
	}
	expectedPath := "/tailnet/test-tailnet/settings"
	if gotPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, gotPath)
	}
}
