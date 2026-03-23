//go:build integration

package integration

import (
	"testing"
	"time"

	"github.com/plexusone/webpilot"
)

// TestTracingBasic tests basic tracing functionality.
func TestTracingBasic(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	t.Run("StartAndStop", func(t *testing.T) {
		// Get tracing controller
		tracing := bt.pilot.Tracing()
		if tracing == nil {
			t.Fatal("Tracing() returned nil")
		}

		// Start tracing
		err := tracing.Start(bt.ctx, &webpilot.TracingStartOptions{
			Screenshots: true,
			Snapshots:   true,
			Title:       "Test Trace",
		})
		if err != nil {
			t.Fatalf("Failed to start tracing: %v", err)
		}

		// Navigate to generate some trace data
		bt.go_("https://example.com")
		time.Sleep(500 * time.Millisecond)

		// Stop tracing and get data
		data, err := tracing.Stop(bt.ctx, nil)
		if err != nil {
			t.Fatalf("Failed to stop tracing: %v", err)
		}

		// Verify we got trace data (should be a zip file)
		if len(data) == 0 {
			t.Error("Expected non-empty trace data")
		}

		// Check for ZIP signature (PK)
		if len(data) >= 2 && (data[0] != 0x50 || data[1] != 0x4B) {
			t.Logf("Warning: trace data may not be ZIP format (first bytes: %x %x)", data[0], data[1])
		}
	})

	t.Run("StartWithOptions", func(t *testing.T) {
		tracing := bt.pilot.Tracing()

		err := tracing.Start(bt.ctx, &webpilot.TracingStartOptions{
			Name:        "custom-trace",
			Screenshots: true,
			Snapshots:   true,
			Sources:     true,
			Title:       "Custom Test Trace",
		})
		if err != nil {
			t.Fatalf("Failed to start tracing with options: %v", err)
		}

		bt.go_("https://example.com")
		time.Sleep(300 * time.Millisecond)

		data, err := tracing.Stop(bt.ctx, nil)
		if err != nil {
			t.Fatalf("Failed to stop tracing: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected non-empty trace data")
		}
	})

	t.Run("StartWithNilOptions", func(t *testing.T) {
		tracing := bt.pilot.Tracing()

		// Start with nil options should work
		err := tracing.Start(bt.ctx, nil)
		if err != nil {
			t.Fatalf("Failed to start tracing with nil options: %v", err)
		}

		bt.go_("https://example.com")

		data, err := tracing.Stop(bt.ctx, nil)
		if err != nil {
			t.Fatalf("Failed to stop tracing: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected non-empty trace data")
		}
	})
}

