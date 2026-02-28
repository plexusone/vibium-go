package rpa

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Resolver handles variable interpolation and expression evaluation.
type Resolver struct {
	variables map[string]any
}

// NewResolver creates a new Resolver with the given variables.
func NewResolver(variables map[string]any) *Resolver {
	if variables == nil {
		variables = make(map[string]any)
	}
	return &Resolver{variables: variables}
}

// Variables returns the underlying variables map.
func (r *Resolver) Variables() map[string]any {
	return r.variables
}

// Set sets a variable value.
func (r *Resolver) Set(name string, value any) {
	r.variables[name] = value
}

// Get retrieves a variable value by path (supports dot notation).
func (r *Resolver) Get(path string) (any, bool) {
	parts := strings.Split(path, ".")
	var current any = r.variables

	for _, part := range parts {
		switch v := current.(type) {
		case map[string]any:
			val, ok := v[part]
			if !ok {
				return nil, false
			}
			current = val
		case map[string]string:
			val, ok := v[part]
			if !ok {
				return nil, false
			}
			current = val
		default:
			return nil, false
		}
	}

	return current, true
}

// GetString retrieves a variable as a string.
func (r *Resolver) GetString(path string) (string, bool) {
	val, ok := r.Get(path)
	if !ok {
		return "", false
	}
	switch v := val.(type) {
	case string:
		return v, true
	case fmt.Stringer:
		return v.String(), true
	default:
		return fmt.Sprintf("%v", v), true
	}
}

// varPattern matches ${varName} or ${varName.nested.path}
var varPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// envPattern matches ${env.VAR_NAME}
var envPattern = regexp.MustCompile(`^\s*env\.(.+)\s*$`)

// Resolve interpolates variables in a string value.
func (r *Resolver) Resolve(value string) (string, error) {
	result := varPattern.ReplaceAllStringFunc(value, func(match string) string {
		// Extract the variable path from ${path}
		path := match[2 : len(match)-1]

		// Check for environment variable reference
		if envMatch := envPattern.FindStringSubmatch(path); envMatch != nil {
			return os.Getenv(envMatch[1])
		}

		// Look up in variables
		if val, ok := r.GetString(path); ok {
			return val
		}

		// Return original if not found
		return match
	})

	return result, nil
}

// ResolveAny resolves variables in any value.
func (r *Resolver) ResolveAny(value any) (any, error) {
	switch v := value.(type) {
	case string:
		return r.Resolve(v)
	case map[string]any:
		return r.ResolveMap(v)
	case []any:
		return r.ResolveSlice(v)
	default:
		return value, nil
	}
}

// ResolveMap interpolates variables in all string values of a map.
func (r *Resolver) ResolveMap(m map[string]any) (map[string]any, error) {
	result := make(map[string]any, len(m))
	for k, v := range m {
		resolved, err := r.ResolveAny(v)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve %s: %w", k, err)
		}
		result[k] = resolved
	}
	return result, nil
}

// ResolveSlice interpolates variables in all string values of a slice.
func (r *Resolver) ResolveSlice(s []any) ([]any, error) {
	result := make([]any, len(s))
	for i, v := range s {
		resolved, err := r.ResolveAny(v)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve index %d: %w", i, err)
		}
		result[i] = resolved
	}
	return result, nil
}

// Evaluator handles condition expressions.
type Evaluator struct {
	resolver *Resolver
}

// NewEvaluator creates a new Evaluator with the given resolver.
func NewEvaluator(resolver *Resolver) *Evaluator {
	return &Evaluator{resolver: resolver}
}

// Evaluate evaluates a simple condition expression.
// Supports:
//   - ${var} - truthy check
//   - ${var} == 'value'
//   - ${var} != 'value'
//   - ${var} > number
//   - ${var} < number
//   - ${var} >= number
//   - ${var} <= number
//   - !${var} - falsy check
func (e *Evaluator) Evaluate(expr string) (bool, error) {
	expr = strings.TrimSpace(expr)

	// Handle negation
	if strings.HasPrefix(expr, "!") {
		result, err := e.Evaluate(expr[1:])
		if err != nil {
			return false, err
		}
		return !result, nil
	}

	// Try to parse comparison expressions
	operators := []string{"==", "!=", ">=", "<=", ">", "<"}
	for _, op := range operators {
		if idx := strings.Index(expr, op); idx > 0 {
			left := strings.TrimSpace(expr[:idx])
			right := strings.TrimSpace(expr[idx+len(op):])
			return e.evaluateComparison(left, op, right)
		}
	}

	// Simple truthy check - resolve the expression and check if truthy
	resolved, err := e.resolver.Resolve(expr)
	if err != nil {
		return false, err
	}
	return isTruthy(resolved), nil
}

// evaluateComparison evaluates a comparison expression.
func (e *Evaluator) evaluateComparison(left, op, right string) (bool, error) {
	// Resolve left side
	leftResolved, err := e.resolver.Resolve(left)
	if err != nil {
		return false, err
	}

	// Resolve right side
	rightResolved, err := e.resolver.Resolve(right)
	if err != nil {
		return false, err
	}

	// Remove quotes from string literals
	rightResolved = unquote(rightResolved)

	// Try numeric comparison first
	leftNum, leftIsNum := parseNumber(leftResolved)
	rightNum, rightIsNum := parseNumber(rightResolved)

	if leftIsNum && rightIsNum {
		switch op {
		case "==":
			return leftNum == rightNum, nil
		case "!=":
			return leftNum != rightNum, nil
		case ">":
			return leftNum > rightNum, nil
		case "<":
			return leftNum < rightNum, nil
		case ">=":
			return leftNum >= rightNum, nil
		case "<=":
			return leftNum <= rightNum, nil
		}
	}

	// String comparison
	switch op {
	case "==":
		return leftResolved == rightResolved, nil
	case "!=":
		return leftResolved != rightResolved, nil
	default:
		return false, fmt.Errorf("cannot compare strings with operator %s", op)
	}
}

// isTruthy checks if a value is truthy.
func isTruthy(v any) bool {
	switch val := v.(type) {
	case nil:
		return false
	case bool:
		return val
	case string:
		s := strings.ToLower(strings.TrimSpace(val))
		return s != "" && s != "false" && s != "0" && s != "null"
	case int, int64, float64:
		return val != 0
	case []any:
		return len(val) > 0
	case map[string]any:
		return len(val) > 0
	default:
		return true
	}
}

// parseNumber attempts to parse a string as a number.
func parseNumber(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f, true
	}
	return 0, false
}

// unquote removes surrounding quotes from a string.
func unquote(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
