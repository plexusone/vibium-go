package vibium

import (
	"context"
	"encoding/json"
)

// BrowserContext represents an isolated browser context (like an incognito window).
// Each context has its own cookies, localStorage, and session storage.
type BrowserContext struct {
	client      *BiDiClient
	clicker     *ClickerProcess
	userContext string
	tracing     *Tracing
}

// NewPage creates a new page in this browser context.
func (c *BrowserContext) NewPage(ctx context.Context) (*Vibe, error) {
	params := map[string]interface{}{
		"userContext": c.userContext,
	}

	result, err := c.client.Send(ctx, "browsingContext.create", params)
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
		client:          c.client,
		clicker:         c.clicker,
		browsingContext: resp.Context,
	}, nil
}

// Close closes the browser context and all pages within it.
func (c *BrowserContext) Close(ctx context.Context) error {
	params := map[string]interface{}{
		"userContext": c.userContext,
	}

	_, err := c.client.Send(ctx, "browser.removeUserContext", params)
	return err
}

// Cookies returns cookies matching the specified URLs.
// If no URLs are specified, returns all cookies for the context.
func (c *BrowserContext) Cookies(ctx context.Context, urls ...string) ([]Cookie, error) {
	params := map[string]interface{}{}

	if len(urls) > 0 {
		params["urls"] = urls
	}

	result, err := c.client.Send(ctx, "storage.getCookies", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Cookies []Cookie `json:"cookies"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, err
	}

	return resp.Cookies, nil
}

// SetCookies sets cookies.
func (c *BrowserContext) SetCookies(ctx context.Context, cookies []SetCookieParam) error {
	params := map[string]interface{}{
		"cookies": cookies,
	}

	_, err := c.client.Send(ctx, "storage.setCookie", params)
	return err
}

// ClearCookies clears all cookies.
func (c *BrowserContext) ClearCookies(ctx context.Context) error {
	params := map[string]interface{}{}

	_, err := c.client.Send(ctx, "storage.deleteCookies", params)
	return err
}

// StorageState returns the storage state including cookies and localStorage.
func (c *BrowserContext) StorageState(ctx context.Context) (*StorageState, error) {
	params := map[string]interface{}{
		"userContext": c.userContext,
	}

	result, err := c.client.Send(ctx, "vibium:context.storageState", params)
	if err != nil {
		return nil, err
	}

	var state StorageState
	if err := json.Unmarshal(result, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

// AddInitScript adds a script that will be evaluated in every page created in this context.
func (c *BrowserContext) AddInitScript(ctx context.Context, script string) error {
	params := map[string]interface{}{
		"userContext": c.userContext,
		"script":      script,
	}

	_, err := c.client.Send(ctx, "vibium:context.addInitScript", params)
	return err
}

// Tracing returns the tracing controller for this context.
func (c *BrowserContext) Tracing() *Tracing {
	if c.tracing == nil {
		c.tracing = &Tracing{
			client:      c.client,
			userContext: c.userContext,
		}
	}
	return c.tracing
}

// GrantPermissions grants the specified permissions.
func (c *BrowserContext) GrantPermissions(ctx context.Context, permissions []string, origin string) error {
	params := map[string]interface{}{
		"userContext": c.userContext,
		"permissions": permissions,
	}

	if origin != "" {
		params["origin"] = origin
	}

	_, err := c.client.Send(ctx, "vibium:context.grantPermissions", params)
	return err
}

// ClearPermissions clears all granted permissions.
func (c *BrowserContext) ClearPermissions(ctx context.Context) error {
	params := map[string]interface{}{
		"userContext": c.userContext,
	}

	_, err := c.client.Send(ctx, "vibium:context.clearPermissions", params)
	return err
}
