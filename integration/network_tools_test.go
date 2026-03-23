//go:build integration

package integration

import (
	"strings"
	"testing"
)

// TestNetworkRequestsCapture tests capturing network requests.
func TestNetworkRequestsCapture(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	// Navigate to a page that makes requests
	bt.go_("https://example.com")

	// Get network requests
	requests, err := bt.pilot.NetworkRequests(bt.ctx, nil)
	if err != nil {
		t.Fatalf("Failed to get network requests: %v", err)
	}

	// Should have at least the main document request
	if len(requests) == 0 {
		t.Error("Expected at least one network request")
	}

	// Check for document request
	foundDoc := false
	for _, req := range requests {
		if strings.Contains(req.URL, "example.com") {
			foundDoc = true
			if req.Method != "GET" {
				t.Errorf("Expected GET method, got %s", req.Method)
			}
			break
		}
	}
	if !foundDoc {
		t.Error("Did not find example.com request")
	}
}

// TestNetworkRequestsClear tests clearing captured network requests.
func TestNetworkRequestsClear(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	// Navigate to capture some requests
	bt.go_("https://example.com")

	// Verify we have requests
	requests, err := bt.pilot.NetworkRequests(bt.ctx, nil)
	if err != nil {
		t.Fatalf("Failed to get network requests: %v", err)
	}
	if len(requests) == 0 {
		t.Skip("No requests captured to clear")
	}

	// Clear requests
	err = bt.pilot.ClearNetworkRequests(bt.ctx)
	if err != nil {
		t.Fatalf("Failed to clear network requests: %v", err)
	}

	// Get requests - should be empty
	requests, err = bt.pilot.NetworkRequests(bt.ctx, nil)
	if err != nil {
		t.Fatalf("Failed to get network requests after clear: %v", err)
	}

	if len(requests) != 0 {
		t.Errorf("Expected 0 requests after clear, got %d", len(requests))
	}
}

// TestRouteAndUnroute tests setting and removing routes.
func TestRouteAndUnroute(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	// Navigate first to have a valid context
	bt.go_("https://example.com")

	// Set up a route
	pattern := "**/api/*"
	err := bt.pilot.Route(bt.ctx, pattern, nil)
	if err != nil {
		t.Fatalf("Failed to set route: %v", err)
	}
	t.Log("Route set successfully")

	// Remove the route
	err = bt.pilot.Unroute(bt.ctx, pattern)
	if err != nil {
		t.Fatalf("Failed to unroute: %v", err)
	}
	t.Log("Route removed successfully")
}

// TestNetworkOffline tests setting network offline state.
func TestNetworkOffline(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	// Navigate first
	bt.go_("https://example.com")

	// Set offline
	err := bt.pilot.SetOffline(bt.ctx, true)
	if err != nil {
		t.Fatalf("Failed to set offline: %v", err)
	}
	t.Log("Network set to offline mode")

	// Set back online
	err = bt.pilot.SetOffline(bt.ctx, false)
	if err != nil {
		t.Fatalf("Failed to set online: %v", err)
	}
	t.Log("Network set back to online mode")
}

// TestNetworkRequestsFiltering tests filtering network requests.
func TestNetworkRequestsFiltering(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	// Navigate to a page
	bt.go_("https://example.com")

	// Get all requests
	requests, err := bt.pilot.NetworkRequests(bt.ctx, nil)
	if err != nil {
		t.Fatalf("Failed to get network requests: %v", err)
	}

	t.Logf("Captured %d network requests", len(requests))
	for i, req := range requests {
		t.Logf("  [%d] %s %s (type: %s)", i, req.Method, req.URL, req.ResourceType)
	}
}
