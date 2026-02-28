package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	vibium "github.com/plexusone/vibium-go"
)

// NewPage tool

type NewPageInput struct{}

type NewPageOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleNewPage(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input NewPageInput,
) (*mcp.CallToolResult, NewPageOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, NewPageOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	_, err = vibe.NewPage(ctx)
	if err != nil {
		return nil, NewPageOutput{}, fmt.Errorf("new page failed: %w", err)
	}

	return nil, NewPageOutput{Message: "New page created"}, nil
}

// GetPages tool

type GetPagesInput struct{}

type GetPagesOutput struct {
	Count int `json:"count"`
}

func (s *Server) handleGetPages(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetPagesInput,
) (*mcp.CallToolResult, GetPagesOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, GetPagesOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	pages, err := vibe.Pages(ctx)
	if err != nil {
		return nil, GetPagesOutput{}, fmt.Errorf("get pages failed: %w", err)
	}

	return nil, GetPagesOutput{Count: len(pages)}, nil
}

// GetCookies tool

type GetCookiesInput struct {
	URLs []string `json:"urls" jsonschema:"description=URLs to get cookies for (optional)"`
}

type GetCookiesOutput struct {
	Cookies []CookieOutput `json:"cookies"`
}

type CookieOutput struct {
	Name     string  `json:"name"`
	Value    string  `json:"value"`
	Domain   string  `json:"domain"`
	Path     string  `json:"path"`
	Expires  float64 `json:"expires"`
	HTTPOnly bool    `json:"httpOnly"`
	Secure   bool    `json:"secure"`
	SameSite string  `json:"sameSite"`
}

func (s *Server) handleGetCookies(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetCookiesInput,
) (*mcp.CallToolResult, GetCookiesOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, GetCookiesOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	browserCtx, err := vibe.NewContext(ctx)
	if err != nil {
		return nil, GetCookiesOutput{}, fmt.Errorf("context not available: %w", err)
	}

	cookies, err := browserCtx.Cookies(ctx, input.URLs...)
	if err != nil {
		return nil, GetCookiesOutput{}, fmt.Errorf("get cookies failed: %w", err)
	}

	output := make([]CookieOutput, len(cookies))
	for i, c := range cookies {
		output[i] = CookieOutput{
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

	return nil, GetCookiesOutput{Cookies: output}, nil
}

// SetCookies tool

type SetCookiesInput struct {
	Cookies []SetCookieInput `json:"cookies" jsonschema:"description=Cookies to set,required"`
}

type SetCookieInput struct {
	Name     string  `json:"name"`
	Value    string  `json:"value"`
	URL      string  `json:"url,omitempty"`
	Domain   string  `json:"domain,omitempty"`
	Path     string  `json:"path,omitempty"`
	Expires  float64 `json:"expires,omitempty"`
	HTTPOnly bool    `json:"httpOnly,omitempty"`
	Secure   bool    `json:"secure,omitempty"`
	SameSite string  `json:"sameSite,omitempty"`
}

type SetCookiesOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleSetCookies(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SetCookiesInput,
) (*mcp.CallToolResult, SetCookiesOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, SetCookiesOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	browserCtx, err := vibe.NewContext(ctx)
	if err != nil {
		return nil, SetCookiesOutput{}, fmt.Errorf("context not available: %w", err)
	}

	cookies := make([]vibium.SetCookieParam, len(input.Cookies))
	for i, c := range input.Cookies {
		cookies[i] = vibium.SetCookieParam{
			Name:     c.Name,
			Value:    c.Value,
			URL:      c.URL,
			Domain:   c.Domain,
			Path:     c.Path,
			Expires:  c.Expires,
			HTTPOnly: c.HTTPOnly,
			Secure:   c.Secure,
			SameSite: c.SameSite,
		}
	}

	err = browserCtx.SetCookies(ctx, cookies)
	if err != nil {
		return nil, SetCookiesOutput{}, fmt.Errorf("set cookies failed: %w", err)
	}

	return nil, SetCookiesOutput{Message: fmt.Sprintf("Set %d cookies", len(input.Cookies))}, nil
}

// ClearCookies tool

type ClearCookiesInput struct{}

type ClearCookiesOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleClearCookies(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ClearCookiesInput,
) (*mcp.CallToolResult, ClearCookiesOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, ClearCookiesOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	browserCtx, err := vibe.NewContext(ctx)
	if err != nil {
		return nil, ClearCookiesOutput{}, fmt.Errorf("context not available: %w", err)
	}

	err = browserCtx.ClearCookies(ctx)
	if err != nil {
		return nil, ClearCookiesOutput{}, fmt.Errorf("clear cookies failed: %w", err)
	}

	return nil, ClearCookiesOutput{Message: "Cookies cleared"}, nil
}

// GetStorageState tool

type GetStorageStateInput struct{}

type GetStorageStateOutput struct {
	State string `json:"state"`
}

func (s *Server) handleGetStorageState(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetStorageStateInput,
) (*mcp.CallToolResult, GetStorageStateOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, GetStorageStateOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	browserCtx, err := vibe.NewContext(ctx)
	if err != nil {
		return nil, GetStorageStateOutput{}, fmt.Errorf("context not available: %w", err)
	}

	state, err := browserCtx.StorageState(ctx)
	if err != nil {
		return nil, GetStorageStateOutput{}, fmt.Errorf("get storage state failed: %w", err)
	}

	jsonBytes, err := json.Marshal(state)
	if err != nil {
		return nil, GetStorageStateOutput{}, fmt.Errorf("json marshal failed: %w", err)
	}

	return nil, GetStorageStateOutput{State: string(jsonBytes)}, nil
}
