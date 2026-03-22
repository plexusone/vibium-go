package mcp

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	vibium "github.com/plexusone/vibium-go"
)

// GetContent tool

type GetContentInput struct{}

type GetContentOutput struct {
	Content string `json:"content"`
}

func (s *Server) handleGetContent(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetContentInput,
) (*mcp.CallToolResult, GetContentOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, GetContentOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	content, err := vibe.Content(ctx)
	if err != nil {
		return nil, GetContentOutput{}, fmt.Errorf("get content failed: %w", err)
	}

	return nil, GetContentOutput{Content: content}, nil
}

// SetContent tool

type SetContentInput struct {
	HTML string `json:"html" jsonschema:"HTML content to set,required"`
}

type SetContentOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleSetContent(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SetContentInput,
) (*mcp.CallToolResult, SetContentOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, SetContentOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	err = vibe.SetContent(ctx, input.HTML)
	if err != nil {
		return nil, SetContentOutput{}, fmt.Errorf("set content failed: %w", err)
	}

	return nil, SetContentOutput{Message: "Content set"}, nil
}

// GetViewport tool

type GetViewportInput struct{}

type GetViewportOutput struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func (s *Server) handleGetViewport(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetViewportInput,
) (*mcp.CallToolResult, GetViewportOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, GetViewportOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	vp, err := vibe.GetViewport(ctx)
	if err != nil {
		return nil, GetViewportOutput{}, fmt.Errorf("get viewport failed: %w", err)
	}

	return nil, GetViewportOutput{Width: vp.Width, Height: vp.Height}, nil
}

// SetViewport tool

type SetViewportInput struct {
	Width  int `json:"width" jsonschema:"Viewport width,required"`
	Height int `json:"height" jsonschema:"Viewport height,required"`
}

type SetViewportOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleSetViewport(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SetViewportInput,
) (*mcp.CallToolResult, SetViewportOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, SetViewportOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	err = vibe.SetViewport(ctx, vibium.Viewport{Width: input.Width, Height: input.Height})
	if err != nil {
		return nil, SetViewportOutput{}, fmt.Errorf("set viewport failed: %w", err)
	}

	return nil, SetViewportOutput{Message: fmt.Sprintf("Viewport set to %dx%d", input.Width, input.Height)}, nil
}

// PDF tool

type PDFInput struct {
	Scale           float64 `json:"scale" jsonschema:"Scale of the PDF (default: 1)"`
	PrintBackground bool    `json:"print_background" jsonschema:"Print background graphics"`
	Landscape       bool    `json:"landscape" jsonschema:"Landscape orientation"`
	Format          string  `json:"format" jsonschema:"Paper format (Letter Legal A4 etc)"`
}

type PDFOutput struct {
	Data string `json:"data"`
}

func (s *Server) handlePDF(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input PDFInput,
) (*mcp.CallToolResult, PDFOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, PDFOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	opts := &vibium.PDFOptions{
		Scale:           input.Scale,
		PrintBackground: input.PrintBackground,
		Landscape:       input.Landscape,
		Format:          input.Format,
	}

	data, err := vibe.PDF(ctx, opts)
	if err != nil {
		return nil, PDFOutput{}, fmt.Errorf("pdf generation failed: %w", err)
	}

	return nil, PDFOutput{Data: base64.StdEncoding.EncodeToString(data)}, nil
}

// BringToFront tool

type BringToFrontInput struct{}

type BringToFrontOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleBringToFront(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input BringToFrontInput,
) (*mcp.CallToolResult, BringToFrontOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, BringToFrontOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	err = vibe.BringToFront(ctx)
	if err != nil {
		return nil, BringToFrontOutput{}, fmt.Errorf("bring to front failed: %w", err)
	}

	return nil, BringToFrontOutput{Message: "Page brought to front"}, nil
}

// ClosePage tool

type ClosePageInput struct{}

type ClosePageOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleClosePage(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ClosePageInput,
) (*mcp.CallToolResult, ClosePageOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, ClosePageOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	err = vibe.Close(ctx)
	if err != nil {
		return nil, ClosePageOutput{}, fmt.Errorf("close page failed: %w", err)
	}

	return nil, ClosePageOutput{Message: "Page closed"}, nil
}

// GetFrames tool

type GetFramesInput struct{}

type GetFramesOutput struct {
	Frames []FrameInfoOutput `json:"frames"`
}

type FrameInfoOutput struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

