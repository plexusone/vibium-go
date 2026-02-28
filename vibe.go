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

	// Input controllers (lazy-initialized)
	keyboard *Keyboard
	mouse    *Mouse
	touch    *Touch
	clock    *Clock
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

	// Add semantic selector options if present
	if opts != nil {
		if opts.Role != "" {
			params["role"] = opts.Role
		}
		if opts.Text != "" {
			params["text"] = opts.Text
		}
		if opts.Label != "" {
			params["label"] = opts.Label
		}
		if opts.Placeholder != "" {
			params["placeholder"] = opts.Placeholder
		}
		if opts.TestID != "" {
			params["testid"] = opts.TestID
		}
		if opts.Alt != "" {
			params["alt"] = opts.Alt
		}
		if opts.Title != "" {
			params["title"] = opts.Title
		}
		if opts.XPath != "" {
			params["xpath"] = opts.XPath
		}
		if opts.Near != "" {
			params["near"] = opts.Near
		}
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

	// Use JavaScript to find all matching elements and return JSON string
	// (BiDi serializes arrays in a complex format, so we JSON.stringify ourselves)
	script := `(selector) => {
		const elements = document.querySelectorAll(selector);
		const result = Array.from(elements).map((el, index) => {
			const rect = el.getBoundingClientRect();
			return {
				index: index,
				tag: el.tagName.toLowerCase(),
				text: (el.textContent || '').trim().substring(0, 100),
				box: { x: rect.x, y: rect.y, width: rect.width, height: rect.height }
			};
		});
		return JSON.stringify(result);
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

	// Parse the outer BiDi response to get the JSON string
	var resp struct {
		Result struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"result"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Parse the JSON string containing element data
	var items []struct {
		Index int         `json:"index"`
		Tag   string      `json:"tag"`
		Text  string      `json:"text"`
		Box   BoundingBox `json:"box"`
	}
	if err := json.Unmarshal([]byte(resp.Result.Value), &items); err != nil {
		return nil, fmt.Errorf("failed to parse elements: %w", err)
	}

	elements := make([]*Element, len(items))
	for i, item := range items {
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

// Keyboard returns the keyboard controller for this page.
func (v *Vibe) Keyboard(ctx context.Context) (*Keyboard, error) {
	if v.keyboard != nil {
		return v.keyboard, nil
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return nil, err
	}

	v.keyboard = NewKeyboard(v.client, browsingCtx)
	return v.keyboard, nil
}

// Mouse returns the mouse controller for this page.
func (v *Vibe) Mouse(ctx context.Context) (*Mouse, error) {
	if v.mouse != nil {
		return v.mouse, nil
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return nil, err
	}

	v.mouse = NewMouse(v.client, browsingCtx)
	return v.mouse, nil
}

// Touch returns the touch controller for this page.
func (v *Vibe) Touch(ctx context.Context) (*Touch, error) {
	if v.touch != nil {
		return v.touch, nil
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return nil, err
	}

	v.touch = NewTouch(v.client, browsingCtx)
	return v.touch, nil
}

// Clock returns the clock controller for this page.
func (v *Vibe) Clock(ctx context.Context) (*Clock, error) {
	if v.clock != nil {
		return v.clock, nil
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return nil, err
	}

	v.clock = NewClock(v.client, browsingCtx)
	return v.clock, nil
}

// Content returns the full HTML content of the page.
func (v *Vibe) Content(ctx context.Context) (string, error) {
	if v.closed {
		return "", ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return "", err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	result, err := v.client.Send(ctx, "vibium:page.content", params)
	if err != nil {
		return "", err
	}

	var resp struct {
		Content string `json:"content"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", err
	}

	return resp.Content, nil
}

// SetContent sets the HTML content of the page.
func (v *Vibe) SetContent(ctx context.Context, html string) error {
	if v.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"html":    html,
	}

	_, err = v.client.Send(ctx, "vibium:page.setContent", params)
	return err
}

