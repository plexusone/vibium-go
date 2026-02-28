package activity

import (
	"context"
	"fmt"
	"time"

	vibium "github.com/agentplexus/vibium-go"
)

// FindActivity finds an element by selector.
type FindActivity struct{}

func (a *FindActivity) Name() string { return "element.find" }

func (a *FindActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	selector := GetString(params, "selector")
	if selector == "" {
		return nil, fmt.Errorf("selector parameter is required")
	}

	timeout := time.Duration(GetIntDefault(params, "timeout", 30000)) * time.Millisecond
	opts := &vibium.FindOptions{Timeout: timeout}

	// Add semantic selector options
	if role := GetString(params, "role"); role != "" {
		opts.Role = role
	}
	if text := GetString(params, "text"); text != "" {
		opts.Text = text
	}
	if label := GetString(params, "label"); label != "" {
		opts.Label = label
	}
	if testID := GetString(params, "testId"); testID != "" {
		opts.TestID = testID
	}

	el, err := env.Vibe.Find(ctx, selector, opts)
	if err != nil {
		return nil, fmt.Errorf("element not found: %w", err)
	}

	// Return element info
	return map[string]any{
		"selector": el.Selector(),
		"tag":      el.Info().Tag,
		"text":     el.Info().Text,
		"box":      el.Info().Box,
	}, nil
}

// FindAllActivity finds all elements matching a selector.
type FindAllActivity struct{}

func (a *FindAllActivity) Name() string { return "element.findAll" }

func (a *FindAllActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	selector := GetString(params, "selector")
	if selector == "" {
		return nil, fmt.Errorf("selector parameter is required")
	}

	elements, err := env.Vibe.FindAll(ctx, selector)
	if err != nil {
		return nil, fmt.Errorf("find failed: %w", err)
	}

	result := make([]map[string]any, len(elements))
	for i, el := range elements {
		result[i] = map[string]any{
			"selector": el.Selector(),
			"tag":      el.Info().Tag,
			"text":     el.Info().Text,
			"box":      el.Info().Box,
		}
	}

	return result, nil
}

// GetTextActivity gets the text content of an element.
type GetTextActivity struct{}

func (a *GetTextActivity) Name() string { return "element.getText" }

func (a *GetTextActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	selector := GetString(params, "selector")
	if selector == "" {
		return nil, fmt.Errorf("selector parameter is required")
	}

	timeout := time.Duration(GetIntDefault(params, "timeout", 30000)) * time.Millisecond
	opts := &vibium.FindOptions{Timeout: timeout}

	el, err := env.Vibe.Find(ctx, selector, opts)
	if err != nil {
		return nil, fmt.Errorf("element not found: %w", err)
	}

	text, err := el.Text(ctx)
	if err != nil {
		return nil, fmt.Errorf("get text failed: %w", err)
	}

	return text, nil
}

// GetValueActivity gets the value of an input element.
type GetValueActivity struct{}

func (a *GetValueActivity) Name() string { return "element.getValue" }

func (a *GetValueActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	selector := GetString(params, "selector")
	if selector == "" {
		return nil, fmt.Errorf("selector parameter is required")
	}

	timeout := time.Duration(GetIntDefault(params, "timeout", 30000)) * time.Millisecond
	opts := &vibium.FindOptions{Timeout: timeout}

	el, err := env.Vibe.Find(ctx, selector, opts)
	if err != nil {
		return nil, fmt.Errorf("element not found: %w", err)
	}

	value, err := el.Value(ctx)
	if err != nil {
		return nil, fmt.Errorf("get value failed: %w", err)
	}

	return value, nil
}

// GetAttributeActivity gets an attribute value from an element.
type GetAttributeActivity struct{}

func (a *GetAttributeActivity) Name() string { return "element.getAttribute" }

func (a *GetAttributeActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	selector := GetString(params, "selector")
	if selector == "" {
		return nil, fmt.Errorf("selector parameter is required")
	}

	name := GetString(params, "name")
	if name == "" {
		return nil, fmt.Errorf("name parameter is required")
	}

	timeout := time.Duration(GetIntDefault(params, "timeout", 30000)) * time.Millisecond
	opts := &vibium.FindOptions{Timeout: timeout}

	el, err := env.Vibe.Find(ctx, selector, opts)
	if err != nil {
		return nil, fmt.Errorf("element not found: %w", err)
	}

	value, err := el.GetAttribute(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("get attribute failed: %w", err)
	}

	return value, nil
}

// WaitForActivity waits for an element to reach a state.
type WaitForActivity struct{}

func (a *WaitForActivity) Name() string { return "element.waitFor" }

func (a *WaitForActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	selector := GetString(params, "selector")
	if selector == "" {
		return nil, fmt.Errorf("selector parameter is required")
	}

	state := GetStringDefault(params, "state", "visible")
	timeout := time.Duration(GetIntDefault(params, "timeout", 30000)) * time.Millisecond
	opts := &vibium.FindOptions{Timeout: timeout}

	el, err := env.Vibe.Find(ctx, selector, opts)
	if err != nil {
		return nil, fmt.Errorf("element not found: %w", err)
	}

	if err := el.WaitUntil(ctx, state, timeout); err != nil {
		return nil, fmt.Errorf("wait failed: %w", err)
	}

	return nil, nil
}

// IsVisibleActivity checks if an element is visible.
type IsVisibleActivity struct{}

func (a *IsVisibleActivity) Name() string { return "element.isVisible" }

func (a *IsVisibleActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	selector := GetString(params, "selector")
	if selector == "" {
		return nil, fmt.Errorf("selector parameter is required")
	}

	timeout := time.Duration(GetIntDefault(params, "timeout", 30000)) * time.Millisecond
	opts := &vibium.FindOptions{Timeout: timeout}

	el, err := env.Vibe.Find(ctx, selector, opts)
	if err != nil {
		// Element not found means not visible
		return false, nil
	}

	visible, err := el.IsVisible(ctx)
	if err != nil {
		return false, err
	}

	return visible, nil
}
