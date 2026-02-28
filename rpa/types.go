// Package rpa provides a Robotic Process Automation platform for browser automation.
package rpa

import (
	"time"
)

// ExecutionStatus represents the state of a workflow or step execution.
type ExecutionStatus string

const (
	StatusPending ExecutionStatus = "pending"
	StatusRunning ExecutionStatus = "running"
	StatusSuccess ExecutionStatus = "success"
	StatusFailure ExecutionStatus = "failure"
	StatusSkipped ExecutionStatus = "skipped"
)

// String returns the string representation of the status.
func (s ExecutionStatus) String() string {
	return string(s)
}

// IsTerminal returns true if the status is a terminal state.
func (s ExecutionStatus) IsTerminal() bool {
	return s == StatusSuccess || s == StatusFailure || s == StatusSkipped
}

// Duration represents a duration that can be unmarshaled from YAML/JSON strings.
type Duration time.Duration

// UnmarshalYAML implements yaml.Unmarshaler.
func (d *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	parsed, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = Duration(parsed)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *Duration) UnmarshalJSON(data []byte) error {
	// Remove quotes
	if len(data) >= 2 && data[0] == '"' && data[len(data)-1] == '"' {
		data = data[1 : len(data)-1]
	}
	parsed, err := time.ParseDuration(string(data))
	if err != nil {
		return err
	}
	*d = Duration(parsed)
	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (d Duration) MarshalYAML() (interface{}, error) {
	return time.Duration(d).String(), nil
}

// MarshalJSON implements json.Marshaler.
func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Duration(d).String() + `"`), nil
}

// Duration returns the underlying time.Duration.
func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}

// DefaultTimeout is the default timeout for operations.
const DefaultTimeout = 30 * time.Second

// DefaultRetryDelay is the default delay between retries.
const DefaultRetryDelay = time.Second

// DefaultMaxRetries is the default maximum number of retry attempts.
const DefaultMaxRetries = 3
