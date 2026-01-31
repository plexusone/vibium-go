// Package vibium provides a Go client for browser automation via the Vibium platform.
// It communicates with the Vibium clicker binary over WebSocket using the WebDriver BiDi protocol.
package vibium

import "time"

// BoundingBox represents the position and size of an element.
type BoundingBox struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// ElementInfo contains metadata about a DOM element.
type ElementInfo struct {
	Tag  string      `json:"tag"`
	Text string      `json:"text"`
	Box  BoundingBox `json:"box"`
}

// LaunchOptions configures browser launch behavior.
type LaunchOptions struct {
	// Headless runs the browser without a visible window.
	Headless bool

	// Port specifies the WebSocket port. If 0, an available port is auto-selected.
	Port int

	// ExecutablePath specifies a custom path to the clicker binary.
	ExecutablePath string
}

// FindOptions configures element finding behavior.
type FindOptions struct {
	// Timeout specifies how long to wait for the element to appear.
	// Default is 30 seconds.
	Timeout time.Duration
}

// ActionOptions configures action behavior (click, type).
type ActionOptions struct {
	// Timeout specifies how long to wait for actionability.
	// Default is 30 seconds.
	Timeout time.Duration
}

// DefaultTimeout is the default timeout for finding elements and waiting for actionability.
const DefaultTimeout = 30 * time.Second
