package rpa

import (
	"os"
	"testing"
)

func TestResolverBasic(t *testing.T) {
	r := NewResolver(map[string]any{
		"name":  "John",
		"count": 42,
	})

	// Test Get
	val, ok := r.Get("name")
	if !ok {
		t.Fatal("Expected to find 'name'")
	}
	if val != "John" {
		t.Errorf("Expected 'John', got '%v'", val)
	}

	// Test GetString
	str, ok := r.GetString("name")
	if !ok {
		t.Fatal("Expected to find 'name'")
	}
	if str != "John" {
		t.Errorf("Expected 'John', got '%s'", str)
	}
}

func TestResolverNestedPath(t *testing.T) {
	r := NewResolver(map[string]any{
		"user": map[string]any{
			"name":  "John",
			"email": "john@example.com",
		},
	})

	val, ok := r.Get("user.name")
	if !ok {
		t.Fatal("Expected to find 'user.name'")
	}
	if val != "John" {
		t.Errorf("Expected 'John', got '%v'", val)
	}
}

func TestResolverResolve(t *testing.T) {
	r := NewResolver(map[string]any{
		"name": "John",
		"url":  "https://example.com",
	})

	// Simple substitution
	result, err := r.Resolve("Hello, ${name}!")
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	if result != "Hello, John!" {
		t.Errorf("Expected 'Hello, John!', got '%s'", result)
	}

	// Multiple substitutions
	result, err = r.Resolve("Visit ${url} as ${name}")
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	if result != "Visit https://example.com as John" {
		t.Errorf("Expected 'Visit https://example.com as John', got '%s'", result)
	}

	// Unknown variable - should leave unchanged
	result, err = r.Resolve("Hello, ${unknown}!")
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	if result != "Hello, ${unknown}!" {
		t.Errorf("Expected unchanged string, got '%s'", result)
	}
}

func TestResolverEnvVariables(t *testing.T) {
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	r := NewResolver(nil)

	result, err := r.Resolve("Value: ${env.TEST_VAR}")
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	if result != "Value: test_value" {
		t.Errorf("Expected 'Value: test_value', got '%s'", result)
	}
}

func TestEvaluatorTruthy(t *testing.T) {
	r := NewResolver(map[string]any{
		"name":     "John",
		"empty":    "",
		"zero":     0,
		"nonzero":  42,
		"trueBool": true,
	})
	e := NewEvaluator(r)

	tests := []struct {
		expr     string
		expected bool
	}{
		{"${name}", true},
		{"${empty}", false},
		{"${trueBool}", true},
		{"!${name}", false},
		{"!${empty}", true},
	}

	for _, tt := range tests {
		result, err := e.Evaluate(tt.expr)
		if err != nil {
			t.Errorf("Evaluate(%s) failed: %v", tt.expr, err)
			continue
		}
		if result != tt.expected {
			t.Errorf("Evaluate(%s) = %v, want %v", tt.expr, result, tt.expected)
		}
	}
}

func TestEvaluatorComparison(t *testing.T) {
	r := NewResolver(map[string]any{
		"count":  42,
		"name":   "John",
		"status": "active",
	})
	e := NewEvaluator(r)

	tests := []struct {
		expr     string
		expected bool
	}{
		{"${count} == 42", true},
		{"${count} != 42", false},
		{"${count} > 10", true},
		{"${count} < 10", false},
		{"${count} >= 42", true},
		{"${count} <= 42", true},
		{"${name} == 'John'", true},
		{"${name} != 'Jane'", true},
		{"${status} == 'active'", true},
	}

	for _, tt := range tests {
		result, err := e.Evaluate(tt.expr)
		if err != nil {
			t.Errorf("Evaluate(%s) failed: %v", tt.expr, err)
			continue
		}
		if result != tt.expected {
			t.Errorf("Evaluate(%s) = %v, want %v", tt.expr, result, tt.expected)
		}
	}
}
