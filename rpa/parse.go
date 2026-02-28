package rpa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ParseFile parses a workflow from a file, auto-detecting format.
func ParseFile(path string) (*Workflow, error) {
	ext := filepath.Ext(path)
	switch ext {
	case ".yaml", ".yml":
		return parseYAMLFile(path)
	case ".json":
		return parseJSONFile(path)
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}
}

// ParseBytes parses a workflow from bytes, auto-detecting format.
// If the data starts with '{' or '[', it's treated as JSON, otherwise YAML.
func ParseBytes(data []byte) (*Workflow, error) {
	// Skip whitespace to detect format
	for i := 0; i < len(data); i++ {
		c := data[i]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			continue
		}
		if c == '{' || c == '[' {
			return parseJSONBytes(data)
		}
		break
	}
	return parseYAMLBytes(data)
}

// ParseReader parses a workflow from a reader.
func ParseReader(r io.Reader, format string) (*Workflow, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read: %w", err)
	}

	switch format {
	case "yaml", "yml":
		return parseYAMLBytes(data)
	case "json":
		return parseJSONBytes(data)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// MustParseFile parses a workflow from a file and panics on error.
func MustParseFile(path string) *Workflow {
	wf, err := ParseFile(path)
	if err != nil {
		panic(fmt.Sprintf("failed to parse workflow %s: %v", path, err))
	}
	return wf
}

// parseYAMLFile parses a workflow from a YAML file.
func parseYAMLFile(path string) (*Workflow, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return parseYAMLBytes(data)
}

// parseYAMLBytes parses a workflow from YAML bytes.
func parseYAMLBytes(data []byte) (*Workflow, error) {
	var wf Workflow
	decoder := yaml.NewDecoder(bytes.NewReader(data))

	if err := decoder.Decode(&wf); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Validate
	if errors := validateWorkflow(&wf); len(errors) > 0 {
		return nil, fmt.Errorf("validation errors: %v", errors)
	}

	return &wf, nil
}

// parseJSONFile parses a workflow from a JSON file.
func parseJSONFile(path string) (*Workflow, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return parseJSONBytes(data)
}

// parseJSONBytes parses a workflow from JSON bytes.
func parseJSONBytes(data []byte) (*Workflow, error) {
	var wf Workflow
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields() // Strict mode - fail on unknown fields

	if err := decoder.Decode(&wf); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Validate
	if errors := validateWorkflow(&wf); len(errors) > 0 {
		return nil, fmt.Errorf("validation errors: %v", errors)
	}

	return &wf, nil
}

// ParserValidationError represents a validation error during parsing.
type ParserValidationError struct {
	Path    string
	Field   string
	Message string
}

func (e ParserValidationError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("%s.%s: %s", e.Path, e.Field, e.Message)
	}
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// validateWorkflow validates a workflow definition.
func validateWorkflow(wf *Workflow) []ParserValidationError {
	var errors []ParserValidationError

	if wf.Name == "" {
		errors = append(errors, ParserValidationError{
			Field:   "name",
			Message: "workflow name is required",
		})
	}

	if len(wf.Steps) == 0 {
		errors = append(errors, ParserValidationError{
			Field:   "steps",
			Message: "workflow must have at least one step",
		})
	}

	for i, step := range wf.Steps {
		stepErrors := validateParserStep(&step, fmt.Sprintf("steps[%d]", i))
		errors = append(errors, stepErrors...)
	}

	return errors
}

func validateParserStep(step *Step, path string) []ParserValidationError {
	var errors []ParserValidationError

	if step.Activity == "" {
		errors = append(errors, ParserValidationError{
			Path:    path,
			Field:   "activity",
			Message: "activity is required",
		})
	}

	if step.ForEach != nil {
		if step.ForEach.Items == "" {
			errors = append(errors, ParserValidationError{
				Path:    path + ".forEach",
				Field:   "items",
				Message: "items is required for forEach",
			})
		}
		if step.ForEach.Variable == "" {
			errors = append(errors, ParserValidationError{
				Path:    path + ".forEach",
				Field:   "as",
				Message: "variable name (as) is required for forEach",
			})
		}
		for i, nested := range step.ForEach.Steps {
			nestedErrors := validateParserStep(&nested, fmt.Sprintf("%s.forEach.steps[%d]", path, i))
			errors = append(errors, nestedErrors...)
		}
	}

	for i, nested := range step.Steps {
		nestedErrors := validateParserStep(&nested, fmt.Sprintf("%s.steps[%d]", path, i))
		errors = append(errors, nestedErrors...)
	}

	return errors
}
