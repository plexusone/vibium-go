package mcp

import (
	"context"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	vibium "github.com/plexusone/vibium-go"
	"github.com/plexusone/vibium-go/mcp/report"
)

// VerifyValue tool - verifies that an input element has the expected value

type VerifyValueInput struct {
	Selector  string `json:"selector" jsonschema:"CSS selector for the input element,required"`
	Expected  string `json:"expected" jsonschema:"Expected value to verify,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"Timeout in milliseconds (default: 5000)"`
}

type VerifyValueOutput struct {
	Passed  bool   `json:"passed"`
	Actual  string `json:"actual"`
	Message string `json:"message"`
}

func (s *Server) handleVerifyValue(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input VerifyValueInput,
) (*mcp.CallToolResult, VerifyValueOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, VerifyValueOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	start := time.Now()
	elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})

	result := report.StepResult{
		ID:     s.session.NextStepID("verify_value"),
		Action: "verify_value",
		Args:   map[string]any{"selector": input.Selector, "expected": input.Expected},
	}

	if err != nil {
		result.DurationMS = time.Since(start).Milliseconds()
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityCritical
		result.Error = &report.StepError{
			Type:        "ElementNotFoundError",
			Message:     err.Error(),
			Selector:    input.Selector,
			TimeoutMS:   int64(input.TimeoutMS),
			Suggestions: s.session.FindSimilarSelectors(ctx, input.Selector),
		}
		result.Context = s.session.CaptureContext(ctx)
		result.Screenshot = s.session.CaptureScreenshot(ctx)
		s.session.RecordStep(result)
		return nil, VerifyValueOutput{}, fmt.Errorf("element not found: %s", input.Selector)
	}

	actual, err := elem.Value(ctx)
	result.DurationMS = time.Since(start).Milliseconds()

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityCritical
		result.Error = &report.StepError{
			Type:     "GetValueError",
			Message:  err.Error(),
			Selector: input.Selector,
		}
		result.Screenshot = s.session.CaptureScreenshot(ctx)
		s.session.RecordStep(result)
		return nil, VerifyValueOutput{}, fmt.Errorf("get value failed: %w", err)
	}

	passed := actual == input.Expected

	if !passed {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityCritical
		result.Error = &report.StepError{
			Type:    "VerifyValueFailed",
			Message: fmt.Sprintf("Expected %q but got %q", input.Expected, actual),
		}
		result.Context = s.session.CaptureContext(ctx)
		result.Screenshot = s.session.CaptureScreenshot(ctx)
		s.session.RecordStep(result)

		return nil, VerifyValueOutput{
			Passed:  false,
			Actual:  actual,
			Message: fmt.Sprintf("Value mismatch: expected %q but got %q", input.Expected, actual),
		}, nil
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	result.Result = map[string]any{"actual": actual, "expected": input.Expected}
	s.session.RecordStep(result)

	return nil, VerifyValueOutput{
		Passed:  true,
		Actual:  actual,
		Message: fmt.Sprintf("Value matches: %q", actual),
	}, nil
}

// VerifyListVisible tool - verifies that a list of items are visible on the page

type VerifyListVisibleInput struct {
	Items     []string `json:"items" jsonschema:"List of text items that should be visible on the page,required"`
	Selector  string   `json:"selector" jsonschema:"Optional CSS selector to scope the search"`
	TimeoutMS int      `json:"timeout_ms" jsonschema:"Timeout in milliseconds (default: 5000)"`
}

type VerifyListVisibleOutput struct {
	Passed  bool     `json:"passed"`
	Found   []string `json:"found"`
	Missing []string `json:"missing"`
	Message string   `json:"message"`
}

