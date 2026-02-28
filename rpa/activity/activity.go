// Package activity provides the activity system for RPA workflow execution.
package activity

import (
	"context"
	"log/slog"

	vibium "github.com/plexusone/vibium-go"
)

// Activity represents an executable RPA action.
type Activity interface {
	// Name returns the activity's unique identifier (e.g., "browser.navigate").
	Name() string

	// Execute runs the activity with the given parameters and environment.
	Execute(ctx context.Context, params map[string]any, env *Environment) (any, error)
}

// Environment provides execution context to activities.
type Environment struct {
	// Vibe is the browser automation interface.
	Vibe *vibium.Vibe

	// Variables contains workflow and step variables.
	Variables map[string]any

	// WorkDir is the working directory for file operations.
	WorkDir string

	// Logger is the structured logger for activity output.
	Logger *slog.Logger

	// Headless indicates if the browser is running in headless mode.
	Headless bool
}

// NewEnvironment creates a new Environment with initialized fields.
func NewEnvironment(vibe *vibium.Vibe, workDir string, logger *slog.Logger) *Environment {
	return &Environment{
		Vibe:      vibe,
		Variables: make(map[string]any),
		WorkDir:   workDir,
		Logger:    logger,
	}
}

// GetString retrieves a string parameter from the params map.
func GetString(params map[string]any, key string) string {
	if v, ok := params[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetStringDefault retrieves a string parameter with a default value.
func GetStringDefault(params map[string]any, key, defaultValue string) string {
	if v, ok := params[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultValue
}

// GetBool retrieves a boolean parameter from the params map.
func GetBool(params map[string]any, key string) bool {
	if v, ok := params[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// GetInt retrieves an integer parameter from the params map.
func GetInt(params map[string]any, key string) int {
	if v, ok := params[key]; ok {
		switch n := v.(type) {
		case int:
			return n
		case int64:
			return int(n)
		case float64:
			return int(n)
		}
	}
	return 0
}

// GetIntDefault retrieves an integer parameter with a default value.
func GetIntDefault(params map[string]any, key string, defaultValue int) int {
	if v, ok := params[key]; ok {
		switch n := v.(type) {
		case int:
			return n
		case int64:
			return int(n)
		case float64:
			return int(n)
		}
	}
	return defaultValue
}

// GetFloat retrieves a float parameter from the params map.
func GetFloat(params map[string]any, key string) float64 {
	if v, ok := params[key]; ok {
		switch n := v.(type) {
		case float64:
			return n
		case float32:
			return float64(n)
		case int:
			return float64(n)
		case int64:
			return float64(n)
		}
	}
	return 0
}

// GetStringSlice retrieves a string slice parameter from the params map.
func GetStringSlice(params map[string]any, key string) []string {
	if v, ok := params[key]; ok {
		if slice, ok := v.([]string); ok {
			return slice
		}
		// Handle []interface{} from YAML/JSON parsing
		if slice, ok := v.([]any); ok {
			result := make([]string, 0, len(slice))
			for _, item := range slice {
				if s, ok := item.(string); ok {
					result = append(result, s)
				}
			}
			return result
		}
	}
	return nil
}

// GetMap retrieves a map parameter from the params map.
func GetMap(params map[string]any, key string) map[string]any {
	if v, ok := params[key]; ok {
		if m, ok := v.(map[string]any); ok {
			return m
		}
	}
	return nil
}
