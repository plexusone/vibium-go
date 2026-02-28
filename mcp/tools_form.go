//nolint:dupl // form handlers have similar patterns with different actions and reporting
package mcp

import (
	"context"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	vibium "github.com/plexusone/vibium-go"
	"github.com/plexusone/vibium-go/mcp/report"
)

// Fill tool

type FillInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the input element,required"`
	Value     string `json:"value" jsonschema:"description=Value to fill,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type FillOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleFill(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input FillInput,
) (*mcp.CallToolResult, FillOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, FillOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	start := time.Now()
	elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})

	result := report.StepResult{
		ID:     s.session.NextStepID("fill"),
		Action: "fill",
		Args:   map[string]any{"selector": input.Selector, "value": truncateString(input.Value, 50)},
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
		result.Screenshot = s.session.CaptureScreenshot(ctx)
		s.session.RecordStep(result)
		return nil, FillOutput{}, fmt.Errorf("element not found: %s", input.Selector)
	}

	err = elem.Fill(ctx, input.Value, &vibium.ActionOptions{Timeout: timeout})
	result.DurationMS = time.Since(start).Milliseconds()

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityCritical
		result.Error = &report.StepError{
			Type:     "FillError",
			Message:  err.Error(),
			Selector: input.Selector,
		}
		result.Screenshot = s.session.CaptureScreenshot(ctx)
		s.session.RecordStep(result)
		return nil, FillOutput{}, fmt.Errorf("fill failed: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	// Record for script export
	s.session.Recorder().RecordFill(input.Selector, input.Value)

	return nil, FillOutput{Message: fmt.Sprintf("Filled %s", input.Selector)}, nil
}

// Press tool

type PressInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the element,required"`
	Key       string `json:"key" jsonschema:"description=Key to press (e.g. Enter Tab ArrowDown),required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type PressOutput struct {
	Message string `json:"message"`
}

func (s *Server) handlePress(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input PressInput,
) (*mcp.CallToolResult, PressOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, PressOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	start := time.Now()
	elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})

	result := report.StepResult{
		ID:     s.session.NextStepID("press"),
		Action: "press",
		Args:   map[string]any{"selector": input.Selector, "key": input.Key},
	}

	if err != nil {
		result.DurationMS = time.Since(start).Milliseconds()
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityCritical
		result.Error = &report.StepError{
			Type:     "ElementNotFoundError",
			Message:  err.Error(),
			Selector: input.Selector,
		}
		s.session.RecordStep(result)
		return nil, PressOutput{}, fmt.Errorf("element not found: %s", input.Selector)
	}

	err = elem.Press(ctx, input.Key, &vibium.ActionOptions{Timeout: timeout})
	result.DurationMS = time.Since(start).Milliseconds()

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityCritical
		result.Error = &report.StepError{
			Type:     "PressError",
			Message:  err.Error(),
			Selector: input.Selector,
		}
		s.session.RecordStep(result)
		return nil, PressOutput{}, fmt.Errorf("press failed: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	// Record for script export
	s.session.Recorder().RecordPress(input.Selector, input.Key)

	return nil, PressOutput{Message: fmt.Sprintf("Pressed %s on %s", input.Key, input.Selector)}, nil
}

// Clear tool

type ClearInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the input element,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type ClearOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleClear(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ClearInput,
) (*mcp.CallToolResult, ClearOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, ClearOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	start := time.Now()
	elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})

	result := report.StepResult{
		ID:     s.session.NextStepID("clear"),
		Action: "clear",
		Args:   map[string]any{"selector": input.Selector},
	}

	if err != nil {
		result.DurationMS = time.Since(start).Milliseconds()
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityCritical
		result.Error = &report.StepError{
			Type:     "ElementNotFoundError",
			Message:  err.Error(),
			Selector: input.Selector,
		}
		s.session.RecordStep(result)
		return nil, ClearOutput{}, fmt.Errorf("element not found: %s", input.Selector)
	}

	err = elem.Clear(ctx, &vibium.ActionOptions{Timeout: timeout})
	result.DurationMS = time.Since(start).Milliseconds()

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityCritical
		result.Error = &report.StepError{
			Type:     "ClearError",
			Message:  err.Error(),
			Selector: input.Selector,
		}
		s.session.RecordStep(result)
		return nil, ClearOutput{}, fmt.Errorf("clear failed: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	// Record for script export
	s.session.Recorder().RecordClear(input.Selector)

	return nil, ClearOutput{Message: fmt.Sprintf("Cleared %s", input.Selector)}, nil
}

// Check tool

type CheckInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the checkbox,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type CheckOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleCheck(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input CheckInput,
) (*mcp.CallToolResult, CheckOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, CheckOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	start := time.Now()
	elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})

	result := report.StepResult{
		ID:     s.session.NextStepID("check"),
		Action: "check",
		Args:   map[string]any{"selector": input.Selector},
	}

	if err != nil {
		result.DurationMS = time.Since(start).Milliseconds()
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityCritical
		result.Error = &report.StepError{
			Type:     "ElementNotFoundError",
			Message:  err.Error(),
			Selector: input.Selector,
		}
		s.session.RecordStep(result)
		return nil, CheckOutput{}, fmt.Errorf("element not found: %s", input.Selector)
	}

	err = elem.Check(ctx, &vibium.ActionOptions{Timeout: timeout})
	result.DurationMS = time.Since(start).Milliseconds()

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityCritical
		result.Error = &report.StepError{
			Type:     "CheckError",
			Message:  err.Error(),
			Selector: input.Selector,
		}
		s.session.RecordStep(result)
		return nil, CheckOutput{}, fmt.Errorf("check failed: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	// Record for script export
	s.session.Recorder().RecordCheck(input.Selector)

	return nil, CheckOutput{Message: fmt.Sprintf("Checked %s", input.Selector)}, nil
}

// Uncheck tool

type UncheckInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the checkbox,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type UncheckOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleUncheck(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input UncheckInput,
) (*mcp.CallToolResult, UncheckOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, UncheckOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	start := time.Now()
	elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})

	result := report.StepResult{
		ID:     s.session.NextStepID("uncheck"),
		Action: "uncheck",
		Args:   map[string]any{"selector": input.Selector},
	}

	if err != nil {
		result.DurationMS = time.Since(start).Milliseconds()
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityCritical
		result.Error = &report.StepError{
			Type:     "ElementNotFoundError",
			Message:  err.Error(),
			Selector: input.Selector,
		}
		s.session.RecordStep(result)
		return nil, UncheckOutput{}, fmt.Errorf("element not found: %s", input.Selector)
	}

	err = elem.Uncheck(ctx, &vibium.ActionOptions{Timeout: timeout})
	result.DurationMS = time.Since(start).Milliseconds()

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityCritical
		result.Error = &report.StepError{
			Type:     "UncheckError",
			Message:  err.Error(),
			Selector: input.Selector,
		}
		s.session.RecordStep(result)
		return nil, UncheckOutput{}, fmt.Errorf("uncheck failed: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	// Record for script export
	s.session.Recorder().RecordUncheck(input.Selector)

	return nil, UncheckOutput{Message: fmt.Sprintf("Unchecked %s", input.Selector)}, nil
}

// SelectOption tool

type SelectOptionInput struct {
	Selector  string   `json:"selector" jsonschema:"description=CSS selector for the select element,required"`
	Values    []string `json:"values" jsonschema:"description=Option values to select"`
	Labels    []string `json:"labels" jsonschema:"description=Option labels to select"`
	Indexes   []int    `json:"indexes" jsonschema:"description=Option indexes to select (0-based)"`
	TimeoutMS int      `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type SelectOptionOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleSelectOption(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SelectOptionInput,
) (*mcp.CallToolResult, SelectOptionOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, SelectOptionOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	start := time.Now()
	elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})

	result := report.StepResult{
		ID:     s.session.NextStepID("select_option"),
		Action: "select_option",
		Args:   map[string]any{"selector": input.Selector, "values": input.Values, "labels": input.Labels},
	}

	if err != nil {
		result.DurationMS = time.Since(start).Milliseconds()
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityCritical
		result.Error = &report.StepError{
			Type:     "ElementNotFoundError",
			Message:  err.Error(),
			Selector: input.Selector,
		}
		s.session.RecordStep(result)
		return nil, SelectOptionOutput{}, fmt.Errorf("element not found: %s", input.Selector)
	}

	selectValues := vibium.SelectOptionValues{
		Values:  input.Values,
		Labels:  input.Labels,
		Indexes: input.Indexes,
	}
	err = elem.SelectOption(ctx, selectValues, &vibium.ActionOptions{Timeout: timeout})
	result.DurationMS = time.Since(start).Milliseconds()

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityCritical
		result.Error = &report.StepError{
			Type:     "SelectOptionError",
			Message:  err.Error(),
			Selector: input.Selector,
		}
		s.session.RecordStep(result)
		return nil, SelectOptionOutput{}, fmt.Errorf("select option failed: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	// Record for script export
	value := ""
	if len(input.Values) > 0 {
		value = input.Values[0]
	} else if len(input.Labels) > 0 {
		value = input.Labels[0]
	}
	s.session.Recorder().RecordSelect(input.Selector, value)

	return nil, SelectOptionOutput{Message: fmt.Sprintf("Selected option in %s", input.Selector)}, nil
}

