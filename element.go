package vibium

import (
	"context"
	"encoding/base64"
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
	params := map[string]interface{}{
		"context":  e.context,
		"selector": e.selector,
	}

	result, err := e.client.Send(ctx, "vibium:el.text", params)
	if err != nil {
		return "", err
	}

	var resp struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", err
	}

	return strings.TrimSpace(resp.Text), nil
}

// GetAttribute returns the value of the specified attribute.
func (e *Element) GetAttribute(ctx context.Context, name string) (string, error) {
	params := map[string]interface{}{
		"context":  e.context,
		"selector": e.selector,
		"name":     name,
	}

	result, err := e.client.Send(ctx, "vibium:el.attr", params)
	if err != nil {
		return "", err
	}

	var resp struct {
		Value *string `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", err
	}

	if resp.Value == nil {
		return "", nil
	}
	return *resp.Value, nil
}

// BoundingBox returns the element's bounding box.
func (e *Element) BoundingBox(ctx context.Context) (BoundingBox, error) {
	params := map[string]interface{}{
		"context":  e.context,
		"selector": e.selector,
	}

	result, err := e.client.Send(ctx, "vibium:el.bounds", params)
	if err != nil {
		return BoundingBox{}, err
	}

	var box BoundingBox
	if err := json.Unmarshal(result, &box); err != nil {
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

// Fill clears the input and fills it with the specified value.
// It waits for the element to be visible, stable, enabled, and editable before filling.
func (e *Element) Fill(ctx context.Context, value string, opts *ActionOptions) error {
	timeout := DefaultTimeout
	if opts != nil && opts.Timeout > 0 {
		timeout = opts.Timeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	params := map[string]interface{}{
		"context":  e.context,
		"selector": e.selector,
		"value":    value,
		"timeout":  timeout.Milliseconds(),
	}

	_, err := e.client.Send(ctx, "vibium:fill", params)
	return err
}

// Press presses a key on the element.
// It waits for the element to be visible, stable, and able to receive events.
func (e *Element) Press(ctx context.Context, key string, opts *ActionOptions) error {
	timeout := DefaultTimeout
	if opts != nil && opts.Timeout > 0 {
		timeout = opts.Timeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	params := map[string]interface{}{
		"context":  e.context,
		"selector": e.selector,
		"key":      key,
		"timeout":  timeout.Milliseconds(),
	}

	_, err := e.client.Send(ctx, "vibium:press", params)
	return err
}

// Clear clears the text content of an input field.
func (e *Element) Clear(ctx context.Context, opts *ActionOptions) error {
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

	_, err := e.client.Send(ctx, "vibium:clear", params)
	return err
}

// Check checks a checkbox element.
// It waits for the element to be visible, stable, and enabled.
func (e *Element) Check(ctx context.Context, opts *ActionOptions) error {
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

	_, err := e.client.Send(ctx, "vibium:check", params)
	return err
}

// Uncheck unchecks a checkbox element.
// It waits for the element to be visible, stable, and enabled.
func (e *Element) Uncheck(ctx context.Context, opts *ActionOptions) error {
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

	_, err := e.client.Send(ctx, "vibium:uncheck", params)
	return err
}

// SelectOption selects an option in a <select> element by value, label, or index.
func (e *Element) SelectOption(ctx context.Context, values SelectOptionValues, opts *ActionOptions) error {
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

	if len(values.Values) > 0 {
		params["values"] = values.Values
	}
	if len(values.Labels) > 0 {
		params["labels"] = values.Labels
	}
	if len(values.Indexes) > 0 {
		params["indexes"] = values.Indexes
	}

	_, err := e.client.Send(ctx, "vibium:selectOption", params)
	return err
}

// Focus focuses the element.
func (e *Element) Focus(ctx context.Context, opts *ActionOptions) error {
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

	_, err := e.client.Send(ctx, "vibium:focus", params)
	return err
}

// Hover moves the mouse over the element.
func (e *Element) Hover(ctx context.Context, opts *ActionOptions) error {
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

	_, err := e.client.Send(ctx, "vibium:hover", params)
	return err
}

// ScrollIntoView scrolls the element into the visible area of the viewport.
func (e *Element) ScrollIntoView(ctx context.Context, opts *ActionOptions) error {
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

	_, err := e.client.Send(ctx, "vibium:scrollIntoView", params)
	return err
}

// DblClick double-clicks on the element.
func (e *Element) DblClick(ctx context.Context, opts *ActionOptions) error {
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

	_, err := e.client.Send(ctx, "vibium:dblclick", params)
	return err
}

// Value returns the value of an input element.
func (e *Element) Value(ctx context.Context) (string, error) {
	params := map[string]interface{}{
		"context":  e.context,
		"selector": e.selector,
	}

	result, err := e.client.Send(ctx, "vibium:el.value", params)
	if err != nil {
		return "", err
	}

	var resp struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", err
	}

	return resp.Value, nil
}

// InnerHTML returns the inner HTML of the element.
func (e *Element) InnerHTML(ctx context.Context) (string, error) {
	params := map[string]interface{}{
		"context":  e.context,
		"selector": e.selector,
	}

	result, err := e.client.Send(ctx, "vibium:el.html", params)
	if err != nil {
		return "", err
	}

	var resp struct {
		HTML string `json:"html"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", err
	}

	return resp.HTML, nil
}

