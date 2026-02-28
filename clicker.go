package vibium

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// ClickerProcess manages the clicker binary subprocess.
type ClickerProcess struct {
	cmd     *exec.Cmd
	port    int
	wsURL   string
	stopped bool
}

// findClickerBinary locates the clicker binary.
func findClickerBinary(customPath string) (string, error) {
	// 1. Check custom path
	if customPath != "" {
		if _, err := os.Stat(customPath); err == nil {
			return customPath, nil
		}
	}

	// 2. Check VIBIUM_CLICKER_PATH environment variable
	if envPath := os.Getenv("VIBIUM_CLICKER_PATH"); envPath != "" {
		//nolint:gosec // G703: envPath is trusted user configuration from environment
		if _, err := os.Stat(envPath); err == nil {
			return envPath, nil
		}
	}

	// 3. Check system PATH
	if path, err := exec.LookPath("clicker"); err == nil {
		return path, nil
	}

	// 4. Check platform cache directory
	cacheDir := getCacheDir()
	binaryName := "clicker"
	if runtime.GOOS == "windows" {
		binaryName = "clicker.exe"
	}
	cachePath := filepath.Join(cacheDir, binaryName)
	if _, err := os.Stat(cachePath); err == nil {
		return cachePath, nil
	}

	// 5. Check local development paths
	localPaths := []string{
		filepath.Join(".", "clicker", "bin", binaryName),
		filepath.Join("..", "..", "clicker", "bin", binaryName),
	}
	for _, p := range localPaths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	return "", ErrClickerNotFound
}

// getCacheDir returns the platform-specific cache directory.
func getCacheDir() string {
	switch runtime.GOOS {
	case "darwin":
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "Library", "Caches", "vibium")
	case "windows":
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			return filepath.Join(localAppData, "vibium")
		}
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "AppData", "Local", "vibium")
	default: // linux and others
		if xdgCache := os.Getenv("XDG_CACHE_HOME"); xdgCache != "" {
			return filepath.Join(xdgCache, "vibium")
		}
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".cache", "vibium")
	}
}

// StartClicker starts the clicker binary and returns a ClickerProcess.
func StartClicker(ctx context.Context, opts LaunchOptions) (*ClickerProcess, error) {
	binaryPath, err := findClickerBinary(opts.ExecutablePath)
	if err != nil {
		return nil, err
	}

	args := []string{"serve"}
	if opts.Port > 0 {
		args = append(args, "--port", strconv.Itoa(opts.Port))
	}
	if opts.Headless {
		args = append(args, "--headless")
	}

	cmd := exec.CommandContext(ctx, binaryPath, args...)
	cmd.Stderr = os.Stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start clicker: %w", err)
	}

	// Wait for server to start and parse WebSocket URL
	wsURL := ""
	port := 0
	scanner := bufio.NewScanner(stdout)
	urlRegex := regexp.MustCompile(`ws://[^:]+:(\d+)`)

	// Give it time to start
	done := make(chan struct{})
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "Server listening on") {
				matches := urlRegex.FindStringSubmatch(line)
				if len(matches) >= 2 {
					wsURL = matches[0]
					port, _ = strconv.Atoi(matches[1])
					close(done)
					return
				}
			}
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(30 * time.Second):
		_ = cmd.Process.Kill()
		return nil, fmt.Errorf("timeout waiting for clicker to start")
	case <-ctx.Done():
		_ = cmd.Process.Kill()
		return nil, ctx.Err()
	}

	if wsURL == "" {
		_ = cmd.Process.Kill()
		return nil, fmt.Errorf("failed to parse WebSocket URL from clicker output")
	}

	return &ClickerProcess{
		cmd:   cmd,
		port:  port,
		wsURL: wsURL,
	}, nil
}

// WebSocketURL returns the WebSocket URL for connecting to the clicker.
func (p *ClickerProcess) WebSocketURL() string {
	return p.wsURL
}

// Port returns the port the clicker is listening on.
func (p *ClickerProcess) Port() int {
	return p.port
}

// Stop gracefully stops the clicker process.
func (p *ClickerProcess) Stop() error {
	if p.stopped {
		return nil
	}
	p.stopped = true

	if p.cmd == nil || p.cmd.Process == nil {
		return nil
	}

	// Try graceful shutdown first
	done := make(chan error, 1)
	go func() {
		done <- p.cmd.Wait()
	}()

	// Send interrupt signal
	_ = p.cmd.Process.Signal(os.Interrupt)

	select {
	case <-done:
		return nil
	case <-time.After(5 * time.Second):
		// Force kill if graceful shutdown fails
		return p.cmd.Process.Kill()
	}
}

// Wait waits for the clicker process to exit.
func (p *ClickerProcess) Wait() error {
	if p.cmd == nil {
		return nil
	}
	return p.cmd.Wait()
}
