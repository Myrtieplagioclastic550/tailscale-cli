package keychain

import (
	"testing"
)

func TestAccountName(t *testing.T) {
	tests := []struct {
		context  string
		expected string
	}{
		{"default", "tailscale-cli:default"},
		{"work", "tailscale-cli:work"},
		{"personal", "tailscale-cli:personal"},
		{"client-a", "tailscale-cli:client-a"},
		{"", "tailscale-cli:"},
	}

	for _, tt := range tests {
		t.Run(tt.context, func(t *testing.T) {
			result := accountName(tt.context)
			if result != tt.expected {
				t.Errorf("accountName(%q) = %q, want %q", tt.context, result, tt.expected)
			}
		})
	}
}

func TestServiceName(t *testing.T) {
	if serviceName != "tailscale-cli" {
		t.Errorf("serviceName = %q, want %q", serviceName, "tailscale-cli")
	}
}
