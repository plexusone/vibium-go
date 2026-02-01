package vibium

import (
	"context"
	"encoding/json"
	"strings"
	"time"
)

// Element represents a DOM element that can be interacted with.
type Element struct {
	client   *BiDiClient
	context  string // browsing context ID
	selector string
	info     ElementInfo
}

// NewElement creates a new Element instance.
func NewElement(client *BiDiClient, browsingContext, selector string, info ElementInfo) *Element {
	return &Element{
		client:   client,
		context:  browsingContext,
		selector: selector,
		info:     info,
	}
}

// Info returns the element's metadata.
func (e *Element) Info() ElementInfo {
	return e.info
}

// Selector returns the CSS selector used to find this element.
func (e *Element) Selector() string {
	return e.selector
}

// Click clicks on the element. It waits for the element to be visible, stable,
// able to receive events, and enabled before clicking.
func (e *Element) Click(ctx context.Context, opts *ActionOptions) error {
	timeout := DefaultTimeout
	if opts != nil && opts.Timeout > 0 {
		timeout = opts.Timeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	params := map[string]interface{}{
		"context":  e.context,
		"selector": e.selector,
		"timeout":  timeout.Milliseconds(),
	}

	_, err := e.client.Send(ctx, "vibium:click", params)
	return err
}

// Type types text into the element. It waits for the element to be visible,
// stable, able to receive events, enabled, and editable before typing.
func (e *Element) Type(ctx context.Context, text string, opts *ActionOptions) error {
	timeout := DefaultTimeout
	if opts != nil && opts.Timeout > 0 {
		timeout = opts.Timeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	params := map[string]interface{}{
		"context":  e.context,
		"selector": e.selector,
		"text":     text,
		"timeout":  timeout.Milliseconds(),
	}

	_, err := e.client.Send(ctx, "vibium:type", params)
	return err
}

// Text returns the text content of the element.
func (e *Element) Text(ctx context.Context) (string, error) {
	script := `(selector) => {
		const el = document.querySelector(selector);
		return el ? (el.textContent || '').trim() : null;
	}`

	params := map[string]interface{}{
		"functionDeclaration": script,
		"target":              map[string]interface{}{"context": e.context},
		"arguments": []interface{}{
			map[string]interface{}{
				"type":  "string",
				"value": e.selector,
			},
		},
		"awaitPromise":    false,
		"resultOwnership": "root",
	}

	result, err := e.client.Send(ctx, "script.callFunction", params)
	if err != nil {
		return "", err
	}

	var resp struct {
		Result struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"result"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", err
	}

	if resp.Result.Type == "null" {
		return "", &ElementNotFoundError{Selector: e.selector}
	}

	return strings.TrimSpace(resp.Result.Value), nil
}

// GetAttribute returns the value of the specified attribute.
func (e *Element) GetAttribute(ctx context.Context, name string) (string, error) {
	script := `(selector, attrName) => {
		const el = document.querySelector(selector);
		return el ? el.getAttribute(attrName) : null;
	}`

	params := map[string]interface{}{
		"functionDeclaration": script,
		"target":              map[string]interface{}{"context": e.context},
		"arguments": []interface{}{
			map[string]interface{}{
				"type":  "string",
				"value": e.selector,
			},
			map[string]interface{}{
				"type":  "string",
				"value": name,
			},
		},
		"awaitPromise":    false,
		"resultOwnership": "root",
	}

	result, err := e.client.Send(ctx, "script.callFunction", params)
	if err != nil {
		return "", err
	}

	var resp struct {
		Result struct {
			Type  string  `json:"type"`
			Value *string `json:"value"`
		} `json:"result"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", err
	}

	if resp.Result.Type == "null" || resp.Result.Value == nil {
		return "", nil
	}
	return *resp.Result.Value, nil
}

// BoundingBox returns the element's bounding box.
func (e *Element) BoundingBox(ctx context.Context) (BoundingBox, error) {
	script := `(selector) => {
		const el = document.querySelector(selector);
		if (!el) return null;
		const rect = el.getBoundingClientRect();
		return JSON.stringify({ x: rect.x, y: rect.y, width: rect.width, height: rect.height });
	}`

	params := map[string]interface{}{
		"functionDeclaration": script,
		"target":              map[string]interface{}{"context": e.context},
		"arguments": []interface{}{
			map[string]interface{}{
				"type":  "string",
				"value": e.selector,
			},
		},
		"awaitPromise":    false,
		"resultOwnership": "root",
	}

	result, err := e.client.Send(ctx, "script.callFunction", params)
	if err != nil {
		return BoundingBox{}, err
	}

	var resp struct {
		Result struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"result"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return BoundingBox{}, err
	}

	if resp.Result.Type == "null" {
		return BoundingBox{}, &ElementNotFoundError{Selector: e.selector}
	}

	var box BoundingBox
	if err := json.Unmarshal([]byte(resp.Result.Value), &box); err != nil {
		return BoundingBox{}, err
	}

	return box, nil
}

// WaitFor waits for the element to appear in the DOM.
func (e *Element) WaitFor(ctx context.Context, timeout time.Duration) error {
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return &TimeoutError{
				Selector: e.selector,
				Timeout:  timeout.Milliseconds(),
				Reason:   "element did not appear",
			}
		case <-ticker.C:
			script := `(selector) => document.querySelector(selector) !== null`
			params := map[string]interface{}{
				"functionDeclaration": script,
				"target":              map[string]interface{}{"context": e.context},
				"arguments": []interface{}{
					map[string]interface{}{
						"type":  "string",
						"value": e.selector,
					},
				},
				"awaitPromise":    false,
				"resultOwnership": "root",
			}

			result, err := e.client.Send(ctx, "script.callFunction", params)
			if err != nil {
				continue
			}

			var resp struct {
				Result struct {
					Value bool `json:"value"`
				} `json:"result"`
			}
			if err := json.Unmarshal(result, &resp); err != nil {
				continue
			}

			if resp.Result.Value {
				return nil
			}
		}
	}
}

// Center returns the center point of the element.
func (e *Element) Center() (x, y float64) {
	return e.info.Box.X + e.info.Box.Width/2, e.info.Box.Y + e.info.Box.Height/2
}
