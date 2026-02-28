package cmd

import (
	"context"
	"fmt"

	"github.com/plexusone/vibium-go/rpa"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate <workflow-file>",
	Short: "Validate a workflow without executing",
	Long: `Validate an RPA workflow file for syntax and semantic errors.

This command parses the workflow file, validates its structure,
and checks that all referenced activities exist.

Examples:
  # Validate a workflow
  vibium-rpa validate workflow.yaml

  # Validate multiple workflows
  vibium-rpa validate workflow1.yaml workflow2.yaml
`,
	Args: cobra.MinimumNArgs(1),
	RunE: validateWorkflow,
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

func validateWorkflow(cmd *cobra.Command, args []string) error {
	hasErrors := false

	for _, path := range args {
		fmt.Printf("Validating: %s\n", path)

		// Parse the workflow
		wf, err := rpa.ParseFile(path)
		if err != nil {
			fmt.Printf("  ✗ Parse error: %v\n", err)
			hasErrors = true
			continue
		}

		// Validate
		executor := rpa.NewExecutor(rpa.ExecutorConfig{})
		errors := executor.Validate(context.Background(), wf)

		if len(errors) == 0 {
			fmt.Printf("  ✓ Valid workflow: %s\n", wf.Name)
			fmt.Printf("    Steps: %d\n", len(wf.Steps))
			if len(wf.Variables) > 0 {
				fmt.Printf("    Variables: %d\n", len(wf.Variables))
			}
		} else {
			fmt.Printf("  ✗ Validation errors:\n")
			for _, e := range errors {
				if e.StepID != "" {
					fmt.Printf("    - Step %s, %s: %s\n", e.StepID, e.Field, e.Message)
				} else {
					fmt.Printf("    - %s: %s\n", e.Field, e.Message)
				}
			}
			hasErrors = true
		}
	}

	if hasErrors {
		return fmt.Errorf("validation failed")
	}

	return nil
}
