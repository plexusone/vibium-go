package vibium

import (
	"context"
)

// Dialog represents a browser dialog (alert, confirm, prompt, beforeunload).
type Dialog struct {
	client  *BiDiClient
	context string
	id      string
	Type    string `json:"type"` // "alert", "confirm", "prompt", "beforeunload"
	Message string `json:"message"`
	Default string `json:"defaultValue,omitempty"` // For prompt dialogs
}

// Accept accepts the dialog.
// For prompt dialogs, optionally provide a text value.
func (d *Dialog) Accept(ctx context.Context, promptText string) error {
	params := map[string]interface{}{
		"context": d.context,
		"id":      d.id,
		"accept":  true,
	}

	if promptText != "" {
		params["userText"] = promptText
	}

	_, err := d.client.Send(ctx, "vibium:dialog.handle", params)
	return err
}

// Dismiss dismisses the dialog (clicks cancel/no).
func (d *Dialog) Dismiss(ctx context.Context) error {
	params := map[string]interface{}{
		"context": d.context,
		"id":      d.id,
		"accept":  false,
	}

	_, err := d.client.Send(ctx, "vibium:dialog.handle", params)
	return err
}
