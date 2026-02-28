package mcp

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	vibium "github.com/plexusone/vibium-go"
	"github.com/plexusone/vibium-go/mcp/report"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// DragTo tool

type DragToInput struct {
	SourceSelector string `json:"source_selector" jsonschema:"description=CSS selector for the element to drag,required"`
	TargetSelector string `json:"target_selector" jsonschema:"description=CSS selector for the drop target,required"`
	TimeoutMS      int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type DragToOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleDragTo(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input DragToInput,
) (*mcp.CallToolResult, DragToOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, DragToOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	start := time.Now()
	source, err := vibe.Find(ctx, input.SourceSelector, &vibium.FindOptions{Timeout: timeout})
	if err != nil {
		return nil, DragToOutput{}, fmt.Errorf("source element not found: %s", input.SourceSelector)
	}

	target, err := vibe.Find(ctx, input.TargetSelector, &vibium.FindOptions{Timeout: timeout})
	if err != nil {
		return nil, DragToOutput{}, fmt.Errorf("target element not found: %s", input.TargetSelector)
	}

	result := report.StepResult{
		ID:     s.session.NextStepID("drag_to"),
		Action: "drag_to",
		Args:   map[string]any{"source": input.SourceSelector, "target": input.TargetSelector},
	}

	err = source.DragTo(ctx, target, &vibium.ActionOptions{Timeout: timeout})
	result.DurationMS = time.Since(start).Milliseconds()

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityCritical
		result.Error = &report.StepError{
			Type:    "DragToError",
			Message: err.Error(),
		}
		s.session.RecordStep(result)
		return nil, DragToOutput{}, fmt.Errorf("drag to failed: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	return nil, DragToOutput{Message: fmt.Sprintf("Dragged %s to %s", input.SourceSelector, input.TargetSelector)}, nil
}

// Tap tool

type TapInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the element,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type TapOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleTap(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input TapInput,
) (*mcp.CallToolResult, TapOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, TapOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})
	if err != nil {
		return nil, TapOutput{}, fmt.Errorf("element not found: %s", input.Selector)
	}

	err = elem.Tap(ctx, &vibium.ActionOptions{Timeout: timeout})
	if err != nil {
		return nil, TapOutput{}, fmt.Errorf("tap failed: %w", err)
	}

	return nil, TapOutput{Message: fmt.Sprintf("Tapped %s", input.Selector)}, nil
}

// DispatchEvent tool

type DispatchEventInput struct {
	Selector  string         `json:"selector" jsonschema:"description=CSS selector for the element,required"`
	EventType string         `json:"event_type" jsonschema:"description=Event type (e.g. click focus blur),required"`
	EventInit map[string]any `json:"event_init" jsonschema:"description=Event initialization options"`
	TimeoutMS int            `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type DispatchEventOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleDispatchEvent(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input DispatchEventInput,
) (*mcp.CallToolResult, DispatchEventOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, DispatchEventOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})
	if err != nil {
		return nil, DispatchEventOutput{}, fmt.Errorf("element not found: %s", input.Selector)
	}

	err = elem.DispatchEvent(ctx, input.EventType, input.EventInit)
	if err != nil {
		return nil, DispatchEventOutput{}, fmt.Errorf("dispatch event failed: %w", err)
	}

	return nil, DispatchEventOutput{Message: fmt.Sprintf("Dispatched %s on %s", input.EventType, input.Selector)}, nil
}

// SetFiles tool

type SetFilesInput struct {
	Selector  string   `json:"selector" jsonschema:"description=CSS selector for the file input,required"`
	Files     []string `json:"files" jsonschema:"description=File paths to set,required"`
	TimeoutMS int      `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type SetFilesOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleSetFiles(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SetFilesInput,
) (*mcp.CallToolResult, SetFilesOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, SetFilesOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})
	if err != nil {
		return nil, SetFilesOutput{}, fmt.Errorf("element not found: %s", input.Selector)
	}

	err = elem.SetFiles(ctx, input.Files, &vibium.ActionOptions{Timeout: timeout})
	if err != nil {
		return nil, SetFilesOutput{}, fmt.Errorf("set files failed: %w", err)
	}

	return nil, SetFilesOutput{Message: fmt.Sprintf("Set %d files on %s", len(input.Files), input.Selector)}, nil
}

// ElementScreenshot tool

type ElementScreenshotInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the element,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type ElementScreenshotOutput struct {
	Data string `json:"data"`
}

func (s *Server) handleElementScreenshot(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ElementScreenshotInput,
) (*mcp.CallToolResult, ElementScreenshotOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, ElementScreenshotOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})
	if err != nil {
		return nil, ElementScreenshotOutput{}, fmt.Errorf("element not found: %s", input.Selector)
	}

	data, err := elem.Screenshot(ctx)
	if err != nil {
		return nil, ElementScreenshotOutput{}, fmt.Errorf("element screenshot failed: %w", err)
	}

	return nil, ElementScreenshotOutput{Data: base64.StdEncoding.EncodeToString(data)}, nil
}

// ElementEval tool

type ElementEvalInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the element,required"`
	Function  string `json:"function" jsonschema:"description=JavaScript function (receives element as first arg),required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type ElementEvalOutput struct {
	Result any `json:"result"`
}

func (s *Server) handleElementEval(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ElementEvalInput,
) (*mcp.CallToolResult, ElementEvalOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, ElementEvalOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})
	if err != nil {
		return nil, ElementEvalOutput{}, fmt.Errorf("element not found: %s", input.Selector)
	}

	result, err := elem.Eval(ctx, input.Function)
	if err != nil {
		return nil, ElementEvalOutput{}, fmt.Errorf("element eval failed: %w", err)
	}

	return nil, ElementEvalOutput{Result: result}, nil
}
