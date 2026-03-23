//go:build integration

package integration

import (
	"testing"
	"time"

	"github.com/plexusone/webpilot"
)

// testPageSemanticSelectors is an HTML page for testing semantic selectors.
// It contains various elements with semantic attributes for comprehensive testing.
const testPageSemanticSelectors = `data:text/html,<!DOCTYPE html>
<html>
<head>
    <title>Semantic Selectors Test Page</title>
</head>
<body>
    <h1>Semantic Selectors Test</h1>

    <!-- Buttons with roles and text -->
    <button id="submit-btn" data-testid="submit-button">Submit Form</button>
    <button id="cancel-btn" data-testid="cancel-button" title="Cancel the operation">Cancel</button>
    <button id="save-btn" aria-label="Save changes">Save</button>

    <!-- Form with labeled inputs -->
    <form id="test-form">
        <div>
            <label for="email-input">Email Address</label>
            <input type="email" id="email-input" placeholder="Enter your email" data-testid="email-field">
        </div>
        <div>
            <label for="password-input">Password</label>
            <input type="password" id="password-input" placeholder="Enter password" data-testid="password-field">
        </div>
        <div>
            <label for="name-input">Full Name</label>
            <input type="text" id="name-input" placeholder="Your full name">
        </div>
        <div>
            <label>
                <input type="checkbox" id="agree-checkbox" data-testid="agree-checkbox">
                I agree to the terms
            </label>
        </div>
    </form>

    <!-- Links -->
    <nav>
        <a href="#home" data-testid="nav-home">Home</a>
        <a href="#about" title="Learn about us">About Us</a>
        <a href="#contact">Contact</a>
    </nav>

    <!-- Images with alt text -->
    <div id="images">
        <img src="data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"
             alt="Company Logo" data-testid="logo-image">
        <img src="data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"
             alt="User Avatar" title="Your profile picture">
    </div>

    <!-- Elements with various roles -->
    <div role="alert" id="alert-div">This is an alert message</div>
    <div role="dialog" id="dialog-div" aria-label="Confirmation Dialog">
        <p>Are you sure?</p>
        <button>Yes</button>
        <button>No</button>
    </div>

    <!-- Nested structure for scoped search -->
    <div id="section-a" data-testid="section-a">
        <h2>Section A</h2>
        <button data-testid="section-btn">Section A Button</button>
        <input type="text" placeholder="Section A Input">
    </div>
    <div id="section-b" data-testid="section-b">
        <h2>Section B</h2>
        <button data-testid="section-btn">Section B Button</button>
        <input type="text" placeholder="Section B Input">
    </div>

    <!-- Multiple similar elements for FindAll testing -->
    <ul id="item-list">
        <li data-testid="list-item">Item 1</li>
        <li data-testid="list-item">Item 2</li>
        <li data-testid="list-item">Item 3</li>
    </ul>

    <!-- Element near another element -->
    <div id="proximity-test">
        <span id="username-label">Username:</span>
        <input type="text" id="username-input" placeholder="Enter username">
        <button id="username-submit">Submit</button>
    </div>
</body>
</html>`

// TestSemanticSelectorsByRole tests finding elements by ARIA role.
func TestSemanticSelectorsByRole(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(testPageSemanticSelectors)
	time.Sleep(200 * time.Millisecond) // Allow page to fully render

	t.Run("FindButtonByRole", func(t *testing.T) {
		elem, err := bt.pilot.Find(bt.ctx, "", &webpilot.FindOptions{
			Role: "button",
			Text: "Submit Form",
		})
		if err != nil {
			t.Fatalf("Failed to find button by role: %v", err)
		}

		info := elem.Info()
		if info.Tag != "button" {
			t.Errorf("Expected tag 'button', got %q", info.Tag)
		}
		assertContains(t, info.Text, "Submit Form")
	})

	t.Run("FindAlertByRole", func(t *testing.T) {
		elem, err := bt.pilot.Find(bt.ctx, "", &webpilot.FindOptions{
			Role: "alert",
		})
		if err != nil {
			t.Fatalf("Failed to find alert by role: %v", err)
		}

		info := elem.Info()
		assertContains(t, info.Text, "alert message")
	})

	t.Run("FindDialogByRole", func(t *testing.T) {
		elem, err := bt.pilot.Find(bt.ctx, "", &webpilot.FindOptions{
			Role: "dialog",
		})
		if err != nil {
			t.Fatalf("Failed to find dialog by role: %v", err)
		}

		info := elem.Info()
		assertContains(t, info.Text, "Are you sure")
	})

	t.Run("FindAllButtonsByRole", func(t *testing.T) {
		elems, err := bt.pilot.FindAll(bt.ctx, "", &webpilot.FindOptions{
			Role: "button",
		})
		if err != nil {
			t.Fatalf("Failed to find all buttons by role: %v", err)
		}

		// Should find multiple buttons
		if len(elems) < 5 {
			t.Errorf("Expected at least 5 buttons, got %d", len(elems))
		}
	})
}

