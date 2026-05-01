package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	repoOwner    = "dimer47"
	repoName     = "tailscale-cli"
	cacheDuration = 24 * time.Hour
	cacheFile     = "update-check.json"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

type updateCache struct {
	CheckedAt     time.Time `json:"checked_at"`
	LatestVersion string    `json:"latest_version"`
}

// CheckResult contains the result of an update check.
type CheckResult struct {
	UpdateAvailable bool
	CurrentVersion  string
	LatestVersion   string
	Message         string
}

// Check verifies if a newer version is available on GitHub.
// It caches the result for 24h to avoid spamming the API.
func Check(currentVersion string) *CheckResult {
	if currentVersion == "dev" || currentVersion == "" {
		return nil
	}

	// Check cache first
	if cached := readCache(); cached != nil {
		if time.Since(cached.CheckedAt) < cacheDuration {
			return compareVersions(currentVersion, cached.LatestVersion)
		}
	}

	// Fetch latest release from GitHub
	latest, err := fetchLatestRelease()
	if err != nil {
		return nil
	}

	// Update cache
	writeCache(&updateCache{
		CheckedAt:     time.Now(),
		LatestVersion: latest,
	})

	return compareVersions(currentVersion, latest)
}

func fetchLatestRelease() (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return strings.TrimPrefix(release.TagName, "v"), nil
}

func compareVersions(current, latest string) *CheckResult {
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")

	if current == latest {
		return nil
	}

	if !isNewer(latest, current) {
		return nil
	}

	return &CheckResult{
		UpdateAvailable: true,
		CurrentVersion:  current,
		LatestVersion:   latest,
		Message:         formatUpdateMessage(current, latest),
	}
}

// isNewer returns true if a is newer than b (simple semver comparison).
func isNewer(a, b string) bool {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	for i := 0; i < len(aParts) && i < len(bParts); i++ {
		if aParts[i] != bParts[i] {
			return aParts[i] > bParts[i]
		}
	}

	return len(aParts) > len(bParts)
}

func formatUpdateMessage(current, latest string) string {
	return fmt.Sprintf(
		"\n  Une nouvelle version de tailscale-cli est disponible : %s → %s\n  Mise à jour :  tailscale-cli self-update\n",
		current, latest,
	)
}

// ClearCache removes the update check cache file.
func ClearCache() {
	path := cachePath()
	if path != "" {
		_ = os.Remove(path)
	}
}

func cachePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".tailscale-cli", cacheFile)
}

func readCache() *updateCache {
	path := cachePath()
	if path == "" {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var cache updateCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil
	}

	return &cache
}

func writeCache(cache *updateCache) {
	path := cachePath()
	if path == "" {
		return
	}

	data, err := json.Marshal(cache)
	if err != nil {
		return
	}

	// Ensure directory exists
	_ = os.MkdirAll(filepath.Dir(path), 0700)
	_ = os.WriteFile(path, data, 0600)
}