// Focus tool

type FocusInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the element,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type FocusOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleFocus(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input FocusInput,
) (*mcp.CallToolResult, FocusOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, FocusOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	start := time.Now()
	elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})

	result := report.StepResult{
		ID:     s.session.NextStepID("focus"),
		Action: "focus",
		Args:   map[string]any{"selector": input.Selector},
	}

	if err != nil {
		result.DurationMS = time.Since(start).Milliseconds()
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:     "ElementNotFoundError",
			Message:  err.Error(),
			Selector: input.Selector,
		}
		s.session.RecordStep(result)
		return nil, FocusOutput{}, fmt.Errorf("element not found: %s", input.Selector)
	}

	err = elem.Focus(ctx, &vibium.ActionOptions{Timeout: timeout})
	result.DurationMS = time.Since(start).Milliseconds()

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:     "FocusError",
			Message:  err.Error(),
			Selector: input.Selector,
		}
		s.session.RecordStep(result)
		return nil, FocusOutput{}, fmt.Errorf("focus failed: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	// Record for script export
	s.session.Recorder().RecordFocus(input.Selector)

	return nil, FocusOutput{Message: fmt.Sprintf("Focused %s", input.Selector)}, nil
}

// Hover tool

type HoverInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the element,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type HoverOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleHover(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input HoverInput,
) (*mcp.CallToolResult, HoverOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, HoverOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	start := time.Now()
	elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})

	result := report.StepResult{
		ID:     s.session.NextStepID("hover"),
		Action: "hover",
		Args:   map[string]any{"selector": input.Selector},
	}

	if err != nil {
		result.DurationMS = time.Since(start).Milliseconds()
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:     "ElementNotFoundError",
			Message:  err.Error(),
			Selector: input.Selector,
		}
		s.session.RecordStep(result)
		return nil, HoverOutput{}, fmt.Errorf("element not found: %s", input.Selector)
	}

	err = elem.Hover(ctx, &vibium.ActionOptions{Timeout: timeout})
	result.DurationMS = time.Since(start).Milliseconds()

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:     "HoverError",
			Message:  err.Error(),
			Selector: input.Selector,
		}
		s.session.RecordStep(result)
		return nil, HoverOutput{}, fmt.Errorf("hover failed: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	// Record for script export
	s.session.Recorder().RecordHover(input.Selector)

	return nil, HoverOutput{Message: fmt.Sprintf("Hovered over %s", input.Selector)}, nil
}

