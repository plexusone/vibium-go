package launcher

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"time"
)

// Options configures the browser launch.
type Options struct {
	// Headless runs the browser in headless mode.
	Headless bool

	// Port specifies the debugging port. If 0, a random port is used.
	Port int

	// ExecutablePath specifies a custom Chrome executable path.
	ExecutablePath string

	// UserDataDir specifies a custom user data directory.
	// If empty, a temporary directory is used.
	UserDataDir string

	// Args specifies additional command-line arguments.
	Args []string

	// AutoInstall automatically downloads Chrome if not found.
	// Defaults to true.
	AutoInstall *bool
}

// Browser represents a running browser instance.
type Browser struct {
	cmd     *exec.Cmd
	wsURL   string
	port    int
	dataDir string
	tempDir bool // whether dataDir is temporary and should be cleaned up
	stopped bool
}

// Launch starts a new Chrome browser with BiDi support.
func Launch(ctx context.Context, opts *Options) (*Browser, error) {
	if opts == nil {
		opts = &Options{}
	}

	// Find Chrome executable
	chromePath, err := findChrome(opts)
	if err != nil {
		return nil, err
	}

	// Determine port
	port := opts.Port
	if port == 0 {
		port, err = findFreePort()
		if err != nil {
			return nil, fmt.Errorf("failed to find free port: %w", err)
		}
	}

	// Determine user data directory
	dataDir := opts.UserDataDir
	tempDir := false
	if dataDir == "" {
		dataDir, err = os.MkdirTemp("", "vibium-chrome-*")
		if err != nil {
			return nil, fmt.Errorf("failed to create temp dir: %w", err)
		}
		tempDir = true
	}

	// Build command-line arguments
	args := buildArgs(port, dataDir, opts)

	// Start Chrome
	cmd := exec.CommandContext(ctx, chromePath, args...)
	cmd.Dir = filepath.Dir(chromePath)

	// Capture stderr for WebSocket URL
	stderr, err := cmd.StderrPipe()
	if err != nil {
		if tempDir {
			os.RemoveAll(dataDir)
		}
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		if tempDir {
			os.RemoveAll(dataDir)
		}
		return nil, fmt.Errorf("failed to start Chrome: %w", err)
	}

	// Parse WebSocket URL from stderr
	wsURL, err := parseWebSocketURL(ctx, stderr)
	if err != nil {
		_ = cmd.Process.Kill()
		if tempDir {
			os.RemoveAll(dataDir)
		}
		return nil, fmt.Errorf("failed to get WebSocket URL: %w", err)
	}

	return &Browser{
		cmd:     cmd,
		wsURL:   wsURL,
		port:    port,
		dataDir: dataDir,
		tempDir: tempDir,
	}, nil
}

// WebSocketURL returns the WebSocket URL for BiDi connection.
func (b *Browser) WebSocketURL() string {
	return b.wsURL
}

// Port returns the debugging port.
func (b *Browser) Port() int {
	return b.port
}

// Stop gracefully stops the browser.
func (b *Browser) Stop() error {
	if b.stopped {
		return nil
	}
	b.stopped = true

	if b.cmd == nil || b.cmd.Process == nil {
		return nil
	}

	// Try graceful shutdown first
	done := make(chan error, 1)
	go func() {
		done <- b.cmd.Wait()
	}()

	// Send interrupt signal
	_ = b.cmd.Process.Signal(os.Interrupt)

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		// Force kill if graceful shutdown fails
		_ = b.cmd.Process.Kill()
		<-done
	}

	// Clean up temp directory
	if b.tempDir && b.dataDir != "" {
		os.RemoveAll(b.dataDir)
	}

	return nil
}

// Wait waits for the browser process to exit.
func (b *Browser) Wait() error {
	if b.cmd == nil {
		return nil
	}
	return b.cmd.Wait()
}

