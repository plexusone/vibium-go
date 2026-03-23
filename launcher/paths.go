// Package launcher provides Chrome browser launching and installation.
package launcher

import (
	"os"
	"path/filepath"
	"runtime"
)

// Platform represents a browser download platform.
type Platform string

const (
	PlatformMacArm64 Platform = "mac-arm64"
	PlatformMacX64   Platform = "mac-x64"
	PlatformLinux64  Platform = "linux64"
	PlatformWin32    Platform = "win32"
	PlatformWin64    Platform = "win64"
)

// GetPlatform returns the current platform string for Chrome for Testing downloads.
func GetPlatform() Platform {
	switch runtime.GOOS {
	case "darwin":
		if runtime.GOARCH == "arm64" {
			return PlatformMacArm64
		}
		return PlatformMacX64
	case "linux":
		return PlatformLinux64
	case "windows":
		if runtime.GOARCH == "386" {
			return PlatformWin32
		}
		return PlatformWin64
	default:
		return PlatformLinux64
	}
}

// GetCacheDir returns the platform-specific cache directory for vibium.
func GetCacheDir() (string, error) {
	var baseDir string

	switch runtime.GOOS {
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		baseDir = filepath.Join(home, "Library", "Caches", "vibium")
	case "windows":
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			baseDir = filepath.Join(localAppData, "vibium")
		} else {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			baseDir = filepath.Join(home, "AppData", "Local", "vibium")
		}
	default: // Linux and others
		if xdgCache := os.Getenv("XDG_CACHE_HOME"); xdgCache != "" {
			baseDir = filepath.Join(xdgCache, "vibium")
		} else {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			baseDir = filepath.Join(home, ".cache", "vibium")
		}
	}

	return baseDir, nil
}

// GetChromeDir returns the directory where Chrome for Testing is installed.
func GetChromeDir() (string, error) {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cacheDir, "chrome-for-testing"), nil
}

// GetChromePath returns the path to the Chrome executable for the given version.
// If version is empty, it looks for any installed version.
func GetChromePath(version string) (string, error) {
	chromeDir, err := GetChromeDir()
	if err != nil {
		return "", err
	}

	// If no version specified, find first available
	if version == "" {
		entries, err := os.ReadDir(chromeDir)
		if err != nil {
			return "", err
		}
		for _, entry := range entries {
			if entry.IsDir() {
				version = entry.Name()
				break
			}
		}
		if version == "" {
			return "", os.ErrNotExist
		}
	}

	versionDir := filepath.Join(chromeDir, version)

	// Platform-specific executable path
	var execPath string
	switch runtime.GOOS {
	case "darwin":
		execPath = filepath.Join(versionDir, "Google Chrome for Testing.app", "Contents", "MacOS", "Google Chrome for Testing")
	case "windows":
		execPath = filepath.Join(versionDir, "chrome.exe")
	default: // Linux
		execPath = filepath.Join(versionDir, "chrome")
	}

	return execPath, nil
}

// FindSystemChrome looks for Chrome installed on the system.
func FindSystemChrome() (string, error) {
	var paths []string

	switch runtime.GOOS {
	case "darwin":
		paths = []string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"/Applications/Google Chrome for Testing.app/Contents/MacOS/Google Chrome for Testing",
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
		}
		// Also check user Applications
		if home, err := os.UserHomeDir(); err == nil {
			paths = append(paths,
				filepath.Join(home, "Applications", "Google Chrome.app", "Contents", "MacOS", "Google Chrome"),
				filepath.Join(home, "Applications", "Google Chrome for Testing.app", "Contents", "MacOS", "Google Chrome for Testing"),
			)
		}
	case "linux":
		paths = []string{
			"/usr/bin/google-chrome",
			"/usr/bin/google-chrome-stable",
			"/usr/bin/chromium",
			"/usr/bin/chromium-browser",
			"/snap/bin/chromium",
		}
	case "windows":
		programFiles := os.Getenv("PROGRAMFILES")
		programFilesX86 := os.Getenv("PROGRAMFILES(X86)")
		localAppData := os.Getenv("LOCALAPPDATA")

		if programFiles != "" {
			paths = append(paths, filepath.Join(programFiles, "Google", "Chrome", "Application", "chrome.exe"))
		}
		if programFilesX86 != "" {
			paths = append(paths, filepath.Join(programFilesX86, "Google", "Chrome", "Application", "chrome.exe"))
		}
		if localAppData != "" {
			paths = append(paths, filepath.Join(localAppData, "Google", "Chrome", "Application", "chrome.exe"))
		}
	}

	for _, p := range paths {
		//nolint:gosec // G703: paths are hardcoded known Chrome installation locations
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	return "", os.ErrNotExist
}
