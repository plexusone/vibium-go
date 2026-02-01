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

	// Start clicker process
	clicker, err := StartClicker(ctx, *opts)
	if err != nil {
		return nil, err
	}

	// Connect BiDi client
	client := NewBiDiClient()
	if err := client.Connect(ctx, clicker.WebSocketURL()); err != nil {
		_ = clicker.Stop()
		return nil, err
	}

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

	return NewElement(v.client, browsingCtx, selector, info), nil
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
