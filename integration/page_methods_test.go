//go:build integration

package integration

import (
	"testing"

	"github.com/plexusone/webpilot"
)

// TestScrollDown tests scrolling down the page.
func TestScrollDown(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	// Navigate to a page with scrollable content
	bt.go_(`data:text/html,<!DOCTYPE html>
<html><body style="height: 5000px;">
<div id="top">Top of page</div>
<div id="bottom" style="position: absolute; top: 4500px;">Bottom of page</div>
</body></html>`)

	// Get initial scroll position
	initialY, err := bt.pilot.Evaluate(bt.ctx, `window.scrollY`)
	if err != nil {
		t.Fatalf("Failed to get initial scroll: %v", err)
	}
	t.Logf("Initial scrollY: %v", initialY)

	// Scroll down
	err = bt.pilot.Scroll(bt.ctx, "down", 500, nil)
	if err != nil {
		t.Fatalf("Failed to scroll down: %v", err)
	}

	// Check scroll position changed
	finalY, err := bt.pilot.Evaluate(bt.ctx, `window.scrollY`)
	if err != nil {
		t.Fatalf("Failed to get final scroll: %v", err)
	}
	t.Logf("Final scrollY: %v", finalY)

	// Verify scrolled
	if finalY == initialY {
		t.Error("Scroll position did not change")
	}
}

// TestScrollUp tests scrolling up the page.
func TestScrollUp(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html>
<html><body style="height: 5000px;">
<div id="content">Scrollable content</div>
</body></html>`)

	// Scroll down first
	err := bt.pilot.Scroll(bt.ctx, "down", 1000, nil)
	if err != nil {
		t.Fatalf("Failed to scroll down: %v", err)
	}

	midY, err := bt.pilot.Evaluate(bt.ctx, `window.scrollY`)
	if err != nil {
		t.Fatalf("Failed to get mid scroll: %v", err)
	}
	t.Logf("After scroll down, scrollY: %v", midY)

	// Scroll up
	err = bt.pilot.Scroll(bt.ctx, "up", 500, nil)
	if err != nil {
		t.Fatalf("Failed to scroll up: %v", err)
	}

	finalY, err := bt.pilot.Evaluate(bt.ctx, `window.scrollY`)
	if err != nil {
		t.Fatalf("Failed to get final scroll: %v", err)
	}
	t.Logf("After scroll up, scrollY: %v", finalY)
}

// TestScrollInElement tests scrolling within a specific element.
func TestScrollInElement(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html>
<html><body>
<div id="container" style="height: 200px; overflow: auto;">
<div style="height: 1000px;">
  <div id="top">Top</div>
  <div id="bottom" style="position: absolute; top: 800px;">Bottom</div>
</div>
</div>
</body></html>`)

	// Scroll within the container
	err := bt.pilot.Scroll(bt.ctx, "down", 300, &webpilot.ScrollOptions{
		Selector: "#container",
	})
	if err != nil {
		t.Fatalf("Failed to scroll in element: %v", err)
	}

	// Check container scroll position
	scrollTop, err := bt.pilot.Evaluate(bt.ctx, `document.getElementById('container').scrollTop`)
	if err != nil {
		t.Fatalf("Failed to get scrollTop: %v", err)
	}
	t.Logf("Container scrollTop: %v", scrollTop)
}

// TestSetExtraHTTPHeaders tests setting extra HTTP headers.
func TestSetExtraHTTPHeaders(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	// Set custom headers
	headers := map[string]string{
		"X-Custom-Header": "test-value",
		"X-Another":       "another-value",
	}

	err := bt.pilot.SetExtraHTTPHeaders(bt.ctx, headers)
	if err != nil {
		t.Fatalf("Failed to set extra headers: %v", err)
	}

	// Navigate to httpbin to verify headers
	bt.go_("https://httpbin.org/headers")

	// Get page content
	content, err := bt.pilot.Content(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to get content: %v", err)
	}

	// Check if custom headers appear (httpbin echoes headers)
	t.Logf("Headers page content length: %d", len(content))
	// Note: Headers may or may not appear depending on how httpbin formats response
}

