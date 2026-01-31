package vibium

import (
	"errors"
	"fmt"
)

var (
	// ErrConnectionFailed is returned when WebSocket connection fails.
	ErrConnectionFailed = errors.New("failed to connect to clicker server")

	// ErrElementNotFound is returned when an element cannot be found.
	ErrElementNotFound = errors.New("element not found")

	// ErrBrowserCrashed is returned when the browser process exits unexpectedly.
	ErrBrowserCrashed = errors.New("browser crashed")

	// ErrClickerNotFound is returned when the clicker binary cannot be found.
	ErrClickerNotFound = errors.New("clicker binary not found")

	// ErrTimeout is returned when an operation times out.
	ErrTimeout = errors.New("operation timed out")

	// ErrConnectionClosed is returned when the WebSocket connection is closed.
	ErrConnectionClosed = errors.New("connection closed")
)

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
	Selector string
	Timeout  int64 // milliseconds
	Reason   string
}

func (e *TimeoutError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("timeout after %dms waiting for '%s': %s", e.Timeout, e.Selector, e.Reason)
	}
	return fmt.Sprintf("timeout after %dms waiting for '%s'", e.Timeout, e.Selector)
}

// ElementNotFoundError represents an element that could not be found.
type ElementNotFoundError struct {
	Selector string
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
