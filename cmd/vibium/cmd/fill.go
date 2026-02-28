//nolint:dupl // fill and type are separate commands with different semantics
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var fillTimeout time.Duration

var fillCmd = &cobra.Command{
	Use:   "fill <selector> <text>",
	Short: "Fill an input element",
	Long: `Clear an input element and fill it with text (replaces existing content).

Use 'type' command if you want to append to existing content.

Examples:
  vibium fill "#email" "user@example.com"
  vibium fill "input[name='password']" "secret123"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]
		text := args[1]

		ctx, cancel := context.WithTimeout(context.Background(), fillTimeout)
		defer cancel()

		vibe := mustGetVibe(ctx)

		el, err := vibe.Find(ctx, selector, nil)
		if err != nil {
			return fmt.Errorf("element not found: %w", err)
		}

		if err := el.Fill(ctx, text, nil); err != nil {
			return fmt.Errorf("fill failed: %w", err)
		}

		fmt.Printf("Filled: %s\n", selector)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(fillCmd)
	fillCmd.Flags().DurationVar(&fillTimeout, "timeout", 10*time.Second, "Timeout for finding element and filling")
}
