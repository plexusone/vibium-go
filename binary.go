package w3pilot

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Environment variables for specifying binary paths.
const (
	// VibiumBinaryEnvVar is the preferred environment variable for the vibium binary path.
	VibiumBinaryEnvVar = "VIBIUM_BIN_PATH"
	// ClickerBinaryEnvVar is the legacy environment variable (still supported).
	ClickerBinaryEnvVar = "CLICKER_BIN_PATH"
)

// binaryNames returns the binary names to search for based on OS.
func binaryNames() []string {
	if runtime.GOOS == "windows" {
		return []string{"clicker.exe", "vibium.exe"}
	}
	return []string{"clicker", "vibium"}
}

// FindClickerBinary locates the clicker binary.
//
// Search order:
//  1. VIBIUM_BIN_PATH environment variable
//  2. CLICKER_BIN_PATH environment variable (legacy)
//  3. PATH lookup for "clicker" or "vibium"
//  4. Go bin directory ($GOPATH/bin or ~/go/bin)
//  5. Common installation paths
func FindClickerBinary() (string, error) {
	// 1. Check VIBIUM_BIN_PATH environment variable
	if envPath := os.Getenv(VibiumBinaryEnvVar); envPath != "" {
		cleanPath := filepath.Clean(envPath)
		if _, err := os.Stat(cleanPath); err == nil {
			return cleanPath, nil
		}
		return "", fmt.Errorf("%s is set to %q but file does not exist", VibiumBinaryEnvVar, envPath)
	}

	// 2. Check CLICKER_BIN_PATH environment variable (legacy)
	if envPath := os.Getenv(ClickerBinaryEnvVar); envPath != "" {
		cleanPath := filepath.Clean(envPath)
		if _, err := os.Stat(cleanPath); err == nil {
			return cleanPath, nil
		}
		return "", fmt.Errorf("%s is set to %q but file does not exist", ClickerBinaryEnvVar, envPath)
	}

	// 3. Check PATH for clicker and vibium
	for _, name := range binaryNames() {
		if path, err := exec.LookPath(name); err == nil {
			return path, nil
		}
	}

	home, _ := os.UserHomeDir()

	// 4. Check Go bin directories
	goPaths := getGoBinPaths(home)
	for _, goPath := range goPaths {
		for _, name := range binaryNames() {
			binPath := filepath.Join(goPath, name)
			if _, err := os.Stat(binPath); err == nil {
				return binPath, nil
			}
		}
	}

	// 5. Check common installation paths
	commonPaths := getCommonBinaryPaths(home)
	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("clicker binary not found; build from https://github.com/VibiumDev/vibium/tree/main/clicker, or set %s", VibiumBinaryEnvVar)
}

// getGoBinPaths returns Go bin directories to search.
func getGoBinPaths(home string) []string {
	var paths []string

	// Check GOPATH/bin first
	if gopath := os.Getenv("GOPATH"); gopath != "" {
		paths = append(paths, filepath.Join(gopath, "bin"))
	}

	// Check GOBIN
	if gobin := os.Getenv("GOBIN"); gobin != "" {
		paths = append(paths, gobin)
	}

	// Default ~/go/bin
	if home != "" {
		paths = append(paths, filepath.Join(home, "go", "bin"))
	}

	return paths
}

// getCommonBinaryPaths returns common installation paths for clicker/vibium binaries.
func getCommonBinaryPaths(home string) []string {
	var paths []string

	names := binaryNames()

	switch runtime.GOOS {
	case "darwin", "linux":
		for _, name := range names {
			paths = append(paths,
				filepath.Join("/usr/local/bin", name),
				filepath.Join("/usr/bin", name),
			)
		}
		if home != "" {
			for _, name := range names {
				paths = append(paths,
					filepath.Join(home, ".local", "bin", name),
					filepath.Join(home, "bin", name),
				)
			}
		}
	case "windows":
		for _, name := range names {
			paths = append(paths,
				filepath.Join("C:\\Program Files\\vibium", name),
				filepath.Join("C:\\Program Files (x86)\\vibium", name),
			)
		}
		if home != "" {
			for _, name := range names {
				paths = append(paths,
					filepath.Join(home, "AppData", "Local", "Programs", "vibium", name),
				)
			}
		}
	}

	return paths
}

// ClickerVersion returns the version of the clicker/vibium binary.
func ClickerVersion(clickerPath string) (string, error) {
	cmd := exec.Command(clickerPath, "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get binary version: %w", err)
	}
	return string(output), nil
}