func (s *Server) handleVerifyListVisible(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input VerifyListVisibleInput,
) (*mcp.CallToolResult, VerifyListVisibleOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, VerifyListVisibleOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}

	start := time.Now()

	result := report.StepResult{
		ID:     s.session.NextStepID("verify_list_visible"),
		Action: "verify_list_visible",
		Args:   map[string]any{"items": input.Items, "selector": input.Selector},
	}

	// Build script to check for each item's visibility
	var found []string
	var missing []string

	for _, item := range input.Items {
		var script string
		if input.Selector != "" {
			script = fmt.Sprintf(`
				(function() {
					const el = document.querySelector(%q);
					return el && el.textContent.includes(%q);
				})()
			`, input.Selector, item)
		} else {
			script = fmt.Sprintf(`document.body.textContent.includes(%q)`, item)
		}

		evalResult, err := vibe.Evaluate(ctx, script)
		if err != nil {
			missing = append(missing, item)
			continue
		}

		if visible, ok := evalResult.(bool); ok && visible {
			found = append(found, item)
		} else {
			missing = append(missing, item)
		}
	}

	result.DurationMS = time.Since(start).Milliseconds()

	passed := len(missing) == 0

	if !passed {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityCritical
		result.Error = &report.StepError{
			Type:    "VerifyListVisibleFailed",
			Message: fmt.Sprintf("Missing items: %v", missing),
		}
		result.Context = s.session.CaptureContext(ctx)
		result.Screenshot = s.session.CaptureScreenshot(ctx)
		s.session.RecordStep(result)

		return nil, VerifyListVisibleOutput{
			Passed:  false,
			Found:   found,
			Missing: missing,
			Message: fmt.Sprintf("Found %d of %d items, missing: %v", len(found), len(input.Items), missing),
		}, nil
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	result.Result = map[string]any{"found": found}
	s.session.RecordStep(result)

	return nil, VerifyListVisibleOutput{
		Passed:  true,
		Found:   found,
		Missing: missing,
		Message: fmt.Sprintf("All %d items visible", len(input.Items)),
	}, nil
}

// GenerateLocator tool - generates a locator string for a given element

type GenerateLocatorInput struct {
	Selector  string `json:"selector" jsonschema:"CSS selector for the element,required"`
	Strategy  string `json:"strategy" jsonschema:"Locator strategy: css xpath testid role text (default: css),enum=css,enum=xpath,enum=testid,enum=role,enum=text"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"Timeout in milliseconds (default: 5000)"`
}

type GenerateLocatorOutput struct {
	Locator  string            `json:"locator"`
	Strategy string            `json:"strategy"`
	Metadata map[string]string `json:"metadata"`
}

