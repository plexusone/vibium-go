package vibium

import (
	"context"
)

// Mouse provides mouse input control.
type Mouse struct {
	client  *BiDiClient
	context string
}

// NewMouse creates a new Mouse controller.
func NewMouse(client *BiDiClient, browsingContext string) *Mouse {
	return &Mouse{
		client:  client,
		context: browsingContext,
	}
}

// MouseButton represents a mouse button.
type MouseButton string

const (
	MouseButtonLeft   MouseButton = "left"
	MouseButtonRight  MouseButton = "right"
	MouseButtonMiddle MouseButton = "middle"
)

// ClickOptions configures mouse click behavior.
type ClickOptions struct {
	Button     MouseButton
	ClickCount int
	Delay      int // milliseconds between mousedown and mouseup
}

// Click clicks at the specified coordinates.
func (m *Mouse) Click(ctx context.Context, x, y float64, opts *ClickOptions) error {
	params := map[string]interface{}{
		"context": m.context,
		"x":       x,
		"y":       y,
	}

	if opts != nil {
		if opts.Button != "" {
			params["button"] = string(opts.Button)
		}
		if opts.ClickCount > 0 {
			params["clickCount"] = opts.ClickCount
		}
		if opts.Delay > 0 {
			params["delay"] = opts.Delay
		}
	}

	_, err := m.client.Send(ctx, "vibium:mouse.click", params)
	return err
}

// DblClick double-clicks at the specified coordinates.
func (m *Mouse) DblClick(ctx context.Context, x, y float64, opts *ClickOptions) error {
	params := map[string]interface{}{
		"context":    m.context,
		"x":          x,
		"y":          y,
		"clickCount": 2,
	}

	if opts != nil {
		if opts.Button != "" {
			params["button"] = string(opts.Button)
		}
		if opts.Delay > 0 {
			params["delay"] = opts.Delay
		}
	}

	_, err := m.client.Send(ctx, "vibium:mouse.click", params)
	return err
}

// Move moves the mouse to the specified coordinates.
func (m *Mouse) Move(ctx context.Context, x, y float64) error {
	params := map[string]interface{}{
		"context": m.context,
		"x":       x,
		"y":       y,
	}

	_, err := m.client.Send(ctx, "vibium:mouse.move", params)
	return err
}

// Down presses the mouse button.
func (m *Mouse) Down(ctx context.Context, button MouseButton) error {
	params := map[string]interface{}{
		"context": m.context,
	}

	if button != "" {
		params["button"] = string(button)
	}

	_, err := m.client.Send(ctx, "vibium:mouse.down", params)
	return err
}

// Up releases the mouse button.
func (m *Mouse) Up(ctx context.Context, button MouseButton) error {
	params := map[string]interface{}{
		"context": m.context,
	}

	if button != "" {
		params["button"] = string(button)
	}

	_, err := m.client.Send(ctx, "vibium:mouse.up", params)
	return err
}

// Wheel scrolls the mouse wheel.
func (m *Mouse) Wheel(ctx context.Context, deltaX, deltaY float64) error {
	params := map[string]interface{}{
		"context": m.context,
		"deltaX":  deltaX,
		"deltaY":  deltaY,
	}

	_, err := m.client.Send(ctx, "vibium:mouse.wheel", params)
	return err
}
