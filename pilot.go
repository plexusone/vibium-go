package w3pilot

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/plexusone/w3pilot/cdp"
)

// Pilot is the main browser control interface.
type Pilot struct {
	client          *BiDiClient
	pipeTransport   *pipeTransport  // Used in pipe mode (default)
	clicker         *ClickerProcess // Used in WebSocket mode
	browsingContext string
	closed          bool

	// CDP client for direct Chrome DevTools Protocol access
	cdpClient *cdp.Client
	cdpPort   int

	// Input controllers (lazy-initialized)
	keyboard *Keyboard
	mouse    *Mouse
	touch    *Touch
	clock    *Clock

	// CDP screencast manager (lazy-initialized)
	screencast *cdp.Screencast

	// CDP coverage manager (lazy-initialized)
	coverage *cdp.Coverage

	// CDP console debugger (lazy-initialized)
	consoleDebugger *cdp.ConsoleDebugger
}

// Browser provides browser launching capabilities.
var Browser = &browserLauncher{}

type browserLauncher struct{}

// Launch starts a new browser instance and returns a Pilot for controlling it.
func (b *browserLauncher) Launch(ctx context.Context, opts *LaunchOptions) (*Pilot, error) {
	if opts == nil {
		opts = &LaunchOptions{}
	}

	// Set up debug logging if enabled
	if logger := NewDebugLogger(); logger != nil {
		ctx = ContextWithLogger(ctx, logger)
		debugLog(ctx, "launching browser", "headless", opts.Headless, "websocket", opts.UseWebSocket)
	}

	if opts.UseWebSocket {
		// WebSocket mode (clicker serve) - for multiple clients or debugging
		return b.launchWebSocket(ctx, opts)
	}

	// Pipe mode (clicker pipe) - default, full vibium:* command support
	return b.launchPipe(ctx, opts)
}

// launchPipe starts the browser using pipe (stdin/stdout) transport.
func (b *browserLauncher) launchPipe(ctx context.Context, opts *LaunchOptions) (*Pilot, error) {
	transport := newPipeTransport()
	pipeOpts := &PipeOptions{
		Headless:       opts.Headless,
		ExecutablePath: opts.ExecutablePath,
	}

	if err := transport.Start(ctx, pipeOpts); err != nil {
		return nil, err
	}
	debugLog(ctx, "clicker pipe started")

	client := NewBiDiClient(transport)

	pilot := &Pilot{
		client:        client,
		pipeTransport: transport,
	}

	connectCDP(ctx, pilot)
	return pilot, nil
}

// launchWebSocket starts the browser using WebSocket transport.
func (b *browserLauncher) launchWebSocket(ctx context.Context, opts *LaunchOptions) (*Pilot, error) {
	clicker, err := StartClicker(ctx, LaunchOptions{
		Headless:       opts.Headless,
		Port:           opts.Port,
		ExecutablePath: opts.ExecutablePath,
	})
	if err != nil {
		return nil, err
	}
	debugLog(ctx, "clicker started", "url", clicker.WebSocketURL())

	// Connect WebSocket transport for BiDi
	wsTransport := newWSTransport()
	if err := wsTransport.Connect(ctx, clicker.WebSocketURL()); err != nil {
		_ = clicker.Stop()
		return nil, err
	}

	// Wait for browser to be ready
	if err := wsTransport.WaitForReady(ctx, 30*time.Second); err != nil {
		_ = wsTransport.Close()
		_ = clicker.Stop()
		return nil, fmt.Errorf("browser not ready: %w", err)
	}
	debugLog(ctx, "browser ready")

	client := NewBiDiClient(wsTransport)

	pilot := &Pilot{
		client:  client,
		clicker: clicker,
	}

	connectCDP(ctx, pilot)
	return pilot, nil
}

// Connect connects to an existing clicker instance via WebSocket URL.
// This is used to reconnect to a browser session that was previously launched.
func (b *browserLauncher) Connect(ctx context.Context, wsURL string) (*Pilot, error) {
	if logger := NewDebugLogger(); logger != nil {
		ctx = ContextWithLogger(ctx, logger)
		debugLog(ctx, "connecting to existing browser", "url", wsURL)
	}

	// Connect WebSocket transport for BiDi
	wsTransport := newWSTransport()
	if err := wsTransport.Connect(ctx, wsURL); err != nil {
		return nil, fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	client := NewBiDiClient(wsTransport)

	pilot := &Pilot{
		client: client,
	}

	// Try to discover existing browsing context
	result, err := client.Send(ctx, "browsingContext.getTree", map[string]interface{}{})
	if err != nil {
		_ = wsTransport.Close()
		return nil, fmt.Errorf("failed to get browsing context: %w", err)
	}

	var tree struct {
		Contexts []struct {
			Context string `json:"context"`
		} `json:"contexts"`
	}
	if err := json.Unmarshal(result, &tree); err == nil && len(tree.Contexts) > 0 {
		pilot.browsingContext = tree.Contexts[0].Context
	}

	// Connect CDP (best-effort)
	connectCDP(ctx, pilot)

	debugLog(ctx, "connected to existing browser", "context", pilot.browsingContext)
	return pilot, nil
}

// Connect is a convenience function that connects to an existing browser.
func Connect(ctx context.Context, wsURL string) (*Pilot, error) {
	return Browser.Connect(ctx, wsURL)
}

// connectCDP discovers and connects the CDP client (best-effort).
func connectCDP(ctx context.Context, pilot *Pilot) {
	// Give Chrome a moment to start and write DevToolsActivePort
	time.Sleep(500 * time.Millisecond)

	cdpPort, wsEndpoint, err := cdp.DiscoverFromRunningChrome()
	if err == nil {
		cdpClient := cdp.NewClient()
		if err := cdpClient.Connect(ctx, wsEndpoint); err == nil {
			pilot.cdpClient = cdpClient
			pilot.cdpPort = cdpPort
			debugLog(ctx, "CDP client connected", "port", cdpPort)
		} else {
			debugLog(ctx, "CDP connection failed (continuing without CDP)", "error", err)
		}
	} else {
		debugLog(ctx, "CDP discovery failed (continuing without CDP)", "error", err)
	}
}

// Launch is a convenience function that launches a browser with default options.
func Launch(ctx context.Context) (*Pilot, error) {
	return Browser.Launch(ctx, nil)
}

// LaunchHeadless is a convenience function that launches a headless browser.
func LaunchHeadless(ctx context.Context) (*Pilot, error) {
	return Browser.Launch(ctx, &LaunchOptions{Headless: true})
}

// getContext returns the browsing context ID, fetching it if necessary.
func (p *Pilot) getContext(ctx context.Context) (string, error) {
	if p.browsingContext != "" {
		return p.browsingContext, nil
	}

	result, err := p.client.Send(ctx, "browsingContext.getTree", map[string]interface{}{})
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

	p.browsingContext = tree.Contexts[0].Context
	return p.browsingContext, nil
}

// Go navigates to the specified URL.
func (p *Pilot) Go(ctx context.Context, url string) error {
	if p.closed {
		return ErrConnectionClosed
	}
	debugLog(ctx, "navigating", "url", url)

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"url":     url,
		"wait":    "complete",
	}

	_, err = p.client.Send(ctx, "browsingContext.navigate", params)
	if err == nil {
		debugLog(ctx, "navigation complete", "url", url)
	}
	return err
}

