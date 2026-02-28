package rpa

// Workflow represents a complete automation workflow.
type Workflow struct {
	// Name is the human-readable name of the workflow.
	Name string `yaml:"name" json:"name"`

	// Description provides additional context about the workflow.
	Description string `yaml:"description,omitempty" json:"description,omitempty"`

	// Version is the semantic version of the workflow definition.
	Version string `yaml:"version,omitempty" json:"version,omitempty"`

	// Browser contains browser-specific configuration.
	Browser BrowserConfig `yaml:"browser,omitempty" json:"browser,omitempty"`

	// Variables defines workflow-level variables with default values.
	// These can be overridden at runtime.
	Variables map[string]string `yaml:"variables,omitempty" json:"variables,omitempty"`

	// Steps is the ordered list of steps to execute.
	Steps []Step `yaml:"steps" json:"steps"`

	// OnError defines error handling behavior for the workflow.
	OnError *ErrorHandler `yaml:"onError,omitempty" json:"onError,omitempty"`
}

// BrowserConfig contains browser-specific configuration options.
type BrowserConfig struct {
	// Headless runs the browser without a visible UI.
	Headless bool `yaml:"headless" json:"headless"`

	// Timeout is the default timeout for browser operations.
	Timeout Duration `yaml:"timeout,omitempty" json:"timeout,omitempty"`

	// Viewport sets the browser viewport dimensions.
	Viewport *ViewportConfig `yaml:"viewport,omitempty" json:"viewport,omitempty"`

	// UserAgent overrides the browser's user agent string.
	UserAgent string `yaml:"userAgent,omitempty" json:"userAgent,omitempty"`

	// IgnoreHTTPSErrors ignores HTTPS certificate errors.
	IgnoreHTTPSErrors bool `yaml:"ignoreHTTPSErrors,omitempty" json:"ignoreHTTPSErrors,omitempty"`
}

// ViewportConfig defines browser viewport dimensions.
type ViewportConfig struct {
	Width  int `yaml:"width" json:"width"`
	Height int `yaml:"height" json:"height"`
}

// Step represents a single step in a workflow.
type Step struct {
	// ID is a unique identifier for the step (optional).
	ID string `yaml:"id,omitempty" json:"id,omitempty"`

	// Name is a human-readable name for the step.
	Name string `yaml:"name,omitempty" json:"name,omitempty"`

	// Activity is the activity type to execute (e.g., "browser.navigate").
	Activity string `yaml:"activity" json:"activity"`

	// Params contains the parameters for the activity.
	Params map[string]interface{} `yaml:"params,omitempty" json:"params,omitempty"`

	// Condition is an expression that must evaluate to true for the step to execute.
	Condition string `yaml:"if,omitempty" json:"if,omitempty"`

	// ForEach enables iteration over a collection.
	ForEach *ForEachConfig `yaml:"forEach,omitempty" json:"forEach,omitempty"`

	// Store specifies a variable name to store the step's output.
	Store string `yaml:"store,omitempty" json:"store,omitempty"`

	// ContinueOnError allows the workflow to continue if this step fails.
	ContinueOnError bool `yaml:"continueOnError,omitempty" json:"continueOnError,omitempty"`

	// Retry configures automatic retry behavior.
	Retry *RetryConfig `yaml:"retry,omitempty" json:"retry,omitempty"`

	// Timeout overrides the default timeout for this step.
	Timeout Duration `yaml:"timeout,omitempty" json:"timeout,omitempty"`

	// Steps contains nested steps (for control flow activities).
	Steps []Step `yaml:"steps,omitempty" json:"steps,omitempty"`
}

// ForEachConfig configures iteration over a collection.
type ForEachConfig struct {
	// Items is the variable name or expression containing the items to iterate.
	Items string `yaml:"items" json:"items"`

	// Variable is the name of the loop variable (available as ${variable}).
	Variable string `yaml:"as" json:"as"`

	// Steps are the steps to execute for each item.
	Steps []Step `yaml:"steps" json:"steps"`
}

// RetryConfig configures automatic retry behavior.
type RetryConfig struct {
	// MaxAttempts is the maximum number of retry attempts.
	MaxAttempts int `yaml:"maxAttempts" json:"maxAttempts"`

	// Delay is the delay between retry attempts.
	Delay Duration `yaml:"delay" json:"delay"`

	// BackoffMultiplier multiplies the delay after each retry (default: 1.0).
	BackoffMultiplier float64 `yaml:"backoffMultiplier,omitempty" json:"backoffMultiplier,omitempty"`
}

// ErrorHandler configures error handling behavior.
type ErrorHandler struct {
	// Screenshot captures a screenshot when an error occurs.
	Screenshot bool `yaml:"screenshot" json:"screenshot"`

	// Steps are optional steps to execute when an error occurs.
	Steps []Step `yaml:"steps,omitempty" json:"steps,omitempty"`
}

// GetID returns the step's ID, generating one from the name if not set.
func (s *Step) GetID() string {
	if s.ID != "" {
		return s.ID
	}
	if s.Name != "" {
		return s.Name
	}
	return s.Activity
}

// GetTimeout returns the step's timeout, or the default if not set.
func (s *Step) GetTimeout(defaultTimeout Duration) Duration {
	if s.Timeout > 0 {
		return s.Timeout
	}
	return defaultTimeout
}

// HasCondition returns true if the step has a conditional expression.
func (s *Step) HasCondition() bool {
	return s.Condition != ""
}

// HasForEach returns true if the step is a forEach loop.
func (s *Step) HasForEach() bool {
	return s.ForEach != nil
}

// HasRetry returns true if the step has retry configuration.
func (s *Step) HasRetry() bool {
	return s.Retry != nil && s.Retry.MaxAttempts > 0
}