// TestSemanticSelectorsByText tests finding elements by text content.
func TestSemanticSelectorsByText(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(testPageSemanticSelectors)
	time.Sleep(200 * time.Millisecond)

	t.Run("FindByExactText", func(t *testing.T) {
		elem, err := bt.pilot.Find(bt.ctx, "", &webpilot.FindOptions{
			Text: "Cancel",
		})
		if err != nil {
			t.Fatalf("Failed to find element by text: %v", err)
		}

		info := elem.Info()
		if info.Tag != "button" {
			t.Errorf("Expected tag 'button', got %q", info.Tag)
		}
	})

	t.Run("FindByPartialText", func(t *testing.T) {
		elem, err := bt.pilot.Find(bt.ctx, "", &webpilot.FindOptions{
			Text: "Section A",
		})
		if err != nil {
			t.Fatalf("Failed to find element by partial text: %v", err)
		}

		info := elem.Info()
		assertContains(t, info.Text, "Section A")
	})

	t.Run("FindLinkByText", func(t *testing.T) {
		elem, err := bt.pilot.Find(bt.ctx, "", &webpilot.FindOptions{
			Role: "link",
			Text: "Home",
		})
		if err != nil {
			t.Fatalf("Failed to find link by text: %v", err)
		}

		info := elem.Info()
		if info.Tag != "a" {
			t.Errorf("Expected tag 'a', got %q", info.Tag)
		}
	})
}

// TestSemanticSelectorsByLabel tests finding elements by associated label.
func TestSemanticSelectorsByLabel(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(testPageSemanticSelectors)
	time.Sleep(200 * time.Millisecond)

	t.Run("FindInputByLabel", func(t *testing.T) {
		elem, err := bt.pilot.Find(bt.ctx, "", &webpilot.FindOptions{
			Label: "Email Address",
		})
		if err != nil {
			t.Fatalf("Failed to find input by label: %v", err)
		}

		info := elem.Info()
		if info.Tag != "input" {
			t.Errorf("Expected tag 'input', got %q", info.Tag)
		}
	})

	t.Run("FindPasswordByLabel", func(t *testing.T) {
		elem, err := bt.pilot.Find(bt.ctx, "", &webpilot.FindOptions{
			Label: "Password",
		})
		if err != nil {
			t.Fatalf("Failed to find password input by label: %v", err)
		}

		// Type into the field to verify it's the correct element
		err = elem.Type(bt.ctx, "secret123", nil)
		if err != nil {
			t.Fatalf("Failed to type into password field: %v", err)
		}

		value, err := elem.Value(bt.ctx)
		if err != nil {
			t.Fatalf("Failed to get value: %v", err)
		}
		if value != "secret123" {
			t.Errorf("Expected 'secret123', got %q", value)
		}
	})
}

// TestSemanticSelectorsByPlaceholder tests finding elements by placeholder text.
func TestSemanticSelectorsByPlaceholder(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(testPageSemanticSelectors)
	time.Sleep(200 * time.Millisecond)

	t.Run("FindByPlaceholder", func(t *testing.T) {
		elem, err := bt.pilot.Find(bt.ctx, "", &webpilot.FindOptions{
			Placeholder: "Enter your email",
		})
		if err != nil {
			t.Fatalf("Failed to find by placeholder: %v", err)
		}

		info := elem.Info()
		if info.Tag != "input" {
			t.Errorf("Expected tag 'input', got %q", info.Tag)
		}
	})

	t.Run("FindByPartialPlaceholder", func(t *testing.T) {
		elem, err := bt.pilot.Find(bt.ctx, "", &webpilot.FindOptions{
			Placeholder: "password",
		})
		if err != nil {
			t.Fatalf("Failed to find by partial placeholder: %v", err)
		}

		// Fill and verify
		err = elem.Fill(bt.ctx, "mypassword", nil)
		if err != nil {
			t.Fatalf("Failed to fill: %v", err)
		}

		value, err := elem.Value(bt.ctx)
		if err != nil {
			t.Fatalf("Failed to get value: %v", err)
		}
		if value != "mypassword" {
			t.Errorf("Expected 'mypassword', got %q", value)
		}
	})
}

