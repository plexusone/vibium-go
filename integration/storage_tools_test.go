//go:build integration

package integration

import (
	"testing"
)

// TestLocalStorageGetSet tests localStorage get and set operations.
func TestLocalStorageGetSet(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html><html><body><p>Storage test</p></body></html>`)

	// Set a value via JavaScript
	_, err := bt.pilot.Evaluate(bt.ctx, `localStorage.setItem('testKey', 'testValue')`)
	if err != nil {
		t.Fatalf("Failed to set localStorage: %v", err)
	}

	// Get value via JavaScript
	result, err := bt.pilot.Evaluate(bt.ctx, `localStorage.getItem('testKey')`)
	if err != nil {
		t.Fatalf("Failed to get localStorage: %v", err)
	}
	if result != "testValue" {
		t.Errorf("Expected 'testValue', got %q", result)
	}
}

// TestLocalStorageList tests listing all localStorage items.
func TestLocalStorageList(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html><html><body><p>Storage test</p></body></html>`)

	// Clear and set multiple values
	_, err := bt.pilot.Evaluate(bt.ctx, `
		localStorage.clear();
		localStorage.setItem('key1', 'value1');
		localStorage.setItem('key2', 'value2');
		localStorage.setItem('key3', 'value3');
	`)
	if err != nil {
		t.Fatalf("Failed to set localStorage items: %v", err)
	}

	// Get length
	result, err := bt.pilot.Evaluate(bt.ctx, `localStorage.length`)
	if err != nil {
		t.Fatalf("Failed to get localStorage length: %v", err)
	}
	length, ok := result.(float64)
	if !ok || int(length) != 3 {
		t.Errorf("Expected 3 items, got %v", result)
	}
}

// TestLocalStorageDelete tests deleting a localStorage item.
func TestLocalStorageDelete(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html><html><body><p>Storage test</p></body></html>`)

	// Set and delete
	_, err := bt.pilot.Evaluate(bt.ctx, `
		localStorage.setItem('toDelete', 'value');
		localStorage.removeItem('toDelete');
	`)
	if err != nil {
		t.Fatalf("Failed to delete localStorage item: %v", err)
	}

	// Verify deleted
	result, err := bt.pilot.Evaluate(bt.ctx, `localStorage.getItem('toDelete')`)
	if err != nil {
		t.Fatalf("Failed to get deleted item: %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil after delete, got %v", result)
	}
}

// TestLocalStorageClear tests clearing all localStorage.
func TestLocalStorageClear(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html><html><body><p>Storage test</p></body></html>`)

	// Set values then clear
	_, err := bt.pilot.Evaluate(bt.ctx, `
		localStorage.setItem('key1', 'value1');
		localStorage.setItem('key2', 'value2');
		localStorage.clear();
	`)
	if err != nil {
		t.Fatalf("Failed to clear localStorage: %v", err)
	}

	// Verify cleared
	result, err := bt.pilot.Evaluate(bt.ctx, `localStorage.length`)
	if err != nil {
		t.Fatalf("Failed to get localStorage length: %v", err)
	}
	length, ok := result.(float64)
	if !ok || int(length) != 0 {
		t.Errorf("Expected 0 items after clear, got %v", result)
	}
}

