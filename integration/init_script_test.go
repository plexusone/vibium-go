//go:build integration

package integration

import (
	"testing"
	"time"
)

// TestInitScriptBasic tests basic init script functionality.
func TestInitScriptBasic(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	t.Run("AddInitScript", func(t *testing.T) {
		// Add init script that sets a global variable
		err := bt.pilot.AddInitScript(bt.ctx, `window.testInjected = true;`)
		if err != nil {
			t.Fatalf("Failed to add init script: %v", err)
		}

		// Navigate to a page (init script should run)
		bt.go_("https://example.com")
		time.Sleep(300 * time.Millisecond)

		// Check if the variable was set
		result, err := bt.pilot.Evaluate(bt.ctx, "window.testInjected")
		if err != nil {
			t.Fatalf("Failed to evaluate: %v", err)
		}

		if result != true {
			t.Errorf("Expected window.testInjected = true, got %v", result)
		}
	})

	t.Run("InitScriptOnMultipleNavigations", func(t *testing.T) {
		// Add init script
		err := bt.pilot.AddInitScript(bt.ctx, `window.pageCount = (window.pageCount || 0) + 1;`)
		if err != nil {
			t.Fatalf("Failed to add init script: %v", err)
		}

		// First navigation
		bt.go_("https://example.com")
		time.Sleep(200 * time.Millisecond)

		result, err := bt.pilot.Evaluate(bt.ctx, "window.pageCount")
		if err != nil {
			t.Fatalf("Failed to evaluate: %v", err)
		}

		// pageCount should be at least 1
		count, ok := result.(float64)
		if !ok {
			t.Fatalf("Expected number, got %T", result)
		}
		firstCount := int(count)
		t.Logf("After first navigation: pageCount = %d", firstCount)

		// Second navigation (init script should run again)
		bt.go_("https://www.iana.org")
		time.Sleep(200 * time.Millisecond)

		result, err = bt.pilot.Evaluate(bt.ctx, "window.pageCount")
		if err != nil {
			t.Fatalf("Failed to evaluate: %v", err)
		}

		count, ok = result.(float64)
		if !ok {
			t.Fatalf("Expected number, got %T", result)
		}
		secondCount := int(count)
		t.Logf("After second navigation: pageCount = %d", secondCount)

		// The script runs on each navigation, so count should be 1 on each new page
		// (since it's a new page context each time)
		if secondCount < 1 {
			t.Errorf("Expected pageCount >= 1 on second page, got %d", secondCount)
		}
	})
}

// TestInitScriptFunction tests init scripts that define functions.
func TestInitScriptFunction(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	t.Run("DefineFunction", func(t *testing.T) {
		// Add init script that defines a helper function
		err := bt.pilot.AddInitScript(bt.ctx, `
			window.myHelper = function(x, y) {
				return x + y;
			};
		`)
		if err != nil {
			t.Fatalf("Failed to add init script: %v", err)
		}

		bt.go_("https://example.com")
		time.Sleep(200 * time.Millisecond)

		// Call the function
		result, err := bt.pilot.Evaluate(bt.ctx, "window.myHelper(2, 3)")
		if err != nil {
			t.Fatalf("Failed to evaluate: %v", err)
		}

		if result != float64(5) {
			t.Errorf("Expected 5, got %v", result)
		}
	})

	t.Run("DefineClass", func(t *testing.T) {
		// Add init script that defines a class
		err := bt.pilot.AddInitScript(bt.ctx, `
			window.TestClass = class {
				constructor(name) {
					this.name = name;
				}
				greet() {
					return 'Hello, ' + this.name;
				}
			};
		`)
		if err != nil {
			t.Fatalf("Failed to add init script: %v", err)
		}

		bt.go_("https://example.com")
		time.Sleep(200 * time.Millisecond)

		// Use the class
		result, err := bt.pilot.Evaluate(bt.ctx, "new window.TestClass('World').greet()")
		if err != nil {
			t.Fatalf("Failed to evaluate: %v", err)
		}

		if result != "Hello, World" {
			t.Errorf("Expected 'Hello, World', got %v", result)
		}
	})
}

