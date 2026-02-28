// Command vibium provides a CLI for browser automation.
package main

import (
	"os"

	"github.com/agentplexus/vibium-go/cmd/vibium/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
