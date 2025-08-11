package chrome

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// ChromeVersion represents a Chrome version with metadata
type ChromeVersion struct {
	Version     string    `json:"version"`
	ReleaseDate time.Time `json:"release_date"`
	Channel     string    `json:"channel"`
	Platform    string    `json:"platform"`
}

// GetLatestStableVersion returns the latest Chrome Stable version
func GetLatestStableVersion() (string, error) {
	// Try multiple sources for reliability
	sources := []func() (string, error){
		getLatestFromOmahaProxy,
		getLatestFromChromeReleases,
		getLatestFromChromeVersionAPI,
	}

	for _, source := range sources {
		if version, err := source(); err == nil {
			return version, nil
		}
	}

	return "", fmt.Errorf("failed to get latest Chrome version from all sources")
}

// getLatestFromOmahaProxy fetches from Google's Omaha Proxy
func getLatestFromOmahaProxy() (string, error) {
	url := "https://omahaproxy.appspot.com/all.json"
	
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch from Omaha Proxy: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var releases []struct {
		OS      string `json:"os"`
		Channel string `json:"channel"`
		Version string `json:"version"`
	}

	if err := json.Unmarshal(body, &releases); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Find stable version for current platform
	platform := getCurrentPlatform()
	for _, release := range releases {
		if release.OS == platform && release.Channel == "stable" {
			return release.Version, nil
		}
	}

	return "", fmt.Errorf("stable version not found for platform %s", platform)
}

// getLatestFromChromeReleases fetches from Chrome Releases blog
func getLatestFromChromeReleases() (string, error) {
	url := "https://chromereleases.googleblog.com/feeds/posts/default"
	
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch from Chrome Releases: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Parse RSS feed for latest stable release
	content := string(body)
	
	// Look for stable version pattern
	re := regexp.MustCompile(`Chrome (\d+\.\d+\.\d+\.\d+) Stable`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1], nil
	}

	return "", fmt.Errorf("stable version not found in RSS feed")
}

// getLatestFromChromeVersionAPI fetches from Chrome Version API
func getLatestFromChromeVersionAPI() (string, error) {
	url := "https://versionhistory.googleapis.com/v1/chrome/platforms/win/channels/stable/versions"
	
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch from Chrome Version API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var response struct {
		Versions []struct {
			Version string `json:"version"`
		} `json:"versions"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	if len(response.Versions) > 0 {
		return response.Versions[0].Version, nil
	}

	return "", fmt.Errorf("no versions found in API response")
}

// getCurrentPlatform returns the current platform identifier
func getCurrentPlatform() string {
	// This would be more sophisticated in practice
	// For now, return a common platform
	return "win"
}

// ParseVersion parses a Chrome version string into components
func ParseVersion(version string) (major, minor, patch, build int, err error) {
	parts := strings.Split(version, ".")
	if len(parts) != 4 {
		return 0, 0, 0, 0, fmt.Errorf("invalid version format: %s", version)
	}

	// Parse each component
	var components [4]int
	for i, part := range parts {
		if components[i], err = parseInt(part); err != nil {
			return 0, 0, 0, 0, fmt.Errorf("invalid version component %d: %w", i+1, err)
		}
	}

	return components[0], components[1], components[2], components[3], nil
}

// parseInt is a helper function to parse integers
func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}

// CompareVersions compares two Chrome versions
func CompareVersions(v1, v2 string) (int, error) {
	major1, minor1, patch1, build1, err := ParseVersion(v1)
	if err != nil {
		return 0, fmt.Errorf("invalid first version: %w", err)
	}

	major2, minor2, patch2, build2, err := ParseVersion(v2)
	if err != nil {
		return 0, fmt.Errorf("invalid second version: %w", err)
	}

	// Compare major version
	if major1 != major2 {
		if major1 > major2 {
			return 1, nil
		}
		return -1, nil
	}

	// Compare minor version
	if minor1 != minor2 {
		if minor1 > minor2 {
			return 1, nil
		}
		return -1, nil
	}

	// Compare patch version
	if patch1 != patch2 {
		if patch1 > patch2 {
			return 1, nil
		}
		return -1, nil
	}

	// Compare build number
	if build1 != build2 {
		if build1 > build2 {
			return 1, nil
		}
		return -1, nil
	}

	return 0, nil // Versions are equal
}