// GetViewport returns the current viewport dimensions.
func (v *Vibe) GetViewport(ctx context.Context) (Viewport, error) {
	if v.closed {
		return Viewport{}, ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return Viewport{}, err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	result, err := v.client.Send(ctx, "vibium:page.viewport", params)
	if err != nil {
		return Viewport{}, err
	}

	var vp Viewport
	if err := json.Unmarshal(result, &vp); err != nil {
		return Viewport{}, err
	}

	return vp, nil
}

// SetViewport sets the viewport dimensions.
func (v *Vibe) SetViewport(ctx context.Context, viewport Viewport) error {
	if v.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"width":   viewport.Width,
		"height":  viewport.Height,
	}

	_, err = v.client.Send(ctx, "vibium:page.setViewport", params)
	return err
}

// GetWindow returns the browser window state.
func (v *Vibe) GetWindow(ctx context.Context) (WindowState, error) {
	if v.closed {
		return WindowState{}, ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return WindowState{}, err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	result, err := v.client.Send(ctx, "vibium:page.window", params)
	if err != nil {
		return WindowState{}, err
	}

	var ws WindowState
	if err := json.Unmarshal(result, &ws); err != nil {
		return WindowState{}, err
	}

	return ws, nil
}

// SetWindow sets the browser window state.
func (v *Vibe) SetWindow(ctx context.Context, opts SetWindowOptions) error {
	if v.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	if opts.X != nil {
		params["x"] = *opts.X
	}
	if opts.Y != nil {
		params["y"] = *opts.Y
	}
	if opts.Width != nil {
		params["width"] = *opts.Width
	}
	if opts.Height != nil {
		params["height"] = *opts.Height
	}
	if opts.State != "" {
		params["state"] = opts.State
	}

	_, err = v.client.Send(ctx, "vibium:page.setWindow", params)
	return err
}

// PDF generates a PDF of the page and returns the bytes.
func (v *Vibe) PDF(ctx context.Context, opts *PDFOptions) ([]byte, error) {
	if v.closed {
		return nil, ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return nil, err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	if opts != nil {
		if opts.Scale != 0 {
			params["scale"] = opts.Scale
		}
		if opts.DisplayHeader {
			params["displayHeader"] = opts.DisplayHeader
		}
		if opts.DisplayFooter {
			params["displayFooter"] = opts.DisplayFooter
		}
		if opts.PrintBackground {
			params["printBackground"] = opts.PrintBackground
		}
		if opts.Landscape {
			params["landscape"] = opts.Landscape
		}
		if opts.PageRanges != "" {
			params["pageRanges"] = opts.PageRanges
		}
		if opts.Format != "" {
			params["format"] = opts.Format
		}
		if opts.Width != "" {
			params["width"] = opts.Width
		}
		if opts.Height != "" {
			params["height"] = opts.Height
		}
		if opts.Margin != nil {
			params["margin"] = map[string]interface{}{
				"top":    opts.Margin.Top,
				"right":  opts.Margin.Right,
				"bottom": opts.Margin.Bottom,
				"left":   opts.Margin.Left,
			}
		}
	}

	result, err := v.client.Send(ctx, "vibium:page.pdf", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data string `json:"data"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(resp.Data)
}

// BringToFront activates the page (brings the browser tab to front).
func (v *Vibe) BringToFront(ctx context.Context) error {
	if v.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	_, err = v.client.Send(ctx, "browsingContext.activate", params)
	return err
}

// Close closes the current page but not the browser.
func (v *Vibe) Close(ctx context.Context) error {
	if v.closed {
		return nil
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	_, err = v.client.Send(ctx, "browsingContext.close", params)
	return err
}

// Frames returns all frames on the page.
func (v *Vibe) Frames(ctx context.Context) ([]FrameInfo, error) {
	if v.closed {
		return nil, ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return nil, err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	result, err := v.client.Send(ctx, "vibium:page.frames", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Frames []FrameInfo `json:"frames"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, err
	}

	return resp.Frames, nil
}

// Frame finds a frame by name or URL pattern.
func (v *Vibe) Frame(ctx context.Context, nameOrURL string) (*Vibe, error) {
	if v.closed {
		return nil, ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return nil, err
	}

	params := map[string]interface{}{
		"context":   browsingCtx,
		"nameOrURL": nameOrURL,
	}

	result, err := v.client.Send(ctx, "vibium:page.frame", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Context string `json:"context"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, err
	}

	return &Vibe{
		client:          v.client,
		clicker:         v.clicker,
		browsingContext: resp.Context,
	}, nil
}

// A11yTree returns the accessibility tree for the page.
func (v *Vibe) A11yTree(ctx context.Context) (interface{}, error) {
	if v.closed {
		return nil, ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return nil, err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	result, err := v.client.Send(ctx, "vibium:page.a11yTree", params)
	if err != nil {
		return nil, err
	}

	var resp interface{}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// EmulateMedia sets the media emulation options.
func (v *Vibe) EmulateMedia(ctx context.Context, opts EmulateMediaOptions) error {
	if v.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	if opts.Media != "" {
		params["media"] = opts.Media
	}
	if opts.ColorScheme != "" {
		params["colorScheme"] = opts.ColorScheme
	}
	if opts.ReducedMotion != "" {
		params["reducedMotion"] = opts.ReducedMotion
	}
	if opts.ForcedColors != "" {
		params["forcedColors"] = opts.ForcedColors
	}

	_, err = v.client.Send(ctx, "vibium:page.emulateMedia", params)
	return err
}

// SetGeolocation overrides the browser's geolocation.
func (v *Vibe) SetGeolocation(ctx context.Context, coords Geolocation) error {
	if v.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context":   browsingCtx,
		"latitude":  coords.Latitude,
		"longitude": coords.Longitude,
	}

	if coords.Accuracy != 0 {
		params["accuracy"] = coords.Accuracy
	}

	_, err = v.client.Send(ctx, "vibium:page.setGeolocation", params)
	return err
}

// AddScript adds a script that will be evaluated in the page context.
func (v *Vibe) AddScript(ctx context.Context, source string) error {
	if v.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"source":  source,
	}

	_, err = v.client.Send(ctx, "vibium:page.addScript", params)
	return err
}

// AddStyle adds a stylesheet to the page.
func (v *Vibe) AddStyle(ctx context.Context, source string) error {
	if v.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"source":  source,
	}

	_, err = v.client.Send(ctx, "vibium:page.addStyle", params)
	return err
}

// Expose exposes a function that can be called from JavaScript in the page.
// Note: The handler function must be registered separately.
func (v *Vibe) Expose(ctx context.Context, name string) error {
	if v.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"name":    name,
	}

	_, err = v.client.Send(ctx, "vibium:page.expose", params)
	return err
}

// WaitForURL waits for the page URL to match the specified pattern.
func (v *Vibe) WaitForURL(ctx context.Context, pattern string, timeout time.Duration) error {
	if v.closed {
		return ErrConnectionClosed
	}

	if timeout == 0 {
		timeout = DefaultTimeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"pattern": pattern,
		"timeout": timeout.Milliseconds(),
	}

	_, err = v.client.Send(ctx, "vibium:page.waitForURL", params)
	return err
}

// WaitForLoad waits for the page to reach the specified load state.
// State can be: "load", "domcontentloaded", "networkidle".
func (v *Vibe) WaitForLoad(ctx context.Context, state string, timeout time.Duration) error {
	if v.closed {
		return ErrConnectionClosed
	}

	if timeout == 0 {
		timeout = DefaultTimeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"state":   state,
		"timeout": timeout.Milliseconds(),
	}

	_, err = v.client.Send(ctx, "vibium:page.waitForLoad", params)
	return err
}

// WaitForFunction waits for a JavaScript function to return a truthy value.
func (v *Vibe) WaitForFunction(ctx context.Context, fn string, timeout time.Duration) error {
	if v.closed {
		return ErrConnectionClosed
	}

	if timeout == 0 {
		timeout = DefaultTimeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"fn":      fn,
		"timeout": timeout.Milliseconds(),
	}

	_, err = v.client.Send(ctx, "vibium:page.waitForFunction", params)
	return err
}

// RouteHandler is called when a request matches a route pattern.
type RouteHandler func(ctx context.Context, route *Route) error

// Route registers a handler for requests matching the URL pattern.
// The pattern can be a glob pattern (e.g., "**/*.png") or regex (e.g., "/api/.*").
func (v *Vibe) Route(ctx context.Context, pattern string, handler RouteHandler) error {
	if v.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"pattern": pattern,
	}

	_, err = v.client.Send(ctx, "vibium:network.route", params)
	return err
}

// Unroute removes a previously registered route handler.
func (v *Vibe) Unroute(ctx context.Context, pattern string) error {
	if v.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"pattern": pattern,
	}

	_, err = v.client.Send(ctx, "vibium:network.unroute", params)
	return err
}

// SetExtraHTTPHeaders sets extra HTTP headers that will be sent with every request.
func (v *Vibe) SetExtraHTTPHeaders(ctx context.Context, headers map[string]string) error {
	if v.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"headers": headers,
	}

	_, err = v.client.Send(ctx, "vibium:network.setHeaders", params)
	return err
}

// RequestHandler is called for each network request.
type RequestHandler func(*Request)

// ResponseHandler is called for each network response.
type ResponseHandler func(*Response)

// ConsoleHandler is called for each console message.
type ConsoleHandler func(*ConsoleMessage)

// DialogHandler is called when a dialog appears.
type DialogHandler func(*Dialog)

// DownloadHandler is called when a download starts.
type DownloadHandler func(*Download)

// OnRequest registers a handler for network requests.
// Note: This is a convenience method; for full control use Route().
func (v *Vibe) OnRequest(ctx context.Context, handler RequestHandler) error {
	if v.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	_, err = v.client.Send(ctx, "vibium:network.onRequest", params)
	return err
}

// OnResponse registers a handler for network responses.
func (v *Vibe) OnResponse(ctx context.Context, handler ResponseHandler) error {
	if v.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	_, err = v.client.Send(ctx, "vibium:network.onResponse", params)
	return err
}

// OnConsole registers a handler for console messages.
func (v *Vibe) OnConsole(ctx context.Context, handler ConsoleHandler) error {
	if v.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	_, err = v.client.Send(ctx, "vibium:console.on", params)
	return err
}

// OnDialog registers a handler for dialogs (alert, confirm, prompt).
func (v *Vibe) OnDialog(ctx context.Context, handler DialogHandler) error {
	if v.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	_, err = v.client.Send(ctx, "vibium:dialog.on", params)
	return err
}

// OnDownload registers a handler for downloads.
func (v *Vibe) OnDownload(ctx context.Context, handler DownloadHandler) error {
	if v.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	_, err = v.client.Send(ctx, "vibium:download.on", params)
	return err
}

// NewPage creates a new page in the default browser context.
func (v *Vibe) NewPage(ctx context.Context) (*Vibe, error) {
	if v.closed {
		return nil, ErrConnectionClosed
	}

	result, err := v.client.Send(ctx, "browsingContext.create", map[string]interface{}{
		"type": "tab",
	})
	if err != nil {
		return nil, err
	}

	var resp struct {
		Context string `json:"context"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, err
	}

	return &Vibe{
		client:          v.client,
		clicker:         v.clicker,
		browsingContext: resp.Context,
	}, nil
}

// NewContext creates a new isolated browser context.
func (v *Vibe) NewContext(ctx context.Context) (*BrowserContext, error) {
	if v.closed {
		return nil, ErrConnectionClosed
	}

	result, err := v.client.Send(ctx, "browser.createUserContext", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	var resp struct {
		UserContext string `json:"userContext"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, err
	}

	return &BrowserContext{
		client:      v.client,
		clicker:     v.clicker,
		userContext: resp.UserContext,
	}, nil
}

// Pages returns all open pages.
func (v *Vibe) Pages(ctx context.Context) ([]*Vibe, error) {
	if v.closed {
		return nil, ErrConnectionClosed
	}

	result, err := v.client.Send(ctx, "browsingContext.getTree", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	var tree struct {
		Contexts []struct {
			Context string `json:"context"`
		} `json:"contexts"`
	}
	if err := json.Unmarshal(result, &tree); err != nil {
		return nil, err
	}

	pages := make([]*Vibe, len(tree.Contexts))
	for i, c := range tree.Contexts {
		pages[i] = &Vibe{
			client:          v.client,
			clicker:         v.clicker,
			browsingContext: c.Context,
		}
	}

	return pages, nil
}

// Context returns the browser context for this page.
// Returns nil if this is the default context.
func (v *Vibe) Context() *BrowserContext {
	// Default context doesn't have a BrowserContext wrapper
	return nil
}

// BrowserVersion returns the browser version string.
func (v *Vibe) BrowserVersion(ctx context.Context) (string, error) {
	if v.closed {
		return "", ErrConnectionClosed
	}

	result, err := v.client.Send(ctx, "browser.getUserContexts", map[string]interface{}{})
	if err != nil {
		// Fallback to just returning a placeholder
		return "", err
	}

	var resp struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", err
	}

	return resp.Version, nil
}