func (s *Server) handleGenerateLocator(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GenerateLocatorInput,
) (*mcp.CallToolResult, GenerateLocatorOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, GenerateLocatorOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	if input.Strategy == "" {
		input.Strategy = "css"
	}

	elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})
	if err != nil {
		return nil, GenerateLocatorOutput{}, fmt.Errorf("element not found: %s", input.Selector)
	}

	metadata := make(map[string]string)
	var locator string

	switch input.Strategy {
	case "css":
		// Generate a unique CSS selector
		script := `
			(function(selector) {
				const el = document.querySelector(selector);
				if (!el) return null;

				// Try to generate a unique selector
				// Priority: id > data-testid > class combination > tag with index

				if (el.id) {
					return '#' + CSS.escape(el.id);
				}

				if (el.dataset.testid) {
					return '[data-testid="' + el.dataset.testid + '"]';
				}

				// Generate path from element
				let path = [];
				let current = el;
				while (current && current.nodeType === Node.ELEMENT_NODE) {
					let selector = current.tagName.toLowerCase();
					if (current.id) {
						selector = '#' + CSS.escape(current.id);
						path.unshift(selector);
						break;
					}

					let sibling = current;
					let nth = 1;
					while (sibling = sibling.previousElementSibling) {
						if (sibling.tagName === current.tagName) nth++;
					}

					if (nth > 1 || current.nextElementSibling?.tagName === current.tagName) {
						selector += ':nth-of-type(' + nth + ')';
					}

					path.unshift(selector);
					current = current.parentElement;
				}

				return path.join(' > ');
			})(%q)
		`
		result, err := vibe.Evaluate(ctx, fmt.Sprintf(script, input.Selector))
		if err != nil {
			return nil, GenerateLocatorOutput{}, fmt.Errorf("generate locator failed: %w", err)
		}
		if result != nil {
			locator = fmt.Sprintf("%v", result)
		} else {
			locator = input.Selector
		}

	case "xpath":
		// Generate XPath for the element
		script := `
			(function(selector) {
				const el = document.querySelector(selector);
				if (!el) return null;

				if (el.id) {
					return '//*[@id="' + el.id + '"]';
				}

				let path = [];
				let current = el;
				while (current && current.nodeType === Node.ELEMENT_NODE) {
					let tag = current.tagName.toLowerCase();
					let sibling = current;
					let index = 1;
					while (sibling = sibling.previousElementSibling) {
						if (sibling.tagName.toLowerCase() === tag) index++;
					}
					path.unshift(tag + '[' + index + ']');
					current = current.parentElement;
				}

				return '/' + path.join('/');
			})(%q)
		`
		result, err := vibe.Evaluate(ctx, fmt.Sprintf(script, input.Selector))
		if err != nil {
			return nil, GenerateLocatorOutput{}, fmt.Errorf("generate locator failed: %w", err)
		}
		if result != nil {
			locator = fmt.Sprintf("%v", result)
		}

	case "testid":
		testID, err := elem.GetAttribute(ctx, "data-testid")
		if err != nil {
			return nil, GenerateLocatorOutput{}, fmt.Errorf("get testid failed: %w", err)
		}
		if testID == "" {
			return nil, GenerateLocatorOutput{}, fmt.Errorf("element has no data-testid attribute")
		}
		locator = fmt.Sprintf("[data-testid=\"%s\"]", testID)
		metadata["testid"] = testID

	case "role":
		role, err := elem.Role(ctx)
		if err != nil {
			return nil, GenerateLocatorOutput{}, fmt.Errorf("get role failed: %w", err)
		}
		if role == "" {
			return nil, GenerateLocatorOutput{}, fmt.Errorf("element has no ARIA role")
		}
		label, _ := elem.Label(ctx)
		if label != "" {
			locator = fmt.Sprintf("role=%s[name=%q]", role, label)
			metadata["label"] = label
		} else {
			locator = fmt.Sprintf("role=%s", role)
		}
		metadata["role"] = role

	case "text":
		text, err := elem.Text(ctx)
		if err != nil {
			return nil, GenerateLocatorOutput{}, fmt.Errorf("get text failed: %w", err)
		}
		if text == "" {
			return nil, GenerateLocatorOutput{}, fmt.Errorf("element has no text content")
		}
		// Truncate long text
		if len(text) > 50 {
			text = text[:50]
		}
		locator = fmt.Sprintf("text=%q", text)
		metadata["text"] = text

	default:
		return nil, GenerateLocatorOutput{}, fmt.Errorf("unknown strategy: %s", input.Strategy)
	}

	return nil, GenerateLocatorOutput{
		Locator:  locator,
		Strategy: input.Strategy,
		Metadata: metadata,
	}, nil
}

// WaitForSelector tool - waits for an element to reach a specific state

type WaitForSelectorInput struct {
	Selector  string `json:"selector" jsonschema:"CSS selector for the element,required"`
	State     string `json:"state" jsonschema:"State to wait for: attached detached visible hidden (default: visible),enum=attached,enum=detached,enum=visible,enum=hidden"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"Timeout in milliseconds (default: 30000)"`
}

type WaitForSelectorOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleWaitForSelector(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input WaitForSelectorInput,
) (*mcp.CallToolResult, WaitForSelectorOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, WaitForSelectorOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 30000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	if input.State == "" {
		input.State = "visible"
	}

	// Find the element first (for attached/visible states) or wait for condition
	switch input.State {
	case "attached", "visible":
		elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})
		if err != nil {
			return nil, WaitForSelectorOutput{}, fmt.Errorf("wait for selector failed: %w", err)
		}
		if input.State == "visible" {
			err = elem.WaitUntil(ctx, "visible", timeout)
			if err != nil {
				return nil, WaitForSelectorOutput{}, fmt.Errorf("wait for visible failed: %w", err)
			}
		}
	case "detached", "hidden":
		err = vibe.WaitForFunction(ctx, fmt.Sprintf(`() => {
			const el = document.querySelector(%q);
			if (%q === "detached") return el === null;
			if (el === null) return true;
			const style = window.getComputedStyle(el);
			return style.display === 'none' || style.visibility === 'hidden' || el.offsetParent === null;
		}`, input.Selector, input.State), timeout)
		if err != nil {
			return nil, WaitForSelectorOutput{}, fmt.Errorf("wait for %s failed: %w", input.State, err)
		}
	default:
		return nil, WaitForSelectorOutput{}, fmt.Errorf("invalid state: %s", input.State)
	}

	return nil, WaitForSelectorOutput{
		Message: fmt.Sprintf("Element %s is %s", input.Selector, input.State),
	}, nil
}

