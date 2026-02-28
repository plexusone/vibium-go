package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	sessionFile string
	verbose     bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "vibium",
	Short: "Browser automation CLI",
	Long: `Vibium is a browser automation tool that provides:

  - MCP server for AI-assisted browser automation
  - CLI commands for scripted browser control
  - YAML script runner for batch operations

Examples:
  # Start MCP server
  vibium mcp --headless

  # Launch browser and run commands
  vibium launch --headless
  vibium go https://example.com
  vibium click "#submit"
  vibium screenshot output.png
  vibium quit

  # Run a script file
  vibium run test.yaml`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&sessionFile, "session", "", "Session file path (default: ~/.vibium/session.json)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}

// getSessionPath returns the session file path
func getSessionPath() string {
	if sessionFile != "" {
		return sessionFile
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ".vibium-session.json"
	}
	return fmt.Sprintf("%s/.vibium/session.json", home)
}
