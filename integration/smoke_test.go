//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/plexusone/w3pilot"
)

// TestSmoke_BasicProtocol is a minimal test that verifies the vibium:* command
// protocol works end-to-end. If this test fails, the basic communication
// between w3pilot and the vibium clicker is broken.
//
// This test specifically exercises:
// - Browser launch via pipe transport
// - Navigation (standard BiDi command)
// - Element finding (vibium:page.find command)
// - Element info retrieval
func TestSmoke_BasicProtocol(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Launch browser - this tests pipe client startup and lifecycle.ready
	pilot, err := w3pilot.Browser.Launch(ctx, &w3pilot.LaunchOptions{
		Headless: true,
	})
	if err != nil {
		t.Fatalf("Failed to launch browser: %v", err)
	}
	defer func() {
		if err := pilot.Quit(ctx); err != nil {
			t.Logf("Warning: failed to quit browser: %v", err)
		}
	}()

	// Navigate to a data URL - tests browsingContext.navigate
	html := `<html><head><title>Smoke Test</title></head><body><h1 id="heading">Hello World</h1></body></html>`
	if err := pilot.Go(ctx, "data:text/html,"+html); err != nil {
		t.Fatalf("Failed to navigate: %v", err)
	}

	// Find element - this is the critical test for vibium:page.find
	elem, err := pilot.Find(ctx, "#heading", nil)
	if err != nil {
		t.Fatalf("Failed to find element with vibium:page.find: %v", err)
	}

	// Verify element info
	info := elem.Info()
	if info.Tag != "h1" && info.Tag != "H1" {
		t.Errorf("Expected tag 'h1' or 'H1', got %q", info.Tag)
	}
	if info.Text != "Hello World" {
		t.Errorf("Expected text 'Hello World', got %q", info.Text)
	}
}

// TestSmoke_ElementInteraction tests that element actions work.
// This verifies vibium:element.click, vibium:element.type, and related commands.
func TestSmoke_ElementInteraction(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pilot, err := w3pilot.Browser.Launch(ctx, &w3pilot.LaunchOptions{
		Headless: true,
	})
	if err != nil {
		t.Fatalf("Failed to launch browser: %v", err)
	}
	defer func() { _ = pilot.Quit(ctx) }()

	// Create a page with an input field
	html := `<html><body>
		<input type="text" id="input" value="">
		<button id="btn" onclick="document.getElementById('input').value='clicked'">Click</button>
	</body></html>`
	if err := pilot.Go(ctx, "data:text/html,"+html); err != nil {
		t.Fatalf("Failed to navigate: %v", err)
	}

	// Test vibium:element.type
	input, err := pilot.Find(ctx, "#input", nil)
	if err != nil {
		t.Fatalf("Failed to find input: %v", err)
	}
	if err := input.Type(ctx, "hello", nil); err != nil {
		t.Fatalf("Failed to type (vibium:element.type): %v", err)
	}

	// Test vibium:element.click
	btn, err := pilot.Find(ctx, "#btn", nil)
	if err != nil {
		t.Fatalf("Failed to find button: %v", err)
	}
	if err := btn.Click(ctx, nil); err != nil {
		t.Fatalf("Failed to click (vibium:element.click): %v", err)
	}

	// Verify click worked by checking input value
	input2, err := pilot.Find(ctx, "#input", nil)
	if err != nil {
		t.Fatalf("Failed to re-find input: %v", err)
	}
	value, err := input2.Value(ctx)
	if err != nil {
		t.Fatalf("Failed to get value (vibium:element.value): %v", err)
	}
	if value != "clicked" {
		t.Errorf("Expected value 'clicked', got %q", value)
	}
}

// TestSmoke_PageContent tests page content methods.
// This verifies vibium:page.content and related commands.
func TestSmoke_PageContent(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pilot, err := w3pilot.Browser.Launch(ctx, &w3pilot.LaunchOptions{
		Headless: true,
	})
	if err != nil {
		t.Fatalf("Failed to launch browser: %v", err)
	}
	defer func() { _ = pilot.Quit(ctx) }()

	html := `<html><body><p>Test content</p></body></html>`
	if err := pilot.Go(ctx, "data:text/html,"+html); err != nil {
		t.Fatalf("Failed to navigate: %v", err)
	}

	// Test vibium:page.content
	content, err := pilot.Content(ctx)
	if err != nil {
		t.Fatalf("Failed to get content (vibium:page.content): %v", err)
	}
	if content == "" {
		t.Error("Expected non-empty content")
	}
	// Content should include our test text
	found := false
	for i := 0; i <= len(content)-len("Test content"); i++ {
		if content[i:i+len("Test content")] == "Test content" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Content should include 'Test content'")
	}
}

// TestSmoke_SemanticSelector tests semantic selectors.
// This verifies that FindOptions with Role, Text, etc. work.
func TestSmoke_SemanticSelector(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pilot, err := w3pilot.Browser.Launch(ctx, &w3pilot.LaunchOptions{
		Headless: true,
	})
	if err != nil {
		t.Fatalf("Failed to launch browser: %v", err)
	}
	defer func() { _ = pilot.Quit(ctx) }()

	// Note: href="#" cannot be used in data URLs because # is interpreted
	// as a fragment identifier, truncating the HTML content. Use a full URL instead.
	html := `<html><body>
		<button>Submit Form</button>
		<a href="https://example.com">Learn More</a>
	</body></html>`
	if err := pilot.Go(ctx, "data:text/html,"+html); err != nil {
		t.Fatalf("Failed to navigate: %v", err)
	}

	// Find by role (button)
	btn, err := pilot.Find(ctx, "", &w3pilot.FindOptions{
		Role: "button",
		Text: "Submit",
	})
	if err != nil {
		t.Fatalf("Failed to find button by role: %v", err)
	}
	info := btn.Info()
	if info.Tag != "button" && info.Tag != "BUTTON" {
		t.Errorf("Expected button tag, got %q", info.Tag)
	}

	// Find by role (link)
	link, err := pilot.Find(ctx, "", &w3pilot.FindOptions{
		Role: "link",
		Text: "Learn More",
	})
	if err != nil {
		t.Fatalf("Failed to find link: %v", err)
	}
	linkInfo := link.Info()
	if linkInfo.Tag != "a" && linkInfo.Tag != "A" {
		t.Errorf("Expected 'a' tag, got %q", linkInfo.Tag)
	}
}