// findChrome locates the Chrome executable.
func findChrome(opts *Options) (string, error) {
	// 1. Check custom path
	if opts.ExecutablePath != "" {
		if _, err := os.Stat(opts.ExecutablePath); err == nil {
			return opts.ExecutablePath, nil
		}
		return "", fmt.Errorf("chrome not found at %s", opts.ExecutablePath)
	}

	// 2. Check environment variable
	if envPath := os.Getenv("CHROME_PATH"); envPath != "" {
		//nolint:gosec // G703: envPath is trusted user configuration from environment
		if _, err := os.Stat(envPath); err == nil {
			return envPath, nil
		}
	}

	// 3. Check installed Chrome for Testing
	if chromePath, err := GetChromePath(""); err == nil {
		if _, err := os.Stat(chromePath); err == nil {
			return chromePath, nil
		}
	}

	// 4. Check system Chrome
	if chromePath, err := FindSystemChrome(); err == nil {
		return chromePath, nil
	}

	// 5. Auto-install if enabled (default)
	autoInstall := opts.AutoInstall == nil || *opts.AutoInstall
	if autoInstall {
		result, err := Install()
		if err != nil {
			return "", fmt.Errorf("Chrome not found and auto-install failed: %w", err)
		}
		return result.ChromePath, nil
	}

	return "", fmt.Errorf("Chrome not found. Install Chrome or set CHROME_PATH environment variable")
}

// buildArgs constructs the Chrome command-line arguments.
func buildArgs(port int, dataDir string, opts *Options) []string {
	args := []string{
		// BiDi WebSocket endpoint
		fmt.Sprintf("--remote-debugging-port=%d", port),

		// User data directory
		fmt.Sprintf("--user-data-dir=%s", dataDir),

		// Disable first run experience
		"--no-first-run",
		"--no-default-browser-check",

		// Disable background features that can interfere
		"--disable-background-networking",
		"--disable-background-timer-throttling",
		"--disable-backgrounding-occluded-windows",
		"--disable-breakpad",
		"--disable-client-side-phishing-detection",
		"--disable-component-extensions-with-background-pages",
		"--disable-default-apps",
		"--disable-dev-shm-usage",
		"--disable-extensions",
		"--disable-hang-monitor",
		"--disable-ipc-flooding-protection",
		"--disable-popup-blocking",
		"--disable-prompt-on-repost",
		"--disable-renderer-backgrounding",
		"--disable-sync",
		"--disable-translate",

		// Performance and stability
		"--metrics-recording-only",
		"--safebrowsing-disable-auto-update",

		// Enable automation features
		"--enable-features=NetworkService,NetworkServiceInProcess",
		"--force-color-profile=srgb",

		// Start with blank page
		"about:blank",
	}

	// Headless mode
	if opts.Headless {
		args = append(args, "--headless=new")
	}

	// Platform-specific flags
	if runtime.GOOS == "linux" {
		args = append(args, "--disable-gpu", "--disable-software-rasterizer")
	}

	// Additional user-specified args
	if len(opts.Args) > 0 {
		args = append(args, opts.Args...)
	}

	return args
}

// parseWebSocketURL reads Chrome's stderr to find the WebSocket URL.
func parseWebSocketURL(ctx context.Context, stderr io.Reader) (string, error) {
	scanner := bufio.NewScanner(stderr)
	urlRegex := regexp.MustCompile(`DevTools listening on (ws://[^\s]+)`)

	resultCh := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			if matches := urlRegex.FindStringSubmatch(line); len(matches) >= 2 {
				resultCh <- matches[1]
				return
			}
		}
		if err := scanner.Err(); err != nil {
			errCh <- err
		} else {
			errCh <- fmt.Errorf("Chrome exited without providing WebSocket URL")
		}
	}()

	select {
	case url := <-resultCh:
		return url, nil
	case err := <-errCh:
		return "", err
	case <-time.After(30 * time.Second):
		return "", fmt.Errorf("timeout waiting for Chrome to start")
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// findFreePort finds an available TCP port.
func findFreePort() (int, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port, nil
}

// MustLaunch is like Launch but panics on error.
func MustLaunch(ctx context.Context, opts *Options) *Browser {
	browser, err := Launch(ctx, opts)
	if err != nil {
		panic(err)
	}
	return browser
}

// LaunchHeadless is a convenience function to launch a headless browser.
func LaunchHeadless(ctx context.Context) (*Browser, error) {
	return Launch(ctx, &Options{Headless: true})
}

// Version returns the installed Chrome version, or empty string if not installed.
func Version() string {
	chromePath, err := GetChromePath("")
	if err != nil {
		return ""
	}
	if _, err := os.Stat(chromePath); err != nil {
		return ""
	}
	return extractVersionFromPath(chromePath)
}

// Uninstall removes the installed Chrome for Testing.
func Uninstall() error {
	chromeDir, err := GetChromeDir()
	if err != nil {
		return err
	}
	return os.RemoveAll(chromeDir)
}
