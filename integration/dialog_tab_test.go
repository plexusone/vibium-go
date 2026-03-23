//go:build integration

package integration

import (
	"testing"
	"time"
)

// TestDialogAlert tests handling alert dialogs.
func TestDialogAlert(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	// Navigate to a page that will show an alert
	bt.go_(`data:text/html,<!DOCTYPE html>
<html><body>
<button id="alertBtn" onclick="alert('Hello World!')">Show Alert</button>
</body></html>`)

	// Set up dialog handler before triggering
	go func() {
		// Give time for alert to appear
		time.Sleep(200 * time.Millisecond)
		err := bt.pilot.HandleDialog(bt.ctx, true, "")
		if err != nil {
			t.Logf("HandleDialog error (may be expected): %v", err)
		}
	}()

	// Click button to trigger alert
	elem := bt.find("#alertBtn")
	err := elem.Click(bt.ctx, nil)
	if err != nil {
		t.Fatalf("Failed to click alert button: %v", err)
	}

	// Wait for dialog handling
	time.Sleep(300 * time.Millisecond)
}

// TestDialogConfirm tests handling confirm dialogs.
func TestDialogConfirm(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html>
<html><body>
<button id="confirmBtn" onclick="document.body.innerHTML = confirm('Proceed?') ? 'yes' : 'no'">Confirm</button>
</body></html>`)

	// Accept the confirm dialog
	go func() {
		time.Sleep(200 * time.Millisecond)
		_ = bt.pilot.HandleDialog(bt.ctx, true, "")
	}()

	elem := bt.find("#confirmBtn")
	err := elem.Click(bt.ctx, nil)
	if err != nil {
		t.Fatalf("Failed to click confirm button: %v", err)
	}

	time.Sleep(300 * time.Millisecond)

	// Check result
	content, err := bt.pilot.Content(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to get content: %v", err)
	}
	// The confirm was accepted, so body should contain "yes"
	t.Logf("Content after confirm: %s", content)
}

// TestDialogPrompt tests handling prompt dialogs.
func TestDialogPrompt(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html>
<html><body>
<button id="promptBtn" onclick="document.body.innerHTML = prompt('Name?') || 'cancelled'">Prompt</button>
</body></html>`)

	// Accept with text
	go func() {
		time.Sleep(200 * time.Millisecond)
		_ = bt.pilot.HandleDialog(bt.ctx, true, "TestUser")
	}()

	elem := bt.find("#promptBtn")
	err := elem.Click(bt.ctx, nil)
	if err != nil {
		t.Fatalf("Failed to click prompt button: %v", err)
	}

	time.Sleep(300 * time.Millisecond)

	content, err := bt.pilot.Content(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to get content: %v", err)
	}
	t.Logf("Content after prompt: %s", content)
}

// TestNewPage tests creating a new page/tab.
func TestNewPage(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	// Navigate initial page
	bt.go_("https://example.com")

	// Create new page
	newPage, err := bt.pilot.NewPage(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to create new page: %v", err)
	}

	// Navigate new page to different URL
	err = newPage.Go(bt.ctx, "https://httpbin.org/html")
	if err != nil {
		t.Fatalf("Failed to navigate new page: %v", err)
	}

	// Check new page URL
	url, err := newPage.URL(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to get new page URL: %v", err)
	}
	t.Logf("New page URL: %s", url)

	// Close the new page
	err = newPage.Close(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to close new page: %v", err)
	}
}

// TestPages tests listing all pages.
func TestPages(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	// Navigate initial page
	bt.go_("https://example.com")

	// Get pages
	pages, err := bt.pilot.Pages(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to get pages: %v", err)
	}

	// Should have at least one page
	if len(pages) == 0 {
		t.Error("Expected at least one page")
	}
	t.Logf("Found %d page(s)", len(pages))

	// Create another page
	newPage, err := bt.pilot.NewPage(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to create new page: %v", err)
	}
	defer newPage.Close(bt.ctx)

	// Get pages again
	pages, err = bt.pilot.Pages(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to get pages after creating new: %v", err)
	}

	if len(pages) < 2 {
		t.Errorf("Expected at least 2 pages, got %d", len(pages))
	}
	t.Logf("Found %d page(s) after creating new", len(pages))
}

// TestBringToFront tests bringing a page to front.
func TestBringToFront(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_("https://example.com")

	// Create and navigate new page
	newPage, err := bt.pilot.NewPage(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to create new page: %v", err)
	}
	defer newPage.Close(bt.ctx)

	err = newPage.Go(bt.ctx, "https://httpbin.org/html")
	if err != nil {
		t.Fatalf("Failed to navigate new page: %v", err)
	}

	// Bring original page to front
	err = bt.pilot.BringToFront(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to bring to front: %v", err)
	}
	t.Log("Brought original page to front")
}

// TestClosePage tests closing a page.
func TestClosePage(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_("https://example.com")

	// Create a new page
	newPage, err := bt.pilot.NewPage(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to create new page: %v", err)
	}

	// Get page count before close
	pagesBefore, err := bt.pilot.Pages(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to get pages: %v", err)
	}

	// Close the new page
	err = newPage.Close(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to close page: %v", err)
	}

	// Give time for close to complete
	time.Sleep(100 * time.Millisecond)

	// Get page count after close
	pagesAfter, err := bt.pilot.Pages(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to get pages after close: %v", err)
	}

	if len(pagesAfter) >= len(pagesBefore) {
		t.Errorf("Expected fewer pages after close: before=%d, after=%d", len(pagesBefore), len(pagesAfter))
	}
}
