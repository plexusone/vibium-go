//go:build integration

package integration

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/plexusone/vibium-go"
)

// testPageStorage is an HTML page for testing storage operations.
const testPageStorage = `data:text/html,<!DOCTYPE html>
<html>
<head><title>Storage Test</title></head>
<body>
    <h1>Storage Test Page</h1>
    <script>
        // Set some initial storage values
        localStorage.setItem('testKey', 'testValue');
        localStorage.setItem('user', 'john');
        sessionStorage.setItem('sessionKey', 'sessionValue');
        sessionStorage.setItem('token', 'abc123');
    </script>
</body>
</html>`

// TestStorageStateGet tests getting storage state.
func TestStorageStateGet(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	t.Run("GetStorageState", func(t *testing.T) {
		// Navigate to set up storage
		bt.go_(testPageStorage)
		time.Sleep(300 * time.Millisecond)

		// Get storage state
		state, err := bt.vibe.StorageState(bt.ctx)
		if err != nil {
			t.Fatalf("Failed to get storage state: %v", err)
		}

		if state == nil {
			t.Fatal("StorageState returned nil")
		}

		// Verify we have origins data
		t.Logf("Got storage state with %d cookies and %d origins", len(state.Cookies), len(state.Origins))

		// Check for localStorage values
		foundLocalStorage := false
		for _, origin := range state.Origins {
			if len(origin.LocalStorage) > 0 {
				foundLocalStorage = true
				t.Logf("Origin %s has %d localStorage items", origin.Origin, len(origin.LocalStorage))
			}
			if len(origin.SessionStorage) > 0 {
				t.Logf("Origin %s has %d sessionStorage items", origin.Origin, len(origin.SessionStorage))
			}
		}

		if !foundLocalStorage {
			t.Log("Warning: No localStorage found (may be expected for data: URL)")
		}
	})

	t.Run("GetStorageStateFromRealSite", func(t *testing.T) {
		// Use a real site to test cookie handling
		bt.go_("https://example.com")
		time.Sleep(500 * time.Millisecond)

		// Set some localStorage via JavaScript
		_, err := bt.vibe.Evaluate(bt.ctx, `
			localStorage.setItem('test_key', 'test_value');
			sessionStorage.setItem('session_test', 'session_value');
		`)
		if err != nil {
			t.Fatalf("Failed to set storage: %v", err)
		}

		// Get storage state
		state, err := bt.vibe.StorageState(bt.ctx)
		if err != nil {
			t.Fatalf("Failed to get storage state: %v", err)
		}

		// Verify structure
		if state == nil {
			t.Fatal("StorageState returned nil")
		}

		t.Logf("Cookies: %d, Origins: %d", len(state.Cookies), len(state.Origins))

		// Look for our test values
		for _, origin := range state.Origins {
			if origin.Origin == "https://example.com" {
				if val, ok := origin.LocalStorage["test_key"]; ok {
					if val != "test_value" {
						t.Errorf("Expected localStorage 'test_key' = 'test_value', got %q", val)
					}
				} else {
					t.Error("localStorage 'test_key' not found")
				}

				if val, ok := origin.SessionStorage["session_test"]; ok {
					if val != "session_value" {
						t.Errorf("Expected sessionStorage 'session_test' = 'session_value', got %q", val)
					}
				} else {
					t.Error("sessionStorage 'session_test' not found")
				}
			}
		}
	})
}

