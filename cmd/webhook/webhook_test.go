package webhook

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

// makeOpts returns a WebhookOptions wired to the given httptest server.
func makeOpts(srv *httptest.Server) WebhookOptions {
	return WebhookOptions{
		GetClient: func() (*api.Client, error) {
			return api.NewClient("tskey-test-token", api.WithBaseURL(srv.URL)), nil
		},
		GetOutputFormat: func() string { return "json" },
		GetTailnet:      func() string { return "test-tailnet" },
	}
}

// execCmd builds and runs the webhook cobra command with the provided arguments.
func execCmd(opts WebhookOptions, args []string) error {
	root := &cobra.Command{Use: "tailscale"}
	root.AddCommand(NewCmdWebhook(opts))
	root.SetArgs(args)
	root.SilenceUsage = true
	root.SilenceErrors = true
	return root.Execute()
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestWebhookList(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/webhook_list.json with mock response")

	fixture := loadFixture(t, "webhook_list.json")

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
	err := execCmd(opts, []string{"webhook", "list"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if gotMethod != http.MethodGet {
		t.Errorf("expected method GET, got %s", gotMethod)
	}
	expectedPath := "/tailnet/test-tailnet/webhooks"
	if gotPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, gotPath)
	}
}

func TestWebhookCreate(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/webhook_list.json with mock response")

	fixture := loadFixture(t, "webhook_list.json")

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
	err := execCmd(opts, []string{"webhook", "create", "--url", "https://example.com/hook", "--events", "nodeCreated,nodeDeleted"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if gotMethod != http.MethodPost {
		t.Errorf("expected method POST, got %s", gotMethod)
	}
	expectedPath := "/tailnet/test-tailnet/webhooks"
	if gotPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, gotPath)
	}
}

func TestWebhookGet(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/webhook_list.json with mock response")

	fixture := loadFixture(t, "webhook_list.json")

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
	err := execCmd(opts, []string{"webhook", "get", "wh-123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if gotMethod != http.MethodGet {
		t.Errorf("expected method GET, got %s", gotMethod)
	}
	expectedPath := "/webhooks/wh-123"
	if gotPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, gotPath)
	}
}

func TestWebhookDelete(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/webhook_list.json with mock response")

	var mu sync.Mutex
	var gotPath string
	var gotMethod string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		gotPath = r.URL.Path
		gotMethod = r.Method
		mu.Unlock()

		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	opts := makeOpts(srv)
	err := execCmd(opts, []string{"webhook", "delete", "wh-123", "--confirm"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if gotMethod != http.MethodDelete {
		t.Errorf("expected method DELETE, got %s", gotMethod)
	}
	expectedPath := "/webhooks/wh-123"
	if gotPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, gotPath)
	}
}

func TestWebhookTest(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/webhook_list.json with mock response")

	var mu sync.Mutex
	var gotPath string
	var gotMethod string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		gotPath = r.URL.Path
		gotMethod = r.Method
		mu.Unlock()

		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	opts := makeOpts(srv)
	err := execCmd(opts, []string{"webhook", "test", "wh-123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if gotMethod != http.MethodPost {
		t.Errorf("expected method POST, got %s", gotMethod)
	}
	expectedPath := "/webhooks/wh-123/test"
	if gotPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, gotPath)
	}
}