// Reload reloads the current page.
func (p *Pilot) Reload(ctx context.Context) error {
	if p.closed {
		return ErrConnectionClosed
	}
	debugLog(ctx, "reloading page")

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"wait":    "complete",
	}

	_, err = p.client.Send(ctx, "browsingContext.reload", params)
	return err
}

// Back navigates back in history.
func (p *Pilot) Back(ctx context.Context) error {
	if p.closed {
		return ErrConnectionClosed
	}
	debugLog(ctx, "navigating back")

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"delta":   -1,
	}

	_, err = p.client.Send(ctx, "browsingContext.traverseHistory", params)
	return err
}

// Forward navigates forward in history.
func (p *Pilot) Forward(ctx context.Context) error {
	if p.closed {
		return ErrConnectionClosed
	}
	debugLog(ctx, "navigating forward")

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"delta":   1,
	}

	_, err = p.client.Send(ctx, "browsingContext.traverseHistory", params)
	return err
}

// Screenshot captures a screenshot of the current page and returns PNG data.
func (p *Pilot) Screenshot(ctx context.Context) ([]byte, error) {
	if p.closed {
		return nil, ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return nil, err
	}

	result, err := p.client.Send(ctx, "browsingContext.captureScreenshot", map[string]interface{}{
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
func (p *Pilot) Find(ctx context.Context, selector string, opts *FindOptions) (*Element, error) {
	if p.closed {
		return nil, ErrConnectionClosed
	}
	debugLog(ctx, "finding element", "selector", selector)

	browsingCtx, err := p.getContext(ctx)
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

	result, err := p.client.Send(ctx, "vibium:page.find", params)
	if err != nil {
		return nil, err
	}

	var info ElementInfo
	if err := json.Unmarshal(result, &info); err != nil {
		return nil, fmt.Errorf("failed to parse element info: %w", err)
	}

	debugLog(ctx, "element found", "selector", selector, "tag", info.Tag)
	return NewElement(p.client, browsingCtx, selector, info), nil
}

// FindAll finds all elements matching the selector and optional semantic options.
// If selector is empty but semantic options are provided, elements are found by those options.
func (p *Pilot) FindAll(ctx context.Context, selector string, opts *FindOptions) ([]*Element, error) {
	if p.closed {
		return nil, ErrConnectionClosed
	}
	debugLog(ctx, "finding all elements", "selector", selector)

	browsingCtx, err := p.getContext(ctx)
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

	result, err := p.client.Send(ctx, "vibium:page.findAll", params)
	if err != nil {
		return nil, err
	}

	// Parse the response containing element data
	// Response format: {"elements": [...], "count": N}
	var resp struct {
		Elements []struct {
			Index    int         `json:"index"`
			Selector string      `json:"selector"`
			Tag      string      `json:"tag"`
			Text     string      `json:"text"`
			Box      BoundingBox `json:"box"`
		} `json:"elements"`
		Count int `json:"count"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse elements: %w", err)
	}
	items := resp.Elements

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
		elements[i] = NewElement(p.client, browsingCtx, elemSelector, info)
	}

	debugLog(ctx, "elements found", "selector", selector, "count", len(elements))
	return elements, nil
}

// MustFind finds an element by CSS selector and panics if not found.
func (p *Pilot) MustFind(ctx context.Context, selector string) *Element {
	elem, err := p.Find(ctx, selector, nil)
	if err != nil {
		panic(err)
	}
	return elem
}

// Evaluate executes JavaScript in the page context and returns the result.
func (p *Pilot) Evaluate(ctx context.Context, script string) (interface{}, error) {
	if p.closed {
		return nil, ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return nil, err
	}

	// Wrap script in arrow function.
	// If script has return/let/const/var/if/for/while/try/throw statements, use block syntax.
	// Otherwise use expression syntax so simple expressions return their value.
	var wrappedScript string
	trimmed := strings.TrimSpace(script)
	if strings.HasPrefix(trimmed, "return ") ||
		strings.HasPrefix(trimmed, "let ") ||
		strings.HasPrefix(trimmed, "const ") ||
		strings.HasPrefix(trimmed, "var ") ||
		strings.HasPrefix(trimmed, "if ") ||
		strings.HasPrefix(trimmed, "if(") ||
		strings.HasPrefix(trimmed, "for ") ||
		strings.HasPrefix(trimmed, "for(") ||
		strings.HasPrefix(trimmed, "while ") ||
		strings.HasPrefix(trimmed, "while(") ||
		strings.HasPrefix(trimmed, "try ") ||
		strings.HasPrefix(trimmed, "try{") ||
		strings.HasPrefix(trimmed, "throw ") ||
		strings.HasPrefix(trimmed, "{") ||
		strings.Contains(trimmed, ";") {
		// Statement(s): use block syntax
		wrappedScript = fmt.Sprintf("() => { %s }", script)
	} else {
		// Expression: use expression syntax for implicit return
		wrappedScript = fmt.Sprintf("() => (%s)", script)
	}

	params := map[string]interface{}{
		"functionDeclaration": wrappedScript,
		"target":              map[string]interface{}{"context": browsingCtx},
		"arguments":           []interface{}{},
		"awaitPromise":        true,
		"resultOwnership":     "root",
	}

	result, err := p.client.Send(ctx, "script.callFunction", params)
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
func (p *Pilot) Title(ctx context.Context) (string, error) {
	result, err := p.Evaluate(ctx, "return document.title")
	if err != nil {
		return "", err
	}
	if s, ok := result.(string); ok {
		return s, nil
	}
	return "", nil
}

// URL returns the current page URL.
func (p *Pilot) URL(ctx context.Context) (string, error) {
	result, err := p.Evaluate(ctx, "return window.location.href")
	if err != nil {
		return "", err
	}
	if s, ok := result.(string); ok {
		return s, nil
	}
	return "", nil
}

// WaitForNavigation waits for a navigation to complete.
func (p *Pilot) WaitForNavigation(ctx context.Context, timeout time.Duration) error {
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
			result, err := p.Evaluate(ctx, "return document.readyState")
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
func (p *Pilot) Quit(ctx context.Context) error {
	if p.closed {
		return nil
	}
	p.closed = true

	// Close the CDP client connection
	if p.cdpClient != nil {
		_ = p.cdpClient.Close()
	}

	// Close the BiDi client connection
	if p.client != nil {
		_ = p.client.Close()
	}

	// Stop the clicker process (WebSocket mode)
	if p.clicker != nil {
		return p.clicker.Stop()
	}

	// Close pipe transport (pipe mode)
	if p.pipeTransport != nil {
		return p.pipeTransport.Close()
	}

	return nil
}

// IsClosed returns whether the browser has been closed.
func (p *Pilot) IsClosed() bool {
	return p.closed
}

// Clicker returns the clicker process, or nil if using pipe mode.
func (p *Pilot) Clicker() *ClickerProcess {
	return p.clicker
}

// BrowsingContext returns the browsing context ID for this page.
func (p *Pilot) BrowsingContext() string {
	return p.browsingContext
}

// CDP returns the Chrome DevTools Protocol client, or nil if not available.
// Use HasCDP() to check availability before calling CDP methods.
func (p *Pilot) CDP() *cdp.Client {
	return p.cdpClient
}

// HasCDP returns true if the CDP client is connected and available.
func (p *Pilot) HasCDP() bool {
	return p.cdpClient != nil && p.cdpClient.IsConnected()
}

// CDPPort returns the CDP port, or 0 if CDP is not available.
func (p *Pilot) CDPPort() int {
	return p.cdpPort
}

// TakeHeapSnapshot captures a V8 heap snapshot and saves it to a file.
// Requires CDP connection. Returns error if CDP is not available.
func (p *Pilot) TakeHeapSnapshot(ctx context.Context, path string) (*cdp.HeapSnapshot, error) {
	if !p.HasCDP() {
		return nil, fmt.Errorf("CDP not available")
	}
	return p.cdpClient.TakeHeapSnapshot(ctx, path)
}

// GetNetworkResponseBody retrieves the body of a network response.
// Requires CDP connection. Returns error if CDP is not available.
func (p *Pilot) GetNetworkResponseBody(ctx context.Context, requestID string, saveTo string) (*cdp.ResponseBody, error) {
	if !p.HasCDP() {
		return nil, fmt.Errorf("CDP not available")
	}
	return p.cdpClient.GetResponseBody(ctx, requestID, saveTo)
}

// EmulateNetwork sets network throttling conditions.
// Use cdp.NetworkSlow3G, cdp.NetworkFast3G, cdp.Network4G, etc.
// Requires CDP connection. Returns error if CDP is not available.
func (p *Pilot) EmulateNetwork(ctx context.Context, conditions cdp.NetworkConditions) error {
	if !p.HasCDP() {
		return fmt.Errorf("CDP not available")
	}
	return p.cdpClient.SetNetworkConditions(ctx, conditions)
}

// ClearNetworkEmulation clears network throttling.
// Requires CDP connection. Returns error if CDP is not available.
func (p *Pilot) ClearNetworkEmulation(ctx context.Context) error {
	if !p.HasCDP() {
		return fmt.Errorf("CDP not available")
	}
	return p.cdpClient.ClearNetworkConditions(ctx)
}

// EmulateCPU sets CPU throttling.
// rate=1 means no throttling, rate=4 means 4x slowdown.
// Use cdp.CPU4xSlowdown, cdp.CPU6xSlowdown, etc.
// Requires CDP connection. Returns error if CDP is not available.
func (p *Pilot) EmulateCPU(ctx context.Context, rate int) error {
	if !p.HasCDP() {
		return fmt.Errorf("CDP not available")
	}
	return p.cdpClient.SetCPUThrottlingRate(ctx, rate)
}

// ClearCPUEmulation clears CPU throttling.
// Requires CDP connection. Returns error if CDP is not available.
func (p *Pilot) ClearCPUEmulation(ctx context.Context) error {
	if !p.HasCDP() {
		return fmt.Errorf("CDP not available")
	}
	return p.cdpClient.ClearCPUThrottling(ctx)
}

// ScreencastFrameHandler is called for each captured screencast frame.
type ScreencastFrameHandler func(frame *cdp.ScreencastFrame)

// StartScreencast begins capturing screen frames.
// The handler is called for each captured frame with base64-encoded image data.
// Requires CDP connection. Returns error if CDP is not available.
func (p *Pilot) StartScreencast(ctx context.Context, opts *cdp.ScreencastOptions, handler ScreencastFrameHandler) error {
	if !p.HasCDP() {
		return fmt.Errorf("CDP not available")
	}
	if p.screencast == nil {
		p.screencast = cdp.NewScreencast(p.cdpClient)
	}
	return p.screencast.Start(ctx, opts, func(frame *cdp.ScreencastFrame) {
		if handler != nil {
			handler(frame)
		}
	})
}

// StopScreencast stops capturing screen frames.
// Requires CDP connection. Returns error if CDP is not available.
func (p *Pilot) StopScreencast(ctx context.Context) error {
	if !p.HasCDP() {
		return fmt.Errorf("CDP not available")
	}
	if p.screencast == nil {
		return nil
	}
	return p.screencast.Stop(ctx)
}

// IsScreencasting returns whether screencast is active.
func (p *Pilot) IsScreencasting() bool {
	if p.screencast == nil {
		return false
	}
	return p.screencast.IsRunning()
}

// ExtensionInfo contains information about a browser extension.
type ExtensionInfo = cdp.ExtensionInfo

// InstallExtension loads an unpacked extension from a directory.
// Returns the extension ID if successful.
// Requires CDP connection. Returns error if CDP is not available.
func (p *Pilot) InstallExtension(ctx context.Context, path string) (string, error) {
	if !p.HasCDP() {
		return "", fmt.Errorf("CDP not available")
	}
	return p.cdpClient.LoadUnpackedExtension(ctx, path)
}

// UninstallExtension removes an extension by ID.
// Requires CDP connection. Returns error if CDP is not available.
func (p *Pilot) UninstallExtension(ctx context.Context, id string) error {
	if !p.HasCDP() {
		return fmt.Errorf("CDP not available")
	}
	return p.cdpClient.UninstallExtension(ctx, id)
}

// ListExtensions returns all installed extensions.
// Requires CDP connection. Returns error if CDP is not available.
func (p *Pilot) ListExtensions(ctx context.Context) ([]ExtensionInfo, error) {
	if !p.HasCDP() {
		return nil, fmt.Errorf("CDP not available")
	}
	return p.cdpClient.GetAllExtensions(ctx)
}

// CoverageReport is an alias for cdp.CoverageReport.
type CoverageReport = cdp.CoverageReport

// CoverageSummary is an alias for cdp.CoverageSummary.
type CoverageSummary = cdp.CoverageSummary

// StartCoverage begins collecting JS and CSS coverage data.
// Requires CDP connection. Returns error if CDP is not available.
func (p *Pilot) StartCoverage(ctx context.Context) error {
	if !p.HasCDP() {
		return fmt.Errorf("CDP not available")
	}
	if p.coverage == nil {
		p.coverage = cdp.NewCoverage(p.cdpClient)
	}
	return p.coverage.Start(ctx)
}

// StartJSCoverage begins collecting JavaScript coverage data.
// callCount: collect execution counts per block
// detailed: collect block-level coverage (vs function-level)
// Requires CDP connection. Returns error if CDP is not available.
func (p *Pilot) StartJSCoverage(ctx context.Context, callCount, detailed bool) error {
	if !p.HasCDP() {
		return fmt.Errorf("CDP not available")
	}
	if p.coverage == nil {
		p.coverage = cdp.NewCoverage(p.cdpClient)
	}
	return p.coverage.StartJS(ctx, callCount, detailed)
}

// StartCSSCoverage begins collecting CSS coverage data.
// Requires CDP connection. Returns error if CDP is not available.
func (p *Pilot) StartCSSCoverage(ctx context.Context) error {
	if !p.HasCDP() {
		return fmt.Errorf("CDP not available")
	}
	if p.coverage == nil {
		p.coverage = cdp.NewCoverage(p.cdpClient)
	}
	return p.coverage.StartCSS(ctx)
}

// StopCoverage stops coverage collection and returns the results.
// Requires CDP connection. Returns error if CDP is not available.
func (p *Pilot) StopCoverage(ctx context.Context) (*CoverageReport, error) {
	if !p.HasCDP() {
		return nil, fmt.Errorf("CDP not available")
	}
	if p.coverage == nil {
		return nil, fmt.Errorf("coverage not started")
	}
	return p.coverage.Stop(ctx)
}

// IsCoverageRunning returns whether coverage collection is active.
func (p *Pilot) IsCoverageRunning() bool {
	if p.coverage == nil {
		return false
	}
	return p.coverage.IsRunning()
}

// ConsoleEntry is an alias for cdp.ConsoleEntry.
type ConsoleEntry = cdp.ConsoleEntry

// ExceptionDetails is an alias for cdp.ExceptionDetails.
type ExceptionDetails = cdp.ExceptionDetails

// LogEntry is an alias for cdp.LogEntry.
type LogEntry = cdp.LogEntry

// EnableConsoleDebugger starts capturing console messages with full stack traces.
// This uses CDP Runtime domain for enhanced debugging compared to BiDi console events.
// Requires CDP connection. Returns error if CDP is not available.
func (p *Pilot) EnableConsoleDebugger(ctx context.Context) error {
	if !p.HasCDP() {
		return fmt.Errorf("CDP not available")
	}
	if p.consoleDebugger == nil {
		p.consoleDebugger = cdp.NewConsoleDebugger(p.cdpClient)
	}
	return p.consoleDebugger.Enable(ctx)
}

// DisableConsoleDebugger stops capturing console messages.
// Requires CDP connection. Returns error if CDP is not available.
func (p *Pilot) DisableConsoleDebugger(ctx context.Context) error {
	if !p.HasCDP() {
		return fmt.Errorf("CDP not available")
	}
	if p.consoleDebugger == nil {
		return nil
	}
	return p.consoleDebugger.Disable(ctx)
}

// ConsoleEntries returns all captured console entries with stack traces.
// Call EnableConsoleDebugger first to start capturing.
func (p *Pilot) ConsoleEntries() []ConsoleEntry {
	if p.consoleDebugger == nil {
		return nil
	}
	return p.consoleDebugger.Entries()
}

// ConsoleExceptions returns all captured JavaScript exceptions.
// Call EnableConsoleDebugger first to start capturing.
func (p *Pilot) ConsoleExceptions() []ExceptionDetails {
	if p.consoleDebugger == nil {
		return nil
	}
	return p.consoleDebugger.Errors()
}

// BrowserLogs returns all captured browser log entries (deprecations, interventions).
// Call EnableConsoleDebugger first to start capturing.
func (p *Pilot) BrowserLogs() []LogEntry {
	if p.consoleDebugger == nil {
		return nil
	}
	return p.consoleDebugger.Logs()
}

// ClearConsoleDebugger clears all captured console entries, exceptions, and logs.
func (p *Pilot) ClearConsoleDebugger() {
	if p.consoleDebugger != nil {
		p.consoleDebugger.Clear()
	}
}

// IsConsoleDebuggerEnabled returns whether the console debugger is active.
func (p *Pilot) IsConsoleDebuggerEnabled() bool {
	if p.consoleDebugger == nil {
		return false
	}
	return p.consoleDebugger.IsEnabled()
}

// Keyboard returns the keyboard controller for this page.
func (p *Pilot) Keyboard(ctx context.Context) (*Keyboard, error) {
	if p.keyboard != nil {
		return p.keyboard, nil
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return nil, err
	}

	p.keyboard = NewKeyboard(p.client, browsingCtx)
	return p.keyboard, nil
}

// Mouse returns the mouse controller for this page.
func (p *Pilot) Mouse(ctx context.Context) (*Mouse, error) {
	if p.mouse != nil {
		return p.mouse, nil
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return nil, err
	}

	p.mouse = NewMouse(p.client, browsingCtx)
	return p.mouse, nil
}

// Touch returns the touch controller for this page.
func (p *Pilot) Touch(ctx context.Context) (*Touch, error) {
	if p.touch != nil {
		return p.touch, nil
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return nil, err
	}

	p.touch = NewTouch(p.client, browsingCtx)
	return p.touch, nil
}

// Clock returns the clock controller for this page.
func (p *Pilot) Clock(ctx context.Context) (*Clock, error) {
	if p.clock != nil {
		return p.clock, nil
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return nil, err
	}

	p.clock = NewClock(p.client, browsingCtx)
	return p.clock, nil
}

// Content returns the full HTML content of the page.
func (p *Pilot) Content(ctx context.Context) (string, error) {
	if p.closed {
		return "", ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return "", err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	result, err := p.client.Send(ctx, "vibium:page.content", params)
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
func (p *Pilot) SetContent(ctx context.Context, html string) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"html":    html,
	}

	_, err = p.client.Send(ctx, "vibium:page.setContent", params)
	return err
}

// GetViewport returns the current viewport dimensions.
func (p *Pilot) GetViewport(ctx context.Context) (Viewport, error) {
	if p.closed {
		return Viewport{}, ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return Viewport{}, err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	result, err := p.client.Send(ctx, "vibium:page.viewport", params)
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
func (p *Pilot) SetViewport(ctx context.Context, viewport Viewport) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"width":   viewport.Width,
		"height":  viewport.Height,
	}

	_, err = p.client.Send(ctx, "vibium:page.setViewport", params)
	return err
}

// GetWindow returns the browser window state.
func (p *Pilot) GetWindow(ctx context.Context) (WindowState, error) {
	if p.closed {
		return WindowState{}, ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return WindowState{}, err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	result, err := p.client.Send(ctx, "vibium:page.window", params)
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
func (p *Pilot) SetWindow(ctx context.Context, opts SetWindowOptions) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
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

	_, err = p.client.Send(ctx, "vibium:page.setWindow", params)
	return err
}

// PDF generates a PDF of the page and returns the bytes.
func (p *Pilot) PDF(ctx context.Context, opts *PDFOptions) ([]byte, error) {
	if p.closed {
		return nil, ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
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

	result, err := p.client.Send(ctx, "vibium:page.pdf", params)
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
func (p *Pilot) BringToFront(ctx context.Context) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	_, err = p.client.Send(ctx, "browsingContext.activate", params)
	return err
}

// Close closes the current page but not the browser.
func (p *Pilot) Close(ctx context.Context) error {
	if p.closed {
		return nil
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	_, err = p.client.Send(ctx, "browsingContext.close", params)
	return err
}

// Frames returns all frames on the page.
func (p *Pilot) Frames(ctx context.Context) ([]FrameInfo, error) {
	if p.closed {
		return nil, ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return nil, err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	result, err := p.client.Send(ctx, "vibium:page.frames", params)
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
func (p *Pilot) Frame(ctx context.Context, nameOrURL string) (*Pilot, error) {
	if p.closed {
		return nil, ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return nil, err
	}

	params := map[string]interface{}{
		"context":   browsingCtx,
		"nameOrURL": nameOrURL,
	}

	result, err := p.client.Send(ctx, "vibium:page.frame", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Context string `json:"context"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, err
	}

	return &Pilot{
		client:          p.client,
		clicker:         p.clicker,
		browsingContext: resp.Context,
	}, nil
}

// A11yTree returns the accessibility tree for the page.
// Options can filter the tree to only interesting nodes or specify a root element.
func (p *Pilot) A11yTree(ctx context.Context, opts *A11yTreeOptions) (interface{}, error) {
	if p.closed {
		return nil, ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
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

	result, err := p.client.Send(ctx, "vibium:page.a11yTree", params)
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
// Since Pilot represents both page and frame in this SDK, it returns itself.
// This method exists for API compatibility with other W3Pilot clients.
func (p *Pilot) MainFrame() *Pilot {
	return p
}

// EmulateMedia sets the media emulation options.
func (p *Pilot) EmulateMedia(ctx context.Context, opts EmulateMediaOptions) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
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

	_, err = p.client.Send(ctx, "vibium:page.emulateMedia", params)
	return err
}

// SetGeolocation overrides the browser's geolocation.
func (p *Pilot) SetGeolocation(ctx context.Context, coords Geolocation) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
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

	_, err = p.client.Send(ctx, "vibium:page.setGeolocation", params)
	return err
}

// AddScript adds a script that will be evaluated in the page context.
func (p *Pilot) AddScript(ctx context.Context, source string) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"source":  source,
	}

	_, err = p.client.Send(ctx, "vibium:page.addScript", params)
	return err
}

// AddStyle adds a stylesheet to the page.
func (p *Pilot) AddStyle(ctx context.Context, source string) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"source":  source,
	}

	_, err = p.client.Send(ctx, "vibium:page.addStyle", params)
	return err
}

