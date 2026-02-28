package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/agentplexus/vibium-go/rpa"
	"github.com/spf13/cobra"
)

var (
	outputFile   string
	outputFormat string
	dryRun       bool
)

var runCmd = &cobra.Command{
	Use:   "run <workflow-file>",
	Short: "Execute a workflow",
	Long: `Execute an RPA workflow from a YAML or JSON file.

The workflow file contains steps that will be executed in order,
with browser automation, file operations, and HTTP requests.

Examples:
  # Run a workflow
  vibium-rpa run workflow.yaml

  # Run in headless mode with variables
  vibium-rpa run workflow.yaml --headless --var username=admin

  # Run and save results to JSON
  vibium-rpa run workflow.yaml --output results.json

  # Dry run (validate without executing)
  vibium-rpa run workflow.yaml --dry-run
`,
	Args: cobra.ExactArgs(1),
	RunE: runWorkflow,
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Save results to file (format from extension)")
	runCmd.Flags().StringVar(&outputFormat, "format", "", "Output format: json, markdown, html, junit")
	runCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Validate workflow without executing")
}

func runWorkflow(cmd *cobra.Command, args []string) error {
	workflowPath := args[0]

	// Set up logging
	logLevel := slog.LevelInfo
	if verbose {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel}))

	// Create executor config
	config := rpa.ExecutorConfig{
		Headless:       headless,
		WorkDir:        getWorkDir(),
		Variables:      parseVariables(),
		DryRun:         dryRun,
		Logger:         logger,
		OnStepStart:    onStepStart,
		OnStepComplete: onStepComplete,
	}

	executor := rpa.NewExecutor(config)

	// Set up context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Fprintln(os.Stderr, "\nReceived interrupt, cancelling...")
		cancel()
	}()

	// Run the workflow
	fmt.Printf("Running workflow: %s\n", workflowPath)
	if dryRun {
		fmt.Println("(dry-run mode - validating only)")
	}

	result, err := executor.RunFile(ctx, workflowPath)
	if err != nil {
		return fmt.Errorf("workflow execution failed: %w", err)
	}

	// Output results
	if outputFile != "" {
		if err := writeOutput(result, outputFile, outputFormat); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		fmt.Printf("Results saved to: %s\n", outputFile)
	}

	// Print summary
	printSummary(result)

	if result.Status != rpa.StatusSuccess {
		return fmt.Errorf("workflow failed: %s", result.Error)
	}

	return nil
}

func onStepStart(step *rpa.Step) {
	if verbose {
		fmt.Printf("  → Starting: %s\n", step.GetID())
	}
}

func onStepComplete(step *rpa.Step, result *rpa.StepResult) {
	status := "✓"
	if result.Status == rpa.StatusFailure {
		status = "✗"
	} else if result.Status == rpa.StatusSkipped {
		status = "○"
	}

	if verbose || result.Status != rpa.StatusSuccess {
		fmt.Printf("  %s %s (%s)\n", status, step.GetID(), result.Duration.Round(1000000))
		if result.Error != "" {
			fmt.Printf("    Error: %s\n", result.Error)
		}
	}
}

func printSummary(result *rpa.WorkflowResult) {
	fmt.Println()
	fmt.Printf("Workflow: %s\n", result.WorkflowName)
	fmt.Printf("Status:   %s\n", result.Status)
	fmt.Printf("Duration: %s\n", result.Duration.Round(1000000))
	fmt.Printf("Steps:    %d total, %d success, %d failed, %d skipped\n",
		result.TotalSteps(),
		result.SuccessCount(),
		result.FailureCount(),
		result.SkippedCount())

	if result.Error != "" {
		fmt.Printf("Error:    %s\n", result.Error)
	}
}

