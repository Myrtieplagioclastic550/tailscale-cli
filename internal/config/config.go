package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the CLI configuration.
type Config struct {
	DefaultContext string             `json:"default_context"`
	Contexts       map[string]Context `json:"contexts"`
}

// Context holds connection details for a Tailscale account.
// APIToken is no longer stored in the config file — it lives in the macOS Keychain.
type Context struct {
	APIToken string `json:"api_token,omitempty"` // Legacy: kept for migration, not written for new contexts
	Tailnet  string `json:"tailnet"`
}

// Load reads the configuration from the given file path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if cfg.Contexts == nil {
		cfg.Contexts = make(map[string]Context)
	}

	return &cfg, nil
}

// Save writes the configuration to the given file path.
// Tokens are NOT written to the config file — they are stored in the Keychain.
func (c *Config) Save(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	// Strip tokens before saving to file
	sanitized := Config{
		DefaultContext: c.DefaultContext,
		Contexts:       make(map[string]Context, len(c.Contexts)),
	}
	for name, ctx := range c.Contexts {
		sanitized.Contexts[name] = Context{
			Tailnet: ctx.Tailnet,
		}
	}

	data, err := json.MarshalIndent(sanitized, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}

// DefaultConfigPath returns the default path for the configuration file.
func DefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".tailscale-cli", "config.json")
}

// GetActiveContext resolves the active context. If contextName is non-empty,
// it looks up that context; otherwise it uses the default context from the config.
func GetActiveContext(config *Config, contextName string) (*Context, string, error) {
	name := contextName
	if name == "" || name == "default" {
		if config.DefaultContext != "" {
			name = config.DefaultContext
		} else {
			name = "default"
		}
	}

	ctx, ok := config.Contexts[name]
	if !ok {
		return nil, name, fmt.Errorf("context %q not found in configuration", name)
	}

	return &ctx, name, nil
}
