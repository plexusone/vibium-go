package vibium

import (
	"context"
)

// Route represents an intercepted network request.
type Route struct {
	client    *BiDiClient
	context   string
	intercept string
	Request   *Request
}

// Request represents a network request.
type Request struct {
	URL                 string            `json:"url"`
	Method              string            `json:"method"`
	Headers             map[string]string `json:"headers"`
	PostData            string            `json:"postData,omitempty"`
	ResourceType        string            `json:"resourceType"`
	IsNavigationRequest bool              `json:"isNavigationRequest"`
}

// Response represents a network response.
type Response struct {
	URL        string            `json:"url"`
	Status     int               `json:"status"`
	StatusText string            `json:"statusText"`
	Headers    map[string]string `json:"headers"`
	Body       []byte            `json:"-"`
}

// FulfillOptions configures how to fulfill a route.
type FulfillOptions struct {
	Status      int
	Headers     map[string]string
	ContentType string
	Body        []byte
	Path        string // Path to file to serve
}

// ContinueOptions configures how to continue a route.
type ContinueOptions struct {
	URL      string
	Method   string
	Headers  map[string]string
	PostData string
}

// Fulfill fulfills the route with the given response.
func (r *Route) Fulfill(ctx context.Context, opts FulfillOptions) error {
	params := map[string]interface{}{
		"context":   r.context,
		"intercept": r.intercept,
	}

	if opts.Status != 0 {
		params["status"] = opts.Status
	}
	if opts.Headers != nil {
		params["headers"] = opts.Headers
	}
	if opts.ContentType != "" {
		params["contentType"] = opts.ContentType
	}
	if opts.Body != nil {
		params["body"] = opts.Body
	}
	if opts.Path != "" {
		params["path"] = opts.Path
	}

	_, err := r.client.Send(ctx, "vibium:network.fulfill", params)
	return err
}

// Continue continues the route with optional modifications.
func (r *Route) Continue(ctx context.Context, opts *ContinueOptions) error {
	params := map[string]interface{}{
		"context":   r.context,
		"intercept": r.intercept,
	}

	if opts != nil {
		if opts.URL != "" {
			params["url"] = opts.URL
		}
		if opts.Method != "" {
			params["method"] = opts.Method
		}
		if opts.Headers != nil {
			params["headers"] = opts.Headers
		}
		if opts.PostData != "" {
			params["postData"] = opts.PostData
		}
	}

	_, err := r.client.Send(ctx, "vibium:network.continue", params)
	return err
}

// Abort aborts the route.
func (r *Route) Abort(ctx context.Context) error {
	params := map[string]interface{}{
		"context":   r.context,
		"intercept": r.intercept,
	}

	_, err := r.client.Send(ctx, "vibium:network.abort", params)
	return err
}

// ConsoleMessage represents a console message from the browser.
type ConsoleMessage struct {
	Type string   `json:"type"`
	Text string   `json:"text"`
	Args []string `json:"args,omitempty"`
	URL  string   `json:"url,omitempty"`
	Line int      `json:"line,omitempty"`
}
