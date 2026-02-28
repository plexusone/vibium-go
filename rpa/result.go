package rpa

import (
	"encoding/json"
	"time"
)

// WorkflowResult contains the results of a workflow execution.
type WorkflowResult struct {
	// WorkflowName is the name of the executed workflow.
	WorkflowName string `json:"workflowName"`

	// Status is the final execution status.
	Status ExecutionStatus `json:"status"`

	// StartTime is when the workflow started.
	StartTime time.Time `json:"startTime"`

	// EndTime is when the workflow finished.
	EndTime time.Time `json:"endTime"`

	// Duration is the total execution time.
	Duration time.Duration `json:"duration"`

	// Steps contains results for each executed step.
	Steps []StepResult `json:"steps"`

	// Variables contains the final state of all variables.
	Variables map[string]interface{} `json:"variables"`

	// Error contains the error message if the workflow failed.
	Error string `json:"error,omitempty"`

	// Screenshots contains any captured screenshots (base64 encoded).
	Screenshots []Screenshot `json:"screenshots,omitempty"`
}

// StepResult contains the result of a single step execution.
type StepResult struct {
	// StepID is the unique identifier of the step.
	StepID string `json:"stepId"`

	// StepName is the human-readable name of the step.
	StepName string `json:"stepName"`

	// Activity is the activity type that was executed.
	Activity string `json:"activity"`

	// Status is the execution status of the step.
	Status ExecutionStatus `json:"status"`

	// StartTime is when the step started.
	StartTime time.Time `json:"startTime"`

	// EndTime is when the step finished.
	EndTime time.Time `json:"endTime"`

	// Duration is the step execution time.
	Duration time.Duration `json:"duration"`

	// Output contains the step's output value.
	Output interface{} `json:"output,omitempty"`

	// Error contains the error message if the step failed.
	Error string `json:"error,omitempty"`

	// Screenshot contains a screenshot if captured (base64 encoded).
	Screenshot string `json:"screenshot,omitempty"`

	// Retries is the number of retry attempts made.
	Retries int `json:"retries,omitempty"`

	// Params contains the resolved parameters used for execution.
	Params map[string]interface{} `json:"params,omitempty"`
}

// Screenshot represents a captured screenshot.
type Screenshot struct {
	// StepID is the ID of the step that triggered the screenshot.
	StepID string `json:"stepId,omitempty"`

	// Timestamp is when the screenshot was captured.
	Timestamp time.Time `json:"timestamp"`

	// Data is the base64-encoded PNG image data.
	Data string `json:"data"`

	// Reason describes why the screenshot was captured.
	Reason string `json:"reason,omitempty"`
}

// NewWorkflowResult creates a new WorkflowResult with initialized fields.
func NewWorkflowResult(workflowName string) *WorkflowResult {
	return &WorkflowResult{
		WorkflowName: workflowName,
		Status:       StatusPending,
		StartTime:    time.Now(),
		Steps:        make([]StepResult, 0),
		Variables:    make(map[string]interface{}),
		Screenshots:  make([]Screenshot, 0),
	}
}

// Complete finalizes the workflow result.
func (r *WorkflowResult) Complete(status ExecutionStatus, err error) {
	r.EndTime = time.Now()
	r.Duration = r.EndTime.Sub(r.StartTime)
	r.Status = status
	if err != nil {
		r.Error = err.Error()
	}
}

// AddStep adds a step result to the workflow result.
func (r *WorkflowResult) AddStep(step StepResult) {
	r.Steps = append(r.Steps, step)
}

// AddScreenshot adds a screenshot to the workflow result.
func (r *WorkflowResult) AddScreenshot(screenshot Screenshot) {
	r.Screenshots = append(r.Screenshots, screenshot)
}

// SuccessCount returns the number of successful steps.
func (r *WorkflowResult) SuccessCount() int {
	count := 0
	for _, s := range r.Steps {
		if s.Status == StatusSuccess {
			count++
		}
	}
	return count
}

// FailureCount returns the number of failed steps.
func (r *WorkflowResult) FailureCount() int {
	count := 0
	for _, s := range r.Steps {
		if s.Status == StatusFailure {
			count++
		}
	}
	return count
}

// SkippedCount returns the number of skipped steps.
func (r *WorkflowResult) SkippedCount() int {
	count := 0
	for _, s := range r.Steps {
		if s.Status == StatusSkipped {
			count++
		}
	}
	return count
}

// TotalSteps returns the total number of executed steps.
func (r *WorkflowResult) TotalSteps() int {
	return len(r.Steps)
}

// IsSuccess returns true if the workflow completed successfully.
func (r *WorkflowResult) IsSuccess() bool {
	return r.Status == StatusSuccess
}

// JSON returns the result as formatted JSON.
func (r *WorkflowResult) JSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

// NewStepResult creates a new StepResult with initialized fields.
func NewStepResult(step *Step) *StepResult {
	return &StepResult{
		StepID:    step.GetID(),
		StepName:  step.Name,
		Activity:  step.Activity,
		Status:    StatusPending,
		StartTime: time.Now(),
	}
}

// Complete finalizes the step result.
func (r *StepResult) Complete(status ExecutionStatus, output interface{}, err error) {
	r.EndTime = time.Now()
	r.Duration = r.EndTime.Sub(r.StartTime)
	r.Status = status
	r.Output = output
	if err != nil {
		r.Error = err.Error()
	}
}

// MarkRunning marks the step as running.
func (r *StepResult) MarkRunning() {
	r.Status = StatusRunning
}

// MarkSkipped marks the step as skipped with a reason.
func (r *StepResult) MarkSkipped(reason string) {
	r.Status = StatusSkipped
	r.Error = reason
	r.EndTime = time.Now()
	r.Duration = r.EndTime.Sub(r.StartTime)
}
