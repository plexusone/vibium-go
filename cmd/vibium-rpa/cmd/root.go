package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	headless bool
	verbose  bool
	workDir  string
	varFlags []string
)

var rootCmd = &cobra.Command{
	Use:   "vibium-rpa",
	Short: "RPA workflow automation tool",
	Long: `vibium-rpa is a Robotic Process Automation tool that executes
YAML/JSON workflow definitions using browser automation.

Examples:
  # Run a workflow
  vibium-rpa run workflow.yaml

  # Run in headless mode
  vibium-rpa run workflow.yaml --headless

  # Run with variables
  vibium-rpa run workflow.yaml --var username=admin --var password=secret

  # Validate a workflow
  vibium-rpa validate workflow.yaml

  # List available activities
  vibium-rpa list activities
`,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&headless, "headless", false, "Run browser in headless mode")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().StringVar(&workDir, "workdir", "", "Working directory (default: current directory)")
	rootCmd.PersistentFlags().StringArrayVar(&varFlags, "var", nil, "Set variable (key=value)")

	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "help",
		Short:  "Help about any command",
		Hidden: true,
	})

	// Set custom usage template
	rootCmd.SetUsageTemplate(`Usage:
  {{.CommandPath}} [command]

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}

Use "{{.CommandPath}} [command] --help" for more information about a command.
`)
}

// parseVariables parses --var flags into a map.
func parseVariables() map[string]string {
	vars := make(map[string]string)
	for _, v := range varFlags {
		for i := 0; i < len(v); i++ {
			if v[i] == '=' {
				vars[v[:i]] = v[i+1:]
				break
			}
		}
	}
	return vars
}

// getWorkDir returns the working directory.
func getWorkDir() string {
	if workDir != "" {
		return workDir
	}
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to get working directory: %v\n", err)
		return "."
	}
	return wd
}
