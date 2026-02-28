package rpa

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"time"

	vibium "github.com/agentplexus/vibium-go"
	"github.com/agentplexus/vibium-go/rpa/activity"
)

// ExecutorConfig configures the workflow executor.
type ExecutorConfig struct {
	// Headless runs the browser in headless mode.
	Headless bool

	// DefaultTimeout is the default timeout for operations.
	DefaultTimeout time.Duration

	// WorkDir is the working directory for file operations.
	WorkDir string

	// Variables contains runtime variable overrides.
	Variables map[string]string

	// DryRun parses and validates without executing.
	DryRun bool

	// Logger is the structured logger.
	Logger *slog.Logger

	// OnStepStart is called when a step starts.
	OnStepStart func(step *Step)

	// OnStepComplete is called when a step completes.
	OnStepComplete func(step *Step, result *StepResult)
}

// Executor runs RPA workflows.
type Executor struct {
	config   ExecutorConfig
	registry *activity.Registry
	logger   *slog.Logger
}

// NewExecutor creates a new workflow executor.
func NewExecutor(config ExecutorConfig) *Executor {
	if config.DefaultTimeout == 0 {
		config.DefaultTimeout = DefaultTimeout
	}
	if config.WorkDir == "" {
		config.WorkDir, _ = os.Getwd()
	}
	if config.Logger == nil {
		config.Logger = slog.Default()
	}

	return &Executor{
		config:   config,
		registry: activity.DefaultRegistry,
		logger:   config.Logger,
	}
}

// RunFile executes a workflow from a file.
func (e *Executor) RunFile(ctx context.Context, path string) (*WorkflowResult, error) {
	wf, err := ParseFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse workflow: %w", err)
	}
	return e.RunWorkflow(ctx, wf)
}

// RunWorkflow executes a parsed workflow.
func (e *Executor) RunWorkflow(ctx context.Context, wf *Workflow) (*WorkflowResult, error) {
	result := NewWorkflowResult(wf.Name)
	result.Status = StatusRunning

	// Initialize resolver with workflow variables
	variables := make(map[string]any)
	for k, v := range wf.Variables {
		variables[k] = v
	}
	// Override with config variables
	for k, v := range e.config.Variables {
		variables[k] = v
	}

	resolver := NewResolver(variables)

	// Dry run - just validate
	if e.config.DryRun {
		errors := e.Validate(ctx, wf)
		if len(errors) > 0 {
			result.Complete(StatusFailure, fmt.Errorf("validation failed: %v", errors))
			return result, nil
		}
		result.Complete(StatusSuccess, nil)
		return result, nil
	}

	// Determine headless mode
	headless := e.config.Headless || wf.Browser.Headless

	// Launch browser
	e.logger.Info("launching browser", "headless", headless)
	launchOpts := &vibium.LaunchOptions{Headless: headless}
	vibe, err := vibium.Browser.Launch(ctx, launchOpts)
	if err != nil {
		result.Complete(StatusFailure, fmt.Errorf("failed to launch browser: %w", err))
		return result, nil
	}
	defer vibe.Quit(ctx)

	// Create execution environment
	env := activity.NewEnvironment(vibe, e.config.WorkDir, e.logger)
	env.Variables = resolver.Variables()
	env.Headless = headless

	// Execute steps
	if err := e.runSteps(ctx, wf.Steps, env, resolver, result); err != nil {
		// Handle error
		if wf.OnError != nil {
			e.handleError(ctx, wf.OnError, env, resolver, result, err)
		}
		result.Complete(StatusFailure, err)
		return result, nil
	}

	result.Variables = resolver.Variables()
	result.Complete(StatusSuccess, nil)
	return result, nil
}

// runSteps executes a list of steps.
func (e *Executor) runSteps(ctx context.Context, steps []Step, env *activity.Environment, resolver *Resolver, result *WorkflowResult) error {
	evaluator := NewEvaluator(resolver)

	for i := range steps {
		step := &steps[i]

		// Check for context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Check condition
		if step.HasCondition() {
			ok, err := evaluator.Evaluate(step.Condition)
			if err != nil {
				return fmt.Errorf("condition evaluation failed for step %s: %w", step.GetID(), err)
			}
			if !ok {
				stepResult := NewStepResult(step)
				stepResult.MarkSkipped("condition not met")
				result.AddStep(*stepResult)
				continue
			}
		}

		// Handle forEach
		if step.HasForEach() {
			if err := e.runForEach(ctx, step, env, resolver, result); err != nil {
				if !step.ContinueOnError {
					return err
				}
			}
			continue
		}

		// Execute step with retries
		stepResult, err := e.executeStepWithRetry(ctx, step, env, resolver)
		result.AddStep(*stepResult)

		// Store output
		if step.Store != "" && stepResult.Output != nil {
			resolver.Set(step.Store, stepResult.Output)
			env.Variables[step.Store] = stepResult.Output
		}

		// Handle errors
		if err != nil && !step.ContinueOnError {
			return err
		}
	}

	return nil
}

// runForEach executes a forEach loop.
func (e *Executor) runForEach(ctx context.Context, step *Step, env *activity.Environment, resolver *Resolver, result *WorkflowResult) error {
	forEach := step.ForEach

	// Get items to iterate
	itemsValue, ok := resolver.Get(forEach.Items)
	if !ok {
		return fmt.Errorf("forEach items not found: %s", forEach.Items)
	}

	items, ok := itemsValue.([]any)
	if !ok {
		return fmt.Errorf("forEach items must be an array")
	}

	// Iterate
	for i, item := range items {
		// Set loop variable
		resolver.Set(forEach.Variable, item)
		resolver.Set(forEach.Variable+"_index", i)
		env.Variables[forEach.Variable] = item
		env.Variables[forEach.Variable+"_index"] = i

		// Execute steps
		if err := e.runSteps(ctx, forEach.Steps, env, resolver, result); err != nil {
			if !step.ContinueOnError {
				return err
			}
		}
	}

	return nil
}

