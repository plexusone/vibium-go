package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var launchHeadless bool

var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch a browser instance",
	Long: `Launch a new browser instance for automation.

The browser will stay open until you run 'vibium quit' or press Ctrl+C.

Examples:
  vibium launch              # Launch visible browser
  vibium launch --headless   # Launch headless browser`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Handle interrupt
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

		vibe, err := launchBrowser(ctx, launchHeadless)
		if err != nil {
			return err
		}

		mode := "visible"
		if launchHeadless {
			mode = "headless"
		}
		fmt.Printf("Browser launched (%s mode)\n", mode)
		fmt.Println("Press Ctrl+C to quit or use 'vibium quit'")

		// Wait for interrupt
		<-sigCh
		fmt.Println("\nShutting down...")

		if err := vibe.Quit(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		}
		if err := clearSession(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(launchCmd)
	launchCmd.Flags().BoolVar(&launchHeadless, "headless", false, "Run browser in headless mode")
}