// InnerText returns the rendered text content of the element.
func (e *Element) InnerText(ctx context.Context) (string, error) {
	params := map[string]interface{}{
		"context":  e.context,
		"selector": e.selector,
	}

	result, err := e.client.Send(ctx, "vibium:el.innerText", params)
	if err != nil {
		return "", err
	}

	var resp struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", err
	}

	return resp.Text, nil
}

// IsVisible returns whether the element is visible.
func (e *Element) IsVisible(ctx context.Context) (bool, error) {
	params := map[string]interface{}{
		"context":  e.context,
		"selector": e.selector,
	}

	result, err := e.client.Send(ctx, "vibium:el.isVisible", params)
	if err != nil {
		return false, err
	}

	var resp struct {
		Visible bool `json:"visible"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return false, err
	}

	return resp.Visible, nil
}

// IsHidden returns whether the element is hidden.
func (e *Element) IsHidden(ctx context.Context) (bool, error) {
	params := map[string]interface{}{
		"context":  e.context,
		"selector": e.selector,
	}

	result, err := e.client.Send(ctx, "vibium:el.isHidden", params)
	if err != nil {
		return false, err
	}

	var resp struct {
		Hidden bool `json:"hidden"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return false, err
	}

	return resp.Hidden, nil
}

// IsEnabled returns whether the element is enabled.
func (e *Element) IsEnabled(ctx context.Context) (bool, error) {
	params := map[string]interface{}{
		"context":  e.context,
		"selector": e.selector,
	}

	result, err := e.client.Send(ctx, "vibium:el.isEnabled", params)
	if err != nil {
		return false, err
	}

	var resp struct {
		Enabled bool `json:"enabled"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return false, err
	}

	return resp.Enabled, nil
}

// IsChecked returns whether a checkbox or radio element is checked.
func (e *Element) IsChecked(ctx context.Context) (bool, error) {
	params := map[string]interface{}{
		"context":  e.context,
		"selector": e.selector,
	}

	result, err := e.client.Send(ctx, "vibium:el.isChecked", params)
	if err != nil {
		return false, err
	}

	var resp struct {
		Checked bool `json:"checked"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return false, err
	}

	return resp.Checked, nil
}

// IsEditable returns whether the element is editable.
func (e *Element) IsEditable(ctx context.Context) (bool, error) {
	params := map[string]interface{}{
		"context":  e.context,
		"selector": e.selector,
	}

	result, err := e.client.Send(ctx, "vibium:el.isEditable", params)
	if err != nil {
		return false, err
	}

	var resp struct {
		Editable bool `json:"editable"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return false, err
	}

	return resp.Editable, nil
}

// Role returns the ARIA role of the element.
func (e *Element) Role(ctx context.Context) (string, error) {
	params := map[string]interface{}{
		"context":  e.context,
		"selector": e.selector,
	}

	result, err := e.client.Send(ctx, "vibium:el.role", params)
	if err != nil {
		return "", err
	}

	var resp struct {
		Role string `json:"role"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", err
	}

	return resp.Role, nil
}

