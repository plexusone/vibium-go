// Package report provides test result tracking and report generation.
package report

import (
	"time"
)

// Status represents the validation status.
type Status string

const (
	StatusGo   Status = "GO"
	StatusWarn Status = "WARN"
	StatusNoGo Status = "NO-GO"
	StatusSkip Status = "SKIP"
)

// Severity represents the impact level of a result.
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
	SeverityInfo     Severity = "info"
)

// StepResult represents the result of executing a single test step.
type StepResult struct {
	// ID is a unique identifier for this step.
	ID string `json:"id"`

	// Action is the tool/action name (e.g., "click", "navigate").
	Action string `json:"action"`

	// Args are the arguments passed to the action.
	Args map[string]any `json:"args,omitempty"`

	// Status is the step status (GO, WARN, NO-GO, SKIP).
	Status Status `json:"status"`

	// Severity is the impact level (critical, high, medium, low, info).
	Severity Severity `json:"severity,omitempty"`

	// DurationMS is the step execution time in milliseconds.
	DurationMS int64 `json:"duration_ms"`

	// Result holds success result data.
	Result any `json:"result,omitempty"`

	// Error holds error details on failure.
	Error *StepError `json:"error,omitempty"`

	// Context holds page state at execution.
	Context *StepContext `json:"context,omitempty"`

	// Console holds browser console log entries.
	Console []ConsoleEntry `json:"console_logs,omitempty"`

	// Network holds failed network requests.
	Network []NetworkError `json:"network_errors,omitempty"`

	// Screenshot holds screenshot reference.
	Screenshot *ScreenshotRef `json:"screenshot,omitempty"`
}

// StepError holds detailed error information.
type StepError struct {
	// Type is the error type name (e.g., "ElementNotFoundError").
	Type string `json:"type"`

	// Message is the full error message.
	Message string `json:"message"`

	// Selector is the CSS selector that failed (if applicable).
	Selector string `json:"selector,omitempty"`

	// TimeoutMS is the timeout that was exceeded (if applicable).
	TimeoutMS int64 `json:"timeout_ms,omitempty"`

	// Suggestions are alternative selectors or fixes.
	Suggestions []string `json:"suggestions,omitempty"`
}

// StepContext holds page state at the time of execution.
type StepContext struct {
	// PageURL is the current page URL.
	PageURL string `json:"page_url"`

	// PageTitle is the current page title.
	PageTitle string `json:"page_title"`

	// VisibleButtons lists visible interactive elements.
	VisibleButtons []string `json:"visible_buttons,omitempty"`

	// DOMSnippet is a relevant DOM fragment.
	DOMSnippet string `json:"dom_snippet,omitempty"`
}

// ConsoleEntry represents a browser console log entry.
type ConsoleEntry struct {
	// Level is the log level (error, warn, info, log).
	Level string `json:"level"`

	// Message is the log message.
	Message string `json:"message"`

	// Source is the log source (javascript, network).
	Source string `json:"source"`

	// URL is the source URL (if applicable).
	URL string `json:"url,omitempty"`
}

// NetworkError represents a failed network request.
type NetworkError struct {
	// URL is the request URL.
	URL string `json:"url"`

	// Method is the HTTP method.
	Method string `json:"method"`

	// StatusCode is the HTTP status code.
	StatusCode int `json:"status"`
}

// ScreenshotRef holds a reference to a screenshot.
type ScreenshotRef struct {
	// Path is the file path (if saved to disk).
	Path string `json:"path,omitempty"`

	// Base64 is the base64-encoded image data.
	Base64 string `json:"base64,omitempty"`
}

// BrowserInfo holds browser information.
type BrowserInfo struct {
	// Name is the browser name (e.g., "chromium").
	Name string `json:"name"`

	// Headless indicates if running in headless mode.
	Headless bool `json:"headless"`

	// Viewport holds viewport dimensions.
	Viewport struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"viewport"`
}

// TestResult holds the complete test execution results.
type TestResult struct {
	// TestPlan is the source test plan file (if applicable).
	TestPlan string `json:"test_plan,omitempty"`

	// Project is the project name.
	Project string `json:"project"`

	// Target is the test target description.
	Target string `json:"target"`

	// Status is the overall status.
	Status Status `json:"status"`

	// DurationMS is the total execution time in milliseconds.
	DurationMS int64 `json:"duration_ms"`

	// Browser holds browser information.
	Browser BrowserInfo `json:"browser"`

	// Steps holds the individual step results.
	Steps []StepResult `json:"steps"`

	// Recommendations are AI-friendly fix suggestions.
	Recommendations []string `json:"recommendations,omitempty"`

	// GeneratedAt is when the report was generated.
	GeneratedAt time.Time `json:"generated_at"`
}

// ComputeOverallStatus computes the overall status from steps.
func ComputeOverallStatus(steps []StepResult) Status {
	hasNoGo := false
	hasWarn := false
	allSkipped := true

	for _, s := range steps {
		if s.Status != StatusSkip {
			allSkipped = false
		}
		switch s.Status {
		case StatusNoGo:
			hasNoGo = true
		case StatusWarn:
			hasWarn = true
		}
	}

	if allSkipped && len(steps) > 0 {
		return StatusSkip
	}
	if hasNoGo {
		return StatusNoGo
	}
	if hasWarn {
		return StatusWarn
	}
	return StatusGo
}

// ComputeTotalDuration computes the total duration from steps.
func ComputeTotalDuration(steps []StepResult) int64 {
	var total int64
	for _, s := range steps {
		total += s.DurationMS
	}
	return total
}
