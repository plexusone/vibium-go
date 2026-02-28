package activity

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// FileReadActivity reads a file's contents.
type FileReadActivity struct{}

func (a *FileReadActivity) Name() string { return "file.read" }

func (a *FileReadActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	path := GetString(params, "path")
	if path == "" {
		return nil, fmt.Errorf("path parameter is required")
	}

	// Resolve relative paths
	if !filepath.IsAbs(path) {
		path = filepath.Join(env.WorkDir, path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read failed: %w", err)
	}

	// Parse as JSON if requested
	if GetBool(params, "json") {
		var result any
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("JSON parse failed: %w", err)
		}
		return result, nil
	}

	return string(data), nil
}

// FileWriteActivity writes content to a file.
type FileWriteActivity struct{}

func (a *FileWriteActivity) Name() string { return "file.write" }

func (a *FileWriteActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	path := GetString(params, "path")
	if path == "" {
		return nil, fmt.Errorf("path parameter is required")
	}

	// Resolve relative paths
	if !filepath.IsAbs(path) {
		path = filepath.Join(env.WorkDir, path)
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Get content to write
	content := params["content"]
	var data []byte

	format := GetStringDefault(params, "format", "text")
	switch format {
	case "json":
		var err error
		data, err = json.MarshalIndent(content, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("JSON encoding failed: %w", err)
		}
	default:
		switch v := content.(type) {
		case string:
			data = []byte(v)
		case []byte:
			data = v
		default:
			// Convert to string
			data = []byte(fmt.Sprintf("%v", v))
		}
	}

	// Write file
	mode := os.FileMode(GetIntDefault(params, "mode", 0644))
	if GetBool(params, "append") {
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, mode)
		if err != nil {
			return nil, fmt.Errorf("open failed: %w", err)
		}
		defer f.Close()
		if _, err := f.Write(data); err != nil {
			return nil, fmt.Errorf("write failed: %w", err)
		}
	} else {
		if err := os.WriteFile(path, data, mode); err != nil {
			return nil, fmt.Errorf("write failed: %w", err)
		}
	}

	return map[string]any{
		"path":  path,
		"bytes": len(data),
	}, nil
}

// FileExistsActivity checks if a file exists.
type FileExistsActivity struct{}

func (a *FileExistsActivity) Name() string { return "file.exists" }

func (a *FileExistsActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	path := GetString(params, "path")
	if path == "" {
		return nil, fmt.Errorf("path parameter is required")
	}

	// Resolve relative paths
	if !filepath.IsAbs(path) {
		path = filepath.Join(env.WorkDir, path)
	}

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return nil, fmt.Errorf("stat failed: %w", err)
	}

	return map[string]any{
		"exists": true,
		"isDir":  info.IsDir(),
		"size":   info.Size(),
	}, nil
}

// FileDeleteActivity deletes a file.
type FileDeleteActivity struct{}

func (a *FileDeleteActivity) Name() string { return "file.delete" }

func (a *FileDeleteActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	path := GetString(params, "path")
	if path == "" {
		return nil, fmt.Errorf("path parameter is required")
	}

	// Resolve relative paths
	if !filepath.IsAbs(path) {
		path = filepath.Join(env.WorkDir, path)
	}

	// Check if exists first
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if GetBool(params, "ignoreNotExist") {
			return nil, nil
		}
		return nil, fmt.Errorf("file not found: %s", path)
	}

	if err := os.Remove(path); err != nil {
		return nil, fmt.Errorf("delete failed: %w", err)
	}

	return nil, nil
}
