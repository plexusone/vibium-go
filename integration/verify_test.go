//go:build integration

package integration

import (
	"strings"
	"testing"
)

// TestVerifyInputValue tests verifying input element values (like verify_value MCP tool).
func TestVerifyInputValue(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html>
<html><body>
<input id="name" type="text" value="John Doe">
<input id="email" type="email" value="john@example.com">
<input id="empty" type="text" value="">
</body></html>`)

	// Verify expected value
	nameInput := bt.find("#name")
	value, err := nameInput.Value(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to get value: %v", err)
	}
	expected := "John Doe"
	if value != expected {
		t.Errorf("Value mismatch: expected %q, got %q", expected, value)
	}

	// Verify another value
	emailInput := bt.find("#email")
	value, err = emailInput.Value(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to get email value: %v", err)
	}
	expected = "john@example.com"
	if value != expected {
		t.Errorf("Email value mismatch: expected %q, got %q", expected, value)
	}

	// Verify empty value
	emptyInput := bt.find("#empty")
	value, err = emptyInput.Value(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to get empty value: %v", err)
	}
	if value != "" {
		t.Errorf("Empty value should be empty string, got %q", value)
	}
}

// TestVerifyListVisible tests verifying multiple items are visible (like verify_list_visible MCP tool).
func TestVerifyListVisible(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html>
<html><body>
<ul id="features">
  <li>Feature A</li>
  <li>Feature B</li>
  <li>Feature C</li>
</ul>
<div id="hidden" style="display:none">Hidden Text</div>
</body></html>`)

	// Items that should be visible
	items := []string{"Feature A", "Feature B", "Feature C"}
	var found []string
	var missing []string

	for _, item := range items {
		// Check if item is in body text
		result, err := bt.pilot.Evaluate(bt.ctx, `document.body.textContent.includes("`+item+`")`)
		if err != nil {
			missing = append(missing, item)
			continue
		}
		if visible, ok := result.(bool); ok && visible {
			found = append(found, item)
		} else {
			missing = append(missing, item)
		}
	}

	if len(missing) > 0 {
		t.Errorf("Missing items: %v", missing)
	}
	if len(found) != len(items) {
		t.Errorf("Expected %d items found, got %d", len(items), len(found))
	}
	t.Logf("Found items: %v", found)

	// Test scoped search within selector
	result, err := bt.pilot.Evaluate(bt.ctx, `
		(function() {
			const el = document.querySelector('#features');
			return el && el.textContent.includes('Feature A');
		})()
	`)
	if err != nil {
		t.Fatalf("Scoped search failed: %v", err)
	}
	if visible, ok := result.(bool); !ok || !visible {
		t.Error("Feature A should be visible in #features")
	}
}

// TestGenerateLocator tests generating locator strings for elements.
func TestGenerateLocator(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html>
<html><body>
<button id="submit-btn" data-testid="submit" role="button">Submit Form</button>
<input id="email" type="email" placeholder="Enter email" aria-label="Email address">
<a href="/about" class="nav-link">About Us</a>
</body></html>`)

	// Test CSS selector works
	elem := bt.find("#submit-btn")
	text, err := elem.Text(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to get text: %v", err)
	}
	if text != "Submit Form" {
		t.Errorf("Expected 'Submit Form', got %q", text)
	}

	// Test data-testid selector
	elem = bt.find("[data-testid='submit']")
	text, err = elem.Text(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to get text by testid: %v", err)
	}
	if text != "Submit Form" {
		t.Errorf("Expected 'Submit Form', got %q", text)
	}

	// Test role selector via evaluate
	result, err := bt.pilot.Evaluate(bt.ctx, `document.querySelector('[role="button"]').textContent`)
	if err != nil {
		t.Fatalf("Failed role query: %v", err)
	}
	if result != "Submit Form" {
		t.Errorf("Role query: expected 'Submit Form', got %v", result)
	}

	// Test xpath via evaluate
	result, err = bt.pilot.Evaluate(bt.ctx, `
		document.evaluate("//button[@id='submit-btn']", document, null,
			XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue.textContent
	`)
	if err != nil {
		t.Fatalf("XPath query failed: %v", err)
	}
	if result != "Submit Form" {
		t.Errorf("XPath query: expected 'Submit Form', got %v", result)
	}

	// Test text content search
	result, err = bt.pilot.Evaluate(bt.ctx, `
		Array.from(document.querySelectorAll('a')).find(a => a.textContent === 'About Us')?.href
	`)
	if err != nil {
		t.Fatalf("Text search failed: %v", err)
	}
	resultStr, ok := result.(string)
	if !ok || !strings.Contains(resultStr, "/about") {
		t.Errorf("Text search: expected href containing '/about', got %v", result)
	}
}

// TestVerifyElementStates tests verifying various element states.
func TestVerifyElementStates(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html>
<html><body>
<button id="enabled">Enabled</button>
<button id="disabled" disabled>Disabled</button>
<input id="readonly" type="text" readonly value="Read only">
<input id="checked" type="checkbox" checked>
<input id="unchecked" type="checkbox">
<div id="visible">Visible</div>
<div id="hidden" style="display:none">Hidden</div>
</body></html>`)

	// Test enabled state
	enabledBtn := bt.find("#enabled")
	enabled, err := enabledBtn.IsEnabled(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to check enabled: %v", err)
	}
	if !enabled {
		t.Error("Button should be enabled")
	}

	// Test disabled state
	disabledBtn := bt.find("#disabled")
	enabled, err = disabledBtn.IsEnabled(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to check disabled: %v", err)
	}
	if enabled {
		t.Error("Button should be disabled")
	}

	// Test checked state
	checkedBox := bt.find("#checked")
	checked, err := checkedBox.IsChecked(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to check checked: %v", err)
	}
	if !checked {
		t.Error("Checkbox should be checked")
	}

	// Test unchecked state
	uncheckedBox := bt.find("#unchecked")
	checked, err = uncheckedBox.IsChecked(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to check unchecked: %v", err)
	}
	if checked {
		t.Error("Checkbox should be unchecked")
	}

	// Test visible state
	visibleDiv := bt.find("#visible")
	visible, err := visibleDiv.IsVisible(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to check visible: %v", err)
	}
	if !visible {
		t.Error("Div should be visible")
	}
}
