package activity

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// HTTPGetActivity performs an HTTP GET request.
type HTTPGetActivity struct{}

func (a *HTTPGetActivity) Name() string { return "http.get" }

func (a *HTTPGetActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	url := GetString(params, "url")
	if url == "" {
		return nil, fmt.Errorf("url parameter is required")
	}

	timeout := time.Duration(GetIntDefault(params, "timeout", 30000)) * time.Millisecond
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	if headers := GetMap(params, "headers"); headers != nil {
		for k, v := range headers {
			if s, ok := v.(string); ok {
				req.Header.Set(k, s)
			}
		}
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	result := map[string]any{
		"status":     resp.StatusCode,
		"statusText": resp.Status,
		"headers":    headerToMap(resp.Header),
	}

	// Parse as JSON if requested or if Content-Type is JSON
	if GetBool(params, "json") || isJSONContentType(resp.Header.Get("Content-Type")) {
		var jsonBody any
		if err := json.Unmarshal(body, &jsonBody); err == nil {
			result["body"] = jsonBody
			return result, nil
		}
	}

	result["body"] = string(body)
	return result, nil
}

// HTTPPostActivity performs an HTTP POST request.
type HTTPPostActivity struct{}

func (a *HTTPPostActivity) Name() string { return "http.post" }

func (a *HTTPPostActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	url := GetString(params, "url")
	if url == "" {
		return nil, fmt.Errorf("url parameter is required")
	}

	timeout := time.Duration(GetIntDefault(params, "timeout", 30000)) * time.Millisecond
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Prepare body
	var bodyReader io.Reader
	contentType := GetStringDefault(params, "contentType", "application/json")

	if body := params["body"]; body != nil {
		switch v := body.(type) {
		case string:
			bodyReader = bytes.NewBufferString(v)
		case []byte:
			bodyReader = bytes.NewBuffer(v)
		default:
			// Encode as JSON
			data, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("failed to encode body: %w", err)
			}
			bodyReader = bytes.NewBuffer(data)
			contentType = "application/json"
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)

	// Add headers
	if headers := GetMap(params, "headers"); headers != nil {
		for k, v := range headers {
			if s, ok := v.(string); ok {
				req.Header.Set(k, s)
			}
		}
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	result := map[string]any{
		"status":     resp.StatusCode,
		"statusText": resp.Status,
		"headers":    headerToMap(resp.Header),
	}

	// Parse as JSON if requested or if Content-Type is JSON
	if GetBool(params, "json") || isJSONContentType(resp.Header.Get("Content-Type")) {
		var jsonBody any
		if err := json.Unmarshal(respBody, &jsonBody); err == nil {
			result["body"] = jsonBody
			return result, nil
		}
	}

	result["body"] = string(respBody)
	return result, nil
}

// HTTPDownloadActivity downloads a file.
type HTTPDownloadActivity struct{}

func (a *HTTPDownloadActivity) Name() string { return "http.download" }

func (a *HTTPDownloadActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	url := GetString(params, "url")
	if url == "" {
		return nil, fmt.Errorf("url parameter is required")
	}

	path := GetString(params, "path")
	if path == "" {
		return nil, fmt.Errorf("path parameter is required")
	}

	// Resolve relative paths
	if !filepath.IsAbs(path) {
		path = filepath.Join(env.WorkDir, path)
	}

	timeout := time.Duration(GetIntDefault(params, "timeout", 60000)) * time.Millisecond
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	if headers := GetMap(params, "headers"); headers != nil {
		for k, v := range headers {
			if s, ok := v.(string); ok {
				req.Header.Set(k, s)
			}
		}
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Create file
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	// Copy response body to file
	n, err := io.Copy(f, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	return map[string]any{
		"path":  path,
		"bytes": n,
	}, nil
}

// headerToMap converts http.Header to a simple map.
func headerToMap(h http.Header) map[string]string {
	m := make(map[string]string)
	for k, v := range h {
		if len(v) > 0 {
			m[k] = v[0]
		}
	}
	return m
}

// isJSONContentType checks if the content type is JSON.
func isJSONContentType(ct string) bool {
	return ct == "application/json" ||
		ct == "application/json; charset=utf-8" ||
		ct == "text/json"
}