// ScrollIntoView tool

type ScrollIntoViewInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the element,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type ScrollIntoViewOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleScrollIntoView(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ScrollIntoViewInput,
) (*mcp.CallToolResult, ScrollIntoViewOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, ScrollIntoViewOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	start := time.Now()
	elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})

	result := report.StepResult{
		ID:     s.session.NextStepID("scroll_into_view"),
		Action: "scroll_into_view",
		Args:   map[string]any{"selector": input.Selector},
	}

	if err != nil {
		result.DurationMS = time.Since(start).Milliseconds()
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:     "ElementNotFoundError",
			Message:  err.Error(),
			Selector: input.Selector,
		}
		s.session.RecordStep(result)
		return nil, ScrollIntoViewOutput{}, fmt.Errorf("element not found: %s", input.Selector)
	}

	err = elem.ScrollIntoView(ctx, &vibium.ActionOptions{Timeout: timeout})
	result.DurationMS = time.Since(start).Milliseconds()

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:     "ScrollIntoViewError",
			Message:  err.Error(),
			Selector: input.Selector,
		}
		s.session.RecordStep(result)
		return nil, ScrollIntoViewOutput{}, fmt.Errorf("scroll into view failed: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	// Record for script export
	s.session.Recorder().RecordScrollIntoView(input.Selector)

	return nil, ScrollIntoViewOutput{Message: fmt.Sprintf("Scrolled %s into view", input.Selector)}, nil
}

// DblClick tool

type DblClickInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the element,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type DblClickOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleDblClick(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input DblClickInput,
) (*mcp.CallToolResult, DblClickOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, DblClickOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	start := time.Now()
	elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})

	result := report.StepResult{
		ID:     s.session.NextStepID("dblclick"),
		Action: "dblclick",
		Args:   map[string]any{"selector": input.Selector},
	}

	if err != nil {
		result.DurationMS = time.Since(start).Milliseconds()
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityCritical
		result.Error = &report.StepError{
			Type:     "ElementNotFoundError",
			Message:  err.Error(),
			Selector: input.Selector,
		}
		s.session.RecordStep(result)
		return nil, DblClickOutput{}, fmt.Errorf("element not found: %s", input.Selector)
	}

	err = elem.DblClick(ctx, &vibium.ActionOptions{Timeout: timeout})
	result.DurationMS = time.Since(start).Milliseconds()

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityCritical
		result.Error = &report.StepError{
			Type:     "DblClickError",
			Message:  err.Error(),
			Selector: input.Selector,
		}
		s.session.RecordStep(result)
		return nil, DblClickOutput{}, fmt.Errorf("double click failed: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	// Record for script export
	s.session.Recorder().RecordDblClick(input.Selector)

	return nil, DblClickOutput{Message: fmt.Sprintf("Double-clicked %s", input.Selector)}, nil
}
