package report

import (
	"encoding/json"
	"fmt"
)

// DiagnosticReport is the full diagnostic report for agent consumption.
// It contains all the details needed for an AI agent to understand
// what happened and why a test failed.
type DiagnosticReport struct {
	TestResult
}

// NewDiagnosticReport creates a DiagnosticReport from a TestResult.
func NewDiagnosticReport(tr *TestResult) *DiagnosticReport {
	return &DiagnosticReport{TestResult: *tr}
}

// JSON serializes the diagnostic report to JSON.
func (r *DiagnosticReport) JSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

// GenerateRecommendations analyzes the test results and generates
// AI-friendly recommendations for fixing failures.
func (r *DiagnosticReport) GenerateRecommendations() {
	recommendations := make([]string, 0)

	for _, step := range r.Steps {
		if step.Status != StatusNoGo || step.Error == nil {
			continue
		}

		// Generate recommendations based on error type
		switch step.Error.Type {
		case "ElementNotFoundError":
			if len(step.Error.Suggestions) > 0 {
				recommendations = append(recommendations,
					"Selector '"+step.Error.Selector+"' not found. Try: "+step.Error.Suggestions[0])
			} else {
				recommendations = append(recommendations,
					"Selector '"+step.Error.Selector+"' not found. Check if the element exists or wait for it to load.")
			}

		case "TimeoutError":
			recommendations = append(recommendations,
				"Operation timed out. Consider increasing the timeout or checking if the page is loading correctly.")

		case "NavigationError":
			recommendations = append(recommendations,
				"Navigation failed. Check if the URL is correct and the server is responding.")

		case "ClickError":
			recommendations = append(recommendations,
				"Click failed. The element may be obscured, not interactable, or outside the viewport.")
		}

		// Add network error recommendations
		if len(step.Network) > 0 {
			for _, ne := range step.Network {
				if ne.StatusCode == 404 {
					recommendations = append(recommendations,
						"Network error: "+ne.URL+" returned 404. This missing resource may affect page functionality.")
				} else if ne.StatusCode >= 500 {
					recommendations = append(recommendations,
						fmt.Sprintf("Server error: %s returned %d. The server may be experiencing issues.", ne.URL, ne.StatusCode))
				}
			}
		}

		// Add console error recommendations
		for _, ce := range step.Console {
			if ce.Level == "error" {
				recommendations = append(recommendations,
					"JavaScript error detected: "+truncate(ce.Message, 100))
			}
		}
	}

	r.Recommendations = recommendations
}

// truncate shortens a string to maxLen, adding "..." if truncated.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
