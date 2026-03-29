package w3pilot

import (
	"errors"
	"fmt"
)

var (
	// ErrConnectionFailed is returned when WebSocket connection fails.
	ErrConnectionFailed = errors.New("failed to connect to browser")

	// ErrElementNotFound is returned when an element cannot be found.
	ErrElementNotFound = errors.New("element not found")

	// ErrBrowserCrashed is returned when the browser process exits unexpectedly.
	ErrBrowserCrashed = errors.New("browser crashed")

	// ErrBrowserNotFound is returned when Chrome cannot be found.
	ErrBrowserNotFound = errors.New("Chrome not found")

	// ErrClickerNotFound is deprecated: use ErrBrowserNotFound instead.
	// Deprecated: This error is no longer used. Use ErrBrowserNotFound.
	ErrClickerNotFound = ErrBrowserNotFound

	// ErrTimeout is returned when an operation times out.
	ErrTimeout = errors.New("operation timed out")

	// ErrConnectionClosed is returned when the WebSocket connection is closed.
	ErrConnectionClosed = errors.New("connection closed")
)

// PageContext provides context about the page state when an error occurred.
// This helps AI agents understand the situation and recover from errors.
type PageContext struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	VisibleText string `json:"visible_text,omitempty"`
}

// ConnectionError represents a WebSocket connection failure.
type ConnectionError struct {
	URL   string
	Cause error
}

func (e *ConnectionError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("failed to connect to %s: %v", e.URL, e.Cause)
	}
	return fmt.Sprintf("failed to connect to %s", e.URL)
}

func (e *ConnectionError) Unwrap() error {
	return e.Cause
}

// TimeoutError represents a timeout waiting for an element or action.
type TimeoutError struct {
	Selector    string       `json:"selector"`
	Timeout     int64        `json:"timeout_ms"` // milliseconds
	Reason      string       `json:"reason,omitempty"`
	PageContext *PageContext `json:"page_context,omitempty"`
	Suggestions []string     `json:"suggestions,omitempty"`
}

func (e *TimeoutError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("timeout after %dms waiting for '%s': %s", e.Timeout, e.Selector, e.Reason)
	}
	return fmt.Sprintf("timeout after %dms waiting for '%s'", e.Timeout, e.Selector)
}

// ElementNotFoundError represents an element that could not be found.
type ElementNotFoundError struct {
	Selector    string       `json:"selector"`
	PageContext *PageContext `json:"page_context,omitempty"`
	Suggestions []string     `json:"suggestions,omitempty"`
}

func (e *ElementNotFoundError) Error() string {
	return fmt.Sprintf("element not found: %s", e.Selector)
}

// BrowserCrashedError represents an unexpected browser exit.
type BrowserCrashedError struct {
	ExitCode int
	Output   string
}

func (e *BrowserCrashedError) Error() string {
	if e.Output != "" {
		return fmt.Sprintf("browser crashed with exit code %d: %s", e.ExitCode, e.Output)
	}
	return fmt.Sprintf("browser crashed with exit code %d", e.ExitCode)
}

// BiDiError represents an error from the BiDi protocol.
type BiDiError struct {
	ErrorType string
	Message   string
}

func (e *BiDiError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s: %s", e.ErrorType, e.Message)
	}
	return e.ErrorType
}
