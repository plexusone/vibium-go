package activity

import (
	"context"
	"fmt"
	"time"
)

// LogActivity logs a message.
type LogActivity struct{}

func (a *LogActivity) Name() string { return "util.log" }

func (a *LogActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	message := GetString(params, "message")
	level := GetStringDefault(params, "level", "info")

	switch level {
	case "debug":
		env.Logger.Debug(message)
	case "info":
		env.Logger.Info(message)
	case "warn":
		env.Logger.Warn(message)
	case "error":
		env.Logger.Error(message)
	default:
		env.Logger.Info(message)
	}

	return nil, nil
}

// WaitActivity waits for a specified duration.
type WaitActivity struct{}

func (a *WaitActivity) Name() string { return "util.wait" }

func (a *WaitActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	// Support both "duration" (string like "5s") and "ms" (integer milliseconds)
	if duration := GetString(params, "duration"); duration != "" {
		d, err := time.ParseDuration(duration)
		if err != nil {
			return nil, fmt.Errorf("invalid duration: %w", err)
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(d):
			return nil, nil
		}
	}

	if ms := GetInt(params, "ms"); ms > 0 {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Duration(ms) * time.Millisecond):
			return nil, nil
		}
	}

	return nil, fmt.Errorf("duration or ms parameter is required")
}

// AssertActivity asserts a condition is true.
type AssertActivity struct{}

func (a *AssertActivity) Name() string { return "util.assert" }

func (a *AssertActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	condition := params["condition"]
	message := GetStringDefault(params, "message", "assertion failed")

	// Check the condition
	var isTrue bool
	switch v := condition.(type) {
	case bool:
		isTrue = v
	case string:
		isTrue = v != ""
	case nil:
		isTrue = false
	default:
		isTrue = true // Non-nil values are truthy
	}

	if !isTrue {
		return nil, fmt.Errorf("%s", message)
	}

	return nil, nil
}

// SetVariableActivity sets a variable value.
type SetVariableActivity struct{}

func (a *SetVariableActivity) Name() string { return "util.setVariable" }

func (a *SetVariableActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	name := GetString(params, "name")
	if name == "" {
		return nil, fmt.Errorf("name parameter is required")
	}

	value := params["value"]
	env.Variables[name] = value

	return value, nil
}