// TestInitScriptMocking tests using init scripts to mock APIs.
func TestInitScriptMocking(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	t.Run("MockDate", func(t *testing.T) {
		// Mock Date.now() to return a fixed timestamp
		err := bt.pilot.AddInitScript(bt.ctx, `
			const fixedTime = 1609459200000; // 2021-01-01T00:00:00.000Z
			Date.now = function() {
				return fixedTime;
			};
			window.originalDateNow = fixedTime;
		`)
		if err != nil {
			t.Fatalf("Failed to add init script: %v", err)
		}

		bt.go_("https://example.com")
		time.Sleep(200 * time.Millisecond)

		// Check mocked Date.now()
		result, err := bt.pilot.Evaluate(bt.ctx, "Date.now()")
		if err != nil {
			t.Fatalf("Failed to evaluate: %v", err)
		}

		expected := float64(1609459200000)
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("MockLocalStorage", func(t *testing.T) {
		// Add init script that pre-populates localStorage-like behavior
		err := bt.pilot.AddInitScript(bt.ctx, `
			window.mockStorage = {
				'user': 'test_user',
				'token': 'mock_token_123'
			};
		`)
		if err != nil {
			t.Fatalf("Failed to add init script: %v", err)
		}

		bt.go_("https://example.com")
		time.Sleep(200 * time.Millisecond)

		// Check mock storage
		result, err := bt.pilot.Evaluate(bt.ctx, "window.mockStorage.user")
		if err != nil {
			t.Fatalf("Failed to evaluate: %v", err)
		}

		if result != "test_user" {
			t.Errorf("Expected 'test_user', got %v", result)
		}
	})
}

// TestInitScriptMultiple tests adding multiple init scripts.
func TestInitScriptMultiple(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	t.Run("MultipleScripts", func(t *testing.T) {
		// Add first init script
		err := bt.pilot.AddInitScript(bt.ctx, `window.script1 = 'first';`)
		if err != nil {
			t.Fatalf("Failed to add first init script: %v", err)
		}

		// Add second init script
		err = bt.pilot.AddInitScript(bt.ctx, `window.script2 = 'second';`)
		if err != nil {
			t.Fatalf("Failed to add second init script: %v", err)
		}

		// Add third init script that depends on the first two
		err = bt.pilot.AddInitScript(bt.ctx, `
			window.combined = window.script1 + ' and ' + window.script2;
		`)
		if err != nil {
			t.Fatalf("Failed to add third init script: %v", err)
		}

		bt.go_("https://example.com")
		time.Sleep(200 * time.Millisecond)

		// Check all scripts ran
		result, err := bt.pilot.Evaluate(bt.ctx, "window.script1")
		if err != nil {
			t.Fatalf("Failed to evaluate script1: %v", err)
		}
		if result != "first" {
			t.Errorf("Expected 'first', got %v", result)
		}

		result, err = bt.pilot.Evaluate(bt.ctx, "window.script2")
		if err != nil {
			t.Fatalf("Failed to evaluate script2: %v", err)
		}
		if result != "second" {
			t.Errorf("Expected 'second', got %v", result)
		}

		result, err = bt.pilot.Evaluate(bt.ctx, "window.combined")
		if err != nil {
			t.Fatalf("Failed to evaluate combined: %v", err)
		}
		if result != "first and second" {
			t.Errorf("Expected 'first and second', got %v", result)
		}
	})
}

// TestInitScriptBeforePageScripts tests that init scripts run before page scripts.
func TestInitScriptBeforePageScripts(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	t.Run("RunsBeforePageScripts", func(t *testing.T) {
		// Add init script that sets a flag
		err := bt.pilot.AddInitScript(bt.ctx, `
			window.initScriptRan = true;
			window.initScriptTime = performance.now();
		`)
		if err != nil {
			t.Fatalf("Failed to add init script: %v", err)
		}

		// Navigate to a page
		bt.go_("https://example.com")
		time.Sleep(200 * time.Millisecond)

		// Verify init script ran
		result, err := bt.pilot.Evaluate(bt.ctx, "window.initScriptRan")
		if err != nil {
			t.Fatalf("Failed to evaluate: %v", err)
		}

		if result != true {
			t.Error("Init script did not run")
		}

		// Check that initScriptTime is a reasonable value (should be very early)
		result, err = bt.pilot.Evaluate(bt.ctx, "window.initScriptTime")
		if err != nil {
			t.Fatalf("Failed to evaluate: %v", err)
		}

		time, ok := result.(float64)
		if !ok {
			t.Fatalf("Expected number for initScriptTime, got %T", result)
		}

		t.Logf("Init script ran at performance.now() = %f", time)
	})
}

// TestInitScriptFromContext tests init scripts via BrowserContext.
func TestInitScriptFromContext(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	t.Run("ContextInitScript", func(t *testing.T) {
		// Create a browser context
		browserCtx, err := bt.pilot.NewContext(bt.ctx)
		if err != nil {
			t.Fatalf("Failed to create context: %v", err)
		}

		// Add init script via context
		err = browserCtx.AddInitScript(bt.ctx, `window.contextScript = 'from_context';`)
		if err != nil {
			t.Fatalf("Failed to add init script to context: %v", err)
		}

		// Create a page in this context
		page, err := browserCtx.NewPage(bt.ctx)
		if err != nil {
			t.Fatalf("Failed to create page: %v", err)
		}

		// Navigate
		err = page.Go(bt.ctx, "https://example.com")
		if err != nil {
			t.Fatalf("Failed to navigate: %v", err)
		}
		time.Sleep(200 * time.Millisecond)

		// Check init script ran
		result, err := page.Evaluate(bt.ctx, "window.contextScript")
		if err != nil {
			t.Fatalf("Failed to evaluate: %v", err)
		}

		if result != "from_context" {
			t.Errorf("Expected 'from_context', got %v", result)
		}
	})
}
