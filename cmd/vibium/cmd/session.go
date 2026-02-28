package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	vibium "github.com/agentplexus/vibium-go"
)

// SessionInfo stores information about a running browser session
type SessionInfo struct {
	WebSocketURL string `json:"websocket_url"`
	Headless     bool   `json:"headless"`
	PID          int    `json:"pid,omitempty"`
}

// saveSession saves session info to disk
func saveSession(info *SessionInfo) error {
	path := getSessionPath()

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create session directory: %w", err)
	}

	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	return nil
}

// loadSession loads session info from disk
//
//nolint:unused // scaffolding for future session reconnection feature
func loadSession() (*SessionInfo, error) {
	path := getSessionPath()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no active session (run 'vibium launch' first)")
		}
		return nil, fmt.Errorf("failed to read session file: %w", err)
	}

	var info SessionInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("failed to parse session file: %w", err)
	}

	return &info, nil
}

// clearSession removes the session file
func clearSession() error {
	path := getSessionPath()
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove session file: %w", err)
	}
	return nil
}

// Global vibe instance for the session
var globalVibe *vibium.Vibe

// getVibe returns a connected Vibe instance, launching if necessary
//
//nolint:unused // scaffolding for future session reconnection feature
func getVibe(_ context.Context) (*vibium.Vibe, error) {
	if globalVibe != nil && !globalVibe.IsClosed() {
		return globalVibe, nil
	}

	// Try to load existing session
	_, err := loadSession()
	if err != nil {
		return nil, err
	}

	// For now, we can't reconnect to an existing session
	// The clicker process must be running from 'launch' command
	return nil, fmt.Errorf("session exists but cannot reconnect (browser may have closed)")
}

// launchBrowser launches a new browser and saves the session
func launchBrowser(ctx context.Context, headless bool) (*vibium.Vibe, error) {
	opts := &vibium.LaunchOptions{
		Headless: headless,
	}

	vibe, err := vibium.Browser.Launch(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	// Save session info
	info := &SessionInfo{
		Headless: headless,
	}
	if err := saveSession(info); err != nil {
		// Non-fatal, just warn
		fmt.Fprintf(os.Stderr, "Warning: could not save session: %v\n", err)
	}

	globalVibe = vibe
	return vibe, nil
}

// quitBrowser closes the browser and clears the session
func quitBrowser(ctx context.Context) error {
	if globalVibe != nil {
		if err := globalVibe.Quit(ctx); err != nil {
			return fmt.Errorf("failed to quit browser: %w", err)
		}
		globalVibe = nil
	}

	if err := clearSession(); err != nil {
		return err
	}

	return nil
}

// mustGetVibe returns a Vibe or exits with an error
func mustGetVibe(_ context.Context) *vibium.Vibe {
	if globalVibe != nil && !globalVibe.IsClosed() {
		return globalVibe
	}
	fmt.Fprintln(os.Stderr, "Error: no active browser session (run 'vibium launch' first)")
	os.Exit(1)
	return nil
}
