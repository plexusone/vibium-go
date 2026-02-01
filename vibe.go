package vibium

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

// Vibe is the main browser control interface.
type Vibe struct {
	client          *BiDiClient
	clicker         *ClickerProcess
	browsingContext string
	closed          bool
}

// Browser provides browser launching capabilities.
var Browser = &browserLauncher{}

type browserLauncher struct{}

// Launch starts a new browser instance and returns a Vibe for controlling it.
func (b *browserLauncher) Launch(ctx context.Context, opts *LaunchOptions) (*Vibe, error) {
	if opts == nil {
		opts = &LaunchOptions{}
	}

	// Set up debug logging if enabled
	if logger := NewDebugLogger(); logger != nil {
		ctx = ContextWithLogger(ctx, logger)
		debugLog(ctx, "launching browser", "headless", opts.Headless, "port", opts.Port)
	}

	// Start clicker process
	clicker, err := StartClicker(ctx, *opts)
	if err != nil {
		return nil, err
	}
	debugLog(ctx, "clicker started", "url", clicker.WebSocketURL())

	// Connect BiDi client
	client := NewBiDiClient()
	if err := client.Connect(ctx, clicker.WebSocketURL()); err != nil {
		_ = clicker.Stop()
		return nil, err
	}
	debugLog(ctx, "BiDi client connected")

	return &Vibe{
		client:  client,
		clicker: clicker,
	}, nil
}

// Launch is a convenience function that launches a browser with default options.
func Launch(ctx context.Context) (*Vibe, error) {
	return Browser.Launch(ctx, nil)
}

// LaunchHeadless is a convenience function that launches a headless browser.
func LaunchHeadless(ctx context.Context) (*Vibe, error) {
	return Browser.Launch(ctx, &LaunchOptions{Headless: true})
}

// getContext returns the browsing context ID, fetching it if necessary.
func (v *Vibe) getContext(ctx context.Context) (string, error) {
	if v.browsingContext != "" {
		return v.browsingContext, nil
	}

	result, err := v.client.Send(ctx, "browsingContext.getTree", map[string]interface{}{})
	if err != nil {
		return "", fmt.Errorf("failed to get browsing context: %w", err)
	}

	var tree struct {
		Contexts []struct {
			Context string `json:"context"`
		} `json:"contexts"`
	}
	if err := json.Unmarshal(result, &tree); err != nil {
		return "", fmt.Errorf("failed to parse browsing context tree: %w", err)
	}

	if len(tree.Contexts) == 0 {
		return "", fmt.Errorf("no browsing context available")
	}

	v.browsingContext = tree.Contexts[0].Context
	return v.browsingContext, nil
}

// Go navigates to the specified URL.
func (v *Vibe) Go(ctx context.Context, url string) error {
	if v.closed {
		return ErrConnectionClosed
	}
	debugLog(ctx, "navigating", "url", url)

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"url":     url,
		"wait":    "complete",
	}

	_, err = v.client.Send(ctx, "browsingContext.navigate", params)
	if err == nil {
		debugLog(ctx, "navigation complete", "url", url)
	}
	return err
}

// Reload reloads the current page.
func (v *Vibe) Reload(ctx context.Context) error {
	if v.closed {
		return ErrConnectionClosed
	}
	debugLog(ctx, "reloading page")

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"wait":    "complete",
	}

	_, err = v.client.Send(ctx, "browsingContext.reload", params)
	return err
}

// Back navigates back in history.
func (v *Vibe) Back(ctx context.Context) error {
	if v.closed {
		return ErrConnectionClosed
	}
	debugLog(ctx, "navigating back")

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"delta":   -1,
	}

	_, err = v.client.Send(ctx, "browsingContext.traverseHistory", params)
	return err
}

// Forward navigates forward in history.
func (v *Vibe) Forward(ctx context.Context) error {
	if v.closed {
		return ErrConnectionClosed
	}
	debugLog(ctx, "navigating forward")

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"delta":   1,
	}

	_, err = v.client.Send(ctx, "browsingContext.traverseHistory", params)
	return err
}