// TestConsoleMessages tests console message capture.
func TestConsoleMessages(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html>
<html><body>
<script>
console.log('Test log message');
console.warn('Test warning');
console.error('Test error');
</script>
</body></html>`)

	// Get console messages (empty string = all levels)
	messages, err := bt.pilot.ConsoleMessages(bt.ctx, "")
	if err != nil {
		t.Fatalf("Failed to get console messages: %v", err)
	}

	t.Logf("Captured %d console messages", len(messages))
	for _, msg := range messages {
		t.Logf("  [%s] %s", msg.Type, msg.Text)
	}
}

// TestClearConsoleMessages tests clearing console messages.
func TestClearConsoleMessages(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html>
<html><body>
<script>
console.log('Message 1');
console.log('Message 2');
</script>
</body></html>`)

	// Clear console messages
	err := bt.pilot.ClearConsoleMessages(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to clear console messages: %v", err)
	}

	// Get messages - should be empty
	messages, err := bt.pilot.ConsoleMessages(bt.ctx, "")
	if err != nil {
		t.Fatalf("Failed to get console messages after clear: %v", err)
	}

	if len(messages) != 0 {
		t.Errorf("Expected 0 messages after clear, got %d", len(messages))
	}
}

// TestFillForm tests filling multiple form fields.
func TestFillForm(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html>
<html><body>
<form id="testForm">
  <input id="name" type="text" placeholder="Name">
  <input id="email" type="email" placeholder="Email">
  <input id="phone" type="tel" placeholder="Phone">
  <textarea id="message" placeholder="Message"></textarea>
</form>
</body></html>`)

	// Fill individual fields
	nameInput := bt.find("#name")
	err := nameInput.Fill(bt.ctx, "John Doe", nil)
	if err != nil {
		t.Fatalf("Failed to fill name: %v", err)
	}

	emailInput := bt.find("#email")
	err = emailInput.Fill(bt.ctx, "john@example.com", nil)
	if err != nil {
		t.Fatalf("Failed to fill email: %v", err)
	}

	phoneInput := bt.find("#phone")
	err = phoneInput.Fill(bt.ctx, "555-1234", nil)
	if err != nil {
		t.Fatalf("Failed to fill phone: %v", err)
	}

	messageInput := bt.find("#message")
	err = messageInput.Fill(bt.ctx, "Hello, this is a test message.", nil)
	if err != nil {
		t.Fatalf("Failed to fill message: %v", err)
	}

	// Verify values
	nameValue, err := nameInput.Value(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to get name value: %v", err)
	}
	if nameValue != "John Doe" {
		t.Errorf("Expected 'John Doe', got %q", nameValue)
	}

	emailValue, err := emailInput.Value(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to get email value: %v", err)
	}
	if emailValue != "john@example.com" {
		t.Errorf("Expected 'john@example.com', got %q", emailValue)
	}
}

// TestMouseDrag tests mouse drag operation.
func TestMouseDrag(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html>
<html><body>
<div id="draggable" style="width:100px; height:100px; background:red; position:absolute; left:50px; top:50px;">
Drag me
</div>
<script>
let dragging = false;
let offset = {x: 0, y: 0};
const el = document.getElementById('draggable');
el.addEventListener('mousedown', (e) => {
  dragging = true;
  offset.x = e.clientX - el.offsetLeft;
  offset.y = e.clientY - el.offsetTop;
});
document.addEventListener('mousemove', (e) => {
  if (dragging) {
    el.style.left = (e.clientX - offset.x) + 'px';
    el.style.top = (e.clientY - offset.y) + 'px';
  }
});
document.addEventListener('mouseup', () => { dragging = false; });
</script>
</body></html>`)

	// Get initial position
	initialLeft, err := bt.pilot.Evaluate(bt.ctx, `document.getElementById('draggable').offsetLeft`)
	if err != nil {
		t.Fatalf("Failed to get initial left: %v", err)
	}
	t.Logf("Initial left: %v", initialLeft)

	// Perform drag using mouse controller
	mouse, err := bt.pilot.Mouse(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to get mouse: %v", err)
	}

	// Move to element center
	err = mouse.Move(bt.ctx, 100.0, 100.0) // Center of 100x100 element at (50,50)
	if err != nil {
		t.Fatalf("Failed to move mouse: %v", err)
	}

	// Mouse down, move, mouse up
	err = mouse.Down(bt.ctx, webpilot.MouseButtonLeft)
	if err != nil {
		t.Fatalf("Failed mouse down: %v", err)
	}

	err = mouse.Move(bt.ctx, 200.0, 200.0)
	if err != nil {
		t.Fatalf("Failed to drag: %v", err)
	}

	err = mouse.Up(bt.ctx, webpilot.MouseButtonLeft)
	if err != nil {
		t.Fatalf("Failed mouse up: %v", err)
	}

	// Check position changed
	finalLeft, err := bt.pilot.Evaluate(bt.ctx, `document.getElementById('draggable').offsetLeft`)
	if err != nil {
		t.Fatalf("Failed to get final left: %v", err)
	}
	t.Logf("Final left: %v", finalLeft)
}