// TestStorageStateRestore tests restoring storage state.
func TestStorageStateRestore(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	t.Run("SetStorageState", func(t *testing.T) {
		// Create a storage state to restore
		state := &vibium.StorageState{
			Cookies: []vibium.Cookie{
				{
					Name:   "test_cookie",
					Value:  "cookie_value",
					Domain: "example.com",
					Path:   "/",
				},
			},
			Origins: []vibium.StorageStateOrigin{
				{
					Origin: "https://example.com",
					LocalStorage: map[string]string{
						"restored_key": "restored_value",
						"user_id":      "12345",
					},
					SessionStorage: map[string]string{
						"session_token": "xyz789",
					},
				},
			},
		}

		// Navigate first (needed to set storage for an origin)
		bt.go_("https://example.com")
		time.Sleep(300 * time.Millisecond)

		// Restore storage state
		err := bt.vibe.SetStorageState(bt.ctx, state)
		if err != nil {
			t.Fatalf("Failed to set storage state: %v", err)
		}

		// Verify localStorage was set
		result, err := bt.vibe.Evaluate(bt.ctx, "localStorage.getItem('restored_key')")
		if err != nil {
			t.Fatalf("Failed to check localStorage: %v", err)
		}
		if result != "restored_value" {
			t.Errorf("Expected 'restored_value', got %v", result)
		}

		// Verify user_id
		result, err = bt.vibe.Evaluate(bt.ctx, "localStorage.getItem('user_id')")
		if err != nil {
			t.Fatalf("Failed to check localStorage: %v", err)
		}
		if result != "12345" {
			t.Errorf("Expected '12345', got %v", result)
		}

		// Verify sessionStorage was set
		result, err = bt.vibe.Evaluate(bt.ctx, "sessionStorage.getItem('session_token')")
		if err != nil {
			t.Fatalf("Failed to check sessionStorage: %v", err)
		}
		if result != "xyz789" {
			t.Errorf("Expected 'xyz789', got %v", result)
		}
	})

	t.Run("RoundTrip", func(t *testing.T) {
		// Navigate and set some data
		bt.go_("https://example.com")
		time.Sleep(300 * time.Millisecond)

		// Set storage via JavaScript
		_, err := bt.vibe.Evaluate(bt.ctx, `
			localStorage.setItem('roundtrip_key', 'roundtrip_value');
			sessionStorage.setItem('roundtrip_session', 'session_data');
		`)
		if err != nil {
			t.Fatalf("Failed to set storage: %v", err)
		}

		// Get storage state
		state, err := bt.vibe.StorageState(bt.ctx)
		if err != nil {
			t.Fatalf("Failed to get storage state: %v", err)
		}

		// Serialize to JSON (like saving to file)
		jsonData, err := json.Marshal(state)
		if err != nil {
			t.Fatalf("Failed to marshal state: %v", err)
		}
		t.Logf("Storage state JSON size: %d bytes", len(jsonData))

		// Clear storage
		err = bt.vibe.ClearStorage(bt.ctx)
		if err != nil {
			t.Fatalf("Failed to clear storage: %v", err)
		}

		// Verify cleared
		result, err := bt.vibe.Evaluate(bt.ctx, "localStorage.getItem('roundtrip_key')")
		if err != nil {
			t.Fatalf("Failed to check localStorage: %v", err)
		}
		if result != nil {
			t.Errorf("Expected null after clear, got %v", result)
		}

		// Deserialize and restore
		var restoredState vibium.StorageState
		if err := json.Unmarshal(jsonData, &restoredState); err != nil {
			t.Fatalf("Failed to unmarshal state: %v", err)
		}

		err = bt.vibe.SetStorageState(bt.ctx, &restoredState)
		if err != nil {
			t.Fatalf("Failed to restore storage state: %v", err)
		}

		// Verify restored
		result, err = bt.vibe.Evaluate(bt.ctx, "localStorage.getItem('roundtrip_key')")
		if err != nil {
			t.Fatalf("Failed to check localStorage: %v", err)
		}
		if result != "roundtrip_value" {
			t.Errorf("Expected 'roundtrip_value' after restore, got %v", result)
		}
	})
}

// TestStorageStateClear tests clearing storage.
func TestStorageStateClear(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	t.Run("ClearStorage", func(t *testing.T) {
		bt.go_("https://example.com")
		time.Sleep(300 * time.Millisecond)

		// Set some storage
		_, err := bt.vibe.Evaluate(bt.ctx, `
			localStorage.setItem('to_clear', 'value1');
			sessionStorage.setItem('to_clear_session', 'value2');
		`)
		if err != nil {
			t.Fatalf("Failed to set storage: %v", err)
		}

		// Verify set
		result, _ := bt.vibe.Evaluate(bt.ctx, "localStorage.getItem('to_clear')")
		if result != "value1" {
			t.Fatalf("Storage not set properly, got %v", result)
		}

		// Clear storage
		err = bt.vibe.ClearStorage(bt.ctx)
		if err != nil {
			t.Fatalf("Failed to clear storage: %v", err)
		}

		// Verify localStorage cleared
		result, err = bt.vibe.Evaluate(bt.ctx, "localStorage.getItem('to_clear')")
		if err != nil {
			t.Fatalf("Failed to check localStorage: %v", err)
		}
		if result != nil {
			t.Errorf("Expected null after clear, got %v", result)
		}

		// Verify sessionStorage cleared
		result, err = bt.vibe.Evaluate(bt.ctx, "sessionStorage.getItem('to_clear_session')")
		if err != nil {
			t.Fatalf("Failed to check sessionStorage: %v", err)
		}
		if result != nil {
			t.Errorf("Expected null after clear, got %v", result)
		}
	})

	t.Run("ClearStorageOnBlankPage", func(t *testing.T) {
		// ClearStorage should not error on blank page
		err := bt.vibe.ClearStorage(bt.ctx)
		if err != nil {
			t.Errorf("ClearStorage should not error on blank page: %v", err)
		}
	})
}

