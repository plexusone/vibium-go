package mcp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	vibium "github.com/agentplexus/vibium-go"
	"github.com/agentplexus/vibium-go/mcp/report"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Tool input/output types

type BrowserLaunchInput struct {
	Headless bool `json:"headless" jsonschema:"description=Run browser without GUI (default: true)"`
}

type BrowserLaunchOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleBrowserLaunch(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input BrowserLaunchInput,
) (*mcp.CallToolResult, BrowserLaunchOutput, error) {
	// Default to configured headless mode, but allow override
	s.session.config.Headless = input.Headless

	start := time.Now()
	err := s.session.LaunchIfNeeded(ctx)
	duration := time.Since(start)

	result := report.StepResult{
		ID:         s.session.NextStepID("browser_launch"),
		Action:     "browser_launch",
		Args:       map[string]any{"headless": input.Headless},
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityCritical
		result.Error = &report.StepError{
			Type:    "LaunchError",
			Message: err.Error(),
		}
		s.session.RecordStep(result)
		return nil, BrowserLaunchOutput{}, fmt.Errorf("failed to launch browser: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	return nil, BrowserLaunchOutput{Message: "Browser launched successfully"}, nil
}

type BrowserQuitInput struct{}

type BrowserQuitOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleBrowserQuit(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input BrowserQuitInput,
) (*mcp.CallToolResult, BrowserQuitOutput, error) {
	start := time.Now()
	err := s.session.Close(ctx)
	duration := time.Since(start)

	result := report.StepResult{
		ID:         s.session.NextStepID("browser_quit"),
		Action:     "browser_quit",
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:    "QuitError",
			Message: err.Error(),
		}
		s.session.RecordStep(result)
		return nil, BrowserQuitOutput{}, fmt.Errorf("failed to quit browser: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	return nil, BrowserQuitOutput{Message: "Browser closed successfully"}, nil
}

type NavigateInput struct {
	URL string `json:"url" jsonschema:"description=The URL to navigate to,required"`
}

type NavigateOutput struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

func (s *Server) handleNavigate(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input NavigateInput,
) (*mcp.CallToolResult, NavigateOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, NavigateOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	start := time.Now()
	err = vibe.Go(ctx, input.URL)
	duration := time.Since(start)

	result := report.StepResult{
		ID:         s.session.NextStepID("navigate"),
		Action:     "navigate",
		Args:       map[string]any{"url": input.URL},
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityCritical
		result.Error = &report.StepError{
			Type:    "NavigationError",
			Message: err.Error(),
		}
		result.Screenshot = s.session.CaptureScreenshot(ctx)
		s.session.RecordStep(result)
		return nil, NavigateOutput{}, fmt.Errorf("navigation failed: %w", err)
	}

	currentURL, _ := vibe.URL(ctx)
	currentTitle, _ := vibe.Title(ctx)

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	result.Result = map[string]any{
		"url":   currentURL,
		"title": currentTitle,
	}
	s.session.RecordStep(result)

	// Record for script export
	s.session.Recorder().RecordNavigate(input.URL)

	return nil, NavigateOutput{
		URL:   currentURL,
		Title: currentTitle,
	}, nil
}

type ClickInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the element to click,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type ClickOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleClick(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ClickInput,
) (*mcp.CallToolResult, ClickOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, ClickOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	start := time.Now()
	elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})

	result := report.StepResult{
		ID:     s.session.NextStepID("click"),
		Action: "click",
		Args:   map[string]any{"selector": input.Selector},
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
		return nil, ClickOutput{}, fmt.Errorf("element not found: %s", input.Selector)
	}

	err = elem.Click(ctx, &vibium.ActionOptions{Timeout: timeout})
	result.DurationMS = time.Since(start).Milliseconds()

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityCritical
		result.Error = &report.StepError{
			Type:     "ClickError",
			Message:  err.Error(),
			Selector: input.Selector,
		}
		result.Screenshot = s.session.CaptureScreenshot(ctx)
		s.session.RecordStep(result)
		return nil, ClickOutput{}, fmt.Errorf("click failed: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	// Record for script export
	s.session.Recorder().RecordClick(input.Selector)

	return nil, ClickOutput{Message: fmt.Sprintf("Clicked %s", input.Selector)}, nil
}

type TypeInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the input element,required"`
	Text      string `json:"text" jsonschema:"description=Text to type,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type TypeOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleType(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input TypeInput,
) (*mcp.CallToolResult, TypeOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, TypeOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	start := time.Now()
	elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})

	result := report.StepResult{
		ID:     s.session.NextStepID("type"),
		Action: "type",
		Args:   map[string]any{"selector": input.Selector, "text": input.Text},
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
		return nil, TypeOutput{}, fmt.Errorf("element not found: %s", input.Selector)
	}

	err = elem.Type(ctx, input.Text, &vibium.ActionOptions{Timeout: timeout})
	result.DurationMS = time.Since(start).Milliseconds()

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityCritical
		result.Error = &report.StepError{
			Type:     "TypeError",
			Message:  err.Error(),
			Selector: input.Selector,
		}
		result.Screenshot = s.session.CaptureScreenshot(ctx)
		s.session.RecordStep(result)
		return nil, TypeOutput{}, fmt.Errorf("type failed: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	// Record for script export
	s.session.Recorder().RecordType(input.Selector, input.Text)

	return nil, TypeOutput{Message: fmt.Sprintf("Typed into %s", input.Selector)}, nil
}

type GetTextInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the element,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type GetTextOutput struct {
	Text string `json:"text"`
}

func (s *Server) handleGetText(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetTextInput,
) (*mcp.CallToolResult, GetTextOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, GetTextOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	start := time.Now()
	elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})

	result := report.StepResult{
		ID:     s.session.NextStepID("get_text"),
		Action: "get_text",
		Args:   map[string]any{"selector": input.Selector},
	}

	if err != nil {
		result.DurationMS = time.Since(start).Milliseconds()
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:        "ElementNotFoundError",
			Message:     err.Error(),
			Selector:    input.Selector,
			TimeoutMS:   int64(input.TimeoutMS),
			Suggestions: s.session.FindSimilarSelectors(ctx, input.Selector),
		}
		s.session.RecordStep(result)
		return nil, GetTextOutput{}, fmt.Errorf("element not found: %s", input.Selector)
	}

	text, err := elem.Text(ctx)
	result.DurationMS = time.Since(start).Milliseconds()

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:     "GetTextError",
			Message:  err.Error(),
			Selector: input.Selector,
		}
		s.session.RecordStep(result)
		return nil, GetTextOutput{}, fmt.Errorf("get text failed: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	result.Result = map[string]any{"text": text}
	s.session.RecordStep(result)

	return nil, GetTextOutput{Text: text}, nil
}

type ScreenshotInput struct {
	Format string `json:"format" jsonschema:"description=Output format: base64 (default) or file,enum=base64,enum=file"`
	Path   string `json:"path" jsonschema:"description=File path (required if format is file)"`
}

type ScreenshotOutput struct {
	Format string `json:"format"`
	Data   string `json:"data,omitempty"`
	Path   string `json:"path,omitempty"`
}

func (s *Server) handleScreenshot(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ScreenshotInput,
) (*mcp.CallToolResult, ScreenshotOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, ScreenshotOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.Format == "" {
		input.Format = "base64"
	}

	start := time.Now()
	data, err := vibe.Screenshot(ctx)
	duration := time.Since(start)

	result := report.StepResult{
		ID:         s.session.NextStepID("screenshot"),
		Action:     "screenshot",
		Args:       map[string]any{"format": input.Format},
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:    "ScreenshotError",
			Message: err.Error(),
		}
		s.session.RecordStep(result)
		return nil, ScreenshotOutput{}, fmt.Errorf("screenshot failed: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	// Record for script export
	s.session.Recorder().RecordScreenshot("screenshot.png", false)

	output := ScreenshotOutput{Format: input.Format}
	if input.Format == "base64" {
		output.Data = base64.StdEncoding.EncodeToString(data)
	}
	// TODO: Handle file format

	return nil, output, nil
}

type GetTitleInput struct{}

type GetTitleOutput struct {
	Title string `json:"title"`
}

func (s *Server) handleGetTitle(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetTitleInput,
) (*mcp.CallToolResult, GetTitleOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, GetTitleOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	title, err := vibe.Title(ctx)
	if err != nil {
		return nil, GetTitleOutput{}, fmt.Errorf("get title failed: %w", err)
	}
	return nil, GetTitleOutput{Title: title}, nil
}

type GetURLInput struct{}

type GetURLOutput struct {
	URL string `json:"url"`
}

func (s *Server) handleGetURL(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetURLInput,
) (*mcp.CallToolResult, GetURLOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, GetURLOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	url, err := vibe.URL(ctx)
	if err != nil {
		return nil, GetURLOutput{}, fmt.Errorf("get url failed: %w", err)
	}
	return nil, GetURLOutput{URL: url}, nil
}

type EvaluateInput struct {
	Script string `json:"script" jsonschema:"description=JavaScript to execute,required"`
}

type EvaluateOutput struct {
	Result any `json:"result"`
}

func (s *Server) handleEvaluate(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input EvaluateInput,
) (*mcp.CallToolResult, EvaluateOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, EvaluateOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	start := time.Now()
	result, err := vibe.Evaluate(ctx, input.Script)
	duration := time.Since(start)

	stepResult := report.StepResult{
		ID:         s.session.NextStepID("evaluate"),
		Action:     "evaluate",
		Args:       map[string]any{"script": truncateString(input.Script, 100)},
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		stepResult.Status = report.StatusNoGo
		stepResult.Severity = report.SeverityMedium
		stepResult.Error = &report.StepError{
			Type:    "EvaluateError",
			Message: err.Error(),
		}
		s.session.RecordStep(stepResult)
		return nil, EvaluateOutput{}, fmt.Errorf("evaluate failed: %w", err)
	}

	stepResult.Status = report.StatusGo
	stepResult.Severity = report.SeverityInfo
	s.session.RecordStep(stepResult)

	// Record for script export
	s.session.Recorder().RecordEval(input.Script)

	return nil, EvaluateOutput{Result: result}, nil
}

type AssertTextInput struct {
	Text     string `json:"text" jsonschema:"description=Text to search for,required"`
	Selector string `json:"selector" jsonschema:"description=Optional: limit search to element matching selector"`
}

type AssertTextOutput struct {
	Found   bool   `json:"found"`
	Message string `json:"message"`
}

func (s *Server) handleAssertText(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input AssertTextInput,
) (*mcp.CallToolResult, AssertTextOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, AssertTextOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	start := time.Now()

	// Build script to check for text
	var script string
	if input.Selector != "" {
		script = fmt.Sprintf(`
			(function() {
				const el = document.querySelector(%q);
				return el && el.textContent.includes(%q);
			})()
		`, input.Selector, input.Text)
	} else {
		script = fmt.Sprintf(`document.body.textContent.includes(%q)`, input.Text)
	}

	result, err := vibe.Evaluate(ctx, script)
	duration := time.Since(start)

	stepResult := report.StepResult{
		ID:         s.session.NextStepID("assert_text"),
		Action:     "assert_text",
		Args:       map[string]any{"text": input.Text, "selector": input.Selector},
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		stepResult.Status = report.StatusNoGo
		stepResult.Severity = report.SeverityCritical
		stepResult.Error = &report.StepError{
			Type:    "AssertTextError",
			Message: err.Error(),
		}
		s.session.RecordStep(stepResult)
		return nil, AssertTextOutput{}, fmt.Errorf("assert text failed: %w", err)
	}

	found, _ := result.(bool)
	if !found {
		stepResult.Status = report.StatusNoGo
		stepResult.Severity = report.SeverityCritical
		stepResult.Error = &report.StepError{
			Type:    "AssertTextFailed",
			Message: fmt.Sprintf("Text %q not found", input.Text),
		}
		stepResult.Context = s.session.CaptureContext(ctx)
		stepResult.Screenshot = s.session.CaptureScreenshot(ctx)
		s.session.RecordStep(stepResult)
		return nil, AssertTextOutput{Found: false, Message: fmt.Sprintf("Text %q not found", input.Text)}, nil
	}

	stepResult.Status = report.StatusGo
	stepResult.Severity = report.SeverityInfo
	s.session.RecordStep(stepResult)

	return nil, AssertTextOutput{Found: true, Message: fmt.Sprintf("Text %q found", input.Text)}, nil
}

type AssertElementInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the element,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type AssertElementOutput struct {
	Found   bool   `json:"found"`
	Message string `json:"message"`
}

func (s *Server) handleAssertElement(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input AssertElementInput,
) (*mcp.CallToolResult, AssertElementOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, AssertElementOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	start := time.Now()
	_, err = vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})
	duration := time.Since(start)

	stepResult := report.StepResult{
		ID:         s.session.NextStepID("assert_element"),
		Action:     "assert_element",
		Args:       map[string]any{"selector": input.Selector},
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		stepResult.Status = report.StatusNoGo
		stepResult.Severity = report.SeverityCritical
		stepResult.Error = &report.StepError{
			Type:        "AssertElementFailed",
			Message:     fmt.Sprintf("Element %q not found", input.Selector),
			Selector:    input.Selector,
			TimeoutMS:   int64(input.TimeoutMS),
			Suggestions: s.session.FindSimilarSelectors(ctx, input.Selector),
		}
		stepResult.Context = s.session.CaptureContext(ctx)
		stepResult.Screenshot = s.session.CaptureScreenshot(ctx)
		s.session.RecordStep(stepResult)
		return nil, AssertElementOutput{Found: false, Message: fmt.Sprintf("Element %q not found", input.Selector)}, nil
	}

	stepResult.Status = report.StatusGo
	stepResult.Severity = report.SeverityInfo
	s.session.RecordStep(stepResult)

	return nil, AssertElementOutput{Found: true, Message: fmt.Sprintf("Element %q found", input.Selector)}, nil
}

type GetTestReportInput struct {
	Format string `json:"format" jsonschema:"description=Report format: box (terminal) or diagnostic (full JSON) or json (multi-agent-spec),enum=box,enum=diagnostic,enum=json"`
}

type GetTestReportOutput struct {
	Report string `json:"report"`
}

func (s *Server) handleGetTestReport(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetTestReportInput,
) (*mcp.CallToolResult, GetTestReportOutput, error) {
	if input.Format == "" {
		input.Format = "box"
	}

	testResult := s.session.GetTestResult()

	switch input.Format {
	case "box":
		rendered, err := report.RenderBoxString(testResult)
		if err != nil {
			return nil, GetTestReportOutput{}, fmt.Errorf("render failed: %w", err)
		}
		return nil, GetTestReportOutput{Report: rendered}, nil

	case "diagnostic":
		diag := report.NewDiagnosticReport(testResult)
		diag.GenerateRecommendations()
		jsonBytes, err := diag.JSON()
		if err != nil {
			return nil, GetTestReportOutput{}, fmt.Errorf("json marshal failed: %w", err)
		}
		return nil, GetTestReportOutput{Report: string(jsonBytes)}, nil

	case "json":
		teamReport := report.ToTeamReport(testResult)
		jsonBytes, err := json.MarshalIndent(teamReport, "", "  ")
		if err != nil {
			return nil, GetTestReportOutput{}, fmt.Errorf("json marshal failed: %w", err)
		}
		return nil, GetTestReportOutput{Report: string(jsonBytes)}, nil

	default:
		return nil, GetTestReportOutput{}, fmt.Errorf("unknown format: %s", input.Format)
	}
}

type ResetSessionInput struct{}

type ResetSessionOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleResetSession(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ResetSessionInput,
) (*mcp.CallToolResult, ResetSessionOutput, error) {
	s.session.Reset()
	return nil, ResetSessionOutput{Message: "Session reset successfully"}, nil
}

type SetTargetInput struct {
	Target string `json:"target" jsonschema:"description=Test target description,required"`
}

type SetTargetOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleSetTarget(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SetTargetInput,
) (*mcp.CallToolResult, SetTargetOutput, error) {
	s.session.SetTarget(input.Target)
	return nil, SetTargetOutput{Message: fmt.Sprintf("Target set to: %s", input.Target)}, nil
}

// truncateString shortens a string to maxLen.
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