// Label returns the accessible label of the element.
func (e *Element) Label(ctx context.Context) (string, error) {
	params := map[string]interface{}{
		"context":  e.context,
		"selector": e.selector,
	}

	result, err := e.client.Send(ctx, "vibium:el.label", params)
	if err != nil {
		return "", err
	}

	var resp struct {
		Label string `json:"label"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", err
	}

	return resp.Label, nil
}

// WaitUntil waits for the element to reach the specified state.
// State can be: "attached", "detached", "visible", "hidden".
func (e *Element) WaitUntil(ctx context.Context, state string, timeout time.Duration) error {
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	params := map[string]interface{}{
		"context":  e.context,
		"selector": e.selector,
		"state":    state,
		"timeout":  timeout.Milliseconds(),
	}

	_, err := e.client.Send(ctx, "vibium:el.waitFor", params)
	return err
}

// DragTo drags this element to the target element.
func (e *Element) DragTo(ctx context.Context, target *Element, opts *ActionOptions) error {
	timeout := DefaultTimeout
	if opts != nil && opts.Timeout > 0 {
		timeout = opts.Timeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	params := map[string]interface{}{
		"context":        e.context,
		"selector":       e.selector,
		"targetSelector": target.selector,
		"timeout":        timeout.Milliseconds(),
	}

	_, err := e.client.Send(ctx, "vibium:dragTo", params)
	return err
}

// Tap performs a touch tap on the element.
func (e *Element) Tap(ctx context.Context, opts *ActionOptions) error {
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

	_, err := e.client.Send(ctx, "vibium:tap", params)
	return err
}

// DispatchEvent dispatches a DOM event on the element.
func (e *Element) DispatchEvent(ctx context.Context, eventType string, eventInit map[string]interface{}) error {
	params := map[string]interface{}{
		"context":   e.context,
		"selector":  e.selector,
		"eventType": eventType,
	}

	if eventInit != nil {
		params["eventInit"] = eventInit
	}

	_, err := e.client.Send(ctx, "vibium:dispatchEvent", params)
	return err
}

// SetFiles sets the files for a file input element.
func (e *Element) SetFiles(ctx context.Context, paths []string, opts *ActionOptions) error {
	timeout := DefaultTimeout
	if opts != nil && opts.Timeout > 0 {
		timeout = opts.Timeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	params := map[string]interface{}{
		"context":  e.context,
		"selector": e.selector,
		"files":    paths,
		"timeout":  timeout.Milliseconds(),
	}

	_, err := e.client.Send(ctx, "vibium:el.setFiles", params)
	return err
}

// Screenshot captures a screenshot of just this element.
func (e *Element) Screenshot(ctx context.Context) ([]byte, error) {
	params := map[string]interface{}{
		"context":  e.context,
		"selector": e.selector,
	}

	result, err := e.client.Send(ctx, "vibium:el.screenshot", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data string `json:"data"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, err
	}

	// Decode base64 PNG data
	return decodeBase64(resp.Data)
}

// Eval evaluates a JavaScript function with this element as the argument.
// The function should accept the element as its first parameter.
func (e *Element) Eval(ctx context.Context, fn string, args ...interface{}) (interface{}, error) {
	params := map[string]interface{}{
		"context":  e.context,
		"selector": e.selector,
		"fn":       fn,
	}

	if len(args) > 0 {
		params["args"] = args
	}

	result, err := e.client.Send(ctx, "vibium:el.eval", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Value interface{} `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, err
	}

	return resp.Value, nil
}

// decodeBase64 decodes a base64 string to bytes.
func decodeBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
