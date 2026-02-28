//nolint:dupl // type and fill are separate commands with different semantics
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var typeTimeout time.Duration

var typeCmd = &cobra.Command{
	Use:   "type <selector> <text>",
	Short: "Type text into an element",
	Long: `Type text into an input element (appends to existing content).

Use 'fill' command if you want to clear existing content first.

Examples:
  vibium type "#search" "hello world"
  vibium type "input[name='query']" "search term"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]
		text := args[1]

		ctx, cancel := context.WithTimeout(context.Background(), typeTimeout)
		defer cancel()

		vibe := mustGetVibe(ctx)

		el, err := vibe.Find(ctx, selector, nil)
		if err != nil {
			return fmt.Errorf("element not found: %w", err)
		}

		if err := el.Type(ctx, text, nil); err != nil {
			return fmt.Errorf("type failed: %w", err)
		}

		fmt.Printf("Typed into: %s\n", selector)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(typeCmd)
	typeCmd.Flags().DurationVar(&typeTimeout, "timeout", 10*time.Second, "Timeout for finding element and typing")
}
