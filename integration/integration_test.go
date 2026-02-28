//go:build integration

// Package integration contains integration tests that run against live websites.
// Run with: go test -tags=integration -v ./integration/...
package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/plexusone/vibium-go"
)

// testTimeout is the default timeout for test operations.
const testTimeout = 60 * time.Second

// TestMain sets up the test environment.
func TestMain(m *testing.M) {
	// Could add global setup here (e.g., check clicker binary exists)
	os.Exit(m.Run())
}

// browserTest is a helper for running browser tests.
type browserTest struct {
	t       *testing.T
	vibe    *vibium.Vibe
	ctx     context.Context
	cancel  context.CancelFunc
	headless bool
}

// newBrowserTest creates a new browser test helper.
func newBrowserTest(t *testing.T) *browserTest {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)

	// Use headless mode in CI, visible mode locally for debugging
	headless := os.Getenv("CI") != "" || os.Getenv("VIBIUM_HEADLESS") == "1"

	vibe, err := vibium.Browser.Launch(ctx, &vibium.LaunchOptions{
		Headless: headless,
	})
	if err != nil {
		cancel()
		t.Fatalf("Failed to launch browser: %v", err)
	}

	return &browserTest{
		t:        t,
		vibe:     vibe,
		ctx:      ctx,
		cancel:   cancel,
		headless: headless,
	}
}

// cleanup closes the browser and cancels the context.
func (bt *browserTest) cleanup() {
	if bt.vibe != nil {
		if err := bt.vibe.Quit(bt.ctx); err != nil {
			bt.t.Logf("Warning: failed to quit browser: %v", err)
		}
	}
	bt.cancel()
}

// go navigates to a URL.
func (bt *browserTest) go_(url string) {
	bt.t.Helper()
	if err := bt.vibe.Go(bt.ctx, url); err != nil {
		bt.t.Fatalf("Failed to navigate to %s: %v", url, err)
	}
}

// find finds an element by selector.
func (bt *browserTest) find(selector string) *vibium.Element {
	bt.t.Helper()
	elem, err := bt.vibe.Find(bt.ctx, selector, nil)
	if err != nil {
		bt.t.Fatalf("Failed to find element %q: %v", selector, err)
	}
	return elem
}

// findAll finds all elements matching the selector.
func (bt *browserTest) findAll(selector string) []*vibium.Element {
	bt.t.Helper()
	elements, err := bt.vibe.FindAll(bt.ctx, selector)
	if err != nil {
		bt.t.Fatalf("Failed to find elements %q: %v", selector, err)
	}
	return elements
}

// screenshot takes a screenshot and verifies it's valid PNG.
func (bt *browserTest) screenshot() []byte {
	bt.t.Helper()
	data, err := bt.vibe.Screenshot(bt.ctx)
	if err != nil {
		bt.t.Fatalf("Failed to take screenshot: %v", err)
	}

	// Verify PNG magic bytes
	if len(data) < 8 {
		bt.t.Fatalf("Screenshot too small: %d bytes", len(data))
	}
	if data[0] != 0x89 || data[1] != 0x50 || data[2] != 0x4E || data[3] != 0x47 {
		bt.t.Fatal("Screenshot is not valid PNG")
	}

	return data
}

// evaluate executes JavaScript and returns the result.
func (bt *browserTest) evaluate(script string) interface{} {
	bt.t.Helper()
	result, err := bt.vibe.Evaluate(bt.ctx, script)
	if err != nil {
		bt.t.Fatalf("Failed to evaluate script: %v", err)
	}
	return result
}

// assertContains checks that a string contains a substring.
func assertContains(t *testing.T, s, substr string) {
	t.Helper()
	if len(s) == 0 {
		t.Errorf("Expected string containing %q, got empty string", substr)
		return
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return
		}
	}
	t.Errorf("Expected string containing %q, got %q", substr, s)
}
