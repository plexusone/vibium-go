package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	screenshotTimeout  time.Duration
	screenshotSelector string
)

var screenshotCmd = &cobra.Command{
	Use:   "screenshot <filename>",
	Short: "Take a screenshot",
	Long: `Capture a screenshot of the current page or a specific element.

Examples:
  vibium screenshot page.png
  vibium screenshot element.png --selector "#main"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), screenshotTimeout)
		defer cancel()

		vibe := mustGetVibe(ctx)

		var data []byte
		var err error

		if screenshotSelector != "" {
			// Element screenshot
			el, findErr := vibe.Find(ctx, screenshotSelector, nil)
			if findErr != nil {
				return fmt.Errorf("element not found: %w", findErr)
			}
			data, err = el.Screenshot(ctx)
		} else {
			// Page screenshot
			data, err = vibe.Screenshot(ctx)
		}

		if err != nil {
			return fmt.Errorf("screenshot failed: %w", err)
		}

		if err := os.WriteFile(filename, data, 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		fmt.Printf("Screenshot saved: %s (%d bytes)\n", filename, len(data))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(screenshotCmd)
	screenshotCmd.Flags().DurationVar(&screenshotTimeout, "timeout", 30*time.Second, "Screenshot timeout")
	screenshotCmd.Flags().StringVar(&screenshotSelector, "selector", "", "Capture specific element instead of page")
}
