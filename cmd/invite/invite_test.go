package invite

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
func newTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *InviteOptions) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	opts := &InviteOptions{
		GetClient: func() (*api.Client, error) {
			return api.NewClient("test-token", api.WithBaseURL(server.URL)), nil
		},
		GetOutputFormat: func() string { return "json" },
		GetTailnet:      func() string { return "test.example.com" },
	}
	return server, opts
}

// ---------------------------------------------------------------------------
// TestInviteUserList
// ---------------------------------------------------------------------------

func TestInviteUserList(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/invite_user_list.json with mock response")

	fixture := loadFixture(t, "invite_user_list.json")
	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/user-invites" {
			t.Errorf("expected path /tailnet/test.example.com/user-invites, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(fixture)
	})
}

// ---------------------------------------------------------------------------
// TestInviteUserCreate
// ---------------------------------------------------------------------------

func TestInviteUserCreate(t *testing.T) {
	t.Skip("TODO: implement invite user create test with mock server")

	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/test.example.com/user-invites" {
			t.Errorf("expected path /tailnet/test.example.com/user-invites, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}

// ---------------------------------------------------------------------------
// TestInviteUserGet
// ---------------------------------------------------------------------------

func TestInviteUserGet(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/invite_user_get.json with mock response")

	fixture := loadFixture(t, "invite_user_get.json")
	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/user-invites/12345" {
			t.Errorf("expected path /user-invites/12345, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(fixture)
	})
}

// ---------------------------------------------------------------------------
// TestInviteUserDelete
// ---------------------------------------------------------------------------

func TestInviteUserDelete(t *testing.T) {
	t.Skip("TODO: implement invite user delete test with mock server")

	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/user-invites/12345" {
			t.Errorf("expected path /user-invites/12345, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
}

// ---------------------------------------------------------------------------
// TestInviteUserResend
// ---------------------------------------------------------------------------

func TestInviteUserResend(t *testing.T) {
	t.Skip("TODO: implement invite user resend test with mock server")

	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/user-invites/12345/resend" {
			t.Errorf("expected path /user-invites/12345/resend, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
}

// ---------------------------------------------------------------------------
// TestInviteDeviceList
// ---------------------------------------------------------------------------

func TestInviteDeviceList(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/invite_device_list.json with mock response")

	fixture := loadFixture(t, "invite_device_list.json")
	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/device/nTEST1234CNTRL/device-invites" {
			t.Errorf("expected path /device/nTEST1234CNTRL/device-invites, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(fixture)
	})
}

// ---------------------------------------------------------------------------
// TestInviteDeviceCreate
// ---------------------------------------------------------------------------

func TestInviteDeviceCreate(t *testing.T) {
	t.Skip("TODO: implement invite device create test with mock server")

	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/device/nTEST1234CNTRL/device-invites" {
			t.Errorf("expected path /device/nTEST1234CNTRL/device-invites, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}

// ---------------------------------------------------------------------------
// TestInviteDeviceGet
// ---------------------------------------------------------------------------

func TestInviteDeviceGet(t *testing.T) {
	t.Skip("TODO: fill testdata/fixtures/invite_device_get.json with mock response")

	fixture := loadFixture(t, "invite_device_get.json")
	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/device-invites/12345" {
			t.Errorf("expected path /device-invites/12345, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(fixture)
	})
}

// ---------------------------------------------------------------------------
// TestInviteDeviceDelete
// ---------------------------------------------------------------------------

func TestInviteDeviceDelete(t *testing.T) {
	t.Skip("TODO: implement invite device delete test with mock server")

	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/device-invites/12345" {
			t.Errorf("expected path /device-invites/12345, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
}

// ---------------------------------------------------------------------------
// TestInviteDeviceResend
// ---------------------------------------------------------------------------

func TestInviteDeviceResend(t *testing.T) {
	t.Skip("TODO: implement invite device resend test with mock server")

	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/device-invites/12345/resend" {
			t.Errorf("expected path /device-invites/12345/resend, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
}

// ---------------------------------------------------------------------------
// TestInviteDeviceAccept
// ---------------------------------------------------------------------------

func TestInviteDeviceAccept(t *testing.T) {
	t.Skip("TODO: implement invite device accept test with mock server")

	_, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/device-invites/-/accept" {
			t.Errorf("expected path /device-invites/-/accept, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
}