func (s *Server) handleGetFrames(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetFramesInput,
) (*mcp.CallToolResult, GetFramesOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, GetFramesOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	frames, err := vibe.Frames(ctx)
	if err != nil {
		return nil, GetFramesOutput{}, fmt.Errorf("get frames failed: %w", err)
	}

	output := make([]FrameInfoOutput, len(frames))
	for i, f := range frames {
		output[i] = FrameInfoOutput{URL: f.URL, Name: f.Name}
	}

	return nil, GetFramesOutput{Frames: output}, nil
}

// SelectFrame tool - switch to a frame by name or URL

type SelectFrameInput struct {
	NameOrURL string `json:"name_or_url" jsonschema:"Frame name or URL pattern to match,required"`
}

type SelectFrameOutput struct {
	Message string `json:"message"`
	URL     string `json:"url"`
	Name    string `json:"name"`
}

func (s *Server) handleSelectFrame(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SelectFrameInput,
) (*mcp.CallToolResult, SelectFrameOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, SelectFrameOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	frame, err := vibe.Frame(ctx, input.NameOrURL)
	if err != nil {
		return nil, SelectFrameOutput{}, fmt.Errorf("frame not found: %w", err)
	}

	// Update the session to use this frame context
	s.session.SetVibe(frame)

	// Get frame info
	url, _ := frame.URL(ctx)
	title, _ := frame.Title(ctx)

	return nil, SelectFrameOutput{
		Message: fmt.Sprintf("Switched to frame: %s", input.NameOrURL),
		URL:     url,
		Name:    title,
	}, nil
}

// SelectMainFrame tool - switch back to the main frame

type SelectMainFrameInput struct{}

type SelectMainFrameOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleSelectMainFrame(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SelectMainFrameInput,
) (*mcp.CallToolResult, SelectMainFrameOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, SelectMainFrameOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	// MainFrame returns the main frame (self in our case)
	mainFrame := vibe.MainFrame()
	s.session.SetVibe(mainFrame)

	return nil, SelectMainFrameOutput{
		Message: "Switched to main frame",
	}, nil
}

// EmulateMedia tool

type EmulateMediaInput struct {
	Media         string `json:"media,omitempty" jsonschema:"Media type: screen or print. Empty to reset."`
	ColorScheme   string `json:"color_scheme,omitempty" jsonschema:"Color scheme preference: light dark or no-preference. For testing dark mode support.,enum=light,enum=dark,enum=no-preference"`
	ReducedMotion string `json:"reduced_motion,omitempty" jsonschema:"Reduced motion preference: reduce or no-preference. For testing animation accessibility.,enum=reduce,enum=no-preference"`
	ForcedColors  string `json:"forced_colors,omitempty" jsonschema:"Forced colors mode: active or none. For testing Windows High Contrast Mode.,enum=active,enum=none"`
	Contrast      string `json:"contrast,omitempty" jsonschema:"Contrast preference: more less or no-preference. For testing low vision accessibility.,enum=more,enum=less,enum=no-preference"`
}

type EmulateMediaOutput struct {
	Message  string   `json:"message"`
	Settings []string `json:"settings"`
}

func (s *Server) handleEmulateMedia(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input EmulateMediaInput,
) (*mcp.CallToolResult, EmulateMediaOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, EmulateMediaOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	err = vibe.EmulateMedia(ctx, vibium.EmulateMediaOptions{
		Media:         input.Media,
		ColorScheme:   input.ColorScheme,
		ReducedMotion: input.ReducedMotion,
		ForcedColors:  input.ForcedColors,
		Contrast:      input.Contrast,
	})
	if err != nil {
		return nil, EmulateMediaOutput{}, fmt.Errorf("emulate media failed: %w", err)
	}

	// Build list of applied settings
	var settings []string
	if input.Media != "" {
		settings = append(settings, "media="+input.Media)
	}
	if input.ColorScheme != "" {
		settings = append(settings, "colorScheme="+input.ColorScheme)
	}
	if input.ReducedMotion != "" {
		settings = append(settings, "reducedMotion="+input.ReducedMotion)
	}
	if input.ForcedColors != "" {
		settings = append(settings, "forcedColors="+input.ForcedColors)
	}
	if input.Contrast != "" {
		settings = append(settings, "contrast="+input.Contrast)
	}

	return nil, EmulateMediaOutput{
		Message:  "Media emulation set",
		Settings: settings,
	}, nil
}

// SetGeolocation tool

