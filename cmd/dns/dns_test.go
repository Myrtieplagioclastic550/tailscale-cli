package dns

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

// makeOpts returns a DnsOptions wired to the given httptest server.
func makeOpts(srv *httptest.Server) DnsOptions {
	return DnsOptions{
		GetClient: func() (*api.Client, error) {
			return api.NewClient("tskey-test-token", api.WithBaseURL(srv.URL)), nil
		},
		GetOutputFormat: func() string { return "json" },
		GetTailnet:      func() string { return "test-tailnet" },
	}
}

// execCmd builds and runs the dns cobra command with the provided arguments.
func execCmd(opts DnsOptions, args []string) error {
	root := &cobra.Command{Use: "tailscale"}
	root.AddCommand(NewCmdDns(opts))
	root.SetArgs(args)
	root.SilenceUsage = true
	root.SilenceErrors = true
	return root.Execute()
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestDnsNameserversList(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/dns_nameservers.json with mock response")

	fixture := loadFixture(t, "dns_nameservers.json")

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
	err := execCmd(opts, []string{"dns", "nameservers", "list"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if gotMethod != http.MethodGet {
		t.Errorf("expected method GET, got %s", gotMethod)
	}
	expectedPath := "/tailnet/test-tailnet/dns/nameservers"
	if gotPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, gotPath)
	}
}

func TestDnsNameserversSet(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/dns_nameservers.json with mock response")

	fixture := loadFixture(t, "dns_nameservers.json")

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
	err := execCmd(opts, []string{"dns", "nameservers", "set", "--nameservers", "8.8.8.8,1.1.1.1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if gotMethod != http.MethodPost {
		t.Errorf("expected method POST, got %s", gotMethod)
	}
	expectedPath := "/tailnet/test-tailnet/dns/nameservers"
	if gotPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, gotPath)
	}
}

func TestDnsPreferencesGet(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/dns_preferences.json with mock response")

	fixture := loadFixture(t, "dns_preferences.json")

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
	err := execCmd(opts, []string{"dns", "preferences", "get"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if gotMethod != http.MethodGet {
		t.Errorf("expected method GET, got %s", gotMethod)
	}
	expectedPath := "/tailnet/test-tailnet/dns/preferences"
	if gotPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, gotPath)
	}
}

func TestDnsPreferencesSet(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/dns_preferences.json with mock response")

	fixture := loadFixture(t, "dns_preferences.json")

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
	err := execCmd(opts, []string{"dns", "preferences", "set", "--magic-dns"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if gotMethod != http.MethodPost {
		t.Errorf("expected method POST, got %s", gotMethod)
	}
	expectedPath := "/tailnet/test-tailnet/dns/preferences"
	if gotPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, gotPath)
	}
}

func TestDnsSearchpathsList(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/dns_searchpaths.json with mock response")

	fixture := loadFixture(t, "dns_searchpaths.json")

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
	err := execCmd(opts, []string{"dns", "searchpaths", "list"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if gotMethod != http.MethodGet {
		t.Errorf("expected method GET, got %s", gotMethod)
	}
	expectedPath := "/tailnet/test-tailnet/dns/searchpaths"
	if gotPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, gotPath)
	}
}

func TestDnsSplitGet(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/dns_split.json with mock response")

	fixture := loadFixture(t, "dns_split.json")

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
	err := execCmd(opts, []string{"dns", "split", "get"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if gotMethod != http.MethodGet {
		t.Errorf("expected method GET, got %s", gotMethod)
	}
	expectedPath := "/tailnet/test-tailnet/dns/split-dns"
	if gotPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, gotPath)
	}
}

func TestDnsConfigGet(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/dns_config.json with mock response")

	fixture := loadFixture(t, "dns_config.json")

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
	err := execCmd(opts, []string{"dns", "config", "get"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if gotMethod != http.MethodGet {
		t.Errorf("expected method GET, got %s", gotMethod)
	}
	expectedPath := "/tailnet/test-tailnet/dns/configuration"
	if gotPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, gotPath)
	}
}
