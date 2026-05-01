package user

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

// makeOpts returns a UserOptions wired to the given httptest server.
func makeOpts(srv *httptest.Server) UserOptions {
	return UserOptions{
		GetClient: func() (*api.Client, error) {
			return api.NewClient("tskey-test-token", api.WithBaseURL(srv.URL)), nil
		},
		GetOutputFormat: func() string { return "json" },
		GetTailnet:      func() string { return "test-tailnet" },
	}
}

// execCmd builds and runs the user cobra command with the provided arguments.
func execCmd(opts UserOptions, args []string) error {
	root := &cobra.Command{Use: "tailscale"}
	root.AddCommand(NewCmdUser(opts))
	root.SetArgs(args)
	root.SilenceUsage = true
	root.SilenceErrors = true
	return root.Execute()
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestUserList(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/user_list.json with mock response")

	fixture := loadFixture(t, "user_list.json")

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
	err := execCmd(opts, []string{"user", "list"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if gotMethod != http.MethodGet {
		t.Errorf("expected method GET, got %s", gotMethod)
	}
	expectedPath := "/tailnet/test-tailnet/users"
	if gotPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, gotPath)
	}
}

func TestUserGet(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/user_get.json with mock response")

	fixture := loadFixture(t, "user_get.json")

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
	err := execCmd(opts, []string{"user", "get", "user-123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if gotMethod != http.MethodGet {
		t.Errorf("expected method GET, got %s", gotMethod)
	}
	expectedPath := "/users/user-123"
	if gotPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, gotPath)
	}
}

func TestUserSetRole(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/user_get.json with mock response")

	fixture := loadFixture(t, "user_get.json")

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
	err := execCmd(opts, []string{"user", "set-role", "user-123", "admin"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if gotMethod != http.MethodPost {
		t.Errorf("expected method POST, got %s", gotMethod)
	}
	expectedPath := "/users/user-123/role"
	if gotPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, gotPath)
	}
}

func TestUserApprove(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/user_get.json with mock response")

	fixture := loadFixture(t, "user_get.json")

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
	err := execCmd(opts, []string{"user", "approve", "user-123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if gotMethod != http.MethodPost {
		t.Errorf("expected method POST, got %s", gotMethod)
	}
	expectedPath := "/users/user-123/approve"
	if gotPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, gotPath)
	}
}

func TestUserSuspend(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/user_get.json with mock response")

	fixture := loadFixture(t, "user_get.json")

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
	err := execCmd(opts, []string{"user", "suspend", "user-123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if gotMethod != http.MethodPost {
		t.Errorf("expected method POST, got %s", gotMethod)
	}
	expectedPath := "/users/user-123/suspend"
	if gotPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, gotPath)
	}
}

func TestUserRestore(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/user_get.json with mock response")

	fixture := loadFixture(t, "user_get.json")

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
	err := execCmd(opts, []string{"user", "restore", "user-123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if gotMethod != http.MethodPost {
		t.Errorf("expected method POST, got %s", gotMethod)
	}
	expectedPath := "/users/user-123/restore"
	if gotPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, gotPath)
	}
}