// Screenshot captures a screenshot of the current page and returns PNG data.
func (v *Vibe) Screenshot(ctx context.Context) ([]byte, error) {
	if v.closed {
		return nil, ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return nil, err
	}

	result, err := v.client.Send(ctx, "browsingContext.captureScreenshot", map[string]interface{}{
		"context": browsingCtx,
	})
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data string `json:"data"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse screenshot response: %w", err)
	}

	// Decode base64 PNG data
	data, err := base64.StdEncoding.DecodeString(resp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode screenshot data: %w", err)
	}

	return data, nil
}

// Find finds an element by CSS selector.
func (v *Vibe) Find(ctx context.Context, selector string, opts *FindOptions) (*Element, error) {
	if v.closed {
		return nil, ErrConnectionClosed
	}
	debugLog(ctx, "finding element", "selector", selector)

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return nil, err
	}

	timeout := DefaultTimeout
	if opts != nil && opts.Timeout > 0 {
		timeout = opts.Timeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	params := map[string]interface{}{
		"context":  browsingCtx,
		"selector": selector,
		"timeout":  timeout.Milliseconds(),
	}

	result, err := v.client.Send(ctx, "vibium:find", params)
	if err != nil {
		return nil, err
	}

	var info ElementInfo
	if err := json.Unmarshal(result, &info); err != nil {
		return nil, fmt.Errorf("failed to parse element info: %w", err)
	}

	debugLog(ctx, "element found", "selector", selector, "tag", info.Tag)
	return NewElement(v.client, browsingCtx, selector, info), nil
}

// FindAll finds all elements matching the CSS selector.
func (v *Vibe) FindAll(ctx context.Context, selector string) ([]*Element, error) {
	if v.closed {
		return nil, ErrConnectionClosed
	}
	debugLog(ctx, "finding all elements", "selector", selector)

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return nil, err
	}

	// Use JavaScript to find all matching elements and get their info
	script := `(selector) => {
		const elements = document.querySelectorAll(selector);
		return Array.from(elements).map((el, index) => {
			const rect = el.getBoundingClientRect();
			return {
				index: index,
				tag: el.tagName.toLowerCase(),
				text: (el.textContent || '').trim().substring(0, 100),
				box: { x: rect.x, y: rect.y, width: rect.width, height: rect.height }
			};
		});
	}`

	params := map[string]interface{}{
		"functionDeclaration": script,
		"target":              map[string]interface{}{"context": browsingCtx},
		"arguments": []interface{}{
			map[string]interface{}{
				"type":  "string",
				"value": selector,
			},
		},
		"awaitPromise":    false,
		"resultOwnership": "root",
	}

	result, err := v.client.Send(ctx, "script.callFunction", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result struct {
			Value []struct {
				Index int         `json:"index"`
				Tag   string      `json:"tag"`
				Text  string      `json:"text"`
				Box   BoundingBox `json:"box"`
			} `json:"value"`
		} `json:"result"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse elements: %w", err)
	}

	elements := make([]*Element, len(resp.Result.Value))
	for i, item := range resp.Result.Value {
		// Create indexed selector for each element
		indexedSelector := fmt.Sprintf("%s:nth-of-type(%d)", selector, item.Index+1)
		info := ElementInfo{
			Tag:  item.Tag,
			Text: item.Text,
			Box:  item.Box,
		}
		elements[i] = NewElement(v.client, browsingCtx, indexedSelector, info)
	}

	debugLog(ctx, "elements found", "selector", selector, "count", len(elements))
	return elements, nil
}

// MustFind finds an element by CSS selector and panics if not found.
func (v *Vibe) MustFind(ctx context.Context, selector string) *Element {
	elem, err := v.Find(ctx, selector, nil)
	if err != nil {
		panic(err)
	}
	return elem
}

// Evaluate executes JavaScript in the page context and returns the result.
func (v *Vibe) Evaluate(ctx context.Context, script string) (interface{}, error) {
	if v.closed {
		return nil, ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return nil, err
	}

	// Wrap script in arrow function
	wrappedScript := fmt.Sprintf("() => { %s }", script)

	params := map[string]interface{}{
		"functionDeclaration": wrappedScript,
		"target":              map[string]interface{}{"context": browsingCtx},
		"arguments":           []interface{}{},
		"awaitPromise":        true,
		"resultOwnership":     "root",
	}

	result, err := v.client.Send(ctx, "script.callFunction", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result struct {
			Type  string      `json:"type"`
			Value interface{} `json:"value"`
		} `json:"result"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, err
	}

	return resp.Result.Value, nil
}

// Title returns the page title.
func (v *Vibe) Title(ctx context.Context) (string, error) {
	result, err := v.Evaluate(ctx, "return document.title")
	if err != nil {
		return "", err
	}
	if s, ok := result.(string); ok {
		return s, nil
	}
	return "", nil
}

// URL returns the current page URL.
func (v *Vibe) URL(ctx context.Context) (string, error) {
	result, err := v.Evaluate(ctx, "return window.location.href")
	if err != nil {
		return "", err
	}
	if s, ok := result.(string); ok {
		return s, nil
	}
	return "", nil
}

// WaitForNavigation waits for a navigation to complete.
func (v *Vibe) WaitForNavigation(ctx context.Context, timeout time.Duration) error {
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Simple implementation: wait for document ready state
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return &TimeoutError{
				Selector: "navigation",
				Timeout:  timeout.Milliseconds(),
				Reason:   "navigation did not complete",
			}
		case <-ticker.C:
			result, err := v.Evaluate(ctx, "return document.readyState")
			if err != nil {
				continue
			}
			if result == "complete" {
				return nil
			}
		}
	}
}

// Quit closes the browser and cleans up resources.
func (v *Vibe) Quit(ctx context.Context) error {
	if v.closed {
		return nil
	}
	v.closed = true

	// Close BiDi connection
	var clientErr error
	if v.client != nil {
		clientErr = v.client.Close()
	}

	// Stop clicker process
	if v.clicker != nil {
		if err := v.clicker.Stop(); err != nil {
			return err
		}
	}

	return clientErr
}

// IsClosed returns whether the browser has been closed.
func (v *Vibe) IsClosed() bool {
	return v.closed
}