// executeStepWithRetry executes a step with retry logic.
func (e *Executor) executeStepWithRetry(ctx context.Context, step *Step, env *activity.Environment, resolver *Resolver) (*StepResult, error) {
	maxAttempts := 1
	delay := DefaultRetryDelay

	if step.HasRetry() {
		maxAttempts = step.Retry.MaxAttempts
		if step.Retry.Delay > 0 {
			delay = step.Retry.Delay.Duration()
		}
	}

	var lastErr error
	stepResult := NewStepResult(step)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		stepResult.Retries = attempt - 1
		stepResult.MarkRunning()

		// Notify step start
		if e.config.OnStepStart != nil {
			e.config.OnStepStart(step)
		}

		output, err := e.executeStep(ctx, step, env, resolver)

		if err == nil {
			stepResult.Complete(StatusSuccess, output, nil)
			if e.config.OnStepComplete != nil {
				e.config.OnStepComplete(step, stepResult)
			}
			return stepResult, nil
		}

		lastErr = err
		e.logger.Warn("step failed",
			"step", step.GetID(),
			"attempt", attempt,
			"maxAttempts", maxAttempts,
			"error", err)

		if attempt < maxAttempts {
			// Apply backoff
			backoffDelay := delay
			if step.Retry != nil && step.Retry.BackoffMultiplier > 0 {
				for i := 1; i < attempt; i++ {
					backoffDelay = time.Duration(float64(backoffDelay) * step.Retry.BackoffMultiplier)
				}
			}

			select {
			case <-ctx.Done():
				stepResult.Complete(StatusFailure, nil, ctx.Err())
				return stepResult, ctx.Err()
			case <-time.After(backoffDelay):
			}
		}
	}

	stepResult.Complete(StatusFailure, nil, lastErr)
	if e.config.OnStepComplete != nil {
		e.config.OnStepComplete(step, stepResult)
	}
	return stepResult, lastErr
}

// executeStep executes a single step.
func (e *Executor) executeStep(ctx context.Context, step *Step, env *activity.Environment, resolver *Resolver) (any, error) {
	// Get the activity
	act, ok := e.registry.Get(step.Activity)
	if !ok {
		return nil, fmt.Errorf("unknown activity: %s", step.Activity)
	}

	// Resolve parameters
	params := make(map[string]any)
	if step.Params != nil {
		resolved, err := resolver.ResolveMap(step.Params)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve params: %w", err)
		}
		params = resolved
	}

	// Apply timeout
	timeout := step.GetTimeout(Duration(e.config.DefaultTimeout)).Duration()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	e.logger.Info("executing step", "step", step.GetID(), "activity", step.Activity)

	// Execute the activity
	output, err := act.Execute(ctx, params, env)
	if err != nil {
		return nil, fmt.Errorf("activity %s failed: %w", step.Activity, err)
	}

	return output, nil
}

// handleError handles workflow error.
func (e *Executor) handleError(ctx context.Context, handler *ErrorHandler, env *activity.Environment, resolver *Resolver, result *WorkflowResult, originalErr error) {
	// Take screenshot if configured
	if handler.Screenshot && env.Vibe != nil {
		data, err := env.Vibe.Screenshot(ctx)
		if err == nil {
			result.AddScreenshot(Screenshot{
				Timestamp: time.Now(),
				Data:      base64.StdEncoding.EncodeToString(data),
				Reason:    "error: " + originalErr.Error(),
			})
		}
	}

	// Execute error handling steps
	if len(handler.Steps) > 0 {
		e.runSteps(ctx, handler.Steps, env, resolver, result)
	}
}

// ValidationError represents a validation error.
type ValidationError struct {
	StepID  string
	Field   string
	Message string
}

// Validate checks a workflow for errors without executing.
func (e *Executor) Validate(ctx context.Context, wf *Workflow) []ValidationError {
	var errors []ValidationError

	if wf.Name == "" {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "workflow name is required",
		})
	}

	if len(wf.Steps) == 0 {
		errors = append(errors, ValidationError{
			Field:   "steps",
			Message: "workflow must have at least one step",
		})
	}

	for i := range wf.Steps {
		stepErrors := e.validateStep(&wf.Steps[i])
		errors = append(errors, stepErrors...)
	}

	return errors
}

// validateStep validates a single step.
func (e *Executor) validateStep(step *Step) []ValidationError {
	var errors []ValidationError

	if step.Activity == "" {
		errors = append(errors, ValidationError{
			StepID:  step.GetID(),
			Field:   "activity",
			Message: "activity is required",
		})
	} else if _, ok := e.registry.Get(step.Activity); !ok {
		errors = append(errors, ValidationError{
			StepID:  step.GetID(),
			Field:   "activity",
			Message: fmt.Sprintf("unknown activity: %s", step.Activity),
		})
	}

	// Validate nested steps
	if step.ForEach != nil {
		for i := range step.ForEach.Steps {
			stepErrors := e.validateStep(&step.ForEach.Steps[i])
			errors = append(errors, stepErrors...)
		}
	}
	for i := range step.Steps {
		stepErrors := e.validateStep(&step.Steps[i])
		errors = append(errors, stepErrors...)
	}

	return errors
}
