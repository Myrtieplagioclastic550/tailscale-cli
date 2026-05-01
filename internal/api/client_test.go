package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Structures fictives conformes au schema API Tailscale v2
// ---------------------------------------------------------------------------

type fakeDevice struct {
	NodeID     string   `json:"nodeId"`
	Hostname   string   `json:"hostname"`
	OS         string   `json:"os"`
	Addresses  []string `json:"addresses"`
	Authorized bool     `json:"authorized"`
	User       string   `json:"user"`
}

type fakeDevicesResponse struct {
	Devices []fakeDevice `json:"devices"`
}

type fakeKey struct {
	ID      string `json:"id"`
	Key     string `json:"key"`
	KeyType string `json:"keyType"`
	Created string `json:"created"`
}

type fakeAPIError struct {
	Message string `json:"message"`
}

// ---------------------------------------------------------------------------
// 1. TestNewClient
// ---------------------------------------------------------------------------

func TestNewClient(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		c := NewClient("tskey-test-token")
		if c.baseURL != defaultBaseURL {
			t.Errorf("expected base URL %q, got %q", defaultBaseURL, c.baseURL)
		}
		if c.apiToken != "tskey-test-token" {
			t.Errorf("expected token %q, got %q", "tskey-test-token", c.apiToken)
		}
		if c.debug {
			t.Error("expected debug to be false by default")
		}
	})

	t.Run("with custom base URL", func(t *testing.T) {
		c := NewClient("tskey-test-token", WithBaseURL("https://custom.api.local"))
		if c.baseURL != "https://custom.api.local" {
			t.Errorf("expected base URL %q, got %q", "https://custom.api.local", c.baseURL)
		}
	})

	t.Run("with debug enabled", func(t *testing.T) {
		c := NewClient("tskey-test-token", WithDebug(true))
		if !c.debug {
			t.Error("expected debug to be true")
		}
	})

	t.Run("with custom HTTP client", func(t *testing.T) {
		custom := &http.Client{}
		c := NewClient("tskey-test-token", WithHTTPClient(custom))
		if c.httpClient != custom {
			t.Error("expected custom HTTP client to be set")
		}
	})

	t.Run("with multiple options", func(t *testing.T) {
		c := NewClient("tskey-test-token",
			WithBaseURL("https://custom.api.local"),
			WithDebug(true),
		)
		if c.baseURL != "https://custom.api.local" {
			t.Errorf("expected base URL %q, got %q", "https://custom.api.local", c.baseURL)
		}
		if !c.debug {
			t.Error("expected debug to be true")
		}
	})
}

// ---------------------------------------------------------------------------
// 2. TestGet_Success
// ---------------------------------------------------------------------------

func TestGet_Success(t *testing.T) {
	expected := fakeDevicesResponse{
		Devices: []fakeDevice{
			{
				NodeID:     "nABCD1234CNTRL",
				Hostname:   "test-server-01",
				OS:         "linux",
				Addresses:  []string{"100.64.0.1"},
				Authorized: true,
				User:       "testuser@example.com",
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/-/devices" {
			t.Errorf("expected path /tailnet/-/devices, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expected)
	}))
	defer srv.Close()

	c := NewClient("tskey-test-token", WithBaseURL(srv.URL))
	body, err := c.Get("/tailnet/-/devices")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got fakeDevicesResponse
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(got.Devices) != 1 {
		t.Fatalf("expected 1 device, got %d", len(got.Devices))
	}
	d := got.Devices[0]
	if d.NodeID != "nABCD1234CNTRL" {
		t.Errorf("expected nodeId %q, got %q", "nABCD1234CNTRL", d.NodeID)
	}
	if d.Hostname != "test-server-01" {
		t.Errorf("expected hostname %q, got %q", "test-server-01", d.Hostname)
	}
	if d.OS != "linux" {
		t.Errorf("expected os %q, got %q", "linux", d.OS)
	}
	if !d.Authorized {
		t.Error("expected device to be authorized")
	}
	if d.User != "testuser@example.com" {
		t.Errorf("expected user %q, got %q", "testuser@example.com", d.User)
	}
}

// ---------------------------------------------------------------------------
// 3. TestGet_APIError
// ---------------------------------------------------------------------------

func TestGet_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(fakeAPIError{Message: "device not found"})
	}))
	defer srv.Close()

	c := NewClient("tskey-test-token", WithBaseURL(srv.URL))
	_, err := c.Get("/device/nINVALID")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	if !strings.Contains(err.Error(), "404") {
		t.Errorf("expected error to contain status code 404, got: %v", err)
	}
	if !strings.Contains(err.Error(), "device not found") {
		t.Errorf("expected error to contain message %q, got: %v", "device not found", err)
	}
}

// ---------------------------------------------------------------------------
// 4. TestGet_ServerError
// ---------------------------------------------------------------------------

