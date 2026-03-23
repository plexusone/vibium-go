//go:build integration

package integration

import (
	"testing"
	"time"
)

// TestTheInternet tests functionality against the-internet.herokuapp.com.
// This site provides various UI patterns for testing.
func TestTheInternet(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_("https://the-internet.herokuapp.com/")

	t.Run("PageTitle", func(t *testing.T) {
		title, err := bt.pilot.Title(bt.ctx)
		if err != nil {
			t.Fatalf("Failed to get title: %v", err)
		}
		assertContains(t, title, "The Internet")
	})

	t.Run("FindHeading", func(t *testing.T) {
		h1 := bt.find("h1.heading")
		info := h1.Info()

		if info.Tag != "h1" {
			t.Errorf("Expected tag 'h1', got %q", info.Tag)
		}
		assertContains(t, info.Text, "Welcome to the-internet")
	})

	t.Run("FindAllLinks", func(t *testing.T) {
		links := bt.findAll("ul li a")
		if len(links) < 10 {
			t.Errorf("Expected at least 10 links, got %d", len(links))
		}
	})

	t.Run("Screenshot", func(t *testing.T) {
		data := bt.screenshot()
		if len(data) < 5000 {
			t.Errorf("Screenshot too small: %d bytes", len(data))
		}
	})
}

// TestTheInternetAddRemoveElements tests dynamic element creation.
func TestTheInternetAddRemoveElements(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_("https://the-internet.herokuapp.com/add_remove_elements/")

	t.Run("AddElement", func(t *testing.T) {
		// Click "Add Element" button
		addBtn := bt.find("button[onclick='addElement()']")
		if err := addBtn.Click(bt.ctx, nil); err != nil {
			t.Fatalf("Failed to click Add Element: %v", err)
		}

		// Wait a moment for element to appear
		time.Sleep(100 * time.Millisecond)

		// Find the newly added Delete button
		deleteBtn := bt.find(".added-manually")
		info := deleteBtn.Info()
		assertContains(t, info.Text, "Delete")
	})

	t.Run("AddMultipleElements", func(t *testing.T) {
		// Add 3 more elements
		addBtn := bt.find("button[onclick='addElement()']")
		for i := 0; i < 3; i++ {
			if err := addBtn.Click(bt.ctx, nil); err != nil {
				t.Fatalf("Failed to click Add Element: %v", err)
			}
			time.Sleep(50 * time.Millisecond)
		}

		// Should have 4 delete buttons now (1 from previous test + 3 new)
		deleteBtns := bt.findAll(".added-manually")
		if len(deleteBtns) < 4 {
			t.Errorf("Expected at least 4 delete buttons, got %d", len(deleteBtns))
		}
	})
}

// TestTheInternetInputs tests form input handling.
func TestTheInternetInputs(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_("https://the-internet.herokuapp.com/inputs")

	t.Run("TypeNumber", func(t *testing.T) {
		input := bt.find("input[type='number']")

		// Type a number
		if err := input.Type(bt.ctx, "12345", nil); err != nil {
			t.Fatalf("Failed to type: %v", err)
		}

		// Verify the value
		result := bt.evaluate("return document.querySelector('input').value")
		value, ok := result.(string)
		if !ok {
			t.Fatalf("Expected string, got %T", result)
		}
		if value != "12345" {
			t.Errorf("Expected '12345', got %q", value)
		}
	})
}

// TestTheInternetCheckboxes tests checkbox interaction.
func TestTheInternetCheckboxes(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_("https://the-internet.herokuapp.com/checkboxes")

	t.Run("FindCheckboxes", func(t *testing.T) {
		checkboxes := bt.findAll("input[type='checkbox']")
		if len(checkboxes) != 2 {
			t.Errorf("Expected 2 checkboxes, got %d", len(checkboxes))
		}
	})

	t.Run("ClickCheckbox", func(t *testing.T) {
		// Get initial state of first checkbox
		initialState := bt.evaluate("return document.querySelectorAll('input[type=checkbox]')[0].checked")
		wasChecked, _ := initialState.(bool)

		// Click first checkbox
		checkbox := bt.find("input[type='checkbox']")
		if err := checkbox.Click(bt.ctx, nil); err != nil {
			t.Fatalf("Failed to click checkbox: %v", err)
		}

		// Verify state changed
		newState := bt.evaluate("return document.querySelectorAll('input[type=checkbox]')[0].checked")
		isChecked, _ := newState.(bool)

		if isChecked == wasChecked {
			t.Error("Checkbox state should have changed after click")
		}
	})
}

// TestTheInternetDropdown tests select element interaction.
func TestTheInternetDropdown(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_("https://the-internet.herokuapp.com/dropdown")

	t.Run("FindDropdown", func(t *testing.T) {
		dropdown := bt.find("#dropdown")
		info := dropdown.Info()

		if info.Tag != "select" {
			t.Errorf("Expected tag 'select', got %q", info.Tag)
		}
	})

	t.Run("FindOptions", func(t *testing.T) {
		options := bt.findAll("#dropdown option")
		if len(options) != 3 {
			t.Errorf("Expected 3 options, got %d", len(options))
		}
	})
}

// TestTheInternetNavigation tests browser navigation (back/forward).
func TestTheInternetNavigation(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	// Start at home page
	bt.go_("https://the-internet.herokuapp.com/")

	// Navigate to inputs page
	bt.go_("https://the-internet.herokuapp.com/inputs")

	url1, _ := bt.pilot.URL(bt.ctx)
	assertContains(t, url1, "inputs")

	t.Run("Back", func(t *testing.T) {
		if err := bt.pilot.Back(bt.ctx); err != nil {
			t.Fatalf("Failed to go back: %v", err)
		}

		// Wait for navigation
		time.Sleep(500 * time.Millisecond)

		url, _ := bt.pilot.URL(bt.ctx)
		if url == url1 {
			t.Error("URL should have changed after going back")
		}
	})

	t.Run("Forward", func(t *testing.T) {
		if err := bt.pilot.Forward(bt.ctx); err != nil {
			t.Fatalf("Failed to go forward: %v", err)
		}

		// Wait for navigation
		time.Sleep(500 * time.Millisecond)

		url, _ := bt.pilot.URL(bt.ctx)
		assertContains(t, url, "inputs")
	})

	t.Run("Reload", func(t *testing.T) {
		if err := bt.pilot.Reload(bt.ctx); err != nil {
			t.Fatalf("Failed to reload: %v", err)
		}

		// Page should still be inputs after reload
		url, _ := bt.pilot.URL(bt.ctx)
		assertContains(t, url, "inputs")
	})
}
