package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// sampleConfig returns a Config with realistic test data.
func sampleConfig() *Config {
	return &Config{
		DefaultContext: "work",
		Contexts: map[string]Context{
			"work":     {Tailnet: "example.com", APIToken: "tskey-api-fake-token"},
			"personal": {Tailnet: "-"},
		},
	}
}

// writeTempConfig writes JSON data to a temporary file and returns its path.
func writeTempConfig(t *testing.T, data []byte) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return path
}

func TestLoad_ValidConfig(t *testing.T) {
	cfg := sampleConfig()
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	path := writeTempConfig(t, data)

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if loaded.DefaultContext != "work" {
		t.Errorf("DefaultContext = %q, want %q", loaded.DefaultContext, "work")
	}
	if len(loaded.Contexts) != 2 {
		t.Fatalf("len(Contexts) = %d, want 2", len(loaded.Contexts))
	}
	if loaded.Contexts["work"].Tailnet != "example.com" {
		t.Errorf("work.Tailnet = %q, want %q", loaded.Contexts["work"].Tailnet, "example.com")
	}
	if loaded.Contexts["personal"].Tailnet != "-" {
		t.Errorf("personal.Tailnet = %q, want %q", loaded.Contexts["personal"].Tailnet, "-")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/tmp/does-not-exist-config-test.json")
	if err == nil {
		t.Fatal("Load should return an error for a missing file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	path := writeTempConfig(t, []byte(`{invalid json`))

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load should return an error for invalid JSON")
	}
}

func TestSave_CreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	nested := filepath.Join(dir, "deep", "nested", "config.json")

	cfg := sampleConfig()
	if err := cfg.Save(nested); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	if _, err := os.Stat(nested); os.IsNotExist(err) {
		t.Fatal("Save did not create the config file in nested directory")
	}
}

func TestSave_StripTokens(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	cfg := sampleConfig()
	if err := cfg.Save(path); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading saved file: %v", err)
	}

	var ondisk Config
	if err := json.Unmarshal(data, &ondisk); err != nil {
		t.Fatalf("unmarshal saved file: %v", err)
	}

	for name, ctx := range ondisk.Contexts {
		if ctx.APIToken != "" {
			t.Errorf("context %q: APIToken should be empty on disk, got %q", name, ctx.APIToken)
		}
	}
}

func TestSave_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	original := sampleConfig()
	if err := original.Save(path); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if loaded.DefaultContext != original.DefaultContext {
		t.Errorf("DefaultContext = %q, want %q", loaded.DefaultContext, original.DefaultContext)
	}
	if len(loaded.Contexts) != len(original.Contexts) {
		t.Fatalf("len(Contexts) = %d, want %d", len(loaded.Contexts), len(original.Contexts))
	}
	for name, origCtx := range original.Contexts {
		loadedCtx, ok := loaded.Contexts[name]
		if !ok {
			t.Errorf("context %q missing after round-trip", name)
			continue
		}
		if loadedCtx.Tailnet != origCtx.Tailnet {
			t.Errorf("context %q: Tailnet = %q, want %q", name, loadedCtx.Tailnet, origCtx.Tailnet)
		}
		// Token must NOT survive the round-trip (stripped by Save)
		if loadedCtx.APIToken != "" {
			t.Errorf("context %q: APIToken should be empty after round-trip, got %q", name, loadedCtx.APIToken)
		}
	}
}

func TestGetActiveContext_ExplicitName(t *testing.T) {
	cfg := sampleConfig()

	ctx, name, err := GetActiveContext(cfg, "personal")
	if err != nil {
		t.Fatalf("GetActiveContext returned error: %v", err)
	}
	if name != "personal" {
		t.Errorf("resolved name = %q, want %q", name, "personal")
	}
	if ctx.Tailnet != "-" {
		t.Errorf("Tailnet = %q, want %q", ctx.Tailnet, "-")
	}
}

func TestGetActiveContext_DefaultFallback(t *testing.T) {
	cfg := sampleConfig()

	tests := []struct {
		input string
		label string
	}{
		{"", "empty string"},
		{"default", "literal default"},
	}

	for _, tt := range tests {
		ctx, name, err := GetActiveContext(cfg, tt.input)
		if err != nil {
			t.Fatalf("[%s] GetActiveContext returned error: %v", tt.label, err)
		}
		if name != "work" {
			t.Errorf("[%s] resolved name = %q, want %q", tt.label, name, "work")
		}
		if ctx.Tailnet != "example.com" {
			t.Errorf("[%s] Tailnet = %q, want %q", tt.label, ctx.Tailnet, "example.com")
		}
	}
}

func TestGetActiveContext_NotFound(t *testing.T) {
	cfg := sampleConfig()

	_, _, err := GetActiveContext(cfg, "nonexistent")
	if err == nil {
		t.Fatal("GetActiveContext should return an error for an unknown context")
	}
}

func TestDefaultConfigPath(t *testing.T) {
	path := DefaultConfigPath()
	expected := filepath.Join(".tailscale-cli", "config.json")
	if !filepath.IsAbs(path) {
		// If UserHomeDir fails the path starts with "." which is not absolute,
		// but it should still contain the expected suffix.
	}
	if !containsSuffix(path, expected) {
		t.Errorf("DefaultConfigPath() = %q, should contain %q", path, expected)
	}
}

func containsSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}