// TestStorageStateWithCookies tests cookie handling in storage state.
func TestStorageStateWithCookies(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	t.Run("CookiesIncluded", func(t *testing.T) {
		bt.go_("https://example.com")
		time.Sleep(300 * time.Millisecond)

		// Get browser context to set cookies
		browserCtx, err := bt.vibe.NewContext(bt.ctx)
		if err != nil {
			t.Fatalf("Failed to get context: %v", err)
		}

		// Set a cookie
		err = browserCtx.SetCookies(bt.ctx, []vibium.SetCookieParam{
			{
				Name:   "test_cookie",
				Value:  "cookie_value",
				Domain: "example.com",
				Path:   "/",
			},
		})
		if err != nil {
			t.Fatalf("Failed to set cookie: %v", err)
		}

		// Get storage state
		state, err := bt.vibe.StorageState(bt.ctx)
		if err != nil {
			t.Fatalf("Failed to get storage state: %v", err)
		}

		// Check for our cookie
		foundCookie := false
		for _, cookie := range state.Cookies {
			if cookie.Name == "test_cookie" && cookie.Value == "cookie_value" {
				foundCookie = true
				break
			}
		}

		if !foundCookie {
			t.Error("Cookie not found in storage state")
		}
	})

	t.Run("RestoreCookies", func(t *testing.T) {
		state := &vibium.StorageState{
			Cookies: []vibium.Cookie{
				{
					Name:   "restored_cookie",
					Value:  "restored_value",
					Domain: "example.com",
					Path:   "/",
				},
			},
			Origins: []vibium.StorageStateOrigin{},
		}

		bt.go_("https://example.com")
		time.Sleep(300 * time.Millisecond)

		err := bt.vibe.SetStorageState(bt.ctx, state)
		if err != nil {
			t.Fatalf("Failed to set storage state: %v", err)
		}

		// Get cookies to verify
		browserCtx, err := bt.vibe.NewContext(bt.ctx)
		if err != nil {
			t.Fatalf("Failed to get context: %v", err)
		}

		cookies, err := browserCtx.Cookies(bt.ctx, "https://example.com")
		if err != nil {
			t.Fatalf("Failed to get cookies: %v", err)
		}

		foundCookie := false
		for _, cookie := range cookies {
			if cookie.Name == "restored_cookie" && cookie.Value == "restored_value" {
				foundCookie = true
				break
			}
		}

		if !foundCookie {
			t.Error("Restored cookie not found")
		}
	})
}

// TestStorageStateMultipleOrigins tests storage state with multiple origins.
func TestStorageStateMultipleOrigins(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	t.Run("MultipleOriginsInState", func(t *testing.T) {
		// Create state with multiple origins
		state := &vibium.StorageState{
			Origins: []vibium.StorageStateOrigin{
				{
					Origin: "https://example.com",
					LocalStorage: map[string]string{
						"key1": "value1",
					},
				},
				{
					Origin: "https://www.iana.org",
					LocalStorage: map[string]string{
						"key2": "value2",
					},
				},
			},
		}

		// Navigate to first origin and restore
		bt.go_("https://example.com")
		time.Sleep(300 * time.Millisecond)

		err := bt.vibe.SetStorageState(bt.ctx, state)
		if err != nil {
			t.Fatalf("Failed to set storage state: %v", err)
		}

		// Verify first origin's storage
		result, err := bt.vibe.Evaluate(bt.ctx, "localStorage.getItem('key1')")
		if err != nil {
			t.Fatalf("Failed to check localStorage: %v", err)
		}
		if result != "value1" {
			t.Errorf("Expected 'value1', got %v", result)
		}
	})
}
