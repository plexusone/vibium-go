package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// Version information (set at build time)
var (
	Version   = "0.1.0"
	GitCommit = "dev"
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  "Print the version, build date, and git commit of vibium-rpa.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("vibium-rpa version %s\n", Version)
		fmt.Printf("  Git commit: %s\n", GitCommit)
		fmt.Printf("  Build date: %s\n", BuildDate)
		fmt.Printf("  Go version: %s\n", runtime.Version())
		fmt.Printf("  OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
