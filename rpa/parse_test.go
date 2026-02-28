package rpa

import (
	"strings"
	"testing"
)

func TestParseYAML(t *testing.T) {
	yaml := `
name: Test Workflow
description: A test workflow
version: "1.0"

browser:
  headless: true
  timeout: 30s

variables:
  username: testuser
  password: testpass

steps:
  - name: Navigate to page
    activity: browser.navigate
    params:
      url: https://example.com

  - name: Click button
    activity: browser.click
    params:
      selector: "#submit"
`

	wf, err := ParseBytes([]byte(yaml))
	if err != nil {
		t.Fatalf("ParseBytes failed: %v", err)
	}

	if wf.Name != "Test Workflow" {
		t.Errorf("Expected name 'Test Workflow', got '%s'", wf.Name)
	}

	if wf.Description != "A test workflow" {
		t.Errorf("Expected description 'A test workflow', got '%s'", wf.Description)
	}

	if !wf.Browser.Headless {
		t.Error("Expected headless to be true")
	}

	if len(wf.Steps) != 2 {
		t.Errorf("Expected 2 steps, got %d", len(wf.Steps))
	}

	if wf.Steps[0].Activity != "browser.navigate" {
		t.Errorf("Expected first step activity 'browser.navigate', got '%s'", wf.Steps[0].Activity)
	}
}

func TestParseJSON(t *testing.T) {
	json := `{
		"name": "JSON Workflow",
		"steps": [
			{
				"name": "Navigate",
				"activity": "browser.navigate",
				"params": {"url": "https://example.com"}
			}
		]
	}`

	wf, err := ParseBytes([]byte(json))
	if err != nil {
		t.Fatalf("ParseBytes failed: %v", err)
	}

	if wf.Name != "JSON Workflow" {
		t.Errorf("Expected name 'JSON Workflow', got '%s'", wf.Name)
	}

	if len(wf.Steps) != 1 {
		t.Errorf("Expected 1 step, got %d", len(wf.Steps))
	}
}

func TestParseValidationError(t *testing.T) {
	// Missing name
	yaml := `
steps:
  - activity: browser.navigate
`

	_, err := ParseBytes([]byte(yaml))
	if err == nil {
		t.Fatal("Expected validation error for missing name")
	}

	if !strings.Contains(err.Error(), "name") {
		t.Errorf("Expected error about name, got: %v", err)
	}
}

func TestParseEmptySteps(t *testing.T) {
	yaml := `
name: Empty Workflow
steps: []
`

	_, err := ParseBytes([]byte(yaml))
	if err == nil {
		t.Fatal("Expected validation error for empty steps")
	}

	if !strings.Contains(err.Error(), "steps") {
		t.Errorf("Expected error about steps, got: %v", err)
	}
}

func TestParseMissingActivity(t *testing.T) {
	yaml := `
name: Missing Activity
steps:
  - name: No activity step
`

	_, err := ParseBytes([]byte(yaml))
	if err == nil {
		t.Fatal("Expected validation error for missing activity")
	}

	if !strings.Contains(err.Error(), "activity") {
		t.Errorf("Expected error about activity, got: %v", err)
	}
}
