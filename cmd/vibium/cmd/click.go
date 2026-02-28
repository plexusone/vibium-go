package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var clickTimeout time.Duration

var clickCmd = &cobra.Command{
	Use:   "click <selector>",
	Short: "Click an element",
	Long: `Click an element identified by CSS selector.

Examples:
  vibium click "#submit"
  vibium click "button.login"
  vibium click "[data-testid='submit-btn']"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), clickTimeout)
		defer cancel()

		vibe := mustGetVibe(ctx)

		el, err := vibe.Find(ctx, selector, nil)
		if err != nil {
			return fmt.Errorf("element not found: %w", err)
		}

		if err := el.Click(ctx, nil); err != nil {
			return fmt.Errorf("click failed: %w", err)
		}

		fmt.Printf("Clicked: %s\n", selector)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(clickCmd)
	clickCmd.Flags().DurationVar(&clickTimeout, "timeout", 10*time.Second, "Timeout for finding and clicking element")
}