func TestGet_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}))
	defer srv.Close()

	c := NewClient("tskey-test-token", WithBaseURL(srv.URL))
	_, err := c.Get("/tailnet/-/devices")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	if !strings.Contains(err.Error(), "500") {
		t.Errorf("expected error to contain status code 500, got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// 5. TestPost_WithBody
// ---------------------------------------------------------------------------

func TestPost_WithBody(t *testing.T) {
	expectedKey := fakeKey{
		ID:      "kTEST1234CNTRL",
		Key:     "tskey-auth-test-secret",
		KeyType: "auth",
		Created: "2026-01-01T00:00:00Z",
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/tailnet/-/keys" {
			t.Errorf("expected path /tailnet/-/keys, got %s", r.URL.Path)
		}

		// Verify request body was sent
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		if len(reqBody) == 0 {
			t.Error("expected non-empty request body")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedKey)
	}))
	defer srv.Close()

	c := NewClient("tskey-test-token", WithBaseURL(srv.URL))
	reqBody := strings.NewReader(`{"capabilities":{"devices":{"create":{"reusable":false,"ephemeral":true}}}}`)
	body, err := c.Post("/tailnet/-/keys", reqBody)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got fakeKey
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if got.ID != "kTEST1234CNTRL" {
		t.Errorf("expected key id %q, got %q", "kTEST1234CNTRL", got.ID)
	}
	if got.Key != "tskey-auth-test-secret" {
		t.Errorf("expected key %q, got %q", "tskey-auth-test-secret", got.Key)
	}
	if got.KeyType != "auth" {
		t.Errorf("expected keyType %q, got %q", "auth", got.KeyType)
	}
	if got.Created != "2026-01-01T00:00:00Z" {
		t.Errorf("expected created %q, got %q", "2026-01-01T00:00:00Z", got.Created)
	}
}

// ---------------------------------------------------------------------------
// 6. TestDoJSON
// ---------------------------------------------------------------------------

func TestDoJSON(t *testing.T) {
	type keyRequest struct {
		Capabilities map[string]interface{} `json:"capabilities"`
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		// Verify the body was JSON-encoded
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}

		var decoded keyRequest
		if err := json.Unmarshal(reqBody, &decoded); err != nil {
			t.Fatalf("request body is not valid JSON: %v", err)
		}
		if decoded.Capabilities == nil {
			t.Error("expected capabilities to be present in request body")
		}

		ct := r.Header.Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("expected Content-Type %q, got %q", "application/json", ct)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fakeKey{
			ID:      "kTEST1234CNTRL",
			Key:     "tskey-auth-test-secret",
			KeyType: "auth",
			Created: "2026-01-01T00:00:00Z",
		})
	}))
	defer srv.Close()

	c := NewClient("tskey-test-token", WithBaseURL(srv.URL))

	payload := keyRequest{
		Capabilities: map[string]interface{}{
			"devices": map[string]interface{}{
				"create": map[string]interface{}{
					"reusable":  false,
					"ephemeral": true,
				},
			},
		},
	}

	body, err := c.DoJSON("POST", "/tailnet/-/keys", payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got fakeKey
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if got.ID != "kTEST1234CNTRL" {
		t.Errorf("expected key id %q, got %q", "kTEST1234CNTRL", got.ID)
	}
}

// ---------------------------------------------------------------------------
// 7. TestAuthHeader
// ---------------------------------------------------------------------------

func TestAuthHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		expected := "Bearer tskey-auth-secret-12345"
		if auth != expected {
			t.Errorf("expected Authorization header %q, got %q", expected, auth)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}))
	defer srv.Close()

	c := NewClient("tskey-auth-secret-12345", WithBaseURL(srv.URL))
	_, err := c.Get("/tailnet/-/devices")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// 8. TestContentTypeHeader
// ---------------------------------------------------------------------------

func TestContentTypeHeader(t *testing.T) {
	t.Run("with body sets Content-Type", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ct := r.Header.Get("Content-Type")
			if ct != "application/json" {
				t.Errorf("expected Content-Type %q, got %q", "application/json", ct)
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{}"))
		}))
		defer srv.Close()

		c := NewClient("tskey-test-token", WithBaseURL(srv.URL))
		_, err := c.Post("/tailnet/-/keys", strings.NewReader(`{"test": true}`))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("without body does not set Content-Type", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ct := r.Header.Get("Content-Type")
			if ct == "application/json" {
				t.Error("expected Content-Type to NOT be set for GET without body")
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{}"))
		}))
		defer srv.Close()

		c := NewClient("tskey-test-token", WithBaseURL(srv.URL))
		_, err := c.Get("/tailnet/-/devices")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

// ---------------------------------------------------------------------------
// 9. TestDelete_Success
// ---------------------------------------------------------------------------

func TestDelete_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/device/nABCD1234" {
			t.Errorf("expected path /device/nABCD1234, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := NewClient("tskey-test-token", WithBaseURL(srv.URL))
	body, err := c.Delete("/device/nABCD1234")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(body) != 0 {
		t.Errorf("expected empty body, got %q", string(body))
	}
}

// ---------------------------------------------------------------------------
// 10. TestPatch_Success
// ---------------------------------------------------------------------------

func TestPatch_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/device/nABCD1234CNTRL/routes" {
			t.Errorf("expected path /device/nABCD1234CNTRL/routes, got %s", r.URL.Path)
		}

		// Verify the request body
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}

		var payload map[string]interface{}
		if err := json.Unmarshal(reqBody, &payload); err != nil {
			t.Fatalf("request body is not valid JSON: %v", err)
		}

		// Return updated routes
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"routes":["10.0.0.0/24","192.168.1.0/24"]}`))
	}))
	defer srv.Close()

	c := NewClient("tskey-test-token", WithBaseURL(srv.URL))
	body, err := c.Patch("/device/nABCD1234CNTRL/routes",
		strings.NewReader(`{"routes":["10.0.0.0/24","192.168.1.0/24"]}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got map[string]interface{}
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	routes, ok := got["routes"].([]interface{})
	if !ok {
		t.Fatal("expected routes to be an array")
	}
	if len(routes) != 2 {
		t.Errorf("expected 2 routes, got %d", len(routes))
	}
}
