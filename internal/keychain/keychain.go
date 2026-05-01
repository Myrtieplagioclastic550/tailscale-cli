package keychain

import (
	"github.com/zalando/go-keyring"
)

const (
	serviceName = "tailscale-cli"
)

// Set stores a token in the system keyring for the given context.
func Set(context, token string) error {
	return keyring.Set(serviceName, accountName(context), token)
}

// Get retrieves a token from the system keyring for the given context.
// Returns empty string and nil error if the entry does not exist.
func Get(context string) (string, error) {
	token, err := keyring.Get(serviceName, accountName(context))
	if err == keyring.ErrNotFound {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return token, nil
}

// Delete removes a token from the system keyring for the given context.
func Delete(context string) error {
	err := keyring.Delete(serviceName, accountName(context))
	if err == keyring.ErrNotFound {
		return nil
	}
	return err
}

// IsAvailable checks if the system keyring is available.
func IsAvailable() bool {
	// Try a Get on a dummy key — if the backend is unavailable,
	// go-keyring returns an error other than ErrNotFound.
	_, err := keyring.Get(serviceName, "__keyring_probe__")
	return err == nil || err == keyring.ErrNotFound
}

func accountName(context string) string {
	return serviceName + ":" + context
}