// TestSemanticSelectorsByTestID tests finding elements by data-testid attribute.
func TestSemanticSelectorsByTestID(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(testPageSemanticSelectors)
	time.Sleep(200 * time.Millisecond)

	t.Run("FindByTestID", func(t *testing.T) {
		elem, err := bt.pilot.Find(bt.ctx, "", &webpilot.FindOptions{
			TestID: "submit-button",
		})
		if err != nil {
			t.Fatalf("Failed to find by testid: %v", err)
		}

		info := elem.Info()
		assertContains(t, info.Text, "Submit Form")
	})

	t.Run("FindInputByTestID", func(t *testing.T) {
		elem, err := bt.pilot.Find(bt.ctx, "", &webpilot.FindOptions{
			TestID: "email-field",
		})
		if err != nil {
			t.Fatalf("Failed to find input by testid: %v", err)
		}

		info := elem.Info()
		if info.Tag != "input" {
			t.Errorf("Expected tag 'input', got %q", info.Tag)
		}
	})

	t.Run("FindAllByTestID", func(t *testing.T) {
		elems, err := bt.pilot.FindAll(bt.ctx, "", &webpilot.FindOptions{
			TestID: "list-item",
		})
		if err != nil {
			t.Fatalf("Failed to find all by testid: %v", err)
		}

		if len(elems) != 3 {
			t.Errorf("Expected 3 list items, got %d", len(elems))
		}
	})
}

// TestSemanticSelectorsByAlt tests finding images by alt text.
func TestSemanticSelectorsByAlt(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(testPageSemanticSelectors)
	time.Sleep(200 * time.Millisecond)

	t.Run("FindImageByAlt", func(t *testing.T) {
		elem, err := bt.pilot.Find(bt.ctx, "", &webpilot.FindOptions{
			Alt: "Company Logo",
		})
		if err != nil {
			t.Fatalf("Failed to find image by alt: %v", err)
		}

		info := elem.Info()
		if info.Tag != "img" {
			t.Errorf("Expected tag 'img', got %q", info.Tag)
		}
	})

	t.Run("FindImageByPartialAlt", func(t *testing.T) {
		elem, err := bt.pilot.Find(bt.ctx, "", &webpilot.FindOptions{
			Alt: "Avatar",
		})
		if err != nil {
			t.Fatalf("Failed to find image by partial alt: %v", err)
		}

		info := elem.Info()
		if info.Tag != "img" {
			t.Errorf("Expected tag 'img', got %q", info.Tag)
		}
	})
}

// TestSemanticSelectorsByTitle tests finding elements by title attribute.
func TestSemanticSelectorsByTitle(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(testPageSemanticSelectors)
	time.Sleep(200 * time.Millisecond)

	t.Run("FindByTitle", func(t *testing.T) {
		elem, err := bt.pilot.Find(bt.ctx, "", &webpilot.FindOptions{
			Title: "Cancel the operation",
		})
		if err != nil {
			t.Fatalf("Failed to find by title: %v", err)
		}

		info := elem.Info()
		assertContains(t, info.Text, "Cancel")
	})

	t.Run("FindLinkByTitle", func(t *testing.T) {
		elem, err := bt.pilot.Find(bt.ctx, "", &webpilot.FindOptions{
			Title: "Learn about us",
		})
		if err != nil {
			t.Fatalf("Failed to find link by title: %v", err)
		}

		info := elem.Info()
		if info.Tag != "a" {
			t.Errorf("Expected tag 'a', got %q", info.Tag)
		}
	})
}

// TestSemanticSelectorsByXPath tests finding elements by XPath.
func TestSemanticSelectorsByXPath(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(testPageSemanticSelectors)
	time.Sleep(200 * time.Millisecond)

	t.Run("FindByXPath", func(t *testing.T) {
		elem, err := bt.pilot.Find(bt.ctx, "", &webpilot.FindOptions{
			XPath: "//button[@id='submit-btn']",
		})
		if err != nil {
			t.Fatalf("Failed to find by xpath: %v", err)
		}

		info := elem.Info()
		assertContains(t, info.Text, "Submit Form")
	})

	t.Run("FindInputByXPath", func(t *testing.T) {
		elem, err := bt.pilot.Find(bt.ctx, "", &webpilot.FindOptions{
			XPath: "//input[@type='email']",
		})
		if err != nil {
			t.Fatalf("Failed to find input by xpath: %v", err)
		}

		info := elem.Info()
		if info.Tag != "input" {
			t.Errorf("Expected tag 'input', got %q", info.Tag)
		}
	})

	t.Run("FindAllByXPath", func(t *testing.T) {
		elems, err := bt.pilot.FindAll(bt.ctx, "", &webpilot.FindOptions{
			XPath: "//li[@data-testid='list-item']",
		})
		if err != nil {
			t.Fatalf("Failed to find all by xpath: %v", err)
		}

		if len(elems) != 3 {
			t.Errorf("Expected 3 items, got %d", len(elems))
		}
	})
}