// VerifyText tool - verifies element text matches expected value

type VerifyTextInput struct {
	Selector  string `json:"selector" jsonschema:"CSS selector for the element,required"`
	Expected  string `json:"expected" jsonschema:"Expected text content,required"`
	Exact     bool   `json:"exact" jsonschema:"Require exact match (default: false uses contains)"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"Timeout in milliseconds (default: 5000)"`
}

type VerifyTextOutput struct {
	Passed  bool   `json:"passed"`
	Actual  string `json:"actual"`
	Message string `json:"message"`
}

func (s *Server) handleVerifyText(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input VerifyTextInput,
) (*mcp.CallToolResult, VerifyTextOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, VerifyTextOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})
	if err != nil {
		return nil, VerifyTextOutput{}, fmt.Errorf("element not found: %s", input.Selector)
	}

	actual, err := elem.Text(ctx)
	if err != nil {
		return nil, VerifyTextOutput{}, fmt.Errorf("get text failed: %w", err)
	}

	var passed bool
	if input.Exact {
		passed = actual == input.Expected
	} else {
		passed = contains(actual, input.Expected)
	}

	if !passed {
		matchType := "contain"
		if input.Exact {
			matchType = "equal"
		}
		return nil, VerifyTextOutput{
			Passed:  false,
			Actual:  actual,
			Message: fmt.Sprintf("Text does not %s expected: got %q, expected %q", matchType, actual, input.Expected),
		}, nil
	}

	return nil, VerifyTextOutput{
		Passed:  true,
		Actual:  actual,
		Message: fmt.Sprintf("Text matches: %q", actual),
	}, nil
}

// contains checks if s contains substr (case-sensitive)
func contains(s, substr string) bool {
	return len(substr) == 0 || (len(s) >= len(substr) && searchString(s, substr))
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// VerifyVisible tool - verifies element is visible

type VerifyVisibleInput struct {
	Selector  string `json:"selector" jsonschema:"CSS selector for the element,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"Timeout in milliseconds (default: 5000)"`
}

type VerifyVisibleOutput struct {
	Passed  bool   `json:"passed"`
	Message string `json:"message"`
}

func (s *Server) handleVerifyVisible(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input VerifyVisibleInput,
) (*mcp.CallToolResult, VerifyVisibleOutput, error) {
	passed, msg, err := s.verifyElementState(ctx, input.Selector, input.TimeoutMS, "visible")
	if err != nil {
		return nil, VerifyVisibleOutput{}, err
	}
	return nil, VerifyVisibleOutput{Passed: passed, Message: msg}, nil
}

// VerifyEnabled tool - verifies element is enabled

type VerifyEnabledInput struct {
	Selector  string `json:"selector" jsonschema:"CSS selector for the element,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"Timeout in milliseconds (default: 5000)"`
}

type VerifyEnabledOutput struct {
	Passed  bool   `json:"passed"`
	Message string `json:"message"`
}

func (s *Server) handleVerifyEnabled(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input VerifyEnabledInput,
) (*mcp.CallToolResult, VerifyEnabledOutput, error) {
	passed, msg, err := s.verifyElementState(ctx, input.Selector, input.TimeoutMS, "enabled")
	if err != nil {
		return nil, VerifyEnabledOutput{}, err
	}
	return nil, VerifyEnabledOutput{Passed: passed, Message: msg}, nil
}

// verifyElementState is a helper that verifies an element's state (visible, hidden, enabled, disabled)
func (s *Server) verifyElementState(ctx context.Context, selector string, timeoutMS int, state string) (bool, string, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return false, "", fmt.Errorf("browser not available: %w", err)
	}

	if timeoutMS == 0 {
		timeoutMS = 5000
	}
	timeout := time.Duration(timeoutMS) * time.Millisecond

	elem, err := vibe.Find(ctx, selector, &vibium.FindOptions{Timeout: timeout})
	if err != nil {
		// For hidden state, element not found could be valid
		if state == "hidden" {
			return true, fmt.Sprintf("Element is hidden (not found): %s", selector), nil
		}
		return false, fmt.Sprintf("Element not found: %s", selector), nil
	}

	var checkResult bool
	var checkErr error
	var expectTrue bool = true // Whether we expect the check to return true

	switch state {
	case "visible":
		checkResult, checkErr = elem.IsVisible(ctx)
	case "hidden":
		checkResult, checkErr = elem.IsVisible(ctx)
		expectTrue = false // For hidden, we expect IsVisible to return false
	case "enabled":
		checkResult, checkErr = elem.IsEnabled(ctx)
	case "disabled":
		checkResult, checkErr = elem.IsEnabled(ctx)
		expectTrue = false // For disabled, we expect IsEnabled to return false
	default:
		return false, "", fmt.Errorf("unknown state: %s", state)
	}

	if checkErr != nil {
		return false, "", fmt.Errorf("check %s failed: %w", state, checkErr)
	}

	// Determine pass/fail based on expected state
	passed := checkResult == expectTrue

	if !passed {
		oppositeState := map[string]string{
			"visible":  "hidden",
			"hidden":   "visible",
			"enabled":  "disabled",
			"disabled": "enabled",
		}
		return false, fmt.Sprintf("Element is %s, expected %s: %s", oppositeState[state], state, selector), nil
	}

	return true, fmt.Sprintf("Element is %s: %s", state, selector), nil
}

