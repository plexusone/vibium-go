package vibium

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/plexusone/vibium-go/launcher"
)

// Vibe is the main browser control interface.
type Vibe struct {
	client          *BiDiClient
	browser         *launcher.Browser
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

	// Start browser process
	launchOpts := &launcher.Options{
		Headless:       opts.Headless,
		Port:           opts.Port,
		ExecutablePath: opts.ExecutablePath,
		UserDataDir:    opts.UserDataDir,
		Args:           opts.Args,
		AutoInstall:    opts.AutoInstall,
	}

	browser, err := launcher.Launch(ctx, launchOpts)
	if err != nil {
		return nil, err
	}
	debugLog(ctx, "browser started", "url", browser.WebSocketURL())

	// Connect BiDi client
	client := NewBiDiClient()
	if err := client.Connect(ctx, browser.WebSocketURL()); err != nil {
		_ = browser.Stop()
		return nil, err
	}
	debugLog(ctx, "BiDi client connected")

	return &Vibe{
		client:  client,
		browser: browser,
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

// FindAll finds all elements matching the selector and optional semantic options.
// If selector is empty but semantic options are provided, elements are found by those options.
func (v *Vibe) FindAll(ctx context.Context, selector string, opts *FindOptions) ([]*Element, error) {
	if v.closed {
		return nil, ErrConnectionClosed
	}
	debugLog(ctx, "finding all elements", "selector", selector)

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return nil, err
	}

	timeout := DefaultTimeout
	if opts != nil && opts.Timeout > 0 {
		timeout = opts.Timeout
	}

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

	result, err := v.client.Send(ctx, "vibium:findAll", params)
	if err != nil {
		return nil, err
	}

	// Parse the response containing element data
	var items []struct {
		Index    int         `json:"index"`
		Selector string      `json:"selector"`
		Tag      string      `json:"tag"`
		Text     string      `json:"text"`
		Box      BoundingBox `json:"box"`
	}
	if err := json.Unmarshal(result, &items); err != nil {
		return nil, fmt.Errorf("failed to parse elements: %w", err)
	}

	elements := make([]*Element, len(items))
	for i, item := range items {
		// Use the selector returned by the server, or create an indexed one
		elemSelector := item.Selector
		if elemSelector == "" {
			elemSelector = fmt.Sprintf("%s:nth-of-type(%d)", selector, item.Index+1)
		}
		info := ElementInfo{
			Tag:  item.Tag,
			Text: item.Text,
			Box:  item.Box,
		}
		elements[i] = NewElement(v.client, browsingCtx, elemSelector, info)
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

	// Stop browser process
	if v.browser != nil {
		if err := v.browser.Stop(); err != nil {
			return err
		}
	}

	return clientErr
}

// IsClosed returns whether the browser has been closed.
func (v *Vibe) IsClosed() bool {
	return v.closed
}

// BrowsingContext returns the browsing context ID for this page.
func (v *Vibe) BrowsingContext() string {
	return v.browsingContext
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
		browser:         v.browser,
		browsingContext: resp.Context,
	}, nil
}

// A11yTree returns the accessibility tree for the page.
// Options can filter the tree to only interesting nodes or specify a root element.
func (v *Vibe) A11yTree(ctx context.Context, opts *A11yTreeOptions) (interface{}, error) {
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
		if opts.InterestingOnly != nil {
			params["interestingOnly"] = *opts.InterestingOnly
		}
		if opts.Root != "" {
			params["root"] = opts.Root
		}
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

// MainFrame returns the main frame of the page.
// Since Vibe represents both page and frame in this SDK, it returns itself.
// This method exists for API compatibility with other Vibium clients.
func (v *Vibe) MainFrame() *Vibe {
	return v
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
	if opts.Contrast != "" {
		params["contrast"] = opts.Contrast
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

// MockRouteOptions configures a static mock response for a route.
type MockRouteOptions struct {
	Status      int               // HTTP status code (default: 200)
	Body        string            // Response body
	ContentType string            // Content-Type header (default: application/json)
	Headers     map[string]string // Additional response headers
}

// MockRoute registers a route that returns a static mock response.
// This is useful for MCP tools and testing without callbacks.
func (v *Vibe) MockRoute(ctx context.Context, pattern string, opts MockRouteOptions) error {
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

	if opts.Status != 0 {
		params["status"] = opts.Status
	} else {
		params["status"] = 200
	}

	if opts.Body != "" {
		params["body"] = opts.Body
	}

	if opts.ContentType != "" {
		params["contentType"] = opts.ContentType
	} else {
		params["contentType"] = "application/json"
	}

	if opts.Headers != nil {
		params["headers"] = opts.Headers
	}

	_, err = v.client.Send(ctx, "vibium:network.mockRoute", params)
	return err
}

// RouteInfo represents information about an active route.
type RouteInfo struct {
	Pattern     string `json:"pattern"`
	Status      int    `json:"status,omitempty"`
	ContentType string `json:"contentType,omitempty"`
}

// ListRoutes returns all active route handlers.
func (v *Vibe) ListRoutes(ctx context.Context) ([]RouteInfo, error) {
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

	result, err := v.client.Send(ctx, "vibium:network.listRoutes", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Routes []RouteInfo `json:"routes"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, err
	}

	return resp.Routes, nil
}

// SetOffline sets the browser's offline mode.
func (v *Vibe) SetOffline(ctx context.Context, offline bool) error {
	if v.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"offline": offline,
	}

	_, err = v.client.Send(ctx, "vibium:network.setOffline", params)
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

// PageErrorHandler is called when a JavaScript error occurs on the page.
type PageErrorHandler func(*PageError)

// PageHandler is called when a new page is created.
type PageHandler func(*Vibe)

// PopupHandler is called when a popup window opens.
type PopupHandler func(*Vibe)

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

	// Register event handler with BiDi client
	v.client.OnEvent("vibium:network.request", func(event *BiDiEvent) {
		var req Request
		if err := json.Unmarshal(event.Params, &req); err != nil {
			debugLog(ctx, "failed to unmarshal request event", "error", err)
			return
		}
		handler(&req)
	})

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

	// Register event handler with BiDi client
	v.client.OnEvent("vibium:network.response", func(event *BiDiEvent) {
		var resp Response
		if err := json.Unmarshal(event.Params, &resp); err != nil {
			debugLog(ctx, "failed to unmarshal response event", "error", err)
			return
		}
		handler(&resp)
	})

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

	// Register event handler with BiDi client
	v.client.OnEvent("vibium:console.entry", func(event *BiDiEvent) {
		var msg ConsoleMessage
		if err := json.Unmarshal(event.Params, &msg); err != nil {
			debugLog(ctx, "failed to unmarshal console event", "error", err)
			return
		}
		handler(&msg)
	})

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

	// Register event handler with BiDi client
	v.client.OnEvent("vibium:dialog.opened", func(event *BiDiEvent) {
		var dialog Dialog
		if err := json.Unmarshal(event.Params, &dialog); err != nil {
			debugLog(ctx, "failed to unmarshal dialog event", "error", err)
			return
		}
		handler(&dialog)
	})

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

	// Register event handler with BiDi client
	v.client.OnEvent("vibium:download.started", func(event *BiDiEvent) {
		var download Download
		if err := json.Unmarshal(event.Params, &download); err != nil {
			debugLog(ctx, "failed to unmarshal download event", "error", err)
			return
		}
		handler(&download)
	})

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	_, err = v.client.Send(ctx, "vibium:download.on", params)
	return err
}

// OnError registers a handler for JavaScript errors on the page.
func (v *Vibe) OnError(ctx context.Context, handler PageErrorHandler) error {
	if v.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	// Register event handler with BiDi client
	v.client.OnEvent("vibium:page.error", func(event *BiDiEvent) {
		var pageErr PageError
		if err := json.Unmarshal(event.Params, &pageErr); err != nil {
			debugLog(ctx, "failed to unmarshal page error event", "error", err)
			return
		}
		handler(&pageErr)
	})

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	_, err = v.client.Send(ctx, "vibium:page.onError", params)
	return err
}

// CollectConsole enables buffered console message collection.
// Messages can be retrieved with ConsoleMessages() and cleared with ClearConsoleMessages().
func (v *Vibe) CollectConsole(ctx context.Context) error {
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

	_, err = v.client.Send(ctx, "vibium:console.collect", params)
	return err
}

// CollectErrors enables buffered page error collection.
// Errors can be retrieved with Errors() and cleared with ClearErrors().
func (v *Vibe) CollectErrors(ctx context.Context) error {
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

	_, err = v.client.Send(ctx, "vibium:page.collectErrors", params)
	return err
}

// Errors retrieves buffered page errors.
// Call CollectErrors() first to enable error collection.
func (v *Vibe) Errors(ctx context.Context) ([]PageError, error) {
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

	result, err := v.client.Send(ctx, "vibium:page.errors", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Errors []PageError `json:"errors"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse errors: %w", err)
	}

	return resp.Errors, nil
}

// ClearErrors clears the buffered page errors.
func (v *Vibe) ClearErrors(ctx context.Context) error {
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

	_, err = v.client.Send(ctx, "vibium:page.clearErrors", params)
	return err
}

// OnPage registers a handler that is called when a new page is created in the browser.
// This includes pages created via NewPage(), window.open(), or clicking links with target="_blank".
func (v *Vibe) OnPage(ctx context.Context, handler PageHandler) error {
	if v.closed {
		return ErrConnectionClosed
	}

	// Register event handler with BiDi client
	v.client.OnEvent("browsingContext.contextCreated", func(event *BiDiEvent) {
		var params struct {
			Context string `json:"context"`
			URL     string `json:"url"`
			Parent  string `json:"parent,omitempty"`
		}
		if err := json.Unmarshal(event.Params, &params); err != nil {
			debugLog(ctx, "failed to unmarshal page created event", "error", err)
			return
		}

		// Create a new Vibe instance for the new page
		newPage := &Vibe{
			client:          v.client,
			browser:         v.browser,
			browsingContext: params.Context,
		}
		handler(newPage)
	})

	// Subscribe to context created events
	_, err := v.client.Send(ctx, "session.subscribe", map[string]interface{}{
		"events": []string{"browsingContext.contextCreated"},
	})
	return err
}

// OnPopup registers a handler that is called when a popup window is opened.
// Popups are typically created via window.open() with specific features.
func (v *Vibe) OnPopup(ctx context.Context, handler PopupHandler) error {
	if v.closed {
		return ErrConnectionClosed
	}

	// Register event handler with BiDi client
	v.client.OnEvent("browsingContext.contextCreated", func(event *BiDiEvent) {
		var params struct {
			Context  string `json:"context"`
			URL      string `json:"url"`
			Parent   string `json:"parent,omitempty"`
			Original string `json:"originalOpener,omitempty"`
		}
		if err := json.Unmarshal(event.Params, &params); err != nil {
			debugLog(ctx, "failed to unmarshal popup event", "error", err)
			return
		}

		// Only call handler for popups (contexts with an opener)
		if params.Original == "" && params.Parent == "" {
			return
		}

		// Create a new Vibe instance for the popup
		popup := &Vibe{
			client:          v.client,
			browser:         v.browser,
			browsingContext: params.Context,
		}
		handler(popup)
	})

	// Subscribe to context created events
	_, err := v.client.Send(ctx, "session.subscribe", map[string]interface{}{
		"events": []string{"browsingContext.contextCreated"},
	})
	return err
}

// RemoveAllListeners removes all registered event listeners.
// This is useful for cleanup when you no longer need to receive events.
func (v *Vibe) RemoveAllListeners() {
	if v.client != nil {
		// Remove all handlers from the BiDi client
		v.client.handlerMu.Lock()
		v.client.handlers = make(map[string][]EventHandler)
		v.client.handlerMu.Unlock()
	}
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
		browser:         v.browser,
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
		browser:     v.browser,
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
			browser:         v.browser,
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

// HandleDialog handles the current dialog by accepting or dismissing it.
// If accept is true, the dialog is accepted. If promptText is provided (for prompt dialogs),
// it will be entered before accepting.
func (v *Vibe) HandleDialog(ctx context.Context, accept bool, promptText string) error {
	if v.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"accept":  accept,
	}

	if accept && promptText != "" {
		params["userText"] = promptText
	}

	_, err = v.client.Send(ctx, "vibium:dialog.handle", params)
	return err
}

// GetDialog returns information about the current dialog, if any.
func (v *Vibe) GetDialog(ctx context.Context) (DialogInfo, error) {
	if v.closed {
		return DialogInfo{}, ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return DialogInfo{}, err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	result, err := v.client.Send(ctx, "vibium:dialog.get", params)
	if err != nil {
		// No dialog open is not an error
		return DialogInfo{HasDialog: false}, nil
	}

	var resp struct {
		Type         string `json:"type"`
		Message      string `json:"message"`
		DefaultValue string `json:"defaultValue"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return DialogInfo{HasDialog: false}, nil
	}

	if resp.Type == "" {
		return DialogInfo{HasDialog: false}, nil
	}

	return DialogInfo{
		HasDialog:    true,
		Type:         resp.Type,
		Message:      resp.Message,
		DefaultValue: resp.DefaultValue,
	}, nil
}

// ConsoleMessages returns buffered console messages from the page.
// The level parameter filters messages by type (log, info, warn, error, debug).
// If level is empty, all messages are returned.
func (v *Vibe) ConsoleMessages(ctx context.Context, level string) ([]ConsoleMessage, error) {
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

	if level != "" {
		params["level"] = level
	}

	result, err := v.client.Send(ctx, "vibium:console.messages", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Messages []ConsoleMessage `json:"messages"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, err
	}

	return resp.Messages, nil
}

// ClearConsoleMessages clears the buffered console messages.
func (v *Vibe) ClearConsoleMessages(ctx context.Context) error {
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

	_, err = v.client.Send(ctx, "vibium:console.clear", params)
	return err
}

// NetworkRequest represents a captured network request with its response.
type NetworkRequest struct {
	URL          string            `json:"url"`
	Method       string            `json:"method"`
	Headers      map[string]string `json:"headers,omitempty"`
	PostData     string            `json:"postData,omitempty"`
	ResourceType string            `json:"resourceType"`
	Status       int               `json:"status,omitempty"`
	StatusText   string            `json:"statusText,omitempty"`
	ResponseSize int64             `json:"responseSize,omitempty"`
	Timestamp    int64             `json:"timestamp,omitempty"`
}

// NetworkRequests returns buffered network requests from the page.
// Options can filter by URL pattern, method, or resource type.
func (v *Vibe) NetworkRequests(ctx context.Context, opts *NetworkRequestsOptions) ([]NetworkRequest, error) {
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
		if opts.URLPattern != "" {
			params["urlPattern"] = opts.URLPattern
		}
		if opts.Method != "" {
			params["method"] = opts.Method
		}
		if opts.ResourceType != "" {
			params["resourceType"] = opts.ResourceType
		}
	}

	result, err := v.client.Send(ctx, "vibium:network.requests", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Requests []NetworkRequest `json:"requests"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, err
	}

	return resp.Requests, nil
}

// NetworkRequestsOptions configures network request filtering.
type NetworkRequestsOptions struct {
	URLPattern   string // Glob or regex pattern to filter URLs
	Method       string // Filter by HTTP method (GET, POST, etc.)
	ResourceType string // Filter by resource type (document, script, xhr, etc.)
}

// ClearNetworkRequests clears the buffered network requests.
func (v *Vibe) ClearNetworkRequests(ctx context.Context) error {
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

	_, err = v.client.Send(ctx, "vibium:network.clearRequests", params)
	return err
}

// ScrollOptions configures scroll behavior.
type ScrollOptions struct {
	Selector string // Optional CSS selector to scroll within
}

// Scroll scrolls the page or a specific element.
// direction can be "up", "down", "left", or "right".
// amount is the number of pixels to scroll (use 0 for full page).
func (v *Vibe) Scroll(ctx context.Context, direction string, amount int, opts *ScrollOptions) error {
	if v.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context":   browsingCtx,
		"direction": direction,
		"amount":    amount,
	}

	if opts != nil && opts.Selector != "" {
		params["selector"] = opts.Selector
	}

	_, err = v.client.Send(ctx, "vibium:page.scroll", params)
	return err
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

// Tracing returns a tracing controller for the default browser context.
// Use this to record traces for debugging and analysis.
func (v *Vibe) Tracing() *Tracing {
	return &Tracing{
		client:      v.client,
		userContext: "", // Empty string uses the default user context
	}
}

// AddInitScript adds a script that will be evaluated in every page before any page scripts.
// This is useful for mocking APIs, injecting test helpers, or setting up authentication.
func (v *Vibe) AddInitScript(ctx context.Context, script string) error {
	if v.closed {
		return ErrConnectionClosed
	}

	// Get the default user context
	userContext, err := v.getDefaultUserContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get user context: %w", err)
	}

	params := map[string]interface{}{
		"userContext": userContext,
		"script":      script,
	}

	_, err = v.client.Send(ctx, "vibium:context.addInitScript", params)
	return err
}

// getDefaultUserContext returns the default user context ID.
func (v *Vibe) getDefaultUserContext(ctx context.Context) (string, error) {
	result, err := v.client.Send(ctx, "browser.getUserContexts", map[string]interface{}{})
	if err != nil {
		return "", err
	}

	var resp struct {
		UserContexts []struct {
			UserContext string `json:"userContext"`
		} `json:"userContexts"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", err
	}

	if len(resp.UserContexts) == 0 {
		return "", fmt.Errorf("no user contexts available")
	}

	// Return the first (default) user context
	return resp.UserContexts[0].UserContext, nil
}

// StorageState returns the complete browser storage state including cookies, localStorage,
// and sessionStorage for the current page's origin. This can be saved and later restored
// using SetStorageState to resume a session.
func (v *Vibe) StorageState(ctx context.Context) (*StorageState, error) {
	if v.closed {
		return nil, ErrConnectionClosed
	}

	// Get base storage state (cookies + localStorage) from context
	browserCtx, err := v.NewContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get browser context: %w", err)
	}

	state, err := browserCtx.StorageState(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage state: %w", err)
	}

	// Capture sessionStorage for the current page's origin
	currentURL, err := v.URL(ctx)
	if err != nil || currentURL == "" || currentURL == "about:blank" {
		// No page loaded, return state without sessionStorage
		return state, nil
	}

	// Get sessionStorage via JavaScript
	sessionStorageScript := `
		(function() {
			const items = {};
			for (let i = 0; i < sessionStorage.length; i++) {
				const key = sessionStorage.key(i);
				items[key] = sessionStorage.getItem(key);
			}
			return JSON.stringify({
				origin: window.location.origin,
				items: items
			});
		})()
	`

	result, err := v.Evaluate(ctx, sessionStorageScript)
	if err != nil {
		// sessionStorage not available (e.g., file:// protocol), return without it
		return state, nil
	}

	// Parse the sessionStorage result
	resultStr, ok := result.(string)
	if !ok {
		return state, nil
	}

	var sessionData struct {
		Origin string            `json:"origin"`
		Items  map[string]string `json:"items"`
	}
	if err := json.Unmarshal([]byte(resultStr), &sessionData); err != nil {
		return state, nil
	}

	// Merge sessionStorage into the appropriate origin
	if len(sessionData.Items) > 0 {
		found := false
		for i := range state.Origins {
			if state.Origins[i].Origin == sessionData.Origin {
				state.Origins[i].SessionStorage = sessionData.Items
				found = true
				break
			}
		}
		if !found {
			// Add new origin entry for sessionStorage
			state.Origins = append(state.Origins, StorageStateOrigin{
				Origin:         sessionData.Origin,
				LocalStorage:   map[string]string{},
				SessionStorage: sessionData.Items,
			})
		}
	}

	return state, nil
}

// SetStorageState restores browser storage state from a previously saved StorageState.
// This includes cookies, localStorage, and sessionStorage. The browser should be on
// a page (or will be navigated to the first origin) for storage to be set correctly.
func (v *Vibe) SetStorageState(ctx context.Context, state *StorageState) error {
	if v.closed {
		return ErrConnectionClosed
	}

	browserCtx, err := v.NewContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get browser context: %w", err)
	}

	// Set cookies
	if len(state.Cookies) > 0 {
		cookies := make([]SetCookieParam, len(state.Cookies))
		for i, c := range state.Cookies {
			cookies[i] = SetCookieParam{
				Name:     c.Name,
				Value:    c.Value,
				Domain:   c.Domain,
				Path:     c.Path,
				Expires:  c.Expires,
				HTTPOnly: c.HTTPOnly,
				Secure:   c.Secure,
				SameSite: c.SameSite,
			}
		}
		if err := browserCtx.SetCookies(ctx, cookies); err != nil {
			return fmt.Errorf("failed to set cookies: %w", err)
		}
	}

	// Set localStorage and sessionStorage for each origin
	for _, origin := range state.Origins {
		hasLocalStorage := len(origin.LocalStorage) > 0
		hasSessionStorage := len(origin.SessionStorage) > 0

		if !hasLocalStorage && !hasSessionStorage {
			continue
		}

		// Check current URL - we may need to navigate to the origin
		currentURL, _ := v.URL(ctx)
		if currentURL == "" || currentURL == "about:blank" {
			// Navigate to the origin to set storage
			if err := v.Go(ctx, origin.Origin); err != nil {
				return fmt.Errorf("failed to navigate to origin %s: %w", origin.Origin, err)
			}
		}

		// Set localStorage
		if hasLocalStorage {
			localStorageJSON, err := json.Marshal(origin.LocalStorage)
			if err != nil {
				return fmt.Errorf("failed to marshal localStorage: %w", err)
			}

			script := fmt.Sprintf(`
				(function() {
					const items = %s;
					for (const [key, value] of Object.entries(items)) {
						localStorage.setItem(key, value);
					}
					return Object.keys(items).length;
				})()
			`, string(localStorageJSON))

			if _, err := v.Evaluate(ctx, script); err != nil {
				return fmt.Errorf("failed to set localStorage for %s: %w", origin.Origin, err)
			}
		}

		// Set sessionStorage
		if hasSessionStorage {
			sessionStorageJSON, err := json.Marshal(origin.SessionStorage)
			if err != nil {
				return fmt.Errorf("failed to marshal sessionStorage: %w", err)
			}

			script := fmt.Sprintf(`
				(function() {
					const items = %s;
					for (const [key, value] of Object.entries(items)) {
						sessionStorage.setItem(key, value);
					}
					return Object.keys(items).length;
				})()
			`, string(sessionStorageJSON))

			if _, err := v.Evaluate(ctx, script); err != nil {
				return fmt.Errorf("failed to set sessionStorage for %s: %w", origin.Origin, err)
			}
		}
	}

	return nil
}

// ClearStorage clears all cookies, localStorage, and sessionStorage.
func (v *Vibe) ClearStorage(ctx context.Context) error {
	if v.closed {
		return ErrConnectionClosed
	}

	browserCtx, err := v.NewContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get browser context: %w", err)
	}

	// Clear cookies
	if err := browserCtx.ClearCookies(ctx); err != nil {
		return fmt.Errorf("failed to clear cookies: %w", err)
	}

	// Clear localStorage and sessionStorage for current page
	script := `
		(function() {
			localStorage.clear();
			sessionStorage.clear();
			return true;
		})()
	`
	if _, err := v.Evaluate(ctx, script); err != nil {
		// Ignore errors (e.g., about:blank page)
		return nil
	}

	return nil
}
