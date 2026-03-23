package launcher

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	// Chrome for Testing API endpoint
	lastKnownGoodURL = "https://googlechromelabs.github.io/chrome-for-testing/last-known-good-versions-with-downloads.json"
)

// VersionInfo represents Chrome for Testing version information.
type VersionInfo struct {
	Version   string                `json:"version"`
	Downloads map[string][]Download `json:"downloads"`
}

// Download represents a download URL for a specific platform.
type Download struct {
	Platform string `json:"platform"`
	URL      string `json:"url"`
}

// LastKnownGoodResponse represents the API response for last known good versions.
type LastKnownGoodResponse struct {
	Channels map[string]VersionInfo `json:"channels"`
}

// InstallResult contains the paths to installed binaries.
type InstallResult struct {
	ChromePath string
	Version    string
}

// Install downloads and installs Chrome for Testing.
// Returns the path to the installed Chrome executable.
// Skips download if already installed.
func Install() (*InstallResult, error) {
	// Check for skip environment variable
	if os.Getenv("VIBIUM_SKIP_BROWSER_DOWNLOAD") == "1" {
		return nil, fmt.Errorf("browser download skipped (VIBIUM_SKIP_BROWSER_DOWNLOAD=1)")
	}

	// Check if already installed
	chromePath, err := GetChromePath("")
	if err == nil {
		if _, statErr := os.Stat(chromePath); statErr == nil {
			version := extractVersionFromPath(chromePath)
			return &InstallResult{
				ChromePath: chromePath,
				Version:    version,
			}, nil
		}
	}

	platform := GetPlatform()

	// Fetch latest stable version info
	versionInfo, err := fetchLatestStableVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch version info: %w", err)
	}

	fmt.Printf("Installing Chrome for Testing v%s...\n", versionInfo.Version)

	// Create version directory
	chromeDir, err := GetChromeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get chrome dir: %w", err)
	}

	versionDir := filepath.Join(chromeDir, versionInfo.Version)
	if err := os.MkdirAll(versionDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create version dir: %w", err)
	}

	// Download and extract Chrome
	chromeURL := findDownloadURL(versionInfo.Downloads["chrome"], string(platform))
	if chromeURL == "" {
		return nil, fmt.Errorf("no Chrome download available for platform %s", platform)
	}

	fmt.Printf("Downloading Chrome from %s...\n", chromeURL)
	if err := downloadAndExtract(chromeURL, versionDir); err != nil {
		return nil, fmt.Errorf("failed to install Chrome: %w", err)
	}

	// Get path to installed binary
	chromePath, err = GetChromePath(versionInfo.Version)
	if err != nil {
		return nil, fmt.Errorf("Chrome installed but not found: %w", err)
	}

	// Make executable on Unix
	if runtime.GOOS != "windows" {
		if err := os.Chmod(chromePath, 0755); err != nil {
			return nil, fmt.Errorf("failed to chmod: %w", err)
		}
	}

	// Remove quarantine attribute on macOS to avoid Gatekeeper prompts
	if runtime.GOOS == "darwin" {
		// Find the .app bundle
		appPath := chromePath
		for !strings.HasSuffix(appPath, ".app") && appPath != "/" {
			appPath = filepath.Dir(appPath)
		}
		if strings.HasSuffix(appPath, ".app") {
			_ = exec.Command("xattr", "-rd", "com.apple.quarantine", appPath).Run()
		}
	}

	fmt.Printf("Chrome for Testing v%s installed successfully.\n", versionInfo.Version)

	return &InstallResult{
		ChromePath: chromePath,
		Version:    versionInfo.Version,
	}, nil
}

// fetchLatestStableVersion fetches the latest stable Chrome for Testing version.
func fetchLatestStableVersion() (*VersionInfo, error) {
	resp, err := http.Get(lastKnownGoodURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var data LastKnownGoodResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	stable, ok := data.Channels["Stable"]
	if !ok {
		return nil, fmt.Errorf("no Stable channel found")
	}

	return &stable, nil
}

// findDownloadURL finds the download URL for the given platform.
func findDownloadURL(downloads []Download, platform string) string {
	for _, d := range downloads {
		if d.Platform == platform {
			return d.URL
		}
	}
	return ""
}

// downloadAndExtract downloads a zip file and extracts it to the destination.
func downloadAndExtract(url, destDir string) error {
	// Download to temp file
	//nolint:gosec // G107: URL is from trusted Chrome for Testing API
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp("", "chrome-*.zip")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// Show download progress
	written, err := io.Copy(tmpFile, resp.Body)
	if err != nil {
		tmpFile.Close()
		return err
	}
	tmpFile.Close()
	fmt.Printf("Downloaded %.1f MB\n", float64(written)/1024/1024)

	// Extract zip
	fmt.Println("Extracting...")
	return extractZip(tmpPath, destDir)
}

// extractZip extracts a zip file to the destination directory.
func extractZip(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		// Strip the top-level directory (e.g. "chrome-mac-arm64/..." → "...")
		name := f.Name
		if i := strings.IndexByte(name, '/'); i >= 0 {
			name = name[i+1:]
		}
		if name == "" {
			continue
		}

		//nolint:gosec // G305: We check for zip slip below
		fpath := filepath.Join(destDir, name)

		// Security check: prevent zip slip
		if !strings.HasPrefix(fpath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		//nolint:gosec // G110: Trusted source (Chrome for Testing), known file sizes
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

// IsInstalled checks if Chrome for Testing is already installed.
func IsInstalled() bool {
	chromePath, err := GetChromePath("")
	if err != nil {
		return false
	}
	_, err = os.Stat(chromePath)
	return err == nil
}

// extractVersionFromPath extracts the version number from a Chrome path.
func extractVersionFromPath(path string) string {
	parts := strings.Split(path, string(os.PathSeparator))
	for i, part := range parts {
		if part == "chrome-for-testing" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return "unknown"
}