// Expose exposes a function that can be called from JavaScript in the page.
// Note: The handler function must be registered separately.
func (p *Pilot) Expose(ctx context.Context, name string) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"name":    name,
	}

	_, err = p.client.Send(ctx, "vibium:page.expose", params)
	return err
}

// WaitForURL waits for the page URL to match the specified pattern.
func (p *Pilot) WaitForURL(ctx context.Context, pattern string, timeout time.Duration) error {
	if p.closed {
		return ErrConnectionClosed
	}

	if timeout == 0 {
		timeout = DefaultTimeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"pattern": pattern,
		"timeout": timeout.Milliseconds(),
	}

	_, err = p.client.Send(ctx, "vibium:page.waitForURL", params)
	return err
}

// WaitForLoad waits for the page to reach the specified load state.
// State can be: "load", "domcontentloaded", "networkidle".
func (p *Pilot) WaitForLoad(ctx context.Context, state string, timeout time.Duration) error {
	if p.closed {
		return ErrConnectionClosed
	}

	if timeout == 0 {
		timeout = DefaultTimeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"state":   state,
		"timeout": timeout.Milliseconds(),
	}

	_, err = p.client.Send(ctx, "vibium:page.waitForLoad", params)
	return err
}

// WaitForFunction waits for a JavaScript function to return a truthy value.
func (p *Pilot) WaitForFunction(ctx context.Context, fn string, timeout time.Duration) error {
	if p.closed {
		return ErrConnectionClosed
	}

	if timeout == 0 {
		timeout = DefaultTimeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"fn":      fn,
		"timeout": timeout.Milliseconds(),
	}

	_, err = p.client.Send(ctx, "vibium:page.waitForFunction", params)
	return err
}

