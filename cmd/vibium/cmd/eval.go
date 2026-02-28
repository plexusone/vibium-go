package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var evalTimeout time.Duration

var evalCmd = &cobra.Command{
	Use:   "eval <javascript>",
	Short: "Execute JavaScript",
	Long: `Execute JavaScript on the page and print the result.

Examples:
  vibium eval "document.title"
  vibium eval "document.querySelectorAll('a').length"
  vibium eval "window.location.href"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		script := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), evalTimeout)
		defer cancel()

		vibe := mustGetVibe(ctx)

		result, err := vibe.Evaluate(ctx, script)
		if err != nil {
			return fmt.Errorf("eval failed: %w", err)
		}

		// Pretty print result
		if result == nil {
			fmt.Println("undefined")
		} else if s, ok := result.(string); ok {
			fmt.Println(s)
		} else {
			jsonBytes, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				fmt.Printf("%v\n", result)
			} else {
				fmt.Println(string(jsonBytes))
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(evalCmd)
	evalCmd.Flags().DurationVar(&evalTimeout, "timeout", 30*time.Second, "Evaluation timeout")
}