// VerifyHidden tool - verifies element is hidden

type VerifyHiddenInput struct {
	Selector  string `json:"selector" jsonschema:"CSS selector for the element,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"Timeout in milliseconds (default: 5000)"`
}

type VerifyHiddenOutput struct {
	Passed  bool   `json:"passed"`
	Message string `json:"message"`
}

func (s *Server) handleVerifyHidden(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input VerifyHiddenInput,
) (*mcp.CallToolResult, VerifyHiddenOutput, error) {
	passed, msg, err := s.verifyElementState(ctx, input.Selector, input.TimeoutMS, "hidden")
	if err != nil {
		return nil, VerifyHiddenOutput{}, err
	}
	return nil, VerifyHiddenOutput{Passed: passed, Message: msg}, nil
}

// VerifyDisabled tool - verifies element is disabled

type VerifyDisabledInput struct {
	Selector  string `json:"selector" jsonschema:"CSS selector for the element,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"Timeout in milliseconds (default: 5000)"`
}

type VerifyDisabledOutput struct {
	Passed  bool   `json:"passed"`
	Message string `json:"message"`
}

func (s *Server) handleVerifyDisabled(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input VerifyDisabledInput,
) (*mcp.CallToolResult, VerifyDisabledOutput, error) {
	passed, msg, err := s.verifyElementState(ctx, input.Selector, input.TimeoutMS, "disabled")
	if err != nil {
		return nil, VerifyDisabledOutput{}, err
	}
	return nil, VerifyDisabledOutput{Passed: passed, Message: msg}, nil
}

// VerifyChecked tool - verifies checkbox/radio is checked

type VerifyCheckedInput struct {
	Selector  string `json:"selector" jsonschema:"CSS selector for the checkbox/radio element,required"`
	Checked   bool   `json:"checked" jsonschema:"Expected checked state (default: true)"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"Timeout in milliseconds (default: 5000)"`
}

type VerifyCheckedOutput struct {
	Passed  bool   `json:"passed"`
	Checked bool   `json:"checked"`
	Message string `json:"message"`
}

func (s *Server) handleVerifyChecked(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input VerifyCheckedInput,
) (*mcp.CallToolResult, VerifyCheckedOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, VerifyCheckedOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	// Default to expecting checked=true if not specified
	// Note: JSON unmarshaling will set Checked to false if not provided,
	// so we check if this is the first call by looking at the raw input
	expectedChecked := input.Checked

	elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})
	if err != nil {
		return nil, VerifyCheckedOutput{
			Passed:  false,
			Message: fmt.Sprintf("Element not found: %s", input.Selector),
		}, nil
	}

	actualChecked, err := elem.IsChecked(ctx)
	if err != nil {
		return nil, VerifyCheckedOutput{}, fmt.Errorf("check checked state failed: %w", err)
	}

	passed := actualChecked == expectedChecked

	if !passed {
		return nil, VerifyCheckedOutput{
			Passed:  false,
			Checked: actualChecked,
			Message: fmt.Sprintf("Element checked state is %v, expected %v", actualChecked, expectedChecked),
		}, nil
	}

	return nil, VerifyCheckedOutput{
		Passed:  true,
		Checked: actualChecked,
		Message: fmt.Sprintf("Element checked state is %v as expected", actualChecked),
	}, nil
}