// TestSemanticSelectorsCombined tests combining CSS selector with semantic options.
func TestSemanticSelectorsCombined(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(testPageSemanticSelectors)
	time.Sleep(200 * time.Millisecond)

	t.Run("CSSWithRole", func(t *testing.T) {
		// Find a button within the form
		elem, err := bt.pilot.Find(bt.ctx, "form", &webpilot.FindOptions{
			Role: "checkbox",
		})
		if err != nil {
			t.Fatalf("Failed to find with CSS + role: %v", err)
		}

		info := elem.Info()
		if info.Tag != "input" {
			t.Errorf("Expected tag 'input', got %q", info.Tag)
		}
	})

	t.Run("CSSWithLabel", func(t *testing.T) {
		// Find input with label within the form
		elem, err := bt.pilot.Find(bt.ctx, "#test-form", &webpilot.FindOptions{
			Label: "Full Name",
		})
		if err != nil {
			t.Fatalf("Failed to find with CSS + label: %v", err)
		}

		// Type to verify correct element
		err = elem.Fill(bt.ctx, "John Doe", nil)
		if err != nil {
			t.Fatalf("Failed to fill: %v", err)
		}

		value, err := elem.Value(bt.ctx)
		if err != nil {
			t.Fatalf("Failed to get value: %v", err)
		}
		if value != "John Doe" {
			t.Errorf("Expected 'John Doe', got %q", value)
		}
	})

	t.Run("CSSWithRoleAndText", func(t *testing.T) {
		// Find specific button within dialog
		elem, err := bt.pilot.Find(bt.ctx, "[role='dialog']", &webpilot.FindOptions{
			Role: "button",
			Text: "Yes",
		})
		if err != nil {
			t.Fatalf("Failed to find with CSS + role + text: %v", err)
		}

		info := elem.Info()
		assertContains(t, info.Text, "Yes")
	})
}

// TestSemanticSelectorsScoped tests Element.Find() for scoped searching.
func TestSemanticSelectorsScoped(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(testPageSemanticSelectors)
	time.Sleep(200 * time.Millisecond)

	t.Run("ScopedFindByTestID", func(t *testing.T) {
		// First find section A
		sectionA, err := bt.pilot.Find(bt.ctx, "", &webpilot.FindOptions{
			TestID: "section-a",
		})
		if err != nil {
			t.Fatalf("Failed to find section A: %v", err)
		}

		// Find button within section A
		btn, err := sectionA.Find(bt.ctx, "", &webpilot.FindOptions{
			TestID: "section-btn",
		})
		if err != nil {
			t.Fatalf("Failed to find button in section A: %v", err)
		}

		info := btn.Info()
		assertContains(t, info.Text, "Section A Button")
	})

	t.Run("ScopedFindByPlaceholder", func(t *testing.T) {
		// Find section B
		sectionB, err := bt.pilot.Find(bt.ctx, "#section-b", nil)
		if err != nil {
			t.Fatalf("Failed to find section B: %v", err)
		}

		// Find input within section B by placeholder
		input, err := sectionB.Find(bt.ctx, "", &webpilot.FindOptions{
			Placeholder: "Section B Input",
		})
		if err != nil {
			t.Fatalf("Failed to find input in section B: %v", err)
		}

		// Type to verify
		err = input.Fill(bt.ctx, "Section B Value", nil)
		if err != nil {
			t.Fatalf("Failed to fill: %v", err)
		}

		value, err := input.Value(bt.ctx)
		if err != nil {
			t.Fatalf("Failed to get value: %v", err)
		}
		if value != "Section B Value" {
			t.Errorf("Expected 'Section B Value', got %q", value)
		}
	})

	t.Run("ScopedFindAll", func(t *testing.T) {
		// Find the item list
		list, err := bt.pilot.Find(bt.ctx, "#item-list", nil)
		if err != nil {
			t.Fatalf("Failed to find item list: %v", err)
		}

		// Find all items within the list
		items, err := list.FindAll(bt.ctx, "", &webpilot.FindOptions{
			TestID: "list-item",
		})
		if err != nil {
			t.Fatalf("Failed to find items in list: %v", err)
		}

		if len(items) != 3 {
			t.Errorf("Expected 3 items, got %d", len(items))
		}
	})
}