type SetGeolocationInput struct {
	Latitude  float64 `json:"latitude" jsonschema:"Latitude,required"`
	Longitude float64 `json:"longitude" jsonschema:"Longitude,required"`
	Accuracy  float64 `json:"accuracy" jsonschema:"Accuracy in meters"`
}

type SetGeolocationOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleSetGeolocation(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SetGeolocationInput,
) (*mcp.CallToolResult, SetGeolocationOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, SetGeolocationOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	err = vibe.SetGeolocation(ctx, vibium.Geolocation{
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
		Accuracy:  input.Accuracy,
	})
	if err != nil {
		return nil, SetGeolocationOutput{}, fmt.Errorf("set geolocation failed: %w", err)
	}

	return nil, SetGeolocationOutput{Message: fmt.Sprintf("Geolocation set to %f, %f", input.Latitude, input.Longitude)}, nil
}

// AddScript tool

type AddScriptInput struct {
	Source string `json:"source" jsonschema:"JavaScript source to inject,required"`
}

type AddScriptOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleAddScript(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input AddScriptInput,
) (*mcp.CallToolResult, AddScriptOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, AddScriptOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	err = vibe.AddScript(ctx, input.Source)
	if err != nil {
		return nil, AddScriptOutput{}, fmt.Errorf("add script failed: %w", err)
	}

	return nil, AddScriptOutput{Message: "Script added"}, nil
}

// AddStyle tool

type AddStyleInput struct {
	Source string `json:"source" jsonschema:"CSS source to inject,required"`
}

type AddStyleOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleAddStyle(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input AddStyleInput,
) (*mcp.CallToolResult, AddStyleOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, AddStyleOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	err = vibe.AddStyle(ctx, input.Source)
	if err != nil {
		return nil, AddStyleOutput{}, fmt.Errorf("add style failed: %w", err)
	}

	return nil, AddStyleOutput{Message: "Style added"}, nil
}

// WaitForURL tool

type WaitForURLInput struct {
	Pattern   string `json:"pattern" jsonschema:"URL pattern to wait for,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"Timeout in milliseconds (default: 30000)"`
}

type WaitForURLOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleWaitForURL(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input WaitForURLInput,
) (*mcp.CallToolResult, WaitForURLOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, WaitForURLOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 30000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	err = vibe.WaitForURL(ctx, input.Pattern, timeout)
	if err != nil {
		return nil, WaitForURLOutput{}, fmt.Errorf("wait for URL failed: %w", err)
	}

	return nil, WaitForURLOutput{Message: fmt.Sprintf("URL matched pattern: %s", input.Pattern)}, nil
}

// WaitForLoad tool

type WaitForLoadInput struct {
	State     string `json:"state" jsonschema:"Load state: load domcontentloaded networkidle,required,enum=load,enum=domcontentloaded,enum=networkidle"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"Timeout in milliseconds (default: 30000)"`
}

type WaitForLoadOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleWaitForLoad(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input WaitForLoadInput,
) (*mcp.CallToolResult, WaitForLoadOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, WaitForLoadOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 30000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	err = vibe.WaitForLoad(ctx, input.State, timeout)
	if err != nil {
		return nil, WaitForLoadOutput{}, fmt.Errorf("wait for load failed: %w", err)
	}

	return nil, WaitForLoadOutput{Message: fmt.Sprintf("Page reached state: %s", input.State)}, nil
}

// WaitForFunction tool

type WaitForFunctionInput struct {
	Function  string `json:"function" jsonschema:"JavaScript function that returns truthy value,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"Timeout in milliseconds (default: 30000)"`
}

type WaitForFunctionOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleWaitForFunction(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input WaitForFunctionInput,
) (*mcp.CallToolResult, WaitForFunctionOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, WaitForFunctionOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 30000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	err = vibe.WaitForFunction(ctx, input.Function, timeout)
	if err != nil {
		return nil, WaitForFunctionOutput{}, fmt.Errorf("wait for function failed: %w", err)
	}

	return nil, WaitForFunctionOutput{Message: "Function returned truthy value"}, nil
}

// WaitForText tool - wait for text to appear on the page

type WaitForTextInput struct {
	Text      string `json:"text" jsonschema:"Text to wait for,required"`
	Selector  string `json:"selector,omitempty" jsonschema:"Optional selector to scope the search"`
	TimeoutMS int    `json:"timeout,omitempty" jsonschema:"Timeout in milliseconds (default 30000)"`
}

type WaitForTextOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleWaitForText(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input WaitForTextInput,
) (*mcp.CallToolResult, WaitForTextOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, WaitForTextOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 30000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	// Build JavaScript function to check for text
	var checkScript string
	if input.Selector != "" {
		// Search within a specific element
		checkScript = fmt.Sprintf(`() => {
			const el = document.querySelector(%q);
			return el && el.textContent.includes(%q);
		}`, input.Selector, input.Text)
	} else {
		// Search entire document body
		checkScript = fmt.Sprintf(`() => document.body.textContent.includes(%q)`, input.Text)
	}

	err = vibe.WaitForFunction(ctx, checkScript, timeout)
	if err != nil {
		return nil, WaitForTextOutput{}, fmt.Errorf("wait for text failed: %w", err)
	}

	return nil, WaitForTextOutput{Message: fmt.Sprintf("Text found: %s", input.Text)}, nil
}

// AccessibilitySnapshot tool - get accessibility tree snapshot

type AccessibilitySnapshotInput struct {
	InterestingOnly *bool  `json:"interesting_only,omitempty" jsonschema:"Only include interesting nodes with semantic meaning (default true)"`
	Root            string `json:"root,omitempty" jsonschema:"CSS selector for root element to snapshot"`
}

type AccessibilitySnapshotOutput struct {
	Snapshot interface{} `json:"snapshot"`
}

func (s *Server) handleAccessibilitySnapshot(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input AccessibilitySnapshotInput,
) (*mcp.CallToolResult, AccessibilitySnapshotOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, AccessibilitySnapshotOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	opts := &vibium.A11yTreeOptions{
		InterestingOnly: input.InterestingOnly,
		Root:            input.Root,
	}

	tree, err := vibe.A11yTree(ctx, opts)
	if err != nil {
		return nil, AccessibilitySnapshotOutput{}, fmt.Errorf("accessibility snapshot failed: %w", err)
	}

	return nil, AccessibilitySnapshotOutput{Snapshot: tree}, nil
}

// Back tool

type BackInput struct{}

type BackOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleBack(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input BackInput,
) (*mcp.CallToolResult, BackOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, BackOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	err = vibe.Back(ctx)
	if err != nil {
		return nil, BackOutput{}, fmt.Errorf("back failed: %w", err)
	}

	return nil, BackOutput{Message: "Navigated back"}, nil
}

// Forward tool

type ForwardInput struct{}

type ForwardOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleForward(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ForwardInput,
) (*mcp.CallToolResult, ForwardOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, ForwardOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	err = vibe.Forward(ctx)
	if err != nil {
		return nil, ForwardOutput{}, fmt.Errorf("forward failed: %w", err)
	}

	return nil, ForwardOutput{Message: "Navigated forward"}, nil
}

// Reload tool

type ReloadInput struct{}

type ReloadOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleReload(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ReloadInput,
) (*mcp.CallToolResult, ReloadOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, ReloadOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	err = vibe.Reload(ctx)
	if err != nil {
		return nil, ReloadOutput{}, fmt.Errorf("reload failed: %w", err)
	}

	return nil, ReloadOutput{Message: "Page reloaded"}, nil
}

// Scroll tool

type ScrollInput struct {
	Direction string `json:"direction" jsonschema:"Scroll direction: up down left right,required,enum=up,enum=down,enum=left,enum=right"`
	Amount    int    `json:"amount" jsonschema:"Amount to scroll in pixels (0 for full page scroll)"`
	Selector  string `json:"selector" jsonschema:"Optional CSS selector to scroll within a specific element"`
}

type ScrollOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleScroll(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ScrollInput,
) (*mcp.CallToolResult, ScrollOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, ScrollOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	var opts *vibium.ScrollOptions
	if input.Selector != "" {
		opts = &vibium.ScrollOptions{Selector: input.Selector}
	}

	err = vibe.Scroll(ctx, input.Direction, input.Amount, opts)
	if err != nil {
		return nil, ScrollOutput{}, fmt.Errorf("scroll failed: %w", err)
	}

	msg := fmt.Sprintf("Scrolled %s", input.Direction)
	if input.Amount > 0 {
		msg = fmt.Sprintf("Scrolled %s %d pixels", input.Direction, input.Amount)
	}
	if input.Selector != "" {
		msg += fmt.Sprintf(" in %s", input.Selector)
	}

	return nil, ScrollOutput{Message: msg}, nil
}
