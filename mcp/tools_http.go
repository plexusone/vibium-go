package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// HTTPRequestInput defines the input for the http_request tool.
type HTTPRequestInput struct {
	URL           string            `json:"url" jsonschema:"URL to fetch (absolute or relative to current page),required"`
	Method        string            `json:"method" jsonschema:"HTTP method (GET POST PUT DELETE PATCH HEAD OPTIONS),default=GET"`
	Headers       map[string]string `json:"headers" jsonschema:"Additional HTTP headers to include"`
	ContentType   string            `json:"content_type" jsonschema:"Content-Type header value (convenience for common types)"`
	Body          string            `json:"body" jsonschema:"Request body (for POST PUT PATCH)"`
	MaxBodyLength int               `json:"max_body_length" jsonschema:"Maximum response body length in characters (0=unlimited),default=8192"`
}

// HTTPRequestOutput defines the output for the http_request tool.
type HTTPRequestOutput struct {
	Status     int               `json:"status"`
	StatusText string            `json:"status_text"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
	Truncated  bool              `json:"truncated"`
	URL        string            `json:"url"`
}

func (s *Server) handleHTTPRequest(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input HTTPRequestInput,
) (*mcp.CallToolResult, HTTPRequestOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, HTTPRequestOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	// Set defaults
	method := input.Method
	if method == "" {
		method = "GET"
	}
	method = strings.ToUpper(method)

	maxBodyLength := input.MaxBodyLength
	if maxBodyLength <= 0 {
		maxBodyLength = 8192
	}

	// Build headers object for JavaScript
	headers := make(map[string]string)
	for k, v := range input.Headers {
		headers[k] = v
	}
	if input.ContentType != "" {
		headers["Content-Type"] = input.ContentType
	}

	headersJSON, err := json.Marshal(headers)
	if err != nil {
		return nil, HTTPRequestOutput{}, fmt.Errorf("failed to marshal headers: %w", err)
	}

	// Build the fetch script
	// This script:
	// 1. Makes the fetch request with credentials included (uses browser's cookies)
	// 2. Collects response status, headers, and body
	// 3. Truncates body if needed
	// 4. Returns a structured result
	script := fmt.Sprintf(`(async () => {
		const url = %q;
		const method = %q;
		const headers = %s;
		const body = %q;
		const maxLen = %d;

		const options = {
			method: method,
			credentials: 'include',
			headers: headers
		};

		if (body && (method === 'POST' || method === 'PUT' || method === 'PATCH')) {
			options.body = body;
		}

		const response = await fetch(url, options);

		// Collect response headers
		const responseHeaders = {};
		response.headers.forEach((value, key) => {
			responseHeaders[key] = value;
		});

		// Get response body with truncation
		let responseBody = await response.text();
		let truncated = false;
		if (maxLen > 0 && responseBody.length > maxLen) {
			responseBody = responseBody.substring(0, maxLen);
			truncated = true;
		}

		return {
			status: response.status,
			statusText: response.statusText,
			headers: responseHeaders,
			body: responseBody,
			truncated: truncated,
			url: response.url
		};
	})()`, input.URL, method, string(headersJSON), input.Body, maxBodyLength)

	result, err := pilot.Evaluate(ctx, script)
	if err != nil {
		return nil, HTTPRequestOutput{}, fmt.Errorf("fetch failed: %w", err)
	}

	if result == nil {
		return nil, HTTPRequestOutput{}, fmt.Errorf("fetch returned nil result")
	}

	// Parse the result
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		return nil, HTTPRequestOutput{}, fmt.Errorf("unexpected result type: %T", result)
	}

	output := HTTPRequestOutput{}

	if status, ok := resultMap["status"].(float64); ok {
		output.Status = int(status)
	}
	if statusText, ok := resultMap["statusText"].(string); ok {
		output.StatusText = statusText
	}
	if headers, ok := resultMap["headers"].(map[string]interface{}); ok {
		output.Headers = make(map[string]string)
		for k, v := range headers {
			if s, ok := v.(string); ok {
				output.Headers[k] = s
			}
		}
	}
	if body, ok := resultMap["body"].(string); ok {
		output.Body = body
	}
	if truncated, ok := resultMap["truncated"].(bool); ok {
		output.Truncated = truncated
	}
	if url, ok := resultMap["url"].(string); ok {
		output.URL = url
	}

	return nil, output, nil
}
