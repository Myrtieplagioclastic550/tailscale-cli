package keychain

import (
	"fmt"
	"os/exec"
	"strings"
)

const (
	serviceName = "tailscale-cli"
)

// Set stores a token in the macOS Keychain for the given context.
func Set(context, token string) error {
	account := accountName(context)

	// Delete any existing entry first (ignore errors if not found)
	_ = Delete(context)

	cmd := exec.Command("security", "add-generic-password",
		"-s", serviceName,
		"-a", account,
		"-w", token,
		"-U", // update if exists
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("impossible de stocker le token dans le Keychain : %s (%w)", strings.TrimSpace(string(out)), err)
	}

	return nil
}

// Get retrieves a token from the macOS Keychain for the given context.
// Returns empty string and nil error if the entry does not exist.
func Get(context string) (string, error) {
	account := accountName(context)

	cmd := exec.Command("security", "find-generic-password",
		"-s", serviceName,
		"-a", account,
		"-w", // output password only
	)
	out, err := cmd.Output()
	if err != nil {
		// Entry not found is not an error for our purposes
		return "", nil
	}

	return strings.TrimSpace(string(out)), nil
}

// Delete removes a token from the macOS Keychain for the given context.
func Delete(context string) error {
	account := accountName(context)

	cmd := exec.Command("security", "delete-generic-password",
		"-s", serviceName,
		"-a", account,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		outStr := strings.TrimSpace(string(out))
		// "not found" is not a real error
		if strings.Contains(outStr, "could not be found") {
			return nil
		}
		return fmt.Errorf("impossible de supprimer le token du Keychain : %s (%w)", outStr, err)
	}

	return nil
}

// List returns all context names that have a token stored in the Keychain.
func List() ([]string, error) {
	cmd := exec.Command("security", "dump-keychain")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("impossible de lire le Keychain : %w", err)
	}

	var contexts []string
	lines := strings.Split(string(out), "\n")
	inTelscaleEntry := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check for our service name
		if strings.Contains(trimmed, fmt.Sprintf(`"svce"<blob>="%s"`, serviceName)) {
			inTelscaleEntry = true
			continue
		}

		// If we're in a tailscale-cli entry, look for the account name
		if inTelscaleEntry && strings.Contains(trimmed, `"acct"<blob>="`) {
			start := strings.Index(trimmed, `"acct"<blob>="`) + len(`"acct"<blob>="`)
			end := strings.LastIndex(trimmed, `"`)
			if start < end {
				account := trimmed[start:end]
				// Strip the "tailscale-cli:" prefix
				if strings.HasPrefix(account, serviceName+":") {
					context := strings.TrimPrefix(account, serviceName+":")
					contexts = append(contexts, context)
				}
			}
			inTelscaleEntry = false
		}
	}

	return contexts, nil
}

// IsAvailable checks if the macOS Keychain is available (i.e. we're on macOS).
func IsAvailable() bool {
	_, err := exec.LookPath("security")
	return err == nil
}

func accountName(context string) string {
	return serviceName + ":" + context
}
