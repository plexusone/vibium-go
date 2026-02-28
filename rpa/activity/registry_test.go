package activity

import (
	"context"
	"testing"
)

// mockActivity is a simple activity for testing.
type mockActivity struct {
	name   string
	result any
	err    error
}

func (a *mockActivity) Name() string { return a.name }
func (a *mockActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	return a.result, a.err
}

func TestRegistry(t *testing.T) {
	r := NewRegistry()

	// Test empty registry
	if r.Count() != 0 {
		t.Errorf("Expected empty registry, got %d activities", r.Count())
	}

	// Register an activity
	mock := &mockActivity{name: "test.mock", result: "test result"}
	r.Register(mock)

	if r.Count() != 1 {
		t.Errorf("Expected 1 activity, got %d", r.Count())
	}

	// Get the activity
	a, ok := r.Get("test.mock")
	if !ok {
		t.Fatal("Expected to find 'test.mock' activity")
	}
	if a.Name() != "test.mock" {
		t.Errorf("Expected name 'test.mock', got '%s'", a.Name())
	}

	// Get non-existent activity
	_, ok = r.Get("nonexistent")
	if ok {
		t.Error("Expected not to find 'nonexistent' activity")
	}
}

func TestRegistryList(t *testing.T) {
	r := NewRegistry()

	r.Register(&mockActivity{name: "browser.click"})
	r.Register(&mockActivity{name: "browser.navigate"})
	r.Register(&mockActivity{name: "file.read"})

	names := r.List()
	if len(names) != 3 {
		t.Errorf("Expected 3 activities, got %d", len(names))
	}

	// Check sorted order
	expected := []string{"browser.click", "browser.navigate", "file.read"}
	for i, name := range names {
		if name != expected[i] {
			t.Errorf("Expected names[%d] = '%s', got '%s'", i, expected[i], name)
		}
	}
}

func TestRegistryListByCategory(t *testing.T) {
	r := NewRegistry()

	r.Register(&mockActivity{name: "browser.click"})
	r.Register(&mockActivity{name: "browser.navigate"})
	r.Register(&mockActivity{name: "file.read"})
	r.Register(&mockActivity{name: "file.write"})

	byCategory := r.ListByCategory()

	if len(byCategory["browser"]) != 2 {
		t.Errorf("Expected 2 browser activities, got %d", len(byCategory["browser"]))
	}
	if len(byCategory["file"]) != 2 {
		t.Errorf("Expected 2 file activities, got %d", len(byCategory["file"]))
	}
}

func TestDefaultRegistry(t *testing.T) {
	// Verify default registry has activities registered
	names := DefaultRegistry.List()
	if len(names) == 0 {
		t.Fatal("Expected default registry to have registered activities")
	}

	// Check for some expected activities
	expectedActivities := []string{
		"browser.navigate",
		"browser.click",
		"browser.fill",
		"element.getText",
		"file.read",
		"file.write",
		"http.get",
		"util.log",
	}

	for _, name := range expectedActivities {
		_, ok := DefaultRegistry.Get(name)
		if !ok {
			t.Errorf("Expected to find activity '%s' in default registry", name)
		}
	}
}

func TestGetHelpers(t *testing.T) {
	params := map[string]any{
		"stringVal": "hello",
		"intVal":    42,
		"floatVal":  3.14,
		"boolVal":   true,
		"slice":     []any{"a", "b", "c"},
	}

	if v := GetString(params, "stringVal"); v != "hello" {
		t.Errorf("GetString expected 'hello', got '%s'", v)
	}

	if v := GetInt(params, "intVal"); v != 42 {
		t.Errorf("GetInt expected 42, got %d", v)
	}

	if v := GetFloat(params, "floatVal"); v != 3.14 {
		t.Errorf("GetFloat expected 3.14, got %f", v)
	}

	if v := GetBool(params, "boolVal"); !v {
		t.Error("GetBool expected true, got false")
	}

	if v := GetStringSlice(params, "slice"); len(v) != 3 {
		t.Errorf("GetStringSlice expected 3 items, got %d", len(v))
	}

	// Test defaults
	if v := GetStringDefault(params, "missing", "default"); v != "default" {
		t.Errorf("GetStringDefault expected 'default', got '%s'", v)
	}

	if v := GetIntDefault(params, "missing", 99); v != 99 {
		t.Errorf("GetIntDefault expected 99, got %d", v)
	}
}