// TestSessionStorageGetSet tests sessionStorage get and set operations.
func TestSessionStorageGetSet(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html><html><body><p>Storage test</p></body></html>`)

	// Set a value via JavaScript
	_, err := bt.pilot.Evaluate(bt.ctx, `sessionStorage.setItem('sessionKey', 'sessionValue')`)
	if err != nil {
		t.Fatalf("Failed to set sessionStorage: %v", err)
	}

	// Get value
	result, err := bt.pilot.Evaluate(bt.ctx, `sessionStorage.getItem('sessionKey')`)
	if err != nil {
		t.Fatalf("Failed to get sessionStorage: %v", err)
	}
	if result != "sessionValue" {
		t.Errorf("Expected 'sessionValue', got %q", result)
	}
}

// TestSessionStorageList tests listing all sessionStorage items.
func TestSessionStorageList(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html><html><body><p>Storage test</p></body></html>`)

	// Clear and set multiple values
	_, err := bt.pilot.Evaluate(bt.ctx, `
		sessionStorage.clear();
		sessionStorage.setItem('skey1', 'sval1');
		sessionStorage.setItem('skey2', 'sval2');
	`)
	if err != nil {
		t.Fatalf("Failed to set sessionStorage items: %v", err)
	}

	// Get length
	result, err := bt.pilot.Evaluate(bt.ctx, `sessionStorage.length`)
	if err != nil {
		t.Fatalf("Failed to get sessionStorage length: %v", err)
	}
	length, ok := result.(float64)
	if !ok || int(length) != 2 {
		t.Errorf("Expected 2 items, got %v", result)
	}
}

// TestSessionStorageDelete tests deleting a sessionStorage item.
func TestSessionStorageDelete(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html><html><body><p>Storage test</p></body></html>`)

	// Set and delete
	_, err := bt.pilot.Evaluate(bt.ctx, `
		sessionStorage.setItem('toDelete', 'value');
		sessionStorage.removeItem('toDelete');
	`)
	if err != nil {
		t.Fatalf("Failed to delete sessionStorage item: %v", err)
	}

	// Verify deleted
	result, err := bt.pilot.Evaluate(bt.ctx, `sessionStorage.getItem('toDelete')`)
	if err != nil {
		t.Fatalf("Failed to get deleted item: %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil after delete, got %v", result)
	}
}

// TestSessionStorageClear tests clearing all sessionStorage.
func TestSessionStorageClear(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html><html><body><p>Storage test</p></body></html>`)

	// Set values then clear
	_, err := bt.pilot.Evaluate(bt.ctx, `
		sessionStorage.setItem('key1', 'value1');
		sessionStorage.setItem('key2', 'value2');
		sessionStorage.clear();
	`)
	if err != nil {
		t.Fatalf("Failed to clear sessionStorage: %v", err)
	}

	// Verify cleared
	result, err := bt.pilot.Evaluate(bt.ctx, `sessionStorage.length`)
	if err != nil {
		t.Fatalf("Failed to get sessionStorage length: %v", err)
	}
	length, ok := result.(float64)
	if !ok || int(length) != 0 {
		t.Errorf("Expected 0 items after clear, got %v", result)
	}
}

// TestStoragePersistenceAcrossNavigation tests storage persistence.
func TestStoragePersistenceAcrossNavigation(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	// Set up initial page and storage
	bt.go_(`data:text/html,<!DOCTYPE html><html><body><p>Page 1</p></body></html>`)
	_, err := bt.pilot.Evaluate(bt.ctx, `
		localStorage.setItem('persistent', 'localStorage survives');
		sessionStorage.setItem('session', 'sessionStorage survives');
	`)
	if err != nil {
		t.Fatalf("Failed to set storage: %v", err)
	}

	// Navigate to another page (same origin - data URLs share origin)
	bt.go_(`data:text/html,<!DOCTYPE html><html><body><p>Page 2</p></body></html>`)

	// Check localStorage persists
	result, err := bt.pilot.Evaluate(bt.ctx, `localStorage.getItem('persistent')`)
	if err != nil {
		t.Fatalf("Failed to get localStorage: %v", err)
	}
	if result != "localStorage survives" {
		t.Errorf("localStorage did not persist: got %v", result)
	}

	// Check sessionStorage persists within same tab
	result, err = bt.pilot.Evaluate(bt.ctx, `sessionStorage.getItem('session')`)
	if err != nil {
		t.Fatalf("Failed to get sessionStorage: %v", err)
	}
	if result != "sessionStorage survives" {
		t.Errorf("sessionStorage did not persist: got %v", result)
	}
}