// TestSemanticSelectorsNear tests finding elements near another element.
func TestSemanticSelectorsNear(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(testPageSemanticSelectors)
	time.Sleep(200 * time.Millisecond)

	t.Run("FindButtonNearInput", func(t *testing.T) {
		elem, err := bt.pilot.Find(bt.ctx, "", &webpilot.FindOptions{
			Role: "button",
			Near: "#username-input",
		})
		if err != nil {
			t.Fatalf("Failed to find button near input: %v", err)
		}

		info := elem.Info()
		assertContains(t, info.Text, "Submit")
	})

	t.Run("FindInputNearLabel", func(t *testing.T) {
		elem, err := bt.pilot.Find(bt.ctx, "", &webpilot.FindOptions{
			Role: "textbox",
			Near: "#username-label",
		})
		if err != nil {
			t.Fatalf("Failed to find input near label: %v", err)
		}

		// Type to verify it's the correct input
		err = elem.Fill(bt.ctx, "testuser", nil)
		if err != nil {
			t.Fatalf("Failed to fill: %v", err)
		}

		value, err := elem.Value(bt.ctx)
		if err != nil {
			t.Fatalf("Failed to get value: %v", err)
		}
		if value != "testuser" {
			t.Errorf("Expected 'testuser', got %q", value)
		}
	})
}

// TestSemanticSelectorsTimeout tests timeout handling for semantic selectors.
func TestSemanticSelectorsTimeout(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(testPageSemanticSelectors)
	time.Sleep(200 * time.Millisecond)

	t.Run("TimeoutOnNotFound", func(t *testing.T) {
		start := time.Now()

		_, err := bt.pilot.Find(bt.ctx, "", &webpilot.FindOptions{
			TestID:  "non-existent-element",
			Timeout: 1 * time.Second,
		})

		elapsed := time.Since(start)

		if err == nil {
			t.Error("Expected error for non-existent element")
		}

		// Should have waited approximately 1 second
		if elapsed < 900*time.Millisecond {
			t.Errorf("Timeout too short: %v", elapsed)
		}
		if elapsed > 3*time.Second {
			t.Errorf("Timeout too long: %v", elapsed)
		}
	})

	t.Run("CustomTimeout", func(t *testing.T) {
		// Should find quickly with longer timeout
		elem, err := bt.pilot.Find(bt.ctx, "", &webpilot.FindOptions{
			TestID:  "submit-button",
			Timeout: 10 * time.Second,
		})
		if err != nil {
			t.Fatalf("Failed to find element: %v", err)
		}

		info := elem.Info()
		assertContains(t, info.Text, "Submit")
	})
}

// TestSemanticSelectorsInteraction tests interacting with elements found by semantic selectors.
func TestSemanticSelectorsInteraction(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(testPageSemanticSelectors)
	time.Sleep(200 * time.Millisecond)

	t.Run("ClickButtonByRole", func(t *testing.T) {
		elem, err := bt.pilot.Find(bt.ctx, "", &webpilot.FindOptions{
			Role: "button",
			Text: "Submit Form",
		})
		if err != nil {
			t.Fatalf("Failed to find button: %v", err)
		}

		// Click should succeed
		err = elem.Click(bt.ctx, nil)
		if err != nil {
			t.Fatalf("Failed to click button: %v", err)
		}
	})

	t.Run("FillInputByLabel", func(t *testing.T) {
		elem, err := bt.pilot.Find(bt.ctx, "", &webpilot.FindOptions{
			Label: "Email Address",
		})
		if err != nil {
			t.Fatalf("Failed to find email input: %v", err)
		}

		err = elem.Fill(bt.ctx, "test@example.com", nil)
		if err != nil {
			t.Fatalf("Failed to fill email: %v", err)
		}

		value, err := elem.Value(bt.ctx)
		if err != nil {
			t.Fatalf("Failed to get value: %v", err)
		}
		if value != "test@example.com" {
			t.Errorf("Expected 'test@example.com', got %q", value)
		}
	})

	t.Run("CheckCheckboxByTestID", func(t *testing.T) {
		elem, err := bt.pilot.Find(bt.ctx, "", &webpilot.FindOptions{
			TestID: "agree-checkbox",
		})
		if err != nil {
			t.Fatalf("Failed to find checkbox: %v", err)
		}

		// Check the checkbox
		err = elem.Check(bt.ctx, nil)
		if err != nil {
			t.Fatalf("Failed to check checkbox: %v", err)
		}

		isChecked, err := elem.IsChecked(bt.ctx)
		if err != nil {
			t.Fatalf("Failed to get checked state: %v", err)
		}
		if !isChecked {
			t.Error("Expected checkbox to be checked")
		}
	})
}
