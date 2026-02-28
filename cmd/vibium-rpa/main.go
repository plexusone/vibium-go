package main

import (
	"os"

	"github.com/agentplexus/vibium-go/cmd/vibium-rpa/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
