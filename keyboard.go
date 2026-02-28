package vibium

import (
	"context"
)

// Keyboard provides keyboard input control.
type Keyboard struct {
	client  *BiDiClient
	context string
}

// NewKeyboard creates a new Keyboard controller.
func NewKeyboard(client *BiDiClient, browsingContext string) *Keyboard {
	return &Keyboard{
		client:  client,
		context: browsingContext,
	}
}

// Press presses a key on the keyboard.
// Key names follow the Playwright key naming convention (e.g., "Enter", "Tab", "ArrowUp").
func (k *Keyboard) Press(ctx context.Context, key string) error {
	params := map[string]interface{}{
		"context": k.context,
		"key":     key,
	}

	_, err := k.client.Send(ctx, "vibium:keyboard.press", params)
	return err
}

// Down holds down a key.
func (k *Keyboard) Down(ctx context.Context, key string) error {
	params := map[string]interface{}{
		"context": k.context,
		"key":     key,
	}

	_, err := k.client.Send(ctx, "vibium:keyboard.down", params)
	return err
}

// Up releases a held key.
func (k *Keyboard) Up(ctx context.Context, key string) error {
	params := map[string]interface{}{
		"context": k.context,
		"key":     key,
	}

	_, err := k.client.Send(ctx, "vibium:keyboard.up", params)
	return err
}

// Type types text character by character.
// This sends individual keypress events for each character.
func (k *Keyboard) Type(ctx context.Context, text string) error {
	params := map[string]interface{}{
		"context": k.context,
		"text":    text,
	}

	_, err := k.client.Send(ctx, "vibium:keyboard.type", params)
	return err
}

// InsertText inserts text directly without keypress events.
// This is faster than Type but doesn't trigger keyboard events.
func (k *Keyboard) InsertText(ctx context.Context, text string) error {
	params := map[string]interface{}{
		"context": k.context,
		"text":    text,
	}

	_, err := k.client.Send(ctx, "vibium:keyboard.insertText", params)
	return err
}
