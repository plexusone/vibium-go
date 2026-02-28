package vibium

import (
	"context"
)

// Touch provides touch input control.
type Touch struct {
	client  *BiDiClient
	context string
}

// NewTouch creates a new Touch controller.
func NewTouch(client *BiDiClient, browsingContext string) *Touch {
	return &Touch{
		client:  client,
		context: browsingContext,
	}
}

// Tap performs a tap at the specified coordinates.
func (t *Touch) Tap(ctx context.Context, x, y float64) error {
	params := map[string]interface{}{
		"context": t.context,
		"x":       x,
		"y":       y,
	}

	_, err := t.client.Send(ctx, "vibium:touch.tap", params)
	return err
}

// Swipe performs a swipe gesture from one point to another.
func (t *Touch) Swipe(ctx context.Context, startX, startY, endX, endY float64) error {
	params := map[string]interface{}{
		"context": t.context,
		"startX":  startX,
		"startY":  startY,
		"endX":    endX,
		"endY":    endY,
	}

	_, err := t.client.Send(ctx, "vibium:touch.swipe", params)
	return err
}

// Pinch performs a pinch gesture.
// Scale < 1 zooms out, scale > 1 zooms in.
func (t *Touch) Pinch(ctx context.Context, x, y float64, scale float64) error {
	params := map[string]interface{}{
		"context": t.context,
		"x":       x,
		"y":       y,
		"scale":   scale,
	}

	_, err := t.client.Send(ctx, "vibium:touch.pinch", params)
	return err
}
