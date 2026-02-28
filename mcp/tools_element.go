package mcp

import (
	"context"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	vibium "github.com/plexusone/vibium-go"
)

// elementOp performs a common element operation pattern:
// 1. Gets vibe from session
// 2. Finds element by selector with timeout
// 3. Calls the provided operation function on the element
func (s *Server) elementOp(ctx context.Context, selector string, timeoutMS int, op func(*vibium.Element) (any, error)) (any, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, fmt.Errorf("browser not available: %w", err)
	}

	if timeoutMS == 0 {
		timeoutMS = 5000
	}
	timeout := time.Duration(timeoutMS) * time.Millisecond

	elem, err := vibe.Find(ctx, selector, &vibium.FindOptions{Timeout: timeout})
	if err != nil {
		return nil, fmt.Errorf("element not found: %s", selector)
	}

	return op(elem)
}

// GetValue tool

type GetValueInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the input element,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type GetValueOutput struct {
	Value string `json:"value"`
}

func (s *Server) handleGetValue(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetValueInput,
) (*mcp.CallToolResult, GetValueOutput, error) {
	result, err := s.elementOp(ctx, input.Selector, input.TimeoutMS, func(elem *vibium.Element) (any, error) {
		return elem.Value(ctx)
	})
	if err != nil {
		return nil, GetValueOutput{}, err
	}
	return nil, GetValueOutput{Value: result.(string)}, nil
}

// GetInnerHTML tool

type GetInnerHTMLInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the element,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type GetInnerHTMLOutput struct {
	HTML string `json:"html"`
}

func (s *Server) handleGetInnerHTML(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetInnerHTMLInput,
) (*mcp.CallToolResult, GetInnerHTMLOutput, error) {
	result, err := s.elementOp(ctx, input.Selector, input.TimeoutMS, func(elem *vibium.Element) (any, error) {
		return elem.InnerHTML(ctx)
	})
	if err != nil {
		return nil, GetInnerHTMLOutput{}, err
	}
	return nil, GetInnerHTMLOutput{HTML: result.(string)}, nil
}

// GetInnerText tool

type GetInnerTextInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the element,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type GetInnerTextOutput struct {
	Text string `json:"text"`
}

func (s *Server) handleGetInnerText(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetInnerTextInput,
) (*mcp.CallToolResult, GetInnerTextOutput, error) {
	result, err := s.elementOp(ctx, input.Selector, input.TimeoutMS, func(elem *vibium.Element) (any, error) {
		return elem.InnerText(ctx)
	})
	if err != nil {
		return nil, GetInnerTextOutput{}, err
	}
	return nil, GetInnerTextOutput{Text: result.(string)}, nil
}

// GetAttribute tool

type GetAttributeInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the element,required"`
	Name      string `json:"name" jsonschema:"description=Attribute name,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type GetAttributeOutput struct {
	Value string `json:"value"`
}

func (s *Server) handleGetAttribute(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetAttributeInput,
) (*mcp.CallToolResult, GetAttributeOutput, error) {
	result, err := s.elementOp(ctx, input.Selector, input.TimeoutMS, func(elem *vibium.Element) (any, error) {
		return elem.GetAttribute(ctx, input.Name)
	})
	if err != nil {
		return nil, GetAttributeOutput{}, err
	}
	return nil, GetAttributeOutput{Value: result.(string)}, nil
}

// IsVisible tool

type IsVisibleInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the element,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type IsVisibleOutput struct {
	Visible bool `json:"visible"`
}

func (s *Server) handleIsVisible(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input IsVisibleInput,
) (*mcp.CallToolResult, IsVisibleOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, IsVisibleOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})
	if err != nil {
		return nil, IsVisibleOutput{Visible: false}, nil // Element not found = not visible
	}

	visible, err := elem.IsVisible(ctx)
	if err != nil {
		return nil, IsVisibleOutput{}, fmt.Errorf("is visible check failed: %w", err)
	}

	return nil, IsVisibleOutput{Visible: visible}, nil
}

// IsHidden tool

type IsHiddenInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the element,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type IsHiddenOutput struct {
	Hidden bool `json:"hidden"`
}

func (s *Server) handleIsHidden(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input IsHiddenInput,
) (*mcp.CallToolResult, IsHiddenOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, IsHiddenOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})
	if err != nil {
		return nil, IsHiddenOutput{Hidden: true}, nil // Element not found = hidden
	}

	hidden, err := elem.IsHidden(ctx)
	if err != nil {
		return nil, IsHiddenOutput{}, fmt.Errorf("is hidden check failed: %w", err)
	}

	return nil, IsHiddenOutput{Hidden: hidden}, nil
}

// IsEnabled tool

type IsEnabledInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the element,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type IsEnabledOutput struct {
	Enabled bool `json:"enabled"`
}

func (s *Server) handleIsEnabled(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input IsEnabledInput,
) (*mcp.CallToolResult, IsEnabledOutput, error) {
	result, err := s.elementOp(ctx, input.Selector, input.TimeoutMS, func(elem *vibium.Element) (any, error) {
		return elem.IsEnabled(ctx)
	})
	if err != nil {
		return nil, IsEnabledOutput{}, err
	}
	return nil, IsEnabledOutput{Enabled: result.(bool)}, nil
}

// IsChecked tool

type IsCheckedInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the checkbox/radio,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type IsCheckedOutput struct {
	Checked bool `json:"checked"`
}

func (s *Server) handleIsChecked(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input IsCheckedInput,
) (*mcp.CallToolResult, IsCheckedOutput, error) {
	result, err := s.elementOp(ctx, input.Selector, input.TimeoutMS, func(elem *vibium.Element) (any, error) {
		return elem.IsChecked(ctx)
	})
	if err != nil {
		return nil, IsCheckedOutput{}, err
	}
	return nil, IsCheckedOutput{Checked: result.(bool)}, nil
}

// IsEditable tool

type IsEditableInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the element,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type IsEditableOutput struct {
	Editable bool `json:"editable"`
}

func (s *Server) handleIsEditable(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input IsEditableInput,
) (*mcp.CallToolResult, IsEditableOutput, error) {
	result, err := s.elementOp(ctx, input.Selector, input.TimeoutMS, func(elem *vibium.Element) (any, error) {
		return elem.IsEditable(ctx)
	})
	if err != nil {
		return nil, IsEditableOutput{}, err
	}
	return nil, IsEditableOutput{Editable: result.(bool)}, nil
}

// GetRole tool

type GetRoleInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the element,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type GetRoleOutput struct {
	Role string `json:"role"`
}

func (s *Server) handleGetRole(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetRoleInput,
) (*mcp.CallToolResult, GetRoleOutput, error) {
	result, err := s.elementOp(ctx, input.Selector, input.TimeoutMS, func(elem *vibium.Element) (any, error) {
		return elem.Role(ctx)
	})
	if err != nil {
		return nil, GetRoleOutput{}, err
	}
	return nil, GetRoleOutput{Role: result.(string)}, nil
}

// GetLabel tool

type GetLabelInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the element,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type GetLabelOutput struct {
	Label string `json:"label"`
}

func (s *Server) handleGetLabel(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetLabelInput,
) (*mcp.CallToolResult, GetLabelOutput, error) {
	result, err := s.elementOp(ctx, input.Selector, input.TimeoutMS, func(elem *vibium.Element) (any, error) {
		return elem.Label(ctx)
	})
	if err != nil {
		return nil, GetLabelOutput{}, err
	}
	return nil, GetLabelOutput{Label: result.(string)}, nil
}

// WaitUntil tool

type WaitUntilInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the element,required"`
	State     string `json:"state" jsonschema:"description=State to wait for: attached detached visible hidden,required,enum=attached,enum=detached,enum=visible,enum=hidden"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 30000)"`
}

type WaitUntilOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleWaitUntil(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input WaitUntilInput,
) (*mcp.CallToolResult, WaitUntilOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, WaitUntilOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 30000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	// First find the element
	elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})
	if err != nil {
		// For "detached" state, not finding element is success
		if input.State == "detached" {
			return nil, WaitUntilOutput{Message: fmt.Sprintf("Element %s is detached", input.Selector)}, nil
		}
		return nil, WaitUntilOutput{}, fmt.Errorf("element not found: %s", input.Selector)
	}

	err = elem.WaitUntil(ctx, input.State, timeout)
	if err != nil {
		return nil, WaitUntilOutput{}, fmt.Errorf("wait until %s failed: %w", input.State, err)
	}

	return nil, WaitUntilOutput{Message: fmt.Sprintf("Element %s is %s", input.Selector, input.State)}, nil
}

// GetBoundingBox tool

type GetBoundingBoxInput struct {
	Selector  string `json:"selector" jsonschema:"description=CSS selector for the element,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 5000)"`
}

type GetBoundingBoxOutput struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

func (s *Server) handleGetBoundingBox(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetBoundingBoxInput,
) (*mcp.CallToolResult, GetBoundingBoxOutput, error) {
	result, err := s.elementOp(ctx, input.Selector, input.TimeoutMS, func(elem *vibium.Element) (any, error) {
		return elem.BoundingBox(ctx)
	})
	if err != nil {
		return nil, GetBoundingBoxOutput{}, err
	}
	box := result.(vibium.BoundingBox)
	return nil, GetBoundingBoxOutput{
		X:      box.X,
		Y:      box.Y,
		Width:  box.Width,
		Height: box.Height,
	}, nil
}