// TestTracingChunks tests trace chunk functionality.
func TestTracingChunks(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	t.Run("ChunkRecording", func(t *testing.T) {
		tracing := bt.pilot.Tracing()

		// Start main trace
		err := tracing.Start(bt.ctx, &webpilot.TracingStartOptions{
			Screenshots: true,
			Title:       "Chunked Trace",
		})
		if err != nil {
			t.Fatalf("Failed to start tracing: %v", err)
		}

		// Navigate first page
		bt.go_("https://example.com")
		time.Sleep(300 * time.Millisecond)

		// Start chunk
		err = tracing.StartChunk(bt.ctx, &webpilot.TracingChunkOptions{
			Name:  "chunk1",
			Title: "First Chunk",
		})
		if err != nil {
			t.Fatalf("Failed to start chunk: %v", err)
		}

		// Some activity in the chunk
		time.Sleep(200 * time.Millisecond)

		// Stop chunk
		chunkData, err := tracing.StopChunk(bt.ctx, nil)
		if err != nil {
			t.Fatalf("Failed to stop chunk: %v", err)
		}

		if len(chunkData) == 0 {
			t.Log("Chunk data is empty (may be expected if no events in chunk)")
		}

		// Stop main trace
		data, err := tracing.Stop(bt.ctx, nil)
		if err != nil {
			t.Fatalf("Failed to stop tracing: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected non-empty trace data")
		}
	})

	t.Run("MultipleChunks", func(t *testing.T) {
		tracing := bt.pilot.Tracing()

		err := tracing.Start(bt.ctx, &webpilot.TracingStartOptions{
			Screenshots: true,
		})
		if err != nil {
			t.Fatalf("Failed to start tracing: %v", err)
		}

		bt.go_("https://example.com")

		// First chunk
		err = tracing.StartChunk(bt.ctx, &webpilot.TracingChunkOptions{Name: "chunk1"})
		if err != nil {
			t.Fatalf("Failed to start chunk 1: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
		_, err = tracing.StopChunk(bt.ctx, nil)
		if err != nil {
			t.Fatalf("Failed to stop chunk 1: %v", err)
		}

		// Second chunk
		err = tracing.StartChunk(bt.ctx, &webpilot.TracingChunkOptions{Name: "chunk2"})
		if err != nil {
			t.Fatalf("Failed to start chunk 2: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
		_, err = tracing.StopChunk(bt.ctx, nil)
		if err != nil {
			t.Fatalf("Failed to stop chunk 2: %v", err)
		}

		// Stop main trace
		data, err := tracing.Stop(bt.ctx, nil)
		if err != nil {
			t.Fatalf("Failed to stop tracing: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected non-empty trace data")
		}
	})
}

// TestTracingGroups tests trace group functionality.
func TestTracingGroups(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	t.Run("GroupRecording", func(t *testing.T) {
		tracing := bt.pilot.Tracing()

		err := tracing.Start(bt.ctx, &webpilot.TracingStartOptions{
			Screenshots: true,
			Title:       "Grouped Trace",
		})
		if err != nil {
			t.Fatalf("Failed to start tracing: %v", err)
		}

		bt.go_("https://example.com")

		// Start a group
		err = tracing.StartGroup(bt.ctx, "Login Flow", &webpilot.TracingGroupOptions{
			Location: "test_file.go:123",
		})
		if err != nil {
			t.Fatalf("Failed to start group: %v", err)
		}

		// Some activity in the group
		time.Sleep(200 * time.Millisecond)

		// Stop group
		err = tracing.StopGroup(bt.ctx)
		if err != nil {
			t.Fatalf("Failed to stop group: %v", err)
		}

		// Stop trace
		data, err := tracing.Stop(bt.ctx, nil)
		if err != nil {
			t.Fatalf("Failed to stop tracing: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected non-empty trace data")
		}
	})

	t.Run("NestedGroups", func(t *testing.T) {
		tracing := bt.pilot.Tracing()

		err := tracing.Start(bt.ctx, nil)
		if err != nil {
			t.Fatalf("Failed to start tracing: %v", err)
		}

		bt.go_("https://example.com")

		// Outer group
		err = tracing.StartGroup(bt.ctx, "Outer Group", nil)
		if err != nil {
			t.Fatalf("Failed to start outer group: %v", err)
		}

		// Inner group
		err = tracing.StartGroup(bt.ctx, "Inner Group", nil)
		if err != nil {
			t.Fatalf("Failed to start inner group: %v", err)
		}

		time.Sleep(100 * time.Millisecond)

		err = tracing.StopGroup(bt.ctx) // Stop inner
		if err != nil {
			t.Fatalf("Failed to stop inner group: %v", err)
		}

		err = tracing.StopGroup(bt.ctx) // Stop outer
		if err != nil {
			t.Fatalf("Failed to stop outer group: %v", err)
		}

		data, err := tracing.Stop(bt.ctx, nil)
		if err != nil {
			t.Fatalf("Failed to stop tracing: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected non-empty trace data")
		}
	})
}

// TestTracingFromContext tests tracing via BrowserContext.
func TestTracingFromContext(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	t.Run("ContextTracing", func(t *testing.T) {
		// Create a browser context
		browserCtx, err := bt.pilot.NewContext(bt.ctx)
		if err != nil {
			t.Fatalf("Failed to create context: %v", err)
		}

		// Get tracing from context
		tracing := browserCtx.Tracing()
		if tracing == nil {
			t.Fatal("BrowserContext.Tracing() returned nil")
		}

		err = tracing.Start(bt.ctx, &webpilot.TracingStartOptions{
			Title: "Context Trace",
		})
		if err != nil {
			t.Fatalf("Failed to start tracing: %v", err)
		}

		time.Sleep(200 * time.Millisecond)

		data, err := tracing.Stop(bt.ctx, nil)
		if err != nil {
			t.Fatalf("Failed to stop tracing: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected non-empty trace data")
		}
	})
}