func writeOutput(result *rpa.WorkflowResult, path, format string) error {
	// Determine format from extension if not specified
	if format == "" {
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".json":
			format = "json"
		case ".md":
			format = "markdown"
		case ".html":
			format = "html"
		case ".xml":
			format = "junit"
		default:
			format = "json"
		}
	}

	var data []byte
	var err error

	switch format {
	case "json":
		data, err = json.MarshalIndent(result, "", "  ")
	case "markdown":
		data = formatMarkdown(result)
	case "html":
		data = formatHTML(result)
	case "junit":
		data = formatJUnit(result)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func formatMarkdown(result *rpa.WorkflowResult) []byte {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# Workflow: %s\n\n", result.WorkflowName))
	sb.WriteString(fmt.Sprintf("**Status:** %s\n\n", result.Status))
	sb.WriteString(fmt.Sprintf("**Duration:** %s\n\n", result.Duration.Round(1000000)))

	sb.WriteString("## Summary\n\n")
	sb.WriteString(fmt.Sprintf("- Total steps: %d\n", result.TotalSteps()))
	sb.WriteString(fmt.Sprintf("- Successful: %d\n", result.SuccessCount()))
	sb.WriteString(fmt.Sprintf("- Failed: %d\n", result.FailureCount()))
	sb.WriteString(fmt.Sprintf("- Skipped: %d\n\n", result.SkippedCount()))

	if result.Error != "" {
		sb.WriteString("## Error\n\n")
		sb.WriteString(fmt.Sprintf("```\n%s\n```\n\n", result.Error))
	}

	sb.WriteString("## Steps\n\n")
	sb.WriteString("| Step | Activity | Status | Duration |\n")
	sb.WriteString("|------|----------|--------|----------|\n")

	for _, step := range result.Steps {
		status := "✓"
		if step.Status == rpa.StatusFailure {
			status = "✗"
		} else if step.Status == rpa.StatusSkipped {
			status = "○"
		}
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
			step.StepID, step.Activity, status, step.Duration.Round(1000000)))
	}

	return []byte(sb.String())
}

func formatHTML(result *rpa.WorkflowResult) []byte {
	var sb strings.Builder

	sb.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
	sb.WriteString("<title>Workflow Results: " + result.WorkflowName + "</title>\n")
	sb.WriteString("<style>\n")
	sb.WriteString("body { font-family: sans-serif; margin: 20px; }\n")
	sb.WriteString("table { border-collapse: collapse; width: 100%; }\n")
	sb.WriteString("th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }\n")
	sb.WriteString("th { background-color: #f2f2f2; }\n")
	sb.WriteString(".success { color: green; }\n")
	sb.WriteString(".failure { color: red; }\n")
	sb.WriteString(".skipped { color: gray; }\n")
	sb.WriteString("</style>\n</head>\n<body>\n")

	sb.WriteString(fmt.Sprintf("<h1>Workflow: %s</h1>\n", result.WorkflowName))

	statusClass := "success"
	if result.Status == rpa.StatusFailure {
		statusClass = "failure"
	}
	sb.WriteString(fmt.Sprintf("<p><strong>Status:</strong> <span class=\"%s\">%s</span></p>\n", statusClass, result.Status))
	sb.WriteString(fmt.Sprintf("<p><strong>Duration:</strong> %s</p>\n", result.Duration.Round(1000000)))

	sb.WriteString("<h2>Steps</h2>\n")
	sb.WriteString("<table>\n<tr><th>Step</th><th>Activity</th><th>Status</th><th>Duration</th><th>Error</th></tr>\n")

	for _, step := range result.Steps {
		statusClass := "success"
		if step.Status == rpa.StatusFailure {
			statusClass = "failure"
		} else if step.Status == rpa.StatusSkipped {
			statusClass = "skipped"
		}
		sb.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td class=\"%s\">%s</td><td>%s</td><td>%s</td></tr>\n",
			step.StepID, step.Activity, statusClass, step.Status, step.Duration.Round(1000000), step.Error))
	}

	sb.WriteString("</table>\n</body>\n</html>")

	return []byte(sb.String())
}

func formatJUnit(result *rpa.WorkflowResult) []byte {
	var sb strings.Builder

	failures := result.FailureCount()
	skipped := result.SkippedCount()

	sb.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	sb.WriteString(fmt.Sprintf("<testsuite name=\"%s\" tests=\"%d\" failures=\"%d\" skipped=\"%d\" time=\"%.3f\">\n",
		result.WorkflowName, result.TotalSteps(), failures, skipped, result.Duration.Seconds()))

	for _, step := range result.Steps {
		sb.WriteString(fmt.Sprintf("  <testcase name=\"%s\" classname=\"%s\" time=\"%.3f\">\n",
			step.StepID, step.Activity, step.Duration.Seconds()))

		if step.Status == rpa.StatusFailure {
			sb.WriteString(fmt.Sprintf("    <failure message=\"%s\"/>\n", step.Error))
		} else if step.Status == rpa.StatusSkipped {
			sb.WriteString("    <skipped/>\n")
		}

		sb.WriteString("  </testcase>\n")
	}

	sb.WriteString("</testsuite>\n")

	return []byte(sb.String())
}
