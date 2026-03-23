//go:build integration

package integration

import (
	"testing"
)

// TestExampleCom tests basic functionality against example.com.
// This is a simple, stable site good for smoke tests.
func TestExampleCom(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_("https://example.com")

	t.Run("PageTitle", func(t *testing.T) {
		title, err := bt.pilot.Title(bt.ctx)
		if err != nil {
			t.Fatalf("Failed to get title: %v", err)
		}
		assertContains(t, title, "Example Domain")
	})

	t.Run("FindHeading", func(t *testing.T) {
		h1 := bt.find("h1")
		info := h1.Info()

		if info.Tag != "h1" {
			t.Errorf("Expected tag 'h1', got %q", info.Tag)
		}
		assertContains(t, info.Text, "Example Domain")
	})

	t.Run("FindLink", func(t *testing.T) {
		link := bt.find("a")
		info := link.Info()

		if info.Tag != "a" {
			t.Errorf("Expected tag 'a', got %q", info.Tag)
		}

		href, err := link.GetAttribute(bt.ctx, "href")
		if err != nil {
			t.Fatalf("Failed to get href: %v", err)
		}
		assertContains(t, href, "iana.org")
	})

	t.Run("GetElementText", func(t *testing.T) {
		link := bt.find("a")
		text, err := link.Text(bt.ctx)
		if err != nil {
			t.Fatalf("Failed to get text: %v", err)
		}
		// example.com link text may vary, just check it's not empty
		if text == "" {
			t.Error("Expected non-empty link text")
		}
	})

	t.Run("Screenshot", func(t *testing.T) {
		data := bt.screenshot()
		if len(data) < 1000 {
			t.Errorf("Screenshot too small: %d bytes", len(data))
		}
	})

	t.Run("BoundingBox", func(t *testing.T) {
		h1 := bt.find("h1")
		box, err := h1.BoundingBox(bt.ctx)
		if err != nil {
			t.Fatalf("Failed to get bounding box: %v", err)
		}

		if box.Width <= 0 || box.Height <= 0 {
			t.Errorf("Expected positive dimensions, got %+v", box)
		}
	})

	t.Run("ElementCenter", func(t *testing.T) {
		h1 := bt.find("h1")
		x, y := h1.Center()

		if x <= 0 || y <= 0 {
			t.Errorf("Expected positive center coordinates, got (%f, %f)", x, y)
		}
	})

	t.Run("EvaluateJS", func(t *testing.T) {
		result := bt.evaluate("return document.querySelectorAll('p').length")
		count, ok := result.(float64)
		if !ok {
			t.Fatalf("Expected float64, got %T", result)
		}
		if count < 1 {
			t.Errorf("Expected at least 1 paragraph, got %v", count)
		}
	})

	t.Run("FindAllParagraphs", func(t *testing.T) {
		paragraphs := bt.findAll("p")
		if len(paragraphs) < 1 {
			t.Errorf("Expected at least 1 paragraph, got %d", len(paragraphs))
		}

		for i, p := range paragraphs {
			info := p.Info()
			if info.Tag != "p" {
				t.Errorf("Paragraph %d: expected tag 'p', got %q", i, info.Tag)
			}
		}
	})
}

// TestExampleComClickLink tests clicking a link on example.com.
// Note: This navigates away from example.com, so it's a separate test.
func TestExampleComClickLink(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_("https://example.com")

	link := bt.find("a")
	if err := link.Click(bt.ctx, nil); err != nil {
		t.Fatalf("Failed to click link: %v", err)
	}

	// After click, should navigate to IANA
	url, err := bt.pilot.URL(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to get URL: %v", err)
	}
	assertContains(t, url, "iana.org")
}