// RouteHandler is called when a request matches a route pattern.
type RouteHandler func(ctx context.Context, route *Route) error

// Route registers a handler for requests matching the URL pattern.
// The pattern can be a glob pattern (e.g., "**/*.png") or regex (e.g., "/api/.*").
func (p *Pilot) Route(ctx context.Context, pattern string, handler RouteHandler) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"pattern": pattern,
	}

	_, err = p.client.Send(ctx, "vibium:network.route", params)
	return err
}

// Unroute removes a previously registered route handler.
func (p *Pilot) Unroute(ctx context.Context, pattern string) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"pattern": pattern,
	}

	_, err = p.client.Send(ctx, "vibium:network.unroute", params)
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
func (p *Pilot) MockRoute(ctx context.Context, pattern string, opts MockRouteOptions) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
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

	_, err = p.client.Send(ctx, "vibium:network.mockRoute", params)
	return err
}

// RouteInfo represents information about an active route.
type RouteInfo struct {
	Pattern     string `json:"pattern"`
	Status      int    `json:"status,omitempty"`
	ContentType string `json:"contentType,omitempty"`
}

// ListRoutes returns all active route handlers.
func (p *Pilot) ListRoutes(ctx context.Context) ([]RouteInfo, error) {
	if p.closed {
		return nil, ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return nil, err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	result, err := p.client.Send(ctx, "vibium:network.listRoutes", params)
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
// Tries BiDi first, falls back to CDP if BiDi doesn't support this command.
// For fine-grained network control (latency, bandwidth), use EmulateNetwork() instead.
func (p *Pilot) SetOffline(ctx context.Context, offline bool) error {
	if p.closed {
		return ErrConnectionClosed
	}

	// Try BiDi first
	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"offline": offline,
	}

	_, err = p.client.Send(ctx, "vibium:network.setOffline", params)
	if err == nil {
		return nil
	}

	// If BiDi doesn't support this command, fall back to CDP
	if IsUnsupportedCommand(err) {
		if !p.HasCDP() {
			return fmt.Errorf("SetOffline: BiDi not supported and CDP not available")
		}
		if offline {
			return p.EmulateNetwork(ctx, cdp.NetworkOffline)
		}
		return p.ClearNetworkEmulation(ctx)
	}

	return err
}

// SetExtraHTTPHeaders sets extra HTTP headers that will be sent with every request.
func (p *Pilot) SetExtraHTTPHeaders(ctx context.Context, headers map[string]string) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
		"headers": headers,
	}

	_, err = p.client.Send(ctx, "vibium:page.setHeaders", params)
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
type PageHandler func(*Pilot)

