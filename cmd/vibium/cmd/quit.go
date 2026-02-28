package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var quitTimeout time.Duration

var quitCmd = &cobra.Command{
	Use:   "quit",
	Short: "Quit the browser",
	Long: `Close the browser and cleanup resources.

Examples:
  vibium quit`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), quitTimeout)
		defer cancel()

		if err := quitBrowser(ctx); err != nil {
			return err
		}

		fmt.Println("Browser closed")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(quitCmd)
	quitCmd.Flags().DurationVar(&quitTimeout, "timeout", 10*time.Second, "Quit timeout")
}
