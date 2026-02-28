package mcp

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	vibium "github.com/agentplexus/vibium-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
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
	HTML string `json:"html" jsonschema:"description=HTML content to set,required"`
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
	Width  int `json:"width" jsonschema:"description=Viewport width,required"`
	Height int `json:"height" jsonschema:"description=Viewport height,required"`
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
	Scale           float64 `json:"scale" jsonschema:"description=Scale of the PDF (default: 1)"`
	PrintBackground bool    `json:"print_background" jsonschema:"description=Print background graphics"`
	Landscape       bool    `json:"landscape" jsonschema:"description=Landscape orientation"`
	Format          string  `json:"format" jsonschema:"description=Paper format (Letter Legal A4 etc)"`
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

// EmulateMedia tool

type EmulateMediaInput struct {
	Media         string `json:"media" jsonschema:"description=Media type: screen print or empty"`
	ColorScheme   string `json:"color_scheme" jsonschema:"description=Color scheme: light dark no-preference or empty"`
	ReducedMotion string `json:"reduced_motion" jsonschema:"description=Reduced motion: reduce no-preference or empty"`
}

type EmulateMediaOutput struct {
	Message string `json:"message"`
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
	})
	if err != nil {
		return nil, EmulateMediaOutput{}, fmt.Errorf("emulate media failed: %w", err)
	}

	return nil, EmulateMediaOutput{Message: "Media emulation set"}, nil
}

// SetGeolocation tool

type SetGeolocationInput struct {
	Latitude  float64 `json:"latitude" jsonschema:"description=Latitude,required"`
	Longitude float64 `json:"longitude" jsonschema:"description=Longitude,required"`
	Accuracy  float64 `json:"accuracy" jsonschema:"description=Accuracy in meters"`
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
	Source string `json:"source" jsonschema:"description=JavaScript source to inject,required"`
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
	Source string `json:"source" jsonschema:"description=CSS source to inject,required"`
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
	Pattern   string `json:"pattern" jsonschema:"description=URL pattern to wait for,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 30000)"`
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
	State     string `json:"state" jsonschema:"description=Load state: load domcontentloaded networkidle,required,enum=load,enum=domcontentloaded,enum=networkidle"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 30000)"`
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
	Function  string `json:"function" jsonschema:"description=JavaScript function that returns truthy value,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Timeout in milliseconds (default: 30000)"`
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