// PopupHandler is called when a popup window opens.
type PopupHandler func(*Pilot)

// OnRequest registers a handler for network requests.
// Note: This is a convenience method; for full control use Route().
func (p *Pilot) OnRequest(ctx context.Context, handler RequestHandler) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	// Register event handler with BiDi client
	p.client.OnEvent("vibium:network.request", func(event *BiDiEvent) {
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

	_, err = p.client.Send(ctx, "vibium:network.onRequest", params)
	return err
}

// OnResponse registers a handler for network responses.
func (p *Pilot) OnResponse(ctx context.Context, handler ResponseHandler) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	// Register event handler with BiDi client
	p.client.OnEvent("vibium:network.response", func(event *BiDiEvent) {
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

	_, err = p.client.Send(ctx, "vibium:network.onResponse", params)
	return err
}

// OnConsole registers a handler for console messages.
func (p *Pilot) OnConsole(ctx context.Context, handler ConsoleHandler) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	// Register event handler with BiDi client
	p.client.OnEvent("vibium:console.entry", func(event *BiDiEvent) {
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

	_, err = p.client.Send(ctx, "vibium:console.on", params)
	return err
}

// OnDialog registers a handler for dialogs (alert, confirm, prompt).
func (p *Pilot) OnDialog(ctx context.Context, handler DialogHandler) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	// Register event handler with BiDi client
	p.client.OnEvent("vibium:dialog.opened", func(event *BiDiEvent) {
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

	_, err = p.client.Send(ctx, "vibium:dialog.on", params)
	return err
}

// OnDownload registers a handler for downloads.
func (p *Pilot) OnDownload(ctx context.Context, handler DownloadHandler) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	// Register event handler with BiDi client
	p.client.OnEvent("vibium:download.started", func(event *BiDiEvent) {
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

	_, err = p.client.Send(ctx, "vibium:download.on", params)
	return err
}

// OnError registers a handler for JavaScript errors on the page.
func (p *Pilot) OnError(ctx context.Context, handler PageErrorHandler) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	// Register event handler with BiDi client
	p.client.OnEvent("vibium:page.error", func(event *BiDiEvent) {
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

	_, err = p.client.Send(ctx, "vibium:page.onError", params)
	return err
}

// CollectConsole enables buffered console message collection.
// Messages can be retrieved with ConsoleMessages() and cleared with ClearConsoleMessages().
func (p *Pilot) CollectConsole(ctx context.Context) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	_, err = p.client.Send(ctx, "vibium:console.collect", params)
	return err
}

// CollectErrors enables buffered page error collection.
// Errors can be retrieved with Errors() and cleared with ClearErrors().
func (p *Pilot) CollectErrors(ctx context.Context) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	_, err = p.client.Send(ctx, "vibium:page.collectErrors", params)
	return err
}

// Errors retrieves buffered page errors.
// Call CollectErrors() first to enable error collection.
func (p *Pilot) Errors(ctx context.Context) ([]PageError, error) {
	if p.closed {
		return nil, ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return nil, err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	result, err := p.client.Send(ctx, "vibium:page.errors", params)
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
func (p *Pilot) ClearErrors(ctx context.Context) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	_, err = p.client.Send(ctx, "vibium:page.clearErrors", params)
	return err
}

// OnPage registers a handler that is called when a new page is created in the browser.
// This includes pages created via NewPage(), window.open(), or clicking links with target="_blank".
func (p *Pilot) OnPage(ctx context.Context, handler PageHandler) error {
	if p.closed {
		return ErrConnectionClosed
	}

	// Register event handler with BiDi client
	p.client.OnEvent("browsingContext.contextCreated", func(event *BiDiEvent) {
		var params struct {
			Context string `json:"context"`
			URL     string `json:"url"`
			Parent  string `json:"parent,omitempty"`
		}
		if err := json.Unmarshal(event.Params, &params); err != nil {
			debugLog(ctx, "failed to unmarshal page created event", "error", err)
			return
		}

		// Create a new Pilot instance for the new page
		newPage := &Pilot{
			client:          p.client,
			clicker:         p.clicker,
			browsingContext: params.Context,
		}
		handler(newPage)
	})

	// Subscribe to context created events
	_, err := p.client.Send(ctx, "session.subscribe", map[string]interface{}{
		"events": []string{"browsingContext.contextCreated"},
	})
	return err
}

// OnPopup registers a handler that is called when a popup window is opened.
// Popups are typically created via window.open() with specific features.
func (p *Pilot) OnPopup(ctx context.Context, handler PopupHandler) error {
	if p.closed {
		return ErrConnectionClosed
	}

	// Register event handler with BiDi client
	p.client.OnEvent("browsingContext.contextCreated", func(event *BiDiEvent) {
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

		// Create a new Pilot instance for the popup
		popup := &Pilot{
			client:          p.client,
			clicker:         p.clicker,
			browsingContext: params.Context,
		}
		handler(popup)
	})

	// Subscribe to context created events
	_, err := p.client.Send(ctx, "session.subscribe", map[string]interface{}{
		"events": []string{"browsingContext.contextCreated"},
	})
	return err
}

// RemoveAllListeners removes all registered event listeners.
// This is useful for cleanup when you no longer need to receive events.
func (p *Pilot) RemoveAllListeners() {
	if p.client != nil {
		// Remove all handlers from the BiDi client
		p.client.handlerMu.Lock()
		p.client.handlers = make(map[string][]EventHandler)
		p.client.handlerMu.Unlock()
	}
}

// NewPage creates a new page in the default browser context.
func (p *Pilot) NewPage(ctx context.Context) (*Pilot, error) {
	if p.closed {
		return nil, ErrConnectionClosed
	}

	result, err := p.client.Send(ctx, "browsingContext.create", map[string]interface{}{
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

	return &Pilot{
		client:          p.client,
		clicker:         p.clicker,
		browsingContext: resp.Context,
	}, nil
}

// NewContext creates a new isolated browser context.
func (p *Pilot) NewContext(ctx context.Context) (*BrowserContext, error) {
	if p.closed {
		return nil, ErrConnectionClosed
	}

	result, err := p.client.Send(ctx, "browser.createUserContext", map[string]interface{}{})
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
		client:      p.client,
		clicker:     p.clicker,
		userContext: resp.UserContext,
	}, nil
}

// Pages returns all open pages.
func (p *Pilot) Pages(ctx context.Context) ([]*Pilot, error) {
	if p.closed {
		return nil, ErrConnectionClosed
	}

	result, err := p.client.Send(ctx, "browsingContext.getTree", map[string]interface{}{})
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

	pages := make([]*Pilot, len(tree.Contexts))
	for i, c := range tree.Contexts {
		pages[i] = &Pilot{
			client:          p.client,
			clicker:         p.clicker,
			browsingContext: c.Context,
		}
	}

	return pages, nil
}

// Context returns the browser context for this page.
// Returns nil if this is the default context.
func (p *Pilot) Context() *BrowserContext {
	// Default context doesn't have a BrowserContext wrapper
	return nil
}

// HandleDialog handles the current dialog by accepting or dismissing it.
// If accept is true, the dialog is accepted. If promptText is provided (for prompt dialogs),
// it will be entered before accepting.
func (p *Pilot) HandleDialog(ctx context.Context, accept bool, promptText string) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	// Use the correct clicker commands: vibium:dialog.accept or vibium:dialog.dismiss
	if accept {
		if promptText != "" {
			params["userText"] = promptText
		}
		_, err = p.client.Send(ctx, "vibium:dialog.accept", params)
	} else {
		_, err = p.client.Send(ctx, "vibium:dialog.dismiss", params)
	}
	return err
}

// GetDialog returns information about the current dialog, if any.
func (p *Pilot) GetDialog(ctx context.Context) (DialogInfo, error) {
	if p.closed {
		return DialogInfo{}, ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return DialogInfo{}, err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	result, err := p.client.Send(ctx, "vibium:dialog.get", params)
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
func (p *Pilot) ConsoleMessages(ctx context.Context, level string) ([]ConsoleMessage, error) {
	if p.closed {
		return nil, ErrConnectionClosed
	}

	// Try BiDi first
	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return nil, err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	if level != "" {
		params["level"] = level
	}

	result, err := p.client.Send(ctx, "vibium:console.messages", params)
	if err == nil {
		var resp struct {
			Messages []ConsoleMessage `json:"messages"`
		}
		if err := json.Unmarshal(result, &resp); err != nil {
			return nil, err
		}
		return resp.Messages, nil
	}

	// If BiDi doesn't support this command, fall back to CDP
	if IsUnsupportedCommand(err) {
		if !p.HasCDP() {
			return nil, fmt.Errorf("ConsoleMessages: BiDi not supported and CDP not available")
		}
		// Ensure console debugger is enabled
		if p.consoleDebugger == nil {
			if err := p.EnableConsoleDebugger(ctx); err != nil {
				return nil, fmt.Errorf("ConsoleMessages: failed to enable console debugger: %w", err)
			}
		}
		// Get entries from CDP and convert to ConsoleMessage format
		entries := p.ConsoleEntries()
		messages := make([]ConsoleMessage, 0, len(entries))
		for _, e := range entries {
			// Filter by level if specified
			if level != "" && string(e.Type) != level {
				continue
			}
			messages = append(messages, ConsoleMessage{
				Type: string(e.Type),
				Text: e.Text,
			})
		}
		return messages, nil
	}

	return nil, err
}

// ClearConsoleMessages clears the buffered console messages.
// Tries BiDi first, falls back to CDP if BiDi doesn't support this command.
func (p *Pilot) ClearConsoleMessages(ctx context.Context) error {
	if p.closed {
		return ErrConnectionClosed
	}

	// Try BiDi first
	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	_, err = p.client.Send(ctx, "vibium:console.clear", params)
	if err == nil {
		return nil
	}

	// If BiDi doesn't support this command, fall back to CDP
	if IsUnsupportedCommand(err) {
		if p.consoleDebugger != nil {
			p.consoleDebugger.Clear()
		}
		return nil
	}

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
func (p *Pilot) NetworkRequests(ctx context.Context, opts *NetworkRequestsOptions) ([]NetworkRequest, error) {
	if p.closed {
		return nil, ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
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

	result, err := p.client.Send(ctx, "vibium:network.requests", params)
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
func (p *Pilot) ClearNetworkRequests(ctx context.Context) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	_, err = p.client.Send(ctx, "vibium:network.clearRequests", params)
	return err
}

// ScrollOptions configures scroll behavior.
type ScrollOptions struct {
	Selector string // Optional CSS selector to scroll within
}

// Scroll scrolls the page or a specific element.
// direction can be "up", "down", "left", or "right".
// amount is the number of pixels to scroll (use 0 for full page).
func (p *Pilot) Scroll(ctx context.Context, direction string, amount int, opts *ScrollOptions) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
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

	_, err = p.client.Send(ctx, "vibium:page.scroll", params)
	return err
}

// BrowserVersion returns the browser version string.
func (p *Pilot) BrowserVersion(ctx context.Context) (string, error) {
	if p.closed {
		return "", ErrConnectionClosed
	}

	result, err := p.client.Send(ctx, "browser.getUserContexts", map[string]interface{}{})
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

// TODO: Tracing requires vibium:tracing.* commands which are not implemented in clicker.
// Uncomment when clicker adds support.
/*
// Tracing returns a tracing controller for the default browser context.
// Use this to record traces for debugging and analysis.
func (p *Pilot) Tracing() *Tracing {
	return &Tracing{
		client:      p.client,
		userContext: "", // Empty string uses the default user context
	}
}
*/

// AddInitScript adds a script that will be evaluated in every page before any page scripts.
// This is useful for mocking APIs, injecting test helpers, or setting up authentication.
func (p *Pilot) AddInitScript(ctx context.Context, script string) error {
	if p.closed {
		return ErrConnectionClosed
	}

	// Get the default user context
	userContext, err := p.getDefaultUserContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get user context: %w", err)
	}

	params := map[string]interface{}{
		"userContext": userContext,
		"script":      script,
	}

	_, err = p.client.Send(ctx, "vibium:context.addInitScript", params)
	return err
}

// getDefaultUserContext returns the default user context ID.
func (p *Pilot) getDefaultUserContext(ctx context.Context) (string, error) {
	result, err := p.client.Send(ctx, "browser.getUserContexts", map[string]interface{}{})
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
func (p *Pilot) StorageState(ctx context.Context) (*StorageState, error) {
	if p.closed {
		return nil, ErrConnectionClosed
	}

	// Get base storage state (cookies + localStorage) from context
	browserCtx, err := p.NewContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get browser context: %w", err)
	}

	state, err := browserCtx.StorageState(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage state: %w", err)
	}

	// Capture sessionStorage for the current page's origin
	currentURL, err := p.URL(ctx)
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

	result, err := p.Evaluate(ctx, sessionStorageScript)
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
func (p *Pilot) SetStorageState(ctx context.Context, state *StorageState) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browserCtx, err := p.NewContext(ctx)
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
		currentURL, _ := p.URL(ctx)
		if currentURL == "" || currentURL == "about:blank" {
			// Navigate to the origin to set storage
			if err := p.Go(ctx, origin.Origin); err != nil {
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

			if _, err := p.Evaluate(ctx, script); err != nil {
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

			if _, err := p.Evaluate(ctx, script); err != nil {
				return fmt.Errorf("failed to set sessionStorage for %s: %w", origin.Origin, err)
			}
		}
	}

	return nil
}

// ClearStorage clears all cookies, localStorage, and sessionStorage.
func (p *Pilot) ClearStorage(ctx context.Context) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browserCtx, err := p.NewContext(ctx)
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
	if _, err := p.Evaluate(ctx, script); err != nil {
		// Ignore errors (e.g., about:blank page)
		return nil
	}

	return nil
}
